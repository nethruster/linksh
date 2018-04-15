package models

import "time"

type Link struct {
	Id string `gorm:"primary_key"`
	Content string
	Hits uint
	UserId string

	CreatedAt time.Time
	UpdatedAt time.Time
}