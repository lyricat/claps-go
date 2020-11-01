package model

type PaginationQ struct {
	Q      string `form:"q" json:"q"`
	Limit  int    `form:"limit" json:"limit"`
	Offset int    `form:"offset" json:"offset"`
	Total  int    `json:"total"`
}
