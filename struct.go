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
