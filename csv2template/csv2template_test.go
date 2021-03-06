package csv2template

import (
	"strings"
	"testing"
)

// A simple terraform template for aws amis in multiple zones
var terraformTemplate = `variable "images" {
    default = {
{{range $index, $row := .Rows}}{{if eq (index $row.Columns 2) "artifact"}}{{if eq (index $row.Columns 4) "id"}}{{ $artifact := (index $row.Columns 5) }}{{ $artifactb := ($artifact | Split ":")}}
        {{index $artifactb 0}} = "{{index $artifactb 1}}"{{end}}{{end}}{{end}}
    }
}`

// A simple packer template for aws amis in multiple zones
var packerTemplate = `{
{{range $index, $row := .Rows}}{{if eq (index $row.Columns 2) "artifact"}}{{if eq (index $row.Columns 4) "id"}}{{ $artifact := (index $row.Columns 5) }}{{ $artifactb := ($artifact | Split ":")}}
    "{{index $artifactb 0}}": "{{index $artifactb 1}}"{{if lt $index ((len $row.Columns) -1)}},{{end}}{{end}}{{end}}{{end}}
}`

func csvToPage(data string) (page TemplatePage) {
	splitData := strings.Split(data, "\n")
	for _, v := range splitData {
		row := Row{strings.Split(v, ",")}
		page.Rows = append(page.Rows, row)
	}
	return page
}

func TestToBasicTemplate(t *testing.T) {
	page := csvToPage(`1432168589,amazon-ebs,artifact-count,2
1432168589,amazon-ebs,artifact,0,builder-id,mitchellh.amazonebs
1432168589,amazon-ebs,artifact,0,id,us-west-1:ami-df79909b
1432168589,amazon-ebs,artifact,0,string,AMIs were created:\n\nus-west-1: ami-df79909b
1432168589,amazon-ebs,artifact,0,files-count,0
1432168589,amazon-ebs,artifact,0,end
1432168589,amazon-ebs,artifact,1,builder-id,mitchellh.amazonebs
1432168589,amazon-ebs,artifact,1,id,us-west-2:ami-df79909c
1432168589,amazon-ebs,artifact,1,string,AMIs were created:\n\nus-west-2: ami-df79909c
1432168589,amazon-ebs,artifact,1,files-count,0
1432168589,amazon-ebs,artifact,1,end`)
	out := `
1432168589, amazon-ebs
1432168589, amazon-ebs
1432168589, amazon-ebs
1432168589, amazon-ebs
1432168589, amazon-ebs
1432168589, amazon-ebs
1432168589, amazon-ebs
1432168589, amazon-ebs
1432168589, amazon-ebs
1432168589, amazon-ebs
1432168589, amazon-ebs`

	doc, err := ToTemplate(page, DefaultTemplate)
	if err != nil {
		t.Log("Template transform produced an error")
		t.Logf("Error: %#v", err)
		t.Fail()
	}
	if doc != out {
		t.Log("Template transform didn't produce correct output")
		t.Logf("Doc: %#v", doc)
		t.Logf("Output: %#v", out)
		t.Fail()
	}
}

func TestToTerraformTemplate(t *testing.T) {
	page := csvToPage(`1432168589,amazon-ebs,artifact-count,2
1432168589,amazon-ebs,artifact,0,builder-id,mitchellh.amazonebs
1432168589,amazon-ebs,artifact,0,id,us-west-1:ami-df79909b
1432168589,amazon-ebs,artifact,0,string,AMIs were created:\n\nus-west-1: ami-df79909b
1432168589,amazon-ebs,artifact,0,files-count,0
1432168589,amazon-ebs,artifact,0,end
1432168589,amazon-ebs,artifact,1,builder-id,mitchellh.amazonebs
1432168589,amazon-ebs,artifact,1,id,us-west-2:ami-df79909c
1432168589,amazon-ebs,artifact,1,string,AMIs were created:\n\nus-west-2: ami-df79909c
1432168589,amazon-ebs,artifact,1,files-count,0
1432168589,amazon-ebs,artifact,1,end`)
	out := `variable "images" {
    default = {

        us-west-1 = "ami-df79909b"
        us-west-2 = "ami-df79909c"
    }
}`

	doc, err := ToTemplate(page, terraformTemplate)
	if err != nil {
		t.Log("Terraform Template transform produced an error")
		t.Logf("Error: %#v", err)
		t.Fail()
	}
	if doc != out {
		t.Log("Terraform Template transform didn't produce correct output")
		t.Logf("Doc: %s", doc)
		t.Logf("Output: %s", out)
		t.Fail()
	}
}

func TestToPackerTemplate(t *testing.T) {
	page := csvToPage(`1432168589,amazon-ebs,artifact-count,2
1432168589,amazon-ebs,artifact,0,builder-id,mitchellh.amazonebs
1432168589,amazon-ebs,artifact,0,id,us-west-1:ami-df79909b
1432168589,amazon-ebs,artifact,0,string,AMIs were created:\n\nus-west-1: ami-df79909b
1432168589,amazon-ebs,artifact,0,files-count,0
1432168589,amazon-ebs,artifact,0,end
1432168589,amazon-ebs,artifact,1,builder-id,mitchellh.amazonebs
1432168589,amazon-ebs,artifact,1,id,us-west-2:ami-df79909c
1432168589,amazon-ebs,artifact,1,string,AMIs were created:\n\nus-west-2: ami-df79909c
1432168589,amazon-ebs,artifact,1,files-count,0
1432168589,amazon-ebs,artifact,1,end`)
	out := `{

    "us-west-1": "ami-df79909b",
    "us-west-2": "ami-df79909c"
}`

	doc, err := ToTemplate(page, packerTemplate)
	if err != nil {
		t.Log("Packer Template transform produced an error")
		t.Logf("Error: %#v", err)
		t.Fail()
	}
	if doc != out {
		t.Log("Packer Template transform didn't produce correct output")
		t.Logf("Doc: %s", doc)
		t.Logf("Output: %s", out)
		t.Fail()
	}
}
