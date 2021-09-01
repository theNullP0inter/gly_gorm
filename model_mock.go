package rdb

type MockModel struct {
	ID   BinID `gorm:"primaryKey" json:"id"`
	Name string
}
