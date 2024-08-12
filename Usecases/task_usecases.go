package usecases

import (
	"errors"
	"strings"

	domain "github.com/Task-Management-go/Domain"
	"github.com/Task-Management-go/Domain/err"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskService struct {
	TaskRepo ITaskRepository
}

func (ts *TaskService) GetTasks() (*[]domain.Task, error) {

	tasks, err := ts.TaskRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return tasks, nil

}

func (ts *TaskService) GetTaskByID(id string) (*domain.Task, error) {
	objectId, e := primitive.ObjectIDFromHex(id)
	if e != nil {
		return nil, err.NewValidation("invalid ID format")
	}
	task, err := ts.TaskRepo.FindOne(objectId)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (ts *TaskService) UpdateItem(ID string, updatedTask domain.Task) error {

	if strings.ToLower(updatedTask.Status) != "in progress" && strings.ToLower(updatedTask.Status) != "completed" && strings.ToLower(updatedTask.Status) != "pending" {
		return errors.New("status error")
	}

	objectId, e := primitive.ObjectIDFromHex(ID)
	if e != nil {
		return err.NewValidation("invalid ID format")
	}

	e = ts.TaskRepo.UpdateOne(objectId, updatedTask)

	if e != nil {
		return e
	}

	return nil

}

func (ts *TaskService) DeleteTask(ID string) (*domain.Task, error) {
	objectId, e := primitive.ObjectIDFromHex(ID)
	if e != nil {
		return nil, err.NewValidation("invalid ID format")
	}

	task, err := ts.TaskRepo.DeleteOne(objectId)

	if err != nil {
		return nil, err
	}
	return task, nil

}

func (ts *TaskService) AddTask(task domain.Task) (*domain.Task, error) {
	if strings.ToLower(task.Status) != "in progress" && strings.ToLower(task.Status) != "completed" && strings.ToLower(task.Status) != "pending" {
		return nil, errors.New("status error")
	}

	t, err := ts.TaskRepo.Save(task)

	if err != nil {
		return nil, err
	}

	return t, nil
}
