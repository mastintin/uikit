package uikit

import (
	"image/png"
	"os"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
)

func TestModernDatePicker_ManualInputAndFormatting(t *testing.T) {
	test.NewApp()

	dp := NewModernDatePicker(time.Date(2026, 10, 3, 0, 0, 0, 0, time.UTC), nil)

	// 1. Initial date
	if dp.Date().Year() != 2026 {
		t.Fatalf("expected initial year 2026, got %d", dp.Date().Year())
	}

	// 2. User writes "03/12/2029" in the text box without pressing enter
	dp.entryDate.SetText("03/12/2029")
	d1 := dp.Date()
	if d1.Year() != 2029 || d1.Month() != time.December || d1.Day() != 3 {
		t.Errorf("expected 2029-12-03, got %s", d1.Format("2006-01-02"))
	}

	// 3. User writes "3/12/2029" (single digit day)
	dp.entryDate.SetText("3/12/2029")
	d2 := dp.Date()
	if d2.Year() != 2029 || d2.Month() != time.December || d2.Day() != 3 {
		t.Errorf("expected 2029-12-03 from '3/12/2029', got %s", d2.Format("2006-01-02"))
	}

	// 4. User writes "3-12-2029"
	dp.entryDate.SetText("3-12-2029")
	d3 := dp.Date()
	if d3.Year() != 2029 || d3.Month() != time.December || d3.Day() != 3 {
		t.Errorf("expected 2029-12-03 from '3-12-2029', got %s", d3.Format("2006-01-02"))
	}

	// 5. User writes "2029-12-03" (ISO)
	dp.entryDate.SetText("2029-12-03")
	d4 := dp.Date()
	if d4.Year() != 2029 || d4.Month() != time.December || d4.Day() != 3 {
		t.Errorf("expected 2029-12-03 from '2029-12-03', got %s", d4.Format("2006-01-02"))
	}
}

func TestModernDatePicker_ShowNavButtons(t *testing.T) {
	test.NewApp()

	dp := NewModernDatePicker(time.Now(), nil)

	// Por defecto ShowNavButtons debe ser true
	if !dp.ShowNavButtons {
		t.Errorf("Expected ShowNavButtons to be true by default")
	}
	if len(dp.Container.Objects) != 5 {
		t.Errorf("Expected 5 objects in Container when ShowNavButtons=true, got %d", len(dp.Container.Objects))
	}

	// Ocultar botones de navegación (+ - mes y año)
	dp.SetShowNavButtons(false)
	if dp.ShowNavButtons {
		t.Errorf("Expected ShowNavButtons to be false after SetShowNavButtons(false)")
	}
	if len(dp.Container.Objects) != 1 {
		t.Errorf("Expected 1 object in Container when ShowNavButtons=false, got %d", len(dp.Container.Objects))
	}

	// Volver a activar botones de navegación
	dp.SetShowNavButtons(true)
	if !dp.ShowNavButtons {
		t.Errorf("Expected ShowNavButtons to be true after SetShowNavButtons(true)")
	}
	if len(dp.Container.Objects) != 5 {
		t.Errorf("Expected 5 objects in Container when ShowNavButtons=true, got %d", len(dp.Container.Objects))
	}
}

func TestModernDatePicker_CalendarActionItem(t *testing.T) {
	test.NewApp()

	dp := NewModernDatePicker(time.Date(2026, 10, 11, 0, 0, 0, 0, time.Local), nil)

	// Debe tener el botón de calendario en ActionItem del textbox
	if dp.entryDate.ActionItem == nil {
		t.Fatalf("Expected entryDate.ActionItem not to be nil")
	}

	dp.entryDate.Resize(fyne.NewSize(140, 38))
	r := test.WidgetRenderer(dp.entryDate)
	r.Layout(fyne.NewSize(140, 38))

	t.Logf("entryDate minSize: %v", dp.entryDate.MinSize())
	t.Logf("ActionItem minSize: %v", dp.entryDate.ActionItem.MinSize())
	t.Logf("ActionItem size: %v, pos: %v, visible: %v", dp.entryDate.ActionItem.Size(), dp.entryDate.ActionItem.Position(), dp.entryDate.ActionItem.Visible())

	btnR := test.WidgetRenderer(dp.btnCalendar)
	t.Logf("btnCalendar objects count: %d", len(btnR.Objects()))
	for i, obj := range btnR.Objects() {
		t.Logf("  obj %d: %T, size: %v, pos: %v, visible: %v", i, obj, obj.Size(), obj.Position(), obj.Visible())
		if img, ok := obj.(*canvas.Image); ok {
			t.Logf("    image resource: %v", img.Resource != nil)
		}
	}

	// Verificar refresco de colores
	dp.Refresh()
}

