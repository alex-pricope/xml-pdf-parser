package render

import (
	"fmt"
	"github.com/alex-pricope/form-parser/models"
)

// Renderer - generic interface that different renderers implement
type Renderer interface {
	// Render - Renders the file and submission data to target
	Render(content *models.ContentNode, submission *models.ContentSubmission) error
}

// GetRenderer - Factory method that creates the renderer based on file type
func GetRenderer(fileType models.FileType, fileName string, dir string) (Renderer, error) {
	switch fileType {
	case models.PDFFileType:
		return NewPDFRenderer(fileName, dir), nil

	default:
		return nil, fmt.Errorf("unimplemented renderer type: %s", fileType)
	}
}
