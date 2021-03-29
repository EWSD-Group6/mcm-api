package common

import (
	"encoding/base64"
	"encoding/json"
	"mcm-api/pkg/apperror"
	"reflect"
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
	Limit int    `query:"limit"`
}

func (c CursorQuery) GetNext(next interface{}) error {
	decoded, err := base64.URLEncoding.DecodeString(c.Next)
	if err != nil {
		return apperror.New(apperror.ErrInvalid, "invalid next param", err)
	}
	return json.Unmarshal(decoded, next)
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

func NewPaginateResponse(data interface{}, total int64, page int, limit int) *PaginateResponse {
	r := &PaginateResponse{
		Total:       total,
		CurrentPage: page,
		LastPage:    calculateLastPage(total, limit),
		PerPage:     limit,
	}
	if !isNil(data) {
		r.Data = data
	} else {
		r.Data = make([]interface{}, 0)
	}
	return r
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
	Next string      `json:"next,omitempty"`
	Data interface{} `json:"data"`
}

func isNil(a interface{}) bool {
	return a == nil || reflect.ValueOf(a).IsNil()
}

func NewCursorResponse(data interface{}, next interface{}) *CursorResponse {
	var nextStr string
	if !isNil(next) {
		nextJson, _ := json.Marshal(next)
		nextStr = base64.StdEncoding.EncodeToString(nextJson)
	}
	return &CursorResponse{
		Next: nextStr,
		Data: data,
	}
}
