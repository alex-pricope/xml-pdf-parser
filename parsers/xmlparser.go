package parsers

import (
	"bytes"
	"encoding/xml"
	myerrors "github.com/alex-pricope/form-parser/errors"
	"github.com/alex-pricope/form-parser/logging"
	"github.com/alex-pricope/form-parser/models"
	"io"
	"strings"
)

type XMLParser struct{}

func (p *XMLParser) Parse(content []byte) (*models.ContentNode, error) {
	if len(content) == 0 {
		return nil, myerrors.ErrEmptyFile
	}

	root, err := p.parseXMLContent(content)
	if err != nil {
		logging.Log.Errorf("XMLParser parse file error: %s", err)
		return nil, err
	}

	return root, nil
}

func (p *XMLParser) parseXMLContent(content []byte) (*models.ContentNode, error) {
	/* Because I want to read and parse the contents while I go through the file content, I need to find the children.
	 A good way to solve this is using a stack.

	 Given a stack [ ]
	 And a file like below
	 <form>
	    <section>
	         <field></field>
	    </section>
	</form>

	I want to start parsing the contents and I use the tokens below.
	* I find first token (open) -> push to the stack [ form ]
	* I find the next token (open) -> push to the stack [ section form ]
	* I find the next (open) -> push again [ field section form ]
	* I find the next (closed) -> pop the stack -> attach to parent [ section form]
	* I find the next (closed) -> pop the stack -> attach to parent [ form ]
	* I find the next (closed) -> pop the stack -> root [ ]

	*/
	dec := xml.NewDecoder(bytes.NewReader(content))
	var root *models.ContentNode

	// Use an array like a stack
	var stack []*models.ContentNode

	for {
		xmlToken, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			// If XML is invalid, log and return error
			logging.Log.Errorf("XMLParser: error decoding XML token: %v", err)
			return nil, err
		}

		// Type switch to get the type of the token
		switch tType := xmlToken.(type) {
		case xml.StartElement:
			// If StartElement eg <form>, push to the stack, create a new node
			node := &models.ContentNode{
				ElementType: models.SafeReadElementType(tType.Name.Local),
				Metadata:    p.extractMetadata(tType.Attr),
			}

			if name, ok := node.Metadata["Name"]; ok {
				node.Name = name
			}

			// If this is the first node => root node, else attach it to children of the prev node (peek)
			if len(stack) == 0 {
				root = node
			} else {
				prev := stack[len(stack)-1]
				prev.Children = append(prev.Children, node)
			}
			stack = append(stack, node)
		case xml.CharData:
			// If its CharData => value for current node (peek)
			if len(stack) > 0 {
				current := stack[len(stack)-1]
				text := strings.TrimSpace(string(tType))
				if len(text) > 0 {
					current.Value += text
				}
			}
		case xml.EndElement:
			// If EndElement => pop the stack
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		}
	}

	return root, nil
}

func (p *XMLParser) extractMetadata(attributes []xml.Attr) map[string]string {
	result := make(map[string]string)

	for _, a := range attributes {
		result[a.Name.Local] = a.Value
	}

	return result
}
