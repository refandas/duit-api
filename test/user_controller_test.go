package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/refandas/duit-api/model/domain"
	"github.com/refandas/duit-api/repository"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCreateUserSuccess(t *testing.T) {
	db := setupTestDB(testUserTableName)
	router := setupRouter(db)

	jsonData := `
	{
		"name": "Test User",
		"email": "test@example.com",
		"password": "secret"
	}
`
	requestBody := strings.NewReader(jsonData)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8000/api/v1/users", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}
	userId := responseBody["data"].(map[string]interface{})["id"]
	defer clearUserDataAfterTest(db, userId.(string))

	assert.Equal(t, http.StatusCreated, int(responseBody["code"].(float64)))
	assert.Equal(t, "CREATED", responseBody["status"])
	assert.Equal(t, "Test User", responseBody["data"].(map[string]interface{})["name"])
	assert.Equal(t, "test@example.com", responseBody["data"].(map[string]interface{})["email"])
}

func TestCreateUserFailed(t *testing.T) {
	db := setupTestDB(testUserTableName)
	router := setupRouter(db)

	jsonData := `
	{
		"name": "Test User",
		"email": "",
		"password": "secret"
	}
`
	requestBody := strings.NewReader(jsonData)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8000/api/v1/users", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusBadRequest, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestUpdateUserSuccess(t *testing.T) {
	db := setupTestDB(testUserTableName)

	router := setupRouter(db)

	user := createUser(db)
	defer clearUserDataAfterTest(db, user.Id)

	jsonData := `
		{
			"name": "%s",
			"email": "%s",
			"password": "%s"
		}
	`
	updatedName := "Test User 2"
	updatedEmail := "another@exmple.com"
	password := []byte("supersecret")

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	jsonData = fmt.Sprintf(jsonData, updatedName, updatedEmail, hashedPassword)

	requestBody := strings.NewReader(jsonData)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:8000/api/v1/users/"+user.Id, requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, user.Id, responseBody["data"].(map[string]interface{})["id"])
	assert.Equal(t, updatedName, responseBody["data"].(map[string]interface{})["name"])
	assert.Equal(t, updatedEmail, responseBody["data"].(map[string]interface{})["email"])
}

func TestUpdateUserWithoutPasswordAttributeSuccess(t *testing.T) {
	db := setupTestDB(testUserTableName)

	user := createUser(db)
	defer clearUserDataAfterTest(db, user.Id)

	router := setupRouter(db)

	jsonData := `
	{
		"name": "%s",
		"email": "%s"
	}
`
	updatedName := "Test User 2"
	updatedEmail := "another@exmple.com"
	jsonData = fmt.Sprintf(jsonData, updatedName, updatedEmail)

	requestBody := strings.NewReader(jsonData)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:8000/api/v1/users/"+user.Id, requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, user.Id, responseBody["data"].(map[string]interface{})["id"])
	assert.Equal(t, updatedName, responseBody["data"].(map[string]interface{})["name"])
	assert.Equal(t, updatedEmail, responseBody["data"].(map[string]interface{})["email"])
}

func TestUpdateUserFailed(t *testing.T) {
	db := setupTestDB(testUserTableName)

	user := createUser(db)
	defer clearUserDataAfterTest(db, user.Id)

	router := setupRouter(db)

	jsonData := `
	{
		"name": "",
		"email": "another@example.com"
	}
`
	jsonData = fmt.Sprintf(jsonData)

	requestBody := strings.NewReader(jsonData)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:8000/api/v1/users/"+user.Id, requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusBadRequest, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestGetUserSuccess(t *testing.T) {
	db := setupTestDB(testUserTableName)

	user := createUser(db)
	defer clearUserDataAfterTest(db, user.Id)

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8000/api/v1/users/"+user.Id, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, user.Id, responseBody["data"].(map[string]interface{})["id"])
	assert.Equal(t, user.Name, responseBody["data"].(map[string]interface{})["name"])
	assert.Equal(t, user.Email, responseBody["data"].(map[string]interface{})["email"])
}

func TestGetUserFailed(t *testing.T) {
	db := setupTestDB(testUserTableName)

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8000/api/v1/users/404", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusNotFound, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestDeleteUserSuccess(t *testing.T) {
	db := setupTestDB(testUserTableName)

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

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8000/api/v1/users/"+user.Id, nil)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return
	}

	assert.Equal(t, http.StatusNoContent, int(responseBody["code"].(float64)))
	assert.Equal(t, "DELETED", responseBody["status"])
}

func TestDeleteUserFailed(t *testing.T) {
	db := setupTestDB(testUserTableName)

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8000/api/v1/users/404", nil)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusNotFound, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}
