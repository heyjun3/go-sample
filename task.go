package gosample

import "fmt"

type TaskInterface interface {
	Task() (ID string, Name string, Status string)
}

type Task struct {
	ID     string
	Name   string
	Status string
}

func (t *Task) Task() (ID string, Name string, Status string) {
	ID = t.ID
	Name = t.Name
	Status = t.Status
	return
}

type DoneTask struct {
	task Task
}
func (t *DoneTask) Task() (ID string, Name string, Status string) {
	ID, Name, Status = t.task.Task()
	return
}
func NewDoneTask(id, name, status string) (*DoneTask, error) {
	if status != "done" {
		return nil, fmt.Errorf("failed init task error")
	}
	return &DoneTask{
		task: Task{
			ID: id,
			Name: name,
			Status: status,
		},
	}, nil
}

type UnDoneTask struct {
	task Task
}

func (t *UnDoneTask) Task() (ID string, Name string, Status string) {
	ID, Name, Status = t.task.Task()
	return
}

func (t *UnDoneTask) Done() *DoneTask {
	return &DoneTask{
		task: Task{
			ID:   t.task.ID,
			Name: t.task.Name,
		},
	}
}

type TaskRepositoryInterface interface {
	Save(task TaskInterface) error
}

type TaskRepository struct{}

func (r *TaskRepository) Save(task TaskInterface) error {
	task.Task()
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
	err := repo.Save(&task)
	if err != nil {
		panic(err)
	}
}
