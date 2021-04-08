package systemdata

import "time"

type ValueType string

const (
	Document ValueType = "document"
	Int      ValueType = "int"
	String   ValueType = "string"
)

type Entity struct {
	Key       string `gorm:"primaryKey"`
	Value     string
	Type      ValueType
	UpdatedAt time.Time
}

func (e *Entity) TableName() string {
	return "system_data"
}
