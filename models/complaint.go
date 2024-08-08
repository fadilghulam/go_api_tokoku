package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

const TableNameComplaints = "tk.complaints"

type StringArray []string

// Value converts StringArray to a PostgreSQL array-compatible format.
func (a StringArray) Value() (driver.Value, error) {
	return fmt.Sprintf("{%s}", strings.Join(a, ",")), nil
}

// Scan converts a PostgreSQL array to StringArray.
func (a *StringArray) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		trimmed := strings.Trim(v, "{}")
		if len(trimmed) == 0 {
			*a = []string{}
			return nil
		}
		*a = strings.Split(trimmed, ",")
	case []byte:
		trimmed := strings.Trim(string(v), "{}")
		if len(trimmed) == 0 {
			*a = []string{}
			return nil
		}
		*a = strings.Split(trimmed, ",")
	default:
		return fmt.Errorf("unsupported data type: %T", v)
	}

	return nil
}

type Complaints struct {
	ID                int64       `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerID        int64       `gorm:"column:customer_id;default:null" json:"customerId"`
	TransactionId     int64       `gorm:"column:transaction_id;default:null" json:"transactionId"`
	TransactionItemId int32       `gorm:"column:transaction_item_id;default:null" json:"transactionItemId"`
	Description       string      `gorm:"column:description;default:null" json:"description"`
	Image             StringArray `gorm:"type:varchar[];column:image;default:null" json:"image"`
	CreatedAt         time.Time   `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt         time.Time   `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	Other             *string     `gorm:"column:other;default:null" json:"other"`
}

func (*Complaints) TableName() string {
	return TableNameComplaints
}
