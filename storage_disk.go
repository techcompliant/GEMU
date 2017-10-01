package gemu // import "gemu.techcompliant.com/gemu"

import (
	"io"
	"os"
	"path/filepath"
)

type DiskStorage struct {
	basepath string
}

func NewDiskStorage(basepath string) Storage {
	return &DiskStorage{basepath: basepath}
}

func (DS *DiskStorage) Exists(Item string) bool {
	_, err := os.Stat(filepath.Join(DS.basepath, Item))
	if err == nil {
		return true
	}
	return false
}

func (DS *DiskStorage) Length(Item string) int {
	s, err := os.Stat(filepath.Join(DS.basepath, Item))
	if err != nil {
		return 0
	}
	return int(s.Size())
}

func (DS *DiskStorage) Read(Item string, offset int, data []byte) {
	if DS.Exists(Item) {
		file, err := os.Open(filepath.Join(DS.basepath, Item))
		if err != nil {
			return
		}
		defer file.Close()
		_, err = file.Seek(int64(offset), os.SEEK_SET)
		if err != nil {
			return
		}
		_, err = io.ReadFull(file, data)
		if err != nil {
			return
		}
	} else {
	}
}

func (DS *DiskStorage) Write(Item string, offset int, data []byte) {
	file, err := os.OpenFile(filepath.Join(DS.basepath, Item), os.O_CREATE | os.O_WRONLY, 0777)
	if err != nil {
		return
	}
	defer file.Close()
	file.Seek(int64(offset), os.SEEK_SET)
	file.Write(data)
}
