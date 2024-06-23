package util

import (
	"fmt"
	"reflect"
)

// NotNil panics if any of the arguments are nil.
func NotNil(objs ...interface{}) {
	if len(objs) <= 0 {
		panic("no arguments provided")
	}
	for _, obj := range objs {
		if obj == nil {
			panic(fmt.Sprintf("%s is nil", obj))
		}

		v := reflect.ValueOf(obj)
		if v.Kind() == reflect.Ptr && v.IsNil() {
			panic(fmt.Sprintf("%s is nil", obj))
		}
	}
}
