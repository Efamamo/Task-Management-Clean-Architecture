package usecases

import (
	"errors"
	"testing"
	"time"

	domain "github.com/Task-Management-go/Domain"
	mocks "github.com/Task-Management-go/mocks"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskSuite struct {
	suite.Suite
	taskService TaskService
	mockRepo    *mocks.MockTaskRepository
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(TaskSuite))
}

func (suite *TaskSuite) SetupTest() {
	suite.mockRepo = new(mocks.MockTaskRepository)
	suite.taskService = TaskService{TaskRepo: suite.mockRepo}
}

func (suite *TaskSuite) TestGetTasks() {
	id := primitive.NewObjectID()
	tasks := []domain.Task{
		{ID: id,
			Title:       "Class work",
			Description: "Do class WOrk",
			DueDate:     time.Now(),
			Status:      "In Progress",
		},
		{ID: id,
			Title:       "Home work",
			Description: "Do Home WOrk",
			DueDate:     time.Now(),
			Status:      "Completed",
		},
		{ID: id,
			Title:       "Exam",
			Description: "Do Exam WOrk",
			DueDate:     time.Now(),
			Status:      "Pending",
		},
	}
	suite.mockRepo.On("FindAll").Return(&tasks, nil)

	t, err := suite.taskService.GetTasks()
	ts := *t

	for i := 0; i < len(ts); i++ {
		suite.Equal(tasks[i].Title, ts[i].Title)
		suite.Equal(tasks[i].Description, ts[i].Description)
		suite.Equal(tasks[i].DueDate, ts[i].DueDate)
		suite.Equal(tasks[i].ID, ts[i].ID)
		suite.Equal(tasks[i].Status, ts[i].Status)
	}
	suite.NoError(err)

}

func (suite *TaskSuite) TestGetSingleTask() {
	id := primitive.NewObjectID()
	task := domain.Task{
		ID:          id,
		Title:       "Class work",
		Description: "Do class WOrk",
		DueDate:     time.Now(),
		Status:      "In Progress",
	}

	suite.mockRepo.On("FindOne", id).Return(&task, nil)

	t, err := suite.taskService.GetTaskByID(id.Hex())

	suite.Equal(task.Title, t.Title)
	suite.Equal(task.Description, t.Description)
	suite.Equal(task.DueDate, t.DueDate)
	suite.Equal(task.ID, t.ID)
	suite.Equal(task.Status, t.Status)
	suite.NoError(err)

}

func (suite *TaskSuite) TestErrorGetSingleTask() {
	id := primitive.NewObjectID()

	suite.mockRepo.On("FindOne", id).Return(&domain.Task{}, errors.New("Task Not Found"))

	_, err := suite.taskService.GetTaskByID(id.Hex())
	suite.EqualError(err, "Task Not Found")

}

func (suite *TaskSuite) TestUpdateTask() {
	id := primitive.NewObjectID()
	task := domain.Task{
		ID:          id,
		Title:       "Class work",
		Description: "Do class WOrk",
		DueDate:     time.Now(),
		Status:      "Completed",
	}

	suite.mockRepo.On("UpdateOne", id, task).Return(nil)

	err := suite.taskService.UpdateItem(id.Hex(), task)
	suite.NoError(err)

}

func (suite *TaskSuite) TestUpdateTask_InputError() {
	id := primitive.NewObjectID()
	task := domain.Task{
		ID:          id,
		Title:       "Class work",
		Description: "Do class WOrk",
		DueDate:     time.Now(),
	}

	suite.mockRepo.On("UpdateOne", id, task).Return(errors.New("status error"))

	err := suite.taskService.UpdateItem(id.Hex(), task)
	suite.EqualError(err, "status error")

}

func (suite *TaskSuite) TestDeleteTask() {
	id := primitive.NewObjectID()
	task := domain.Task{
		ID:          id,
		Title:       "Class work",
		Description: "Do class WOrk",
		DueDate:     time.Now(),
		Status:      "Completed",
	}

	suite.mockRepo.On("DeleteOne", id).Return(&task, nil)

	t, err := suite.taskService.DeleteTask(id.Hex())

	suite.Equal(task.Title, t.Title)
	suite.Equal(task.Description, t.Description)
	suite.Equal(task.DueDate, t.DueDate)
	suite.Equal(task.ID, t.ID)
	suite.Equal(task.Status, t.Status)
	suite.NoError(err)

	suite.NoError(err)

}

func (suite *TaskSuite) TestDeleteTask_NotFound() {
	id := primitive.NewObjectID()

	suite.mockRepo.On("DeleteOne", id).Return(&domain.Task{}, errors.New("Task Not Found"))

	_, err := suite.taskService.DeleteTask(id.Hex())

	suite.EqualError(err, "Task Not Found")

}

func (suite *TaskSuite) TestSaveTask() {
	id := primitive.NewObjectID()
	task := domain.Task{
		ID:          id,
		Title:       "Class work",
		Description: "Do class WOrk",
		DueDate:     time.Now(),
		Status:      "Completed",
	}

	suite.mockRepo.On("Save", task).Return(&task, nil)

	t, err := suite.taskService.AddTask(task)

	suite.Equal(task.Title, t.Title)
	suite.Equal(task.Description, t.Description)
	suite.Equal(task.DueDate, t.DueDate)
	suite.Equal(task.ID, t.ID)
	suite.Equal(task.Status, t.Status)
	suite.NoError(err)

	suite.NoError(err)

}

func (suite *TaskSuite) TestSaveTask_InvalidData() {
	id := primitive.NewObjectID()
	task := domain.Task{
		ID:          id,
		Description: "Do class WOrk",
		DueDate:     time.Now(),
		Status:      "Completed",
	}

	suite.mockRepo.On("Save", task).Return(&task, errors.New("Title is Required"))

	_, err := suite.taskService.AddTask(task)

	suite.EqualError(err, "Title is Required")

}
