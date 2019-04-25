package main

import (
	"image"
	"image/jpeg"
	"os"
	"sync"
)

type cmpImg struct {
	pixels     [][]pixel
	xMax, yMax int

	sync.Mutex
}

func newCmpImage() *cmpImg {
	return &cmpImg{}
}

func (m *cmpImg) init(xMax, yMax int) {
	m.Lock()
	defer m.Unlock()

	pixels := make([][]pixel, yMax)
	for y := 0; y < yMax; y++ {
		pixels[y] = make([]pixel, xMax)
		for x := 0; x < xMax; x++ {
			pixels[y][x] = newPixelAverage()
		}
	}
	m.pixels = pixels

	m.xMax = xMax
	m.yMax = yMax
}

func (m *cmpImg) isZero() bool {
	m.Lock()
	defer m.Unlock()

	if m.xMax == 0 && m.yMax == 0 && m.pixels == nil {
		return true
	}

	return false
}

func (m *cmpImg) AddImage(img image.Image) {
	yMax := img.Bounds().Max.Y
	xMax := img.Bounds().Max.X

	if m.isZero() {
		m.init(xMax, yMax)
	}

	for y := 0; y < yMax; y++ {
		for x := 0; x < xMax; x++ {
			m.pixels[y][x].Update(img.At(x, y))
		}
	}
}

func (m *cmpImg) Image() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, m.xMax, m.yMax))

	for y := 0; y < m.yMax; y++ {
		for x := 0; x < m.xMax; x++ {
			img.Set(x, y, m.pixels[y][x].Color())
		}
	}

	return img
}

func (m *cmpImg) Save(path string, quality int) error {
	m.Lock()
	defer m.Unlock()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return jpeg.Encode(file, m.Image(), &jpeg.Options{Quality: quality})
}
