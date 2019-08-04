package model

import (
	"database/sql/driver"
	"fmt"
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
	if string(id) == "" {
		return int64(0), nil
	}
	i, err := strconv.Atoi(string(id))
	if err != nil {
		return nil, fmt.Errorf("Unable to convert %v of %T to int", id, id)
	}
	return int64(i), nil
}

// Scan is custom type for sql for support ID as int
func (id *ID) Scan(value interface{}) error {
	valueT, ok := value.(int64)
	if !ok {
		return fmt.Errorf("Unable to convert %v of %T to int64", value, value)
	}
	*id = ID(strconv.Itoa(int(valueT)))
	return nil
}
