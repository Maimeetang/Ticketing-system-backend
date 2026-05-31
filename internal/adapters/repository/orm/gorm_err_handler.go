package orm

import (
	"errors"
	"ticketing-system/internal/apperror"

	"gorm.io/gorm"
)

func handleError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return apperror.NewInternalServerError("database error: " + err.Error())
}