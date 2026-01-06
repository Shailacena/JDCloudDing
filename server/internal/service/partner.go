package service

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/internal/repository"
	"apollo/server/internal/service/partnerx"
	"apollo/server/pkg/data"
	"apollo/server/pkg/headerx"
	"apollo/server/pkg/timex"
	"apollo/server/pkg/totpx"
	"apollo/server/pkg/util"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

var (
	Partner = new(PartnerService)
)

type PartnerService struct {
}

type SyncGoodReqData struct {
	Fields    string `url:"fields"`
	Timestamp int64  `url:"timestamp"`
	Sign      string `url:"sign"`
	PageNo    int    `url:"pageno"`
}

type GoodsData struct {
	ApproveStatus string `json:"ApproveStatus"`
	NumIid        int    `json:"NumIid"`
	SkuId         int    `json:"SkuId"`
	Nick          string `json:"Nick"`
	Num           int    `json:"Num"`
	Price         string `json:"Price"`
	Title         string `json:"Title"`
}

type GoodsSubSKUData struct {
	Quantity int `json:"quantity"`
	SkuId    int `json:"skuId"`
}

type SyncGoodRspData struct {
	IsSuccess  bool        `json:"IsSuccess"`
	Data       []GoodsData `json:"Data"`
	Error_Code int         `json:"Error_Code"`
	Error_Msg  string      `json:"Error_Msg"`
}

type SyncAnssyGoodRspData struct {
	Code int         `json:"code"`
	Data []GoodsData `json:"data"`
}

type SyncAnssyGoodSubSKURspData struct {
	Code int               `json:"code"`
	Data []GoodsSubSKUData `json:"data"`
}

func (s *PartnerService) Register(c echo.Context, req *v1.PartnerRegisterReq) (*v1.PartnerRegisterResp, error) {
	var partnerGenerator partnerx.IPartnerGenerator

	switch req.Type {
	case model.PartnerTypeAgiso:
		partnerGenerator = partnerx.NewAgiso(req.Type, req)
	case model.PartnerTypeAnssy:
		partnerGenerator = partnerx.NewAnssy(req.Type, req)
	default:
		return nil, errors.New("合作商类型不存在")
	}

	err := partnerGenerator.CheckParams(c)
	if err != nil {
		return nil, err
	}

	p := partnerGenerator.GenPartner()

	header := headerx.GetDataFromHeader(c)
	adminId := header.AdminId
	p.ParentId = adminId

	creator, err := repository.Admin.GetById(c, adminId)
	if err != nil {
		return nil, err
	}
	p.MasterId = creator.MasterId

	newPartner, err := repository.Partner.Register(c, &p)
	if err != nil {
		return nil, err
	}

	secret, urlPath, err := totpx.Generate(p.Username)
	if err != nil {
		return nil, err
	}

	newPartner.SecretKey = secret
	newPartner.UrlKey = urlPath

	_, err = repository.Partner.Update(c, newPartner.ID, &v1.PartnerUpdateReq{
		Secret:       secret,
		UrlPath:      urlPath,
		AqsToken:     req.AqsToken,
		AqsAppSecret: req.AqsAppSecret,
	})
	if err != nil {
		return nil, err
	}

	s.ResetGoodsWeight(c, newPartner.ID)

	return &v1.PartnerRegisterResp{
		Nickname: newPartner.Nickname,
		Password: newPartner.Password,
	}, nil
}

