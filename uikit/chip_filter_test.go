package uikit

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestModernChipFilter_Basic(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	options := []string{"Todos los años", "2026", "2025", "2024"}
	var lastSelected string

	chip := NewChipFilter("Año:", options, "Todos los años", func(sel string) {
		lastSelected = sel
	})

	if chip.Selected() != "Todos los años" {
		t.Errorf("Esperaba 'Todos los años', obtuve %s", chip.Selected())
	}
	if !chip.IsAtDefault() {
		t.Errorf("Esperaba que IsAtDefault() fuera true")
	}

	// Cambiar selección
	chip.SetSelected("2026")
	if chip.Selected() != "2026" {
		t.Errorf("Esperaba '2026', obtuve %s", chip.Selected())
	}
	if lastSelected != "2026" {
		t.Errorf("Callback no recibió '2026', recibió: %s", lastSelected)
	}
	if chip.IsAtDefault() {
		t.Errorf("Esperaba que IsAtDefault() fuera false tras cambio")
	}

	// Restablecer
	chip.Reset()
	if chip.Selected() != "Todos los años" {
		t.Errorf("Esperaba 'Todos los años' tras Reset, obtuve %s", chip.Selected())
	}
	if !chip.IsAtDefault() {
		t.Errorf("Esperaba IsAtDefault() == true tras Reset")
	}
}

func TestModernChipFilter_SetOptions(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	chip := NewChipFilter("Año:", []string{"Todos"}, "Todos", nil)
	chip.SetOptions([]string{"Todos", "2026", "2025"})

	if len(chip.options) != 3 {
		t.Errorf("Esperaba 3 opciones, obtuve %d", len(chip.options))
	}
}

func TestModernChipFilter_TapClear(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	chip := NewChipFilter("Año:", []string{"Todos", "2026"}, "Todos", nil)
	chip.Resize(fyne.NewSize(120, 28))
	chip.SetSelected("2026")

	if chip.IsAtDefault() {
		t.Fatalf("Debería no estar en default")
	}

	// Clic en la esquina derecha donde está el icono '✕'
	chip.Tapped(&fyne.PointEvent{
		Position: fyne.NewPos(110, 14),
	})

	if !chip.IsAtDefault() {
		t.Errorf("Hacer clic en la '✕' debería haber restablecido el filtro a 'Todos'")
	}
}
