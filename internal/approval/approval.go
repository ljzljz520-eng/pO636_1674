package approval

import (
	"fmt"
	"repairdesk.local/internal/model"
	"repairdesk.local/internal/store"
)

type Service struct{ db *store.Store }

func New(db *store.Store) *Service { return &Service{db: db} }
func (s *Service) Queue() (model.ApprovalQueue, error) {
	rs, e := s.db.ListRequests()
	if e != nil {
		return model.ApprovalQueue{}, e
	}
	q := model.ApprovalQueue{}
	for _, r := range rs {
		if r.Status == model.StatusPendingApproval {
			q.Requests = append(q.Requests, r)
			q.Total++
			if r.HighValue {
				q.HighValue++
			}
		}
	}
	return q, nil
}
func (s *Service) Approve(id, supervisor, comment string) (model.ApprovalRecord, error) {
	r, e := s.db.GetRequest(id)
	if e != nil {
		return model.ApprovalRecord{}, e
	}
	if r.Status != model.StatusPendingApproval {
		return model.ApprovalRecord{}, model.ErrInvalidTransition
	}
	r.Status = model.StatusApproved
	if e = s.db.SaveRequest(r); e != nil {
		return model.ApprovalRecord{}, e
	}
	a := model.ApprovalRecord{ID: store.NewID("approval"), RequestID: id, SupervisorID: supervisor, Decision: "approved", Comment: comment}
	if e = s.db.SaveApproval(a); e != nil {
		return model.ApprovalRecord{}, e
	}
	return a, nil
}
func (s *Service) Reject(id, supervisor, comment string) (model.ApprovalRecord, error) {
	r, e := s.db.GetRequest(id)
	if e != nil {
		return model.ApprovalRecord{}, e
	}
	if r.Status != model.StatusPendingApproval {
		return model.ApprovalRecord{}, model.ErrInvalidTransition
	}
	r.Status = model.StatusRejected
	if e = s.db.SaveRequest(r); e != nil {
		return model.ApprovalRecord{}, e
	}
	a := model.ApprovalRecord{ID: store.NewID("approval"), RequestID: id, SupervisorID: supervisor, Decision: "rejected", Comment: comment}
	if e = s.db.SaveApproval(a); e != nil {
		return model.ApprovalRecord{}, e
	}
	return a, nil
}
func (s *Service) Cancel(id, actor string) error {
	r, e := s.db.GetRequest(id)
	if e != nil {
		return e
	}
	if r.IsTerminal() {
		return fmt.Errorf("request already complete")
	}
	r.Status = model.StatusCancelled
	return s.db.SaveRequest(r)
}
func (s *Service) History(id string) ([]model.ApprovalRecord, error) {
	all, e := s.db.ListApprovals()
	if e != nil {
		return nil, e
	}
	out := []model.ApprovalRecord{}
	for _, a := range all {
		if a.RequestID == id {
			out = append(out, a)
		}
	}
	return out, nil
}
