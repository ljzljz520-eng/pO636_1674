package model

import "time"

type SparePart struct {
	ID, Name, Description  string
	UnitCost               float64
	ReorderLevel, Quantity int
	Active                 bool
}
type Equipment struct {
	ID, Name, Location, Owner string
	Active                    bool
}
type SpareIssueRequest struct {
	ID, EquipmentID, PartID, EngineerID, FaultDescription, Status string
	Quantity                                                      int
	HighValue                                                     bool
	CreatedAt                                                     time.Time
}
type ApprovalRecord struct {
	ID, RequestID, SupervisorID, Decision, Comment string
	DecidedAt                                      time.Time
}
type InventoryTransaction struct {
	ID, RequestID, PartID, Kind string
	Quantity                    int
	Operator                    string
	CreatedAt                   time.Time
}
type Notification struct {
	ID, RequestID, Recipient, Message, Severity string
	CreatedAt                                   time.Time
	Read                                        bool
}
type AuditEvent struct {
	ID, RequestID, Actor, Action, Details string
	CreatedAt                             time.Time
}

const (
	StatusSubmitted          = "submitted"
	StatusInventoryConfirmed = "inventory_confirmed"
	StatusPendingApproval    = "pending_approval"
	StatusApproved           = "approved"
	StatusRejected           = "rejected"
	StatusIssued             = "issued"
	StatusCancelled          = "cancelled"
)

func (r SpareIssueRequest) IsTerminal() bool {
	return r.Status == StatusIssued || r.Status == StatusRejected || r.Status == StatusCancelled
}
func (r SpareIssueRequest) Validate() error {
	if r.EquipmentID == "" || r.PartID == "" || r.EngineerID == "" {
		return ErrInvalidRequest
	}
	if r.Quantity <= 0 {
		return ErrInvalidQuantity
	}
	if r.FaultDescription == "" {
		return ErrMissingFault
	}
	return nil
}

type DomainError string

func (e DomainError) Error() string { return string(e) }

var (
	ErrInvalidRequest    = DomainError("equipment, part, and engineer are required")
	ErrInvalidQuantity   = DomainError("quantity must be positive")
	ErrMissingFault      = DomainError("fault description is required")
	ErrInsufficientStock = DomainError("insufficient stock")
	ErrNotFound          = DomainError("record not found")
	ErrInvalidTransition = DomainError("invalid status transition")
)
