package cotprocessor

import (
	"encoding/csv"
	"fmt"
	"github.com/paulschick/cot-excel-processor/pkg/fs"
	"github.com/paulschick/cot-excel-processor/pkg/models"
	grate "github.com/pbnjay/grate/xls"
	"log"
	"os"
	"strings"
)

type Processor struct {
	OrderedColNames  []string
	ColumnNamesIndex map[string]int
	Writer           *csv.Writer
	OutputFile       *os.File
	XlsDir           string
	CsvDir           string
	FileNameManager  *models.FileNameManager
}

func NewProcessor(xlsDir, csvDir string) *Processor {
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
		FileNameManager:  models.NewFileNameManager(xlsDir, csvDir),
	}
}

func (p *Processor) createOutputFile(fileNames *models.FileNames) error {
	err := fs.EnsureDirExists(p.FileNameManager.CsvDir)
	if err != nil {
		return err
	}

	p.OutputFile, err = os.Create(fileNames.CsvPath)
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

func (p *Processor) ProcessXLS(fileNames *models.FileNames) error {
	if err := p.createOutputFile(fileNames); err != nil {
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

	wb, err := grate.Open(fileNames.XlsPath)
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
			rowModel := models.NewRow(row, p.ColumnNamesIndex)
			csvRow := rowModel.GetCsvRow()
			err = p.Writer.Write(csvRow)
			if err != nil {
				return fmt.Errorf("failed writing to CSV: %v", err)
			}
		}
	}
	return nil
}

func (p *Processor) ProcessXLSFiles() error {
	err := fs.EnsureDirExists(p.FileNameManager.XlsDir)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(p.FileNameManager.XlsDir)
	if err != nil {
		return err
	}

	for fileNo, file := range files {
		err := p.processFile(fileNo, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Processor) processFile(fileNo int, file os.DirEntry) error {
	log.Println("Processing file (", fileNo+1, ") ", file.Name())
	if !file.IsDir() && strings.HasSuffix(file.Name(), ".xls") {
		fileNames := p.FileNameManager.GetFileNames(file.Name())
		err := p.ProcessXLS(fileNames)
		if err != nil {
			return err
		}
	}
	return nil
}
