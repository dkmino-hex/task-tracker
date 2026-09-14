package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование:")
		fmt.Println("task-tracker add \"Название\" \"Описание\"")
		fmt.Println("task-tracker list [todo|in-progress|done]")
		fmt.Println("task-tracker update ID \"Новое название\" \"Новое описание\"")
		fmt.Println("task-tracker delete ID")
		fmt.Println("task-tracker mark-done ID")
		fmt.Println("task-tracker mark-in-progress ID")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 4 {
			fmt.Println("Укажите название и описание задачи.")
			return
		}
		if err := addTask(os.Args[2], os.Args[3]); err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

	case "list":
		var filter TaskStatus
		if len(os.Args) > 3 {
			fmt.Println("Использование: task-tracker list [todo|in-progress|done]")
			return
		}

		if len(os.Args) == 3 {
			filter = TaskStatus(os.Args[2])

			if filter != todo && filter != inProgress && filter != done {
				fmt.Println("Неизвестный статус. Используйте: todo, in-progress или done.")
				return
			}
		}

		if err := listTasks(filter); err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Укажите ID задачи.")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("ID должен быть числом.")
			return
		}
		if err := deleteTask(id); err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

	case "update":
		if len(os.Args) < 5 {
			fmt.Println("Укажите ID задачи, новое название и новое описание.")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("ID должен быть числом.")
			return
		}
		if err := updateTask(id, os.Args[3], os.Args[4]); err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("Укажите ID задачи.")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("ID должен быть числом.")
			return
		}

		if err := markTaskDone(id); err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("Укажите ID задачи.")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("ID должен быть числом.")
			return
		}

		if err := markTaskInProgress(id); err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

	default:
		fmt.Println("Неизвестная команда:", command)
	}
}
