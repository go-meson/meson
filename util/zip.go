package util

import "archive/zip"
import "os"
import "io"
import "io/ioutil"
import "path/filepath"

// ExtractZipFile extract files in the specified zip file
func ExtractZipFile(src string, dst string) error {
	sz, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer sz.Close()
	for _, f := range sz.File {
		name := f.Name
		mode := f.Mode()
		t := mode & os.ModeType
		path := filepath.Join(dst, name)
		if (t & os.ModeDir) != 0 {
			os.MkdirAll(path, f.Mode())
			continue
		}
		if (t & os.ModeSymlink) != 0 {
			err := func() error {
				sr, err := f.Open()
				if err != nil {
					return err
				}
				defer sr.Close()
				data, err := ioutil.ReadAll(sr)
				if err != nil {
					return err
				}
				return os.Symlink(string(data), path)
			}()
			if err != nil {
				return err
			}
			continue
		}
		os.MkdirAll(filepath.Dir(path), f.Mode())
		if err := func() error {
			rr, err := f.Open()
			if err != nil {
				return err
			}
			defer rr.Close()
			ff, err := os.OpenFile(path,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				f.Mode())
			if err != nil {
				return err
			}
			defer ff.Close()
			_, err = io.Copy(ff, rr)
			return err
		}(); err != nil {
			return err
		}
	}
	return nil
}
