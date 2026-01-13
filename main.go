package main

import (
	"fmt"
	"image"
	_ "image/jpeg" // Нужно для декодирования jpg
	"image/png"    // Нужно для сохранения png
	"log"
	"os"
	"time"

	"Dogs/pkg/processor" // Замени "Dogs" на название твоего модуля из go.mod
)

func main() {
	fmt.Println("--- Запуск Dog Image Processor (Manual Mode) ---")

	// 1. Открываем файл
	file, err := os.Open("test_photos/3_shiba_inu_original.png")
	if err != nil {
		log.Fatalf("Не удалось открыть файл: %v", err)
	}

	// 2. Декодируем изображение в объект image.Image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Ошибка декодирования: %v", err)
	}

	err = file.Close()
	if err != nil {
		log.Fatalf("Не удалось закрыть файл: %v", err)
	}

	// 3. Создаем наш ручной процессор
	proc := &processor.ManualProcessor{}

	// --- ТЕСТ 1: ПОИСК ГРАНИЦ ---
	fmt.Println("Обработка: поиск границ...")
	startEdges := time.Now()
	edgeImg := proc.EdgeDetection(img)
	fmt.Printf("Границы найдены за %v\n", time.Since(startEdges))

	saveToFile("results/output_edges_manual.png", edgeImg)

	procCV := &processor.LibraryProcessor{}

	fmt.Println("Обработка: поиск границ...")
	startEdges = time.Now()
	edgeImg = procCV.EdgeDetection(img)
	fmt.Printf("Границы найдены за %v\n", time.Since(startEdges))

	saveToFile("results/output_edges_lib.png", edgeImg)
}

// Вспомогательная функция для сохранения результата
func saveToFile(filename string, img image.Image) {
	out, err := os.Create(filename)
	if err != nil {
		log.Printf("Ошибка создания файла %s: %v", filename, err)
		return
	}

	err = png.Encode(out, img)
	if err != nil {
		log.Printf("Ошибка сохранения %s: %v", filename, err)
	}

	err = out.Close()
	if err != nil {
		log.Printf("Ошибка закрытия файла %s: %v", filename, err)
	}
}
