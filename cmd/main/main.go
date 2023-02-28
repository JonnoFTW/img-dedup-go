package main

import (
	"flag"
	"fmt"
	image "image"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"ssim/pkg/utils"
	"sync"
)

var (
	directory  = flag.String("directory", "", "directory path")
	hashMethod = flag.String("hashMethod", "phash", "Hash method, defaults to phash")
)

type ImagePath string
type HashedImage struct {
	path ImagePath
	hash utils.ImageHash
}
type Duplicates map[utils.ImageHash][]ImagePath

func getHash(path ImagePath, c chan<- *HashedImage, wg *sync.WaitGroup) {
	defer wg.Done()
	f, err := os.Open(string(path))
	if err != nil {
		fmt.Errorf("Failed to read %s", f)
		return
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	// put the image hash on the queue
	if err != nil {
		fmt.Println("Error decoding", path, err.Error())
		return
	}
	hash := utils.Hash(img, utils.Average)
	hashedImage := HashedImage{
		path,
		hash,
	}
	c <- &hashedImage
}
func checkMagic(path string, c chan<- ImagePath, wg *sync.WaitGroup) {
	// Check magic bytes at the start of the file to see if they are png or jpeg
	defer wg.Done()
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer f.Close()
	buff := make([]byte, 4)
	var readBytes int
	readBytes, err = f.Read(buff)
	if readBytes != 8 || reflect.DeepEqual(buff, utils.JPEG_MAGIC_BYTES) || reflect.DeepEqual(buff, utils.PNG_MAGIC_BYTES) {
		c <- ImagePath(path)
	}
}
func monitorWorker[T any](wg *sync.WaitGroup, cs chan T) {
	wg.Wait()
	close(cs)
}
func findImages(directory string, wg *sync.WaitGroup) []ImagePath {
	images := make([]ImagePath, 0)
	imagePathChannel := make(chan ImagePath)

	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		if !info.IsDir() {
			wg.Add(1)
			go checkMagic(path, imagePathChannel, wg)
		}
		return nil
	})
	go monitorWorker(wg, imagePathChannel)
	for path := range imagePathChannel {
		images = append(images, path)
	}
	fmt.Println("Found", len(images), "images")
	return images
}
func findDuplicates(directory string) {
	// recursively iterate the folder and find all images
	duplicates := make(Duplicates, 0)
	imagesChannel := make(chan *HashedImage)
	wg := new(sync.WaitGroup)
	images := findImages(directory, wg)
	wg.Add(len(images))
	for _, img := range images {
		go getHash(img, imagesChannel, wg)
	}
	go monitorWorker(wg, imagesChannel)
	for hash := range imagesChannel {
		if _, ok := duplicates[hash.hash]; !ok {
			duplicates[hash.hash] = make([]ImagePath, 0)
		}
		duplicates[hash.hash] = append(duplicates[hash.hash], hash.path)
	}
	for hash, vals := range duplicates {
		fmt.Println("Hash:", hash)
		for _, path := range vals {
			fmt.Println("\t path=", path)
		}
	}
}
func main() {
	flag.Parse()
	if *directory == "" {
		log.Fatalf("Please supply a directory")

	}
	fmt.Println("Detecting duplicates in", *directory)
	//fmt.Println("Using hash method", *hashMethod)
	findDuplicates(*directory)
}
