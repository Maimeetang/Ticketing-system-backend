package domain

type UserRole string

const (
	RoleAdmin   UserRole = "ADMIN"
	RoleManager UserRole = "MANAGER"
	RoleCashier UserRole = "CASHIER"
	RoleScanner UserRole = "SCANNER"
)

type ShiftStatus string

const (
	ShiftOpen   ShiftStatus = "OPEN"
	ShiftClosed ShiftStatus = "CLOSED"
)

type PaymentMethod string

const (
	Cash PaymentMethod = "CASH"
	Epay PaymentMethod = "EPAY"
)

type OrderStatus string

const (
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type TicketStatus string

const (
	TicketActive    TicketStatus = "ACTIVE"
	TicketUsed      TicketStatus = "USED"
	TicketCancelled TicketStatus = "CANCELLED"
)