package models

import (
	"github.com/jinzhu/gorm"
)

//DB is a wrapper for gorm, which include some usfeful methods
type DB struct {
	*gorm.DB
}

//OpenDatabase returns a new DB instance connected to the specified data source
func OpenDatabase(dataSourceName string) (*DB, error) {
	db, err := gorm.Open(dataSourceName)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
