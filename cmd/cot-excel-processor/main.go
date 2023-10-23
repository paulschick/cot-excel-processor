package main

import (
	"flag"
	"fmt"
	downloaderService "github.com/paulschick/cot-downloader/pkg/service"
	"github.com/paulschick/cot-excel-processor/pkg/config"
	"github.com/paulschick/cot-excel-processor/pkg/cotprocessor"
	"github.com/paulschick/cot-excel-processor/pkg/fs"
	"os"
	"time"
)

func main() {
	var (
		environmentFile = flag.Bool("env", true, "Load configuration from .env file. If this is true, all other flags are ignored.")
		shouldDownload  = flag.Bool("download", true, "Download reports before processing.")
		startYear       = flag.Int("startYear", time.Now().Year(), "Start year for report dates")
		endYear         = flag.Int("endYear", time.Now().Year(), "End year for report dates")
		downloadDir     = flag.String("downloadDir", "./data/xls", "Directory for downloaded reports")
		outputDir       = flag.String("outputDir", "./data/csv", "Directory for CSV exports")
	)

	flag.Parse()

	if *environmentFile {
		fmt.Println("Loading configuration from .env file")
		conf := config.Configurations()
		shouldDownload = &conf.Download
		startYear = &conf.StartYear
		endYear = &conf.EndYear
		downloadDir = &conf.DownloadDir
		outputDir = &conf.OutputDir
	} else {
		fmt.Println("Loading configuration from command line flags")
	}

	fmt.Println("shouldDownload:", *shouldDownload)
	fmt.Println("startYear:", *startYear)
	fmt.Println("endYear:", *endYear)
	fmt.Println("downloadDir:", *downloadDir)
	fmt.Println("outputDir:", *outputDir)

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
	processor := cotprocessor.NewProcessor(dataDir, outputDir)

	err := processor.ProcessXLSFiles()
	if err != nil {
		panic(err)
	}
}

func initReports(startYear int, endYear int, dataDir string) {
	downloaderService.DownloadReports(startYear, endYear, "legacy", dataDir, 500)
}
