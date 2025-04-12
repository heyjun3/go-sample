package gosample

import (
	"fmt"
)

type Task struct {
	ID     string
	Name   string
	Status string
}

type UnDoneTask struct {
	task Task
}
func (t *UnDoneTask) Task() Task {
	return t.task
}
func (t *UnDoneTask) Done() *DoneTask {
	return &DoneTask{
		task: Task{
			ID:   t.task.ID,
			Name: t.task.Name,
		},
	}
}
func NewUnDoneTask(id, name, status string) (*UnDoneTask, error) {
	if status != "undone" {
		return nil, fmt.Errorf("invalid status: %s", status)
	}
	return &UnDoneTask{
		task: Task{
			ID:     id,
			Name:   name,
			Status: status,
		},
	}, nil
}

type DoneTask struct {
	task Task
}
func (t *DoneTask) Task() Task {
	return t.task
}
func (t *DoneTask) UnDone() *UnDoneTask {
	return &UnDoneTask{
		task: Task{
			ID:   t.task.ID,
			Name: t.task.Name,
		},
	}
}
func NewDoneTask(id, name, status string) (*DoneTask, error) {
	if status != "done" {
		return nil, fmt.Errorf("invalid status: %s", status)
	}
	return &DoneTask{
		task: Task{
			ID:     id,
			Name:   name,
			Status: status,
		},
	}, nil
}

type TaskRepositoryInterface interface {
	Save(task Task) error
}

type TaskRepository struct{}

func (r *TaskRepository) Save(task Task) error {
	return nil
}

func SaveTask() {
	repo := &TaskRepository{}
	task := UnDoneTask{
		task: Task{
			ID:   "1",
			Name: "task1",
		},
	}
	err := repo.Save(task.Task())
	if err != nil {
		panic(err)
	}
}