func (s *PartnerService) ResetGoodsWeight(c echo.Context, id uint) {
	err := data.Instance().Transaction(func(tx *gorm.DB) error {
		p, err := repository.Partner.FindPartner(c, id)
		if err != nil {
			return err
		}

		partners, err := repository.Partner.FindPartnerByChannelId(c, p.ChannelId)
		if err != nil {
			return err
		}

		ids := lo.Map(partners, func(p *model.Partner, _ int) uint {
			return p.ID
		})

		err = repository.Goods.ResetWeight(c, tx, ids)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		c.Logger().Errorf("Update Transaction error=%s", err)
	}
}

func (s *PartnerService) Login(c echo.Context, req *v1.PartnerLoginReq) (*v1.PartnerLoginResp, error) {
	p, err := repository.Partner.Login(c, req.Username, req.Password, req.VerifiCode)
	if err != nil {
		return nil, err
	}

	return &v1.PartnerLoginResp{
		Id:       p.ID,
		Token:    p.Token,
		Nickname: p.Nickname,
		Level:    p.Level,
	}, nil
}

func (s *PartnerService) Logout(c echo.Context, req *v1.PartnerLogoutReq, token string) (*v1.PartnerLogoutResp, error) {
	err := repository.Partner.Logout(c, token)
	if err != nil {
		return nil, err
	}

	return &v1.PartnerLogoutResp{}, nil
}

func (s *PartnerService) ResetPassword(c echo.Context, req *v1.PartnerResetPasswordReq) (*v1.PartnerResetPasswordResp, error) {
	user, err := repository.Partner.ResetPassword(c, req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.PartnerResetPasswordResp{
		Password: user.Password,
	}, nil
}

func (s *PartnerService) Delete(c echo.Context, req *v1.PartnerDeleteReq) (*v1.PartnerDeleteResp, error) {
	_, err := repository.Partner.Delete(c, req.Id)
	if err != nil {
		return nil, err
	}
	repository.Partner.DeleteAllGoods(c, req.Id)
	return &v1.PartnerDeleteResp{}, nil
}

func (s *PartnerService) List(c echo.Context, req *v1.ListPartnerReq) (*v1.ListPartnerResp, error) {
	var parentIds []uint
	if req.PartnerId == 0 {
		parentIds, _ = Admin.FindParentIds(c)
	}

	partners, total, err := repository.Partner.List(c, req, parentIds)
	if err != nil {
		return nil, err
	}

	ids := lo.Map(partners, func(item *model.Partner, _ int) uint {
		return item.ID
	})

	db := data.Instance()

	dataMap := make(map[uint]repository.QueryPartnerAmountResult)
	if !req.IgnoreStatistics {
		results, err := repository.Order.QueryResultByPartner(c, db, ids, timex.GetPRCNowTime().Carbon2Time())
		if err != nil {
			return nil, err
		}

		dataMap = lo.SliceToMap(results, func(item repository.QueryPartnerAmountResult) (uint, repository.QueryPartnerAmountResult) {
			return item.PartnerId, item
		})
	}

	list := make([]*v1.Partner, 0, len(partners))

	lastData := repository.Order.QueryLastOrderResult(c, db, ids)

	for _, p := range partners {
		var todayAmount, todayOrderNum, todaySuccessAmount, todaySuccessOrderNum float64
		item, ok := dataMap[p.ID]
		if ok {
			todayAmount = item.TodayOrderAmount
			todayOrderNum = item.TodayOrderNum
			todaySuccessAmount = item.TodaySuccessAmount
			todaySuccessOrderNum = item.TodaySuccessOrderNum
		}

		lastDataItem := lastData[p.ID]

		var anssyExpiredAt int64
		if p.AnssyExpiredAt != nil && !p.AnssyExpiredAt.IsZero() {
			anssyExpiredAt = p.AnssyExpiredAt.Unix()
		}

		list = append(list, &v1.Partner{
			Id:            p.ID,
			Nickname:      p.Nickname,
			ChannelId:     p.ChannelId,
			PayType:       p.PayType,
			Balance:       util.ToDecimal(p.Balance),
			Priority:      p.Priority,
			SuperiorAgent: p.SuperiorAgent,
			Level:         p.Level,
			StockAmount:   p.StockAmount,
			RechargeTime:  p.RechargeTime,
			PrivateKey:    p.PrivateKey,
			AqsAppSecret:  p.AqsAppSecret,
			AqsToken:      p.AqsToken,
			Enable:        int(p.Enable),
			Remark:        p.Remark,
			Type:          p.Type,
			UrlKey:        p.UrlKey,
			ParentId:      p.ParentId,
			DarkNumberLength: p.DarkNumberLength,

			AnssyAppSecret: p.AnssyAppSecret,
			AnssyToken:     p.AnssyToken,
			AnssyExpiredAt: anssyExpiredAt,

			TodayOrderNum:        todayOrderNum,
			TodayOrderAmount:     todayAmount,
			TodaySuccessAmount:   todaySuccessAmount,
			TodaySuccessOrderNum: todaySuccessOrderNum,

			Last1HourTotal:       lastDataItem.Last1HourTotal,
			Last1HourSuccess:     lastDataItem.Last1HourSuccess,
			Last30MinutesTotal:   lastDataItem.Last30MinutesTotal,
			Last30MinutesSuccess: lastDataItem.Last30MinutesSuccess,
		})
	}

	return &v1.ListPartnerResp{
		ListTableData: v1.ListTableData[v1.Partner]{
			List:  list,
			Total: total,
		},
	}, nil
}

func (s *PartnerService) SetPassword(c echo.Context, req *v1.PartnerSetPasswordReq, token string) (*v1.PartnerSetPasswordResp, error) {
	if len(req.NewPassword) < 6 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "密码应大于6位")
	}

	_, err := repository.Partner.SetPassword(c, token, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, err
	}

	return &v1.PartnerSetPasswordResp{}, nil
}

