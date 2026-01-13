package processor

import "image"

type ImageProcessor interface {
	EdgeDetection(img image.Image) image.Image
	CornerDetection(img image.Image) image.Image
	CircleDetection(img image.Image) image.Image
}
