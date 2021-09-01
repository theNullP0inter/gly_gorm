package rdb

import (
	"github.com/stretchr/testify/mock"
	"github.com/theNullP0inter/googly/resource"
	"gorm.io/gorm"
)

type MockRdbListQueryBuilder struct {
	mock.Mock
	Rdb *gorm.DB
}

func (b *MockRdbListQueryBuilder) ListQuery(resource.ListQuery) (*gorm.DB, error) {
	b.Called()
	return b.Rdb, nil
}
