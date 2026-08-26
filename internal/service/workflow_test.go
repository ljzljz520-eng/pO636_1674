package service

import (
	"repairdesk.local/internal/store"
	"testing"
)

func TestWorkflowOne(t *testing.T) {
	db, _ := store.Open(t.TempDir() + "/x.db")
	defer db.Close()
	s := New(db)
	if e := s.Seed(); e != nil {
		t.Fatal(e)
	}
	r, e := s.Submit("EQ-1", "SP-FUSE", "eng", "fuse failure", 2)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.ConfirmInventory(r.ID); e != nil {
		t.Fatal(e)
	}
	if _, e = s.RouteForApproval(r.ID); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Approve(r.ID, "sup", "ok"); e != nil {
		t.Fatal(e)
	}
	out, e := s.Issue(r.ID, "stock")
	if e != nil {
		t.Fatal(e)
	}
	if out.Request.Status != "issued" {
		t.Fatal("not issued")
	}
}
func TestWorkflowTwo(t *testing.T) {
	db, _ := store.Open(t.TempDir() + "/x.db")
	defer db.Close()
	s := New(db)
	_ = s.Seed()
	r, e := s.Submit("EQ-2", "SP-1674", "eng", "controller failure", 3)
	if e != nil {
		t.Fatal(e)
	}
	if !r.HighValue {
		t.Fatal("expected high value")
	}
	_, _ = s.ConfirmInventory(r.ID)
	_, _ = s.RouteForApproval(r.ID)
	q, e := s.appr.Queue()
	if e != nil || q.HighValue != 1 {
		t.Fatalf("queue %+v %v", q, e)
	}
}
func TestWorkflowThree(t *testing.T) {
	db, _ := store.Open(t.TempDir() + "/x.db")
	defer db.Close()
	s := New(db)
	_ = s.Seed()
	r, _ := s.Submit("EQ-3", "SP-FUSE", "eng", "fuse", 1)
	_, _ = s.ConfirmInventory(r.ID)
	_, _ = s.RouteForApproval(r.ID)
	_, e := s.Reject(r.ID, "sup", "unsafe")
	if e != nil {
		t.Fatal(e)
	}
	got, _ := s.Get(r.ID)
	if got.Status != "rejected" {
		t.Fatal(got.Status)
	}
}
