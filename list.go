package gvm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gizmo-ds/gvm/utils"
)

// ListRemote returns the available versions of Go from the official download
// page, filtered by the stable parameter. If the stable parameter is omitted,
// all versions are returned. The returned versions are sorted by their version
// number, from the newest to the oldest.
func ListRemote(stable ...bool) ([]Version, error) {
	versions, err := fetchVersions(true)
	if err != nil {
		return nil, err
	}
	var filteredVersions []Version
	for _, v := range versions {
		var filteredFiles []File
		for _, f := range v.Files {
			if f.Arch != runtime.GOARCH || f.Os != runtime.GOOS {
				continue
			}
			filteredFiles = append(filteredFiles, f)
		}
		if len(filteredFiles) == 0 {
			continue
		}
		v.Files = filteredFiles
		if len(stable) > 0 && v.Stable != stable[0] {
			continue
		}
		filteredVersions = append(filteredVersions, v)
	}
	if len(filteredVersions) == 0 {
		return nil, fmt.Errorf("no matching versions found")
	}
	return filteredVersions, nil
}

// LatestStableVersion returns the latest stable version of the go.
func LatestStableVersion() (Version, error) {
	versions, err := fetchVersions(false)
	if err != nil {
		return Version{}, err
	}
	if len(versions) == 0 {
		return Version{}, fmt.Errorf("no versions found")
	}
	version := versions[0]
	var filteredFiles []File
	for _, f := range version.Files {
		if f.Arch != runtime.GOARCH || f.Os != runtime.GOOS {
			continue
		}
		filteredFiles = append(filteredFiles, f)
	}
	if len(filteredFiles) == 0 {
		return Version{}, fmt.Errorf("no matching files found")
	}
	version.Files = filteredFiles
	return version, nil
}

func fetchVersions(all bool) ([]Version, error) {
	query := url.Values{}
	query.Set("mode", "json")
	if all {
		query.Set("include", "all")
	}
	resp, err := http.Get("https://go.dev/dl/?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}
	var versions []Version
	if err = json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, err
	}
	return versions, nil
}

// ListLocalVersions returns a list of all locally installed versions.
func ListLocalVersions() ([]string, error) {
	dir, err := utils.GvmDir()
	if err != nil {
		return nil, err
	}

	versionsDir := filepath.Join(dir, "versions")

	df, err := os.ReadDir(versionsDir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, f := range df {
		if !f.IsDir() {
			continue
		}
		names = append(names, f.Name())
	}
	return names, nil
}
