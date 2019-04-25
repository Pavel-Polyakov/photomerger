package main

import (
	"image/color"
)

type pixel interface {
	Update(c color.Color)
	Color() color.Color
}

type pixelAverage struct {
	R, G, B, A uint32
	N          uint32
}

func newPixelAverage() *pixelAverage {
	return &pixelAverage{}
}

func (p *pixelAverage) Update(c color.Color) {
	p.N += 1

	r, g, b, a := c.RGBA()

	p.R = p.R + r
	p.G = p.G + g
	p.B = p.B + b
	p.A = p.A + a
}

func (p *pixelAverage) Color() color.Color {
	return &color.RGBA64{
		R: uint16(p.R / p.N),
		G: uint16(p.G / p.N),
		B: uint16(p.B / p.N),
		A: uint16(p.A / p.N),
	}
}

type pixelMax struct {
	R, G, B, A uint32
}

func newPixelMax() *pixelMax {
	return &pixelMax{}
}

func (p *pixelMax) Update(c color.Color) {
	r, g, b, a := c.RGBA()

	if r > p.R {
		p.R = r
	}
	if g > p.G {
		p.G = g
	}
	if b > p.B {
		p.B = b
	}
	if a > p.A {
		p.A = a
	}
}

func (p *pixelMax) Color() color.Color {
	return &color.RGBA64{
		R: uint16(p.R),
		G: uint16(p.G),
		B: uint16(p.B),
		A: uint16(p.A),
	}
}
