package main

import (
	"image/color"
)

type pixel interface {
	Update(c color.Color)
	Color() color.Color
}

type pixelAverage struct {
	r, g, b, a, n uint32
}

func newPixelAverage() *pixelAverage {
	return &pixelAverage{}
}

func (p *pixelAverage) Update(c color.Color) {
	p.n += 1

	r, g, b, a := c.RGBA()

	p.r = p.r + r
	p.g = p.g + g
	p.b = p.b + b
	p.a = p.a + a
}

func (p *pixelAverage) Color() color.Color {
	return &color.RGBA64{
		R: uint16(p.r / p.n),
		G: uint16(p.g / p.n),
		B: uint16(p.b / p.n),
		A: uint16(p.a / p.n),
	}
}

type pixelMax struct {
	r, g, b, a uint32
}

func newPixelMax() *pixelMax {
	return &pixelMax{}
}

func (p *pixelMax) Update(c color.Color) {
	r, g, b, a := c.RGBA()

	if r > p.r {
		p.r = r
	}
	if g > p.g {
		p.g = g
	}
	if b > p.b {
		p.b = b
	}
	if a > p.a {
		p.a = a
	}
}

func (p *pixelMax) Color() color.Color {
	return &color.RGBA64{
		R: uint16(p.r),
		G: uint16(p.g),
		B: uint16(p.b),
		A: uint16(p.a),
	}
}
