package router

import (
	"apollo/server/internal/handler"
	"apollo/server/internal/middleware"
	"apollo/server/internal/repository"

	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo) {
	apiGroup := e.Group("/web_api")

	adminTokenChecker := middleware.GenAuthHandler(repository.Admin)
	// adminRoleChecker := middleware.CheckRoleHandler(model.SuperAdminRole)
	adminGroup := apiGroup.Group("/admin", adminTokenChecker(), middleware.HandleOperationLogMiddleware())
	adminGroupWithoutAuth := apiGroup.Group("/admin")
	{
		adminGroupWithoutAuth.POST("/login", handler.Admin.Login)
		adminGroupWithoutAuth.POST("/register11", handler.Admin.Register11)
		adminGroup.POST("/register", handler.Admin.Register)
		adminGroup.GET("/list", handler.Admin.List)
		adminGroupWithoutAuth.POST("/setPassword", handler.Admin.SetPassword)
		adminGroup.POST("/resetPassword", handler.Admin.ResetPassword)
		adminGroup.POST("/delete", handler.Admin.Delete)
		adminGroup.POST("/update", handler.Admin.Update)
		adminGroup.POST("/enable", handler.Admin.Enable)
		adminGroup.POST("/resetVerifiCode", handler.Admin.ResetVerifiCode)
		adminGroupWithoutAuth.POST("/logout", handler.Admin.Logout)
		adminGroupWithoutAuth.GET("/operationLog", handler.OperationLog.List)
		adminGroup.POST("/master/income", handler.Admin.GetMasterIncome)
		// 服务器状态与工具
		adminGroup.GET("/server/status", handler.Status.GetServerStatus)
		adminGroup.GET("/order/trend/today", handler.Status.TodayTrend)
	}

	partnerGroup := apiGroup.Group("/partner", adminTokenChecker(), middleware.HandleOperationLogMiddleware())
	{
		partnerGroup.POST("/register", handler.Partner.Register)
		partnerGroup.POST("/resetPassword", handler.Partner.ResetPassword)
		partnerGroup.POST("/delete", handler.Partner.Delete)
		partnerGroup.GET("/list", handler.Partner.List)
		partnerGroup.POST("/update", handler.Partner.Update)
		partnerGroup.POST("/updateBalance", handler.Partner.UpdateBalance)
		partnerGroup.GET("/listBalanceBill", handler.Partner.ListBalanceBill)
		partnerGroup.POST("/syncGoods", handler.Partner.SyncGoods)
		partnerGroup.POST("/resetVerifiCode", handler.Partner.ResetVerifiCode)
	}

	merchantGroup := apiGroup.Group("/merchant", adminTokenChecker(), middleware.HandleOperationLogMiddleware())
	{

		merchantGroup.POST("/update", handler.Merchant.Update)
		merchantGroup.POST("/updateBalance", handler.Merchant.UpdateBalance)
		merchantGroup.POST("/register", handler.Merchant.Register)
		merchantGroup.POST("/resetPassword", handler.Merchant.ResetPassword)
		merchantGroup.GET("/list", handler.Merchant.List)
		merchantGroup.GET("/listBalanceBill", handler.Merchant.ListBalanceBill)
		merchantGroup.POST("/enable", handler.Merchant.Enable)
		merchantGroup.POST("/resetVerifiCode", handler.Merchant.ResetVerifiCode)
	}

	realNameAccountGroup := apiGroup.Group("/realNameAccount", adminTokenChecker(), middleware.HandleOperationLogMiddleware())
	{
		realNameAccountGroup.POST("/create", handler.RealNameAccount.Create)
		realNameAccountGroup.GET("/list", handler.RealNameAccount.List)
	}

	jdAccountGroup := apiGroup.Group("/jdAccount", adminTokenChecker(), middleware.HandleOperationLogMiddleware())
	{
		jdAccountGroup.POST("/create", handler.JDAccount.Create)
		jdAccountGroup.POST("/enable", handler.JDAccount.Enable)
		jdAccountGroup.GET("/list", handler.JDAccount.List)
		jdAccountGroup.POST("/delete", handler.JDAccount.Delete)
		jdAccountGroup.POST("/reset", handler.JDAccount.Reset)
	}

	statisticsGroup := apiGroup.Group("/statistics", adminTokenChecker(), middleware.HandleOperationLogMiddleware())
	{
		statisticsGroup.GET("/listDailyBill", handler.Statistics.ListDailyBill)
		statisticsGroup.GET("/listDailyBillByPartner", handler.Statistics.ListDailyBillByPartner)
		statisticsGroup.GET("/listDailyBillByMerchant", handler.Statistics.ListDailyBillByMerchant)
	}

	goodsGroup := apiGroup.Group("/goods", adminTokenChecker(), middleware.HandleOperationLogMiddleware())
	{
		goodsGroup.POST("/create", handler.Goods.Create)
		goodsGroup.POST("/update", handler.Goods.Update)
		goodsGroup.POST("/delete", handler.Goods.Delete)
		goodsGroup.POST("/list", handler.Goods.List)
	}

	orderGroup := apiGroup.Group("/order", adminTokenChecker(), middleware.HandleOperationLogMiddleware())
	{
		orderGroup.GET("/list", handler.Order.List)
		orderGroup.GET("/summary", handler.Order.GetOrderSummary)
		orderGroup.POST("/confirm", handler.Order.Confirm)
		orderGroup.POST("/archive", handler.Order.Archive)
	}

	// 合作商后台
	partner1TokenChecker := middleware.GenAuthHandler(repository.Partner)
	partner1Group := apiGroup.Group("/partner1", partner1TokenChecker())
	partner1GroupWithoutAuth := apiGroup.Group("/partner1")

	{
		partner1GroupWithoutAuth.POST("/login", handler.Partner.Login)
		partner1GroupWithoutAuth.POST("/logout", handler.Partner.Logout)
		partner1Group.POST("/setPassword", handler.Partner.SetPassword)
		partner1Group.POST("/syncGoods", handler.Partner.SyncGoods)
		partner1Group.GET("/listBalanceBill", handler.Partner.ListBalanceBill)

		partner1StatisticsGroup := partner1Group.Group("/statistics")
		{
			partner1StatisticsGroup.GET("/listBill", handler.Statistics.ListDailyBill)
		}

		partner1OrderGroup := partner1Group.Group("/order")
		{
			partner1OrderGroup.GET("/list", handler.Order.List)
		}

		partner1GoodsGroup := partner1Group.Group("/goods")
		{
			partner1GoodsGroup.POST("/list", handler.Goods.List)
		}
	}

	// 商户后台
	merchant1TokenChecker := middleware.GenAuthHandler(repository.Merchant)
	merchant1Group := apiGroup.Group("/merchant1", merchant1TokenChecker())
	merchant1GroupWithoutAuth := apiGroup.Group("/merchant1")
	{
		merchant1GroupWithoutAuth.POST("/login", handler.Merchant.Login)
		merchant1GroupWithoutAuth.POST("/logout", handler.Merchant.Logout)
		merchant1Group.POST("/setPassword", handler.Merchant.SetPassword)
		merchant1Group.GET("/listBalanceBill", handler.Merchant.ListBalanceBill)
		merchant1Group.GET("/getBalance", handler.Merchant.GetBalance)

		merchant1StatisticsGroup := merchant1Group.Group("/statistics")
		{
			merchant1StatisticsGroup.GET("/listBill", handler.Statistics.ListDailyBill)
		}

		merchant1OrderGroup := merchant1Group.Group("/order")
		{
			merchant1OrderGroup.GET("/list", handler.Order.List)
		}
	}
}
