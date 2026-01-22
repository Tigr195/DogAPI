package pkg

import "image"

// DogImage инкапсулирует данные о собаке
type DogImage struct {
	Image image.Image
	Breed string
	URL   string
}
