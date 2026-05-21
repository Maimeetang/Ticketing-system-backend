package orm

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"
	"time"

	"gorm.io/gorm"
)

type GormOrderRepository struct {
	db *gorm.DB
}

func NewGormOrderRepository(db *gorm.DB) port.OrderRepository {
	return &GormOrderRepository{db: db}
}

func (r *GormOrderRepository) GetTicketTypeByID(id uint) (*domain.TicketType, error) {
	var ticketType domain.TicketType
	err := r.db.First(&ticketType, id).Error
	if err != nil {
		return nil, handleOrderError(err)
	}
	return &ticketType, nil
}

// Create handles the full transaction into database atomicity (Orders + Items + Tickets) [GORM]
func (r *GormOrderRepository) Create(order *domain.Order) error {
	// Start a database transaction context blocks [GORM]
	return r.db.Transaction(func(tx *gorm.DB) error {
		
		if err := tx.Create(order).Error; err != nil {
			return err // Triggers automatic Rollback [GORM]
		}

		var ticketsToCreate []domain.Ticket
		for _, item := range order.OrderItems {
			for i := 0; i < item.Quantity; i++ {
				code, err := generateSecureTicketCode(12)
				if err != nil {
					return fmt.Errorf("failed to generate ticket code asset: %w", err)
				}

				ticket := domain.Ticket{
					OrderItemID: item.ID, // References the auto-incremented Item ID [GORM]
					TicketCode:  code,
					Status:      domain.TicketUnused,
					SyncStatus:  domain.SyncPending,
					CreatedAt:   time.Now(),
				}
				ticketsToCreate = append(ticketsToCreate, ticket)
			}
		}

		// Create Ticket
		if len(ticketsToCreate) > 0 {
			if err := tx.Create(&ticketsToCreate).Error; err != nil {
				return err // Triggers automatic Rollback if constraint fails (e.g. Code duplicate) [GORM]
			}
		}

		return nil // Commits the transaction securely [GORM]
	})
}

func handleOrderError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.NewNotFound("requested database records not found")
	}
	return apperror.NewInternalServerError("order transaction processing aborted: " + err.Error())
}

// Helper function to generate high-entropy secure alphanumeric ticket codes
func generateSecureTicketCode(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
