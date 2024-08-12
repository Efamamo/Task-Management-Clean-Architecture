package main

import (
	"context"
	"log"
	"os"

	controller "github.com/Task-Management-go/Delivery/controllers"
	"github.com/Task-Management-go/Delivery/router"
	repository "github.com/Task-Management-go/Repository"
	usecases "github.com/Task-Management-go/Usecases"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	er := godotenv.Load("../.env")
	if er != nil {
		log.Fatal("Error loading .env file")
	}

	clientOptions := options.Client().ApplyURI(os.Getenv("ConnectionString"))

	// Connect to MongoDB
	var err error
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	var TaskRepository usecases.ITaskRepository = repository.NewTaskRepo(client)
	taskService := usecases.TaskService{TaskRepo: TaskRepository}
	taskController := controller.TaskController{Service: taskService}

	var UserRepository usecases.IUserRepository = repository.NewUserRepo(client)
	UserService := usecases.UserService{UserRepo: UserRepository}
	userController := controller.UserController{Service: UserService}

	router.SetUpRouter(taskController, userController)

}
