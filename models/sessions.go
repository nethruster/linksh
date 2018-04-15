package models

import "time"

type Session struct {
	Id string `gorm:"primary_key; type:char(36)" json:"id"`
	UserId string
	LastUsedAt time.Time `json:"lastUsedAt"`
	CreatedAt time.Time `json:"createdAt"`
}


