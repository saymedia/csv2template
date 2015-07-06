# csv2template

csv2template transforms CSV files into [Go Templated](http://golang.org/pkg/text/template/) text output. For example, you have transform Packer build output that produces Terraform config files.

[![travis build status for csv2template](https://travis-ci.org/saymedia/csv2template.svg)](https://travis-ci.org/saymedia/csv2template) [![Coverage Status](https://coveralls.io/repos/saymedia/csv2template/badge.svg?branch=master)](https://coveralls.io/r/saymedia/csv2template?branch=master)

## Usage

csv2template reads from STDIN and writes to STDOUT.

    csv2template -f [input filename] -template [template filename]

## Example

    packer -machine-readable build app.json | csv2template > app.tfvars
    
Or:

    csv2template -f packer_out.csv -template tmpl.tfvars > app.tfvars

Given this CSV input:

    1432168589,amazon-ebs,artifact-count,2
    1432168589,amazon-ebs,artifact,0,builder-id,mitchellh.amazonebs
    1432168589,amazon-ebs,artifact,0,id,us-west-1:ami-df79909b
    1432168589,amazon-ebs,artifact,0,string,AMIs were created:\n\nus-west-1: ami-df79909b
    1432168589,amazon-ebs,artifact,0,files-count,0
    1432168589,amazon-ebs,artifact,0,end
    1432168589,amazon-ebs,artifact,1,builder-id,mitchellh.amazonebs
    1432168589,amazon-ebs,artifact,1,id,us-west-2:ami-df79909c
    1432168589,amazon-ebs,artifact,1,string,AMIs were created:\n\nus-west-2: ami-df79909c
    1432168589,amazon-ebs,artifact,1,files-count,0
    1432168589,amazon-ebs,artifact,1,end

And this template:

    variable "images" {
        default = {
    {{range $index, $row := .Rows}}{{if eq (index $row.Columns 2) "artifact"}}{{if eq (index $row.Columns 4) "id"}}{{ $artifact := (index $row.Columns 5) }}{{ $artifactb := ($artifact | Split ":")}}
            {{index $artifactb 0}} = "{{index $artifactb 1}}"{{end}}{{end}}{{end}}
        }
    }

csv2template will produce this output:

    variable "images" {
        default = {

            us-west-1 = "ami-df79909b"
            us-west-2 = "ami-df79909c"
        }
    }

## Install

    go get github.com/saymedia/csv2template

## Test

    go test ./...

Running `./test.sh` additionally tests using `go vet`, `golint`, `gocyclo`, `gofmt` and `go build`, which should be done before a commit.

## License

Copyright © 2015 Say Media Ltd. All Rights Reserved. See the LICENSE file for distribution terms.
