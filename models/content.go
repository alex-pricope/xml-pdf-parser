package models

/* We cannot have a static model here. Each form is dynamic, it can have multiple elements
   But one thing is for sure: the file (when valid) will have different elements that have opening and closing tags
   * example: <Field> ... </Field>

   We can use a graph data structure to keep this.
   [node] <- root
		[node] [/node] <- child 1
		[node] <- child 2
             [node] [/node] <- child 1 of child 1
		[/node]
  [/node]

  Each node is a struct. We can hold the metadata, node type (label, section, etc.)
  At the end we can send this to the renderer.
  The metadata for each node can be: Name, Type, Optional, FieldType

 I decided to elevate Name to properties on the struct since this is used more often.
 The rest I kept in the metadata since only the renderer will use them.
*/

type ContentNode struct {
	ElementType ElementType
	Metadata    map[string]string
	Value       string
	Name        string
	Children    []*ContentNode
}

/* Since I am not sure how the submission values get here, and the XML does not have the user data,
   I decided to use a separate JSON file that will hold that.

   {
    "program_language": "B", <- (Name of field : Name of selection (not direct value))
    "other": "Rust, Python, C++", <- (Name of field : Value)
    "code_repos": "repo.zip" <- (Name of field : Value)
	}
*/

type ContentSubmission map[string]string
