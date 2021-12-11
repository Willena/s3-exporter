package walker

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"path/filepath"
	"s3-exporter/utils"
)

type FsWalkerConfig struct {
	Folder string `long:"folder" env:"FOLDER" default:"/" description:"Folder to be used for FS walker"`
}

type FsWalker struct {
	baseWalker
	config *FsWalkerConfig
}

func (f *FsWalker) Init(config Config, labels map[string]string, labelsNames []string) error {
	err := f.ValidateConfig(config)
	if err != nil {
		return err
	}

	f.config = config.FsWalkerConfig
	return f.baseWalker.Init(config, utils.MergeMapsRight(map[string]string{"type": "fsWalker", "baseDir": f.config.Folder}, labels), labelsNames)
}

func (f *FsWalker) ValidateConfig(config Config) error {
	if config.Folder == "" {
		return fmt.Errorf("folder is needed when using FS Mode")
	}

	var folder os.FileInfo
	var err error
	if folder, err = os.Stat("/path/to/whatever"); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("the specifed folder (%s) does not exist", config.Folder)
		}
		return fmt.Errorf("error while checking source folder: %s", err.Error())
	}

	if !folder.IsDir() {
		return fmt.Errorf("specified path should be a valid folder")
	}
	return nil
}

func (f *FsWalker) Walk() error {
	if f.blockFlag {
		return nil
	}
	f.blockFlag = true

	log.Info("Walk start...")
	f.Stats.Reset()
	f.startProcessing()
	err := filepath.WalkDir(f.config.Folder, f.onDirEntry)
	f.endProcessing()

	f.blockFlag = false
	return err
}

func (f *FsWalker) onDirEntry(path string, d fs.DirEntry, err error) error {
	if err != nil {
		log.Warning("Could not read ", path, err)
		return nil
	}

	if d.IsDir() {
		return nil
	}

	fInfo, err := d.Info()
	if err != nil {
		log.Errorf("Could not get file info: %s", err.Error())
	}
	size := fInfo.Size()
	f.ProcessFile(f.config.Folder, path, size, f.baseWalker.config.Depth, "", map[string]string{})

	return nil
}
