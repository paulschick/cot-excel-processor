package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/extrame/xls"
	downloaderService "github.com/paulschick/cot-downloader/pkg/service"
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

	filePath, err := getFilePath(files[0], dataDir)
	if err != nil {
		panic(err)
	}

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

	workbook, err := OpenExcel(filePath)
	if err != nil {
		panic(err)
	}
	sheet, err := GetFirstSheet(workbook)
	if err != nil {
		panic(err)
	}
	columnNamesIndex := getColumnIndices(sheet, orderedColNames)

	fmt.Println(columnNamesIndex)

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

	file, err := xls.Open(filePath, "utf-8")
	if err != nil {
		panic(err)
	}

	if sheet := file.GetSheet(0); sheet != nil {
		for i := 0; i <= int(sheet.MaxRow); i++ {
			row := sheet.Row(i)
			var csvRow []string
			for _, colName := range orderedColNames {
				colIdx := columnNamesIndex[colName]
				//fmt.Println(colName, colIdx)
				if colIdx != -1 {
					csvRow = append(csvRow, row.Col(colIdx))
				} else {
					csvRow = append(csvRow, "")
				}
				fmt.Println(csvRow)
			}
			err = writer.Write(csvRow)
			if err != nil {
				panic(err)
			}
		}
	}

	writer.Flush()
	fmt.Println("Done")
}

func getFilePath(file os.DirEntry, basePath string) (string, error) {
	if file == nil {
		return "", fmt.Errorf("file is nil")
	}
	if file.IsDir() {
		return "", fmt.Errorf("file is a directory")
	}
	return filepath.Join(basePath, file.Name()), nil
}

func OpenExcel(filePath string) (*xls.WorkBook, error) {
	file, err := xls.Open(filePath, "utf-8")
	if err != nil {
		return nil, err
	}
	return file, nil
}

func GetFirstSheet(workbook *xls.WorkBook) (*xls.WorkSheet, error) {
	if workbook == nil {
		return nil, fmt.Errorf("workbook is nil")
	}
	if sheet := workbook.GetSheet(0); sheet != nil {
		return sheet, nil
	}
	return nil, fmt.Errorf("sheet is nil")
}

func getColumnIndices(sheet *xls.WorkSheet, orderedColNames []string) map[string]int {
	columnNamesIndex := make(map[string]int)
	for _, colName := range orderedColNames {
		columnNamesIndex[colName] = -1
	}

	row1 := sheet.Row(0)
	foundColumns := 0
	for idx := row1.FirstCol(); idx < row1.LastCol(); idx++ {
		colName := row1.Col(idx)
		if colIndex, ok := columnNamesIndex[colName]; ok && colIndex == -1 {
			columnNamesIndex[colName] = idx
			foundColumns++
			if foundColumns == len(columnNamesIndex) {
				break
			}
		}
	}
	return columnNamesIndex
}

func initReports() {
	downloaderService.DownloadReports(2022, 2023, "legacy", "./data", 500)
}
