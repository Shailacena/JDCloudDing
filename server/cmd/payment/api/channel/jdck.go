package channel

import (
	"apollo/server/cmd/payment/api/common"
	"apollo/server/cmd/payment/pkg/errorx"
	"apollo/server/cmd/payment/pkg/types"
	"apollo/server/internal/model"
	"apollo/server/pkg/config"
	"apollo/server/pkg/data"
	"apollo/server/pkg/util"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"net"
	"time"
)

type JDCKPay struct {
	MerchantId int32
	Merchant   *model.Merchant
	db         *gorm.DB
	CD         int
	WxPayUrl   string
}

func NewJDCKPay(merchantId int32) (*JDCKPay, error) {
	db := data.Instance()
	merchant, err := common.FindMerchant(db, merchantId)
	if err != nil {
		return nil, err
	}

	cd := 10
	if !config.IsProd() {
		cd = 2
	}

	return &JDCKPay{
		db:         db,
		MerchantId: merchantId,
		Merchant:   merchant,
		CD:         cd,
	}, nil
}

func (jd *JDCKPay) GetMerchant(c echo.Context) *model.Merchant {
	return jd.Merchant
}

func (jd *JDCKPay) FindGoods(c echo.Context, channelId string, merchantId int32, amount float64) (*common.FindGoodsResult, error) {
	db := jd.db

	rows, err := common.FindGoods(c, db, channelId, merchantId, amount)
	if err != nil {
		c.Logger().Errorf("common.FindGoods error=%s", err)
		return nil, errorx.ErrGoodsNotFound
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var res common.FindGoodsResult
		_ = db.ScanRows(rows, &res)

		if len(res.SkuId) == 0 {
			continue
		}

		err = db.Model(&model.Goods{}).Where("id = ?", res.Id).Update("weight", gorm.Expr("weight + ?", 1)).Error
		if err != nil {
			c.Logger().Error("update goods weight error", err.Error())
		}
		return &res, nil
	}

	return nil, errorx.ErrGoodsNotFound
}

func (jd *JDCKPay) GenOrder(c echo.Context, req types.CreateOrderReq, goods *common.FindGoodsResult, orderId string) (model.Order, error) {
	skuId := goods.SkuId
	partnerId := goods.PartnerId
	shop := goods.ShopName
	payType := goods.PayType
	merchant := jd.Merchant

	partnerOrderId, wxPayUrl, err := jd.checkout(c, skuId)
	if err != nil {
		c.Logger().Error("checkout skuId error", err.Error())
		return model.Order{}, err
	}
	jd.WxPayUrl = wxPayUrl

	o := model.Order{
		OrderId:         orderId,
		ChannelId:       model.ChannelId(req.ChannelId),
		MerchantId:      uint(req.MerchantId),
		MerchantName:    merchant.Nickname,
		MerchantOrderId: req.MerchantTradeNo,
		Amount:          util.ToDecimal(req.Amount),
		SkuId:           skuId,
		IP:              c.RealIP(),
		NotifyUrl:       req.NotifyUrl,
		Shop:            shop,
		PartnerId:       partnerId,
		PartnerName:     goods.PartnerNickname,
		PayType:         payType,
		DarkNumber:      partnerOrderId,
	}

	return o, nil
}

