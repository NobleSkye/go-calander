// File: main.go

package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Task struct {
	Title       string
	Description string
	Priority    string
}

var tasks []Task
var mu sync.Mutex

// Load tasks from a file
func loadTasks(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return // File doesn't exist or can't be opened
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var title, desc, priority string
		fmt.Sscanf(line, "%s|%s|%s", &title, &desc, &priority)
		tasks = append(tasks, Task{Title: title, Description: desc, Priority: priority})
	}
}

// Save tasks to a file
func saveTasks(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	for _, task := range tasks {
		fmt.Fprintf(file, "%s|%s|%s\n", task.Title, task.Description, task.Priority)
	}
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Task Calendar")

	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Task Title")

	descriptionEntry := widget.NewMultiLineEntry()
	descriptionEntry.SetPlaceHolder("Task Description")

	priorityEntry := widget.NewEntry()
	priorityEntry.SetPlaceHolder("Task Priority (high, medium, low)")

	taskList := widget.NewList(
		func() int {
			mu.Lock()
			defer mu.Unlock()
			return len(tasks)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			mu.Lock()
			defer mu.Unlock()
			obj.(*widget.Label).SetText(fmt.Sprintf("%s [%s]", tasks[i].Title, tasks[i].Priority))
		})

	addButton := widget.NewButton("Add Task", func() {
		mu.Lock()
		defer mu.Unlock()
		newTask := Task{
			Title:       titleEntry.Text,
			Description: descriptionEntry.Text,
			Priority:    priorityEntry.Text,
		}
		tasks = append(tasks, newTask)
		taskList.Refresh()
		saveTasks("tasks.txt") // Save tasks to file
		titleEntry.SetText("")
		descriptionEntry.SetText("")
		priorityEntry.SetText("")
	})

	// Load existing tasks
	loadTasks("tasks.txt")

	myWindow.SetContent(container.NewVBox(
		widget.NewLabel("Task Title:"),
		titleEntry,
		widget.NewLabel("Task Description:"),
		descriptionEntry,
		widget.NewLabel("Task Priority:"),
		priorityEntry,
		addButton,
		taskList,
	))

	myWindow.Resize(fyne.NewSize(400, 600))
	myWindow.ShowAndRun()
}
