package uikit

import (
	"testing"
)

func TestTypographyStylesAndApply(t *testing.T) {
	// Test Hero
	h := NewHeroText("1.234,56 €", ColorSuccess)
	if h.TextSize != 26 {
		t.Errorf("Expected Hero size 26, got %f", h.TextSize)
	}
	if !h.TextStyle.Bold {
		t.Errorf("Expected Hero to be bold")
	}

	// Test Title1
	t1 := NewTitle1Text("Simulador")
	if t1.TextSize != 22 {
		t.Errorf("Expected Title1 size 22, got %f", t1.TextSize)
	}

	// Test Title2
	t2 := NewTitle2Text("Subtítulo", ColorWhite)
	if t2.TextSize != 18 {
		t.Errorf("Expected Title2 size 18, got %f", t2.TextSize)
	}

	// Test Section
	sec := NewSectionText("Sección 1")
	if sec.TextSize != 14.5 {
		t.Errorf("Expected Section size 14.5, got %f", sec.TextSize)
	}

	// Test Value (regular and bold)
	vReg := NewValueText("100 €", false)
	if vReg.TextSize != 13 || vReg.TextStyle.Bold {
		t.Errorf("Expected Value regular size 13 not bold")
	}
	vBold := NewValueText("100 €", true)
	if vBold.TextSize != 13 || !vBold.TextStyle.Bold {
		t.Errorf("Expected Value bold size 13 bold")
	}

	// Test Modifiers
	base := NewBodyText("Base text")
	if base.TextSize != 12.5 {
		t.Errorf("Expected Body size 12.5, got %f", base.TextSize)
	}
	ApplyHero(base)
	if base.TextSize != 26 {
		t.Errorf("Expected applied Hero size 26, got %f", base.TextSize)
	}

	ApplyTitle2(base)
	if base.TextSize != 18 {
		t.Errorf("Expected applied Title2 size 18, got %f", base.TextSize)
	}
}
