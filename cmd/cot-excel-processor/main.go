package main

import (
	"flag"
	"fmt"
	downloaderService "github.com/paulschick/cot-downloader/pkg/service"
	"github.com/paulschick/cot-excel-processor/pkg/cotprocessor"
)

func main() {
	fmt.Println("COT Excel Processor")
	shouldDownload := flag.Bool("download", false, "Download reports before processing.")

	if *shouldDownload {
		initReports()
	}

	dataDir := "./data"
	outputDir := "./output"
	processor := cotprocessor.NewProcessor()

	err := processor.ProcessXLSFiles(dataDir, outputDir)
	if err != nil {
		panic(err)
	}
}

func initReports() {
	downloaderService.DownloadReports(2022, 2023, "legacy", "./data", 500)
}
