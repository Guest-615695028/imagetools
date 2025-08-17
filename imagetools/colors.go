package imagetools

import (
	"image/color"
	"math/rand"
	"unsafe"
)

func RandomRGBA64() color.RGBA64 {
	u := rand.Uint64()
	return *(*color.RGBA64)(unsafe.Pointer(&u))
}
func RandomRGBA() color.RGBA {
	u := rand.Uint32()
	return *(*color.RGBA)(unsafe.Pointer(&u))
}
func SamePalettes(p, q color.Palette) bool {
	if len(p) != len(q) {
		return false
	}
	for i, v := range p {
		pr, pg, pb, pa := v.RGBA()
		qr, qg, qb, qa := q[i].RGBA()
		if pr != qr || pg != qg || pb != qb || pa != qa {
			return false
		}
	}
	return true
}

func CompareColors(a, b color.Color) int {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return Compare(uint64(ar)<<48|uint64(ag)<<32|uint64(ab)<<16|uint64(aa),
		uint64(br)<<48|uint64(bg)<<32|uint64(bb)<<16|uint64(ba))
}
