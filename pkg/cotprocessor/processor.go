package cotprocessor

import (
	"github.com/paulschick/cot-excel-processor/pkg/fs"
	"github.com/paulschick/cot-excel-processor/pkg/models"
	"os"
	"strings"
)

type Processor struct {
	FileNameManager *models.FileNameManager
}

func NewProcessor(xlsDir, csvDir string) *Processor {
	return &Processor{
		FileNameManager: models.NewFileNameManager(xlsDir, csvDir),
	}
}

func (p *Processor) ProcessXLSFiles() error {
	files, err := p.getFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".xls") {
			fileNames := p.FileNameManager.GetFileNames(file.Name())
			xlsProcessor, err := NewXlsProcessor(fileNames)
			if err != nil {
				return err
			}
			err = xlsProcessor.ProcessXlsToCsv()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Processor) ProcessXLSForYearRange(startYear, endYear string) error {
	files, err := p.getFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".xls") &&
			(strings.Contains(file.Name(), startYear) || strings.Contains(file.Name(), endYear)) {
			fileNames := p.FileNameManager.GetFileNames(file.Name())
			xlsProcessor, err := NewXlsProcessor(fileNames)
			if err != nil {
				return err
			}
			err = xlsProcessor.ProcessXlsToCsv()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Processor) getFiles() ([]os.DirEntry, error) {
	err := fs.EnsureDirExists(p.FileNameManager.XlsDir)
	if err != nil {
		return nil, err
	}
	err = fs.EnsureDirExists(p.FileNameManager.CsvDir)
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(p.FileNameManager.XlsDir)
	if err != nil {
		return nil, err
	}

	return files, nil
}
