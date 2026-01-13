package processor

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"
)

// ManualProcessor реализует интерфейс ImageProcessor через ручные алгоритмы
type ManualProcessor struct{}

// EdgeDetection реализует поиск границ (Sobel + Threshold) аналогично твоему коду на Python
func (p *ManualProcessor) EdgeDetection(img image.Image) image.Image {
	startTime := time.Now()
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// 1. Конвертация в Grayscale (аналог np.dot с коэффициентами)
	grayMatrix := p.convertToFloatMatrix(img)

	// 2. Ядра Собеля
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

	// 3. Свертка (аналог _convolve_2d_numba)
	gradX := p.convolve2D(grayMatrix, sobelX)
	gradY := p.convolve2D(grayMatrix, sobelY)

	// 4. Величина градиента и порог (Threshold = 99)
	const threshold = 99.0
	resultImg := image.NewGray(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gx := gradX[y][x]
			gy := gradY[y][x]
			// magnitude = sqrt(gx^2 + gy^2)
			magnitude := math.Sqrt(float64(gx*gx + gy*gy))

			if magnitude >= threshold {
				resultImg.SetGray(x, y, color.Gray{Y: 255})
			} else {
				resultImg.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	fmt.Printf("[EdgeDetection] Время выполнения: %v\n", time.Since(startTime))
	return resultImg
}

// convolve2D - ручная реализация свертки (аналог твоих функций на Numba)
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
					// Имитация паддинга "reflect" через ограничение (clamp)
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

// convertToFloatMatrix переводит картинку в матрицу float32 (Grayscale)
func (p *ManualProcessor) convertToFloatMatrix(img image.Image) [][]float32 {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	matrix := make([][]float32, h)

	for y := 0; y < h; y++ {
		matrix[y] = make([]float32, w)
		for x := 0; x < w; x++ {
			r, g, b, _ := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			// Переводим из 0-65535 в 0-255 и применяем веса
			val := 0.299*float32(r/257) + 0.587*float32(g/257) + 0.114*float32(b/257)
			matrix[y][x] = val
		}
	}
	return matrix
}

// CornerDetection и CircleDetection (заглушки для соблюдения интерфейса)
func (p *ManualProcessor) CornerDetection(img image.Image) image.Image { return img }
func (p *ManualProcessor) CircleDetection(img image.Image) image.Image { return img }
