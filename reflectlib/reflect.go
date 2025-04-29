package reflectlib

import (
	"fmt"
	"github.com/dixonwhitmire/golib/errorlib"
	"reflect"
)

// FieldMetadata stores a struct field's Kind, ordinal position (index), and Name.
type FieldMetadata struct {
	Kind  reflect.Kind
	Index int
	Name  string
}

// ParseStructFields returns []FieldMetadata for each field within a struct or a struct type.
// Supported types include: struct, struct pointer, or a struct type object.
func ParseStructFields(s any) ([]FieldMetadata, error) {
	// ensure we have a reflect.Type
	t, ok := s.(reflect.Type)
	if !ok {
		t = reflect.TypeOf(s)
	}

	if t.Kind() != reflect.Struct && !(t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct) {
		return nil, errorlib.CreateError("ParseStructFields",
			fmt.Sprintf("Type [%T] is not a struct or struct pointer", s))
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fields := make([]FieldMetadata, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		fieldMetadata := FieldMetadata{
			Kind:  t.Field(i).Type.Kind(),
			Index: i,
			Name:  t.Field(i).Name,
		}
		fields = append(fields, fieldMetadata)
	}

	return fields, nil
}