func (s *PartnerService) Update(c echo.Context, req *v1.PartnerUpdateReq) (*v1.PartnerUpdateResp, error) {
	_, err := repository.Partner.Update(c, req.Id, req)
	if err != nil {
		return nil, err
	}

	if req.Enable == int(model.Enabled) {
		s.ResetGoodsWeight(c, req.Id)
	}

	return &v1.PartnerUpdateResp{}, nil
}

func (s *PartnerService) UpdateBalance(c echo.Context, req *v1.PartnerUpdateBalanceReq) (*v1.PartnerUpdateBalanceResp, error) {
	if req.ChangeAmount == 0 {
		return &v1.PartnerUpdateBalanceResp{}, nil
	}

	from := model.BalanceFromTypeSystemAdd
	if req.ChangeAmount < 0 {
		from = model.BalanceFromTypeSystemDeduct
	}

	err := repository.Admin.CheckPassword(c, req.AdminId, req.Password)
	if err != nil {
		return nil, err
	}

	db := data.Instance()
	err = repository.Partner.UpdateBalance(c, db, req.PartnerId, "", req.ChangeAmount, from)
	if err != nil {
		return nil, err
	}

	return &v1.PartnerUpdateBalanceResp{}, nil
}

type FieldVale struct {
	Name  string
	Value any
}

func makeSign(aqsAppSecret string, params SyncGoodReqData) string {
	var fieldVales []FieldVale
	v := reflect.ValueOf(params)

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		fieldVales = append(fieldVales, FieldVale{Name: field.Name, Value: value.Interface()})
	}
	sort.Slice(fieldVales, func(i, j int) bool {
		return fieldVales[i].Name < fieldVales[j].Name
	})

	fieldList := lo.Map(fieldVales, func(f FieldVale, index int) string {
		if f.Value != "" {
			return fmt.Sprintf("%s%s", strings.ToLower(f.Name), cast.ToString(f.Value))
		} else {
			return ""
		}
	})
	fieldList = append([]string{aqsAppSecret}, fieldList...)
	fieldList = append(fieldList, aqsAppSecret)
	fieldStr := strings.Join(fieldList, "")
	firstHash := md5.Sum([]byte(fieldStr))
	hashString := hex.EncodeToString(firstHash[:])
	hashString = strings.ToUpper(hashString)

	return hashString
}

