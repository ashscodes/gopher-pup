package main

import "reflect"

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
