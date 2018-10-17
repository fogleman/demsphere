package demsphere

import (
	"image"
	"image/draw"
	"math"

	"github.com/fogleman/fauxgl"
)

type Texture struct {
	W   int
	H   int
	Pix []float64
}

func NewTexture(im image.Image) *Texture {
	gray := ensureGray16(im)
	w := gray.Bounds().Size().X
	h := gray.Bounds().Size().Y
	data := gray16ToFloat64s(gray)
	return &Texture{w, h, data}
}

func (t *Texture) BilinearSample(u, v float64) float64 {
	u -= math.Floor(u)
	v -= math.Floor(v)
	x := u * float64(t.W)
	y := v * float64(t.H)
	x0 := int(x)
	y0 := int(y)
	x1 := x0 + 1
	y1 := y0 + 1
	x -= float64(x0)
	y -= float64(y0)
	if x0 >= t.W {
		x0 = 0
	}
	if y0 >= t.H {
		y0 = 0
	}
	if x1 >= t.W {
		x1 = 0
	}
	if y1 >= t.H {
		y1 = 0
	}
	var d float64
	d += t.Pix[x0+y0*t.W] * ((1 - x) * (1 - y))
	d += t.Pix[x0+y1*t.W] * ((1 - x) * y)
	d += t.Pix[x1+y0*t.W] * (x * (1 - y))
	d += t.Pix[x1+y1*t.W] * (x * y)
	return d
}

func (t *Texture) SphericalSample(spherical fauxgl.Vector) float64 {
	lat := math.Acos(spherical.Z)
	lng := math.Atan2(spherical.Y, spherical.X)
	u := (lng + math.Pi) / (2 * math.Pi)
	v := lat / math.Pi
	return t.BilinearSample(u, v)
}

func (t *Texture) Displace(spherical fauxgl.Vector, lo, hi float64) fauxgl.Vector {
	return spherical.MulScalar(lo + t.SphericalSample(spherical)*(hi-lo))
}

func ensureGray16(im image.Image) *image.Gray16 {
	switch im := im.(type) {
	case *image.Gray16:
		return im
	default:
		dst := image.NewGray16(im.Bounds())
		draw.Draw(dst, im.Bounds(), im, image.ZP, draw.Src)
		return dst
	}
}

func gray16ToFloat64s(im *image.Gray16) []float64 {
	w := im.Bounds().Size().X
	h := im.Bounds().Size().Y
	buf := make([]float64, w*h)
	index := 0
	for y := 0; y < h; y++ {
		i := im.PixOffset(0, y)
		for x := 0; x < w; x++ {
			v := (int(im.Pix[i]) << 8) | int(im.Pix[i+1])
			buf[index] = float64(v) / 0xffff
			index += 1
			i += 2
		}
	}
	return buf
}
