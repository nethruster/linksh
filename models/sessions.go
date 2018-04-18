package models

import (
	"time"
	"github.com/matoous/go-nanoid"
)

type Session struct {
	Id         string    `gorm:"primary_key; type:char(36)" json:"id"`
	UserId     string
	ExpiresAt  time.Time `json:"lastUsedAt"`
	LastUsedAt time.Time `json:"lastUsedAt"`
	CreatedAt  time.Time `json:"createdAt"`
}

func GenerateSessionId() (string, error) {
	return gonanoid.Nanoid(36)
}