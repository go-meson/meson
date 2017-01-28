package util

import "os"
import "path/filepath"

func getApplicationAssetsPath(bundlePath string) string {
	if bundlePath == "" {
		//TODO: change os.Executable in Go1.8
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		return filepath.Join(wd, "assets")
	}
	return filepath.Join(bundlePath, "Contents", "Resources", "assets")
}
