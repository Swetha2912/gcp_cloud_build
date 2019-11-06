package helpers

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ConcurrentSlice -- threadsafe concurrent slice avoiding race conditions
type ConcurrentSlice struct {
	sync.RWMutex
	Items []interface{}
}

// ConcurrentSliceItem --
type ConcurrentSliceItem struct {
	Index int
	Value interface{}
}

// Append -- append to concurrentslice created to avoid race conditions
func (cs *ConcurrentSlice) Append(item interface{}) {
	cs.Lock()
	defer cs.Unlock()

	cs.Items = append(cs.Items, item)
}

// Iter -- concurrentslice interator with range support like a normal slice
func (cs *ConcurrentSlice) Iter() <-chan ConcurrentSliceItem {
	c := make(chan ConcurrentSliceItem)

	f := func() {
		cs.Lock()
		defer cs.Unlock()
		for index, value := range cs.Items {
			c <- ConcurrentSliceItem{index, value}
		}
		close(c)
	}
	go f()

	return c
}

// MakeDirAll --
func MakeDirAll(path string) {
	newpath := filepath.Join(".", path)
	os.MkdirAll(newpath, os.ModePerm)
}

// WriteFile will take file_name and content and create a file
func WriteFile(filePath string, content string) error {
	err := ioutil.WriteFile(filePath, []byte(content), 777)
	return err
}

// ZipDir -- Zip folder
func ZipDir(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}
