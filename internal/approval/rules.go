package approval

import (
	"repairdesk.local/internal/model"
	"repairdesk.local/internal/store"
)

func RequiresApproval(r model.SpareIssueRequest, cost float64) bool {
	return r.HighValue || cost >= 500
}
func AllowedSupervisor(id string) bool { return id != "" }
func BuildRecord(id, request, supervisor, decision, comment string) model.ApprovalRecord {
	return model.ApprovalRecord{ID: id, RequestID: request, SupervisorID: supervisor, Decision: decision, Comment: comment}
}
func PendingCount(db *store.Store) (int, error) { q, e := New(db).Queue(); return q.Total, e }
func IsDecisionValid(d string) bool             { return d == "approved" || d == "rejected" }
