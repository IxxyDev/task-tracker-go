package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const filename = "tasks.json"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	args := os.Args[1:]
	cmd := args[0]

	switch cmd {
	case "add":
		if len(args) < 2 {
			return errors.New("There is no description for task")
		}
		return addTask(args[1])
	case "update":
		if len(args) < 3 {
			return errors.New("You should enter ID of the task and its description")
		}
		return updateTask(args[1], args[2])
	case "mark-done":
		if len(args) < 2 {
			return errors.New("You should enter ID of the task you want make done")
		}
		return updateTaskStatus(args[1], "done")
	case "mark-in-progress":
		if len(args) < 2 {
			return errors.New("You should enter ID of the task you want make in-progress")
		}
		return updateTaskStatus(args[1], "in-progress")
	case "list":
		status := "all"
		if len(args) > 1 {
			status = args[1]
		}
		return listTasks(status)
	case "delete":
		if len(args) < 2 {
			return errors.New("You should enter ID of the task you want to delete")
		}
		return deleteTask(args[1])
	default:
		return fmt.Errorf("unknown command %s", cmd)
	}
}

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdateAt    time.Time `json:"updateAt"`
}

func readTasksFromFile() ([]Task, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func writeTasksToFile(tasks []Task) error {
	jsonData, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func addTask(description string) error {
	tasks, err := readTasksFromFile()
	if err != nil {
		return err
	}
	var maxID int

	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	newTask := Task{
		ID:          maxID + 1,
		Description: description,
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdateAt:    time.Now(),
	}

	tasks = append(tasks, newTask)
	fmt.Printf("Task has been added (ID: %d)\n", newTask.ID)
	return writeTasksToFile(tasks)
}

func listTasks(statusFilter string) error {
	tasks, err := readTasksFromFile()
	if err != nil {
		return err
	}

	fmt.Println("------------------------")
	found := false

	for _, task := range tasks {
		if statusFilter == "all" || task.Status == statusFilter {
			fmt.Printf("%d. [%s] %s\n", task.ID, task.Status, task.Description)
		}
	}

	if !found && statusFilter != "all" {
		fmt.Printf("Tasks with status '%s' are not found", statusFilter)
	}
	if len(tasks) == 0 {
		fmt.Println("Task list is empty")
	}
	fmt.Println("------------------------")

	return nil
}

func updateTask(idStr string, description string) error {
	tasks, err := readTasksFromFile()
	if err != nil {
		return err
	}
	taskIndex, err := findTaskIndexByID(tasks, idStr)
	if err != nil {
		return err
	}

	tasks[taskIndex].Description = description
	tasks[taskIndex].UpdateAt = time.Now()
	fmt.Printf("The task %d has been updated.\n", tasks[taskIndex].ID)
	return writeTasksToFile(tasks)
}

func updateTaskStatus(idStr string, status string) error {
	tasks, err := readTasksFromFile()
	if err != nil {
		return err
	}

	taskIndex, err := findTaskIndexByID(tasks, idStr)
	if err != nil {
		return err
	}

	tasks[taskIndex].UpdateAt = time.Now()
	tasks[taskIndex].Status = "done"
	fmt.Printf("Status of the task %d has been changed to %s.\n", tasks[taskIndex].ID, status)
	return writeTasksToFile(tasks)
}

func deleteTask(idStr string) error {
	tasks, err := readTasksFromFile()
	if err != nil {
		return err
	}
	ID, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("incorrect ID format: %s", idStr)
	}

	var newTasks []Task
	found := false
	for _, task := range tasks {
		if task.ID != ID {
			newTasks = append(newTasks, task)
		} else {
			found = true

		}
	}

	if !found {
		return fmt.Errorf("The task with ID %d is not found", ID)
	}

	fmt.Printf("The task %d has been removed.\n", ID)
	return writeTasksToFile(newTasks)
}

func findTaskIndexByID(tasks []Task, idStr string) (int, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return -1, fmt.Errorf("incorrect ID format: %s", idStr)
	}

	for i, task := range tasks {
		if task.ID == id {
			return i, nil
		}
	}

	return -1, fmt.Errorf("The task with ID %d is not found", id)
}

func printHelp() {
	fmt.Println("Usage: task-cli <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  add <description>         - Add new task")
	fmt.Println("  list [status]             - Show tasks (all, todo, in-progress, done)")
	fmt.Println("  update <ID> <description> - Update description of the task")
	fmt.Println("  delete <ID>               - Delete the task")
	fmt.Println("  mark-done <ID>            - Mark task as 'done'")
	fmt.Println("  mark-in-progress <ID>     - Mark task as 'in-progress'")
}
