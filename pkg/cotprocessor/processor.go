package cotprocessor

import (
	"encoding/csv"
	"fmt"
	"github.com/paulschick/cot-excel-processor/pkg/fs"
	grate "github.com/pbnjay/grate/xls"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Processor struct {
	OrderedColNames  []string
	ColumnNamesIndex map[string]int
	Writer           *csv.Writer
	OutputFile       *os.File
}

func NewProcessor() *Processor {
	return &Processor{
		OrderedColNames: []string{
			"Market_and_Exchange_Names",
			"Report_Date_as_MM_DD_YYYY",
			"Open_Interest_All",
			"NonComm_Positions_Long_All",
			"NonComm_Positions_Short_All",
			"Comm_Positions_Long_All",
			"Comm_Positions_Short_All",
			"NonRept_Positions_Long_All",
			"NonRept_Positions_Short_All",
		},
		ColumnNamesIndex: make(map[string]int),
	}
}

func (p *Processor) createOutputFile(xlsFileName string, outputDir string) error {
	outputFileName := strings.TrimSuffix(xlsFileName, ".xls") + "_processed.csv"
	outputFilePath := filepath.Join(outputDir, outputFileName)
	err := fs.EnsureDirExists(outputDir)
	if err != nil {
		return err
	}

	p.OutputFile, err = os.Create(outputFilePath)
	if err != nil {
		return err
	}

	p.Writer = csv.NewWriter(p.OutputFile)
	return nil
}

func (p *Processor) initializeColumnIndices(headerRow []string) {
	for _, colName := range p.OrderedColNames {
		p.ColumnNamesIndex[colName] = -1
	}
	for colIdx := 0; colIdx < len(headerRow); colIdx++ {
		col := headerRow[colIdx]
		if _, ok := p.ColumnNamesIndex[col]; ok {
			p.ColumnNamesIndex[col] = colIdx
		}
	}
}

func (p *Processor) processRow(row []string) error {
	csvRow := make([]string, len(p.OrderedColNames))

	for idx, colName := range p.OrderedColNames {
		colIdx, exists := p.ColumnNamesIndex[colName]
		if exists && colIdx != -1 {
			csvRow[idx] = row[colIdx]
		} else {
			csvRow[idx] = ""
		}
	}

	return p.Writer.Write(csvRow)
}

func (p *Processor) ProcessXLS(filePath string, outputDir string) error {
	if err := p.createOutputFile(filepath.Base(filePath), outputDir); err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}

	defer func() {
		p.Writer.Flush()
		if err := p.Writer.Error(); err != nil {
			panic(err)
		}
		err := p.OutputFile.Close()
		if err != nil {
			panic(err)
		}
	}()

	wb, err := grate.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening xls file: %v", err)
	}

	sheets, err := wb.List()
	if err != nil {
		return fmt.Errorf("error listing sheets: %v", err)
	}

	initializedHeaders := false
	for _, s := range sheets {
		sheet, err := wb.Get(s)
		if err != nil {
			return fmt.Errorf("error getting sheet: %v", err)
		}

		for sheet.Next() {
			row := sheet.Strings()
			if !initializedHeaders {
				p.initializeColumnIndices(row)
				headers := append(
					p.OrderedColNames,
					"NonComm_Positions_Net_All",
					"Comm_Positions_Net_All",
					"NonRept_Positions_Net_All",
				)
				err := p.Writer.Write(headers)
				if err != nil {
					return fmt.Errorf("error writing headers to CSV: %v", err)
				}
				initializedHeaders = true
				continue
			}
			csvRow := make([]string, len(p.OrderedColNames))

			for idx, colName := range p.OrderedColNames {
				colIdx := p.ColumnNamesIndex[colName]
				if colIdx != -1 && colIdx < len(row) {
					csvRow[idx] = row[colIdx]
				} else {
					csvRow[idx] = ""
				}
			}

			// Calculating the new columns
			if nonCommLong, nonCommShort, err := p.getValuesFromRow(row, "NonComm_Positions_Long_All", "NonComm_Positions_Short_All"); err == nil {
				csvRow = append(csvRow, strconv.Itoa(nonCommLong-nonCommShort))
			} else {
				csvRow = append(csvRow, "")
			}

			if commLong, commShort, err := p.getValuesFromRow(row, "Comm_Positions_Long_All", "Comm_Positions_Short_All"); err == nil {
				csvRow = append(csvRow, strconv.Itoa(commLong-commShort))
			} else {
				csvRow = append(csvRow, "")
			}

			if nonReptLong, nonReptShort, err := p.getValuesFromRow(row, "NonRept_Positions_Long_All", "NonRept_Positions_Short_All"); err == nil {
				csvRow = append(csvRow, strconv.Itoa(nonReptLong-nonReptShort))
			} else {
				csvRow = append(csvRow, "")
			}

			err = p.Writer.Write(csvRow)
			if err != nil {
				return fmt.Errorf("failed writing to CSV: %v", err)
			}
		}
	}
	return nil
}

func (p *Processor) ProcessXLSFiles(dataDir string, outputDir string) error {
	err := fs.EnsureDirExists(dataDir)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(dataDir)
	if err != nil {
		return err
	}

	for fileNo, file := range files {
		fmt.Println("Processing file (", fileNo+1, ") ", file.Name())

		if !file.IsDir() && strings.HasSuffix(file.Name(), ".xls") {
			filePath := filepath.Join(dataDir, file.Name())

			err = p.ProcessXLS(filePath, outputDir)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Processor) getValuesFromRow(row []string, col1, col2 string) (int, int, error) {
	index1, ok1 := p.ColumnNamesIndex[col1]
	index2, ok2 := p.ColumnNamesIndex[col2]

	if !ok1 || !ok2 {
		return 0, 0, fmt.Errorf("one or both columns not found: %s, %s", col1, col2)
	}

	if index1 >= len(row) || index2 >= len(row) {
		return 0, 0, fmt.Errorf("row does not contain values for one or both columns: %s, %s", col1, col2)
	}

	v1, err := strconv.Atoi(row[index1])
	if err != nil {
		return 0, 0, err
	}

	v2, err := strconv.Atoi(row[index2])
	if err != nil {
		return 0, 0, err
	}

	return v1, v2, nil
}
