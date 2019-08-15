package commands

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path"
)

func Extract(filePath, extractPath, fileName string) error {
	file, err := os.Open(path.Join(filePath, fileName))
	if err != nil {
		return (err)
	}
	uncompressedStream, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(uncompressedStream)

	for {
		entry, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		p := path.Join(extractPath, path.Dir(fileName), entry.Name)
		if entry.Typeflag == tar.TypeReg {
			ow, err := writeFile(p, os.FileMode(entry.Mode))
			defer ow.Close()
			if err != nil {
				return err
			}
			if _, err := io.Copy(ow, tarReader); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeFile(mpath string, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(mpath, os.O_RDWR|os.O_TRUNC, perm)
	if err != nil {
		err = os.MkdirAll(path.Dir(mpath), perm)
		if err != nil {
			return f, err
		}
		f, err = os.Create(mpath)
		if err != nil {
			return f, err
		}
	}
	return f, nil
}
