package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Extract(file, dest string) error {
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", dest, err)
	}
	switch {
	case strings.HasSuffix(file, ".zip"):
		return unzip(file, dest)
	case strings.HasSuffix(file, ".tar.gz"), strings.HasSuffix(file, ".tgz"):
		return untar(file, dest)
	default:
		return fmt.Errorf("unsupported file type: %s", file)
	}
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		err = func(f *zip.File, dest string) error {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			name := strings.TrimPrefix(f.Name, "go/")
			fi := f.FileInfo()
			if fi.IsDir() {
				if name != "" {
					if err = os.Mkdir(filepath.Join(dest, name), f.Mode()); err != nil && !errors.Is(err, os.ErrExist) {
						return err
					}
				}
				return nil
			}
			if err = os.MkdirAll(filepath.Join(dest, filepath.Dir(name)), fi.Mode()); err != nil {
				return err
			}
			_f, err := os.OpenFile(filepath.Join(dest, name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer _f.Close()
			_, err = io.Copy(_f, rc)
			return err
		}(f, dest)
		if err != nil {
			return err
		}
	}
	return nil
}

func untar(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	var r io.ReadCloser = f
	if strings.HasSuffix(src, ".gz") || strings.HasSuffix(src, ".tgz") {
		if r, err = gzip.NewReader(f); err != nil {
			return err
		}
	}
	tr := tar.NewReader(r)

	for {
		h, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		name := strings.TrimPrefix(h.Name, "go/")

		switch h.Typeflag {
		case tar.TypeDir:
			if name != "" {
				if err = os.Mkdir(filepath.Join(dest, name), os.FileMode(h.Mode)); err != nil && !errors.Is(err, os.ErrExist) {
					return err
				}
			}

		case tar.TypeReg:
			file, err := os.Create(filepath.Join(dest, name))
			if err != nil {
				return err
			}
			defer file.Close()
			if _, err := io.Copy(file, tr); err != nil {
				return err
			}
			if err := file.Chmod(os.FileMode(h.Mode)); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unsupported type flag %v for file %q", h.Typeflag, h.Name)
		}
	}
	return nil
}
