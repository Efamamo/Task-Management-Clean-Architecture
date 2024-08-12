package repository

import (
	"context"
	"time"

	domain "github.com/Task-Management-go/Domain"
	"github.com/Task-Management-go/Domain/err"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository struct {
	collection mongo.Collection
}

func NewTaskRepo(client *mongo.Client) *TaskRepository {
	taskCollection := client.Database("task-management").Collection("tasks")
	return &TaskRepository{
		collection: *taskCollection,
	}
}

func (tr *TaskRepository) FindAll() (*[]domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, e := tr.collection.Find(ctx, bson.D{})
	if e != nil {
		return nil, err.NewUnexpected("Server Error")
	}

	tasks := make([]domain.Task, 0)
	for cur.Next(ctx) {
		var task domain.Task

		e := cur.Decode(&task)

		if e != nil {
			return nil, err.NewUnexpected("Server Error")
		}
		tasks = append(tasks, task)
	}

	if cur.Err() != nil {
		return nil, err.NewUnexpected("Server Error")
	}

	cur.Close(ctx)

	return &tasks, nil

}

func (tr *TaskRepository) FindOne(id primitive.ObjectID) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	var task domain.Task
	e := tr.collection.FindOne(ctx, filter).Decode(&task)
	if e != nil {
		if e == mongo.ErrNoDocuments {
			return nil, err.NewNotFound("task not found")
		}
		return nil, err.NewUnexpected(e.Error())
	}

	return &task, nil
}

func (tr *TaskRepository) UpdateOne(ID primitive.ObjectID, updatedTask domain.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": ID}

	update := bson.M{
		"$set": bson.M{
			"title":       updatedTask.Title,
			"description": updatedTask.Description,
			"due_date":    updatedTask.DueDate,
			"status":      updatedTask.Status,
		},
	}

	// Update the document that matches the filter
	_, e := tr.collection.UpdateOne(ctx, filter, update)
	if e != nil {
		return err.NewUnauthorized("Server Error")
	}
	return nil

}

func (tr *TaskRepository) DeleteOne(ID primitive.ObjectID) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	task, e := tr.FindOne(ID)
	if e != nil {
		return nil, err.NewNotFound("Task Not Found")
	}

	filter := bson.M{"_id": ID}

	// Delete the document that matches the filter
	_, e = tr.collection.DeleteOne(ctx, filter)

	if e != nil {
		return nil, err.NewUnexpected("iServer Error")
	}

	return task, nil
}

func (tr *TaskRepository) Save(task domain.Task) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for {
		task.ID = primitive.NewObjectID()

		_, e := tr.collection.InsertOne(ctx, task)

		if mongo.IsDuplicateKeyError(e) {
			continue
		} else if e != nil {
			return nil, err.NewUnexpected("Server Error")
		}

		return &task, nil
	}
}
