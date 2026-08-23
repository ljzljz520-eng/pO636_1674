package model

import "testing"

func TestRequestValidation(t *testing.T) {
	if (SpareIssueRequest{EquipmentID: "e", PartID: "p", EngineerID: "g", FaultDescription: "f", Quantity: 1}).Validate() != nil {
		t.Fatal("valid")
	}
	if (SpareIssueRequest{Quantity: 0}).Validate() == nil {
		t.Fatal("invalid")
	}
}
