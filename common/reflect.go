package common

import (
	"log"
	"reflect"
)

func PrintFromDesc(pref string, s interface{}) {
	v := reflect.ValueOf(s)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		p, ok := t.Field(i).Tag.Lookup("desc")
		if !ok {
			continue
		}
		field := v.Field(i).Interface()
		if field == "" {
			continue
		}
		if field == "" {
			continue
		}
		log.Println(pref, p, "=", field)
	}
}
