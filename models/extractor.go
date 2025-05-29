package models

type ExtractResult struct {
	CoverageBlocks []*CoverageBlock
	ByFile         map[string]CoverageFileSummary
	ByDir          map[string]CoverageDirSummary
	Total          *CoverageSummary
}
