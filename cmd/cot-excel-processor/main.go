package main

import (
	"flag"
	"fmt"
	downloaderService "github.com/paulschick/cot-downloader/pkg/service"
)

func main() {
	fmt.Println("COT Excel Processor")
	shouldDownload := flag.Bool("download", false, "Download reports before processing.")

	if *shouldDownload {
		initReports()
	}
}

func initReports() {
	downloaderService.DownloadReports(2022, 2023, "legacy", "./data", 500)
}
