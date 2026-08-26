package store

import (
	"repairdesk.local/internal/model"
	"sort"
)

func (s *Store) RequestsByStatus(status string) ([]model.SpareIssueRequest, error) {
	all, e := s.ListRequests()
	if e != nil {
		return nil, e
	}
	out := []model.SpareIssueRequest{}
	for _, r := range all {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out, nil
}
func (s *Store) TransactionsForRequest(id string) ([]model.InventoryTransaction, error) {
	all, e := s.ListTransactions()
	if e != nil {
		return nil, e
	}
	out := []model.InventoryTransaction{}
	for _, v := range all {
		if v.RequestID == id {
			out = append(out, v)
		}
	}
	return out, nil
}
func (s *Store) SortPartsByQuantity() ([]model.SparePart, error) {
	all, e := s.ListParts()
	if e != nil {
		return nil, e
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Quantity < all[j].Quantity })
	return all, nil
}
func (s *Store) MarkNotificationRead(id string) error {
	n, e := get[model.Notification](s, buckets[5], id)
	if e != nil {
		return e
	}
	n.Read = true
	return put(s, buckets[5], id, n)
}
func (s *Store) RequestCount() int {
	all, e := s.ListRequests()
	if e != nil {
		return 0
	}
	return len(all)
}
