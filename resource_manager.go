package rdb

import (
	"errors"
	"reflect"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/copier"
	"github.com/theNullP0inter/googly/db"
	"github.com/theNullP0inter/googly/logger"
	"github.com/theNullP0inter/googly/resource"
	"gorm.io/gorm"
)

type RdbResourceManager interface {
	resource.DbResourceManager
}

type BaseRdbResourceManager struct {
	*resource.BaseResourceManager
	Db               *gorm.DB
	Model            db.BaseModel
	ListQueryBuilder RdbListQueryBuilder
}

// handleGormError converts gorm Errors to ResourceErrors
func handleGormError(err error) error {

	if err == gorm.ErrRecordNotFound {
		return resource.ErrResourceNotFound
	} else if err == gorm.ErrInvalidTransaction {
		return resource.ErrInvalidTransaction
	}

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return resource.ErrUniqueConstraint
	}
	return resource.ErrInvalidQuery
}

// GetModel will get you the model resource
func (s *BaseRdbResourceManager) GetResource() resource.Resource {
	return s.Model
}

// GetModel will get you the model resource
func (s *BaseRdbResourceManager) GetModel() resource.Resource {
	return s.Model
}

// Create creates an entry in with given data
func (s *BaseRdbResourceManager) Create(m resource.DataInterface) (resource.DataInterface, error) {
	item := reflect.New(reflect.TypeOf(s.GetModel())).Interface()
	copier.Copy(item, m)
	result := s.Db.Create(item)

	if result.Error != nil {
		return nil, handleGormError(result.Error)
	}
	return item, nil
}

// Get gets 1 item with given id
func (s *BaseRdbResourceManager) Get(id resource.DataInterface) (resource.DataInterface, error) {
	strId := id.(string)
	item := reflect.New(reflect.TypeOf(s.GetModel())).Interface()
	binId, err := StringToBinID(strId)
	if err != nil {
		return nil, resource.ErrInvalidFormat
	}
	bId, _ := binId.MarshalBinary()
	err = s.Db.Where("id = ?", bId).First(item).Error
	if err != nil {
		return nil, handleGormError(err)
	}
	return item, nil
}

// Update updates 1 item with given id & given data/update_set
func (s *BaseRdbResourceManager) Update(id resource.DataInterface, data resource.DataInterface) error {

	m := reflect.New(reflect.TypeOf(s.GetModel())).Interface()
	strId := id.(string)
	binId, err := StringToBinID(strId)
	if err != nil {
		return resource.ErrInvalidFormat
	}
	bId, _ := binId.MarshalBinary()

	copier.Copy(m, data)

	result := s.Db.Model(s.GetModel()).Where("id = ?", bId).Updates(m)

	if result.Error != nil {
		return handleGormError(result.Error)
	}

	return nil
}

// Delete will delete 1 item with given _id
func (s *BaseRdbResourceManager) Delete(id resource.DataInterface) error {
	item, err := s.Get(id)
	if err != nil {
		return err
	}
	s.Db.Delete(item)
	return nil
}

// List will get you a list of items
//
// it uses QueryBuilder.ListQuery() to filter the throuh rows
func (s *BaseRdbResourceManager) List(parameters resource.DataInterface) (resource.DataInterface, error) {

	items := reflect.New(reflect.SliceOf(reflect.TypeOf(s.GetModel()))).Interface()
	result, err := s.ListQueryBuilder.ListQuery(parameters)
	if err != nil {
		return nil, resource.ErrInternal
	}
	result = result.Find(items)
	if result.Error != nil {
		return nil, handleGormError(result.Error)
	}

	return items, nil
}

// NewRdbResourceManager creates a new RdbResourceManager
func NewRdbResourceManager(
	db *gorm.DB,
	logger logger.GooglyLoggerInterface,
	model db.BaseModel,
	queryBuilder RdbListQueryBuilder,
) *BaseRdbResourceManager {
	resourceManager := resource.NewBaseResourceManager(logger, model)
	return &BaseRdbResourceManager{
		BaseResourceManager: resourceManager,
		Db:                  db,
		Model:               model,
		ListQueryBuilder:    queryBuilder,
	}

}
