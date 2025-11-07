package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Task represents a single task with its properties
// JSON tags are used for serialization/deserialization.
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAT"`
}

const (
	tasksFile   = "tasks.json" // The name of the saved task file.
	statusTodo  = "todo"
	statusDone  = "done"
	statusDoing = "doing"
)

func main() {
	// Check for a command argument
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// The first argument is the command (e.g., "add", "list")
	command := os.Args[1]

	var err error

	switch command {
	case "add":
		// Usage: task add "Description"
		if len(os.Args) < 3 {
			fmt.Println("Usage: task add <description>")
			printUsage()
			os.Exit(1)
		}
		err = addTask(os.Args[2])

	case "update":
		// Usage: task update ID "New Description"
		if len(os.Args) < 4 {
			fmt.Println("Usage: task update <id> <new description>")
			os.Exit(1)
		}
		id, parseErr := strconv.Atoi(os.Args[2])
		if parseErr != nil {
			fmt.Printf("Error: Invalid task ID '%s'.\n", os.Args[2])
			os.Exit(1)
		}
		err = updateTask(id, os.Args[3])

	case "delete":
		// Usage: task delete ID
		if len(os.Args) < 3 {
			fmt.Println("Usage: task delete <id>")
			os.Exit(1)
		}
		id, parseErr := strconv.Atoi(os.Args[2])
		if parseErr != nil {
			fmt.Printf("Error: Invalid task ID '%s'.\n", os.Args[2])
			os.Exit(1)
		}
		err = deleteTask(id)

	case "mark":
		// Usage: task mark <status> <id>
		if len(os.Args) < 4 {
			fmt.Println("Usage: task mark <status> <id>")
			os.Exit(1)
		}

		status := strings.ToLower(os.Args[2])
		if status != statusDone && status != statusTodo && status != statusDoing {
			fmt.Printf("Invalid mark status '%s'. Use 'todo', 'doing', or 'done'.\n", status)
			os.Exit(1)
		}

		id, parseErr := strconv.Atoi(os.Args[3])
		if parseErr != nil {
			fmt.Printf("Error: Invalid task ID '%s'.\n", os.Args[2])
			os.Exit(1)
		}

		err = updateTaskStatus(id, status)
		if err == nil {
			fmt.Printf("Task ID %d marked as %s.\n", id, status)
		}

	case "list":
		// Usage: task list <status>
		filter := ""
		if len(os.Args) == 3 {
			filter = os.Args[2]
			// Basic validation for list filters
			if filter != statusDone && filter != statusTodo && filter != statusDoing {
				fmt.Printf("Invalid list status filter '%s'. Use 'done', 'todo', or 'doing'.\n", filter)
				os.Exit(1)
			}
		}
		err = listTasks(filter)

	default:
		fmt.Printf("Error: Unknown command '%s'\n", command)
		printUsage()
		os.Exit(1)
	}

	// Graceful error handling for task operations
	if err != nil {
		fmt.Printf("Operation Failed: %v\n", err)
		os.Exit(1)
	}
}

// printUsage displays the application's command usage instructions.
func printUsage() {
	fmt.Println("\nUsage: task <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  add \"<description>\"                    - Add a new task")
	fmt.Println("  update <ID> \"<new description>\"        - Update a task's description")
	fmt.Println("  delete <ID>                            - Delete a task")
	fmt.Println("  mark <status> <ID>                     - Mark a task with a status (todo, doing, done)")
	fmt.Println("  list <status>                          - List all tasks or filter by status (todo, doing, done)")
	fmt.Println()
}

// loadTasks reads tasks from the saved JSON file.
func loadTasks() ([]Task, error) {
	if _, err := os.Stat(tasksFile); os.IsNotExist(err) {
		return []Task{}, nil
	}

	data, err := os.ReadFile(tasksFile)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// If file is empty or only contains whitespace, return an empty slice
	if len(data) == 0 || len(data) > 0 && data[0] == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("error unmarshlling JSON: %w", err)
	}

	return tasks, nil
}

// getNextID creates a new ID.
func getNextID(tasks []Task) int {
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

// addTask adds a new task with "todo" status.
func addTask(description string) error {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	now := time.Now()

	newTask := Task{
		ID:          getNextID(tasks),
		Description: description,
		Status:      statusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tasks = append(tasks, newTask)
	err = saveTasks(tasks)
	if err != nil {
		return err
	}

	fmt.Printf("Task added successfully (ID: %d)\n", newTask.ID)
	return nil
}

// saveTasks writes the tasks slice to the JSON file.
func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	err = os.WriteFile(tasksFile, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}

// updateTask updates the description of a task by ID.
func updateTask(id int, description string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Description = description
			tasks[i].UpdatedAt = time.Now()
			return saveTasks(tasks)
		}
	}

	return fmt.Errorf("task with ID %d not found", id)
}

// deleteTask deletes a task by ID.
func deleteTask(id int) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == id {
			// Remove the task by slicing (efficient way to delete from a slice)
			tasks = append(tasks[:i], tasks[i+1:]...)
			if err := saveTasks(tasks); err != nil {
				return err
			}
			fmt.Printf("Task ID %d deleted successfully\n", id)
			return nil
		}
	}

	return fmt.Errorf("Task with ID %d not found", id)
}

// updateTaskStatus changes the status of a task by ID.
func updateTaskStatus(id int, newStatus string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = newStatus
			tasks[i].UpdatedAt = time.Now()
			return saveTasks(tasks)
		}
	}

	return fmt.Errorf("task with ID %d not found", id)
}

// listTasks prints tasks based on the filter.
func listTasks(filter string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	var filteredTasks []Task
	for _, task := range tasks {
		if filter == "" || task.Status == filter {
			filteredTasks = append(filteredTasks, task)
		}
	}

	if len(filteredTasks) == 0 {
		statusMsg := "all"
		if filter != "" {
			statusMsg = filter
		}
		fmt.Printf("No tasks found with status: %s\n", statusMsg)
		return nil
	}

	fmt.Println("--- Task List ---")
	for _, task := range filteredTasks {
		// Use a simple formatting for date/time
		createdAt := task.CreatedAt.Format("2006-01-02 15:04:05")
		updatedAt := task.UpdatedAt.Format("2006-01-02 15:04:05")

		fmt.Printf("[ID: %d] [%s] %s\n", task.ID, task.Status, task.Description)
		fmt.Printf("  Created: %s | Updated: %s\n", createdAt, updatedAt)
	}
	fmt.Println("-----------------")

	return nil
}
