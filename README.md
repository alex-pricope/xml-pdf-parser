[![Build Status](https://github.com/alex-pricope/form-parser/actions/workflows/ci.yml/badge.svg)](https://github.com/alex-pricope/form-parser/actions/workflows/ci.yml)

## Problem
We need to parse a given `input file` (XML initially), together with user submission data, and generate an `output file` (PDF) for now. 

The input XML is:

``` XML
<Form>
  <Field Name="program_language" Type="Enumeration(A,B,C)" Optional="False" FieldType="Select">
    <Caption>Pick your programing language</Caption>
    <Labels>
       <Label Name="A">A(+)</Label>
       <Label Name="B">B</Label>
       <Label Name="C">C (All flavors except C#)</Label>
    </Labels>
  </Field>
  <Section Name="experience" Optional="False">
    <Title>Regarding your experience</Title>
    <Contents>
      <Field Name="other" Type="Text([0,200],Lines:4)" Optional="True" FieldType="TextBox">
        <Caption>Other programming experiences</Caption>
      </Field>
      <Field Name="code_repos" Type="File" Optional="True" FieldType="File">
        <Caption>Upload your code repo's in ZIP.</Caption>
      </Field>
    </Contents>
  </Section>
</Form>
```
The result should be a PDF file rendering the components and user submission data.

## What I did
### Command Utility
I decided to create a utility `command` using Cobra library. This small program can run from anywhere amd has self explaining flags.

![image](https://github.com/user-attachments/assets/25064638-8d3c-4a55-8421-6e670753abb2)


### Usage:
* Step 1: Build the binary (run the below from repo root) - this will build the `parser` binary in `/bin` folder
  * > make all
* Step 2: Run the binary (the command has a nice help output to help with the arguments)
  * > ./parser -f=some_file.xml -s=submission --from=xml --to=pdf
* Step 3: Hopefully the command ran OK and it created the PDF. 

### Arguments:
* `-f, --file`: input file for the parser.
* `-s, --sub`: submission file (`JSON`) - this contains the user submitted data (more explanation below).
* `--from`: **from** **file type** - tells the utility what is the type of the input file. 
* `--to`: **to** **file type** - tells the utility what is the type of the output file.
*  `-o, --out`: optional **output** folder - if unspecified will use the current folder.

### Design
#### Generic components
I wanted to have `extensibility, simplicity and testability` so, for this, I used 3 major components:
* **_Reader_** (interface) - Reads files
* **_Parser_** (interface) + **_Factory pattern_** - This allows me to have multiple Parsers, for this exercise I only did the `XMLParser`
* **_Renderer_** (interface) + **_Factory pattern_** - Same as above, we can have multiple Renderers, I only implemented the `PDFRenderer`

These 3 components are used in a simple **_ParseFormCommandHandler_**, we can have many other commands. 

Because the utility takes in FROM and TO (file types), we can create the correct components where needed. 

#### Graph structure in the parser
The complex part of this, is to support a _dynamic structure_, where fields and sections can be mixed and generate content.
For this, I used a `graph structure` inside the **_Parser_**, with lots of documentation comments on how it works.
The basic idea is to build the content graph with parents and children, so we can traverse later. 

The structure is simple:
* `ElementType` (enum) is extracted from the XML data: `form`, `field`, `caption`, `labels`, `label`, etc.
* `Metadata`: `[Name="program_language", Type="Enumeration(A,B,C)", Optional="False", FieldType="Select"]`
* `Value`: the value of an element, if present - e.g. _Pick your programing language_ inside `Caption`
* `Name`: The name of the element used later to link user submission data to the actual item - e.g: _Name="program_language"_ (should be unique)
* `Children`: The collection of children

``` golang
type ContentNode struct {
	ElementType ElementType
	Metadata    map[string]string
	Value       string
	Name        string
	Children    []*ContentNode
}
```

On top of this, I used a `Stack` approach to traverse the XML since that is one of the best way to do this. 

#### User submission file
I did not know how to deal with this since the XML does not have the user submission inside. That's why I decided to have a separate JSON file 
that contains this needed data. 

The structure is simple, it will have the `Name` of the element (Caption) and the selected `values`. 
Example:
``` json
{
    "program_language": "B",
    "other": "Rust, Python, C++",
    "code_repos": "repo.zip"
}
```

#### Comments in the code
I wanted to have the code explained so that's why I left a lot of comments inside. 

#### Testing
To keep the time for this exercise at a normal level, I **did not include unit tests for the PDFRenderer** (and maybe some other types)

But I did include
* Unit tests for most of the components (usually same folder files with same name but _test)
* Integration tests 
  * I also included a more complex scenario with a different XML structure

#### Makefile
I love makefiles and I included one for this as well. This allows anyone tio quickly build the binary, while running tests, linting, etc 

Targets
* all - does everything `[lint test build]`
* individual `build`, `lint`, `test`
  * NOTE: `lint` has a dependency to `golangci-lint` - it will output this and the website to download it from 

#### Assumptions
* User submission file - like I explained above
* I did not do a dropdown list in the output PDF - that would complicate the renderer more
* I forced the order on the `labels` - map in go does not guarantee item insertion ordering 
* I did not inject the logger into the components to keep them on the lighter side
* I did not write all the tests, it would take a lot of time
* The items on the XML have unique names - the only way to get to them with separate user submission file

#### Libraries 
* [Cobra](https://github.com/spf13/cobra) - for building great CLI
* [Logrus](https://github.com/sirupsen/logrus) - great logging library
* [Testify](https://github.com/stretchr/testify) - for great assertions
* [gofpdf](https://github.com/jung-kurt/gofpdf) - picked this one to render the PDF, archived and not maintained anymore but good for this exercise

#### Github Actions pipeline
I also included a `GHA pipeline` that executes `make all` after every push. 
