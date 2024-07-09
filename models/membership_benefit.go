package model

import (
	"time"
)

const TableNameMembershipBenefit = "tk.membership_benefit"

type MembershipBenefit struct {
	ID                int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	MembershipId      int32     `gorm:"column:membership_id;default:null" json:"membershipId"`
	Type              string    `gorm:"column:type;default:null" json:"type"`
	RefTable          string    `gorm:"column:ref_table;default:null" json:"refTable"`
	RefId             int64     `gorm:"column:ref_id;default:null" json:"refId"`
	Value             float64   `gorm:"column:value;default:null" json:"value"`
	Quota             float64   `gorm:"column:quota;default:null" json:"quota"`
	IsPercentageValue int16     `gorm:"column:is_percentage_value;default:0" json:"isPercentageValue"`
	CreatedAt         time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

func (*MembershipBenefit) TableName() string {
	return TableNameMembershipBenefit
}
