package processor

import (
	"image"
	"image/color"
	"math"
)

type ManualProcessor struct{}

func (p *ManualProcessor) EdgeDetection(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	grayMatrix := p.convertToFloatMatrix(img)

	sobelX := [][]float32{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}
	sobelY := [][]float32{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	gradX := p.convolve2D(grayMatrix, sobelX)
	gradY := p.convolve2D(grayMatrix, sobelY)

	const threshold = 99.0
	resultImg := image.NewGray(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gx := gradX[y][x]
			gy := gradY[y][x]
			magnitude := math.Sqrt(float64(gx*gx + gy*gy))

			if magnitude >= threshold {
				resultImg.SetGray(x, y, color.Gray{Y: 255})
			} else {
				resultImg.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	return resultImg
}

func (p *ManualProcessor) convolve2D(input [][]float32, kernel [][]float32) [][]float32 {
	h := len(input)
	w := len(input[0])
	kh := len(kernel)
	kw := len(kernel[0])
	padH, padW := kh/2, kw/2

	output := make([][]float32, h)
	for i := range output {
		output[i] = make([]float32, w)
	}

	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			var sum float32
			for ki := 0; ki < kh; ki++ {
				for kj := 0; kj < kw; kj++ {
					ii := i + ki - padH
					jj := j + kj - padW

					if ii < 0 {
						ii = 0
					} else if ii >= h {
						ii = h - 1
					}
					if jj < 0 {
						jj = 0
					} else if jj >= w {
						jj = w - 1
					}

					sum += input[ii][jj] * kernel[ki][kj]
				}
			}
			output[i][j] = sum
		}
	}
	return output
}

func (p *ManualProcessor) convertToFloatMatrix(img image.Image) [][]float32 {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	matrix := make([][]float32, h)

	for y := 0; y < h; y++ {
		matrix[y] = make([]float32, w)
		for x := 0; x < w; x++ {
			r, g, b, _ := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			val := 0.299*float32(r/257) + 0.587*float32(g/257) + 0.114*float32(b/257)
			matrix[y][x] = val
		}
	}
	return matrix
}

func (p *ManualProcessor) CornerDetection(img image.Image) image.Image { return img }
func (p *ManualProcessor) CircleDetection(img image.Image) image.Image { return img }
