package model

import (
	"time"
)

const TableNameMembership = "tk.membership"

type Membership struct {
	ID                 int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name               string    `gorm:"column:name;default:null" json:"name"`
	Description        string    `gorm:"column:description;default:null" json:"description"`
	FormulaRequirement string    `gorm:"column:formula_requirement;default:null" json:"formulaRequirement"`
	FormulaBenefit     string    `gorm:"column:formula_benefit;default:null" json:"formulaBenefit"`
	Orders             int16     `gorm:"column:orders" json:"orders"`
	CreatedAt          time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

func (*Membership) TableName() string {
	return TableNameMembership
}
