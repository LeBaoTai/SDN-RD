package model

type Device struct {
	ID      string `gorm:"type:uuid;primaryKey"`
	Name    string `gorm:"type:varchar(50);not null"`
	Address string `gorm:"type:inet;not null"`
	Port    int    `gorm:"type:integer;not null"`
}

func (Device) TableName() string {
	return "public.devices"
}
