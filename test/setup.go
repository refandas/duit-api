package test

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/refandas/duit-api/app"
	"github.com/refandas/duit-api/controller"
	"github.com/refandas/duit-api/helper"
	"github.com/refandas/duit-api/model/domain"
	"github.com/refandas/duit-api/repository"
	"github.com/refandas/duit-api/service"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

const testUserTableName = "TestUsers"
const testSpendingTableName = "TestSpending"

func setupTestDB(tableName string) *helper.DynamoDB {
	client := app.SetupClient(context.TODO())
	db := &helper.DynamoDB{Client: client}
	db.TableName = tableName

	if tableName == testUserTableName {
		app.CreateTable(context.Background(), db, app.CreateTableUser)
	}
	if tableName == testSpendingTableName {
		app.CreateTable(context.Background(), db, app.CreateTableSpending)
	}
	return db
}

func setupRouter(db *helper.DynamoDB) http.Handler {
	validate := validator.New()

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db, validate)
	userController := controller.NewUserController(userService)

	spendingRepository := repository.NewSpendingRepository()
	spendingService := service.NewSpendingService(spendingRepository, db, validate)
	spendingController := controller.NewSpendingController(spendingService)

	registerRouter := app.Router{
		UserController:     userController,
		SpendingController: spendingController,
	}
	router := registerRouter.NewRouter()
	return router
}

func clearUserDataAfterTest(db *helper.DynamoDB, id string) {
	userRepository := repository.NewUserRepository()
	userRepository.Delete(context.Background(), db, domain.User{
		Id: id,
	})
}

func clearSpendingDataAfterTest(db *helper.DynamoDB, id string) {
	spendingRepository := repository.NewSpendingRepository()
	spendingRepository.Delete(context.Background(), db, domain.Spending{
		Id: id,
	})
}

// createUser creates a user then return the user's data
func createUser(db *helper.DynamoDB) domain.User {
	userRepository := repository.NewUserRepository()
	userId, _ := uuid.NewRandom()
	password := []byte("secret")
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	user := userRepository.Save(context.Background(), db, domain.User{
		Id:        userId.String(),
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  string(hashedPassword),
		CreatedAt: time.Now().UnixMilli(),
	})
	return user
}

// createUser creates three users then return the user's data
func createUsers(db *helper.DynamoDB) []domain.User {
	userRepository := repository.NewUserRepository()
	password := []byte("secret")
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	users := []domain.User{
		domain.User{
			Name:      "User 1",
			Email:     "user1@example.com",
			Password:  string(hashedPassword),
			CreatedAt: time.Now().UnixMilli(),
		},
		domain.User{
			Name:      "User 2",
			Email:     "user2@example.com",
			Password:  string(hashedPassword),
			CreatedAt: time.Now().UnixMilli(),
		},
		domain.User{
			Name:      "User 3",
			Email:     "user3@example.com",
			Password:  string(hashedPassword),
			CreatedAt: time.Now().UnixMilli(),
		},
	}

	for i := 0; i < 3; i++ {
		userId, _ := uuid.NewRandom()
		users[i].Id = userId.String()
		users[i] = userRepository.Save(context.Background(), db, users[i])
	}
	return users
}

func createSpending(db *helper.DynamoDB, userId string) domain.Spending {
	spendingRepository := repository.NewSpendingRepository()
	spendingId, _ := uuid.NewRandom()

	spending := spendingRepository.Save(context.Background(), db, domain.Spending{
		Id:          spendingId.String(),
		UserId:      userId,
		Title:       "Makan malam",
		Date:        1701795600000,
		Amount:      50000,
		Category:    "food",
		Description: "Makan malam dengan sate kambing",
		CreatedAt:   time.Now().UnixMilli(),
	})
	return spending
}

func createSpendings(db *helper.DynamoDB, userId string) []domain.Spending {
	spendingRepository := repository.NewSpendingRepository()

	spendings := []domain.Spending{
		domain.Spending{
			UserId:      userId,
			Title:       "Makan malam",
			Date:        1702141200000,
			Amount:      25000,
			Category:    "food",
			Description: "Makan malam dengan ayam bakar",
			CreatedAt:   time.Now().UnixMilli(),
		},
		domain.Spending{
			UserId:      userId,
			Title:       "Makan malam",
			Date:        1702227600000,
			Amount:      35000,
			Category:    "food",
			Description: "Makan malam dengan nasi goreng",
			CreatedAt:   time.Now().UnixMilli(),
		},
		domain.Spending{
			UserId:      userId,
			Title:       "Makan malam",
			Date:        1702314000000,
			Amount:      50000,
			Category:    "food",
			Description: "Makan malam dengan sate kambing",
			CreatedAt:   time.Now().UnixMilli(),
		},
	}

	for i := 0; i < 3; i++ {
		spendingId, _ := uuid.NewRandom()
		spendings[i].Id = spendingId.String()
		spendings[i] = spendingRepository.Save(context.Background(), db, spendings[i])
	}
	return spendings
}
