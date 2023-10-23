package downloader

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// ExtractYear - this is probably going to replace the other
// 2003 and earlier has one URL and all from 2004 and on have another, which is what I'm using currently.
func ExtractYear(filename string) (int, error) {
	r := regexp.MustCompile(`dea(?:fut|_fut)_xls_(\d{4})\.(?:zip|xls)`)
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

type Rename struct {
	OldPath string
	NewPath string
}

func ExtractDownload(zipFilePath string, destDir string) error {
	arch, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}

	var renameList []Rename
	for _, f := range arch.File {
		filePath := filepath.Join(destDir, f.Name)

		err := ValidFilePath(filePath, destDir)
		if err != nil {
			return err
		}

		err = MakeOutputDir(destDir)
		if err != nil {
			return err
		}

		dstFile, err := CopyArchiveFileToDest(f, filePath)
		if err != nil {
			return err
		}

		newName, err := GetNewFileName(zipFilePath, dstFile)
		if err != nil {
			return err
		}

		renameList = append(renameList, Rename{OldPath: dstFile.Name(), NewPath: newName})

		err = dstFile.Close()
		if err != nil {
			return err
		}
	}

	for _, rename := range renameList {
		err = os.Rename(rename.OldPath, rename.NewPath)
		if err != nil {
			return err
		}
	}

	err = arch.Close()
	if err != nil {
		return err
	}

	// Delete temp zip file
	if err := os.Remove(zipFilePath); err != nil {
		return fmt.Errorf("removing zip file: %w", err)
	}

	return nil
}
