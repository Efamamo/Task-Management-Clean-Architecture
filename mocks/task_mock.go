package mocks

import (
	"testing"

	domain "github.com/Task-Management-go/Domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockTaskRepository struct {
	mock.Mock
	t *testing.T
}

func NewMockTaskRepository(t *testing.T) *MockTaskRepository {
	return &MockTaskRepository{t: t}
}

func (m *MockTaskRepository) FindAll() (*[]domain.Task, error) {
	args := m.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	tasks, ok := result.(*[]domain.Task)
	assert.True(m.t, ok, "Expected *[]domain.Task")
	return tasks, args.Error(1)
}

func (m *MockTaskRepository) FindOne(id primitive.ObjectID) (*domain.Task, error) {
	assert.NotEmpty(m.t, id, "ID should not be empty")
	args := m.Called(id)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	task, ok := result.(*domain.Task)
	assert.True(m.t, ok, "Expected *domain.Task")
	return task, args.Error(1)
}

func (m *MockTaskRepository) UpdateOne(id primitive.ObjectID, updatedTask domain.Task) error {
	assert.NotEmpty(m.t, id, "ID should not be empty")
	assert.NotNil(m.t, updatedTask, "Updated task should not be nil")
	args := m.Called(id, updatedTask)
	return args.Error(0)
}

func (m *MockTaskRepository) DeleteOne(id primitive.ObjectID) (*domain.Task, error) {
	assert.NotEmpty(m.t, id, "ID should not be empty")
	args := m.Called(id)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	task, ok := result.(*domain.Task)
	assert.True(m.t, ok, "Expected *domain.Task")
	return task, args.Error(1)
}

func (m *MockTaskRepository) Save(task domain.Task) (*domain.Task, error) {
	assert.NotNil(m.t, task, "Task should not be nil")
	args := m.Called(task)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	savedTask, ok := result.(*domain.Task)
	assert.True(m.t, ok, "Expected *domain.Task")
	return savedTask, args.Error(1)
}
