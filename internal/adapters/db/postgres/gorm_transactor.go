package orm

import (
	"context"

	"gorm.io/gorm"
)

type GormATransactor struct {
	db *gorm.DB
}

func NewGormTransactor(db *gorm.DB) *GormATransactor {
	return &GormATransactor{db: db}
}

func (t *GormATransactor) WithTransaction(
	ctx context.Context, 
	fn func(ctx context.Context) error,
) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}
