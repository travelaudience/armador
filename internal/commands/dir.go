package commands

import (
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/travelaudience/armador/internal/logger"
)

var armadorPath = ".armador"

type Dirs struct {
	Tmp      TmpDirs
	Cache    CacheDirs
	Snapshot string
}

type CacheDirs struct {
	Root      string
	Charts    string
	Overrides string
}

type TmpDirs struct {
	Root      string
	Extracted string
	Overrides string
	Hold      string
}

func CreateDirs() (Dirs, error) {
	cacheDirs := CacheDirs{}
	tmpDirs := TmpDirs{}
	home, err := homedir.Dir()
	if err != nil {
		logger.GetLogger().Errorf("Unable to establish the users home dir: %s", err)
		return Dirs{}, err
	}
	armadorDir, err := createDir(home, armadorPath)
	if err != nil {
		return Dirs{}, err
	}
	tmpDirs.Root, err = createDir(armadorDir, "tmp")
	if err != nil {
		return Dirs{}, err
	}
	cacheDirs.Root, err = createDir(armadorDir, "cache")
	if err != nil {
		return Dirs{}, err
	}
	snapshot, err := createDir(armadorDir, "snapshot")
	if err != nil {
		return Dirs{}, err
	}
	tmpDirs.Extracted, err = createDir(tmpDirs.Root, "extracted")
	if err != nil {
		return Dirs{}, err
	}
	tmpDirs.Overrides, err = createDir(tmpDirs.Root, "overrides")
	if err != nil {
		return Dirs{}, err
	}
	tmpDirs.Hold, err = createDir(tmpDirs.Root, "hold")
	if err != nil {
		return Dirs{}, err
	}
	cacheDirs.Charts, err = createDir(cacheDirs.Root, "charts")
	if err != nil {
		return Dirs{}, err
	}
	cacheDirs.Overrides, err = createDir(cacheDirs.Root, "overrides")
	if err != nil {
		return Dirs{}, err
	}

	dirs := Dirs{Tmp: tmpDirs, Cache: cacheDirs, Snapshot: snapshot}
	return dirs, nil
}

func createDir(parent, newFolder string) (string, error) {
	dir := filepath.Join(parent, newFolder)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		logger.GetLogger().Errorf("Failed to create %s dir: %s", dir, err)
		return "", err
	}
	return dir, nil
}

func CleanDirs() error {
	home, err := homedir.Dir()
	if err != nil {
		logger.GetLogger().Errorf("Unable to establish the users home dir: %s", err)
		return err
	}
	tmpDir := filepath.Join(home, armadorPath, "tmp")
	cacheDir := filepath.Join(home, armadorPath, "cache")

	err = removeContents(tmpDir)
	if err != nil {
		return err
	}
	err = removeContents(cacheDir)
	if err != nil {
		return err
	}
	return nil
}

func removeContents(dir string) error {
	logger.GetLogger().Debugf("clearing out %s", dir)
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		logger.GetLogger().Warnf("Path to '%s' does not exist, will create it now", dir)
	}
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		logger.GetLogger().Errorf("Failed to create %s dir: %s", dir, err)
		return err
	}
	return nil
}
