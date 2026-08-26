package approval

import (
	"repairdesk.local/internal/model"
	"repairdesk.local/internal/store"
	"testing"
)

func TestApprovalDecisions(t *testing.T) {
	db, _ := store.Open(t.TempDir() + "/a.db")
	defer db.Close()
	_ = db.SaveRequest(model.SpareIssueRequest{ID: "r", Status: model.StatusPendingApproval})
	s := New(db)
	if _, e := s.Approve("r", "sup", "go"); e != nil {
		t.Fatal(e)
	}
	r, _ := db.GetRequest("r")
	if r.Status != model.StatusApproved {
		t.Fatal(r.Status)
	}
}
