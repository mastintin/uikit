package uikit

import (
	"testing"
	"time"
)

func TestValToStr(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"int", 42, "42"},
		{"float", 3.14, "3.14"},
		{"string", "hola mundo", "hola mundo"},
		{"bool", true, "true"},
		{"nil", nil, "<nil>"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := ValToStr(tc.input)
			if res != tc.expected {
				t.Errorf("ValToStr(%v) = %q; want %q", tc.input, res, tc.expected)
			}
		})
	}
}

func TestDateToStr(t *testing.T) {
	d := time.Date(2026, 10, 16, 0, 0, 0, 0, time.UTC)

	eu := DateToStrEu(d)
	if eu != "16/10/2026" {
		t.Errorf("DateToStrEu = %q; want 16/10/2026", eu)
	}

	us := DateToStrUs(d)
	if us != "10/16/2026" {
		t.Errorf("DateToStrUs = %q; want 10/16/2026", us)
	}

	baseEu := dateToStrBase(d, true)
	if baseEu != "16/10/2026" {
		t.Errorf("dateToStrBase(true) = %q; want 16/10/2026", baseEu)
	}

	baseUs := dateToStrBase(d, false)
	if baseUs != "10/16/2026" {
		t.Errorf("dateToStrBase(false) = %q; want 10/16/2026", baseUs)
	}
}

func TestConCatStr(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{"empty", []string{}, ""},
		{"single", []string{"hola"}, "hola"},
		{"multiple", []string{"Go", "Jubila", " ", "2026"}, "GoJubila 2026"},
		{"symbols", []string{"1.234", ",", "56", " €"}, "1.234,56 €"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := ConCatStr(tc.input...)
			if res != tc.expected {
				t.Errorf("ConCatStr(%v) = %q; want %q", tc.input, res, tc.expected)
			}
		})
	}
}

func TestAAddStr(t *testing.T) {
	var lines []string
	AAddStr(&lines, "linea 1")
	AAddStr(&lines, "linea 2")
	AAddStr(&lines, "linea 3")

	if len(lines) != 3 {
		t.Fatalf("Expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "linea 1" || lines[1] != "linea 2" || lines[2] != "linea 3" {
		t.Errorf("Unexpected content: %v", lines)
	}
}

func TestCRLF(t *testing.T) {
	if CRLF != "\n" {
		t.Errorf("CRLF = %q; want %q", CRLF, "\n")
	}
}

