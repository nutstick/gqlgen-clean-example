package utils

import (
	"reflect"
	"time"

	"github.com/imdario/mergo"
)

type timeTransfomer struct {
	overwrite bool
}

func (t timeTransfomer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ == reflect.TypeOf(time.Time{}) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				if t.overwrite {
					isZero := src.MethodByName("IsZero")
					result := isZero.Call([]reflect.Value{})
					if !result[0].Bool() {
						dst.Set(src)
					}
				} else {
					isZero := dst.MethodByName("IsZero")
					result := isZero.Call([]reflect.Value{})
					if result[0].Bool() {
						dst.Set(src)
					}
				}
			}
			return nil
		}
	}
	return nil
}

func Merge(dst, src interface{}) error {
	return mergo.Merge(dst, src, mergo.WithOverride, mergo.WithTransformers(timeTransfomer{true}))
}
