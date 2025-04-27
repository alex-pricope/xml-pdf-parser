package models

import "strings"

type FileType string

// Using a custom enum to keep the file types centrally
const (
	XMLFileType  FileType = "xml"
	JSonFileType FileType = "json"
	PDFFileType  FileType = "pdf"
	HTMLFileType FileType = "html"

	UnknownFileType FileType = "unknown"
)

// SafeReadFileFormat - read the file type in a safe way to avoid a panic.
func SafeReadFileFormat(name string) FileType {
	switch strings.ToLower(name) {
	case "xml":
		return XMLFileType
	case "json":
		return JSonFileType
	case "pdf":
		return PDFFileType
	case "html":
		return HTMLFileType

	default:
		return UnknownFileType
	}
}
