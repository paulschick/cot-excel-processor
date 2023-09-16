package cotprocessor

import (
	"encoding/csv"
	"fmt"
	"github.com/paulschick/cot-excel-processor/pkg/fs"
	grate "github.com/pbnjay/grate/xls"
	"os"
	"path/filepath"
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

	for _, s := range sheets {
		sheet, err := wb.Get(s)
		if err != nil {
			return fmt.Errorf("error getting sheet: %v", err)
		}

		i := 0
		for sheet.Next() {
			row := sheet.Strings()
			if i == 0 {
				p.initializeColumnIndices(row)
			}

			if err := p.processRow(row); err != nil {
				return fmt.Errorf("error processing row: %v", err)
			}

			i++
		}
	}
	p.Writer.Flush()
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
