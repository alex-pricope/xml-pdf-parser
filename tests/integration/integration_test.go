package integration

import (
	"github.com/alex-pricope/form-parser/config"
	"github.com/alex-pricope/form-parser/handlers"
	"github.com/alex-pricope/form-parser/logging"
	"github.com/alex-pricope/form-parser/parsers"
	"github.com/alex-pricope/form-parser/reader"
	"github.com/alex-pricope/form-parser/render"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	logging.Log = logrus.New()
	logging.Log.SetLevel(logrus.FatalLevel)

	os.Exit(m.Run())
}

func TestParseXMLForm_CreatePDF(t *testing.T) {
	tests := []struct {
		name            string
		options         *config.CommandOptions
		expectedPDFPath string
	}{
		{
			name: "DirSpecified",
			options: &config.CommandOptions{
				Filename:           "../../tests/payload/valid_xml_tag",
				SubmissionFileName: "../../tests/payload/valid_submission",
				OutputDir:          "./out",
				FromType:           "xml",
				ToType:             "pdf",
			},
			expectedPDFPath: "./out/valid_xml_tag.pdf",
		},
		{
			name: "DirNotSpecified",
			options: &config.CommandOptions{
				Filename:           "../../tests/payload/valid_xml_tag",
				SubmissionFileName: "../../tests/payload/valid_submission",
				OutputDir:          "",
				FromType:           "xml",
				ToType:             "pdf",
			},
			expectedPDFPath: "../../tests/payload/valid_xml_tag.pdf",
		},
		{
			name: "DirSpecifiedStripExtension",
			options: &config.CommandOptions{
				Filename:           "../../tests/payload/valid_xml.xml",
				SubmissionFileName: "../../tests/payload/valid_submission",
				OutputDir:          "./out",
				FromType:           "xml",
				ToType:             "pdf",
			},
			expectedPDFPath: "./out/valid_xml.pdf",
		},
		{
			name: "ComplexStructureXML",
			options: &config.CommandOptions{
				Filename:           "../../tests/payload/complex_valid_xml",
				SubmissionFileName: "../../tests/payload/complex_valid_submission",
				OutputDir:          "./out",
				FromType:           "xml",
				ToType:             "pdf",
			},
			expectedPDFPath: "./out/valid_xml.pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			aReader := &reader.FileReader{}
			aParser, err := parsers.GetParser(tt.options.FromType)
			require.NoError(t, err)
			aRenderer, err := render.GetRenderer(tt.options.ToType, tt.options.Filename, tt.options.OutputDir)
			require.NoError(t, err)

			commandHandler := handlers.NewParseFormCommandHandler(aReader, aParser, aRenderer, tt.options)
			require.NotNil(t, commandHandler)

			// Act
			err = commandHandler.Handle()

			// Assert
			require.NoError(t, err)

			info, err := os.Stat(tt.expectedPDFPath)
			require.NoError(t, err)
			require.False(t, info.IsDir())
			require.Greater(t, info.Size(), int64(0))
		})
	}
}
