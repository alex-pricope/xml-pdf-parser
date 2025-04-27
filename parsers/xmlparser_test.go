package parsers

import (
	"encoding/xml"
	myerrors "github.com/alex-pricope/form-parser/errors"
	"github.com/alex-pricope/form-parser/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestXMLParser_Parse_HappyPath(t *testing.T) {
	// Arrange
	parser := &XMLParser{}
	content, err := os.ReadFile("../tests/payload/valid_xml_tag")
	require.NoErrorf(t, err, "expected no error reading file, got %v", err)

	// Act
	rootNode, err := parser.Parse(content)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, rootNode)
	validateContentNode(t, rootNode)
}

func validateContentNode(t *testing.T, root *models.ContentNode) {
	// Validate root form node
	assert.Equal(t, models.FormElementType, root.ElementType)
	assert.NotNil(t, root.Children)
	assert.Len(t, root.Children, 2)

	// Validate field node
	fieldNode := root.Children[0]
	assert.Equal(t, models.FieldElementType, fieldNode.ElementType)
	assert.Equal(t, "program_language", fieldNode.Metadata["Name"])
	assert.Equal(t, "Enumeration(A,B,C)", fieldNode.Metadata["Type"])
	assert.Equal(t, "False", fieldNode.Metadata["Optional"])
	assert.Equal(t, "Select", fieldNode.Metadata["FieldType"])
	assert.NotNil(t, fieldNode.Children)
	assert.Len(t, fieldNode.Children, 2)

	// Validate caption under field
	captionNode := fieldNode.Children[0]
	assert.Equal(t, models.CaptionElementType, captionNode.ElementType)
	assert.Empty(t, captionNode.Metadata)
	assert.Equal(t, "Pick your programing language", captionNode.Value)
	assert.Nil(t, captionNode.Children)

	// Validate labels under field
	labelsNode := fieldNode.Children[1]
	assert.Equal(t, models.LabelsElementType, labelsNode.ElementType)
	assert.Empty(t, labelsNode.Metadata)
	assert.NotNil(t, labelsNode.Children)
	assert.Len(t, labelsNode.Children, 3)

	// Validate each label under labels
	label1 := labelsNode.Children[0]
	assert.Equal(t, models.LabelElementType, label1.ElementType)
	assert.Equal(t, "A", label1.Metadata["Name"])
	assert.Equal(t, "A(+)", label1.Value)
	assert.Nil(t, label1.Children)

	label2 := labelsNode.Children[1]
	assert.Equal(t, models.LabelElementType, label2.ElementType)
	assert.Equal(t, "B", label2.Metadata["Name"])
	assert.Equal(t, "B", label2.Value)
	assert.Nil(t, label2.Children)

	label3 := labelsNode.Children[2]
	assert.Equal(t, models.LabelElementType, label3.ElementType)
	assert.Equal(t, "C", label3.Metadata["Name"])
	assert.Equal(t, "C (All flavors except C#)", label3.Value)
	assert.Nil(t, label3.Children)

	// Validate section node
	sectionNode := root.Children[1]
	assert.Equal(t, models.SectionElementType, sectionNode.ElementType)
	assert.Equal(t, "experience", sectionNode.Metadata["Name"])
	assert.Equal(t, "False", sectionNode.Metadata["Optional"])
	assert.NotNil(t, sectionNode.Children)
	assert.Len(t, sectionNode.Children, 2)

	// Validate title under section
	titleNode := sectionNode.Children[0]
	assert.Equal(t, models.TitleElementType, titleNode.ElementType)
	assert.Empty(t, titleNode.Metadata)
	assert.Equal(t, "Regarding your experience", titleNode.Value)
	assert.Nil(t, titleNode.Children)

	// Validate contents under section
	contentsNode := sectionNode.Children[1]
	assert.Equal(t, models.ContentsElementType, contentsNode.ElementType)
	assert.Empty(t, contentsNode.Metadata)
	assert.NotNil(t, contentsNode.Children)
	assert.Len(t, contentsNode.Children, 2)

	// Validate first field under contents
	contentField1 := contentsNode.Children[0]
	assert.Equal(t, models.FieldElementType, contentField1.ElementType)
	assert.Equal(t, "other", contentField1.Metadata["Name"])
	assert.Equal(t, "Text([0,200],Lines:4)", contentField1.Metadata["Type"])
	assert.Equal(t, "True", contentField1.Metadata["Optional"])
	assert.Equal(t, "TextBox", contentField1.Metadata["FieldType"])
	assert.NotNil(t, contentField1.Children)
	assert.Len(t, contentField1.Children, 1)

	// Validate caption under first field
	contentCaption1 := contentField1.Children[0]
	assert.Equal(t, models.CaptionElementType, contentCaption1.ElementType)
	assert.Empty(t, contentCaption1.Metadata)
	assert.Equal(t, "Other programming experiences", contentCaption1.Value)
	assert.Nil(t, contentCaption1.Children)

	// Validate second field under contents
	contentField2 := contentsNode.Children[1]
	assert.Equal(t, models.FieldElementType, contentField2.ElementType)
	assert.Equal(t, "code_repos", contentField2.Metadata["Name"])
	assert.Equal(t, "File", contentField2.Metadata["Type"])
	assert.Equal(t, "True", contentField2.Metadata["Optional"])
	assert.Equal(t, "File", contentField2.Metadata["FieldType"])
	assert.NotNil(t, contentField2.Children)
	assert.Len(t, contentField2.Children, 1)

	// Validate caption under second field
	contentCaption2 := contentField2.Children[0]
	assert.Equal(t, models.CaptionElementType, contentCaption2.ElementType)
	assert.Empty(t, contentCaption2.Metadata)
	assert.Equal(t, "Upload your code repo's in ZIP.", contentCaption2.Value)
	assert.Nil(t, contentCaption2.Children)
}

func TestExtractMetadata(t *testing.T) {
	// Arrange
	parser := &XMLParser{}

	attributes := []xml.Attr{
		{Name: xml.Name{Local: "Name"}, Value: "program_language"},
		{Name: xml.Name{Local: "Type"}, Value: "Enumeration(A,B,C)"},
		{Name: xml.Name{Local: "Optional"}, Value: "False"},
		{Name: xml.Name{Local: "FieldType"}, Value: "Select"},
	}

	// Act
	result := parser.extractMetadata(attributes)

	// Assert
	assert.Len(t, result, 4)
	assert.Equal(t, "program_language", result["Name"])
	assert.Equal(t, "Enumeration(A,B,C)", result["Type"])
	assert.Equal(t, "False", result["Optional"])
	assert.Equal(t, "Select", result["FieldType"])
}

func TestXMLParser_Parse_EmptyContent(t *testing.T) {
	// Arrange
	parser := &XMLParser{}
	var emptyContent []byte

	// Act
	rootNode, err := parser.Parse(emptyContent)

	// Assert
	require.Error(t, err)
	assert.Nil(t, rootNode)
	assert.ErrorIs(t, err, myerrors.ErrEmptyFile)
}
