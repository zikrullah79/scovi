package service

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/zikrullah79/scovi/helper"
	"github.com/zikrullah79/scovi/models"
)

type ExtractorService interface {
	ExtractCoverageFile(filename string) (*models.ExtractResult, error)
	SummarizeCoverage(result *models.ExtractResult) *models.ParsedCoverageInformation
}

type extractorService struct {
}

func NewParseService() ExtractorService {
	return &extractorService{}
}
func (s *extractorService) ExtractCoverageFile(filename string) (*models.ExtractResult, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var blocks []*models.CoverageBlock
	byFile := make(map[string]models.CoverageFileSummary)
	byDir := make(map[string]models.CoverageDirSummary)
	total := &models.CoverageSummary{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		block := helper.ExtractLine(line)
		if block == nil {
			continue // skip invalid lines
		}

		blocks = append(blocks, block)

		file := block.FilePath
		dir := filepath.Dir(file)

		// Update file-level
		fileSum := byFile[file]
		fileSum.Directory = dir
		fileSum.TotalStatements += block.NumStatements
		if block.Count > 0 {
			fileSum.CoveredStatements += block.NumStatements
		}
		byFile[file] = fileSum

		// Update dir-level
		dirSum := byDir[dir]
		dirSum.TotalFiles++
		dirSum.TotalStatements += block.NumStatements
		if block.Count > 0 {
			dirSum.CoveredStatements += block.NumStatements
		}
		byDir[dir] = dirSum

		// Update total
		total.TotalStatements += block.NumStatements
		if block.Count > 0 {
			total.CoveredStatements += block.NumStatements
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &models.ExtractResult{
		CoverageBlocks: blocks,
		ByFile:         byFile,
		ByDir:          byDir,
		Total:          total,
	}, nil
}

// TODO: need to improve this function.
func (s *extractorService) SummarizeCoverage(result *models.ExtractResult) *models.ParsedCoverageInformation {
	res := models.NewParsedCoverageInformation()

	res.CoveragePercentage = result.Total.Percent()

	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		for k, v := range result.ByFile {
			percent := v.Percent()
			mu.Lock()

			filePathSplit := strings.Split(k, "/")
			fileName := filePathSplit[len(filePathSplit)-1]
			// res.CoveragePercentageByFile[k] = percent
			res.CoveragePercentageByFileString += fmt.Sprintf(`
				<tr class="border-t border-t-[#dce0e5]">
					<td class="table-d717db04-a7ff-4f22-9c5d-6769a2b0b51d-column-120 h-[72px] px-4 py-2 w-[400px] text-[#111418] text-sm font-normal leading-normal">%v</td>
					<td class="table-d717db04-a7ff-4f22-9c5d-6769a2b0b51d-column-240 h-[72px] px-4 py-2 w-[400px] text-[#637588] text-sm font-normal leading-normal">%v</td>
					<td class="table-d717db04-a7ff-4f22-9c5d-6769a2b0b51d-column-360 h-[72px] px-4 py-2 w-[400px] text-[#637588] text-sm font-normal leading-normal">%v</td>
				</tr>
					  `, fileName, v.Directory, strconv.FormatFloat(percent, 'f', 2, 64))
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		for k, v := range result.ByDir {
			percent := v.Percent()
			mu.Lock()
			// res.CoveragePercentageByDirectory[k] = percent
			res.AllDirectoriesTag += fmt.Sprintf(`
				<div class="flex h-8 shrink-0 items-center justify-center gap-x-2 rounded-xl bg-[#f0f2f4] pl-4 pr-4">
                	<p class="text-[#111418] text-sm font-medium leading-normal">%v</p>
              	</div>`, k)
			res.CoveragePercentageByDirectoryString += fmt.Sprintf(`
				<tr class="border-t border-t-[#dce0e5]">
					<td class="table-152b368a-9747-410a-b3ca-e108bd48c085-column-120 h-[72px] px-4 py-2 w-[400px] text-[#111418] text-sm font-normal leading-normal">%v</td>
                    <td class="table-152b368a-9747-410a-b3ca-e108bd48c085-column-240 h-[72px] px-4 py-2 w-[400px] text-[#637588] text-sm font-normal leading-normal">%v</td>
                    <td class="table-152b368a-9747-410a-b3ca-e108bd48c085-column-360 h-[72px] px-4 py-2 w-[400px] text-[#637588] text-sm font-normal leading-normal">%v</td>
				</tr>
					  `, k, v.TotalFiles, strconv.FormatFloat(percent, 'f', 2, 64))
			mu.Unlock()
		}
	}()

	wg.Wait()
	return res
}
