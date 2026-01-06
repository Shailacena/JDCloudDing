package v1

type GoodsInfo struct {
	PartnerId uint    `json:"partnerId"`
	SkuId     string  `json:"skuId"`
	BrandId   string  `json:"brandId"`
	Price     float64 `json:"price"`
	RealPrice float64 `json:"realPrice"`
	ShopName  string  `json:"shopName"`
	Status    int     `json:"status"`
}

// 商品创建
type GoodsCreateReq struct {
	GoodsInfo
}

type GoodsCreateResp struct {
}

// 商品创建
type GoodsUpdateReq struct {
	Id uint `json:"id" validate:"required"`
	GoodsInfo
}

type GoodsUpdateResp struct {
}

// 删除商品
type GoodsDeleteReq struct {
	Id uint `json:"id" validate:"required"`
}

type GoodsDeleteResp struct {
}

// 商品列表
type ListGoodsReq struct {
	PartnerId uint   `json:"partnerId"`
	SkuId     string `json:"skuId"`
	Pagination
}

type ListGoodsResp struct {
	ListTableData[Goods]
}

type Goods struct {
	Id uint `json:"id"`
	GoodsInfo
	CreateAt int64 `json:"createAt"`
}

// 合作商同步商品
type PartnerSyncGoodsReq struct {
	Id uint `json:"id" validate:"required"`
}

type PartnerSyncGoodsResp struct {
}
