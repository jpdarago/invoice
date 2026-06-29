package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpandNote(t *testing.T) {
	meta := map[string]string{"ADDRESS": "bc1qexample", "NAME": "Acme"}
	cases := []struct {
		name, in, want string
	}{
		{"braces", "pay: ${ADDRESS}", "pay: bc1qexample"},
		{"bare", "pay: $ADDRESS", "pay: bc1qexample"},
		{"mixed", "$NAME at ${ADDRESS}", "Acme at bc1qexample"},
		{"unknown bare left verbatim", "cost $50 now", "cost $50 now"},
		{"unknown key left verbatim", "ref ${MISSING}", "ref ${MISSING}"},
		{"no placeholders", "plain text", "plain text"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := expandNote(c.in, meta); got != c.want {
				t.Errorf("expandNote(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

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

func extractPDFText(t *testing.T, pdfPath string) string {
	t.Helper()
	out, err := exec.Command("pdftotext", pdfPath, "-").Output()
	if err != nil {
		t.Fatalf("pdftotext failed: %v", err)
	}
	return string(out)
}

func assertContains(t *testing.T, text, substr string) {
	t.Helper()
	if !strings.Contains(text, substr) {
		t.Errorf("expected text to contain %q, but it did not\ntext:\n%s", substr, text)
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

func TestInvoiceContent(t *testing.T) {
	if _, err := exec.LookPath("pdftotext"); err != nil {
		t.Skip("pdftotext not found, skipping content test")
	}

	dir := t.TempDir()
	outPath := filepath.Join(dir, "content-test.pdf")

	origFile := file
	origOutput := output
	t.Cleanup(func() {
		file = origFile
		output = origOutput
	})

	file = Invoice{
		Id:         "INV-2024-TEST",
		Title:      "INVOICE",
		From:       "Acme Widgets Inc.",
		To:         "Global Buyers LLC",
		Date:       "Jan 15, 2024",
		Due:        "Jan 29, 2024",
		Items:      []string{"Consulting", "Development"},
		Quantities: []int{5, 10},
		Rates:      []float64{150.00, 200.00},
		Tax:        0.1,
		Discount:   0.05,
		Currency:   "USD",
		Note:       "Payment via wire transfer",
	}
	output = outPath

	err := generateCmd.RunE(generateCmd, []string{})
	if err != nil {
		t.Fatalf("generateCmd.RunE failed: %v", err)
	}

	text := extractPDFText(t, outPath)

	for _, s := range []string{
		// Identity
		"INV-2024-TEST",
		"INVOICE",
		// Parties
		"Acme Widgets Inc.",
		"Global Buyers LLC",
		// Dates
		"Jan 15, 2024",
		"Jan 29, 2024",
		// Items
		"Consulting",
		"Development",
		// Quantities
		"5",
		"10",
		// Rates
		"$150.00",
		"$200.00",
		// Line amounts
		"$750.00",
		"$2000.00",
		// Totals
		"$2750.00",
		"$275.00",
		"$137.50",
		"$2887.50",
		// Labels
		"Subtotal",
		"Tax",
		"Discount",
		"Total",
		// Note
		"Payment via wire transfer",
		// Header labels
		"BILL TO",
		"ITEM",
		"QTY",
		"RATE",
		"AMOUNT",
		"NOTES",
	} {
		assertContains(t, text, s)
	}
}
