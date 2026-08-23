package model

type RequestFilter struct {
	Status, EngineerID, EquipmentID, PartID string
	HighValueOnly                           bool
}
type InventorySnapshot struct {
	Part      SparePart
	Reserved  int
	Available int
}
type Dashboard struct {
	Submitted, PendingApproval, Approved, Issued, Rejected, HighValue int
	LowStock                                                          []SparePart
}
type IssueResult struct {
	Request     SpareIssueRequest
	Transaction InventoryTransaction
	Remaining   int
}
type ApprovalQueue struct {
	Requests  []SpareIssueRequest
	Total     int
	HighValue int
}
type Page struct{ Limit, Offset, Total int }

func (f RequestFilter) Matches(r SpareIssueRequest) bool {
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.EngineerID != "" && r.EngineerID != f.EngineerID {
		return false
	}
	if f.EquipmentID != "" && r.EquipmentID != f.EquipmentID {
		return false
	}
	if f.PartID != "" && r.PartID != f.PartID {
		return false
	}
	if f.HighValueOnly && !r.HighValue {
		return false
	}
	return true
}
