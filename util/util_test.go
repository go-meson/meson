package util

import (
	"os"
	"testing"
)

func TestSystemDirectory(t *testing.T) {
	t.Log(GetSystemDirectoryPath(UserCacheDirectory))
	t.Log(GetSystemDirectoryPath(DocumentDirectory))
	t.Log(GetSystemDirectoryPath(DesktopDirectory))
	t.Log(os.TempDir())
}
