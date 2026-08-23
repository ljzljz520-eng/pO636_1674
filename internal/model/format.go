package model

import "fmt"

func (p SparePart) Value() float64     { return p.UnitCost * float64(p.Quantity) }
func (p SparePart) NeedsReorder() bool { return p.Active && p.Quantity <= p.ReorderLevel }
func (r SpareIssueRequest) Summary() string {
	return fmt.Sprintf("%s requests %d of %s for %s", r.EngineerID, r.Quantity, r.PartID, r.EquipmentID)
}
func (r SpareIssueRequest) CanApprove() bool { return r.Status == StatusPendingApproval }
func (r SpareIssueRequest) CanIssue() bool   { return r.Status == StatusApproved }
func (a ApprovalRecord) Accepted() bool      { return a.Decision == "approved" }
func (t InventoryTransaction) IsIssue() bool { return t.Kind == "issue" }
func NormalizeStatus(s string) string {
	switch s {
	case "submitted", "new":
		return StatusSubmitted
	case "approved":
		return StatusApproved
	case "issued", "completed":
		return StatusIssued
	case "rejected":
		return StatusRejected
	}
	return s
}
