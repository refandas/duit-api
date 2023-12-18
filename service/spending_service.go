package service

import (
	"context"
	"github.com/refandas/duit-api/model/web"
)

type SpendingService interface {
	Create(ctx context.Context, request web.SpendingCreateRequest) web.SpendingResponse
	Update(ctx context.Context, request web.SpendingUpdateRequest) web.SpendingResponse
	Delete(ctx context.Context, spendingId string)
	FindById(ctx context.Context, spendingId string) web.SpendingResponse
	FindByUserId(ctx context.Context, userId string) []web.SpendingResponse
}
