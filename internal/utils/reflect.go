package utils

import (
	"fmt"
	"reflect"
)

func UpdateStruct(src interface{}, dest interface{}) error {
	srcValue := reflect.ValueOf(src)
	destValue := reflect.ValueOf(dest)

	// Ensure src is a pointer to a struct
	if srcValue.Kind() != reflect.Ptr || srcValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("src must be a pointer to a struct")
	}

	// Ensure dest is a pointer to a struct
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dest must be a pointer to a struct")
	}

	// Dereference the pointers to access the struct values
	srcValue = srcValue.Elem()
	destValue = destValue.Elem()

	// Iterate over the fields of the src struct
	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		srcFieldName := srcValue.Type().Field(i).Name
		destField := destValue.FieldByName(srcFieldName)

		// Check if the field exists in dest and is settable
		if !destField.IsValid() || !destField.CanSet() {
			continue
		}

		// Check if the source field is non-zero
		if srcField.IsZero() {
			continue
		}

		// Ensure the field types match before setting
		if srcField.Type() == destField.Type() {
			destField.Set(srcField)
		}
	}

	return nil
}

func UpdateStructField(s interface{}, fieldName string, newValue interface{}) error {
	ref := reflect.ValueOf(s)

	// Ensure we have a pointer to a struct
	if ref.Kind() != reflect.Ptr || ref.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected a pointer to a struct")
	}

	// Get the struct value
	ref = ref.Elem()

	// Get the field by name
	field := ref.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("no such field: %s in struct", fieldName)
	}

	// Check if the field is settable
	if !field.CanSet() {
		return fmt.Errorf("cannot set field: %s", fieldName)
	}

	// Set the field value (type must match)
	fieldValue := reflect.ValueOf(newValue)
	if field.Type() != fieldValue.Type() {
		return fmt.Errorf("provided value type doesn't match field type")
	}

	field.Set(fieldValue)
	return nil
}