func TestDatePickerSnapshot(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	dp := NewModernDatePicker(time.Date(2026, 10, 11, 0, 0, 0, 0, time.Local), nil)
	w := test.NewWindow(dp.Container)
	w.Resize(fyne.NewSize(400, 100))

	img := w.Canvas().Capture()
	t.Logf("Captured image bounds: %v", img.Bounds())

	os.MkdirAll("/Users/manuel/.gemini/antigravity/brain/1a2b4731-1e5b-4352-b5d1-c6f5e84045da/scratch", 0755)
	f, err := os.Create("/Users/manuel/.gemini/antigravity/brain/1a2b4731-1e5b-4352-b5d1-c6f5e84045da/scratch/datepicker_snapshot.png")
	if err == nil {
		defer f.Close()
		_ = png.Encode(f, img)
	}

	dpNoNav := NewModernDatePicker(time.Date(2026, 10, 11, 0, 0, 0, 0, time.Local), nil).SetShowNavButtons(false)
	wNoNav := test.NewWindow(dpNoNav.Container)
	wNoNav.Resize(fyne.NewSize(200, 100))
	imgNoNav := wNoNav.Canvas().Capture()
	f2, err2 := os.Create("/Users/manuel/.gemini/antigravity/brain/1a2b4731-1e5b-4352-b5d1-c6f5e84045da/scratch/datepicker_nonav_snapshot.png")
	if err2 == nil {
		defer f2.Close()
		_ = png.Encode(f2, imgNoNav)
	}
}

func TestModernCalendar_NavigationAndSelection(t *testing.T) {
	test.NewApp()

	initialDate := time.Date(2026, 10, 15, 0, 0, 0, 0, time.Local)
	var selectedDate time.Time

	cal := NewModernCalendar(initialDate, func(selected time.Time) {
		selectedDate = selected
	}, nil)

	if cal == nil {
		t.Fatal("NewModernCalendar returned nil")
	}
	if cal.lblMonthYear == nil || cal.lblMonthYear.Text != "Octubre 2026" {
		t.Errorf("Expected 'Octubre 2026', got '%s'", cal.lblMonthYear.Text)
	}

	// Navegar un mes adelante
	cal.navigate(1)
	if cal.lblMonthYear.Text != "Noviembre 2026" {
		t.Errorf("Expected 'Noviembre 2026', got '%s'", cal.lblMonthYear.Text)
	}

	// Navegar un año atrás
	cal.navigate(-12)
	if cal.lblMonthYear.Text != "Noviembre 2025" {
		t.Errorf("Expected 'Noviembre 2025', got '%s'", cal.lblMonthYear.Text)
	}

	// Seleccionar fecha
	target := time.Date(2025, 11, 20, 0, 0, 0, 0, time.Local)
	cal.selectDate(target)

	if !isSameDay(selectedDate, target) {
		t.Errorf("Expected selected date %v, got %v", target, selectedDate)
	}
}

func TestModernDatePicker_ThemeChangeUpdatesButtons(t *testing.T) {
	test.NewApp()

	SetAppTheme(ThemeLight)
	dp := NewModernDatePicker(time.Now(), nil)

	r := test.WidgetRenderer(dp.btnPrevYear).(*modernButtonRenderer)
	lightBg := ColorToHex(r.bg.FillColor)

	// Cambiar a tema Dark
	SetAppTheme(ThemeDark)

	darkBg := ColorToHex(r.bg.FillColor)
	if darkBg == lightBg {
		t.Errorf("Expected button background to change on theme change, but remained %s", darkBg)
	}

	palDark := GetThemePalette(ThemeDark)
	expectedDarkBg := ColorToHex(palDark.Buttons.Default.Bg)
	if darkBg != expectedDarkBg {
		t.Errorf("Expected dark button bg %s, got %s", expectedDarkBg, darkBg)
	}

	// Restaurar a tema Light
	SetAppTheme(ThemeLight)
	restoredBg := ColorToHex(r.bg.FillColor)
	if restoredBg != lightBg {
		t.Errorf("Expected restored bg %s, got %s", lightBg, restoredBg)
	}
}

func TestCalculateSmartPopupPosition_FlipUp(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	target := canvas.NewRectangle(nil)
	target.SetMinSize(fyne.NewSize(150, 36))
	target.Resize(fyne.NewSize(150, 36))

	w := a.NewWindow("Test Smart Popup")
	defer w.Close()

	root := container.NewWithoutLayout(target)
	w.SetContent(root)
	w.Resize(fyne.NewSize(800, 600))
	c := w.Canvas()

	// Caso 1: Target en la parte inferior de la ventana (ej. Y = 520)
	target.Move(fyne.NewPos(100, 520))
	popupSize := fyne.NewSize(280, 290)

	t.Logf("Canvas Size: %v, Target Size: %v, Target Pos: %v", c.Size(), target.Size(), target.Position())

	posBottom := calculateSmartPopupPosition(c, target, popupSize)
	// Debe recolocarse HACIA ARRIBA (posY < target.Position().Y)
	if posBottom.Y >= target.Position().Y {
		t.Errorf("Expected popup to flip UP above %v, got Y = %v", target.Position().Y, posBottom.Y)
	}
	if posBottom.Y < 4 {
		t.Errorf("Expected popup Y to be within screen bounds, got %v", posBottom.Y)
	}

	// Caso 2: Target en la parte superior de la ventana (ej. Y = 50)
	target.Move(fyne.NewPos(100, 50))
	posTop := calculateSmartPopupPosition(c, target, popupSize)
	// Debe abrirse HACIA ABAJO (posY > 50 + 36)
	if posTop.Y < 50+36 {
		t.Errorf("Expected popup to open DOWN below target, got Y = %v", posTop.Y)
	}
}


