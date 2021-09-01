package rdb

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBaseModel(t *testing.T) {

	b := &BaseModel{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.Equal(t, BinID(uuid.Nil), b.ID)
	rdb, _ := GetMockRdb(t)
	b.BeforeCreate(rdb)

	assert.NotNil(t, b.ID)
}
