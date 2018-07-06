package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/matoous/go-nanoid"
)

type Session struct {
	Id         string    `gorm:"primary_key; type:char(36)" json:"id"`
	UserId     string    `gorm:"type:char(36)"`
	ExpiresAt  time.Time `json:"expiresAt"`
	LastUsedAt time.Time `json:"lastUsedAt"`
	CreatedAt  time.Time `json:"createdAt"`
}

func CreateSession(db *gorm.DB, owner User) (Session, error) {
	id, err := GenerateSessionId()
	if err != nil {
		return Session{}, err
	}
	session := Session{
		Id:         id,
		UserId:     owner.Id,
		LastUsedAt: time.Now(),
	}
	err = db.Create(&session).Error
	if err != nil {
		return session, err
	}

	return session, nil
}
func GenerateSessionId() (string, error) {
	return gonanoid.Nanoid(36)
}
func UpdateSessionLastUsed(db *gorm.DB, session Session) error {
	return db.Model(&session).Update("LastUsedAt", time.Now()).Error
}
