package repository

import (
	"context"
	"github.com/refandas/duit-api/helper"
	"github.com/refandas/duit-api/model/domain"
)

type SpendingRepository interface {
	Save(ctx context.Context, db *helper.DynamoDB, spending domain.Spending) domain.Spending
	Update(ctx context.Context, db *helper.DynamoDB, spending domain.Spending) domain.Spending
	Delete(ctx context.Context, db *helper.DynamoDB, spending domain.Spending)
	FindById(ctx context.Context, db *helper.DynamoDB, spendingId string) (domain.Spending, error)
	FindByUserId(ctx context.Context, db *helper.DynamoDB, userId string) []domain.Spending
}
