package generic_conv

import (
	"errors"
	"git.singularity-ai.com/backend/library/v2/xContext/loggers/xlog"
	"reflect"
)

func Conv[T any](s string, unmarshal func(string, interface{}) error) (newT T, err error) {
	var t T
	if reflect.TypeOf(t).Kind() == reflect.Uintptr || reflect.TypeOf(t).Kind() == reflect.Pointer {
		v := reflect.New(reflect.TypeOf(t).Elem())
		ptr, ok := v.Interface().(T)
		if !ok {
			err = errors.New("failed conv ptr")
			xlog.Error(err)
			return t, err
		}
		err = unmarshal(s, ptr)
		if err != nil {
			xlog.Error(err)
			return t, err
		}
		return ptr, err
	} else {
		v := reflect.New(reflect.TypeOf(t))
		val, ok := v.Interface().(*T)
		if !ok {
			err = errors.New("failed conv base type")
			xlog.Error(err)
			return t, err
		}
		err = unmarshal(s, val)
		if err != nil {
			xlog.Error(err)
			return t, err
		}
		return *val, err
	}
}

func DefaultValue[T any]() (newT T, err error) {
	var t T
	if reflect.TypeOf(t).Kind() == reflect.Uintptr || reflect.TypeOf(t).Kind() == reflect.Pointer {
		v := reflect.New(reflect.TypeOf(t).Elem())
		ptr, ok := v.Interface().(T)
		if !ok {
			err = errors.New("failed conv ptr")
			xlog.Error(err)
			return t, err
		}
		return ptr, err
	} else {
		v := reflect.New(reflect.TypeOf(t))
		val, ok := v.Interface().(*T)
		if !ok {
			err = errors.New("failed conv base type")
			xlog.Error(err)
			return t, err
		}
		return *val, err
	}
}
