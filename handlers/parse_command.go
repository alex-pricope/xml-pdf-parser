package handlers

import (
	"errors"
	"fmt"
	"github.com/alex-pricope/form-parser/config"
	"github.com/alex-pricope/form-parser/logging"
	"github.com/alex-pricope/form-parser/models"
	"github.com/alex-pricope/form-parser/parsers"
	"github.com/alex-pricope/form-parser/reader"
	"github.com/alex-pricope/form-parser/render"
)

type CommandHandler interface {
	Handle() error
}

type ParseFormCommandHandler struct {
	Config   *config.CommandOptions
	Reader   reader.Reader
	Parser   parsers.Parser
	Renderer render.Renderer
}

func NewParseFormCommandHandler(reader reader.Reader, parser parsers.Parser, renderer render.Renderer, config *config.CommandOptions) *ParseFormCommandHandler {
	return &ParseFormCommandHandler{
		Config:   config,
		Reader:   reader,
		Parser:   parser,
		Renderer: renderer,
	}
}

func (r *ParseFormCommandHandler) Handle() error {
	// Read both files
	fileContent, submission, err := r.readFiles(r.Config.Filename, r.Config.SubmissionFileName)
	if err != nil {
		return err
	}

	if len(fileContent) == 0 {
		logging.Log.Errorf("Found empty (non-existing) file: %v", err)

	}
	// Parse the file and submission
	parsedFile, err := r.Parser.Parse(fileContent)
	if err != nil {
		logging.Log.Errorf("Error parsing file: %v", err)
		return err
	}

	// Render to target directory
	err = r.Renderer.Render(parsedFile, submission)
	if err != nil {
		logging.Log.Errorf("Error rendering file %s to %s: %v", r.Config.Filename, r.Config.ToType, err)
		return err
	}

	return nil
}

// readFiles - Read the content of the input file and the submission data
func (r *ParseFormCommandHandler) readFiles(fileName, submissionFileName string) ([]byte, *models.ContentSubmission, error) {
	fileContent, err := r.Reader.ReadBinary(fileName)
	if err != nil {
		logging.Log.Errorf("error reading file: %v", err)
		return nil, nil, err
	}

	if len(fileContent) == 0 {
		var message = fmt.Sprintf("file %s is empty", fileName)
		logging.Log.Error(message)
		return nil, nil, errors.New(message)
	}

	submission, err := r.Reader.ReadSubmissionFile(submissionFileName)
	if err != nil {
		logging.Log.Errorf("error reading file: %v", err)
		return nil, nil, err
	}

	if submission == nil {
		var message = fmt.Sprintf("submission file %s is empty", fileName)
		logging.Log.Error(message)
		return nil, nil, errors.New(message)
	}

	return fileContent, submission, nil
}
