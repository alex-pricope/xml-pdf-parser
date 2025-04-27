package render

import (
	"testing"

	"github.com/alex-pricope/form-parser/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRenderer_PDFFileType(t *testing.T) {
	// Arrange
	renderer, err := GetRenderer(models.PDFFileType, "output.pdf", "output")
	require.NoError(t, err)
	assert.NotNil(t, renderer)

	// Act
	_, ok := renderer.(*PDFRenderer)

	// Assert
	assert.True(t, ok, "expected type *PDFRenderer")
}

func TestGetRenderer_UnknownFileType(t *testing.T) {
	// Arrange Act Assert
	renderer, err := GetRenderer(models.UnknownFileType, "output.pdf", "output")
	require.Error(t, err)
	assert.Nil(t, renderer)
	assert.Contains(t, err.Error(), "unimplemented renderer type")
}
