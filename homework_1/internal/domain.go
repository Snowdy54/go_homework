package internal

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

var filePath string

type ID int
type TaskStatus string

type Identifier struct {
	Identifier ID
}

type Storage struct {
	IdNext ID     `json:"idnext"`
	Tasks  []Task `json:"tasks"`
}

func (s *Storage) GetNextIdentifier() ID {
	s.IdNext++
	return s.IdNext
}

func CreateStorage(pathOfFile string) (*Storage, error) {
	storage := Storage{IdNext: 0, Tasks: []Task{}}

	filePath = pathOfFile

	fmt.Println(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("Ошибка при создании хранилища данных")
		}
	}

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	decoder.Decode(&storage)

	file.Close()

	changeFile, err2 := os.Create(filePath)
	if err2 != nil {
		return nil, errors.New("Ошибка при создании хранилища данных")
	}

	writer := bufio.NewWriter(changeFile)
	encoder := json.NewEncoder(writer)
	encoder.Encode(storage)
	writer.Flush()

	return &storage, nil
}

func (id *Identifier) GetIdentifier() ID {
	id.Identifier++
	return id.Identifier
}

const (
	StatusPlanned    TaskStatus = "planned"
	StatusInProgress TaskStatus = "in_progress"
	StatusCanceled   TaskStatus = "canceled"
	StatusDone       TaskStatus = "done"
)

type Task struct {
	Identifier ID         `json:"identifier"`
	TaskName   string     `json:"taskName"`
	DeadLine   time.Time  `json:"deadLine"`
	Status     TaskStatus `json:"status"`
}
