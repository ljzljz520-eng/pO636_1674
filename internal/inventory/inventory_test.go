package inventory_test

import (
	"repairdesk.local/internal/inventory"
	"repairdesk.local/internal/model"
	"repairdesk.local/internal/service"
	"repairdesk.local/internal/store"
	"testing"
)

func TestSpareIssueRejectsInsufficientStock(t *testing.T) {
	db, _ := store.Open(t.TempDir() + "/x.db")
	defer db.Close()
	inv := inventory.New(db)
	_ = inv.AddPart(model.SparePart{ID: "SP-1674", Name: "Relay", Quantity: 3, Active: true})
	svc := service.New(db)
	r, e := svc.Submit("EQ-1674", "SP-1674", "eng-1", "relay failed", 8)
	if e != nil {
		t.Fatal(e)
	}
	r.Status = model.StatusApproved
	_ = db.SaveRequest(r)
	_, e = svc.Issue(r.ID, "stock-1")
	if e == nil {
		t.Fatal("expected inventory rejection")
	}
}
