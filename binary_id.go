package rdb

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

type BinID uuid.UUID

// StringToBinID parses string to BinID
func StringToBinID(s string) (BinID, error) {
	id, err := uuid.Parse(s)
	return BinID(id), err
}

// String String Representation of Binary16
func (b BinID) String() string {
	return uuid.UUID(b).String()
}

// GormDataType sets type to binary(16)
func (b BinID) GormDataType() string {
	return "binary(16)"
}

// MarshalBinary converts BinId to binary bytes
func (b BinID) MarshalBinary() ([]byte, error) {
	s := uuid.UUID(b)
	return s.MarshalBinary()
}

// MarshalJSON converts BinId to binary bytes
func (b BinID) MarshalJSON() ([]byte, error) {
	s := uuid.UUID(b)
	str := "\"" + s.String() + "\""
	return []byte(str), nil
}

// UnmarshalJSON sets the object from json bytes
func (b *BinID) UnmarshalJSON(by []byte) error {
	s, err := uuid.ParseBytes(by)
	*b = BinID(s)
	return err
}

// Scan  tells GORM how to receive from the database
func (b *BinID) Scan(value interface{}) error {

	bytes, _ := value.([]byte)
	parseByte, err := uuid.FromBytes(bytes)
	*b = BinID(parseByte)
	return err
}

// Value tells GORM how to save into the database
func (b BinID) Value() (driver.Value, error) {

	if b == BinID(uuid.Nil) {
		b = BinID(uuid.New())
	}
	return uuid.UUID(b).MarshalBinary()
}
