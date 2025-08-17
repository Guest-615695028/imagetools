package imagetools

import "math"

type LongDouble struct {
	e [2]uint32
	f [2]uint64
}

func NewInt[I SignedInteger](i I) (z LongDouble) {
	math.Float64bits(0)
	e0, f0 := uint32(0), uint64(Abs(i))
	if i < 0 {
		e0 = 0xc0000000
	} else {
		e0 = 0x40000000
	}
	return LongDouble{e: [2]uint32{e0, 0}, f: [2]uint64{f0, 0}}
}
func NewUint[I UnsignedInteger](i I) (z LongDouble) {
	math.Float64bits(0)
	e0, f0 := uint32(0), uint64(i)
	if i < 0 {
		e0 = 0xc0000000
	} else {
		e0 = 0x40000000
	}
	return LongDouble{e: [2]uint32{e0, 0}, f: [2]uint64{f0, 0}}
}

func (z LongDouble) Normalize() LongDouble {
	for i := range 2 {
		e, f := z.e[i], z.f[i]
		if e<<1 >= 0xFFFFFFFE { // Inf or Nan
			continue
		} else if f == 0 {
			e = 0
		} else {
			for e<<1 != 0 && f<<1 > f {
				f <<= 1
				e--
			}
		}
		z.e[i] = e
	}
	return z
}
