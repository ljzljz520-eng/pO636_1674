package inventory

import (
	"repairdesk.local/internal/model"
	"repairdesk.local/internal/store"
)

func (s *Service) SeedDefaults() error {
	parts := []model.SparePart{{ID: "SP-1674", Name: "Control Relay", Description: "24V relay", UnitCost: 180, ReorderLevel: 5, Quantity: 3, Active: true}, {ID: "SP-FUSE", Name: "Fuse", Description: "Fast fuse", UnitCost: 12, ReorderLevel: 10, Quantity: 40, Active: true}}
	for _, p := range parts {
		if _, e := s.db.GetPart(p.ID); e == model.ErrNotFound {
			if e = s.db.SavePart(p); e != nil {
				return e
			}
		}
	}
	return nil
}
func (s *Service) Reconcile() error {
	parts, e := s.db.ListParts()
	if e != nil {
		return e
	}
	for _, p := range parts {
		if p.Quantity < 0 {
			p.Quantity = 0
		}
		if !p.Active {
			p.Quantity = 0
		}
		if e = s.db.SavePart(p); e != nil {
			return e
		}
	}
	return nil
}
func StockLabel(p model.SparePart) string {
	if p.Quantity == 0 {
		return "out of stock"
	}
	if p.Quantity <= p.ReorderLevel {
		return "reorder soon"
	}
	return "available"
}
func IsHighValue(p model.SparePart, qty int) bool { return p.UnitCost*float64(qty) >= 500 }
func PartCost(p model.SparePart, qty int) float64 {
	if qty <= 0 {
		return 0
	}
	return p.UnitCost * float64(qty)
}
func ensurePart(db *store.Store, id string) (model.SparePart, error) { return db.GetPart(id) }
