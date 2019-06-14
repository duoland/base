package fields

import (
	"reflect"
	"strings"
)

// CollectNonEmptyFields collects the field names of a struct which are not holding none empty values
func CollectNoneEmptyFields(v interface{}) (fields []string) {
	fields = make([]string, 0)
	val := reflect.ValueOf(v)
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldName := val.Type().Field(i).Name
		// if fieldName is id or updated_at, update it always
		fieldNameLower := strings.ToLower(fieldName)
		if fieldNameLower == "id" {
			continue
		} else if fieldNameLower == "updated_at" || fieldNameLower == "updatedat" {
			fields = append(fields, fieldName)
			continue
		}

		fieldKind := fieldVal.Kind()
		switch fieldKind {
		case reflect.Struct, reflect.Bool:
			{
				// IGNORE
			}
		case reflect.Ptr:
			if !fieldVal.IsNil() {
				fields = append(fields, fieldName)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if fieldVal.Int() != 0 {
				fields = append(fields, fieldName)
			}
		case reflect.Float32, reflect.Float64:
			if fieldVal.Float() != 0 {
				fields = append(fields, fieldName)
			}
		case reflect.String:
			if fieldVal.String() != "" {
				fields = append(fields, fieldName)
			}
		}
	}
	return
}

// TrimFieldSpace trims the space in the field values
func TrimFieldSpace(v interface{}) {
	val := reflect.Indirect(reflect.ValueOf(v))

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldKind := fieldVal.Kind()
		switch fieldKind {
		case reflect.String:
			fieldVal.SetString(strings.TrimSpace(fieldVal.String()))
		}
	}

	return
}
