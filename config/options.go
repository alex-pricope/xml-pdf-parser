package config

import "github.com/alex-pricope/form-parser/models"

type CommandOptions struct {
	Filename           string
	SubmissionFileName string
	OutputDir          string

	FromType models.FileType
	ToType   models.FileType
}
