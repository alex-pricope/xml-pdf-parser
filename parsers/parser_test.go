package parsers

import (
	"testing"

	"github.com/alex-pricope/form-parser/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetParser_XMLFileType(t *testing.T) {
	// Arrange
	parser, err := GetParser(models.XMLFileType)
	require.NoError(t, err)
	assert.NotNil(t, parser)

	// Act
	_, ok := parser.(*XMLParser)

	// Assert
	assert.True(t, ok, "expected type *XMLParser")
}

func TestGetParser_UnknownFileType(t *testing.T) {
	// Arrange Act Assert
	parser, err := GetParser(models.UnknownFileType)
	require.Error(t, err)
	assert.Nil(t, parser)
	assert.Contains(t, err.Error(), "unimplemented parser type")
}
