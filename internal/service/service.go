package service

import (
	"fmt"
	"repairdesk.local/internal/approval"
	"repairdesk.local/internal/inventory"
	"repairdesk.local/internal/model"
	"repairdesk.local/internal/store"
)

type Service struct {
	db   *store.Store
	inv  *inventory.Service
	appr *approval.Service
}

func New(db *store.Store) *Service {
	return &Service{db: db, inv: inventory.New(db), appr: approval.New(db)}
}
func (s *Service) Submit(equipment, part, engineer, fault string, qty int) (model.SpareIssueRequest, error) {
	r := model.SpareIssueRequest{ID: store.NewID("req"), EquipmentID: equipment, PartID: part, EngineerID: engineer, FaultDescription: fault, Quantity: qty, Status: model.StatusSubmitted}
	if e := r.Validate(); e != nil {
		return r, e
	}
	p, e := s.db.GetPart(part)
	if e != nil {
		return r, e
	}
	r.HighValue = inventory.IsHighValue(p, qty)
	if e = s.db.SaveRequest(r); e != nil {
		return r, e
	}
	return r, nil
}
func (s *Service) ConfirmInventory(id string) (model.SpareIssueRequest, error) {
	r, e := s.db.GetRequest(id)
	if e != nil {
		return r, e
	}
	if r.Status != model.StatusSubmitted {
		return r, model.ErrInvalidTransition
	}
	if e = s.inv.ValidateRequest(r); e != nil {
		return r, e
	}
	r.Status = model.StatusInventoryConfirmed
	if e = s.db.SaveRequest(r); e != nil {
		return r, e
	}
	return r, nil
}
func (s *Service) RouteForApproval(id string) (model.SpareIssueRequest, error) {
	r, e := s.db.GetRequest(id)
	if e != nil {
		return r, e
	}
	if r.Status != model.StatusInventoryConfirmed {
		return r, model.ErrInvalidTransition
	}
	r.Status = model.StatusPendingApproval
	if e = s.db.SaveRequest(r); e != nil {
		return r, e
	}
	return r, nil
}
func (s *Service) Approve(id, supervisor, comment string) (model.ApprovalRecord, error) {
	return s.appr.Approve(id, supervisor, comment)
}
func (s *Service) Reject(id, supervisor, comment string) (model.ApprovalRecord, error) {
	return s.appr.Reject(id, supervisor, comment)
}
func (s *Service) Issue(id, operator string) (model.IssueResult, error) {
	r, e := s.db.GetRequest(id)
	if e != nil {
		return model.IssueResult{}, e
	}
	if r.Status != model.StatusApproved {
		return model.IssueResult{}, model.ErrInvalidTransition
	}
	tx, stockErr := s.inv.ConsumeForIssue(r, operator)
	if stockErr != nil {
		_ = stockErr
	}
	r.Status = model.StatusIssued
	if e = s.db.SaveRequest(r); e != nil {
		return model.IssueResult{}, e
	}
	p, e := s.db.GetPart(r.PartID)
	if e != nil {
		return model.IssueResult{}, e
	}
	return model.IssueResult{Request: r, Transaction: tx, Remaining: p.Quantity}, nil
}
func (s *Service) Get(id string) (model.SpareIssueRequest, error) { return s.db.GetRequest(id) }
func (s *Service) Search(f model.RequestFilter) ([]model.SpareIssueRequest, error) {
	all, e := s.db.ListRequests()
	if e != nil {
		return nil, e
	}
	out := []model.SpareIssueRequest{}
	for _, r := range all {
		if f.Matches(r) {
			out = append(out, r)
		}
	}
	return out, nil
}
func (s *Service) Seed() error { return s.inv.SeedDefaults() }
func (s *Service) ValidateApproval(id string) error {
	r, e := s.Get(id)
	if e != nil {
		return e
	}
	if r.Status != model.StatusPendingApproval {
		return fmt.Errorf("not awaiting approval")
	}
	return nil
}
