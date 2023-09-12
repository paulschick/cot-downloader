package service

import (
	"fmt"
	"github.com/paulschick/cot-downloader/pkg/downloader"
	"github.com/paulschick/cot-downloader/pkg/urlgen"
	"path/filepath"
	"time"
)

func DownloadReports(start int, end int, reportType string, outputDir string, delay time.Duration) {
	urls := urlgen.GenerateURLs(reportType, start, end)

	// Convert to absolute
	absoluteOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		fmt.Println("Error converting path to absolute:", err)
		return
	}

	for _, url := range urls {
		if err := DownloadAndExtract(url, absoluteOutputDir); err != nil {
			fmt.Printf("error downloading and extracting: %v\n", err)
		}
		time.Sleep(delay)
	}
}

func DownloadAndExtract(url, outputDir string) error {
	fmt.Println("Downloading file:", url)
	zipPath, err := downloader.DownloadFile(url)
	if err != nil {
		return err
	}
	if err := downloader.ExtractDownload(zipPath, outputDir); err != nil {
		fmt.Println("Error unzipping file (unzipDownload):", err)
		return err
	}
	return nil
}
