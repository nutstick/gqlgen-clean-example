package model

import "time"

// Admin is the structure representing a admin.
type Admin struct {
	tableName struct{}  `sql:"admins"`
	ID        ID        `json:"id" sql:",pk" bson:"_id"`
	Email     string    `json:"email" sql:",unique" bson:"email"`
	Password  string    `json:"-" bson:"password"`
	Name      string    `json:"name" bson:"name"`
	Avatar    *string   `json:"avatar" bson:"avatar"`
	Roles     []string  `json:"roles" bson:"roles"`
	CreateAt  time.Time `json:"createAt" bson:"createAt"`
	UpdateAt  time.Time `json:"updateAt" bson:"updateAt"`
}

// IsNode is Node type interface method
func (Admin) IsNode() {}
