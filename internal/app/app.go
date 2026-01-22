package app

import (
	"Dogs/internal/api"
	"Dogs/pkg/processor"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type App struct {
	apiClient *api.DogClient
	manual    *processor.ManualProcessor
	library   *processor.LibraryProcessor
}

func New() *App {
	return &App{
		apiClient: api.NewDogClient(),
		manual:    &processor.ManualProcessor{},
		library:   &processor.LibraryProcessor{},
	}
}

func (a *App) Run(limit int) error {
	items, err := a.apiClient.GetDogs(limit)
	if err != nil {
		return fmt.Errorf("ошибка получения данных из API: %w", err)
	}

	fmt.Printf("Получено ссылок: %d. Начинаю параллельную обработку...\n", len(items))

	jobs := make(chan api.DogItem, len(items))
	var wg sync.WaitGroup

	for w := 0; w < limit; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for item := range jobs {
				a.processItem(item, w)
			}
		}(w)
	}

	for _, item := range items {
		jobs <- item
	}
	close(jobs)

	wg.Wait()
	fmt.Println("--- Все задачи выполнены ---")
	return nil
}

func (a *App) processItem(item api.DogItem, number int) {

	fmt.Printf("[Загрузка] Номер: %d, URL: %s\n", number, item.URL)

	img, err := a.downloadImage(item.URL)
	if err != nil {
		log.Printf("Ошибка загрузки %s: %v", item.URL, err)
		return
	}

	fileName := fmt.Sprintf("results/%d_%s_original.png", number, item.Id)
	if err := a.saveImage(fileName, img); err != nil {
		log.Printf("Ошибка сохранения %s: %v", fileName, err)
	} else {
		fmt.Printf("[Готово] Сохранено: %s\n", fileName)
	}

	processedImg := a.manual.EdgeDetection(img)

	fileName = fmt.Sprintf("results/%d_%s_manual.png", number, item.Id)
	if err := a.saveImage(fileName, processedImg); err != nil {
		log.Printf("Ошибка сохранения %s: %v", fileName, err)
	} else {
		fmt.Printf("[Готово] Сохранено: %s\n", fileName)
	}

	processedImg = a.library.EdgeDetection(img)

	fileName = fmt.Sprintf("results/%d_%s_lib.png", number, item.Id)
	if err := a.saveImage(fileName, processedImg); err != nil {
		log.Printf("Ошибка сохранения %s: %v", fileName, err)
	} else {
		fmt.Printf("[Готово] Сохранено: %s\n", fileName)
	}
}

func (a *App) downloadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("статус ответа: %d", resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	return img, err
}

func (a *App) saveImage(path string, img image.Image) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return jpeg.Encode(f, img, nil)
}
