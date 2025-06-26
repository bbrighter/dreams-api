package repository

import (
	"fmt"
	"reflect"
	"strings"
)

func validateTypes(model any, updates map[string]any) error {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}
	for key, value := range updates {
		field, ok := modelType.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, key)
		})
		if !ok {
			return fmt.Errorf("invalid field %s", key)
		}
		expectedType := field.Type
		if expectedType.Kind() == reflect.Pointer {
			expectedType = expectedType.Elem()
		}

		if value == nil {
			continue
		}

		valueType := reflect.TypeOf(value)
		if valueType.Kind() == reflect.Pointer {
			valueType = valueType.Elem()
		}

		if !valueType.AssignableTo(expectedType) && !valueType.ConvertibleTo(expectedType) {
			return fmt.Errorf("invalid type for field %s, expected %s and got %s", key, expectedType, valueType)
		}
	}
	return nil
}
