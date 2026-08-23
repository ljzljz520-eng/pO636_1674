package service

import (
	"fmt"
	"repairdesk.local/internal/model"
)

func ValidateEngineer(id string) error {
	if len(id) < 2 {
		return fmt.Errorf("engineer id is too short")
	}
	return nil
}
func ValidateEquipment(id string) error {
	if len(id) < 2 {
		return fmt.Errorf("equipment id is too short")
	}
	return nil
}
func CanEdit(r model.SpareIssueRequest, actor string) bool {
	return r.EngineerID == actor && r.Status == model.StatusSubmitted
}
func CanCancel(r model.SpareIssueRequest, actor string) bool { return actor != "" && !r.IsTerminal() }
func NeedsEscalation(r model.SpareIssueRequest) bool         { return r.HighValue || r.Quantity >= 10 }
func StatusLabel(status string) string {
	switch status {
	case model.StatusSubmitted:
		return "Submitted"
	case model.StatusInventoryConfirmed:
		return "Inventory confirmed"
	case model.StatusPendingApproval:
		return "Pending supervisor approval"
	case model.StatusApproved:
		return "Approved"
	case model.StatusRejected:
		return "Rejected"
	case model.StatusIssued:
		return "Issued"
	case model.StatusCancelled:
		return "Cancelled"
	}
	return "Unknown"
}
func AllowedTransition(from, to string) bool {
	switch from {
	case model.StatusSubmitted:
		return to == model.StatusInventoryConfirmed || to == model.StatusCancelled
	case model.StatusInventoryConfirmed:
		return to == model.StatusPendingApproval || to == model.StatusCancelled
	case model.StatusPendingApproval:
		return to == model.StatusApproved || to == model.StatusRejected
	case model.StatusApproved:
		return to == model.StatusIssued || to == model.StatusCancelled
	}
	return false
}
func WorkflowComplete(r model.SpareIssueRequest) bool {
	return r.Status == model.StatusIssued || r.Status == model.StatusRejected
}
