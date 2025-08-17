package imagetools

import (
	"fmt"
	types "imagetools/types"
)

/*
Format the matrix

Flags:
  - '#': print the type and dimension information
  - '+': print '+' if the number is not negative
*/
func (m Matrix[T]) Format(f fmt.State, r rune) {
	defer func() { recover() }()
	m.reval()
	if f.Flag('#') {
		fmt.Fprintf(f, "%T(%d,%d)", m, m.x, m.y)
	}
	s := fmt.Sprintf("%"+string(r), T(0))
	if s[0] == '%' && s[1] == '!' {
		f.Write([]byte("[" + s + "]"))
		return
	}
	if m.x <= 0 || m.y <= 0 {
		f.Write([]byte("[]"))
		return
	}
	prec, ok := f.Precision()
	if !ok {
		prec = 6
	}
	wid, _ := f.Width()
	wid = min(wid, m.y)
	if wid > 0 { // multi-line output
		outs := make([][]string, m.y)
		ls := make([]int, wid)
		for i := range m.y {
			outs[i] = make([]string, m.x)
			for j := range m.x {
				s = types.FormatNumber(m.val[i*m.x+j], byte(r), prec)
				outs[i][j] = s
				ls[i%wid] = max(ls[i%wid], len(s))
			}
		}
		for i, row := range outs {
			if i == 0 {
				f.Write([]byte{'['})
			} else if wid < m.y {
				f.Write([]byte{'\n'})
			}
			for j, v := range row {
				if j > 0 {
					f.Write([]byte{','})
				}
				if j%wid == 0 {
					f.Write([]byte{'\n', '\t'})
				}
				f.Write(Fill[byte](' ', ls[j%wid]-len(v)))
				f.Write([]byte(v))
			}
		}
	} else { //one-line output
		for i, v := range m.val {
			if i == 0 {
				f.Write([]byte{'['})
			} else if i%m.x == 0 {
				f.Write([]byte{';'})
			} else {
				f.Write([]byte{','})
			}
			f.Write([]byte(types.FormatNumber(v, byte(r), prec)))
		}
	}
	f.Write([]byte{']'})
}
