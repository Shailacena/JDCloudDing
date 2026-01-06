package v1

type ListTableData[T any] struct {
	List  []*T  `json:"list"`
	Total int64 `json:"total"`
}

type Pagination struct {
	CurrentPage int `query:"currentPage"`
	PageSize    int `query:"pageSize"`
}

func (p *Pagination) Offset() int {
	return (p.CurrentPage - 1) * p.PageSize
}

func (p *Pagination) Limit() int {
	return p.PageSize
}
