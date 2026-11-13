package f

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func J(v interface{}) string {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%+v", v)
	}

	return string(jsonBytes)
}

func JB(v interface{}) []byte {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return []byte(fmt.Sprintf("%+v", v))
	}

	return jsonBytes
}

type KeyValue struct {
	Name  string
	Value interface{}
}

func NonEmptyStructValues(data interface{}) []KeyValue {

	structIterator := reflect.ValueOf(data)
	if structIterator.Kind() != reflect.Struct {
		return nil
	}

	var keyValues []KeyValue
	for i := 0; i < structIterator.NumField(); i++ {
		field := structIterator.Type().Field(i).Name
		val := structIterator.Field(i).Interface()

		// Check if the field is zero-valued, meaning it won't be updated
		if !reflect.DeepEqual(val, reflect.Zero(structIterator.Field(i).Type()).Interface()) {
			// fmt.Printf("%v is non-zero, adding to update\n", field)
			keyValues = append(keyValues, KeyValue{
				Name:  field,
				Value: val,
			})
		}
	}

	return keyValues
}
