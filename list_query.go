package rdb

import (
	"fmt"
	"reflect"

	"github.com/theNullP0inter/googly/logger"
	"github.com/theNullP0inter/googly/resource"
	"gorm.io/gorm"
)

const DefaultPageSize = 30

// RdbListQueryBuilder converts resource.ListQuery to filters needed for gorm
// ListQuery() should facilitate this
type RdbListQueryBuilder interface {
	ListQuery(resource.ListQuery) (*gorm.DB, error)
}

// PaginatedRdbListQueryBuilder is a pagination implementation for RdbListQueryBuilder
type PaginatedRdbListQueryBuilder interface {
	RdbListQueryBuilder
	PaginationQuery(resource.ListQuery) *gorm.DB
}

// BasePaginatedRdbListQueryBuilder is a base implementation for PaginatedRdbListQueryBuilder
type BasePaginatedRdbListQueryBuilder struct {
	Rdb    *gorm.DB
	Logger logger.GooglyLoggerInterface
}

// ListQuery should be implemented for RdbListQueryBuilder.
func (qb *BasePaginatedRdbListQueryBuilder) ListQuery(parameters resource.ListQuery) (*gorm.DB, error) {
	return qb.PaginationQuery(parameters), nil
}

// PaginationQuery Converts resource.ListQuery to pagination filters needed for gorm
func (qb *BasePaginatedRdbListQueryBuilder) PaginationQuery(parameters resource.ListQuery) *gorm.DB {
	query := qb.Rdb

	// if params do not match the required format, log and return empty filter
	val := reflect.ValueOf(parameters).Elem()
	if val.Kind() != reflect.Struct {
		qb.Logger.Errorf("Unexpected type of parameters for PaginationQuery")
		return query
	}

	paginationParameters := val.FieldByName("PaginatedListQuery")
	hasPaginationParams := paginationParameters.IsValid() && !paginationParameters.IsNil()

	// Parsing page number
	var page int64
	page = 0
	if hasPaginationParams {
		pageValue := val.FieldByName("Page")
		if !pageValue.IsValid() || pageValue.Kind() != reflect.Int {
			qb.Logger.Info("Page in in invalid format, Using default value")
		} else {
			page = pageValue.Int()
		}
	}

	// Parsing page size
	var pageSize int64
	pageSize = DefaultPageSize
	if hasPaginationParams {
		pageSizeValue := val.FieldByName("PageSize")
		if !pageSizeValue.IsValid() || pageSizeValue.Kind() != reflect.Int {
			qb.Logger.Info("PageSize in in invalid format, Using default value")
		} else {
			pageSize = pageSizeValue.Int()
		}
	}

	// Applying Pagination to query
	limit := pageSize
	offset := page * pageSize
	query = query.Offset(int(offset)).Limit(int(limit))

	// Parsing orderBy
	var orderBy = ""
	if hasPaginationParams {
		orderByValue := val.FieldByName("OrderBy")
		if orderByValue.IsValid() && orderByValue.Kind() == reflect.String {
			orderBy = orderByValue.String()
		}
	}

	// Return if order by is not present
	if orderBy == "" {
		return query
	}

	// Parsing orderBy descending or ascending
	var orderDesc = false
	if hasPaginationParams {
		orderDescValue := val.FieldByName("OrderDesc")
		if orderDescValue.IsValid() && orderDescValue.Kind() == reflect.Bool {
			orderDesc = orderDescValue.Bool()
		}
	}

	// Adding sort to query
	if len(orderBy) > 0 {
		if orderDesc {
			query = query.Order(fmt.Sprintf("%s DESC", orderBy))
		} else {
			query = query.Order(fmt.Sprintf("%s ASC", orderBy))
		}
	}

	return query
}

// NewBasePaginatedRdbListQueryBuilder creates a new BasePaginatedRdbListQueryBuilder
func NewBasePaginatedRdbListQueryBuilder(db *gorm.DB, logger logger.GooglyLoggerInterface) *BasePaginatedRdbListQueryBuilder {
	return &BasePaginatedRdbListQueryBuilder{Rdb: db, Logger: logger}
}
