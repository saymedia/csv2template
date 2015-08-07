package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/saymedia/csv2template/csv2template"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func help() {
	log.Println(`Usage csv2template [options...]
csv2template turns Packer's machine-readable output into a Terraform-readable format.

Options:
    -f Filename of the input CSV. Alternatively use STDIN.
    -h This help information.
    -template Filename of the template to use in the output.

Example:
    packer -machine-readable build app.json | \
        csv2template -template templates/amazon-ebs.hcl > app.tfvars
`)
}

func main() {

	tmpl := flag.String("template", "", "a template file")
	csv := flag.String("f", "", "a csv file")
	helpMe := flag.Bool("h", false, "help")

	flag.Parse()

	if *helpMe {
		help()
		os.Exit(0)
	}

	// Read a file or use STDIN
	var csvFile io.Reader
	if len(*csv) > 0 {
		f, err := os.Open(*csv)
		if err != nil {
			log.Fatalf("CSV file read failed %s", err)
		}
		csvFile = bufio.NewReader(f)
	} else if stat, _ := os.Stdin.Stat(); stat.Mode()&os.ModeCharDevice == 0 {
		// has STDIN data
		csvFile = bufio.NewReader(os.Stdin)
	} else {
		// No input data
		help()
		os.Exit(0)
	}

	// Get the CSV as a string array
	rawData, err := csv2template.ReadCSV(csvFile)
	if err != nil {
		log.Fatalf("CSV read failed %s", err)
	}

	page := csv2template.TemplatePage{
		Rows: []csv2template.Row{},
	}
	for _, v := range rawData {
		row := csv2template.Row{
			Columns: v,
		}
		page.Rows = append(page.Rows, row)
	}

	// Print rows using a template
	var templateString string
	if len(*tmpl) == 0 {
		templateString = csv2template.DefaultTemplate
	} else {
		buf, err := ioutil.ReadFile(*tmpl)
		if err != nil {
			log.Fatalf("Template file read failed: %s", err)
		}
		templateString = string(buf)
	}
	doc, err := csv2template.ToTemplate(page, templateString)
	if err != nil {
		log.Fatalf("Template render failed: %s", err)
	}
	fmt.Println(doc)

	// Done
	os.Exit(0)

}
