package storage

import (
	"errors"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrFieldNotFound       = errors.New("field not found")
	ErrUnexpectedLength    = errors.New("unexpected length")
	ErrViewRefFieldMissing = errors.New("view ref field missing")
	ErrTableNotFoundInView = errors.New("table not found in view")
	ErrTypeAssert          = errors.New("type assert failed")
)

func ContainsField(fields []string, field string) bool {
	for _, f := range fields {
		if f == field {
			return true
		}
	}
	return false
}
