package service

import (
	"repairdesk.local/internal/model"
	"repairdesk.local/internal/store"
)

func (s *Service) RecordAudit(req, actor, action, details string) error {
	return s.db.SaveAudit(model.AuditEvent{ID: store.NewID("audit"), RequestID: req, Actor: actor, Action: action, Details: details})
}
func (s *Service) Notify(req, recipient, message, severity string) error {
	return s.db.SaveNotification(model.Notification{ID: store.NewID("note"), RequestID: req, Recipient: recipient, Message: message, Severity: severity})
}
func (s *Service) Notifications() ([]model.Notification, error) { return s.db.ListNotifications() }
func (s *Service) Audit(req string) ([]model.AuditEvent, error) {
	all, e := s.db.ListAudit()
	if e != nil {
		return nil, e
	}
	out := []model.AuditEvent{}
	for _, a := range all {
		if a.RequestID == req {
			out = append(out, a)
		}
	}
	return out, nil
}
func (s *Service) CompleteWorkflow(id, supervisor, operator string) (model.IssueResult, error) {
	if _, e := s.ConfirmInventory(id); e != nil {
		return model.IssueResult{}, e
	}
	if _, e := s.RouteForApproval(id); e != nil {
		return model.IssueResult{}, e
	}
	if _, e := s.Approve(id, supervisor, "routine approval"); e != nil {
		return model.IssueResult{}, e
	}
	return s.Issue(id, operator)
}
func (s *Service) Reopen(path string) (*Service, error) {
	db, e := store.Open(path)
	if e != nil {
		return nil, e
	}
	return New(db), nil
}
