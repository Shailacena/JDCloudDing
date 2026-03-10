package v1

type CardInfo struct {
	CardNo    string  `json:"cardNo"`
	Password  string  `json:"password"`
	CardGroup string  `json:"cardGroup"`
	Amount    float64 `json:"amount"`
}

type CardCreateReq struct {
	Cards []CardInfo `json:"cards" validate:"required"`
}

type CardCreateResp struct {
	Count int `json:"count"`
}

type VirtualCardGenerateReq struct {
	Prefix      string  `json:"prefix" validate:"required"`
	CardNoLen   int     `json:"cardNoLen" validate:"required,min=12"`
	PasswordLen int     `json:"passwordLen" validate:"required"`
	CardGroup   string  `json:"cardGroup" validate:"required"`
	Amount      float64 `json:"amount" validate:"required"`
	Count       int     `json:"count" validate:"required,min=1"`
}

type VirtualCardGenerateResp struct {
	Count int `json:"count"`
}

type ListCardReq struct {
	CardNo    string `json:"cardNo"`
	CardGroup string `json:"cardGroup"`
	BatchNo   string `json:"batchNo"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	CardType  string `json:"cardType"`
	Pagination
}

type ListCardResp struct {
	ListTableData[PriceCard]
}

type PriceCard struct {
	Id         uint   `json:"id"`
	CardNo     string `json:"cardNo"`
	Password   string `json:"password"`
	CardGroup  string `json:"cardGroup"`
	Amount     float64 `json:"amount"`
	CardType   string `json:"cardType"`
	BatchNo    string `json:"batchNo"`
	CardStatus string `json:"cardStatus"`
	OrderId    string `json:"orderId"`
	Remark     string `json:"remark"`
	UsedAt     int64  `json:"usedAt"`
	CreateAt   int64  `json:"createAt"`
}

type DeleteCardReq struct {
	Ids []uint `json:"ids" validate:"required"`
}

type DeleteCardResp struct {
	Count int `json:"count"`
}
