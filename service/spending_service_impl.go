package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/refandas/duit-api/helper"
	"github.com/refandas/duit-api/model/domain"
	"github.com/refandas/duit-api/model/web"
	"github.com/refandas/duit-api/repository"
)

type SpendingServiceImpl struct {
	SpendingRepository repository.SpendingRepository
	DB                 *helper.DynamoDB
	Validator          *validator.Validate
}

func NewSpendingService(spendingRepository repository.SpendingRepository, DB *helper.DynamoDB, validator *validator.Validate) SpendingService {
	return &SpendingServiceImpl{
		SpendingRepository: spendingRepository,
		DB:                 DB,
		Validator:          validator,
	}
}

func (service *SpendingServiceImpl) Create(ctx context.Context, request web.SpendingCreateRequest) web.SpendingResponse {
	err := service.Validator.Struct(request)
	if err != nil {
		panic(err)
	}

	spending := domain.Spending{
		Id:          request.Id,
		UserId:      request.UserId,
		Title:       request.Title,
		Description: request.Description,
		Category:    request.Category,
		Date:        request.Date,
		Amount:      request.Amount,
		CreatedAt:   request.CreatedAt,
	}

	spendingResponse := service.SpendingRepository.Save(ctx, service.DB, spending)
	return helper.ToSpendingResponse(spendingResponse)
}

func (service *SpendingServiceImpl) Update(ctx context.Context, request web.SpendingUpdateRequest) web.SpendingResponse {
	err := service.Validator.Struct(request)
	if err != nil {
		panic(err)
	}

	spending, err := service.SpendingRepository.FindById(ctx, service.DB, request.Id)
	if err != nil {
		panic(err)
	}

	spending.Title = request.Title
	spending.Date = request.Date
	spending.Description = request.Description
	spending.Amount = request.Amount
	spending.Category = request.Category

	response := service.SpendingRepository.Update(ctx, service.DB, spending)
	return helper.ToSpendingResponse(response)
}

func (service *SpendingServiceImpl) Delete(ctx context.Context, spendingId string) {
	spending, err := service.SpendingRepository.FindById(ctx, service.DB, spendingId)
	if err != nil {
		panic(err)
	}
	service.SpendingRepository.Delete(ctx, service.DB, spending)
}

func (service *SpendingServiceImpl) FindById(ctx context.Context, spendingId string) web.SpendingResponse {
	spending, err := service.SpendingRepository.FindById(ctx, service.DB, spendingId)
	if err != nil {
		panic(err)
	}
	return helper.ToSpendingResponse(spending)
}

func (service *SpendingServiceImpl) FindByUserId(ctx context.Context, userId string) []web.SpendingResponse {
	spendings := service.SpendingRepository.FindByUserId(ctx, service.DB, userId)
	return helper.ToSpendingResponses(spendings)
}
