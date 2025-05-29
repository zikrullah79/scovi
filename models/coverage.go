package models

import "math"

type CoverageBlock struct {
	FilePath      string
	StartLine     int
	StartCol      int
	EndLine       int
	EndCol        int
	NumStatements int
	Count         int
}

type ParsedCoverageInformation struct {
	// Title of the report
	Title string
	// all files
	CoveragePercentage float64

	AllDirectoriesTag                   string
	CoveragePercentageByDirectoryString string
	CoveragePercentageByFileString      string
}

func NewParsedCoverageInformation() *ParsedCoverageInformation {
	p := &ParsedCoverageInformation{
		Title: "Coverage Report",
	}
	return p
}

type CoverageSummary struct {
	TotalStatements   int
	CoveredStatements int
}

type CoverageDirSummary struct {
	CoverageSummary
	TotalFiles int
}

type CoverageFileSummary struct {
	CoverageSummary
	Directory string
}

func (c CoverageSummary) Percent() float64 {
	if c.TotalStatements == 0 {
		return 100.0
	}
	// return with 2 decimal places.
	return math.Round((float64(c.CoveredStatements)/float64(c.TotalStatements))*100*100) / 100
}
