package downloader

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func DownloadAndExtractAll(urls []string, delay time.Duration, outputDir string) error {
	for _, url := range urls {
		if err := downloadAndExtract(url, outputDir); err != nil {
			return err
		}
		time.Sleep(delay)
	}
	return nil
}

func downloadAndExtract(url, outputDir string) error {
	fmt.Println("Downloading file:", url)
	zipPath, err := downloadFile(url)
	if err != nil {
		return err
	}
	if err := unzipDownload(zipPath, outputDir); err != nil {
		fmt.Println("Error unzipping file (unzipDownload):", err)
		return err
	}
	return nil
}

// ExtractYear - this is probably going to replace the other
// I need to handle the older files, which have a slightly different URL
// 2003 and earlier has one URL and all from 2004 and on have another, which is what I'm using currently.
func ExtractYear(filename string) (int, error) {
	//r := regexp.MustCompile(`dea_fut_xls_(\d{4}).zip`)
	r := regexp.MustCompile(`dea(?:fut|_fut)_xls_(\d{4}).zip`)
	matches := r.FindStringSubmatch(filename)
	if matches != nil && len(matches) > 1 {
		year, err := strconv.Atoi(matches[1])
		if err != nil {
			fmt.Println("Error converting year to int:", err)
			return 0, err
		}
		return year, nil
	}
	return 0, fmt.Errorf("could not extract year from file name - Invalid file name format")
}

func downloadFile(url string) (tmpFilePath string, err error) {
	urlFileName := path.Base(url)

	tmpFilePath = path.Join(os.TempDir(), urlFileName)
	tmpFile, err := os.Create(tmpFilePath)

	if err != nil {
		return "", err
	}
	defer func() {
		closeErr := tmpFile.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	// Handle HTTP Status Codes
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

func renameExtractedFile(oldName, outputDirPath, year string) error {
	outName := filepath.Join(outputDirPath, year+"_"+filepath.Base(oldName))
	return os.Rename(oldName, outName)
}

func unzipDownload(fp string, dest string) (err error) {
	fmt.Println("Unzipping file:", fp)

	arch, err := zip.OpenReader(fp)
	if err != nil {
		fmt.Println("Error opening zip file (unzipDownload - 90): ", err)
		return err
	}
	defer func() {
		closeErr := arch.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	for _, f := range arch.File {
		filePath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", filePath)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		err = fileInArchive.Close()
		if err != nil {
			return err
		}

		year, err := ExtractYear(fp)
		if err != nil {
			return err
		}
		yearStr := strconv.Itoa(year)

		err = renameExtractedFile(dstFile.Name(), dest, yearStr)
		if err != nil {
			return err
		}

		err = dstFile.Close()
		if err != nil {
			return err
		}
	}

	// Delete temp zip file
	if err := os.Remove(fp); err != nil {
		return fmt.Errorf("removing zip file: %w", err)
	}

	return nil
}
