package transport

import (
	"net/http/httptest"
	"repairdesk.local/internal/service"
	"repairdesk.local/internal/store"
	"strings"
	"testing"
)

func TestFriendlyErrorMessage(t *testing.T) {
	if ErrorMessage(fmtErr{}) != "demo" {
		t.Fatal("mapping")
	}
	db, _ := store.Open(t.TempDir() + "/h.db")
	defer db.Close()
	h := New(service.New(db))
	req := httptest.NewRequest("POST", "/requests", strings.NewReader("bad"))
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)
	if res.Code != 400 {
		t.Fatal(res.Code)
	}
}

type fmtErr struct{}

func (fmtErr) Error() string { return "demo" }
