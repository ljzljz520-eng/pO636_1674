package transport

import (
	"encoding/json"
	"net/http"
	"repairdesk.local/internal/model"
	"repairdesk.local/internal/service"
	"strings"
)

type Handler struct{ svc *service.Service }

func New(s *service.Service) *Handler { return &Handler{svc: s} }
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	if path == "health" {
		write(w, 200, map[string]string{"status": "ok"})
		return
	}
	if path == "requests" && r.Method == http.MethodPost {
		var in struct {
			EquipmentID, PartID, EngineerID, FaultDescription string
			Quantity                                          int
		}
		if json.NewDecoder(r.Body).Decode(&in) != nil {
			writeError(w, 400, "Please provide valid request data")
			return
		}
		out, e := h.svc.Submit(in.EquipmentID, in.PartID, in.EngineerID, in.FaultDescription, in.Quantity)
		if e != nil {
			writeError(w, 400, e.Error())
			return
		}
		write(w, 201, out)
		return
	}
	if strings.HasPrefix(path, "requests/") {
		id := strings.TrimPrefix(path, "requests/")
		if r.Method == http.MethodGet {
			out, e := h.svc.Get(id)
			if e != nil {
				writeError(w, 404, "Request was not found")
				return
			}
			write(w, 200, out)
			return
		}
	}
	writeError(w, 404, "That repair workflow endpoint was not found")
}
func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, status int, msg string) {
	write(w, status, map[string]string{"error": msg})
}
func ErrorMessage(e error) string {
	if e == nil {
		return ""
	}
	switch e {
	case model.ErrInsufficientStock:
		return "Not enough stock is available for this request."
	case model.ErrInvalidTransition:
		return "This request is not ready for that action."
	}
	return e.Error()
}
