package csv2template

import (
	"bytes"
	"encoding/csv"
	"io"
	"strings"
	"text/template"
)

// Row is our representation of a Packer.Row
type Row struct {
	Columns []string
}

// AddColumns constructs a Row
func (r *Row) AddColumns(strings []string) {
	for _, v := range strings {
		r.Columns = append(r.Columns, v)
	}
}

// From the Packer docs, this represents:
// 1 index, 2 subtype, 3..n subtype data

// TemplatePage represents a page of rows and columns
type TemplatePage struct {
	Rows []Row
}

// AddRows constructs a TemplatePage
func (t *TemplatePage) AddRows(rows []Row) {
	for _, v := range rows {
		t.Rows = append(t.Rows, v)
	}
}

// A very basic output for tests
var DefaultTemplate = `{{range .Rows}}
{{index .Columns 0}}, {{index .Columns 1}}{{end}}`

// A simple terraform template for aws amis in zones
var TerraformTemplate = `variable "images" {
    default = {
{{range $index, $row := .Rows}}{{if eq (index $row.Columns 2) "artifact"}}{{if eq (index $row.Columns 4) "id"}}{{ $artifact := (index $row.Columns 5) }}{{ $artifactb := ($artifact | Split ":")}}
        {{index $artifactb 0}} = "{{index $artifactb 1}}"{{end}}{{end}}{{end}}
    }
}`

// ReadCSV converts the csv files into a data structure we can use
func ReadCSV(csvReader io.Reader) (ret [][]string, err error) {
	reader := csv.NewReader(csvReader)
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	return reader.ReadAll()
}

// this is a special string split that works for templates
func split(s string, d string) []string {
	// panic(d)
	arr := strings.Split(d, s)
	return arr
}

// ToTemplate applies the Rows to a given template string
func ToTemplate(page TemplatePage, tmpl string) (ret string, err error) {

	tplFuncMap := make(template.FuncMap)

	tplFuncMap["Split"] = split

	t := template.Must(template.New("tmpl").Funcs(tplFuncMap).Parse(tmpl))

	var doc bytes.Buffer
	err = t.Execute(&doc, page)
	ret = doc.String()

	// bad template, throw error
	if err != nil {
		panic(err)
	}

	return ret, err
}
