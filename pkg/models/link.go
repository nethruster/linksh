package models

import (
	"time"

	"github.com/matoous/go-nanoid"
)

//Link type
type Link struct {
	ID      string `gorm:"primary_key" json:"id"`
	Content string `json:"content"`
	Hits    uint   `gorm:"DEFAULT:0"`
	UserID  string `gorm:"type:char(36)"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

//GetLink return the requested link
func (db *DB) GetLink(id string) (Link, error) {
	var link Link
	err := db.Where("id = ?", id).Take(&link).Error

	return link, err
}

//GetLinks returns an array of links
//if an ownerID is provided, the result will be limited to the links which belongs to that user
//a '0' value in offset or limit will be ignored
func (db *DB) GetLinks(ownerID string, offset, limit int) ([]Link, error) {
	var links []Link
	query := db.DB

	if offset != 0 {
		query = query.Offset(offset)
	}
	if limit != 0 {
		query = query.Limit(limit)
	}
	if ownerID != "" {
		query = query.Where("user_id = ?", ownerID)
	}

	err := query.Find(&links).Error

	return links, err
}

//CreateLink creates a new link and stores in the database
func (db *DB) CreateLink(ownerID, customID, content string) (Link, error) {
	var id string
	if customID == "" {
		gid, err := GenerateLinkID()
		if err != nil {
			return Link{}, err
		}
		id = gid
	} else {
		id = customID
	}

	link := Link{
		ID:      id,
		Content: content,
		UserID:  ownerID,
	}

	err := db.Create(&link).Error
	if err != nil && customID == "" && err.Error() == "Duplicate entry" {
		return db.CreateLink(ownerID, customID, content)
	}

	return link, err
}

//UpdateLink modifies the content of the selected link
func (db *DB) UpdateLink(linkID string, content string) error {
	return db.Model(&Link{ID: linkID}).Update("content", content).Error
}

//DeleteLink remove the selected link from the database
func (db *DB) DeleteLink(linkID string) error {
	return db.Delete(&Link{ID: linkID}).Error
}

//GenerateLinkID retuns a randdom valid ID for a link
func GenerateLinkID() (string, error) {
	return gonanoid.Nanoid(6)
}
