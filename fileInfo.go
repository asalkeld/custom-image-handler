package main

import (
	"io/fs"
	"time"
)

var _ fs.FileInfo = &imageFile{}

func (i *imageFile) Name() string {
	return i.name
}

func (i *imageFile) Size() int64 {
	return 1109524480 // TODO we don't know the size before opening the stream..
}

func (i *imageFile) Mode() fs.FileMode {
	return 0444
}

func (i *imageFile) ModTime() time.Time {
	return time.Now()
}

func (i *imageFile) IsDir() bool {
	return false
}

func (i *imageFile) Sys() interface{} {
	return nil
}
