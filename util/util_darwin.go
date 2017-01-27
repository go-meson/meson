package util

import "os"
import "path/filepath"

func getApplicationAssetsPath(bundlePath string) string {
	if bundlePath == "" {
		bundlePath = filepath.Base(os.Args[0])
	}
	return filepath.Join(bundlePath, "Contents", "Resources", "assets")
}
