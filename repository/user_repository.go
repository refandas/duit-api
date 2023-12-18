package repository

import (
	"context"
	"github.com/refandas/duit-api/helper"
	"github.com/refandas/duit-api/model/domain"
)

type UserRepository interface {
	Save(ctx context.Context, db *helper.DynamoDB, user domain.User) domain.User
	Update(ctx context.Context, db *helper.DynamoDB, user domain.User) domain.User
	Delete(ctx context.Context, db *helper.DynamoDB, user domain.User)
	FindById(ctx context.Context, db *helper.DynamoDB, userId string) (domain.User, error)
}
