package model

import (
	"database/sql/driver"
	"strconv"

	"gopkg.in/mgo.v2/bson"
)

// ID is type alias for cross database ID, for File or SQL will be string
// for mongodb will be bson.ObjectId
type ID string

// GetBSON is custom bson serialize function for support ID as ObjectID
func (id ID) GetBSON() (interface{}, error) {
	return bson.ObjectIdHex(string(id)), nil
}

// SetBSON is custom bson serialize function for support ID as ObjectID
func (id *ID) SetBSON(raw bson.Raw) error {
	var decoded bson.ObjectId
	bsonErr := raw.Unmarshal(decoded)

	if bsonErr == nil {
		*id = ID(decoded.Hex())
		return nil
	}
	return bsonErr
}

// Value is custom type for sql for support ID as int
func (id ID) Value() (driver.Value, error) {
	return strconv.Atoi(string(id))
}

// Scan is custom type for sql for support ID as int
func (id *ID) Scan(value interface{}) error {
	valueT, _ := value.(int)
	*id = ID(valueT)
	return nil
}
