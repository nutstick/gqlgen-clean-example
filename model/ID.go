package model

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// ID is type alias for cross database ID, for File or SQL will be string
// for mongodb will be bson.ObjectId
type ID string

// MarshalBSONValue is custom bson serialize function for support ID as ObjectID
func (id ID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	b, err := hex.DecodeString(string(id))
	if err != nil {
		return bsontype.ObjectID, bsoncore.AppendObjectID(nil, primitive.NilObjectID), err
	}
	if len(b) != 12 {
		return bsontype.ObjectID, bsoncore.AppendObjectID(nil, primitive.NilObjectID), primitive.ErrInvalidHex
	}
	// Enforce `byte` to 12 bytes type
	var oid [12]byte
	copy(oid[:], b[:])
	return bsontype.ObjectID, bsoncore.AppendObjectID(nil, oid), err
}

// UnmarshalBSONValue implements the bsoncodec.ValueUnmarshaler interface.
func (id *ID) UnmarshalBSONValue(t bsontype.Type, val []byte) error {
	// Enforce `byte` to 12 bytes type
	var oid [12]byte
	copy(oid[:], val[:])

	*id = ID(primitive.ObjectID(oid).Hex())
	return nil
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