func jdCkHttp[T any](c echo.Context, url string, params map[string]string) (T, error) {
	var result T

	client := resty.New()
	client = client.SetTimeout(10 * time.Second)
	resp, err := client.R().SetQueryParams(params).Get(url)
	if err != nil {
		return result, err
	}

	body := resp.Body()
	c.Logger().Infof("%T, jdCkHttp, params=%+v, body=%s", result, params, string(body))
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (jd *JDCKPay) getProxyIp(c echo.Context) (string, error) {
	url := "http://proxy.siyetian.com/apis_get2.html?token=0QjM2YTM2YjN2ETLNpWQ59ERJdXTqV1dORUS00keRl3TB1STqFUeORVQw0kanpnTqFENPR0a14kejp3TUlke.wN0cTMygTN0cTM&limit=1&type=0&time=&split=0&split_text="
	client := resty.New()
	client = client.SetTimeout(10 * time.Second)
	resp, err := client.R().Get(url)
	if err != nil {
		return "", err
	}
	ip := fmt.Sprintf("16666166244:pow123456@%s", string(resp.Body()))
	c.Logger().Info("getProxyIp resp:", ip)
	return fmt.Sprintf("16666166244:pow123456@%s", ip), nil
}

func (jd *JDCKPay) checkout(c echo.Context, skuId string) (string, string, error) {
	url := "http://127.0.0.1:8008/api"

	proxyIp, err := jd.getProxyIp(c)
	if err != nil {
		c.Logger().Error("getProxyIp error", err)
	}

	accounts, err := jd.findCk(jd.Merchant.MasterId, 20)
	if err != nil {
		return "", "", err
	}

	c.Logger().Info("accounts=", accounts)

	for i := 0; i < len(accounts); i++ {
		account := accounts[i]
		ckId := account.id
		ck := account.ck
		if len(proxyIp) == 0 {
			proxyIp, err = jd.getProxyIp(c)
			if err != nil {
				c.Logger().Error("getProxyIp error", err)
			}
		}

		params := map[string]string{
			"method": "is_cookie_valid",
			"cookie": ck,
		}
		validCookieResp, err := jdCkHttp[ValidCookieResp](c, url, params)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				proxyIp = ""
				err := jd.IncrWeight(ckId)
				if err != nil {
					c.Logger().Error("IncrWeight failed", ckId, err)
				}
			}
			c.Logger().Error("ValidCookie failed", err)
			continue
		}
		if validCookieResp.Code != 200 {
			c.Logger().Error("ValidCookie failed, ck不可用", ck, validCookieResp)
			err := common.UpdateCKStatus(jd.db, ckId, map[string]any{
				"status": model.JDAccountStatusInvalid,
				"remark": validCookieResp.Message,
			})
			if err != nil {
				c.Logger().Error("UpdateCKStatus failed, ck不可用", ckId, ck, err)
			}
			continue
		}

		once := true
	submitAgain:
		params = map[string]string{
			"method":      "app_Submitorder",
			"skuid":       skuId,
			"num":         "1",
			"skuSource":   "0",
			"paymentType": "4",
			"ip":          proxyIp,
			"cookie":      ck,
		}
		submitOrderResp, err := jdCkHttp[SubmitOrderResp](c, url, params)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				proxyIp = ""
				err := jd.IncrWeight(ckId)
				if err != nil {
					c.Logger().Error("IncrWeight failed", ckId, err)
				}
			}
			c.Logger().Error("SubmitOrder failed", err)
			continue
		}
		if submitOrderResp.Code != "0" {
			c.Logger().Error("SubmitOrder failed", submitOrderResp)

			switch submitOrderResp.Code {
			case "400":
				if once {
					once = false
					address, err := jd.findAddress()
					if err != nil {
						continue
					}
					params = map[string]string{
						"method":  "app_addAddress",
						"text":    address,
						"default": "1",
						"ip":      proxyIp,
						"cookie":  ck,
					}
					addAddressResp, err := jdCkHttp[AddAddressResp](c, url, params)
					if err != nil {
						c.Logger().Error("AddAddress failed", err)
					}
					if err == nil && addAddressResp.Code != "0" {
						c.Logger().Error("AddAddress failed resp", addAddressResp)
					}
					if err != nil || addAddressResp.Code != "0" {
						err1 := common.UpdateCKStatus(jd.db, ckId, map[string]any{
							"status": model.JDAccountStatusAddAddressErr,
							"remark": "添加地址失败",
						})
						if err1 != nil {
							c.Logger().Error("UpdateCKStatus failed, 添加地址失败", ckId, ck, err1)
						}
						continue
					}

					goto submitAgain
				}

			case "601":
				err1 := common.UpdateCKStatus(jd.db, ckId, map[string]any{
					"status": model.JDAccountStatusHot,
					"remark": submitOrderResp.Message,
					"weight": gorm.Expr("weight + ?", 1),
				})
				if err1 != nil {
					c.Logger().Error("UpdateCKStatus failed, 火热", ckId, ck, err1)
				}
				continue
			default:
				err1 := common.UpdateCKStatus(jd.db, ckId, map[string]any{
					"status": model.JDAccountStatusSubmitOrderErr,
					"remark": submitOrderResp.Message,
				})
				if err1 != nil {
					c.Logger().Error("UpdateCKStatus failed, 下单失败", ckId, ck, err1)
				}
				continue
			}
		}
		orderId := cast.ToString(submitOrderResp.SubmitOrder.OrderID)

		params = map[string]string{
			"method":  "app_order_wx_pay",
			"orderId": orderId,
			"ip":      proxyIp,
			"cookie":  ck,
		}

		appOrderWxPayResp, err := jdCkHttp[AppOrderWxPayResp](c, url, params)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				proxyIp = ""
				err := jd.IncrWeight(ckId)
				if err != nil {
					c.Logger().Error("IncrWeight failed", ckId, err)
				}
			}
			c.Logger().Error("AppOrderWxPay failed", err)
			continue
		}
		if appOrderWxPayResp.Code != 200 {
			c.Logger().Error("AppOrderWxPay failed", appOrderWxPayResp)
			m := map[string]any{
				"status": model.JDAccountStatusGetWxPayErr,
				"remark": appOrderWxPayResp.Message,
			}
			if appOrderWxPayResp.Code == 400 {
				m = map[string]any{
					"status": model.JDAccountStatusInvalid,
					"remark": appOrderWxPayResp.Message,
				}
			}
			err1 := common.UpdateCKStatus(jd.db, ckId, m)
			if err1 != nil {
				c.Logger().Error("AppOrderWxPayResp UpdateCKStatus failed, 提码失败", ckId, ck, err1)
			}
			continue
		}

		err = common.UpdateCKUseCount(jd.db, ckId)
		if err != nil {
			c.Logger().Error("UpdateCKUseCount failed", ckId, err)
		}

		return orderId, appOrderWxPayResp.Data, nil
	}

	return "", "", errors.New("没有可用的ck")
}

