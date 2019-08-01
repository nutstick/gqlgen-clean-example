package model

import "time"

// Admin is the structure representing a admin.
type Admin struct {
	ID       string    `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"-" bson:"password"`
	Name     string    `json:"name"`
	Avatar   *string   `json:"avatar"`
	Roles    []string  `json:"roles"`
	CreateAt time.Time `json:"createAt"`
	UpdateAt time.Time `json:"updateAt"`
}

// IsNode is Node type interface method
func (Admin) IsNode() {}
