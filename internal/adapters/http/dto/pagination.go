package dto

import (
	"fmt"
	"math"
)

type PaginationParam struct {
	Page    int `json:"page" query:"page"`
	PerPage int `json:"per_page" query:"per_page"`
}

type metadata struct {
	Page       int                 `json:"page"`
	PerPage    int                 `json:"per_page"`
	PageCount  int                 `json:"page_count"`
	TotalCount int64               `json:"total_count"`
	Links      []map[string]string `json:"Links"`
}

type paginationResponse[T any] struct {
	Metadata metadata `json:"_metadata"`
	Records  []T      `json:"records"`
}

func NewPaginationResponse[T any](records []T, totalCount int64, param PaginationParam, basePath string) *paginationResponse[T] {
	if param.Page < 1 {
		param.Page = 1
	}
	if param.PerPage <= 0 {
		param.PerPage = 10
	}

	pageCount := int(math.Ceil(float64(totalCount) / float64(param.PerPage)))
	if pageCount == 0 {
		pageCount = 1
	}

	if param.Page > pageCount {
		param.Page = pageCount
	}

	links := []map[string]string{
		{"self": fmt.Sprintf("%s?page=%d&per_page=%d", basePath, param.Page, param.PerPage)},
		{"first": fmt.Sprintf("%s?page=1&per_page=%d", basePath, param.PerPage)},
	}

	if param.Page > 1 {
		links = append(links, map[string]string{"previous": fmt.Sprintf("%s?page=%d&per_page=%d", basePath, param.Page-1, param.PerPage)})
	} else {
		links = append(links, map[string]string{"previous": ""})
	}

	if param.Page < pageCount {
		links = append(links, map[string]string{"next": fmt.Sprintf("%s?page=%d&per_page=%d", basePath, param.Page+1, param.PerPage)})
	} else {
		links = append(links, map[string]string{"next": ""})
	}

	links = append(links, map[string]string{"last": fmt.Sprintf("%s?page=%d&per_page=%d", basePath, pageCount, param.PerPage)})

	return &paginationResponse[T]{
		Metadata: metadata{
			Page:       param.Page,
			PerPage:    param.PerPage,
			PageCount:  pageCount,
			TotalCount: totalCount,
			Links:      links,
		},
		Records: records,
	}
}