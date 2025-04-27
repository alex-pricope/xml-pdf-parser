package models

import "strings"

type FieldType string

const (
	SelectFieldType  FieldType = "select"
	TextboxFieldType FieldType = "textbox"
	FileFieldType    FieldType = "file"

	UnknownFieldType FieldType = "unknown"
)

// SafeReadFieldType - read the field type in a safe way to avoid panics.
func SafeReadFieldType(name string) FieldType {
	switch strings.ToLower(name) {
	case "select":
		return SelectFieldType
	case "textbox":
		return TextboxFieldType
	case "file":
		return FileFieldType

	default:
		return UnknownFieldType
	}
}
