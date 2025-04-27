package render

import (
	"fmt"
	"github.com/alex-pricope/form-parser/logging"
	"github.com/alex-pricope/form-parser/models"
	"github.com/jung-kurt/gofpdf"
	"path/filepath"
	"sort"
	"strings"
)

var orientation, unit, size, font = "P", "mm", "A4", "Arial"
var defaultFontSize float64 = 12
var titleFontSize float64 = 14
var missingCaptionTextValue = "(missing caption)"
var selectedMarkerValue = "(selected)"
var missingAnswerTextValue = "(missing answer)"

type PDFRenderer struct {
	Filename string
	Dir      string

	pdf *gofpdf.Fpdf
}

func NewPDFRenderer(fileName, dir string) *PDFRenderer {
	return &PDFRenderer{
		Filename: fileName,
		Dir:      dir,
	}
}

func (r *PDFRenderer) Render(content *models.ContentNode, submission *models.ContentSubmission) error {
	r.pdf = gofpdf.New(orientation, unit, size, "")
	r.useNormalFont(defaultFontSize)
	r.pdf.AddPage()

	// Traverse the graph and render the needed elements
	err := r.renderNode(content, submission)
	if err != nil {
		return err
	}

	// Write the PDF file
	err = r.writeFile()
	if err != nil {
		return err
	}

	return nil
}

// renderNode - renders a content node
func (r *PDFRenderer) renderNode(node *models.ContentNode, submission *models.ContentSubmission) error {
	switch node.ElementType {

	case models.SectionElementType:
		r.renderSection(node, submission)

	case models.FieldElementType:
		r.renderField(node, submission)

	case models.FormElementType:
		logging.Log.Info("(skip)Form content node type")

	case models.CaptionElementType, models.LabelsElementType, models.LabelElementType, models.TitleElementType, models.ContentsElementType:
		// No-op here - these are handled inside the parents

	default:
		logging.Log.Warnf("(skip)Unknown content node type: %s", node.ElementType)
	}

	for _, child := range node.Children {
		err := r.renderNode(child, submission)
		if err != nil {
			return err
		}
	}

	return nil
}

// writeFile - write the PDF file based on Dir or Filename path
func (r *PDFRenderer) writeFile() error {
	// If Dir is set, use the dir/filename.pdf
	// If not, use the filename_path/filename.pdf
	basePdfName := strings.TrimSuffix(filepath.Base(r.Filename), filepath.Ext(r.Filename)) + ".pdf"
	var outputPath string
	if r.Dir != "" {
		outputPath = filepath.Join(r.Dir, basePdfName)
	} else {
		outputPath = filepath.Join(filepath.Dir(r.Filename), basePdfName)
	}
	err := r.pdf.OutputFileAndClose(outputPath)
	if err != nil {
		logging.Log.Errorf("Error writing file: %v", err)
		return err
	}

	return nil
}

// findCaption - search the nodes for the Caption.
func (r *PDFRenderer) findCaption(node *models.ContentNode) string {
	for _, child := range node.Children {
		if child.ElementType == models.CaptionElementType {
			return child.Value
		}
	}
	return missingCaptionTextValue
}

// findTitle - search the nodes for the Title
func (r *PDFRenderer) findTitle(node *models.ContentNode) string {
	for _, child := range node.Children {
		if child.ElementType == models.TitleElementType {
			return child.Value
		}
	}
	return ""
}

// renderField - generic method that will render the field
func (r *PDFRenderer) renderField(node *models.ContentNode, submission *models.ContentSubmission) {

	// Read the FieldType - decided not to transform this in the parser to keep it simple
	fieldTypeStr, ok := node.Metadata["FieldType"]
	var fieldType models.FieldType
	if !ok {
		fieldType = models.UnknownFieldType
	} else {
		fieldType = models.SafeReadFieldType(fieldTypeStr)
	}

	switch fieldType {

	// For simplicity, File will render like a normal textbox
	case models.FileFieldType, models.TextboxFieldType:
		r.renderTextBoxFieldType(node, submission)

	case models.SelectFieldType:
		r.renderSelectFieldType(node, submission)

	// For Unknown, just skip and log
	case models.UnknownFieldType:
		logging.Log.Warnf("(skip)Unknown field type: %s", fieldType)
		return
	}
}

// renderSection - renders a Section. E.g. <section> ... </field>
func (r *PDFRenderer) renderSection(node *models.ContentNode, submission *models.ContentSubmission) {
	// Render title if any in a bigger and bolder text
	r.renderTitle(node)
}

