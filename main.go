package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

func openImg(path string) (image.Image, error) {
	f, err := os.Open(path)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	return jpeg.Decode(f)
}

func isJPG(v string) bool {
	return strings.Contains(strings.ToLower(v), ".jpg")
}

func cleanJPG(v string) string {
	return strings.Replace(strings.ToLower(v), ".jpg", "", 1)
}

func main() {
	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores)

	filesFolder := flag.String("folder", "", "Path to folder with JPG images to use")
	quality := flag.Int("quality", 100, "Final JPG quality")
	maxFiles := flag.Int("max-files", 0, "Limit files number")
	flag.Parse()

	if *filesFolder == "" {
		panic("folder arguments must be set")
	}

	files, err := ioutil.ReadDir(*filesFolder)
	if err != nil {
		panic(err)
	}

	var filesLimit int
	if *maxFiles != 0 {
		if *maxFiles >= len(files) {
			filesLimit = len(files)
		} else {
			filesLimit = *maxFiles
		}
	}

	var wg sync.WaitGroup
	var tokens = make(chan struct{}, cores)

	cmpImg := newCmpImage()

	processed := 0

	logProgress := func() {
		percent := (float64(processed) / float64(filesLimit)) * 100
		fmt.Printf("%s: [%d/%d]: %f percent \n", time.Now(), processed, filesLimit, percent)
	}

	logProgress()

	for i, file := range files {
		if !isJPG(file.Name()) {
			continue
		}

		if i >= filesLimit {
			break
		}

		wg.Add(1)
		tokens <- struct{}{}

		go func() {
			defer func() {
				wg.Done()
				<-tokens
			}()

			path := *filesFolder + "/" + file.Name()
			img, err := openImg(path)
			if err != nil {
				fmt.Printf("%s: err is %s\n", time.Now(), err)
				return
			}

			cmpImg.AddImage(img)

			processed++
			logProgress()
		}()
	}

	wg.Wait()

	firstImgName := cleanJPG(files[0].Name())
	lastImgName := cleanJPG(files[filesLimit].Name())
	finalImageName := strings.ToUpper(firstImgName+"-"+lastImgName) + ".jpg"

	fmt.Printf("%s: [final]: saving %s\n", time.Now(), finalImageName)
	if err := cmpImg.Save(finalImageName, *quality); err != nil {
		panic(err)
	}
	fmt.Printf("%s: [final]: completed\n", time.Now())
}
