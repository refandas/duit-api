package test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateSpendingSuccess(t *testing.T) {
	userDb := setupTestDB(testUserTableName)
	spendingDb := setupTestDB(testSpendingTableName)
	router := setupRouter(spendingDb)

	user := createUser(userDb)
	defer clearUserDataAfterTest(userDb, user.Id)

	// Create the spending data
	jsonData := `
	{
		"user_id": "%s",
		"amount": 50000,
		"date": 1701795600000,
		"category": "food",
		"title": "Makan malam",
		"description": "Makan malam dengan sate kambing"
	}
`
	jsonData = fmt.Sprintf(jsonData, user.Id)

	requestBody := strings.NewReader(jsonData)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8000/api/v1/spendings", requestBody)
	request.Header.Add("Content-Type", "application/json")

	// Apply mock testing
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
	spendingId := responseBody["data"].(map[string]interface{})["id"]
	defer clearSpendingDataAfterTest(spendingDb, spendingId.(string))

	assert.Equal(t, http.StatusCreated, int(responseBody["code"].(float64)))
	assert.Equal(t, "CREATED", responseBody["status"])
	assert.Equal(t, user.Id, responseBody["data"].(map[string]interface{})["user_id"])
	assert.Equal(t, "Makan malam", responseBody["data"].(map[string]interface{})["title"])
	assert.Equal(t, float64(50000), responseBody["data"].(map[string]interface{})["amount"])
	assert.Equal(t, int64(1701795600000), int64(responseBody["data"].(map[string]interface{})["date"].(float64)))
	assert.Equal(t, "food", responseBody["data"].(map[string]interface{})["category"])
	assert.Equal(t, "Makan malam dengan sate kambing", responseBody["data"].(map[string]interface{})["description"])
}

