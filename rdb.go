package rdb

import (
	"gorm.io/gorm"
)

// NewRdb creates a new instance of gorm.DB
//
// It takes dialect as input. So, it can support any db supported by gorm
func NewRdb(dialect gorm.Dialector) *gorm.DB {
	db, err := gorm.Open(dialect)
	if err != nil {
		panic("failed to connect to the database")
	}

	return db
}
