package downloader

import (
	"io"
	"os"
	"path"
)

func DownloadFile(url string) (string, error) {
	urlFileName := path.Base(url)
	tempFilePath := path.Join(os.TempDir(), urlFileName)
	tempFile, err := os.Create(tempFilePath)

	if err != nil {
		return "", err
	}

	defer func() {
		closeErr := tempFile.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	resp, err := RequestFile(url)
	if err != nil {
		return "", err
	}

	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return "", err
	}
	return tempFile.Name(), nil
}
