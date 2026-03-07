package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultInvoice(t *testing.T) {
	inv := DefaultInvoice()

	if inv.Title != "INVOICE" {
		t.Errorf("Title = %q, want %q", inv.Title, "INVOICE")
	}
	if inv.Currency != "USD" {
		t.Errorf("Currency = %q, want %q", inv.Currency, "USD")
	}
	if len(inv.Items) != 1 {
		t.Errorf("len(Items) = %d, want 1", len(inv.Items))
	}
	if len(inv.Rates) != 1 {
		t.Errorf("len(Rates) = %d, want 1", len(inv.Rates))
	}
	if len(inv.Quantities) != 1 {
		t.Errorf("len(Quantities) = %d, want 1", len(inv.Quantities))
	}
	if inv.Date == "" {
		t.Error("Date should not be empty")
	}
	if inv.Due == "" {
		t.Error("Due should not be empty")
	}
	if inv.Id == "" {
		t.Error("Id should not be empty")
	}
}

func TestGenerateInvoice(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "test-invoice.pdf")

	// Save and restore global state
	origFile := file
	origOutput := output
	t.Cleanup(func() {
		file = origFile
		output = origOutput
	})

	file = DefaultInvoice()
	file.From = "Test Sender"
	file.To = "Test Receiver"
	file.Items = []string{"Test Item"}
	file.Quantities = []int{1}
	file.Rates = []float64{100.0}
	output = outPath

	err := generateCmd.RunE(generateCmd, []string{})
	if err != nil {
		t.Fatalf("generateCmd.RunE failed: %v", err)
	}

	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("output file not found: %v", err)
	}
	if info.Size() == 0 {
		t.Error("output PDF file is empty")
	}
}
