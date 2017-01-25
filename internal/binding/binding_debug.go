// +build framework_debug

package binding

import (
	"errors"
	"github.com/go-meson/meson/provision"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func resolveFrameworkPath() (string, error) {
	cmd := exec.Command("go", "list", "-f", "{{.Dir}}", "github.com/go-meson/meson")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	pkgdir := strings.TrimSpace(string(out))
	frameworkDir := filepath.Join(filepath.Dir(pkgdir), "framework", "out", "D")
	stat, err := os.Stat(frameworkDir)
	if err != nil {
		return "", err
	}
	if !stat.IsDir() {
		return "", errors.New("framework is not build yet")
	}
	frameworkDir = provision.GetFrameworkPathFromRootPath(frameworkDir)
	stat, err = os.Stat(frameworkDir)
	if err != nil {
		return "", err
	}
	if !stat.IsDir() {
		return "", errors.New("framework path is broken")
	}
	return frameworkDir, nil
}