func fetchAqsGoodsPage(secret string, token string, pageNo int) (*SyncGoodRspData, error) {
	reqUrl := "http://gw.api.agiso.com/alds/Item/OnSaleGet"

	syncGoodReqData := SyncGoodReqData{
		Fields:    "approve_status,num_iid,title,nick,type,cid,pic_rul,num,price",
		Timestamp: time.Now().Unix(),
		PageNo:    pageNo,
	}

	sign := makeSign(secret, syncGoodReqData)

	syncGoodReqData.Sign = sign

	v, _ := query.Values(syncGoodReqData)

	client := &http.Client{}

	syncReq, err := http.NewRequest("POST", reqUrl, strings.NewReader(v.Encode()))

	if err != nil {
		return nil, nil
	}

	syncReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	syncReq.Header.Add("Authorization", "Bearer "+token)
	syncReq.Header.Add("ApiVersion", "1")

	resp, _ := client.Do(syncReq)

	defer func() {
		resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, nil
	}

	var datas SyncGoodRspData

	err = json.Unmarshal(body, &datas)
	if err != nil {
		fmt.Printf("Error parsing JSON: %s\n", err)
		return nil, nil
	}

	return &datas, nil
}

func fetchAqsGoods(secret string, token string) (*SyncGoodRspData, error) {

	var pageNo = 1

	datas, err := fetchAqsGoodsPage(secret, token, pageNo)

	if err != nil {
		return nil, err
	}

	if datas.IsSuccess && len(datas.Data) == 200 {
		for {
			pageNo++

			resp, err := fetchAqsGoodsPage(secret, token, pageNo)

			if err != nil {
				return nil, err
			}

			datas.Data = append(datas.Data, resp.Data...)

			if !datas.IsSuccess || len(resp.Data) < 200 {
				fmt.Printf("icccccccccccccccccccccc: pageNo %d\n", pageNo)
				fmt.Printf("icccccccccccccccccccccc: len resp.Data %d\n", len(resp.Data))
				fmt.Printf("icccccccccccccccccccccc: len datas.Data %d\n", len(datas.Data))
				break
			}
		}
	}

	return datas, nil
}

func (s *PartnerService) SyncGoods(c echo.Context, req *v1.PartnerSyncGoodsReq) (*v1.PartnerSyncGoodsResp, error) {
	p, err := repository.Partner.FindPartner(c, req.Id)
	if err != nil {
		return nil, err
	}

	if p.Type == model.PartnerTypeAnssy {
		return s.SyncAnssyGoods(c, p)
	}

	if p.ChannelId != model.ChannelTBPay {
		return nil, errors.New("仅支持淘宝直付通道")
	}

	if p.AqsAppSecret == "" || p.AqsToken == "" {
		return nil, errors.New("合作商阿奇索配置不全")
	}

	datas, err := fetchAqsGoods(p.AqsAppSecret, p.AqsToken)

	if err != nil {
		return nil, err
	}

	if datas.IsSuccess {
		_, err := repository.Partner.DeleteAllGoods(c, p.ID)

		if err != nil {
			return nil, err
		}

		for _, value := range datas.Data {
			status := model.GoodsStatusDisabled
			if value.ApproveStatus == "onsale" {
				status = model.GoodsStatusEnabled
			}

			price, _ := strconv.ParseFloat(value.Price, 64)
			goods := &model.Goods{
				PartnerId:  p.ID,
				SkuId:      fmt.Sprint(value.NumIid),
				Amount:     util.ToDecimal(price),
				RealAmount: util.ToDecimal(price),
				ShopName:   p.Nickname,
				Status:     status,
			}

			err := repository.Goods.Create(c, goods, false)
			if err != nil {
				return nil, err
			}
		}

	} else {
		return nil, errors.New(strconv.Itoa(datas.Error_Code) + datas.Error_Msg)
	}

	return &v1.PartnerSyncGoodsResp{}, nil
}

