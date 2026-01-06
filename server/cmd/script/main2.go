package main

import (
	// mw "apollo/server/internal/middleware"
	// "apollo/server/internal/model"
	// "apollo/server/internal/router"
	// "apollo/server/pkg/app"
	// "apollo/server/pkg/config"
	// "apollo/server/pkg/data"
	// "apollo/server/pkg/validate"
	// "fmt"

	// "github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
	// "github.com/labstack/gommon/log"

	"fmt"
	"os/exec"
	"strings"
)

func main() {
	// getpayurl 获取支付链接
	// checkorder 查询支付状态

	//======================
	// action := "getpayurl"
	// ck := "pin=jd_44754e08b8767;wskey=AAJnqh-tADAzXrPjDuVgolEya6lscGQnnbW2IxYvV_Lzwj1_9aoHgxL6zd1GEWMuHVtky94yqzs;"
	// sku := "10077221265581"
	// adress := "李先生 13756376578 江西省宜春地区宜春市 建设路11号19A"
	// // ip := "http://211.95.152.42:11641"
	// ip := ""
	// orderid := "10000001"
	// jdorderid := ""

	//===========================
	action := "getpayurl"
	ck := "pin=jd_rnQayWzSzxie;wskey=AAJnjFGIAECZzqoiHl80yojPuCakPwEASKc7OrjjUfCj6NDFgdJsXoy-x9e5QgzXrf5L1pvHF3ymS33_VTrel-tpFtTwNO8U;"
	sku := "10077221265581"
	adress := "李先生 13756376578 江西省宜春地区宜春市 建设路11号19A"
	ip := "211.95.152.52:12882"
	orderid := "df7643b5586d43b49bc3ce17487f687c"
	jdorderid := ""

	cmd := exec.Command("python3", "../jd/main.py", action, ck, sku, orderid, jdorderid, adress, ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// 解析输出结果，查找异常信息
	outputStr := string(output)
	if strings.Contains(outputStr, "发生错误:") {
		fmt.Println("Python脚本抛出异常:", strings.TrimPrefix(outputStr, "发生错误:"))
	} else {
		fmt.Println("Python脚本输出:", outputStr)
	}

	// conf := config.New("configs/config.yaml")

	// db := data.Init(conf.MysqlConfig)
	// model.InitMigrate(db)

	// e := app.Engine()

	// e.Logger.SetLevel(log.INFO)

	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	// e.Use(middleware.RequestID())

	// e.Validator = validate.NewReqValidator()

	// e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	// 	c.Logger().Info("reqBody:", string(reqBody))
	// 	c.Logger().Info("resBody:", string(resBody))
	// }))

	// e.Use(mw.HandleErrorMiddleware())

	// router.Init(e)

	// e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.HttpConfig.Port)))
}
