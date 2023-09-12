package main

import (
	"cot/downloader"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Historical Compressed files
// https://www.cftc.gov/MarketReports/CommitmentsofTraders/HistoricalCompressed/index.htm
func main() {
	fmt.Println("COT Downloads")
	// CLI flags
	var (
		reportType = flag.String("report", "legacy", `Report type. Available options: "legacy". Default is "legacy".`)
		startYear  = flag.Int("start", 2001, "Specify the start year.")
		endYear    = flag.Int("end", 2023, "Specify the end year.")
		outputDir  = flag.String("output", "./output", "Specify the output directory.")
		delay      = flag.Int("delay", 500, "Specify the delay (ms) between downloads.")
	)

	flag.Parse()

	delayDuration := time.Duration(*delay) * time.Millisecond

	downloadReports(*startYear, *endYear, *reportType, *outputDir, delayDuration)
}

func generateURLs(reportType string, start, end int) []string {
	var urls []string

	for i := start; i <= end; i++ {
		// TODO
		if reportType == "legacy" {
			url := GetLegacyArchiveURLForYear(i)
			urls = append(urls, url)
		}
	}

	return urls
}

func GetLegacyArchiveURLForYear(year int) string {
	if year <= 2003 {
		return fmt.Sprintf("https://www.cftc.gov/files/dea/history/deafut_xls_%d.zip", year)
	} else {
		return fmt.Sprintf("https://www.cftc.gov/files/dea/history/dea_fut_xls_%d.zip", year)
	}
}

func downloadReports(start int, end int, reportType string, outputDir string, delay time.Duration) {
	urls := generateURLs(reportType, start, end)

	// Convert to absolute
	absoluteOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		fmt.Println("Error converting path to absolute:", err)
		return
	}

	// Check if directory exists
	if _, err := os.Stat(absoluteOutputDir); os.IsNotExist(err) {
		// Create directory if it doesn't exist
		err := os.MkdirAll(absoluteOutputDir, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	} else if err != nil {
		fmt.Println("Error checking directory:", err)
		return
	}

	err = downloader.DownloadAndExtractAll(urls, delay, absoluteOutputDir)
	if err != nil {
		fmt.Printf("error downloading and extracting: %v\n", err)
	}
}
