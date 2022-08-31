package entity

import "fmt"

const itemsPerPage uint64 = 5

func NewPagination(page uint64, descName, descCreated bool) Pagination {

	return Pagination{
		Page:         page,
		DescName:     descName,
		DescCreated:  descCreated,
		ItemsPerPage: itemsPerPage,
		Key:          fmt.Sprintf("%v_%v_%v", page, descName, descCreated),
	}
}

type Pagination struct {
	Page         uint64 `json:"page" form:"num"`
	DescName     bool   `json:"name" form:"name"`
	DescCreated  bool   `json:"created" form:"created"`
	ItemsPerPage uint64
	Key          string `json:"key" form:"key"` // for cache
}
