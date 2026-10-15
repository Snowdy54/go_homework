package internal

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func isOverdue(deadLine time.Time) bool {
	return time.Now().Before(time.Time(deadLine))
}

func (s *Storage) AddTask(taskName string, deadLine time.Time) (*Storage, error) {
	if strings.TrimSpace(taskName) == "" {
		return s, errors.New("У задачи должен быть заголовок")
	}

	if !isOverdue(deadLine) {
		return s, errors.New("Дедлайн не может быть в прошлом")
	}

	identifierValue := s.GetNextIdentifier()
	task := Task{Identifier: identifierValue, TaskName: taskName,
		DeadLine: deadLine, Status: StatusPlanned}

	s, err := s.SaveTask(task)
	if err != nil {
		return s, errors.New("Ошибка при создании задачи")
	}
	return s, nil
}

func (s *Storage) SaveTask(task Task) (*Storage, error) {
	file, err := os.Open(filePath)
	if err != nil {
		panic("Ошибка при открытии базы данных для записи!")
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)

	err1 := decoder.Decode(s)
	if err1 != nil {
		return nil, errors.New("Ошибка при чтении файла")
	}

	file.Close()

	changeFile, err2 := os.Create(filePath)
	if err2 != nil {
		return s, errors.New("Ошибка при чтении файла")
	}
	defer changeFile.Close()

	encoder := json.NewEncoder(changeFile)

	s.Tasks = append(s.Tasks, task)
	s.IdNext = s.GetNextIdentifier()
	err3 := encoder.Encode(s)
	if err3 != nil {
		log.Fatalf("Ошибка при сохранении файла %d", err3)
	}

	return s, nil
}

func (s *Storage) GetTaskList() error {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)

	var storage Storage

	err1 := decoder.Decode(&storage)
	if err1 != nil {
		return errors.New("Файл пуст!")
	}

	fmt.Println("id    название       дедлайн                статус")

	for _, task := range storage.Tasks {
		fmt.Println(task)
	}
	return nil
}

func (s *Storage) GetDefiniteTask(id int) error {

	file, err := os.Open(filePath)
	if err != nil {
		panic("Ошибка при открытии базы данных для чтения!")
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)

	err1 := decoder.Decode(&s)
	if err1 != nil {
		return errors.New("Ошибка при чтении файла")
	}

	taskID := ID(id)

	for _, task := range s.Tasks {
		if task.Identifier == taskID {
			fmt.Println("id    название       дедлайн                статус")
			fmt.Println(task.Identifier, task.TaskName, task.DeadLine, task.Status)
		}
	}

	return nil
}

func (s *Storage) ChangeTaskStatus(id int, newStatus TaskStatus) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.New("Ошибка при чтении базы данных")
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)

	err1 := decoder.Decode(&s)
	if err1 != nil {
		return errors.New("Ошибка при записи файла")
	}

	file.Close()

	changeFile, err2 := os.Create(filePath)
	if err2 != nil {
		log.Fatalf("Ошибка при открытии файла для изменения статуса %v", err2)
	}

	writer := bufio.NewWriter(changeFile)
	encoder := json.NewEncoder(writer)

	taskId := ID(id)

	for i, task := range s.Tasks {
		if task.Identifier == taskId {
			if (task.Status == StatusCanceled) || (task.Status == StatusDone) {
				return errors.New("Задачу в статусе done или canceled менять нельзя — её можно только удалить")
			}
			if task.Status == newStatus {
				return errors.New("Нельзя перевести задачу в статус, который у неё уже стоит")
			}
			s.Tasks[i].Status = newStatus
			encoder.Encode(s)
			writer.Flush()
		}
	}
	return nil
}

func (s *Storage) ChangeTaskInfo(id int, newName string, newDeadLine time.Time) error {

	file, err := os.Open(filePath)
	if err != nil {
		return errors.New("Ошибка при открытии файла")
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)

	err1 := decoder.Decode(&s)
	if err1 != nil {
		return errors.New("Ошибка при чтении файла")
	}

	file.Close()

	changeFile, err2 := os.Create(filePath)
	if err2 != nil {
		return errors.New("Ошибка при открытии файла")
	}

	writer := bufio.NewWriter(changeFile)
	encoder := json.NewEncoder(writer)

	taskId := ID(id)

	for i, task := range s.Tasks {
		if task.Identifier == taskId {
			if !isOverdue(newDeadLine) {
				return errors.New("Ошибка: дедлайн не может стоять в прошлом")
			}
			if strings.TrimSpace(newName) == "" {
				return errors.New("Имя задачи не может быть пустой")
			}
			s.Tasks[i].TaskName = newName
			s.Tasks[i].DeadLine = newDeadLine
			encoder.Encode(s)
			writer.Flush()
		}
	}
	return nil
}

func (s *Storage) DeleteTask(id int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.New("Ошибка при чтении файла")
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)

	err1 := decoder.Decode(&s)
	if err1 != nil {
		return errors.New("Ошибка при чтении файла")
	}

	if len(s.Tasks) == 0 {
		return fmt.Errorf("Ошибка!")
	}

	file.Close()

	changeFile, err2 := os.Create(filePath)
	if err2 != nil {
		return errors.New("Ошибка при открытии файла")
	}

	writer := bufio.NewWriter(changeFile)
	encoder := json.NewEncoder(writer)

	taskId := ID(id)

	for i, task := range s.Tasks {
		if task.Identifier == taskId {
			s.Tasks = append(s.Tasks[:i], s.Tasks[i+1:]...)
			encoder.Encode(s)
			writer.Flush()
		}
	}
	return nil
}
