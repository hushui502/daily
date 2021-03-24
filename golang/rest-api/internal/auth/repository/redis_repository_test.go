package repository

import (
	"context"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"log"
	"rest-api/internal/auth"
	"rest-api/internal/models"
	"testing"
)

func SetupRedis() auth.RedisRepository {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	authRedisRepo := NewAuthRedisRepo(client)
	return authRedisRepo
}

func TestNewAuthRedisRepo_GetByIDCtx(t *testing.T) {
	t.Parallel()

	authRedisRepo := SetupRedis()

	t.Run("GetByIDCtx", func(t *testing.T) {
		key := uuid.New().String()
		userID := uuid.New()
		u := &models.User{
			UserID:    userID,
			FirstName: "Alex",
			LastName:  "Bryksin",
		}

		err := authRedisRepo.SetUserCtx(context.Background(), key, 10, u)
		require.NoError(t, err)
		require.Nil(t, err)

		user, err := authRedisRepo.GetByIDCtx(context.Background(), key)
		require.NoError(t, err)
		require.NotNil(t, user)
	})
}

func TestNewAuthRedisRepo_SetUserCtx(t *testing.T) {
	t.Parallel()

	authRedisRepo := SetupRedis()

	t.Run("SetUserCtx", func(t *testing.T) {
		key := uuid.New().String()
		userID := uuid.New()
		u := &models.User{
			UserID: userID,
			FirstName: "Alex",
			LastName:  "Bryksin",
		}

		err := authRedisRepo.SetUserCtx(context.Background(), key, 10, u)
		require.NoError(t, err)
		require.Nil(t, err)
	})
}






































