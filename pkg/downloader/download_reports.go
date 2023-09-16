package downloader

import (
	downloaderService "github.com/paulschick/cot-downloader/pkg/service"
	"time"
)

func InitReportsDefault(start int, end int, outDir string) {
	downloaderService.DownloadReports(start, end, "legacy", outDir, 500)
}

func UpdateCurrentYear(outDir string) {
	startEnd := time.Now().Year()
	downloaderService.DownloadReports(startEnd, startEnd, "legacy", outDir, 500)
}
