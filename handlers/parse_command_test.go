package handlers

import (
	"errors"
	"github.com/alex-pricope/form-parser/config"
	"github.com/alex-pricope/form-parser/logging"
	"github.com/alex-pricope/form-parser/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type fakeReader struct {
	fileError       error
	submissionError error
	fileContent     []byte
	submissionData  *models.ContentSubmission
}

func (r *fakeReader) ReadBinary(_ string) ([]byte, error) {
	return r.fileContent, r.fileError
}

func (r *fakeReader) ReadSubmissionFile(_ string) (*models.ContentSubmission, error) {
	return r.submissionData, r.submissionError
}

type fakeParser struct {
	parseError error
}

func (p *fakeParser) Parse(_ []byte) (*models.ContentNode, error) {
	if p.parseError != nil {
		return nil, p.parseError
	}
	return &models.ContentNode{}, nil
}

type fakeRenderer struct {
	renderError error
}

func (r *fakeRenderer) Render(_ *models.ContentNode, _ *models.ContentSubmission) error {
	return r.renderError
}

func TestMain(m *testing.M) {
	logging.Log = logrus.New()
	logging.Log.SetLevel(logrus.FatalLevel)

	os.Exit(m.Run())
}

func TestHandle_HappyPath(t *testing.T) {
	// Arrange
	handler := &ParseFormCommandHandler{
		Reader:   &fakeReader{fileContent: []byte("some xml"), submissionData: &models.ContentSubmission{}},
		Parser:   &fakeParser{},
		Renderer: &fakeRenderer{},
		Config: &config.CommandOptions{
			Filename:           "some xml",
			SubmissionFileName: "some name",
			OutputDir:          "",
			FromType:           "xml",
			ToType:             "pdf",
		},
	}

	// Act
	err := handler.Handle()

	// Assert
	require.NoError(t, err)
}

func TestHandle_ReadFilesError(t *testing.T) {
	// Arrange
	handler := &ParseFormCommandHandler{
		Reader:   &fakeReader{fileError: errors.New("read file error")},
		Parser:   &fakeParser{},
		Renderer: &fakeRenderer{},
		Config:   &config.CommandOptions{},
	}

	// Act
	err := handler.Handle()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read file error")
}

func TestHandle_ParseError(t *testing.T) {
	// Arrange
	handler := &ParseFormCommandHandler{
		Reader:   &fakeReader{fileContent: []byte("some xml"), submissionData: &models.ContentSubmission{}},
		Parser:   &fakeParser{parseError: errors.New("parse error")},
		Renderer: &fakeRenderer{},
		Config:   &config.CommandOptions{},
	}

	// Act
	err := handler.Handle()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parse error")
}

func TestHandle_RenderError(t *testing.T) {
	// Arrange
	handler := &ParseFormCommandHandler{
		Reader:   &fakeReader{fileContent: []byte("some xml"), submissionData: &models.ContentSubmission{}},
		Parser:   &fakeParser{},
		Renderer: &fakeRenderer{renderError: errors.New("render error")},
		Config:   &config.CommandOptions{},
	}

	// Act
	err := handler.Handle()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "render error")
}
