package rdb

import (
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/theNullP0inter/googly/logger"
	"github.com/theNullP0inter/googly/resource"
	"gorm.io/gorm"
)

func TestHandleGormError(t *testing.T) {

	assert.Equal(
		t,
		resource.ErrResourceNotFound,
		handleGormError(gorm.ErrRecordNotFound),
	)

	assert.Equal(
		t,
		resource.ErrInvalidTransaction,
		handleGormError(gorm.ErrInvalidTransaction),
	)

	assert.Equal(
		t,
		resource.ErrUniqueConstraint,
		handleGormError(&mysql.MySQLError{Number: 1062, Message: "test unique constraint"}),
	)

}

//  Test Create
func TestBaseRdbResourceManagerCreate(t *testing.T) {

	rdb, mock := GetMockRdb(t)

	l := logger.NewGooglyLogger()
	r := &MockModel{}
	qb := new(MockRdbListQueryBuilder)
	rm := NewRdbResourceManager(rdb, l, r, qb)

	// Test Create
	create_request := &MockModel{ID: BinID(uuid.New()), Name: "mock"}
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `mock_models`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	_, err := rm.Create(create_request)
	if err != nil {
		t.Errorf("Error creating a new row %s", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		fmt.Println(rdb.Statement.Table)
		t.Errorf("there were unfulfilled expectations at create: %s", err)
	}

	// Test Create with error

	create_request = &MockModel{ID: BinID(uuid.New()), Name: "mock"}
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `mock_models`").WillReturnError(errors.New("mock error"))
	mock.ExpectRollback()
	_, err = rm.Create(create_request)
	if err == nil {
		t.Error("Expected Error, received nil")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		fmt.Println(rdb.Statement.Table)
		t.Errorf("there were unfulfilled expectations at create: %s", err)
	}

}

// Test Get
func TestBaseRdbResourceManagerGet(t *testing.T) {

	rdb, mock := GetMockRdb(t)

	l := logger.NewGooglyLogger()
	r := &MockModel{}
	qb := new(MockRdbListQueryBuilder)
	rm := NewRdbResourceManager(rdb, l, r, qb)

	// Test Get With Invalid BinId Format
	id := ""
	_, err := rm.Get(id)
	if err == nil {
		t.Errorf("Expected Error when passed invalid BinId: %s", id)
	}

	// Test Get
	id = uuid.New().String()
	bin_id, _ := StringToBinID(id)

	mock_query_response := mock.NewRows([]string{"id", "name"})
	mock_query_response.AddRow(bin_id, "mock name")
	mock.ExpectQuery("SELECT").WillReturnRows(mock_query_response)

	_, err = rm.Get(id)

	if err != nil {
		t.Errorf("Error Getting Resource %s", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		fmt.Println(rdb.Statement.Table)
		t.Errorf("there were unfulfilled expectations at Get: %s", err)
	}

	// Test Get with Error
	mock.ExpectQuery("SELECT").WillReturnError(errors.New("mock error"))

	_, err = rm.Get(id)

	if err == nil {
		t.Error("Expected Error for the query. Received nil")
	}

}

// Test Delete
func TestBaseRdbResourceManagerDelete(t *testing.T) {

	rdb, mock := GetMockRdb(t)

	l := logger.NewGooglyLogger()
	r := &MockModel{}
	qb := new(MockRdbListQueryBuilder)
	rm := NewRdbResourceManager(rdb, l, r, qb)

	// Test Delete With Invalid BinId Format
	id := ""
	err := rm.Delete(id)
	if err == nil {
		t.Errorf("Expected Error when passed invalid BinId: %s", id)
	}

	// Test Delete
	id = uuid.New().String()
	bin_id, _ := StringToBinID(id)

	mock_query_response := mock.NewRows([]string{"id", "name"})
	mock_query_response.AddRow(bin_id, "mock name")
	mock.ExpectQuery("SELECT").WillReturnRows(mock_query_response)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(bin_id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = rm.Delete(id)
	if err != nil {
		t.Errorf("Error Deleting row %s", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		fmt.Println(rdb.Statement.Table)
		t.Errorf("there were unfulfilled expectations at Delete: %s", err)
	}

}

//  Test Update
func TestBaseRdbResourceManagerUpdate(t *testing.T) {

	rdb, mock := GetMockRdb(t)

	l := logger.NewGooglyLogger()
	r := &MockModel{}
	qb := new(MockRdbListQueryBuilder)
	rm := NewRdbResourceManager(rdb, l, r, qb)

	id := uuid.New().String()

	// Test Update
	update_request := &MockModel{Name: "new mock"}
	mock.ExpectBegin()
	mock.ExpectCommit()
	err := rm.Update(id, update_request)
	if err != nil {
		t.Errorf("Error Updating %s", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		fmt.Println(rdb.Statement.Table)
		t.Errorf("there were unfulfilled expectations at update: %s", err)
	}
}

// Test List
func TestBaseRdbResourceManagerList(t *testing.T) {
	rdb, _ := GetMockRdb(t)

	l := logger.NewGooglyLogger()
	r := &MockModel{}
	qb := new(MockRdbListQueryBuilder)
	rm := NewRdbResourceManager(rdb, l, r, qb)

	// Test Get With Invalid BinId Format
	id := ""
	_, err := rm.Get(id)
	if err == nil {
		t.Errorf("Expected Error when passed invalid BinId: %s", id)
	}

	// Test Get
	bin_id_1, _ := StringToBinID(uuid.New().String())
	bin_id_2, _ := StringToBinID(uuid.New().String())

	rdb, db_mock := GetMockRdb(t)
	qb.Rdb = rdb
	qb.On("ListQuery", mock.Anything).Return(mock.Anything, nil)

	mock_query_response := db_mock.NewRows([]string{"id", "name"})
	mock_query_response.AddRow(bin_id_1, "mock name")
	mock_query_response.AddRow(bin_id_2, "mock name")
	db_mock.ExpectQuery("SELECT").WillReturnRows(mock_query_response)

	_, err = rm.List(nil)

	if err != nil {
		t.Errorf("Error Getting List of Resources %s", err)
	}

	if err = db_mock.ExpectationsWereMet(); err != nil {
		fmt.Println(rdb.Statement.Table)
		t.Errorf("there were unfulfilled expectations at List: %s", err)
	}
}
