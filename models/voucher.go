package model

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const TableNameVoucher = "tk.voucher"

type Int32Array []int32

// Value converts Int32Array to a PostgreSQL array-compatible format.
func (a Int32Array) Value() (driver.Value, error) {
	// Convert []int32 to []interface{}
	var arr = make([]interface{}, len(a))
	for i, v := range a {
		arr[i] = v
	}
	return arr, nil
}

// Scan converts a PostgreSQL array to Int32Array.
func (a *Int32Array) Scan(value interface{}) error {
	var ints []int32

	switch v := value.(type) {
	case string:
		// Handle the case where the array is returned as a string
		trimmed := strings.Trim(v, "{}")
		if len(trimmed) == 0 {
			*a = []int32{}
			return nil
		}
		strElements := strings.Split(trimmed, ",")
		for _, strElem := range strElements {
			i, err := strconv.Atoi(strElem)
			if err != nil {
				return err
			}
			ints = append(ints, int32(i))
		}
	case []byte:
		// Handle the case where the array is returned as []byte
		trimmed := strings.Trim(string(v), "{}")
		if len(trimmed) == 0 {
			*a = []int32{}
			return nil
		}
		strElements := strings.Split(trimmed, ",")
		for _, strElem := range strElements {
			i, err := strconv.Atoi(strElem)
			if err != nil {
				return err
			}
			ints = append(ints, int32(i))
		}
	default:
		return fmt.Errorf("unsupported data type: %T", v)
	}

	*a = ints
	return nil
}

type Voucher struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	// ProdukId []int `gorm:"column:produk_id" json:"produk_id"`
	// ProdukId     datatypes.JSONType[int64] `gorm:"column:produk_id" json:"produk_id"`
	ProdukId     Int32Array `gorm:"column:produk_id;type:int[]" json:"produk_id"`
	DateStart    string     `gorm:"column:date_start;not null;default:now()" json:"date_start"`
	DateEnd      string     `gorm:"column:date_end;not null;default:now()" json:"date_end"`
	Code         string     `gorm:"column:code" json:"code"`
	Icon         string     `gorm:"column:icon" json:"icon"`
	Diskon       float64    `gorm:"column:diskon" json:"diskon"`
	IsPercentage int32      `gorm:"column:is_percentage" json:"is_percentage"`
	MinCost      float64    `gorm:"column:min_diskon" json:"min_cost"`
	MaxDiskon    float64    `gorm:"column:max_diskon" json:"max_diskon"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	Note         string     `gorm:"column:note" json:"note"`
	Amount       int16      `gorm:"column:amount" json:"amount"`
}

func (*Voucher) TableName() string {
	return TableNameVoucher
}
