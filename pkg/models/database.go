package models

import (
	"github.com/jinzhu/gorm"
)

//A wrapper for gorm, which include some usfeful methods
type DB struct {
    *gorm.DB
}

func OpenDatabase(dataSourceName string) (*DB, error) {
    db, err := gorm.Open(dataSourceName)
    if err != nil {
        return nil, err
    }

    return &DB{db}, nil
}
