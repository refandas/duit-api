package helper

import (
	"github.com/refandas/duit-api/model/domain"
	"github.com/refandas/duit-api/model/web"
)

// ToUserResponse converts a domain.User struct to a web.UserResponse struct.
//
// This function is commonly used to transform user information retrieved
// from process the data on the database into a format suitable for sending
// back as a response in API endpoints.
func ToUserResponse(user domain.User) web.UserResponse {
	return web.UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

// ToUserResponses converts a slice of domain.User struct to a slice of
// web.UserResponse struct.
//
// This function is commonly used to transform multiple user records retrieved
// from process the data on the database into a format suitable for sending
// back as a response in API endpoints.
func ToUserResponses(users []domain.User) []web.UserResponse {
	var userResponses []web.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, ToUserResponse(user))
	}
	return userResponses
}

// ToSpendingResponse converts a domain.Spending struct to a
// web.SpendingResponse struct.
//
// This function is commonly used to transform user's spending information
// retrieved from process the data on the database into a format suitable for
// sending back as a response in API endpoints.
func ToSpendingResponse(spending domain.Spending) web.SpendingResponse {
	return web.SpendingResponse{
		Id:          spending.Id,
		UserId:      spending.UserId,
		Title:       spending.Title,
		Amount:      spending.Amount,
		Description: spending.Description,
		Category:    spending.Category,
		Date:        spending.Date,
		CreatedAt:   spending.CreatedAt,
	}
}

// ToSpendingResponses converts a slice of domain.Spending struct to a
// slice of web.SpendingResponse struct.
//
// This function is commonly used to transform multiple user's spending records
// retrieved from process the data on the database into a format suitable for
// sending back as a response in API endpoints.
func ToSpendingResponses(spendings []domain.Spending) []web.SpendingResponse {
	var spendingResponses []web.SpendingResponse
	for _, spending := range spendings {
		spendingResponses = append(spendingResponses, ToSpendingResponse(spending))
	}
	return spendingResponses
}
