package main

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"fmt"
	"log"
	"os"
)

//Image struct is used to capture the value of an image
type Image struct {
	path        string
	name        string
	hash        string
	resWidthPx  int
	resHeightPx int
	dateTaken   string
}

//Images contain multiple instance of Image objects
type Images struct {
	imageList []Image
}

func (images *Images) append(image *Image) {
	images.imageList = append(images.imageList, *image)
}

func scanZip(zipPath string) {
	var images Images
	zipReader, err := zip.OpenReader(zipPath)

	defer func() {
		if err := zipReader.Close(); err != nil {
			panic(err)
		}
	}()

	if err != nil {
		log.Fatalf("Could not open zip file " + zipPath)
	}

	for _, file := range zipReader.File {
		image, err := imageFromFile(*file)

		if err != nil {
			continue
		}

		(&images).append(image)
	}
}

func imageFromFile(file zip.File) (img *Image, err error) {
	fileRdr, err := file.Open()
	defer fileRdr.Close()

	if err != nil {
		log.Printf("Could not open file " + file.Name)
		return nil, fmt.Errorf("Could not open file %s in the archive ", file.Name)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(fileRdr)
	log.Printf("Scanning file %s with md5 checksum %x", file.Name, md5.Sum(buf.Bytes()))

	return &Image{name: file.Name}, nil
}

func main() {
	scanZip(os.Args[1])
}
