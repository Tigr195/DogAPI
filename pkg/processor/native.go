package processor

import (
	"github.com/anthonynsimon/bild/effect"
	"image"
)

type LibraryProcessor struct{}

func (p *LibraryProcessor) EdgeDetection(img image.Image) image.Image {
	// Библиотечный Собель (аналог твоего EdgeDetection)
	// Параметр 1.0 — это радиус/интенсивность
	return effect.EdgeDetection(img, 1.0)
}

func (p *LibraryProcessor) CornerDetection(img image.Image) image.Image {
	// В bild нет прямого Харриса, но есть сегментация или пороговые фильтры.
	// Оставим пока так, чтобы не усложнять.
	return img
}

func (p *LibraryProcessor) CircleDetection(img image.Image) image.Image {
	return img
}
