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
	err := r.db.Preload("Ticket").
				Preload("Ticket.TicketInfos").
				Preload("Ticket.TicketLogs").
				First(&order, id).Error
	if err != nil {
		return nil, handleError(err)
	}
	return &order, nil
}

func (r *GormOrderRepository) GetByTicketCode(code string) (*domain.Order, error) {
	var order domain.Order

	err := r.db.
		Preload("Ticket").
		Preload("Ticket.TicketLogs").
		Joins("Ticket").
		Where("ticket_code = ?", code).
		First(&order).Error

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
		query = query.Preload("Ticket").Preload("Ticket.TicketInfos")
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
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.Order{}).
			Where("id = ?", order.ID).
			Updates(map[string]any{
				"status":      order.Status,
				"total_price": order.TotalPrice,
			}).Error; err != nil {
			return err
		}

		if err := tx.Model(&domain.Ticket{}).
			Where("id = ?", order.Ticket.ID).
			Updates(map[string]any{
				"status":      order.Ticket.Status,
				"total_price": order.Ticket.TotalPrice,
			}).Error; err != nil {
			return err
		}

		var newLogs []domain.TicketLog
		for _, log := range order.Ticket.TicketLogs {
			if log.ID == 0 {
				log.TicketID = order.Ticket.ID
				newLogs = append(newLogs, log)
			}
		}

		if len(newLogs) > 0 {
			if err := tx.Create(&newLogs).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, handleError(err)
	}

	return order, nil
}