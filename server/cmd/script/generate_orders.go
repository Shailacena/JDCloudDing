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

// 使用项目中的model结构体

// 使用项目中的constants

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
	dsn := "root:abc123!!!@tcp(127.0.0.1:3306)/gva?charset=utf8mb4&parseTime=True&loc=Local"
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

	// 今天的开始和结束时间
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	fmt.Printf("开始生成 %d 条今天的订单数据...\n", 10000)

	// 渠道ID列表
	channelIds := []string{
		"JD000000", // 京东游戏
		"JD111111", // 京东实物
		"TB000000", // 淘宝店铺
		"TB111111", // 淘宝码上收
		"TB222222", // 淘宝复制
		"TB668888", // 淘宝代付
		"TB686088", // 天猫直付
		"TB087888", // 淘宝电子券
		"JS000000", // 京东复制
		"JS111111", // 京东ck
	}

	// 批量生成订单
	batchSize := 1000
	for i := 0; i < 10000; i += batchSize {
		var orders []model.Order
		end := i + batchSize
		if end > 10000 {
			end = 10000
		}

		for j := i; j < end; j++ {
			// 随机选择商户和合作商
			merchant := merchants[rand.Intn(len(merchants))]
			partner := partners[rand.Intn(len(partners))]

			// 生成订单ID
			orderId, err := util.SFlake.GenString()
			if err != nil {
				log.Printf("生成订单ID失败: %v", err)
				continue
			}

			// 随机生成订单时间（今天内）
			randomSeconds := rand.Int63n(int64(24 * time.Hour / time.Second))
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
				PayAccount:      fmt.Sprintf("user%d@example.com", rand.Intn(10000)),
				Status:          []model.OrderStatus{model.OrderStatusUnpaid, model.OrderStatusPaid, model.OrderStatusFinish}[rand.Intn(3)],
				SkuId:           fmt.Sprintf("SKU%d", rand.Intn(1000)),
				PartnerId:       partner.ID,
				MerchantName:    merchant.Nickname,
				PartnerName:     partner.Nickname,
				Shop:            fmt.Sprintf("店铺%d", rand.Intn(100)),
				NotifyStatus:    model.NotNotify,
				NotifyUrl:       "https://example.com/notify",
				PayUrl:          "https://example.com/pay",
				ExtParam:        "{}",
				IP:              fmt.Sprintf("192.168.1.%d", rand.Intn(255)),
				Device:          []string{"mobile", "pc", "tablet"}[rand.Intn(3)],
				DarkNumber:      fmt.Sprintf("DN%d", rand.Intn(100000)),
				Remark:          fmt.Sprintf("测试订单%d", j+1),
				ConfirmStatus:   model.ConfirmStatusAuto,
			}
			
			// 设置创建时间
			order.CreatedAt = orderTime
			order.UpdatedAt = orderTime

			// 如果订单已支付，设置支付时间和实收金额
			if order.Status >= model.OrderStatusPaid {
				order.PayAt = orderTime.Add(time.Duration(rand.Intn(3600)) * time.Second)
				order.ReceivedAmount = order.Amount
			}

			orders = append(orders, order)
		}

		// 批量插入订单
		if err := db.CreateInBatches(orders, batchSize).Error; err != nil {
			log.Printf("批量插入订单失败 (批次 %d-%d): %v", i+1, end, err)
			continue
		}

		fmt.Printf("已生成订单: %d/%d\n", end, 10000)
	}

	fmt.Printf("订单数据生成完成！\n")

	// 验证插入的数据
	var count int64
	db.Model(&model.Order{}).Where("DATE(created_at) = ?", startOfDay.Format("2006-01-02")).Count(&count)
	fmt.Printf("今天总共有 %d 条订单记录\n", count)
}