type ckObj struct {
	id uint
	ck string
}

func (jd *JDCKPay) findCk(masterId uint, limit int) ([]ckObj, error) {
	jdAccounts, err := common.FindCK(jd.db, masterId, limit)
	if err != nil {
		return nil, err
	}

	accounts := lo.Map(jdAccounts, func(account *model.JDAccount, index int) ckObj {
		return ckObj{
			id: account.ID,
			ck: fmt.Sprintf("pin=%s;wskey=%s;", account.Account, account.WsKey),
		}
	})
	return accounts, nil

	// return "pin=%E5%88%98%E5%BF%83%E5%A8%B4;wskey=AAJoDP0ZAEBLeCU7d_tMomSLN8TKDknHWqQfPPlHkUREfwH65zyqJqbkYyHTV2MP6dkVZHpuuEcaZ65mdHhQP-ls59n-yLhx;", nil
}

func (jd *JDCKPay) findAddress() (string, error) {
	realNameAccount, err := common.FindAddress(jd.db, jd.Merchant.MasterId)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s,%s,%s", realNameAccount.Address, realNameAccount.Name, realNameAccount.Mobile), nil
}

func (jd *JDCKPay) IncrWeight(id uint) error {
	return jd.db.Model(&model.JDAccount{}).Where("id = ?", id).Update("weight", gorm.Expr("weight + ?", 1)).Error
}

func (jd *JDCKPay) GenPayUrl(c echo.Context, baseUrl string, o model.Order, numId string) string {
	return fmt.Sprintf("%s?merchantOrderId=%s&orderId=%s&price=%f&sku=%s&time=%d&ts=%d&wxPayUrl=%s", baseUrl, o.MerchantOrderId, o.OrderId, o.Amount, o.SkuId, jd.CD*60, time.Now().Unix(), jd.WxPayUrl)
}
