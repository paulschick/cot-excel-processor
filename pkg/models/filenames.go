package models

import (
	"path/filepath"
	"strings"
)

type FileNameManager struct {
	XlsDir string
	CsvDir string
}

func NewFileNameManager(xlsDir, csvDir string) *FileNameManager {
	return &FileNameManager{
		XlsDir: xlsDir,
		CsvDir: csvDir,
	}
}

func (f *FileNameManager) GetFileNames(xlsName string) *FileNames {
	xlsPath := filepath.Join(f.XlsDir, xlsName)
	csvName := strings.TrimSuffix(xlsName, ".xls") + "_processed.csv"
	csvPath := filepath.Join(f.CsvDir, csvName)
	return NewFileNames(xlsPath, csvPath)
}

type FileNames struct {
	XlsPath string
	CsvPath string
}

func NewFileNames(xlsPath, csvPath string) *FileNames {
	return &FileNames{
		XlsPath: xlsPath,
		CsvPath: csvPath,
	}
}

func (f *FileNames) GetXlsDir() string {
	return filepath.Dir(f.XlsPath)
}

func (f *FileNames) GetCsvDir() string {
	return filepath.Dir(f.CsvPath)
}
