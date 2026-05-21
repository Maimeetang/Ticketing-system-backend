package domain

type SyncStatus string

const (
	SyncPending SyncStatus = "PENDING" // Data is stored locally but not yet pushed to cloud
	SyncSynced  SyncStatus = "SYNCED"  // Data successfully replicated to central cloud
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleManager UserRole = "manager"
	RoleCashier UserRole = "cashier"
	RoleScanner UserRole = "scanner"
)

type PaymentMethod string

const (
	Cash PaymentMethod = "cash"
	Epay PaymentMethod = "epay"
)

type TicketStatus string

const (
	TicketUnused    TicketStatus = "UNUSED"
	TicketUsed      TicketStatus = "USED"
	TicketCancelled TicketStatus = "CANCELLED"
)