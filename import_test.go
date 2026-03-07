package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
)

func TestImportJson(t *testing.T) {
	data := []byte(`{
		"id": "INV-001",
		"from": "Acme Corp",
		"to": "Client LLC",
		"items": ["Widget", "Gadget"],
		"quantities": [3, 5],
		"rates": [10.50, 20.00],
		"tax": 0.1
	}`)

	inv := &Invoice{}
	err := importJson(data, inv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if inv.Id != "INV-001" {
		t.Errorf("Id = %q, want %q", inv.Id, "INV-001")
	}
	if inv.From != "Acme Corp" {
		t.Errorf("From = %q, want %q", inv.From, "Acme Corp")
	}
	if inv.To != "Client LLC" {
		t.Errorf("To = %q, want %q", inv.To, "Client LLC")
	}
	if len(inv.Items) != 2 || inv.Items[0] != "Widget" {
		t.Errorf("Items = %v, want [Widget Gadget]", inv.Items)
	}
	if len(inv.Rates) != 2 || inv.Rates[0] != 10.50 {
		t.Errorf("Rates = %v, want [10.5 20]", inv.Rates)
	}
	if len(inv.Quantities) != 2 || inv.Quantities[0] != 3 {
		t.Errorf("Quantities = %v, want [3 5]", inv.Quantities)
	}
	if inv.Tax != 0.1 {
		t.Errorf("Tax = %f, want 0.1", inv.Tax)
	}
}

func TestImportJsonInvalid(t *testing.T) {
	err := importJson([]byte(`{not valid json`), &Invoice{})
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestImportYaml(t *testing.T) {
	data := []byte(`
id: "INV-002"
from: "Yaml Corp"
to: "Yaml Client"
items:
  - "Service A"
  - "Service B"
quantities:
  - 1
  - 2
rates:
  - 100.0
  - 200.0
tax: 0.05
`)

	inv := &Invoice{}
	err := importYaml(data, inv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if inv.Id != "INV-002" {
		t.Errorf("Id = %q, want %q", inv.Id, "INV-002")
	}
	if inv.From != "Yaml Corp" {
		t.Errorf("From = %q, want %q", inv.From, "Yaml Corp")
	}
	if len(inv.Items) != 2 || inv.Items[1] != "Service B" {
		t.Errorf("Items = %v, want [Service A Service B]", inv.Items)
	}
	if inv.Tax != 0.05 {
		t.Errorf("Tax = %f, want 0.05", inv.Tax)
	}
}

func TestImportYamlInvalid(t *testing.T) {
	data := []byte(":\n\t- :\n\t\t[invalid")
	err := importYaml(data, &Invoice{})
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestImportData(t *testing.T) {
	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "test.json")
	content := []byte(`{
		"id": "INV-100",
		"from": "File Corp",
		"to": "File Client",
		"items": ["Consulting"],
		"quantities": [10],
		"rates": [150.0]
	}`)
	if err := os.WriteFile(jsonPath, content, 0644); err != nil {
		t.Fatal(err)
	}

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("from", "", "")
	_ = flags.Set("from", "CLI Override Corp")

	inv := &Invoice{}
	err := importData(jsonPath, inv, flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if inv.Id != "INV-100" {
		t.Errorf("Id = %q, want %q", inv.Id, "INV-100")
	}
	// CLI flag should override the file value
	if inv.From != "CLI Override Corp" {
		t.Errorf("From = %q, want %q (CLI override)", inv.From, "CLI Override Corp")
	}
	if len(inv.Items) != 1 || inv.Items[0] != "Consulting" {
		t.Errorf("Items = %v, want [Consulting]", inv.Items)
	}
}
