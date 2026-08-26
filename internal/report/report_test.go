package report

import (
	"repairdesk.local/internal/model"
	"repairdesk.local/internal/store"
	"testing"
)

func TestDashboardCounts(t *testing.T) {
	db, _ := store.Open(t.TempDir() + "/r.db")
	defer db.Close()
	_ = db.SaveRequest(model.SpareIssueRequest{ID: "1", Status: model.StatusIssued, HighValue: true})
	d, e := Dashboard(db)
	if e != nil || d.Issued != 1 || d.HighValue != 1 {
		t.Fatalf("%+v %v", d, e)
	}
}
