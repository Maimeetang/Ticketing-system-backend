package repository

import "ticketing-system/internal/core/model"

type UserRepository interface {
	Add(user model.User) error
	Edit(user model.User) error
	DeleteByID(id uint) error
	FindByID(id uint) (model.User, error)
	FindByUsername(username string) (model.User, error)
}
