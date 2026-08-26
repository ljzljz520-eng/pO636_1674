package store

import (
	"repairdesk.local/internal/model"
	"testing"
)

func TestStoreRoundTrip(t *testing.T) {
	db, e := Open(t.TempDir() + "/s.db")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	_ = db.SaveEquipment(model.Equipment{ID: "eq", Name: "Pump", Active: true})
	got, e := db.GetEquipment("eq")
	if e != nil || got.Name != "Pump" {
		t.Fatal(got, e)
	}
}
