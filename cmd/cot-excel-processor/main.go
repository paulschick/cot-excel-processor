package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	downloaderService "github.com/paulschick/cot-downloader/pkg/service"
	grate "github.com/pbnjay/grate/xls"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("COT Excel Processor")
	shouldDownload := flag.Bool("download", false, "Download reports before processing.")

	if *shouldDownload {
		initReports()
	}

	dataDir := "./data"
	files, err := os.ReadDir(dataDir)
	if err != nil {
		panic(err)
	}

	filePath, err := GetFilePath(files[0], dataDir)
	if err != nil {
		panic(err)
	}

	outputFile, err := os.Create("output.csv")
	if err != nil {
		panic(err)
	}

	defer func() {
		err := outputFile.Close()
		if err != nil {
			panic(err)
		}
	}()

	writer := csv.NewWriter(outputFile)

	fmt.Println("Reading ", filePath)

	orderedColNames := []string{
		"Market_and_Exchange_Names",
		"Report_Date_as_MM_DD_YYYY",
		"Open_Interest_All",
		"NonComm_Positions_Long_All",
		"NonComm_Positions_Short_All",
		"Comm_Positions_Long_All",
		"Comm_Positions_Short_All",
		"NonRept_Positions_Long_All",
		"NonRept_Positions_Short_All",
	}
	columnNamesIndex := make(map[string]int)

	for _, colName := range orderedColNames {
		columnNamesIndex[colName] = -1
	}

	wb, _ := grate.Open(filePath)
	sheets, _ := wb.List()
	for _, s := range sheets {
		sheet, _ := wb.Get(s)
		i := 0
		for sheet.Next() {
			row := sheet.Strings()
			if i == 0 {
				for colIdx := 0; colIdx < len(row); colIdx++ {
					col := row[colIdx]
					if currentColIndex, ok := columnNamesIndex[col]; ok && currentColIndex == -1 {
						columnNamesIndex[col] = colIdx
					}
				}
			}
			csvRow := make([]string, 0)

			for _, colName := range orderedColNames {
				colIdx := columnNamesIndex[colName]
				if colIdx != -1 {
					csvRow = append(csvRow, row[colIdx])
				} else {
					csvRow = append(csvRow, "")
				}
			}
			err = writer.Write(csvRow)
			if err != nil {
				panic(err)
			}

			i++
		}
	}
}

func GetFilePath(file os.DirEntry, basePath string) (string, error) {
	if file == nil {
		return "", fmt.Errorf("file is nil")
	}
	if file.IsDir() {
		return "", fmt.Errorf("file is a directory")
	}
	return filepath.Join(basePath, file.Name()), nil
}

func initReports() {
	downloaderService.DownloadReports(2022, 2023, "legacy", "./data", 500)
}
