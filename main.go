package main

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"imagetools"
	"os"
	"time"
)

const (
	cameraRadius     = 0.0254 / 6
	pixelWidth       = 2e-6
	pixelLength      = 2e-6
	lensDistance     = 0.0325
	SignalNoiseRatio = 6760.83
)

func main() {
	imgname := "example"
	rgba, _ := OpenRGBA(imgname + ".jpg")
	os.Mkdir(imgname, os.ModeDir)
	names := [2]string{"L", "R"}
	for l, img := range imagetools.Split2(rgba, false) {
		n0 := imgname + "/" + names[l]
		WritePNG(n0, img)
		ms := imagetools.RGBA2Matrices(img)
		edges := ms
		for n, m := range ms {
			ms[n] = imagetools.HistogramizeMatrix(m)
			conv, _ := imagetools.ConvertMatrix[int](m).Conv(imagetools.Laplace12, 1, 1)
			edges[n] = imagetools.Absolutize(conv)
		}
		go WritePNG(n0+"H", imagetools.Matrices2RGB(ms[:]))
		go WritePNG(n0+"E", imagetools.Matrices2RGB(edges[:]))
	}
	time.Sleep(1 * time.Second)
}

func OpenImage(name string) (image.Image, error) {
	f, err := os.Open(name)
	if f == nil {
		return nil, err
	}
	defer f.Close()
	defer func() { recover() }()
	i := len(name)
	for i > 0 {
		if i--; name[i] == '.' {
			break
		}
	}
	switch name[i+1] {
	case 'P', 'p':
		return png.Decode(f)
	case 'G', 'g':
		return gif.Decode(f)
	case 'J', 'j':
		return jpeg.Decode(f)
	}
	m, _, err := image.Decode(f)
	return m, err
}

func OpenRGBA(name string) (*image.RGBA, error) {
	img, err := OpenImage(name)
	return imagetools.RGBA(img), err
}

func WritePNG(name string, im image.Image) error {
	f, err := os.Create(name + ".png")
	if f == nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, im)
}
func WriteJPEG(name string, im image.Image, o jpeg.Options) error {
	f, err := os.Create(name + ".jpeg")
	if f == nil {
		return err
	}
	defer f.Close()
	return jpeg.Encode(f, im, &o)
}
func WriteGIF(name string, im image.Image, o gif.Options) error {
	f, err := os.Create(name + ".gif")
	if f == nil {
		return err
	}
	defer f.Close()
	return gif.Encode(f, im, &o)
}
