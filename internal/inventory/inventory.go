package inventory

import (
	"fmt"
	"repairdesk.local/internal/model"
	"repairdesk.local/internal/store"
)

type Service struct{ db *store.Store }

func New(db *store.Store) *Service { return &Service{db: db} }
func (s *Service) AddPart(p model.SparePart) error {
	if p.ID == "" || p.Name == "" {
		return fmt.Errorf("part id and name required")
	}
	if p.Quantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}
	p.Active = true
	return s.db.SavePart(p)
}
func (s *Service) UpdateStock(id string, delta int) error {
	p, e := s.db.GetPart(id)
	if e != nil {
		return e
	}
	if p.Quantity+delta < 0 {
		return model.ErrInsufficientStock
	}
	p.Quantity += delta
	return s.db.SavePart(p)
}
func (s *Service) CheckAvailability(id string, qty int) (bool, int, error) {
	p, e := s.db.GetPart(id)
	if e != nil {
		return false, 0, e
	}
	if qty <= 0 {
		return false, p.Quantity, model.ErrInvalidQuantity
	}
	return p.Quantity >= qty, p.Quantity, nil
}
func (s *Service) Reserve(id string, qty int) error {
	ok, _, e := s.CheckAvailability(id, qty)
	if e != nil {
		return e
	}
	if !ok {
		return model.ErrInsufficientStock
	}
	return s.UpdateStock(id, -qty)
}
func (s *Service) Release(id string, qty int) error { return s.UpdateStock(id, qty) }
func (s *Service) Snapshot(id string) (model.InventorySnapshot, error) {
	p, e := s.db.GetPart(id)
	if e != nil {
		return model.InventorySnapshot{}, e
	}
	return model.InventorySnapshot{Part: p, Available: p.Quantity}, nil
}
func (s *Service) LowStock() ([]model.SparePart, error) {
	all, e := s.db.ListParts()
	if e != nil {
		return nil, e
	}
	out := []model.SparePart{}
	for _, p := range all {
		if p.Active && p.Quantity <= p.ReorderLevel {
			out = append(out, p)
		}
	}
	return out, nil
}
func (s *Service) ValidateRequest(r model.SpareIssueRequest) error {
	ok, _, e := s.CheckAvailability(r.PartID, r.Quantity)
	if e != nil {
		return e
	}
	if !ok {
		return model.ErrInsufficientStock
	}
	return nil
}
func (s *Service) ConsumeForIssue(r model.SpareIssueRequest, operator string) (model.InventoryTransaction, error) {
	if e := s.Reserve(r.PartID, r.Quantity); e != nil {
		return model.InventoryTransaction{}, e
	}
	tx := model.InventoryTransaction{ID: store.NewID("txn"), RequestID: r.ID, PartID: r.PartID, Kind: "issue", Quantity: r.Quantity, Operator: operator}
	if e := s.db.SaveTransaction(tx); e != nil {
		_ = s.Release(r.PartID, r.Quantity)
		return model.InventoryTransaction{}, e
	}
	return tx, nil
}
