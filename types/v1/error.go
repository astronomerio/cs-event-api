package v1

import "fmt"

// FieldError used to represent errors returned when a field was
// not initialized properly in a structure passed as argument to one of the
// functions of this package.
type FieldError struct {

	// The human-readable representation of the type of structure that wasn't
	// initialized properly.
	Type string

	// The name of the field that wasn't properly initialized.
	Name string

	// The value of the field that wasn't properly initialized.
	Value interface{}
}

func (e FieldError) Error() string {
	return fmt.Sprintf("%s.%s: invalid field value: %#v", e.Type, e.Name, e.Value)
}
