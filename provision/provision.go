package provision

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/go-meson/meson/util"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
)

const mesonFrameworkReleaseURL = "https://github.com/go-meson/framework/releases/download/"

func cacheRootPath() string {
	return filepath.Join(util.GetSystemDirectoryPath(util.UserCacheDirectory),
		"com.github.go-meson")
}

func downloadCachePath() string {
	return filepath.Join(cacheRootPath(), "download")
}

func FrameworkBasePath(version string) string {
	return filepath.Join(cacheRootPath(), version)
}

//https://github.com/go-meson/framework/releases/download/v0.0.1/meson-v0.0.1-darwin-x64.zip
func frameworkURL(version string, suffix string) (string, error) {
	var platform string
	var arch string
	switch runtime.GOOS {
	case "windows":
		platform = "win32"
	case "linux":
		platform = "linux"
	case "darwin":
		platform = "darwin"
	default:
		return "", errors.New("unknown platform")
	}
	switch runtime.GOARCH {
	case "386":
		arch = "ia32"
	case "amd64":
		arch = "x64"
	default:
		return "", fmt.Errorf("unknown architecture: %s", runtime.GOARCH)
	}
	return mesonFrameworkReleaseURL + fmt.Sprintf("%s/meson-%s-%s-%s%s", version, version, platform, arch, suffix), nil
}

func downloadSha256(version string) ([]byte, error) {
	sumURL, err := frameworkURL(version, ".zip.sha256sum")
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(sumURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d(%s)", resp.StatusCode, resp.Status)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(bytes) != 64 {
		return nil, fmt.Errorf("invalid sha256 file length : %d", len(bytes))
	}
	sha256 := make([]byte, 32)
	for i := 0; i < 32; i++ {
		val, err := strconv.ParseUint(string(bytes[i*2:i*2+2]), 16, 32)
		if err != nil {
			return nil, err
		}
		sha256[i] = byte(val)
	}
	return sha256, nil
}

func makeDirIfNeed(path string) error {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.IsDir() {
			// already exists
			return nil
		}
		err = os.Remove(path)
		if err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(path, os.ModePerm)
}

func downloadFramework(version string) (string, error) {
	downloadPath := filepath.Join(downloadCachePath(), version)
	downloadURL, err := frameworkURL(version, ".zip")
	_, fileName := path.Split(downloadURL)
	filePath := filepath.Join(downloadPath, fileName)
	stat, err := os.Stat(filePath)
	if err == nil && !stat.IsDir() {
		// already cached
		return filePath, nil
	}

	err = makeDirIfNeed(downloadPath)

	if err != nil {
		return "", err
	}
	sum, err := downloadSha256(version)
	if err != nil {
		return "", err
	}
	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status code: %d(%s)", resp.StatusCode, resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	dataSum := sha256.Sum256(data)
	for i := 0; i < len(dataSum); i++ {
		if sum[i] != dataSum[i] {
			return "", errors.New("checksum error")
		}
	}
	err = ioutil.WriteFile(filePath, data, os.ModePerm)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// FetchFramework fetch meson framework
//
// FetchFramework() download meson-framework from github release and cache your machine.
func FetchFramework(version string) error {
	rootPath := GetFrameworkRootPath(version)
	if stat, err := os.Stat(rootPath); err == nil && stat.IsDir() {
		return nil
	}
	basePath := FrameworkBasePath(version)
	if err := makeDirIfNeed(basePath); err != nil {
		return nil
	}
	zipPath, err := downloadFramework(version)
	if err != nil {
		return err
	}
	if err = util.ExtractZipFile(zipPath, basePath); err != nil {
		return err
	}

	return nil
}
