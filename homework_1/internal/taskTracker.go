package internal

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

func StartTaskTracker(storage Storage) {
	layout := "02.01.2006 15:04"
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Err()
	fmt.Println("Добро пожаловать в таск-трекер!")
	for {
		scanner.Scan()
		if scanner.Text() != "добавить" &&
			scanner.Text() != "список" &&
			scanner.Text() != "найти" &&
			scanner.Text() != "статус" &&
			scanner.Text() != "редактировать" &&
			scanner.Text() != "удалить" &&
			scanner.Text() != "help" &&
			scanner.Text() != "exit" {
			fmt.Println("Неверная команда. Чтобы узнать доступные команды введите help")
		}
		if scanner.Text() == "help" {
			fmt.Println("Доступны следующие команды")
			fmt.Println("1)\"добавить\" - добавить задачу\n2)\"список\" - получить список задач")
			fmt.Println("3)\"найти\" - найти определённую задачу\n4)\"статус\" - изменить статус задачи")
			fmt.Println("5)\"редактировать\" - позволяет редактировать выбранную задачу")
			fmt.Println("6)\"удалить\" - позволяет удалить выбранную задачу")
			fmt.Println("7)\"exit\" - для выхода из программы")
		}
		if scanner.Text() == "добавить" {
			fmt.Print("Введите название задачи: ")
			scanner.Scan()
			name := scanner.Text()
			for {
				fmt.Print("Дедлайн - введите в формате ДД.ММ.ГГГГ ЧЧ:ММ : ")
				scanner.Scan()
				deadline, err := time.ParseInLocation(layout, scanner.Text(), time.Local)
				if err != nil {
					fmt.Println("Ошибка. Ставьте проблей между датой и временем!")
					continue
				}
				storage.AddTask(name, deadline)
				fmt.Println("Задача успешно создана!")
				break
			}
		}
		if scanner.Text() == "список" {
			storage.GetTaskList()
			continue
		}
		if scanner.Text() == "найти" {
			for {
				fmt.Print("Введите id задачи: ")
				scanner.Scan()
				id, err2 := strconv.Atoi(scanner.Text())
				if err2 != nil {
					fmt.Println("id задачи указан неверно")
					continue
				}
				storage.GetDefiniteTask(id)
				break
			}
		}
		if scanner.Text() == "статус" {
			for {
				fmt.Println("Введите id задачи, статус которой хотите поменять")
				fmt.Print(": ")
				scanner.Scan()
				id, err := strconv.Atoi(scanner.Text())
				if err != nil {
					fmt.Println("Ошибка при вводе id")
					continue
				}
				for {
					fmt.Println("Введите новый статус задачи: planned, in_progress, canceled, done")
					var newStatus TaskStatus
					scanner.Scan()
					if scanner.Text() == "planned" {
						newStatus = StatusPlanned
					}
					if scanner.Text() == "in_progress" {
						newStatus = StatusInProgress
					}
					if scanner.Text() == "canceled" {
						newStatus = StatusCanceled
					}
					if scanner.Text() == "done" {
						newStatus = StatusDone
					}
					if scanner.Text() != string(StatusPlanned) &&
						scanner.Text() != string(StatusInProgress) &&
						scanner.Text() != string(StatusCanceled) &&
						scanner.Text() != string(StatusDone) {
						fmt.Println("Ошибка в указании статуса. Укажите один из предложенных")
						continue
					}
					storage.ChangeTaskStatus(id, newStatus)
					fmt.Println("Статус успешно изменён!")
					break
				}
				break
			}
		}
		if scanner.Text() == "редактировать" {
			for {
				fmt.Print("Введите id задачи: ")
				scanner.Scan()
				id, err2 := strconv.Atoi(scanner.Text())
				if err2 != nil {
					fmt.Println("id задачи указан неверно")
					continue
				}
				fmt.Println("Введите новое название задачи: ")
				scanner.Scan()
				newName := scanner.Text()
				fmt.Println("Введите новый дедлайн в формате ДД.ММ.ГГГГ ЧЧ:ММ : ")
				scanner.Scan()
				deadline, err := time.ParseInLocation(layout, scanner.Text(), time.Local)
				if err != nil {
					fmt.Println("Ошибка. Ставьте проблей между датой и временем!")
					continue
				}
				storage.ChangeTaskInfo(id, newName, deadline)
				fmt.Println("Задача успешно отредактирована!")
				break
			}
		}
		if scanner.Text() == "удалить" {
			for {
				fmt.Print("Введите id задачи для удаления: ")
				scanner.Scan()
				id, err2 := strconv.Atoi(scanner.Text())
				if err2 != nil {
					fmt.Println("id задачи указан неверно")
					continue
				}
				storage.DeleteTask(id)
				fmt.Println("Задача успешно удалена!")
				break
			}
		}
		if scanner.Text() == "exit" {
			os.Exit(0)
		}
	}
}
