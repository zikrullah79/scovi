package helper

import (
	"strconv"
	"strings"

	"github.com/zikrullah79/scovi/models"
)

func ExtractLine(line string) *models.CoverageBlock {

	// Skip mode line
	if strings.HasPrefix(line, "mode:") {
		return nil
	}

	// Example line:
	// github.com/example/project/math.go:10.21,12.37 1 1
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return nil // skip invalid lines
	}

	fileAndRange := parts[0]
	numStatements, err1 := strconv.Atoi(parts[1])
	count, err2 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil {
		return nil
	}

	fileSplit := strings.Split(fileAndRange, ":")
	if len(fileSplit) < 2 {
		return nil
	}

	// get last part of string after last slash
	filePath := fileSplit[0]
	rangePart := strings.TrimPrefix(line, filePath+":")
	rangePart = strings.Fields(rangePart)[0] // get "10.21,12.37"

	ranges := strings.Split(rangePart, ",")
	if len(ranges) != 2 {
		return nil
	}

	start := strings.Split(ranges[0], ".")
	end := strings.Split(ranges[1], ".")
	if len(start) != 2 || len(end) != 2 {
		return nil
	}

	startLine, _ := strconv.Atoi(start[0])
	startCol, _ := strconv.Atoi(start[1])
	endLine, _ := strconv.Atoi(end[0])
	endCol, _ := strconv.Atoi(end[1])

	return &models.CoverageBlock{
		FilePath:      filePath,
		StartLine:     startLine,
		StartCol:      startCol,
		EndLine:       endLine,
		EndCol:        endCol,
		NumStatements: numStatements,
		Count:         count,
	}
}
