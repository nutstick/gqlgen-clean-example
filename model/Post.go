package model

import "time"

// Post is the structure respresenting a post
type Post struct {
	tableName   struct{}  `sql:"posts"`
	ID          ID        `json:"id" gorm:"type:serial" bson:"_id"`
	Description string    `json:"description" bson:"description"`
	ImageURL    string    `json:"imageUrl" bson:"imageUrl"`
	CreateAt    time.Time `json:"createAt" sql:"timestamp" bson:"createAt"`
	UpdateAt    time.Time `json:"updateAt" sql:"timestamp" bson:"updateAt"`
}
