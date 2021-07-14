package util

import (
	"ch2es/log"
	"fmt"
	"reflect"
)

func PrintFromDesc(pref string, s interface{}) {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	for i := 0; i < v.NumField(); i++ {
		p, ok := t.Field(i).Tag.Lookup("desc")
		if !ok {
			continue
		}
		field := v.Field(i).Interface()
		if field == "" {
			continue
		}

		if t.Field(i).Type.Kind() == reflect.Ptr {
			continue
		}

		log.Info(fmt.Sprint(pref, p, "=", field))
	}
}
