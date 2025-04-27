package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeReadFieldType(t *testing.T) {
	tests := []struct {
		input    string
		expected FieldType
	}{
		{"Select", SelectFieldType},
		{"TextBox", TextboxFieldType},
		{"File", FileFieldType},
		{"UnknownType", UnknownFieldType},
		{"", UnknownFieldType},
	}

	for _, tt := range tests {
		t.Run("FieldType_"+tt.input, func(t *testing.T) {
			result := SafeReadFieldType(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSafeReadElementType(t *testing.T) {
	tests := []struct {
		input    string
		expected ElementType
	}{
		{"Form", FormElementType},
		{"Field", FieldElementType},
		{"Section", SectionElementType},
		{"Caption", CaptionElementType},
		{"Labels", LabelsElementType},
		{"Label", LabelElementType},
		{"Title", TitleElementType},
		{"Contents", ContentsElementType},
		{"Unknown", UnknownElementType},
		{"", UnknownElementType},
	}

	for _, tt := range tests {
		t.Run("ElementType_"+tt.input, func(t *testing.T) {
			result := SafeReadElementType(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSafeReadFileFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected FileType
	}{
		{"PDF", PDFFileType},
		{"HTML", HTMLFileType},
		{"Unknown", UnknownFileType},
		{"", UnknownFileType},
	}

	for _, tt := range tests {
		t.Run("FileFormat_"+tt.input, func(t *testing.T) {
			result := SafeReadFileFormat(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
