package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"apollo/server/internal/model"
	"apollo/server/pkg/util"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 获取商户数据
func getMerchants(db *gorm.DB) ([]model.Merchant, error) {
	var merchants []model.Merchant
	err := db.Find(&merchants).Error
	return merchants, err
}

// 获取合作商数据
func getPartners(db *gorm.DB) ([]model.Partner, error) {
	var partners []model.Partner
	err := db.Find(&partners).Error
	return partners, err
}

func main() {
	// 数据库连接
	dsn := "root:12345678@tcp(127.0.0.1:3306)/apollo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 获取现有的商户和合作商
	merchants, err := getMerchants(db)
	if err != nil {
		log.Fatal("获取商户数据失败:", err)
	}

	partners, err := getPartners(db)
	if err != nil {
		log.Fatal("获取合作商数据失败:", err)
	}

	if len(merchants) == 0 || len(partners) == 0 {
		log.Fatal("没有找到商户或合作商数据，请先添加基础数据")
	}

	fmt.Printf("找到 %d 个商户和 %d 个合作商\n", len(merchants), len(partners))

	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 今天的开始时间
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	fmt.Printf("开始生成 100 条今天8点前的订单数据...\n")

	// 渠道ID列表
	channelIds := []string{
		"JD000000", // 京东游戏
		"JD111111", // 京东实物
		"TB000000", // 淘宝店铺
		"TB111111", // 淘宝码上收
		"TB222222", // 淘宝复制
	}

	// 生成8点前的订单（0:00-7:59）
	var orders []model.Order
	for i := 0; i < 100; i++ {
		// 随机选择商户和合作商
		merchant := merchants[rand.Intn(len(merchants))]
		partner := partners[rand.Intn(len(partners))]

		// 生成订单ID
		orderId, err := util.SFlake.GenString()
		if err != nil {
			log.Printf("生成订单ID失败: %v", err)
			continue
		}

		// 随机生成8点前的订单时间（0:00-7:59）
		randomSeconds := rand.Int63n(int64(8 * time.Hour / time.Second)) // 8小时内的随机秒数
		orderTime := startOfDay.Add(time.Duration(randomSeconds) * time.Second)

		// 随机生成订单数据
		order := model.Order{
			OrderId:         orderId,
			ChannelId:       model.ChannelId(channelIds[rand.Intn(len(channelIds))]),
			MerchantId:      merchant.ID,
			MerchantOrderId: fmt.Sprintf("M%s", orderId),
			PartnerOrderId:  fmt.Sprintf("P%s", orderId),
			Amount:          float64(rand.Intn(10000)+100) / 100, // 1-100元
			ReceivedAmount:  0,
			PayType:         []model.PayType{model.AppPay, model.WebPay, model.AliPay}[rand.Intn(3)],
			PayAccount:      fmt.Sprintf("early_user%d@example.com", rand.Intn(1000)),
			Status:          []model.OrderStatus{model.OrderStatusUnpaid, model.OrderStatusPaid, model.OrderStatusFinish}[rand.Intn(3)],
			SkuId:           fmt.Sprintf("EARLY_SKU%d", rand.Intn(100)),
			PartnerId:       partner.ID,
			MerchantName:    merchant.Nickname,
			PartnerName:     partner.Nickname,
			Shop:            fmt.Sprintf("早期店铺%d", rand.Intn(10)),
			NotifyStatus:    model.NotNotify,
			NotifyUrl:       "https://example.com/notify",
			PayUrl:          "https://example.com/pay",
			ExtParam:        "{}",
			IP:              fmt.Sprintf("192.168.1.%d", rand.Intn(255)),
			Device:          []string{"mobile", "pc", "tablet"}[rand.Intn(3)],
			DarkNumber:      fmt.Sprintf("EARLY_DN%d", rand.Intn(10000)),
			Remark:          fmt.Sprintf("8点前测试订单%d - %s", i+1, orderTime.Format("15:04:05")),
			ConfirmStatus:   model.ConfirmStatusAuto,
		}
		
		// 设置创建时间
		order.CreatedAt = orderTime
		order.UpdatedAt = orderTime

		// 如果订单已支付，设置支付时间和实收金额
		if order.Status >= model.OrderStatusPaid {
			order.PayAt = orderTime.Add(time.Duration(rand.Intn(1800)) * time.Second) // 30分钟内支付
			order.ReceivedAmount = order.Amount
		}

		orders = append(orders, order)
	}

	// 批量插入订单
	if err := db.CreateInBatches(orders, 50).Error; err != nil {
		log.Printf("批量插入订单失败: %v", err)
		return
	}

	fmt.Printf("8点前订单数据生成完成！\n")

	// 验证插入的数据
	var count int64
	db.Model(&model.Order{}).Where("DATE(created_at) = ? AND HOUR(created_at) < 8", startOfDay.Format("2006-01-02")).Count(&count)
	fmt.Printf("今天8点前总共有 %d 条订单记录\n", count)
}