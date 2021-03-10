package common

import (
	"encoding/base64"
	"encoding/json"
	"mcm-api/pkg/apperror"
)

const (
	limitDefault = 20
	limitMax     = 100
)

type PaginateQuery struct {
	Limit int `query:"limit"`
	Page  int `query:"page"`
}

type CursorQuery struct {
	Next  string `query:"next"`
	Prev  string `query:"prev"`
	Limit int    `query:"limit"`
}

func (c CursorQuery) GetNext(next interface{}) error {
	var str []byte
	_, err := base64.URLEncoding.Decode(str, []byte(c.Next))
	if err != nil {
		return apperror.New(apperror.ErrInvalid, "invalid next param", err)
	}
	return json.Unmarshal(str, next)
}

func (c CursorQuery) GetPrev(prev interface{}) error {
	var str []byte
	_, err := base64.URLEncoding.Decode(str, []byte(c.Prev))
	if err != nil {
		return apperror.New(apperror.ErrInvalid, "invalid prev param", err)
	}
	return json.Unmarshal(str, prev)
}

func (c CursorQuery) GetLimit() int {
	if c.Limit > limitMax {
		return limitMax
	}
	if c.Limit <= 0 {
		return limitDefault
	}
	return c.Limit
}

func (q PaginateQuery) GetOffSet() int {
	return q.GetLimit() * q.Page
}

func (q PaginateQuery) GetLimit() int {
	if q.Limit > limitMax {
		return limitMax
	}
	if q.Limit <= 0 {
		return limitDefault
	}
	return q.Limit
}

type PaginateResponse struct {
	Total       int64       `json:"total"`
	CurrentPage int         `json:"currentPage"`
	LastPage    int         `json:"lastPage"`
	PerPage     int         `json:"perPage"`
	Data        interface{} `json:"data"`
}

func NewEmptyPaginateResponse() *PaginateResponse {
	return &PaginateResponse{
		Total:       0,
		CurrentPage: 0,
		LastPage:    0,
		PerPage:     limitDefault,
		Data:        nil,
	}
}

func NewPaginateResponse(data interface{}, total int64, page int, limit int) *PaginateResponse {
	return &PaginateResponse{
		Total:       total,
		CurrentPage: page,
		LastPage:    calculateLastPage(total, limit),
		PerPage:     limit,
		Data:        data,
	}
}

func calculateLastPage(total int64, limit int) int {
	if total == 0 {
		return 0
	}
	if total%int64(limit) > 0 {
		return int(total / int64(limit))
	}
	return int(total/int64(limit)) - 1
}

type CursorResponse struct {
	Next string      `json:"next"`
	Prev string      `json:"prev"`
	Data interface{} `json:"data"`
}

func NewCursorResponse(data interface{}, next interface{}, prev interface{}) *CursorResponse {
	nextStr, _ := json.Marshal(next)
	prevStr, _ := json.Marshal(prev)
	var next64, prev64 []byte
	base64.URLEncoding.Encode(next64, nextStr)
	base64.URLEncoding.Encode(prev64, prevStr)
	return &CursorResponse{
		Next: string(next64),
		Prev: string(prevStr),
		Data: data,
	}
}
