package report

import (
	"repairdesk.local/internal/model"
	"strings"
)

func FilterByEquipment(rs []model.SpareIssueRequest, id string) []model.SpareIssueRequest {
	out := []model.SpareIssueRequest{}
	for _, r := range rs {
		if r.EquipmentID == id {
			out = append(out, r)
		}
	}
	return out
}
func FilterByEngineer(rs []model.SpareIssueRequest, id string) []model.SpareIssueRequest {
	out := []model.SpareIssueRequest{}
	for _, r := range rs {
		if r.EngineerID == id {
			out = append(out, r)
		}
	}
	return out
}
func SearchText(rs []model.SpareIssueRequest, q string) []model.SpareIssueRequest {
	q = strings.ToLower(q)
	out := []model.SpareIssueRequest{}
	for _, r := range rs {
		if strings.Contains(strings.ToLower(r.FaultDescription), q) || strings.Contains(strings.ToLower(r.PartID), q) {
			out = append(out, r)
		}
	}
	return out
}
func Paginate(rs []model.SpareIssueRequest, page model.Page) []model.SpareIssueRequest {
	start := page.Offset
	if start > len(rs) {
		return []model.SpareIssueRequest{}
	}
	end := start + page.Limit
	if page.Limit <= 0 || end > len(rs) {
		end = len(rs)
	}
	return rs[start:end]
}
func GroupByPart(rs []model.SpareIssueRequest) map[string][]model.SpareIssueRequest {
	m := map[string][]model.SpareIssueRequest{}
	for _, r := range rs {
		m[r.PartID] = append(m[r.PartID], r)
	}
	return m
}
func GroupByStatus(rs []model.SpareIssueRequest) map[string][]model.SpareIssueRequest {
	m := map[string][]model.SpareIssueRequest{}
	for _, r := range rs {
		m[r.Status] = append(m[r.Status], r)
	}
	return m
}
func IsAged(r model.SpareIssueRequest, threshold int) bool {
	return threshold > 0 && len(r.FaultDescription) >= threshold
}
