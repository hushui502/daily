//go:generate mockgen -source redis_repository.go -destination mock/redis_repository_mock.go -package mock
package auth

import (
	"context"
	"rest-api/internal/models"
)

type RedisRepository interface {
	GetByIDCtx(ctx context.Context, key string) (*models.User, error)
	SetUserCtx(ctx context.Context, key string, seconds int, user *models.User) error
	DeleteUserCtx(ctx context.Context, key string) error
}
