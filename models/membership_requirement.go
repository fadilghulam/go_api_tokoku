package model

import (
	"time"
)

const TableNameMembershipRequirement = "tk.membership_requirement"

type MembershipRequirement struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Type      string    `gorm:"column:type;default:null" json:"type"`
	MinNumber float64   `gorm:"column:min_number;default:null" json:"minNumber"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

func (*MembershipRequirement) TableName() string {
	return TableNameMembershipRequirement
}
