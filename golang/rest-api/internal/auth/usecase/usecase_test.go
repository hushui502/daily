package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"rest-api/config"
	"rest-api/internal/auth/mock"
	"rest-api/internal/models"
	"rest-api/pkg/logger"
	"rest-api/pkg/utils"
)

func TestAuthUC_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockUseCase(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, nil, nil, apiLogger)
}
