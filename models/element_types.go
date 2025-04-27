package models

import "strings"

type ElementType string

const (
	FormElementType     ElementType = "form"
	FieldElementType    ElementType = "field"
	CaptionElementType  ElementType = "caption"
	LabelsElementType   ElementType = "labels"
	LabelElementType    ElementType = "label"
	SectionElementType  ElementType = "section"
	TitleElementType    ElementType = "title"
	ContentsElementType ElementType = "contents"

	UnknownElementType ElementType = "unknown"
)

// SafeReadElementType - read the element type in a safe way to avoid panics.
func SafeReadElementType(name string) ElementType {
	switch strings.ToLower(name) {
	case "form":
		return FormElementType
	case "field":
		return FieldElementType
	case "section":
		return SectionElementType
	case "caption":
		return CaptionElementType
	case "label":
		return LabelElementType
	case "title":
		return TitleElementType
	case "contents":
		return ContentsElementType
	case "labels":
		return LabelsElementType
	default:
		return UnknownElementType
	}
}
