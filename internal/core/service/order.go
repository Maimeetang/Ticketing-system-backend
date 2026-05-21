package service

import (
	"fmt"
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"
	"time"

	"github.com/google/uuid"
)

type orderServiceImpl struct {
	repo port.OrderRepository
}

func NewOrderService(repo port.OrderRepository) port.OrderService {
	return &orderServiceImpl{repo: repo}
}

func (s *orderServiceImpl) CreateOrder(userID uint, shiftID uint, paymentMethod domain.PaymentMethod, inputItems []domain.OrderItem) (*domain.Order, error) {
	if len(inputItems) == 0 {
		return nil, apperror.NewBadRequest("cannot process empty order items")
	}

	var totalAmount float64
	var finalOrderItems []domain.OrderItem

	// Verify product prices and calculate amounts securely
	for _, item := range inputItems {
		if item.Quantity <= 0 {
			return nil, apperror.NewBadRequest("item quantity must be greater than zero")
		}

		// Fetch official ticket configuration from master table to get genuine unit price
		ticketType, err := s.repo.GetTicketTypeByID(item.TicketTypeID)
		if err != nil {
			return nil, apperror.NewBadRequest(fmt.Sprintf("invalid ticket type ID: %d", item.TicketTypeID))
		}

		subtotal := ticketType.Price * float64(item.Quantity)
		totalAmount += subtotal

		orderItem := domain.OrderItem{
			TicketTypeID: ticketType.ID,
			Quantity:     item.Quantity,
			UnitPrice:    ticketType.Price,
			Subtotal:     subtotal,
		}
		finalOrderItems = append(finalOrderItems, orderItem)
	}

	// Instantiate structural Order document schema with dynamic UUID
	newOrder := &domain.Order{
		UUID:          uuid.New().String(), // Ensure unique index identifier across distributed POS machines
		ShiftID:       shiftID,
		TotalAmount:   totalAmount,
		PaymentMethod: paymentMethod,
		SyncStatus:    domain.SyncPending,
		CreatedBy:     userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		OrderItems:    finalOrderItems,
	}

	// Persist into database storage layer
	if err := s.repo.Create(newOrder); err != nil {
		return nil, err
	}

	return newOrder, nil
}