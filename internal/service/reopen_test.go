package service

import (
	"repairdesk.local/internal/model"
	"repairdesk.local/internal/store"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/persist.db"
	db, _ := store.Open(path)
	_ = db.SavePart(model.SparePart{ID: "persist-part", Name: "Bearing", Quantity: 7, Active: true})
	_ = db.Close()
	db2, e := store.Open(path)
	if e != nil {
		t.Fatal(e)
	}
	defer db2.Close()
	p, e := db2.GetPart("persist-part")
	if e != nil || p.Quantity != 7 {
		t.Fatalf("%+v %v", p, e)
	}
}
