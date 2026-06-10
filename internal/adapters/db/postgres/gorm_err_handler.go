package orm

import (
	"errors"
	e "ticketing-system/internal/core/error"

	"gorm.io/gorm"
)

func handleError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return e.NewInternalServerError("database error: " + err.Error())
}