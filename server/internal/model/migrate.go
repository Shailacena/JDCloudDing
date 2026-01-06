package model

import "gorm.io/gorm"

func InitMigrate(db *gorm.DB) {
	db.Set("gorm:table_options", "AUTO_INCREMENT=10000").AutoMigrate(&SysUser{})
	db.Set("gorm:table_options", "AUTO_INCREMENT=10000").AutoMigrate(&Partner{})
	db.Set("gorm:table_options", "AUTO_INCREMENT=10000").AutoMigrate(&Merchant{})
	db.Set("gorm:table_options", "AUTO_INCREMENT=10000").AutoMigrate(&JDAccount{})
	db.AutoMigrate(&RealNameAccount{})
	// db.AutoMigrate(&PartnerBill{})
	// db.AutoMigrate(&DailyBill{})
	db.AutoMigrate(&Goods{})
	db.AutoMigrate(&Order{})
	db.AutoMigrate(&Notify{})
	db.AutoMigrate(&PartnerBalanceBill{})
	db.AutoMigrate(&MerchantBalanceBill{})
	db.AutoMigrate(&OperationLog{})
	db.AutoMigrate(&OrderArchive{})
}
