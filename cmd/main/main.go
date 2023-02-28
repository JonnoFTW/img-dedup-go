package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"ssim/pkg/utils"
	"sync"
)

var (
	directoryArg  = flag.String("directory", "", "directory path")
	hashMethodArg = flag.String("hashMethod", "Average", "Hash method, defaults to phash")
	hashMethods   = map[string]utils.HashMethod{
		"Average":    utils.Average,
		"Perceptual": utils.Perceptual,
	}
)

type ImagePath string
type HashedImage struct {
	path ImagePath
	hash utils.ImageHash
}
type Duplicates map[utils.ImageHash][]ImagePath

func getHash(path ImagePath, hashMethod utils.HashMethod, c chan<- *HashedImage, wg *sync.WaitGroup) {
	defer wg.Done()
	f, err := os.Open(string(path))
	if err != nil {
		log.Print(fmt.Errorf("failed to read %s", err))
		return
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	// put the image hash on the queue
	if err != nil {
		fmt.Println("Error decoding", path, err.Error())
		return
	}
	hash := utils.Hash(img, hashMethod)
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

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		if !info.IsDir() {
			wg.Add(1)
			go checkMagic(path, imagePathChannel, wg)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("failed to walk directory %s: %s", directory, err.Error())
	}
	go monitorWorker(wg, imagePathChannel)
	for path := range imagePathChannel {
		images = append(images, path)
	}
	fmt.Println("Found", len(images), "images")
	return images
}
func findDuplicates(directory string, hashMethod utils.HashMethod) {
	// recursively iterate the folder and find all images
	duplicates := make(Duplicates, 0)
	imagesChannel := make(chan *HashedImage)
	wg := new(sync.WaitGroup)
	images := findImages(directory, wg)
	wg.Add(len(images))

	for _, img := range images {
		go getHash(img, hashMethod, imagesChannel, wg)
	}
	go monitorWorker(wg, imagesChannel)
	for hash := range imagesChannel {
		if _, ok := duplicates[hash.hash]; !ok {
			duplicates[hash.hash] = make([]ImagePath, 0)
		}
		duplicates[hash.hash] = append(duplicates[hash.hash], hash.path)
	}
	for hash, duplicateImages := range duplicates {
		fmt.Println("Hash:", hash)
		for _, path := range duplicateImages {
			fmt.Println("\t path=", path)
		}
	}
}
func main() {
	flag.Parse()
	if *directoryArg == "" {
		log.Fatalf("Please supply a directory")
	}
	hashMethod, ok := hashMethods[*hashMethodArg]
	if !ok {
		log.Fatalf("Invalid hashmethod %s", *hashMethodArg)
	}
	fmt.Println("Detecting duplicates in", *directoryArg, "with method", *hashMethodArg)
	//fmt.Println("Using hash method", *hashMethod)
	findDuplicates(*directoryArg, hashMethod)
}
