package service

import (
	"context"
	"github.com/refandas/duit-api/model/web"
)

type UserService interface {
	Create(ctx context.Context, request web.UserCreateRequest) web.UserResponse
	Update(ctx context.Context, request web.UserUpdateRequest) web.UserResponse
	Delete(ctx context.Context, userId string)
	FindById(ctx context.Context, userId string) web.UserResponse
}
