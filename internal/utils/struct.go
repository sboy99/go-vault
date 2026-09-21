package utils

import (
	"fmt"
	"reflect"
	"time"
)

// GetStructFields converts the fields of a struct to a string array
func GetStructFields(s interface{}) ([]string, error) {
	// Get the type of the struct //
	t := reflect.TypeOf(s)

	// Check if the type is a struct //
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("GetStructFields: %v is not a struct", t)
	}

	// Get the number of fields in the struct //
	numFields := t.NumField()

	// Create a slice to store the field names //
	fields := make([]string, 0, numFields)

	// Loop through the fields and get the names //
	for i := 0; i < numFields; i++ {
		field := t.Field(i)
		fields = append(fields, CamelCaseToTitleCase(field.Name))
	}

	return fields, nil
}

// GetStructValues converts the values of a struct to a string array
func GetStructValues(s interface{}) ([]string, error) {
	// Get the type and value of the struct
	val := reflect.ValueOf(s)

	// Check if the type is a struct //
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("GetStructFields: %v is not a struct", val)
	}

	// Iterate over the struct fields and collect their values as strings
	var result []string
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Handle field types (e.g., format time values)
		switch field.Kind() {
		case reflect.String:
			result = append(result, field.String())
		case reflect.Int, reflect.Int64:
			result = append(result, fmt.Sprintf("%d", field.Int()))
		case reflect.Float64, reflect.Float32:
			result = append(result, fmt.Sprintf("%f", field.Float()))
		case reflect.Bool:
			result = append(result, fmt.Sprintf("%t", field.Bool()))
		case reflect.Struct: // Handle specific struct types like `time.Time`
			if field.Type() == reflect.TypeOf(time.Time{}) {
				result = append(result, field.Interface().(time.Time).Format("2006-01-02 15:04:05"))
			} else {
				result = append(result, fmt.Sprintf("%v", field.Interface()))
			}
		default:
			result = append(result, fmt.Sprintf("%v", field.Interface()))
		}
	}

	return result, nil
}
