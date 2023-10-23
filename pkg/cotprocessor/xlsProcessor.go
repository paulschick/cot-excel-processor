package cotprocessor

import (
	"encoding/csv"
	"github.com/paulschick/cot-excel-processor/pkg/models"
	"github.com/pbnjay/grate"
	grateXls "github.com/pbnjay/grate/xls"
	"os"
)

type XlsProcessor struct {
	FileNames          *models.FileNames
	Wb                 grate.Source
	CsvFile            *os.File
	CsvWriter          *csv.Writer
	ColumnIndices      *models.ColumnNameIndices
	initializedHeaders bool
}

func NewXlsProcessor(fileNames *models.FileNames) (*XlsProcessor, error) {
	wb, err := grateXls.Open(fileNames.XlsPath)
	if err != nil {
		return nil, err
	}
	xlsProcessor := &XlsProcessor{
		FileNames:          fileNames,
		Wb:                 wb,
		ColumnIndices:      models.NewColumnNameIndices(),
		initializedHeaders: false,
	}
	err = xlsProcessor.CreateCsvFile()
	if err != nil {
		return nil, err
	}
	return xlsProcessor, nil
}

func (x *XlsProcessor) CloseOrPanic() {
	x.CsvWriter.Flush()
	err := x.CsvWriter.Error()
	if err != nil {
		panic(err)
	}

	err = x.Wb.Close()
	if err != nil {
		panic(err)
	}

	err = x.CsvFile.Close()
	if err != nil {
		panic(err)
	}
}

func (x *XlsProcessor) GetSheets() ([]string, error) {
	return x.Wb.List()
}

// CreateCsvFile Creates a CSV file for the XLS file and opens a new CSV Writer
func (x *XlsProcessor) CreateCsvFile() error {
	var err error
	x.CsvFile, err = os.Create(x.FileNames.CsvPath)
	if err != nil {
		return err
	}
	x.CsvWriter = csv.NewWriter(x.CsvFile)
	return nil
}

func (x *XlsProcessor) WriteHeaders(headerRow []string) error {
	x.ColumnIndices.InitializeColumnIndices(headerRow)
	err := x.CsvWriter.Write(x.ColumnIndices.FullOrderedColumnNames)
	if err != nil {
		return err
	}
	x.initializedHeaders = true
	return nil
}

func (x *XlsProcessor) processCollection(collection grate.Collection) error {
	rowSlice := collection.Strings()
	if !x.initializedHeaders {
		err := x.WriteHeaders(rowSlice)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
	row := models.NewRow(rowSlice, x.ColumnIndices.ColumnNamesIndex)
	csvRow := row.GetCsvRow()
	err := x.CsvWriter.Write(csvRow)
	if err != nil {
		return err
	}
	return nil
}

func (x *XlsProcessor) ProcessXlsToCsv() error {
	defer func() {
		x.CloseOrPanic()
	}()

	sheets, err := x.GetSheets()
	if err != nil {
		return err
	}
	for _, s := range sheets {
		sheet, err := x.Wb.Get(s)
		if err != nil {
			return err
		}
		for sheet.Next() {
			err = x.processCollection(sheet)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
