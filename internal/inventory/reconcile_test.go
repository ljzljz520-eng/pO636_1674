package inventory

import (
	"repairdesk.local/internal/model"
	"repairdesk.local/internal/store"
	"testing"
)

func TestReconcileInactive(t *testing.T) {
	db, _ := store.Open(t.TempDir() + "/q.db")
	defer db.Close()
	s := New(db)
	_ = db.SavePart(model.SparePart{ID: "x", Name: "x", Quantity: 2, Active: false})
	if e := s.Reconcile(); e != nil {
		t.Fatal(e)
	}
	p, _ := db.GetPart("x")
	if p.Quantity != 0 {
		t.Fatal(p.Quantity)
	}
}
