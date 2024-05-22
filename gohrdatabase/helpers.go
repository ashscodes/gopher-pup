package main

import (
	"fmt"
	"reflect"
	"strings"
)

func (p Person) Clone() *Person {
	clone := p
	return &clone
}

func (p People) ConvertToInterface() []interface{} {
	v := reflect.ValueOf(p)
	if v.Kind() != reflect.Slice {
		panic("input is not a slice")
	}

	result := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i] = v.Index(i).Interface()
	}
	return result
}

func (p People) ConvertToSlice() []*Person {
	var result []*Person
	for i := range p {
		result = append(result, &p[i])
	}

	return result
}

func ConvertToSlice(personMap map[string]Person) []*Person {
	peopleSlice := make([]*Person, 0, len(personMap))
	for _, person := range personMap {
		peopleSlice = append(peopleSlice, &person)
	}

	return peopleSlice
}

// SetFieldByReflection sets a field value of a struct using reflection.
func SetFieldByReflection(item interface{}, fieldPath string, value interface{}) error {
	fields := strings.Split(fieldPath, ".")
	field, err := getField(item, fields)
	if err != nil {
		return err
	}

	return setFieldValue(field, value)
}

func getField(obj interface{}, fields []string) (reflect.Value, error) {
	v := reflect.ValueOf(obj).Elem()

	for _, field := range fields {
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		v = v.FieldByName(field)
		if !v.IsValid() {
			return reflect.Value{}, fmt.Errorf("no such field: %s in obj", field)
		}
	}

	return v, nil
}

func setFieldValue(field reflect.Value, value interface{}) error {
	if !field.CanSet() {
		return fmt.Errorf("cannot set field")
	}

	val := reflect.ValueOf(value)
	if field.Type() != val.Type() {
		return fmt.Errorf("provided value type didn't match obj field type")
	}

	field.Set(val)
	return nil
}