// renderSelectFieldType - renders a Select FieldType. E.g. <field FieldType="Select"> ... </field>
func (r *PDFRenderer) renderSelectFieldType(node *models.ContentNode, submission *models.ContentSubmission) {
	// Step 1: Find the Caption if exists
	caption := r.findCaption(node)

	// Step 2: Find the submitted value
	selectedValue := getSubmittedValue(submission, node.Name)

	// Step 3: Render the Caption and submitted answer
	line := fmt.Sprintf("%s: %s", caption, selectedValue)
	r.useBoldFont(defaultFontSize)
	r.writeCellLn(10, 10, line)

	r.useNormalFont(defaultFontSize)

	// Step 4: Find all Labels and construct the dropdown - also select the submitted value
	validLabels := make(map[string]string)

	for _, child := range node.Children {
		if child.ElementType == models.LabelsElementType {
			for _, labelNode := range child.Children {
				if labelNode.ElementType == models.LabelElementType {
					labelName, ok := labelNode.Metadata["Name"]
					if !ok {
						logging.Log.Warnf("(skip)Label node without Name metadata: %+v", labelNode)
						continue
					}
					labelText := labelNode.Value
					validLabels[labelName] = labelText
					continue
				}
				logging.Log.Warnf("(skip)Unknown label node type: %s", child.ElementType)
			}
		}
	}

	// Step 5: Check if submitted value matches any label
	if _, ok := validLabels[selectedValue]; !ok && selectedValue != "" {
		logging.Log.Warnf("Submitted value '%s' for field '%s' not found in labels", selectedValue, node.Name)
	}

	// Step 6: Render all labels, marking the selected one with bold
	// Apply a forced sorting
	labelNames := make([]string, 0, len(validLabels))
	for labelName := range validLabels {
		labelNames = append(labelNames, labelName)
	}
	sort.Strings(labelNames)

	for _, labelName := range labelNames {
		labelText := validLabels[labelName]

		selectMarker := ""
		if labelName == selectedValue {
			selectMarker = selectedMarkerValue
		}

		optionLine := fmt.Sprintf("- %s %s", labelText, selectMarker)

		if labelName == selectedValue {
			// Selected option
			r.useHighlightColor()
			r.useBoldFont(defaultFontSize)

			r.pdf.CellFormat(0, 8, optionLine, "", 1, "", true, 0, "")

			// Reset the styling to default
			r.useNormalFont(defaultFontSize)
			r.useNormalColor()
		} else {
			// Normal option - nothing special
			r.writeCellLn(8, 8, optionLine)
		}
	}
}

// renderTextBoxFieldType - renders a Textbox FieldType. E.g. <field FieldType="TextBox"> ... </field>
func (r *PDFRenderer) renderTextBoxFieldType(node *models.ContentNode, submission *models.ContentSubmission) {
	// Step 1: Find the Caption if exists
	caption := r.findCaption(node)

	// Step 2: Find the submitted value
	submittedValue := getSubmittedValue(submission, node.Name)

	// If missing, insert placeholder text
	if submittedValue == "" {
		submittedValue = missingAnswerTextValue
	}

	// Step 3: Render the Caption and the value
	r.useBoldFont(defaultFontSize)
	r.writeCellLn(10, 10, caption)

	r.useNormalFont(defaultFontSize)
	r.useHighlightColor()
	r.pdf.MultiCell(0, 8, submittedValue, "", "", true)
	r.useNormalColor()
	r.pdf.Ln(5)
}

func (r *PDFRenderer) renderTitle(node *models.ContentNode) {
	title := r.findTitle(node)

	r.useBoldFont(titleFontSize)
	r.writeCellLn(10, 12, title)
	r.useNormalFont(defaultFontSize)
}

func getSubmittedValue(submission *models.ContentSubmission, fieldName string) string {
	if submission == nil {
		return missingAnswerTextValue
	}
	return (*submission)[fieldName]
}

func (r *PDFRenderer) useNormalColor() {
	r.pdf.SetFillColor(255, 255, 255)
}

func (r *PDFRenderer) useHighlightColor() {
	r.pdf.SetFillColor(220, 220, 220)
}

func (r *PDFRenderer) useNormalFont(size float64) {
	r.pdf.SetFont(font, "", size)
}

func (r *PDFRenderer) useBoldFont(size float64) {
	r.pdf.SetFont(font, "B", size)
}

func (r *PDFRenderer) writeCellLn(hLn, hCell float64, text string) {
	r.pdf.Cell(0, hCell, text)
	r.pdf.Ln(hLn)
}
