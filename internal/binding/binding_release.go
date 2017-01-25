// +build !framework_debug

package binding

import (
	"github.com/go-meson/meson/provision"
	"log"
)

func resolveFrameworkPath() (string, error) {
	log.Printf("version: %s\n", MesonFrameworkVersion())
	err := provision.FetchFramework(MesonFrameworkVersion())
	if err != nil {
		return "", err
	}
	path, err := provision.GetFrameworkPath(MesonFrameworkVersion())
	if err != nil {
		return "", err
	}
	return path, nil
}
