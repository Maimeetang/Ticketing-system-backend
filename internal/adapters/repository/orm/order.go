package orm

import (
	"errors"
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"

	"gorm.io/gorm"
)

type GormOrderRepository struct {
	db *gorm.DB
}

func NewGormOrderRepository(db *gorm.DB) port.OrderRepository {
	return &GormOrderRepository{db: db}
}

func (r *GormOrderRepository) CreateOrder(order *domain.Order) (*domain.Order, error) {
	err := r.db.Create(order).Error
	if err != nil {
		return nil, handleOrderError(err)
	}
	return order, nil
}

func (r *GormOrderRepository) GetOrderByID(id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.First(&order, id).Error
	if err != nil {
		return nil, handleOrderError(err)
	}
	return &order, nil
}

func (r *GormOrderRepository) ListOrders() ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Find(&orders).Error
	if err != nil {
		return nil, handleOrderError(err)
	}
	return orders, nil
}

func (r *GormOrderRepository) UpdateOrder(order *domain.Order) (*domain.Order, error) {
	err := r.db.Save(order).Error
	if err != nil {
		return nil, handleOrderError(err)
	}
	return order, nil
}

func handleOrderError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.NewNotFound("requested database records not found")
	}
	return apperror.NewInternalServerError("database error: " + err.Error())
}