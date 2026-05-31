package orm

import (
	"strings"
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
		return nil, handleError(err)
	}
	return order, nil
}

func (r *GormOrderRepository) GetOrderByID(id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Tickets").
				Preload("Tickets.TicketInfos").
				Preload("Tickets.TicketLogs").
				First(&order, id).Error
	if err != nil {
		return nil, handleError(err)
	}
	return &order, nil
}

func (r *GormOrderRepository) ListOrders(filter domain.OrderFilter) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var totalCount int64

	// initial query
	query := r.db.Model(&domain.Order{})
	if filter.IncludeTickets{
		query = query.Preload("Tickets").Preload("Tickets.TicketInfos")
	}

	// add filters to query
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.PaymentMethod != nil {
		query = query.Where("payment_method = ?", *filter.PaymentMethod)
	}
	if filter.CashierID != nil {
		query = query.Where("cashier_id = ?", *filter.CashierID)
	}
	if filter.ShiftID != nil {
		query = query.Where("shift_id = ?", *filter.ShiftID)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, handleError(err)
	}

	// sort
	sort := "created_at DESC"
	if strings.ToUpper(filter.Sort) == "ASC" {
		sort = "created_at ASC"
	}

	// List orders
	offset := (filter.Page - 1) * filter.Limit
	err := query.Order(sort).
		Limit(filter.Limit).Offset(offset).Find(&orders).Error

	if err != nil {
		return nil, 0, handleError(err)
	}

	return orders, totalCount, nil
}

func (r *GormOrderRepository) UpdateOrder(order *domain.Order) (*domain.Order, error) {
	err := r.db.Save(order).Error
	if err != nil {
		return nil, handleError(err)
	}
	return order, nil
}