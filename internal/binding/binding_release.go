// +build !framework_debug

package binding

import (
	"github.com/go-meson/meson/provision"
)

func resolveFrameworkPath() (string, error) {
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
