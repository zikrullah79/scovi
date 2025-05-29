package service

import (
	"os"
	"path"
	"text/template"

	"github.com/zikrullah79/scovi/models"
)

type ParserService interface {
	ParseHtmlFile(parsedCoverageInformation *models.ParsedCoverageInformation, filename string) error
}
type parserService struct {
}

func NewParserService() ParserService {
	return &parserService{}
}

func (p *parserService) ParseHtmlFile(data *models.ParsedCoverageInformation, filename string) error {
	if filename == "" {
		filename = "coverage_report.html" // Default filename if none is provided
	}

	var filepath = path.Join("template", "report.html")
	tmp, err := template.ParseFiles(filepath)
	if err != nil {
		return err
	}

	// Create a new file with the same name as the project name
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// Execute the template and write to the file
	err = tmp.Execute(file, data)
	if err != nil {
		return err
	}

	return nil
}
