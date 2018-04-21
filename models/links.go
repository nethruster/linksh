package models

import (
    "time"
    "github.com/jinzhu/gorm"
    "github.com/matoous/go-nanoid"
)

type Link struct {
	Id      string `gorm:"primary_key"`
	Content string
    Hits    uint   `gorm:"DEFAULT:0"`
    UserId  string `gorm:"type:char(36)"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func CreateLink(db *gorm.DB, ownerId, customId, content string) (Link, error) {
    var id string
    if customId == "" {
        gid, err := GenerateLinkId()
        if err != nil {
            return Link{}, err
        }
        id = gid
    } else {
        id = ownerId
    }

    link := Link{
        Id:      id,
        Content: content,
        UserId:  ownerId,
    }

    err := db.Create(&link).Error
    return link, err
}
func GenerateLinkId() (string, error) {
    return gonanoid.Nanoid(6)
}
