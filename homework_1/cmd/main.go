package main

import (
	"log"

	"github.com/Snowdy54/go_homework/tree/main/homework_1/internal"
)

func main() {
	const fileName = "../data/data.json"
	storage, err := internal.CreateStorage(fileName)
	if err != nil {
		log.Fatalf("Ошибка при создании базы %v", err)
	}

	internal.StartTaskTracker(*storage)
}
