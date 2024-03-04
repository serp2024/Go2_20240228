package taskstore

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

// TaskStore is a simple in-memory database of tasks;
type TaskStore struct {
	sync.Mutex

	tasks  map[int]Task
	nextId int
}

// TaskStore constructor
func New() *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[int]Task)
	ts.nextId = 1
	return ts
}

// CreateTask create a new task in the store
func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) int {
	ts.Lock()
	defer ts.Unlock()

	task := Task{
		Id:   ts.nextId,
		Text: text,
		Due:  due}
	task.Tags = make([]string, len(tags))
	copy(task.Tags, tags)
	// Сохранили task в TaskStore
	ts.tasks[ts.nextId] = task
	ts.nextId++
	return task.Id
}

// GetTask retrieves the task from taskstore by given id
func (ts *TaskStore) GetTask(id int) (Task, error) {
	ts.Lock()
	defer ts.Unlock()

	t, ok := ts.tasks[id]
	if ok {
		return t, nil
	} else {
		return Task{}, fmt.Errorf("task with id=%d not found", id)
	}

}

// GetAllTask retrieves all task from taskstore, in arbitrary order
func (ts *TaskStore) GetAllTasks() []Task {
	ts.Lock()
	defer ts.Unlock()

	allTasks := make([]Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		allTasks = append(allTasks, task)
	}

	return allTasks

}

// DeleteAllTasks deletes all tasks in the taskstore
func (ts *TaskStore) DeleteAllTasks() error {
	ts.Lock()
	defer ts.Unlock()

	ts.tasks = make(map[int]Task)
	return nil
}

// DeleteTask deletes the task from taskstore by given id. If no such id exists, return Error
func (ts *TaskStore) DeleteTask(id int) error {
	ts.Lock()
	defer ts.Unlock()

	if _, ok := ts.tasks[id]; !ok {
		return fmt.Errorf("task with id=%d not found", id)
	}
	delete(ts.tasks, id)
	return nil

}

func (ts *TaskStore) GetTaskByTag(zTag string) []Task {
	ts.Lock()
	defer ts.Unlock()

	//	allTasks := make([]Task, 0, len(ts.tasks))
	var allTasks []Task
	for _, task := range ts.tasks {
		for _, tag := range task.Tags {
			if tag == zTag {
				allTasks = append(allTasks, task)
				break
			}
		}
	}

	return allTasks
}

func (ts *TaskStore) GetTaskByYMD(year int, month int, day int) []Task {
	ts.Lock()
	defer ts.Unlock()

	zMonth := time.Month(month)
	var allTasks []Task
	for _, task := range ts.tasks {
		y, m, d := task.Due.Date()
		if y == year && m == zMonth && d == day {
			allTasks = append(allTasks, task)
		}
	}
	return allTasks
}

func (ts *TaskStore) GetTaskByYM(year int, month int) []Task {
	ts.Lock()
	defer ts.Unlock()

	zMonth := time.Month(month)
	var allTasks []Task
	for _, task := range ts.tasks {
		y, m, _ := task.Due.Date()
		if y == year && m == zMonth {
			allTasks = append(allTasks, task)
		}
	}
	return allTasks
}

func (ts *TaskStore) GetTaskByY(year int) []Task {
	ts.Lock()
	defer ts.Unlock()

	var allTasks []Task
	for _, task := range ts.tasks {
		y, _, _ := task.Due.Date()
		if y == year {
			allTasks = append(allTasks, task)
		}
	}
	return allTasks
}
