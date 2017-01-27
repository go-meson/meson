package provision

import (
	"errors"
	"github.com/go-meson/meson/util"
	"os"
	"path/filepath"
)

func GetFrameworkRootPath(version string) string {
	return GetFrameworkPathFromRootPath(filepath.Join(FrameworkBasePath(version)))
}

func GetFrameworkPathFromRootPath(rootPath string) string {
	return filepath.Join(rootPath, "Meson.framework")
}

func GetFrameworkPath(version string) (string, error) {
	var frameworkRoot string
	if bundlePath := util.ApplicationBundlePath; bundlePath != "" {
		if stat, err := os.Stat(bundlePath); err != nil || !stat.IsDir() {
			return "", errors.New("invalid application structure")
		}
		frameworkRoot = GetFrameworkPathFromRootPath(filepath.Join(bundlePath, "Contents", "Frameworks"))
	} else {
		// provision path
		err := FetchFramework(version)
		if err != nil {
			return "", err
		}
		frameworkRoot = GetFrameworkRootPath(version)
	}
	if _, err := os.Stat(frameworkRoot); err != nil {
		return "", err
	}
	return frameworkRoot, nil
}
