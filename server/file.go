package server

import (
	"io"
	"log"
	"os"
)

const bufferSize = 512 * 1024

var fileUtils FileUtils = &FileUtilsService{}

type FileUtils interface {
	read(fileName string, consumer func(*[]byte)) error
}

type FileUtilsService struct{}

func (f *FileUtilsService) read(fileName string, consumer func(*[]byte)) error {
	file, err := os.OpenFile(fileName, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		log.Println("cannot read response file:", fileName, " => ", err.Error())
		return err
	}

	defer file.Close()
	buf := make([]byte, bufferSize)

	for {
		n, err := file.Read(buf)
		if err == io.EOF || n <= 0 {
			break
		}
		if err == nil {
			data := buf[:n]
			consumer(&data)
		} else {
			log.Println("cannot read from a file:", file, err.Error())
			return err
		}
	}
	return nil
}

type cachableState int

const canBeCached cachableState = 1
const cannotBeCached cachableState = 2

func getCacheableState(filename string, maxCacheableSize int64) cachableState {
	info, err := os.Stat(filename)

	if err != nil {
		log.Println("cannot get file stat to determine if possible to cache: ", filename, err.Error())
		return cannotBeCached
	}

	size := info.Size()
	if size <= maxCacheableSize {
		return canBeCached
	} else {
		return cannotBeCached
	}
}
