package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/refandas/duit-api/app"
	"github.com/refandas/duit-api/controller"
	"github.com/refandas/duit-api/repository"
	"github.com/refandas/duit-api/service"
	"net/http"
)

func main() {
	db := app.SetupDatabase(context.Background())
	validate := validator.New()

	// Users configuration
	dbUsers := db
	dbUsers.TableName = "Users"
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, &dbUsers, validate)
	userController := controller.NewUserController(userService)

	// Spending configuration
	dbSpending := db
	dbSpending.TableName = "Spending"
	spendingRepository := repository.NewSpendingRepository()
	spendingService := service.NewSpendingService(spendingRepository, &dbSpending, validate)
	spendingController := controller.NewSpendingController(spendingService)

	router := app.Router{
		UserController:     userController,
		SpendingController: spendingController,
	}

	server := http.Server{
		Addr:    "localhost:8000",
		Handler: router.NewRouter(),
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
