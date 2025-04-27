package reader

import (
	"github.com/alex-pricope/form-parser/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadSubmissionFile_HappyPath(t *testing.T) {
	// Arrange
	reader := &FileReader{}
	expected := &models.ContentSubmission{
		"program_language": "B",
		"other":            "Rust, Python, C++",
		"code_repos":       "repo.zip",
	}

	// Act
	result, err := reader.ReadSubmissionFile("../tests/payload/valid_submission")

	// Assert
	require.NoErrorf(t, err, "expected no error, got %v", err)
	require.NotNil(t, result)
	assert.Equal(t, expected, result)
}

func TestReadSubmissionFile_FileNotFound(t *testing.T) {
	// Arrange
	reader := &FileReader{}

	// Act
	_, err := reader.ReadSubmissionFile("non_existent_file.xml")

	// Assert
	require.Error(t, err)
}

func TestReadSubmissionFile_InvalidJson(t *testing.T) {
	// Arrange
	reader := &FileReader{}

	// Act
	_, err := reader.ReadSubmissionFile("../tests/payload/invalid_submission")

	// Assert
	require.Error(t, err)
}

func TestReadFileAsBinary_HappyPath(t *testing.T) {
	// Arrange
	reader := &FileReader{}

	// Act
	result, err := reader.ReadBinary("../tests/payload/valid_xml")

	// Assert
	require.NoErrorf(t, err, "expected no error, got %v", err)
	require.NotEmpty(t, result)
}

func TestReadFileAsBinary_FileNotFound(t *testing.T) {
	// Arrange
	reader := &FileReader{}

	// Act
	_, err := reader.ReadBinary("non_existent_file.xml")

	// Assert
	require.Error(t, err)
}

func TestReadFileAsBinary_FileIsEmpty(t *testing.T) {
	// Arrange
	reader := &FileReader{}

	// Act
	result, err := reader.ReadBinary("../tests/payload/empty")

	// Assert
	require.NoError(t, err)
	require.Empty(t, result)
	require.NotNil(t, result)
}
