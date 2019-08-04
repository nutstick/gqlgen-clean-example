package model

import (
	"time"
)

// Admin is the structure representing a admin.
type Admin struct {
	tableName struct{}    `sql:"admins"`
	ID        ID          `json:"id" gorm:"type:serial" bson:"_id"`
	Email     string      `json:"email" sql:",unique" bson:"email"`
	Password  string      `json:"-" bson:"password"`
	Name      string      `json:"name" bson:"name"`
	Avatar    *string     `json:"avatar" bson:"avatar"`
	Roles     StringArray `json:"roles" sql:",array" gorm:"type:varchar(64)[]" bson:"roles"`
	CreateAt  time.Time   `json:"createAt" sql:"timestamp" bson:"createAt"`
	UpdateAt  time.Time   `json:"updateAt" sql:"timestamp" bson:"updateAt"`
}

// IsNode is Node type interface method
func (Admin) IsNode() {}
