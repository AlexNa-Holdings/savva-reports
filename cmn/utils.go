package cmn

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
)

func LoadImage(cid string) (image.Image, error) {
	content := C.IPFS(cid)
	if content == nil {
		return nil, fmt.Errorf("failed to load content for post %s", cid)
	}

	img, _, err := image.Decode(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image for post %s", cid)
	}

	img = EnsureRGBA(img)
	return img, nil
}

// EnsureRGBA checks if the image is already an 8-bit RGBA-compatible image,
// and converts it only if needed.
func EnsureRGBA(src image.Image) *image.RGBA {
	// If already *image.RGBA, return as-is
	if rgba, ok := src.(*image.RGBA); ok {
		return rgba
	}

	// If already *image.NRGBA (which Go also handles fine), convert safely
	if nrgba, ok := src.(*image.NRGBA); ok {
		bounds := nrgba.Bounds()
		dst := image.NewRGBA(bounds)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				dst.Set(x, y, nrgba.At(x, y))
			}
		}
		return dst
	}

	// If other formats (possibly 16-bit or unsupported), convert pixel-by-pixel
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			dst.Set(x, y, color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			})
		}
	}
	return dst
}
