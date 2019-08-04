package model

import (
	"database/sql/driver"

	"github.com/lib/pq"
)

// StringArray represents a one-dimensional array of the PostgreSQL character types and bson.Array of the MonogDB.
type StringArray pq.StringArray

func (a StringArray) Value() (driver.Value, error) {
	return pq.StringArray(a).Value()
}
