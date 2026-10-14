package uikit

import (
	"image/color"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestNewSpinner(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	s := NewSpinner()
	if s == nil {
		t.Fatal("esperaba que NewSpinner() no fuera nil")
	}
	if s.size != SpinnerMedium {
		t.Errorf("esperaba tamaño por defecto %v, obtuve %v", SpinnerMedium, s.size)
	}
	if !s.showTrack {
		t.Errorf("esperaba que showTrack fuera true por defecto")
	}
	if s.IsRunning() {
		t.Errorf("el spinner no debería estar corriendo al crearse")
	}
}

func TestSpinnerConstructors(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	sSmall := NewSpinnerSmall()
	if sSmall.size != SpinnerSmall {
		t.Errorf("esperaba %v, obtuve %v", SpinnerSmall, sSmall.size)
	}

	sMed := NewSpinnerMedium()
	if sMed.size != SpinnerMedium {
		t.Errorf("esperaba %v, obtuve %v", SpinnerMedium, sMed.size)
	}

	sLarge := NewSpinnerLarge()
	if sLarge.size != SpinnerLarge {
		t.Errorf("esperaba %v, obtuve %v", SpinnerLarge, sLarge.size)
	}

	sPrimary := NewPrimarySpinner()
	if sPrimary.color == nil {
		t.Errorf("esperaba color definido en NewPrimarySpinner")
	}

	sSuccess := NewSuccessSpinner()
	if sSuccess.color == nil {
		t.Errorf("esperaba color definido en NewSuccessSpinner")
	}

	sDanger := NewDangerSpinner()
	if sDanger.color == nil {
		t.Errorf("esperaba color definido en NewDangerSpinner")
	}

	sAccent := NewAccentSpinner()
	if sAccent.color == nil {
		t.Errorf("esperaba color definido en NewAccentSpinner")
	}
}

func TestSpinnerFluentAPI(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	s := NewSpinner().
		WithSize(32).
		WithHexColor("#3b82f6").
		WithTrack(false).
		WithStrokeRatio(0.80)

	if s.size != 32 {
		t.Errorf("esperaba tamaño 32, obtuve %v", s.size)
	}
	if s.showTrack != false {
		t.Errorf("esperaba showTrack false")
	}
	if s.strokeRatio != 0.80 {
		t.Errorf("esperaba strokeRatio 0.80, obtuve %v", s.strokeRatio)
	}

	// MinSize
	minSize := s.MinSize()
	if minSize.Width != 32 || minSize.Height != 32 {
		t.Errorf("MinSize incorrecto: %v", minSize)
	}
}

func TestSpinnerStartStopToggle(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	s := NewSpinner()
	if s.IsRunning() {
		t.Errorf("inicialmente no debería estar corriendo")
	}

	s.Start()
	if !s.IsRunning() {
		t.Errorf("tras Start(), IsRunning debería ser true")
	}

	s.Stop()
	if s.IsRunning() {
		t.Errorf("tras Stop(), IsRunning debería ser false")
	}

	s.Toggle()
	if !s.IsRunning() {
		t.Errorf("tras Toggle() sobre detenido, IsRunning debería ser true")
	}

	s.Toggle()
	if s.IsRunning() {
		t.Errorf("tras Toggle() sobre activo, IsRunning debería ser false")
	}
}

func TestSpinnerRenderer(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	s := NewPrimarySpinner(SpinnerLarge)
	w := test.NewWindow(s)
	defer w.Close()

	w.Resize(fyne.NewSize(100, 100))

	r := test.WidgetRenderer(s)
	if r == nil {
		t.Fatal("el renderer no puede ser nil")
	}

	objs := r.Objects()
	if len(objs) != 2 {
		t.Fatalf("esperaba 2 objetos canvas (track y spin arc), obtuve %d", len(objs))
	}

	s.Start()
	time.Sleep(20 * time.Millisecond)
	s.Refresh()

	s.Stop()
	s.Refresh()

	s.WithTrackColor(color.NRGBA{R: 200, G: 200, B: 200, A: 255})
	s.Refresh()

	r.Destroy()
	if s.IsRunning() {
		t.Errorf("tras Destroy() el spinner debería estar detenido")
	}
}

func TestEaseInOut(t *testing.T) {
	if easeInOut(0) != 0 {
		t.Errorf("easeInOut(0) debería ser 0, obtuve %v", easeInOut(0))
	}
	if easeInOut(1) != 1 {
		t.Errorf("easeInOut(1) debería ser 1, obtuve %v", easeInOut(1))
	}
	if easeInOut(-0.5) != 0 {
		t.Errorf("easeInOut negativo debería truncarse a 0")
	}
	if easeInOut(1.5) != 1 {
		t.Errorf("easeInOut > 1 debería truncarse a 1")
	}
	mid := easeInOut(0.5)
	if mid != 0.5 {
		t.Errorf("easeInOut(0.5) debería ser 0.5, obtuve %v", mid)
	}
}
