package main

import (
	"image/color"
)

func ColorHex(c color.Color) string {
	const hex = "0123456789abcdef"
	r, g, b, a := c.RGBA()
	s := [8]byte{}
	d := [4]uint32{r, g, b, a}
	for i, v := range d {
		v >>= 8
		v &= 255
		s[2*i+1], s[2*i+2] = hex[v/16], hex[v%16]
	}
	return string(s[:])
}
