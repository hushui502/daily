package usecase

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"rest-api/config"
	"rest-api/internal/models"
	"rest-api/internal/session"
)

type sessionUC struct {
	sessionRepo session.SessRepository
	cfg *config.Config
}

func NewSessionUseCase(sessionRepo session.SessRepository, cfg *config.Config) session.UCSession {
	return &sessionUC{sessionRepo: sessionRepo, cfg: cfg}
}

func (u *sessionUC) CreateSession(ctx context.Context, session *models.Session, expire int) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionUC.CreateSession")
	defer span.Finish()

	return u.sessionRepo.CreateSession(ctx, session, expire)
}

func (u *sessionUC) DeleteByID(ctx context.Context, sessionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionUC.DeleteByID")
	defer span.Finish()

	return u.sessionRepo.DeleteByID(ctx, sessionID)
}

func (u *sessionUC) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionUC.GetSessionByID")
	defer span.Finish()

	return u.sessionRepo.GetSessionByID(ctx, sessionID)
}
