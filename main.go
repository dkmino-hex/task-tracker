package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

const dataFile = "task.json"

type TaskStatus string

const (
	todo       TaskStatus = "todo"
	inProgress TaskStatus = "in-progress"
	done       TaskStatus = "done"
)

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func loadTasks() ([]Task, error) {
	data, err := os.ReadFile(dataFile)
	if os.IsNotExist(err) {
		return []Task{}, nil
	}
	if err != nil {
		return nil, err
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

func addtask(title, description string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	newTask := Task{
		ID:          maxID + 1,
		Title:       title,
		Description: description,
		Status:      todo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tasks = append(tasks, newTask)
	newTask.CreatedAt = time.Now()
	return saveTasks(tasks)
}

func listTasks(filter TaskStatus) error {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return err
	}

	if len(tasks) == 0 {
		fmt.Println("Нет задач")
		return nil
	}

	shown := 0

	for _, task := range tasks {
		if filter != "" && task.Status != filter {
			continue
		}

		fmt.Printf(
			"%d. [%s] %s — %s\n",
			task.ID,
			task.Status,
			task.Title,
			task.Description,
		)

		shown++
	}

	if shown == 0 {
		fmt.Printf("Задач со статусом %q нет.\n", filter)
	}
	return nil
}

func deleteTask(id int) error {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return err
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Printf("Задача с ID %d удалена\n", id)
			return saveTasks(tasks)

		}
	}
	return fmt.Errorf("Задача с ID %d не найдена", id)
}

func updateTask(id int, title, description string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Title = title
			tasks[i].Description = description
			tasks[i].UpdatedAt = time.Now()
			return saveTasks(tasks)
		}
	}
	return fmt.Errorf("Задача с ID %d не найдена", id)
}

func markTaskDone(id int) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = done
			tasks[i].UpdatedAt = time.Now()

			if err := saveTasks(tasks); err != nil {
				fmt.Println("Ошибка сохранения:", err)
				return
			}

			fmt.Println("Задача отмечена как выполненная.")
			return
		}
	}
	fmt.Println("Задача с указанным ID не найдена.")
}

func markTaskInProgress(id int) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = inProgress
			tasks[i].UpdatedAt = time.Now()

			if err := saveTasks(tasks); err != nil {
				fmt.Println("Ошибка сохранения:", err)
				return
			}

			fmt.Println("Задача отмечена как в процессе.")
			return
		}
	}
	fmt.Println("Задача с указанным ID не найдена.")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование:")
		fmt.Println("task-tracker add \"Название\" \"Описание\"")
		fmt.Println("task-tracker list")
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
		if err := addtask(os.Args[2], os.Args[3]); err != nil {
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

	case "updateTask":
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

		number, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("ID должен быть числом.")
			return
		}

		markTaskDone(number)

	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("Укажите ID задачи.")
			return
		}

		number, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("ID должен быть числом.")
			return
		}

		markTaskInProgress(number)

	default:
		fmt.Println("Неизвестная команда:", command)
	}
}
