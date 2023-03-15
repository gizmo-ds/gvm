package gvm

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cheggaaa/pb/v3"
)

// Download downloads the Go binary for the specified version.
func Download(version, savepath string, bar *pb.ProgressBar) (string, error) {
	if f, err := os.Stat(savepath); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(savepath, os.ModePerm); err != nil {
				return "", err
			}
		}
	} else if !f.IsDir() {
		return "", fmt.Errorf("%s is not a directory", savepath)
	}

	var extension string
	switch runtime.GOOS {
	case "windows":
		extension = ".zip"
	default:
		extension = ".tar.gz"
	}
	filename := fmt.Sprintf("%s.%s-%s%s", version, runtime.GOOS, runtime.GOARCH, extension)

	savefile := filepath.Join(savepath, filename)
	if _, err := os.Stat(savefile); err == nil {
		if err = os.Remove(savefile); err != nil {
			return "", err
		}
	}

	url := "https://go.dev/dl/" + filename

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", response.StatusCode)
	}
	contentLength := response.ContentLength
	if contentLength < 0 {
		return "", fmt.Errorf("invalid content length")
	}

	file, err := os.Create(savefile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if bar != nil {
		pr := bar.SetTotal(contentLength).
			Set("prefix", filename).
			NewProxyReader(response.Body)
		bar.Start()
		_, err = io.Copy(file, pr)
		bar.Finish()
		return savefile, err
	}
	_, err = io.Copy(file, response.Body)
	return savefile, err
}
