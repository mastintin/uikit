package uikit

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestModernMetricCard(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	card := NewMetricCard(
		"TOTAL NETO",
		"lucide-wallet",
		"#818cf8",
		"#1e1b4b",
		ColorSuccess,
	)

	if card == nil {
		t.Fatal("NewMetricCard devolvió nil")
	}

	if card.lblTitle.Text() != "TOTAL NETO" {
		t.Errorf("Esperado título 'TOTAL NETO', obtenido '%s'", card.lblTitle.Text())
	}

	if card.lblValue.Text() != "--" {
		t.Errorf("Esperado valor inicial '--', obtenido '%s'", card.lblValue.Text())
	}

	// Probar SetValue
	card.SetValue("3.175,00 €")
	if card.lblValue.Text() != "3.175,00 €" {
		t.Errorf("SetValue falló, obtenido '%s'", card.lblValue.Text())
	}

	// Probar SetFooters
	red := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	card.SetFooters("14 pagas", "+3.2%", red)
	if card.lblFootL.Text() != "14 pagas" {
		t.Errorf("SetFooters left falló, obtenido '%s'", card.lblFootL.Text())
	}
	if card.lblFootR.Text() != "+3.2%" {
		t.Errorf("SetFooters right falló, obtenido '%s'", card.lblFootR.Text())
	}
	if card.lblFootR.textColor != red {
		t.Errorf("SetFooters rightColor no se aplicó")
	}

	// Probar SetTitle
	card.SetTitle("NUEVO TÍTULO")
	if card.lblTitle.Text() != "NUEVO TÍTULO" {
		t.Errorf("SetTitle falló")
	}

	// Probar SetValueColor
	card.SetValueColor(ColorWarning)
	if card.valueColor != ColorWarning {
		t.Errorf("SetValueColor falló")
	}
}
