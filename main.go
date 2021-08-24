package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/openshift/assisted-image-service/pkg/isoeditor"
)

const isoFile = "/home/angus/go/src/github.com/asalkeld/custom-image-handler/cached-rhcos-49.84.202107010027-0-openstack.x86_64.qcow2"

// imageFileSystem is an http.FileSystem that creates a virtual filesystem of
// host images. These *could* be later cached as real files.
type imageFileSystem struct {
	images []*imageFile
}

// imageFile is the http.File use in imageFileSystem.
// It is used to wrap the Readdir method of http.File so that we can
// only show images that we are aware of.
type imageFile struct {
	io.ReadSeekCloser
	name              string
	ignitionContent   []byte
	rhcosStreamReader io.ReadSeeker
}

// Readdir is a wrapper around the Readdir method of the embedded File
// that filters out all files that start with a period in their name.
func (f imageFileSystem) Readdir(n int) ([]fs.FileInfo, error) {
	result := []fs.FileInfo{}
	for _, im := range f.images {
		result = append(result, im)
	}
	fmt.Println("ReadDir", result)
	return result, nil
}

func (f imageFileSystem) Open(name string) (http.File, error) {
	fmt.Println("Open", name)
	if name == "/" {
		return nil, fs.ErrPermission
	}
	err := fs.ErrNotExist
	for _, im := range f.images {
		if "/"+im.name == name {
			im.rhcosStreamReader, err = isoeditor.NewRHCOSStreamReader(isoFile, im.ignitionContent)
			if err != nil {
				fmt.Print(err)
				return nil, err
			}
			return im, nil
		}
	}
	return nil, err
}

func main() {
	fsys := imageFileSystem{
		images: []*imageFile{
			{
				name:            "host-image-3.qcow",
				ignitionContent: []byte("someignitioncontent"),
			},
		},
	}
	fmt.Println("main", fsys.images[0])

	http.Handle("/", http.FileServer(fsys))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
