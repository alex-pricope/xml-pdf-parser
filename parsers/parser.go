package parsers

import (
	"fmt"
	"github.com/alex-pricope/form-parser/models"
)

// Parser - generic interface that different parsers will implement
type Parser interface {
	//Parse - Parses the file content into a domain model
	Parse(fileContent []byte) (*models.ContentNode, error)
}

// GetParser - Factory method that creates the parser based on file type
func GetParser(fileType models.FileType) (Parser, error) {
	//The idea is to use this factory method to create different parsers. Right now there is only 1
	switch fileType {
	case models.XMLFileType:
		return &XMLParser{}, nil

	default:
		return nil, fmt.Errorf("unimplemented parser type: %s", fileType)
	}
}
