package main

import (
	"fmt"
	downloaderService "github.com/paulschick/cot-downloader/pkg/service"
)

func main() {
	fmt.Println("COT Excel Processor")
	initReports()
}

func initReports() {
	downloaderService.DownloadReports(2022, 2023, "legacy", "./data", 500)
}
