package uikit

import (
	"image/color"
	"strings"
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

func TestModernButton_PrimaryAndDangerStyling(t *testing.T) {
	a := test.NewApp()
	a.Settings().SetTheme(NewAppTheme().SetForcedVariant(theme.VariantLight))

	btnPrimary := NewPrimaryButton("Nuevo", theme.ContentAddIcon(), func() {})
	if btnPrimary.Variant != ButtonPrimary {
		t.Fatalf("expected ButtonPrimary, got %v", btnPrimary.Variant)
	}

	btnDanger := NewDangerButton("Borrar", theme.DeleteIcon(), func() {})
	if btnDanger.Variant != ButtonDanger {
		t.Fatalf("expected ButtonDanger, got %v", btnDanger.Variant)
	}

	// Probar renderer de botón primario
	rPrimary := btnPrimary.CreateRenderer().(*modernButtonRenderer)
	rPrimary.applyStyle()

	// El fondo normal del botón primario debe ser azul suave (#eff6ff)
	expectedFill := color.NRGBA{R: 239, G: 246, B: 255, A: 255}
	if rPrimary.bg.FillColor != expectedFill {
		t.Errorf("expected primary fill %v, got %v", expectedFill, rPrimary.bg.FillColor)
	}

	// El texto debe ser azul vibrante (#2563eb)
	expectedText := color.NRGBA{R: 37, G: 99, B: 235, A: 255}
	if rPrimary.label.Color != expectedText {
		t.Errorf("expected primary text color %v, got %v", expectedText, rPrimary.label.Color)
	}

	// El borde debe ser azul claro (#bfdbfe)
	expectedStroke := color.NRGBA{R: 191, G: 219, B: 254, A: 255}
	if rPrimary.bg.StrokeColor != expectedStroke {
		t.Errorf("expected primary stroke %v, got %v", expectedStroke, rPrimary.bg.StrokeColor)
	}

	// Probar estado hover
	btnPrimary.hovered = true
	rPrimary.applyStyle()
	expectedHoverFill := color.NRGBA{R: 219, G: 234, B: 254, A: 255} // #dbeafe
	if rPrimary.bg.FillColor != expectedHoverFill {
		t.Errorf("expected primary hover fill %v, got %v", expectedHoverFill, rPrimary.bg.FillColor)
	}

	// Probar estado pressed
	btnPrimary.hovered = false
	btnPrimary.pressed = true
	rPrimary.applyStyle()
	expectedPressFill := color.NRGBA{R: 191, G: 219, B: 254, A: 255} // #bfdbfe
	if rPrimary.bg.FillColor != expectedPressFill {
		t.Errorf("expected primary press fill %v, got %v", expectedPressFill, rPrimary.bg.FillColor)
	}

	// Probar estado disabled
	btnPrimary.pressed = false
	btnPrimary.Disable()
	rPrimary.applyStyle()
	expectedDisabledFill := ThemeColor(ColorNameBtnPrimaryBg)
	if rPrimary.bg.FillColor != expectedDisabledFill {
		t.Errorf("expected primary disabled fill %v, got %v", expectedDisabledFill, rPrimary.bg.FillColor)
	}
}

func TestModernButton_DimensionsAndIconResolution(t *testing.T) {
	btn := NewButton("Default", "Inicio", "dashboard", false, func() {})
	if btn.Icon == nil {
		t.Fatal("expected dashboard icon to resolve, got nil")
	}

	btn.SetHeight(36)
	btn.SetMinWidth(120)

	r := btn.CreateRenderer().(*modernButtonRenderer)
	ms := r.MinSize()
	if ms.Height != 36 {
		t.Errorf("expected height 36, got %v", ms.Height)
	}
	if ms.Width < 120 {
		t.Errorf("expected width >= 120, got %v", ms.Width)
	}

	// Verificar tintado dinámico del icono SVG según el texto
	btn.SetCustomColors(&ModernButtonColors{
		NormalFg: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	})
	r.applyStyle()
	if r.icon.Resource == nil {
		t.Fatal("expected icon resource to be present, got nil")
	}
	content := string(r.icon.Resource.Content())
	if !strings.Contains(content, "#ffffff") {
		t.Errorf("expected SVG to contain #ffffff stroke, got %s", content)
	}
}

func TestModernButton_DarkThemeGoJubilaStyling(t *testing.T) {
	a := test.NewApp()
	a.Settings().SetTheme(NewDarkTheme())

	// Botón Primario en modo oscuro (Indigo 600, White, Indigo 500)
	btnPrimary := NewButton("Primary", "Simulador Completo", "calculator", false, func() {})
	rPrimary := btnPrimary.CreateRenderer().(*modernButtonRenderer)
	rPrimary.applyStyle()

	expectedPrimaryBg := color.NRGBA{R: 79, G: 70, B: 229, A: 255}
	if rPrimary.bg.FillColor != expectedPrimaryBg {
		t.Errorf("expected dark primary bg %v, got %v", expectedPrimaryBg, rPrimary.bg.FillColor)
	}
	expectedPrimaryFg := ColorWhite
	if rPrimary.label.Color != expectedPrimaryFg {
		t.Errorf("expected dark primary text %v, got %v", expectedPrimaryFg, rPrimary.label.Color)
	}
	expectedPrimaryBorder := color.NRGBA{R: 99, G: 102, B: 241, A: 255}
	if rPrimary.bg.StrokeColor != expectedPrimaryBorder {
		t.Errorf("expected dark primary border %v, got %v", expectedPrimaryBorder, rPrimary.bg.StrokeColor)
	}

	// Botón Default en modo oscuro (Slate 800, Slate 100, Slate 700)
	btnDefault := NewButton("Default", "Exportar PDF", "file-text", false, func() {})
	rDefault := btnDefault.CreateRenderer().(*modernButtonRenderer)
	rDefault.applyStyle()

	expectedDefaultBg := color.NRGBA{R: 30, G: 41, B: 59, A: 255}
	if rDefault.bg.FillColor != expectedDefaultBg {
		t.Errorf("expected dark default bg %v, got %v", expectedDefaultBg, rDefault.bg.FillColor)
	}
	expectedDefaultFg := color.NRGBA{R: 241, G: 245, B: 249, A: 255}
	if rDefault.label.Color != expectedDefaultFg {
		t.Errorf("expected dark default text %v, got %v", expectedDefaultFg, rDefault.label.Color)
	}
	expectedDefaultBorder := color.NRGBA{R: 51, G: 65, B: 85, A: 255}
	if rDefault.bg.StrokeColor != expectedDefaultBorder {
		t.Errorf("expected dark default border %v, got %v", expectedDefaultBorder, rDefault.bg.StrokeColor)
	}
}

func TestModernButton_GlassConfiguration(t *testing.T) {
	// Asegurar estado global limpio
	SetGlobalGlassButtons(false)
	if IsGlobalGlass() {
		t.Errorf("Expected IsGlobalGlass to be false initially")
	}

	btn := NewButton("Primary", "Calcular", "calculator", false, func() {})
	// Por defecto hereda global (false)
	if btn.IsGlass() {
		t.Errorf("Expected btn.IsGlass to be false by default")
	}

	// 1. Activar variable global
	SetGlobalGlassButtons(true)
	if !btn.IsGlass() {
		t.Errorf("Expected btn.IsGlass to be true when global is enabled")
	}

	// 2. Sobrescribir localmente en el botón (lGlass = false)
	btn.SetGlass(false)
	if btn.IsGlass() {
		t.Errorf("Expected btn.IsGlass to be false when explicitly overridden to false")
	}

	// 3. Sobrescribir localmente a true cuando global es false
	SetGlobalGlassButtons(false)
	btn.SetGlass(true)
	if !btn.IsGlass() {
		t.Errorf("Expected btn.IsGlass to be true when explicitly overridden to true")
	}

	// 4. ResetGlass vuelve a heredar de global
	btn.ResetGlass()
	if btn.IsGlass() {
		t.Errorf("Expected btn.IsGlass to return to false after ResetGlass()")
	}

	// 5. Verificar renderer con Glass en LightTheme
	app := test.NewApp()
	app.Settings().SetTheme(NewLightTheme())
	btn.SetGlass(true)
	r := btn.CreateRenderer().(*modernButtonRenderer)
	r.applyStyle()

	if !r.button.IsGlass() {
		t.Errorf("Expected renderer button to be in glass mode")
	}
	_, _, _, a := r.bg.FillColor.RGBA()
	if a == 0xffff {
		t.Errorf("Expected translucent glass fill, got solid fill")
	}
	// Restablecer global
	SetGlobalGlassButtons(false)
}

func TestModernButton_LoadingState(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	tapped := false
	btn := NewPrimaryButton("Guardar", theme.DocumentSaveIcon(), func() {
		tapped = true
	})

	if btn.IsLoading() {
		t.Errorf("El botón no debería estar en loading por defecto")
	}

	// Probar clic normal
	test.Tap(btn)
	if !tapped {
		t.Errorf("El botón debería haber ejecutado el callback OnTapped")
	}

	// Activar modo loading
	tapped = false
	btn.SetLoading(true)
	if !btn.IsLoading() {
		t.Errorf("El botón debería estar en estado loading")
	}

	// Al estar en loading, el clic debe quedar bloqueado
	test.Tap(btn)
	if tapped {
		t.Errorf("El botón NO debería haber ejecutado el callback mientras loading=true")
	}

	// Verificar renderer
	r := btn.CreateRenderer().(*modernButtonRenderer)
	r.applyStyle()
	r.Layout(btn.MinSize())

	if !r.spinner.Visible() {
		t.Errorf("El spinner debería ser visible cuando loading=true")
	}
	if !r.spinner.IsRunning() {
		t.Errorf("El spinner debería estar en ejecución cuando loading=true")
	}
	if r.icon.Visible() {
		t.Errorf("El icono normal debería estar oculto mientras loading=true")
	}

	// Desactivar loading
	btn.SetLoading(false)
	r.applyStyle()
	r.Layout(btn.MinSize())

	if btn.IsLoading() {
		t.Errorf("El botón ya no debería estar en loading")
	}
	if r.spinner.IsRunning() {
		t.Errorf("El spinner debería haberse detenido tras SetLoading(false)")
	}
	if r.spinner.Visible() {
		t.Errorf("El spinner debería estar oculto tras SetLoading(false)")
	}

	// Clic vuelve a funcionar
	test.Tap(btn)
	if !tapped {
		t.Errorf("El botón debería volver a responder a clics tras SetLoading(false)")
	}

	r.Destroy()
}

