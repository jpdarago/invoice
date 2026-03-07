package main

import "testing"

func TestCurrencySymbols(t *testing.T) {
	cases := []struct {
		code   string
		symbol string
	}{
		{"USD", "$"},
		{"EUR", "€"},
		{"GBP", "£"},
		{"JPY", "¥"},
	}

	for _, tc := range cases {
		got := currencySymbols[tc.code]
		if got != tc.symbol {
			t.Errorf("currencySymbols[%q] = %q, want %q", tc.code, got, tc.symbol)
		}
	}
}

func TestCurrencySymbolsUnknown(t *testing.T) {
	got := currencySymbols["XYZ"]
	if got != "" {
		t.Errorf("currencySymbols[\"XYZ\"] = %q, want empty string", got)
	}
}
