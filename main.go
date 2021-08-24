package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	data, err := ioutil.ReadFile("go.sum")
	if err != nil {
		panic(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// https://releases-art-rhcos.svc.ci.openshift.org/art/storage/releases/rhcos-4.9/49.84.202107010027-0/x86_64/rhcos-49.84.202107010027-0-live.x86_64.iso
	imageFS := NewImageFileSystem(path.Join(cwd, "rhcos-49.84.202107010027-0-live.x86_64.iso"))

	for i := 0; i < 100; i++ {
		imageFS.images = append(imageFS.images, &imageFile{
			name:            fmt.Sprintf("host-image-%d.qcow", i),
			ignitionContent: []byte(data),
		})
	}
	fmt.Println("use below curl command, replacing '<num>' with 1-100")
	fmt.Println("curl -X GET localhost:8080/host-image-<num>.qcow --output host-image-<num>.qcow")

	// why use a FileServer?
	// 1. it streams files efficiently
	// 2. if we cache these images, then that will be an easy change.
	http.Handle("/", http.FileServer(imageFS))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
