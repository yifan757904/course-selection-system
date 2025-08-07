package model

type Pagination struct {
	Page     int `form:"page" binding:"min=1"`              // 当前页码，最小为1
	PageSize int `form:"page_size" binding:"min=5,max=100"` // 每页数量，范围5-100
}

// 允许排序的字段白名单
var AllowedSortFields = map[string]bool{
	"id":        true,
	"hours":     true,
	"startdate": true,
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p Pagination) Limit() int {
	return p.PageSize
}

type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`       // 总记录数
	Page       int   `json:"page"`        // 当前页码
	PageSize   int   `json:"page_size"`   // 每页数量
	TotalPages int   `json:"total_pages"` // 总页数
}
