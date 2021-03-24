package server

import (
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/minio/minio-go/v7"
	"rest-api/config"
	"rest-api/pkg/logger"
)

const (
	certFile = "ssl/Server.crt"
	keyFile = "ssl/Server.pem"
	maxHeaderBytes = 1 << 20
	ctxTimeout = 5
)

// Server struct
type Server struct {
	echo *echo.Echo
	cfg *config.Config
	db *sqlx.DB
	redisClient *redis.Client
	awsClient *minio.Client
	logger logger.Logger
}

// NewServer New Server constructor
