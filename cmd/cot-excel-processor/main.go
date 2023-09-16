package main

import (
	"flag"
	"fmt"
	downloaderService "github.com/paulschick/cot-downloader/pkg/service"
	"github.com/paulschick/cot-excel-processor/pkg/cotprocessor"
	"github.com/paulschick/cot-excel-processor/pkg/fs"
	"os"
)

func main() {
	var (
		shouldDownload = flag.Bool("download", false, "Download reports before processing.")
		startYear      = flag.Int("startYear", 0, "Start year for report dates")
		endYear        = flag.Int("endYear", 0, "End year for report dates")
		downloadDir    = flag.String("downloadDir", "./data", "Directory for downloaded reports")
		outputDir      = flag.String("outputDir", "./output", "Directory for CSV exports")
	)

	flag.Parse()

	if *shouldDownload {
		if *startYear == 0 || *endYear == 0 {
			fmt.Println("Please provide both startYear and endYear when using the download option")
			os.Exit(1)
		}
		if *startYear > *endYear {
			fmt.Println("startYear should be less than endYear")
			os.Exit(1)
		}
		err := fs.EnsureDirExists(*downloadDir)
		if err != nil {
			fmt.Println("Error ensuring download directory exists:", err)
			os.Exit(1)
		}
		initReports(*startYear, *endYear, *downloadDir)
	}

	err := fs.EnsureDirExists(*outputDir)
	if err != nil {
		fmt.Println("Error ensuring output directory exists:", err)
		os.Exit(1)
	}

	process(*downloadDir, *outputDir)
}

func process(dataDir string, outputDir string) {
	processor := cotprocessor.NewProcessor()

	err := processor.ProcessXLSFiles(dataDir, outputDir)
	if err != nil {
		panic(err)
	}
}

func initReports(startYear int, endYear int, dataDir string) {
	downloaderService.DownloadReports(startYear, endYear, "legacy", dataDir, 500)
}
