package downloader

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ValidFilePath(fp, dest string) error {
	if !strings.HasPrefix(fp, filepath.Clean(dest)+string(filepath.Separator)) {
		return fmt.Errorf("%s: illegal file path", fp)
	}
	return nil
}

func MakeOutputDir(destDir string) error {
	absoluteDestDir, err := filepath.Abs(destDir)
	if err != nil {
		return err
	}

	// Check if directory exists
	if _, err := os.Stat(absoluteDestDir); os.IsNotExist(err) {
		// Create directory if it doesn't exist
		err := os.MkdirAll(absoluteDestDir, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return err
		}
	} else if err != nil {
		fmt.Println("Error checking directory:", err)
		return err
	}
	return nil
}

func RenameExtractedFile(zipFilePath string, file *os.File) (string, error) {
	fp := file.Name()
	fmt.Println("file path: ", zipFilePath)
	parentDir := filepath.Dir(fp)
	year, err := ExtractYear(zipFilePath)
	if err != nil {
		return "", err
	}
	yearStr := strconv.Itoa(year)
	ext := filepath.Ext(fp)
	newName := yearStr + "_annual" + ext
	outName := filepath.Join(parentDir, newName)
	err = os.Rename(fp, outName)
	if err != nil {
		return "", err
	}
	return outName, nil
}

func CopyArchiveFileToDest(zipFile *zip.File, fp string) (*os.File, error) {
	dstFile, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())

	if err != nil {
		return nil, err
	}

	fileInArchive, err := zipFile.Open()
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := fileInArchive.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	if _, err := io.Copy(dstFile, fileInArchive); err != nil {
		return nil, err
	}

	return dstFile, nil
}
