package main

import (
	"flag"
	"fmt"
	"github.com/paulschick/cot-downloader/pkg/service"
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

	service.DownloadReports(*startYear, *endYear, *reportType, *outputDir, delayDuration)
}