func (s *PartnerService) ListBalanceBill(c echo.Context, req *v1.ListPartnerBalanceBillReq) (*v1.ListPartnerBalanceBillResp, error) {
	var partnerIds []uint
	if req.PartnerId > 0 {
		partnerIds = append(partnerIds, req.PartnerId)
	} else {
		parentIds, _ := Admin.FindParentIds(c)

		partners, _, err := repository.Partner.List(c, &v1.ListPartnerReq{}, parentIds)
		if err != nil {
			return nil, err
		}

		partnerIds = lo.Map(partners, func(item *model.Partner, _ int) uint {
			return item.ID
		})

		fmt.Println(parentIds, partnerIds)
		if len(partnerIds) == 0 {
			return &v1.ListPartnerBalanceBillResp{}, nil
		}
	}

	bills, total, err := repository.Partner.ListBalanceBill(c, req, partnerIds)
	if err != nil {
		return nil, err
	}

	list := make([]*v1.PartnerBalanceBill, 0, len(bills))
	for _, m := range bills {
		list = append(list, &v1.PartnerBalanceBill{
			Id:           m.ID,
			PartnerId:    m.PartnerId,
			Nickname:     m.Nickname,
			OrderId:      m.OrderId,
			From:         int(m.From),
			Balance:      util.ToDecimal(m.Balance),
			ChangeAmount: util.ToDecimal(m.ChangeAmount),
			CreateAt:     m.CreatedAt.Unix(),
		})
	}

	return &v1.ListPartnerBalanceBillResp{
		ListTableData: v1.ListTableData[v1.PartnerBalanceBill]{
			List:  list,
			Total: total,
		},
	}, nil
}

