package models

import "time"

type Link struct {
	Id      string `gorm:"primary_key"`
	Content string
	Hits    uint
    UserId  string `gorm:"type:char(36)"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
