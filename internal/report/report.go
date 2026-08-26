package report

import (
	"repairdesk.local/internal/model"
	"repairdesk.local/internal/store"
	"sort"
)

func Dashboard(db *store.Store) (model.Dashboard, error) {
	rs, e := db.ListRequests()
	if e != nil {
		return model.Dashboard{}, e
	}
	d := model.Dashboard{}
	for _, r := range rs {
		switch r.Status {
		case model.StatusSubmitted:
			d.Submitted++
		case model.StatusPendingApproval:
			d.PendingApproval++
		case model.StatusApproved:
			d.Approved++
		case model.StatusIssued:
			d.Issued++
		case model.StatusRejected:
			d.Rejected++
		}
		if r.HighValue {
			d.HighValue++
		}
	}
	parts, e := db.ListParts()
	if e != nil {
		return d, e
	}
	for _, p := range parts {
		if p.Quantity <= p.ReorderLevel {
			d.LowStock = append(d.LowStock, p)
		}
	}
	return d, nil
}
func SortRequests(rs []model.SpareIssueRequest) []model.SpareIssueRequest {
	out := append([]model.SpareIssueRequest{}, rs...)
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out
}
func StatusCounts(rs []model.SpareIssueRequest) map[string]int {
	m := map[string]int{}
	for _, r := range rs {
		m[r.Status]++
	}
	return m
}
func HighValue(rs []model.SpareIssueRequest) []model.SpareIssueRequest {
	out := []model.SpareIssueRequest{}
	for _, r := range rs {
		if r.HighValue {
			out = append(out, r)
		}
	}
	return out
}
func ExportLine(r model.SpareIssueRequest) string {
	return r.ID + "," + r.EquipmentID + "," + r.PartID + "," + r.Status
}
func TotalCost(rs []model.SpareIssueRequest, parts map[string]model.SparePart) float64 {
	var total float64
	for _, r := range rs {
		if p, ok := parts[r.PartID]; ok {
			total += p.UnitCost * float64(r.Quantity)
		}
	}
	return total
}
