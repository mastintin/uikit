package uikit

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestNewSlider(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	s := NewSlider(0, 100)
	if s == nil {
		t.Fatal("NewSlider no debe ser nil")
	}
	if s.Value() != 0 {
		t.Errorf("valor inicial esperado 0, obtuve %v", s.Value())
	}
	if s.min != 0 || s.max != 100 {
		t.Errorf("rango esperado [0, 100], obtuve [%v, %v]", s.min, s.max)
	}
}

func TestSliderSemanticConstructors(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	age := NewAgeSlider(60, 70)
	if age.suffix != " años" {
		t.Errorf("sufijo esperado ' años', obtuve '%v'", age.suffix)
	}
	if age.decimals != 0 {
		t.Errorf("decimales esperados 0, obtuve %v", age.decimals)
	}

	pct := NewPercentageSlider(0, 100)
	if pct.suffix != " %" {
		t.Errorf("sufijo esperado ' %%', obtuve '%v'", pct.suffix)
	}

	curr := NewCurrencySlider(100, 1000, 25)
	if curr.suffix != " €" {
		t.Errorf("sufijo esperado ' €', obtuve '%v'", curr.suffix)
	}
	if curr.decimals != 2 {
		t.Errorf("decimales esperados 2, obtuve %v", curr.decimals)
	}
}

func TestSliderFormatValue(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	s := NewSlider(0, 100).
		WithPrefix("Ref: ").
		WithSuffix(" unidades").
		WithPrecision(1)

	got := s.formatValue(42.56)
	expected := "Ref: 42.6 unidades"
	if got != expected {
		t.Errorf("esperaba '%v', obtuve '%v'", expected, got)
	}

	// Con formateador personalizado
	s.WithFormat(func(v float64) string {
		return fmt.Sprintf("Edad: %.0f años", v)
	})
	gotCustom := s.formatValue(65)
	if gotCustom != "Edad: 65 años" {
		t.Errorf("esperaba 'Edad: 65 años', obtuve '%v'", gotCustom)
	}
}

func TestSliderSetValueAndClamping(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	s := NewSlider(10, 50).WithStep(5)

	// Valor dentro del rango con paso
	s.SetValue(23) // más cercano a 25
	if s.Value() != 25 {
		t.Errorf("esperaba alineación a step 25, obtuve %v", s.Value())
	}

	// Valor menor que el mínimo
	s.SetValue(-10)
	if s.Value() != 10 {
		t.Errorf("esperaba clamp a mínimo 10, obtuve %v", s.Value())
	}

	// Valor mayor que el máximo
	s.SetValue(999)
	if s.Value() != 50 {
		t.Errorf("esperaba clamp a máximo 50, obtuve %v", s.Value())
	}
}

func TestSliderCallbacks(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	s := NewSlider(0, 10)
	changedVal := -1.0
	endedVal := -1.0

	s.OnChanged(func(v float64) {
		changedVal = v
	})
	s.OnChangeEnded(func(v float64) {
		endedVal = v
	})

	r := s.CreateRenderer().(*modernSliderRenderer)
	r.innerSlider.OnChanged(7.0)
	if changedVal != 7.0 {
		t.Errorf("OnChanged callback falló, valor %v", changedVal)
	}

	r.innerSlider.OnChangeEnded(8.0)
	if endedVal != 8.0 {
		t.Errorf("OnChangeEnded callback falló, valor %v", endedVal)
	}
}

func TestSliderRendererLayout(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	// 1. Layout Inline sin título
	s := NewPercentageSlider(0, 100).WithValue(50)
	r := s.CreateRenderer().(*modernSliderRenderer)
	minSize := r.MinSize()
	if minSize.Width <= 0 || minSize.Height <= 0 {
		t.Errorf("MinSize inválido: %v", minSize)
	}

	objs := r.Objects()
	if len(objs) != 3 { // innerSlider, badgeBg, valueText
		t.Errorf("esperaba 3 objetos canvas en layout inline, obtuve %d", len(objs))
	}

	r.Layout(fyne.NewSize(250, 40))

	// 2. Layout con Título superior
	s.WithTitle("Porcentaje de reducción")
	r.Refresh()

	objsTitle := r.Objects()
	if len(objsTitle) != 4 { // titleText, innerSlider, badgeBg, valueText
		t.Errorf("esperaba 4 objetos canvas con título, obtuve %d", len(objsTitle))
	}

	r.Layout(fyne.NewSize(300, 60))

	// 3. Estado Deshabilitado
	s.Disable()
	r.Refresh()
	if !r.innerSlider.Disabled() {
		t.Errorf("el innerSlider debería estar deshabilitado al desactivar ModernSlider")
	}

	s.Enable()
	r.Refresh()
	if r.innerSlider.Disabled() {
		t.Errorf("el innerSlider debería estar habilitado tras Enable()")
	}

	r.Destroy()
}
