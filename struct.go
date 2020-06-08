package utils

import (
	"reflect"

	"github.com/fatih/structs"
)

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func Structs2Maps(list []interface{}) []map[string]interface{} {
	ret := make([]map[string]interface{}, 0)
	for _, item := range list {
		ret = append(ret, structs.Map(item))
	}
	return ret
}

func Struct2MapOmitEmpty(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj).Elem()
	v := reflect.ValueOf(obj).Elem()

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		if !v.Field(i).IsZero() {
			if t.Field(i).Name == "ID" {
				data["_id"] = v.Field(i).Interface()
			} else {
				data[t.Field(i).Name] = v.Field(i).Interface()
			}

		}

	}
	return data
}