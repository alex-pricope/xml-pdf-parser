package reader

import (
	"encoding/json"
	myerrors "github.com/alex-pricope/form-parser/errors"
	"github.com/alex-pricope/form-parser/models"
	"os"
)

type Reader interface {
	ReadBinary(fileName string) ([]byte, error)
	ReadSubmissionFile(fileName string) (*models.ContentSubmission, error)
}

type FileReader struct {
}

func (r *FileReader) ReadBinary(fileName string) ([]byte, error) {
	if fileName == "" {
		return nil, myerrors.ErrEmptyPathProvided
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (r *FileReader) ReadSubmissionFile(fileName string) (*models.ContentSubmission, error) {
	if fileName == "" {
		return nil, myerrors.ErrEmptyPathProvided
	}

	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var submission models.ContentSubmission
	err = json.Unmarshal(file, &submission)
	if err != nil {
		return nil, err
	}

	return &submission, nil
}
