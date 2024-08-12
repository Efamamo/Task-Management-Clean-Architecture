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
	task1Id, _ := primitive.ObjectIDFromHex("66b90802a6d354d828cd8f1a")
	task2Id, _ := primitive.ObjectIDFromHex("66b918ff59c53bc04a1d2594")
	task3Id, _ := primitive.ObjectIDFromHex("66b9ab13c65b0b36f13c5a5c")
	var tasks = []domain.Task{
		{ID: task1Id, Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
		{ID: task2Id, Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
		{ID: task3Id, Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
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
	taskId, _ := primitive.ObjectIDFromHex("66b90802a6d354d828cd8f1a")
	task := domain.Task{ID: taskId, Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"}

	suite.mockRepo.On("FindOne", taskId).Return(&task, nil)

	t, err := suite.taskService.GetTaskByID(string(taskId.Hex()))

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
	taskId, _ := primitive.ObjectIDFromHex("66b90802a6d354d828cd8f1a")
	task := domain.Task{ID: taskId, Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"}

	suite.mockRepo.On("UpdateOne", taskId, task).Return(nil)

	err := suite.taskService.UpdateItem(taskId.Hex(), task)
	suite.NoError(err)

}

func (suite *TaskSuite) TestUpdateTask_InputError() {
	id, _ := primitive.ObjectIDFromHex("66b90802a6d354d828cd8f1a")
	task := domain.Task{ID: id, Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"}

	suite.mockRepo.On("UpdateOne", id, task).Return(errors.New("status error"))

	err := suite.taskService.UpdateItem(id.Hex(), task)
	suite.EqualError(err, "status error")

}

func (suite *TaskSuite) TestDeleteTask() {
	id, _ := primitive.ObjectIDFromHex("66b90802a6d354d828cd8f1a")
	task := domain.Task{ID: id, Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"}

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
	taskId, _ := primitive.ObjectIDFromHex("66b90802a6d354d828cd8f1f")
	task := domain.Task{ID: taskId, Title: "Task 4", Description: "Fourth task", DueDate: time.Now(), Status: "Pending"}

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
	taskId, _ := primitive.ObjectIDFromHex("66b90802a6d354d828cd8f1b")
	task := domain.Task{ID: taskId, Description: "First task", DueDate: time.Now(), Status: "Pending"}

	suite.mockRepo.On("Save", task).Return(&task, errors.New("Title is Required"))

	_, err := suite.taskService.AddTask(task)
	suite.EqualError(err, "Title is Required")

}
