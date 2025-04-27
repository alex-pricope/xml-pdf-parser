package cmd

import (
	"github.com/alex-pricope/form-parser/config"
	"github.com/alex-pricope/form-parser/handlers"
	"github.com/alex-pricope/form-parser/logging"
	"github.com/alex-pricope/form-parser/models"
	"github.com/alex-pricope/form-parser/parsers"
	"github.com/alex-pricope/form-parser/reader"
	"github.com/alex-pricope/form-parser/render"
	"github.com/spf13/cobra"
)

// ParseCommand will parse the file and generate the output
func ParseCommand(cmd *cobra.Command, _ []string) {
	// Read config from user input
	conf, err := readCommandOptions(cmd)
	if err != nil {
		logging.Log.Errorf("Error while reading command parameter: %v", err)
	}

	parse, err := parsers.GetParser(conf.FromType)
	if err != nil {
		logging.Log.Errorf("Error creating parser: %v", err)
		return
	}

	renderer, err := render.GetRenderer(conf.ToType, conf.Filename, conf.OutputDir)
	if err != nil {
		logging.Log.Errorf("Error creating renderer: %v", err)
		return
	}

	// Execute the handler
	handler := handlers.NewParseFormCommandHandler(&reader.FileReader{}, parse, renderer, conf)
	err = handler.Handle()
	if err != nil {
		logging.Log.Errorf("Error while parsing form command: %v", err)
		return
	}

	err = handler.Handle()
	if err != nil {
		logging.Log.Errorf("Error while executing command: %v", err)
		return
	}
}

// readCommandOptions - gather the inputs of the command
func readCommandOptions(cmd *cobra.Command) (*config.CommandOptions, error) {
	filePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return nil, err
	}

	submissionFilePath, err := cmd.Flags().GetString("sub")
	if err != nil {
		return nil, err
	}

	fromFormat, err := cmd.Flags().GetString("from")
	if err != nil {
		return nil, err
	}

	toFormat, err := cmd.Flags().GetString("to")
	if err != nil {
		return nil, err
	}

	outputFolder, err := cmd.Flags().GetString("out")
	if err != nil {
		return nil, err
	}

	return &config.CommandOptions{
		Filename:           filePath,
		SubmissionFileName: submissionFilePath,
		OutputDir:          outputFolder,
		FromType:           models.SafeReadFileFormat(fromFormat),
		ToType:             models.SafeReadFileFormat(toFormat),
	}, nil
}
