package v1

type ListOperationLogReq struct {
	Pagination
}

type ListOperationLogResp struct {
	ListTableData[OperationLog]
}

type OperationLog struct {
	Id           uint   `json:"id"`
	IP           string `json:"ip"`
	Operation    string `json:"operation"`
	OperationStr string `json:"operationStr"`
	Operator     uint   `json:"operator"`
	OperatorName string `json:"operatorName"`
	CreateAt     int64  `json:"createAt"`
}
