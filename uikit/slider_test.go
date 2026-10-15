package uikit

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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
	if !age.showSteppers {
		t.Errorf("se esperaban steppers activados en NewAgeSlider")
	}
	if !age.showMinMax {
		t.Errorf("se esperaba min/max activado en NewAgeSlider")
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

func TestSliderSteppersAndScroll(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	s := NewSlider(0, 10).WithStep(1).WithValue(5)

	s.StepUp()
	if s.Value() != 6 {
		t.Errorf("tras StepUp esperaba 6, obtuve %v", s.Value())
	}

	s.StepDown()
	if s.Value() != 5 {
		t.Errorf("tras StepDown esperaba 5, obtuve %v", s.Value())
	}

	// Scroll wheel hacia arriba
	s.Scrolled(&fyne.ScrollEvent{
		Scrolled: fyne.Delta{DY: 10},
	})
	if s.Value() != 6 {
		t.Errorf("tras Scrolled up esperaba 6, obtuve %v", s.Value())
	}

	// Scroll wheel hacia abajo
	s.Scrolled(&fyne.ScrollEvent{
		Scrolled: fyne.Delta{DY: -10},
	})
	if s.Value() != 5 {
		t.Errorf("tras Scrolled down esperaba 5, obtuve %v", s.Value())
	}
}

func TestSliderScrollInWindow(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	w := test.NewWindow(nil)
	defer w.Close()

	s := NewPercentageSlider(0, 100).WithStep(5).WithValue(50)
	vbox := container.NewVBox(s)
	vscroll := container.NewVScroll(vbox)
	w.SetContent(vscroll)
	w.Resize(fyne.NewSize(400, 400))

	pos := fyne.NewPos(100, 20)
	test.Scroll(w.Canvas(), pos, 0, 10)
	t.Logf("Valor tras test.Scroll at (100, 20) DY=10: %v", s.Value())

	posTrack := fyne.NewPos(100, 35)
	test.Scroll(w.Canvas(), posTrack, 0, 10)
	t.Logf("Valor tras test.Scroll at (100, 35) DY=10: %v", s.Value())

	// Also test when vscroll has scrolled down
	vscroll.Offset.Y = 50
	vscroll.Refresh()
	test.Scroll(w.Canvas(), fyne.NewPos(100, 35), 0, 10)
	t.Logf("Valor tras test.Scroll with Offset.Y=50: %v", s.Value())
}

func TestSliderEditableBadge(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	s := NewSlider(0, 100).WithValue(20).WithEditableBadge(true)
	r := s.CreateRenderer().(*modernSliderRenderer)

	r.startEditing()
	if !r.isEditing {
		t.Fatalf("se esperaba isEditing = true tras startEditing")
	}

	// Simular escritura de nuevo valor y Enter
	r.valueEntry.OnSubmitted("75")
	if r.isEditing {
		t.Errorf("se esperaba isEditing = false tras OnSubmitted")
	}
	if s.Value() != 75 {
		t.Errorf("se esperaba que el valor se actualizara a 75, obtuve %v", s.Value())
	}
}

func TestSliderRendererLayout(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	// 1. Layout Inline simple sin título ni steppers
	s := NewSlider(0, 100).WithValue(50)
	r := s.CreateRenderer().(*modernSliderRenderer)
	minSize := r.MinSize()
	if minSize.Width <= 0 || minSize.Height <= 0 {
		t.Errorf("MinSize inválido: %v", minSize)
	}

	objs := r.Objects()
	// themedSlider, badgeBg, badgeTap, valueText = 4
	if len(objs) < 3 {
		t.Errorf("esperaba al menos 3 objetos canvas en layout inline, obtuve %d", len(objs))
	}

	r.Layout(fyne.NewSize(250, 40))

	// 2. Layout con Título superior, Steppers y MinMax
	s.WithTitle("Porcentaje de reducción").
		WithSteppers(true).
		WithMinMax(true).
		WithRangeLabels("Mínimo", "Máximo")
	r.Refresh()

	objsTitle := r.Objects()
	// titleText, stepperMinus, stepperPlus, themedSlider, minText, maxText, badgeBg, badgeTap, valueText
	if len(objsTitle) < 8 {
		t.Errorf("esperaba al menos 8 objetos canvas con todas las extensiones, obtuve %d", len(objsTitle))
	}

	r.Layout(fyne.NewSize(320, 70))

	// 3. Estado Deshabilitado
	s.Disable()
	r.Refresh()
	if !r.innerSlider.Disabled() {
		t.Errorf("el innerSlider debería estar deshabilitado al desactivar ModernSlider")
	}
	if !r.stepperMinus.disabled {
		t.Errorf("stepperMinus debería estar deshabilitado")
	}

	s.Enable()
	r.Refresh()
	if r.innerSlider.Disabled() {
		t.Errorf("el innerSlider debería estar habilitado tras Enable()")
	}

	r.Destroy()
}
