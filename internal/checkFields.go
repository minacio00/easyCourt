package internal

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

func CheckMissingFields(data interface{}) error {
	v := reflect.ValueOf(data).Elem()
	typeOfData := v.Type()

	var missingFields []string
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := typeOfData.Field(i).Name
		zeroValue := reflect.Zero(field.Type()).Interface()

		if reflect.DeepEqual(field.Interface(), zeroValue) {
			fieldName = strcase.ToLowerCamel(fieldName)
			missingFields = append(missingFields, fieldName)
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ", "))
	}

	return nil
}
