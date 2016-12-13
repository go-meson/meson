//go:generate go-bindata -o bindata.go info.plist.tpl

package main

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type options struct {
	Help             bool   `short:"h" long:"help" description:"Show help"`
	Output           string `short:"o" long:"output"`
	BundleIdentifier string `short:"b" long:"bundle_identifier"`
	Icon             string `short:"i" long:"icons" description:"Path to a .icns file or a .iconset dir"`
	executable       string
}

func parseOptions(argv []string) options {
	opts := options{}

	p := flags.NewParser(&opts, flags.PrintErrors)
	args, err := p.ParseArgs(argv[1:])
	if err != nil || len(args) != 1 {
		p.WriteHelp(os.Stdout)
		os.Exit(1)
	}
	if opts.Help {
		p.WriteHelp(os.Stdout)
		os.Exit(0)
	}
	opts.executable = args[0]

	return opts
}

func must(err error, info ...interface{}) {
	if err != nil {
		fmt.Println(append(info, err.Error())...)
		os.Exit(1)
	}
}

func copyFile(dst, src string) error {
	st, err := os.Stat(src)
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, buf, st.Mode())
}

func copyTree(dst, src string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		// re-stat the path so that we can tell whether it is a symlink
		info, err = os.Lstat(path)
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targ := filepath.Join(dst, rel)

		switch {
		case info.IsDir():
			return os.Mkdir(targ, 0777)
		case info.Mode()&os.ModeSymlink != 0:
			referent, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(referent, targ)
		default:
			return copyFile(targ, path)
		}
	})
}

func makeMesonDirs() []string {
	cmd := exec.Command("go", "list", "-f", "{{.Dir}}", "github.com/go-meson/meson")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("go list github.com/go-meson/meson failed:\n%s\n\n", string(out))
		os.Exit(1)
	}

	mesonFiles := []string{
		"Meson.framework",
		//"Meson Helper EH.app",
		"Meson Helper.app",
		//"Meson Helper NP.app",
	}
	basedir := strings.TrimSpace(string(out))
	dirs := make([]string, len(mesonFiles))
	for i, file := range mesonFiles {
		fw := filepath.Join(basedir, "dist", file)
		st, err := os.Stat(fw)
		if err != nil {
			fmt.Printf("framework not found at %s: %v\n", fw, err)
			os.Exit(1)
		}
		if !st.IsDir() {
			fmt.Printf("%s is not a directly\n", fw)
			os.Exit(1)
		}
		dirs[i] = fw
	}
	return dirs
}

func main() {
	opts := parseOptions(os.Args)
	var bundleName string
	if opts.Output == "" {
		bundleName = filepath.Base(opts.executable)
		opts.Output = bundleName + ".app"
	} else if !strings.HasSuffix(opts.Output, ".app") {
		fmt.Println("output must end with .app")
		os.Exit(1)
	} else {
		bundleName = strings.TrimSuffix(filepath.Base(opts.Output), ".app")
	}
	if opts.BundleIdentifier == "" {
		opts.BundleIdentifier = bundleName
	}

	// extras for the Info.plist
	extraProps := make(map[string]string)

	dirs := makeMesonDirs()

	tmpBundle, err := ioutil.TempDir("", "")
	must(err)
	must(os.MkdirAll(tmpBundle, 0777))
	fwDst := filepath.Join(tmpBundle, "Contents", "Frameworks")
	must(os.MkdirAll(filepath.Dir(fwDst), 0777))

	for _, dir := range dirs {
		dst := filepath.Join(fwDst, filepath.Base(dir))
		must(os.MkdirAll(filepath.Dir(dst), 0777))
		must(copyTree(dst, dir))
	}
	exeDst := filepath.Join(tmpBundle, "Contents", "MacOs", bundleName)
	must(os.MkdirAll(filepath.Dir(exeDst), 0777))
	must(copyFile(exeDst, opts.executable))

	// TODO: copy icons

	// Write Info.plist
	tpl, err := template.New("info.plist.tpl").Parse(string(MustAsset("info.plist.tpl")))
	must(err)

	plistDst := filepath.Join(tmpBundle, "Contents", "Info.plist")
	w, err := os.Create(plistDst)
	must(err)

	tpl.Execute(w, map[string]interface{}{
		"BundleName":       bundleName,
		"BundleIdentifier": opts.BundleIdentifier,
		"Extras":           extraProps,
	})
	must(w.Close())

	// Write PkgInfo.(APPL????)
	pkginfo := []byte{0x41, 0x50, 0x50, 0x4c, 0x3f, 0x3f, 0x3f, 0x3f}
	pkginfoDst := filepath.Join(tmpBundle, "Contents", "PkgInfo")
	must(ioutil.WriteFile(pkginfoDst, pkginfo, 0777))

	// Delete the bundle.app dir if it already exists
	must(os.RemoveAll(opts.Output))

	// Move the temporary dir to the bundle.app location
	must(os.Rename(tmpBundle, opts.Output))
}