func fetchAnssyGoodsPage(taobaoUserId string, token string, pageNo, pageSize int) (*SyncAnssyGoodRspData, error) {
	baseURL := "https://tao.anssy.com/charge/product/onsaleget"

	params := url.Values{}
	params.Set("taobaoUserId", taobaoUserId)
	params.Set("accessToken", token)
	params.Set("pageNum", strconv.Itoa(pageNo))
	params.Set("pageSize", strconv.Itoa(pageSize))
	params.Set("fields", "approve_status,num_iid,title,nick,type,cid,pic_rul,num,price")

	reqURL := baseURL + "?" + params.Encode()

	fmt.Printf("icccccccccccc reqURL %s", reqURL)

	client := &http.Client{}

	syncReq, err := http.NewRequest("GET", reqURL, nil)

	if err != nil {
		return nil, err
	}

	resp, _ := client.Do(syncReq)

	defer func() {
		resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var datas SyncAnssyGoodRspData

	err = json.Unmarshal(body, &datas)
	if err != nil {
		fmt.Printf("Error parsing JSON: %s\n", err)
		return nil, nil
	}

	// 交换id
	for i := 0; i < len(datas.Data); i++ {
		datas.Data[i].SkuId = datas.Data[i].NumIid
		datas.Data[i].NumIid = 0
	}

	return &datas, nil
}

func fetchAnssyGoodsSubSKU(taobaoUserId string, token string, numiid int) (*SyncAnssyGoodSubSKURspData, error) {
	baseURL := "https://tao.anssy.com/charge/product/getSku"

	params := url.Values{}
	params.Set("taobaoUserId", taobaoUserId)
	params.Set("accessToken", token)
	params.Set("numiid", strconv.Itoa(numiid))

	reqURL := baseURL + "?" + params.Encode()

	fmt.Printf("icccccccccccc reqURL %s", reqURL)

	client := &http.Client{}

	syncReq, err := http.NewRequest("GET", reqURL, nil)

	if err != nil {
		return nil, err
	}

	resp, _ := client.Do(syncReq)

	defer func() {
		resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var datas SyncAnssyGoodSubSKURspData

	err = json.Unmarshal(body, &datas)
	if err != nil {
		fmt.Printf("Error parsing JSON: %s\n", err)
		return nil, nil
	}

	return &datas, nil
}

func fetchAnssyGoods(taobaoUserId string, token string, isSubSku bool) (*SyncAnssyGoodRspData, error) {

	var pageNo = 1
	var pageSize = 40

	datas, err := fetchAnssyGoodsPage(taobaoUserId, token, pageNo, pageSize)

	if err != nil {
		return nil, err
	}

	if datas.Code == 200 && len(datas.Data) == pageSize {
		for {
			pageNo++

			resp, err := fetchAnssyGoodsPage(taobaoUserId, token, pageNo, pageSize)

			if err != nil {
				return nil, err
			}

			datas.Data = append(datas.Data, resp.Data...)

			if resp.Code != 200 || len(resp.Data) < pageSize {
				fmt.Printf("icccccccccccccccccccccc: pageNo %d\n", pageNo)
				fmt.Printf("icccccccccccccccccccccc: len resp.Data %d\n", len(resp.Data))
				fmt.Printf("icccccccccccccccccccccc: len datas.Data %d\n", len(datas.Data))
				break
			}
		}
	}
	fmt.Printf("icccccccccccccccccccccc: len total %d\n", len(datas.Data))

	if isSubSku {

		subDatas := make([]GoodsData, 0)

		for i := 0; i < len(datas.Data); i++ {
			numId := datas.Data[i].SkuId
			resp, err := fetchAnssyGoodsSubSKU(taobaoUserId, token, numId)

			if err != nil {
				return nil, err
			}

			if resp.Data != nil && resp.Code == 200 {
				for j := 0; j < len(resp.Data); j++ {
					subData := GoodsData{
						ApproveStatus: datas.Data[i].ApproveStatus,
						Nick:          datas.Data[i].Nick,
						Num:           resp.Data[j].Quantity,
						NumIid:        numId,
						SkuId:         resp.Data[j].SkuId,
						Price:         datas.Data[i].Price,
						Title:         datas.Data[i].Title,
					}
					subDatas = append(subDatas, subData)
				}
			}
		}

		fmt.Printf("icccccccccccccccccccccc: len subDatas %d\n", len(subDatas))

		if len(subDatas) > 0 {
			datas = &SyncAnssyGoodRspData{
				Code: 200,
				Data: subDatas,
			}
		}
	}

	return datas, nil
}

func (s *PartnerService) SyncAnssyGoods(c echo.Context, p *model.Partner) (*v1.PartnerSyncGoodsResp, error) {
	// partner, err := repository.Partner.FindPartner(c, req.Id)
	// if err != nil {
	// 	return nil, err
	// }

	if p.ChannelId != model.ChannelTBPay {
		return nil, errors.New("仅支持淘宝直付通道")
	}

	// TODO
	// 测试安式拉商品，暂时屏蔽

	if p.AnssyAppSecret == "" || p.AnssyToken == "" {
		return nil, errors.New("合作商配置不全")
	}

	var datas *SyncAnssyGoodRspData

	// TODO
	// 这里要换成数据库中合作商对应的taobaoUserId
	token := p.AnssyToken
	taobaoUserId := p.AnssyTbUserId
	isSubSku := true // 收否多规格
	datas, err := fetchAnssyGoods(taobaoUserId, token, isSubSku)
	if err != nil {
		return nil, err
	}

	if datas != nil && datas.Code == 200 {
		_, err := repository.Partner.DeleteAllGoods(c, p.ID)

		if err != nil {
			return nil, err
		}

		for _, value := range datas.Data {
			status := model.GoodsStatusDisabled
			if value.ApproveStatus == "onsale" {
				status = model.GoodsStatusEnabled
			}

			price, _ := strconv.ParseFloat(value.Price, 64)
			goods := &model.Goods{
				PartnerId:  p.ID,
				SkuId:      fmt.Sprint(value.SkuId),
				NumId:      fmt.Sprint(value.NumIid),
				Amount:     util.ToDecimal(price),
				RealAmount: util.ToDecimal(price),
				ShopName:   p.Nickname,
				Status:     status,
			}

			err := repository.Goods.Create(c, goods, false)
			if err != nil {
				c.Logger().Errorf("SyncAnssyGoods error=%s", err)
				continue
			}
		}

	} else {
		return nil, nil
	}

	return &v1.PartnerSyncGoodsResp{}, nil
}

func (s *PartnerService) ResetVerifiCode(c echo.Context, req *v1.PartnerResetVerifiCodeReq) (*v1.PartnerResetPasswordResp, error) {
	user, err := repository.Partner.ResetVerifiCode(c, req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.PartnerResetPasswordResp{
		Password: user.Password,
	}, nil
}
