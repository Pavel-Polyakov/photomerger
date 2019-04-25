package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

func addImg(path string, compositeImg *cmpImg) error {
	f, err := os.Open(path)
	defer f.Close()

	if err != nil {
		return err
	}

	img, _ := jpeg.Decode(f)

	compositeImg.AddImage(img)
	return nil
}

func isJPG(file os.FileInfo) bool {
	return strings.Contains(strings.ToLower(file.Name()), ".jpg")
}

func cleanJPG(value string) string {
	return strings.Replace(strings.ToLower(value), ".jpg", "", 1)
}

func main() {
	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores)

	filesFolder := flag.String("folder", "", "Path to JPG images to use")
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

	var tokens = make(chan struct{}, cores)
	var wg sync.WaitGroup
	processed := 0
	img := newCmpImage()

	fmt.Printf("%s: [0/%d]: 0 percent \n", time.Now(), filesLimit)

	for i, file := range files {
		if !isJPG(file) {
			continue
		}

		if i >= filesLimit {
			break
		}

		tokens <- struct{}{}

		path := *filesFolder + "/" + file.Name()

		wg.Add(1)

		go func() {
			defer wg.Done()
			defer func() {
				processed++
				percent := (float64(processed) / float64(filesLimit)) * 100
				fmt.Printf("%s: [%d/%d]: %f percent \n", time.Now(), processed, filesLimit, percent)
				<-tokens
			}()

			if err := addImg(path, img); err != nil {
				fmt.Printf("%s: err is %s\n", time.Now(), err)
			}

		}()
	}

	wg.Wait()

	finalImageName := strings.ToUpper(cleanJPG(files[0].Name()) + "-" + cleanJPG(files[filesLimit].Name()) + ".jpg")

	fmt.Printf("%s: [%s]: saving\n", time.Now(), finalImageName)

	if err := img.Save(finalImageName, *quality); err != nil {
		panic(err)
	}
}
