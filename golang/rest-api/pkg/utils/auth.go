package utils

import (
	"context"
	"rest-api/pkg/httpErrors"
	"rest-api/pkg/logger"
)

func ValidateIsOwner(ctx context.Context, creatorID string, logger logger.Logger) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return err
	}

	if user.UserID.String() != creatorID {
		logger.Errorf(
			"ValidateIsOwner, userID: %v, creatorID: %V",
			user.UserID.String(),
			creatorID,
		)

		return httpErrors.Forbidden
	}

	return nil
}
