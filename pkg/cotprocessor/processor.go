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
	err := fs.EnsureDirExists(p.FileNameManager.XlsDir)
	if err != nil {
		return err
	}
	err = fs.EnsureDirExists(p.FileNameManager.CsvDir)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(p.FileNameManager.XlsDir)
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
