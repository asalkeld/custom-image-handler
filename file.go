package main

import "io/fs"

var _ fs.File = &imageFile{}

func (f *imageFile) Read(p []byte) (n int, err error) {
	return f.rhcosStreamReader.Read(p)
}

func (f *imageFile) Seek(offset int64, whence int) (int64, error) {
	return f.rhcosStreamReader.Seek(offset, whence)
}

func (f *imageFile) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (f *imageFile) Stat() (fs.FileInfo, error) {
	return fs.FileInfo(f), nil
}

func (f *imageFile) Close() error {
	return nil
}

func (f *imageFile) Readdir(count int) ([]fs.FileInfo, error) {
	return []fs.FileInfo{}, nil
}
