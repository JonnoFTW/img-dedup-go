package main

import (
	"flag"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"ssim/pkg/hash"
	"strings"
	"sync"
)

var (
	directoryArg  = flag.String("directory", "", "directory path")
	hashMethodArg = flag.String("hashMethod", "Average", "Hash method, defaults to phash")
	hashMethods   = map[string]hash.HashMethod{
		"Average":    hash.Average,
		"Perceptual": hash.Perceptual,
		"Difference": hash.Difference,
	}
)

type ImagePath string
type HashedImage struct {
	path ImagePath
	hash hash.ImageHash
}
type Duplicates map[hash.ImageHash][]ImagePath

func getHash(path ImagePath, hashMethod hash.HashMethod, c chan<- *HashedImage, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
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
	imgHash := hash.Hash(&img, hashMethod)
	hashedImage := HashedImage{
		path,
		imgHash,
	}
	bar.Add(1)
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
	if readBytes == 4 && (reflect.DeepEqual(buff, hash.JPEG_MAGIC_BYTES) || reflect.DeepEqual(buff, hash.JPEG_EXIF_MAGIC_BYTES) || reflect.DeepEqual(buff, hash.PNG_MAGIC_BYTES)) {
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
		if !info.IsDir() && (strings.HasSuffix(strings.ToLower(path), "png") || strings.HasSuffix(strings.ToLower(path), "jpg")) {
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
func findDuplicates(directory string, hashMethod hash.HashMethod) {
	// recursively iterate the folder and find all images
	duplicates := make(Duplicates, 0)
	imagesChannel := make(chan *HashedImage)
	wg := new(sync.WaitGroup)
	images := findImages(directory, wg)
	wg.Add(len(images))
	bar := progressbar.Default(int64(len(images)))
	for _, img := range images {
		go getHash(img, hashMethod, imagesChannel, wg, bar)
	}
	go monitorWorker(wg, imagesChannel)
	for imgHash := range imagesChannel {
		if _, ok := duplicates[imgHash.hash]; !ok {
			duplicates[imgHash.hash] = make([]ImagePath, 0)
		}
		duplicates[imgHash.hash] = append(duplicates[imgHash.hash], imgHash.path)
	}
	fmt.Println("Potential Duplicates:")
	for imgHash, duplicateImages := range duplicates {
		if len(duplicateImages) > 1 {
			fmt.Printf("Hash: %064b\n", imgHash)
			for _, path := range duplicateImages {
				fmt.Println("\t path=", path)
			}
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
		log.Fatalf("Invalid hash method '%s' must be one of %s", *hashMethodArg, hashMethods)
	}
	fmt.Println("Detecting duplicates in", *directoryArg, "with method", *hashMethodArg)
	findDuplicates(*directoryArg, hashMethod)
}
