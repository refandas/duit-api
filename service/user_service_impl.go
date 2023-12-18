package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/refandas/duit-api/helper"
	"github.com/refandas/duit-api/model/domain"
	"github.com/refandas/duit-api/model/web"
	"github.com/refandas/duit-api/repository"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *helper.DynamoDB
	Validate       *validator.Validate
}

func NewUserService(userRepository repository.UserRepository, DB *helper.DynamoDB, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DB:             DB,
		Validate:       validate,
	}
}

func (service *UserServiceImpl) Create(ctx context.Context, request web.UserCreateRequest) web.UserResponse {
	err := service.Validate.Struct(request)
	if err != nil {
		panic(err)
	}

	user := domain.User{
		Id:        request.Id,
		Name:      request.Name,
		Email:     request.Email,
		Password:  request.Password,
		CreatedAt: request.CreatedAt,
	}

	userResponse := service.UserRepository.Save(ctx, service.DB, user)
	return helper.ToUserResponse(userResponse)
}

func (service *UserServiceImpl) Update(ctx context.Context, request web.UserUpdateRequest) web.UserResponse {
	err := service.Validate.Struct(request)
	if err != nil {
		panic(err)
	}

	user, err := service.UserRepository.FindById(ctx, service.DB, request.Id)
	if err != nil {
		panic(err)
	}

	// update field
	if request.Name != "" {
		user.Name = request.Name
	}
	if request.Email != "" {
		user.Email = request.Email
	}
	if request.Password != "" {
		user.Password = request.Password
	}

	response := service.UserRepository.Update(ctx, service.DB, user)
	return helper.ToUserResponse(response)
}

func (service *UserServiceImpl) Delete(ctx context.Context, userId string) {
	user, err := service.UserRepository.FindById(ctx, service.DB, userId)
	if err != nil {
		panic(err)
	}
	service.UserRepository.Delete(ctx, service.DB, user)
}

func (service *UserServiceImpl) FindById(ctx context.Context, userId string) web.UserResponse {
	user, err := service.UserRepository.FindById(ctx, service.DB, userId)
	if err != nil {
		panic(err)
	}
	return helper.ToUserResponse(user)
}
