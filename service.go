package main

import (
	"fmt"
	"time"
)

func addTask(title, description string) error {
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

	now := time.Now()

	newTask := Task{
		ID:          maxID + 1,
		Title:       title,
		Description: description,
		Status:      todo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tasks = append(tasks, newTask)
	if err := saveTasks(tasks); err != nil {
		return err
	}
	fmt.Printf("Задача добавлена с ID %d\n", newTask.ID)
	return nil
}

func listTasks(filter TaskStatus) error {
	tasks, err := loadTasks()
	if err != nil {
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

func findTaskIndex(tasks []Task, id int) (int, error) {
	for i, task := range tasks {
		if task.ID == id {
			return i, nil
		}
	}

	return -1, fmt.Errorf("задача с ID %d не найдена", id)
}

func deleteTask(id int) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	index, err := findTaskIndex(tasks, id)
	if err != nil {
		return err
	}

	tasks = append(tasks[:index], tasks[index+1:]...)
	if err := saveTasks(tasks); err != nil {
		return err
	}

	fmt.Printf("Задача с ID %d удалена\n", id)
	return nil
}

func updateTask(id int, title, description string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	index, err := findTaskIndex(tasks, id)
	if err != nil {
		return err
	}

	tasks[index].Title = title
	tasks[index].Description = description
	tasks[index].UpdatedAt = time.Now()
	if err := saveTasks(tasks); err != nil {
		return err
	}

	fmt.Printf("Задача обновлена с ID %d\n", id)
	return nil
}

func markTaskDone(id int) error {
	return changeTaskStatus(id, done)
}

func markTaskInProgress(id int) error {
	return changeTaskStatus(id, inProgress)
}

func changeTaskStatus(id int, status TaskStatus) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	index, err := findTaskIndex(tasks, id)
	if err != nil {
		return err
	}

	tasks[index].Status = status
	tasks[index].UpdatedAt = time.Now()
	if err := saveTasks(tasks); err != nil {
		return err
	}

	switch status {
	case done:
		fmt.Println("Задача отмечена как выполненная.")
	case inProgress:
		fmt.Println("Задача отмечена как в процессе.")
	}

	return nil
}
