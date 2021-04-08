package systemdata

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type DataUpdateReq struct {
	Value string `json:"value"`
}

func (d DataUpdateReq) Validate() error {
	return validation.ValidateStruct(&d, validation.Field(&d.Value, validation.Required))
}

type DataRes struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Type      ValueType `json:"type" enums:"document,int,string"`
	UpdatedAt time.Time `json:"updatedAt"`
}
