package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Download    bool
	StartYear   int
	EndYear     int
	DownloadDir string
	OutputDir   string
}

func Configurations() *Config {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	return &Config{
		Download: func() bool {
			download := os.Getenv("DOWNLOAD")
			if download == "" || download == "false" {
				return false
			}
			return true
		}(),
		StartYear: func() int {
			startYear := os.Getenv("START_YEAR")
			return handleYears(startYear)
		}(),
		EndYear: func() int {
			endYear := os.Getenv("END_YEAR")
			return handleYears(endYear)
		}(),
		DownloadDir: func() string {
			downloadDir := os.Getenv("DOWNLOAD_DIR")
			if downloadDir == "" {
				return "./data/xls"
			}
			return downloadDir
		}(),
		OutputDir: func() string {
			outputDir := os.Getenv("OUTPUT_DIR")
			if outputDir == "" {
				return "./data/csv"
			}
			return outputDir
		}(),
	}
}

func handleYears(yearString string) int {
	now := time.Now()
	year := now.Year()
	if yearString == "" {
		return year
	}
	yearInt, err := strconv.Atoi(yearString)
	if err != nil {
		log.Printf("Error converting start year to int: %v", err)
		return year
	}
	if yearInt > year {
		log.Printf("Start year is greater than current year. Setting start year to current year.")
		return year
	}
	return yearInt
}
