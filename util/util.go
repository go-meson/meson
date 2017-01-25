package util

/*
#cgo darwin CFLAGS: -mmacosx-version-min=10.9
#cgo darwin LDFLAGS: -mmacosx-version-min=10.9
#cgo darwin LDFLAGS: -framework Foundation -framework AppKit

#include <stdlib.h>

extern char* mesonGetBundlePath(void);
extern char* mesonGetSystemDirectoryPath(int);
*/
import "C"

import "unsafe"

// SystemDirectoryType specify the location of variety of directories by the GetSystemFolderPath() function.
type SystemDirectoryType int

const (
	_ SystemDirectoryType = iota
	// UserCacheDirectory is location of discardable cache files for current user
	UserCacheDirectory
	// DocumentDirectory is Document's directory
	DocumentDirectory
	// DesktopDirectory is location of user's desktop directory
	DesktopDirectory
)

func GetApplicationBundlePath() string {
	cstr := C.mesonGetBundlePath()
	if cstr == nil {
		return ""
	}
	str := C.GoString(cstr)
	C.free(unsafe.Pointer(cstr))
	return str
}

// GetSystemDirectoryPath return the specified common directory
func GetSystemDirectoryPath(dir SystemDirectoryType) string {
	cdir := C.int(dir)
	cstr := C.mesonGetSystemDirectoryPath(cdir)
	if cstr == nil {
		return ""
	}
	str := C.GoString(cstr)
	C.free(unsafe.Pointer(cstr))
	return str
}