func TestCreateSpendingFailed(t *testing.T) {
	userDb := setupTestDB(testUserTableName)
	spendingDb := setupTestDB(testSpendingTableName)
	router := setupRouter(spendingDb)

	user := createUser(userDb)
	defer clearUserDataAfterTest(userDb, user.Id)

	// Create the spending data
	jsonData := `
	{
		"user_id": "%s",
		"amount": 50000,
		"date": 1701795600000,
		"category": "food",
		"title": "",
		"description": "Makan malam dengan sate kambing"
	}
`
	jsonData = fmt.Sprintf(jsonData, user.Id)

	requestBody := strings.NewReader(jsonData)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8000/api/v1/spendings", requestBody)
	request.Header.Add("Content-Type", "application/json")

	// Apply mock testing
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

func TestUpdateSpendingSuccess(t *testing.T) {
	userDb := setupTestDB(testUserTableName)
	spendingDb := setupTestDB(testSpendingTableName)

	router := setupRouter(spendingDb)

	user := createUser(userDb)
	defer clearUserDataAfterTest(userDb, user.Id)

	spending := createSpending(spendingDb, user.Id)
	defer clearSpendingDataAfterTest(spendingDb, spending.Id)

	jsonData := `
	{
		"user_id": "%s",
		"amount": %d,
		"date": 1701795600000,
		"category": "food",
		"title": "Makan malam",
		"description": "%s"
	}
`
	updatedAmount := 75000
	updatedDescription := "Makan malam dengan sate kambing dan es teh"

	jsonData = fmt.Sprintf(jsonData, user.Id, updatedAmount, updatedDescription)

	requestBody := strings.NewReader(jsonData)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:8000/api/v1/spendings/"+spending.Id, requestBody)
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
	assert.Equal(t, user.Id, responseBody["data"].(map[string]interface{})["user_id"])
	assert.Equal(t, "Makan malam", responseBody["data"].(map[string]interface{})["title"])
	assert.Equal(t, float64(updatedAmount), responseBody["data"].(map[string]interface{})["amount"])
	assert.Equal(t, int64(1701795600000), int64(responseBody["data"].(map[string]interface{})["date"].(float64)))
	assert.Equal(t, "food", responseBody["data"].(map[string]interface{})["category"])
	assert.Equal(t, updatedDescription, responseBody["data"].(map[string]interface{})["description"])
}

func TestUpdateSpendingFailed(t *testing.T) {
	userDb := setupTestDB(testUserTableName)
	spendingDb := setupTestDB(testSpendingTableName)

	router := setupRouter(spendingDb)

	user := createUser(userDb)
	defer clearUserDataAfterTest(userDb, user.Id)

	spending := createSpending(spendingDb, user.Id)
	defer clearSpendingDataAfterTest(spendingDb, spending.Id)

	jsonData := `
	{
		"user_id": "%s",
		"amount": %d,
		"date": 1701795600000,
		"category": "food",
		"title": "",
		"description": "%s"
	}
`
	updatedAmount := 75000
	updatedDescription := "Makan malam dengan sate kambing dan es teh"

	jsonData = fmt.Sprintf(jsonData, user.Id, updatedAmount, updatedDescription)

	requestBody := strings.NewReader(jsonData)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:8000/api/v1/spendings/"+spending.Id, requestBody)
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

func TestGetSpendingSuccess(t *testing.T) {
	userDb := setupTestDB(testUserTableName)
	spendingDb := setupTestDB(testSpendingTableName)

	// create user data
	user := createUser(userDb)
	defer clearUserDataAfterTest(userDb, user.Id)

	// create user's spending data
	spending := createSpending(spendingDb, user.Id)
	defer clearSpendingDataAfterTest(spendingDb, spending.Id)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8000/api/v1/spendings/"+spending.Id, nil)
	recorder := httptest.NewRecorder()

	router := setupRouter(spendingDb)
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
	assert.Equal(t, user.Id, responseBody["data"].(map[string]interface{})["user_id"])
	assert.Equal(t, "Makan malam", responseBody["data"].(map[string]interface{})["title"])
	assert.Equal(t, float64(50000), responseBody["data"].(map[string]interface{})["amount"])
	assert.Equal(t, int64(1701795600000), int64(responseBody["data"].(map[string]interface{})["date"].(float64)))
	assert.Equal(t, "food", responseBody["data"].(map[string]interface{})["category"])
	assert.Equal(t, "Makan malam dengan sate kambing", responseBody["data"].(map[string]interface{})["description"])
}

func TestGetSpendingFailed(t *testing.T) {
	spendingDb := setupTestDB(testSpendingTableName)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8000/api/v1/spendings/100", nil)
	recorder := httptest.NewRecorder()

	router := setupRouter(spendingDb)
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

// The route to be tested is /api/v1/{user_id}/spendings
func TestGetListOfUserSpendingSuccess(t *testing.T) {
	userDb := setupTestDB(testUserTableName)
	spendingDb := setupTestDB(testSpendingTableName)

	user := createUser(userDb)
	defer clearUserDataAfterTest(userDb, user.Id)

	spendings := createSpendings(spendingDb, user.Id)
	defer clearSpendingDataAfterTest(spendingDb, spendings[0].Id)
	defer clearSpendingDataAfterTest(spendingDb, spendings[1].Id)
	defer clearSpendingDataAfterTest(spendingDb, spendings[2].Id)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8000/api/v1/users/"+user.Id+"/spendings", nil)
	recorder := httptest.NewRecorder()

	router := setupRouter(spendingDb)
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])

	spendingResponses := responseBody["data"].([]interface{})

	spendingResponse1 := spendingResponses[0].(map[string]interface{})
	spendingResponse2 := spendingResponses[1].(map[string]interface{})
	spendingResponse3 := spendingResponses[2].(map[string]interface{})

	assert.Equal(t, spendings[0].Id, spendingResponse1["id"])
	assert.Equal(t, spendings[0].UserId, spendingResponse1["user_id"])
	assert.Equal(t, spendings[0].Title, spendingResponse1["title"])
	assert.Equal(t, spendings[0].Description, spendingResponse1["description"])
	assert.Equal(t, spendings[0].Category, spendingResponse1["category"])
	assert.Equal(t, spendings[0].Amount, spendingResponse1["amount"])
	assert.Equal(t, spendings[0].Date, int64(spendingResponse1["date"].(float64)))

	assert.Equal(t, spendings[1].Id, spendingResponse2["id"])
	assert.Equal(t, spendings[1].UserId, spendingResponse2["user_id"])
	assert.Equal(t, spendings[1].Title, spendingResponse2["title"])
	assert.Equal(t, spendings[1].Description, spendingResponse2["description"])
	assert.Equal(t, spendings[1].Category, spendingResponse2["category"])
	assert.Equal(t, spendings[1].Amount, spendingResponse2["amount"])
	assert.Equal(t, spendings[1].Date, int64(spendingResponse2["date"].(float64)))

	assert.Equal(t, spendings[2].Id, spendingResponse3["id"])
	assert.Equal(t, spendings[2].UserId, spendingResponse3["user_id"])
	assert.Equal(t, spendings[2].Title, spendingResponse3["title"])
	assert.Equal(t, spendings[2].Description, spendingResponse3["description"])
	assert.Equal(t, spendings[2].Category, spendingResponse3["category"])
	assert.Equal(t, spendings[2].Amount, spendingResponse3["amount"])
	assert.Equal(t, spendings[2].Date, int64(spendingResponse3["date"].(float64)))
}

// TestGetListOfUserSpendingFailed test to get user's spending data
// but the user's id is not found.
// The route to be tested is /api/v1/{user_id}/spendings
func TestGetListOfUserSpendingFailed(t *testing.T) {
	spendingDb := setupTestDB(testSpendingTableName)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8000/api/v1/users/404/spendings", nil)
	recorder := httptest.NewRecorder()

	router := setupRouter(spendingDb)
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusNotFound, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestDeleteSpendingSuccess(t *testing.T) {
	userDb := setupTestDB(testUserTableName)
	spendingDb := setupTestDB(testSpendingTableName)

	// create user data
	user := createUser(userDb)
	defer clearUserDataAfterTest(userDb, user.Id)

	// create user's spending data
	spending := createSpending(spendingDb, user.Id)
	defer clearSpendingDataAfterTest(spendingDb, spending.Id)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8000/api/v1/spendings/"+spending.Id, nil)
	request.Header.Add("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router := setupRouter(spendingDb)
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

	assert.Equal(t, http.StatusNoContent, int(responseBody["code"].(float64)))
	assert.Equal(t, "DELETED", responseBody["status"])
}

func TestDeleteSpendingFailed(t *testing.T) {
	spendingDb := setupTestDB(testSpendingTableName)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8000/api/v1/spendings/100", nil)
	request.Header.Add("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router := setupRouter(spendingDb)
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
