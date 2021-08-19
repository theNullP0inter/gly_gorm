package rdb

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel for rdb models
type BaseModel struct {
	ID        BinID     `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// BaseModel for rdb models with soft delete
//
// Refer: https://gorm.io/docs/delete.html#Soft-Delete
type RdbSoftDeleteBaseModel struct {
	BaseModel
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Setting uuid before creation
//
// Refer: https://gorm.io/docs/hooks.html#Creating-an-object
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	b.ID = BinID(id)
	return err
}
