package uikit

import (
	"image/color"
	"image/png"
	"os"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestModernHeaderCard_CreationAndLayout(t *testing.T) {
	_ = test.NewApp()

	header := widget.NewLabel("Banner Header")
	body := widget.NewLabel("Body Content")

	card := NewHeaderCard(header, body)
	if card == nil {
		t.Fatal("NewHeaderCard returned nil")
	}

	card.Resize(fyne.NewSize(400, 200))
	minSize := card.MinSize()
	if minSize.Width <= 0 || minSize.Height <= 0 {
		t.Errorf("Expected positive MinSize, got %v", minSize)
	}

	// Test Fluent API
	card.WithBannerHex("#107c41").
		WithCornerRadius(12).
		WithBannerPadding(10, 10, 15, 15).
		WithBodyPadding(12, 12, 16, 16)

	if card.cornerRadius != 12 {
		t.Errorf("Expected cornerRadius 12, got %f", card.cornerRadius)
	}
	if !card.hasCustomBanner {
		t.Errorf("Expected hasCustomBanner to be true after WithBannerHex")
	}

	card.WithBodyColor(color.RGBA{R: 240, G: 240, B: 240, A: 255})
	if !card.hasCustomBody {
		t.Errorf("Expected hasCustomBody to be true after WithBodyColor")
	}

	// Test Shadow
	if !card.hasShadow {
		t.Errorf("Expected hasShadow to be true by default")
	}
	if card.bgOuter.Shadow.BlurRadius != 6 {
		t.Errorf("Expected shadow BlurRadius 6, got %f", card.bgOuter.Shadow.BlurRadius)
	}

	card.WithElevation(2)
	if card.bgOuter.Shadow.BlurRadius != 10 {
		t.Errorf("Expected shadow BlurRadius 10 at elevation 2, got %f", card.bgOuter.Shadow.BlurRadius)
	}

	card.WithoutShadow()
	if card.hasShadow {
		t.Errorf("Expected hasShadow to be false after WithoutShadow")
	}

	card.Refresh()
}

func TestModernHeaderCardWithTitle(t *testing.T) {
	_ = test.NewApp()

	body := container.NewHBox(widget.NewLabel("Controles de fecha"))
	card := NewHeaderCardWithTitle(
		"calculator",
		"Simulador de Jubilación Ordinaria y Anticipada",
		"Cálculo y proyecciones de pensión ordinaria y anticipada.",
		body,
	)

	if card == nil {
		t.Fatal("NewHeaderCardWithTitle returned nil")
	}
	if card.TitleText == nil || card.TitleText.Text != "Simulador de Jubilación Ordinaria y Anticipada" {
		t.Errorf("Unexpected TitleText: %v", card.TitleText)
	}
	if card.SubtitleText == nil || card.SubtitleText.Text != "Cálculo y proyecciones de pensión ordinaria y anticipada." {
		t.Errorf("Unexpected SubtitleText: %v", card.SubtitleText)
	}

	card.Resize(fyne.NewSize(500, 160))
	minSize := card.MinSize()
	if minSize.Width <= 0 || minSize.Height <= 0 {
		t.Errorf("Expected positive MinSize, got %v", minSize)
	}

	// Verificar soporte Glassmorphism
	if !card.IsGlass() {
		t.Errorf("Expected IsGlass to be true by default")
	}
	if card.iconGlassPillBg == nil {
		t.Errorf("Expected iconGlassPillBg to be initialized")
	}

	card.WithGlass(false)
	if card.IsGlass() {
		t.Errorf("Expected IsGlass to be false after WithGlass(false)")
	}

	card.WithGlass(true)
	if !card.IsGlass() {
		t.Errorf("Expected IsGlass to be true after WithGlass(true)")
	}

	// Cambiar a diferentes temas y verificar Refresh()
	SetAppTheme(ThemeLight)
	card.Refresh()

	SetAppTheme(ThemeDark)
	card.Refresh()

	// Snapshot visual (Light Theme)
	SetAppTheme(ThemeLight)
	card.Refresh()
	w := test.NewWindow(container.NewPadded(card))
	w.Resize(fyne.NewSize(700, 220))
	img := w.Canvas().Capture()
	_ = os.MkdirAll("/Users/manuel/.gemini/antigravity/brain/1a2b4731-1e5b-4352-b5d1-c6f5e84045da/scratch", 0755)
	if f, err := os.Create("/Users/manuel/.gemini/antigravity/brain/1a2b4731-1e5b-4352-b5d1-c6f5e84045da/scratch/header_card_snapshot.png"); err == nil {
		defer f.Close()
		_ = png.Encode(f, img)
	}

	// Snapshot visual (Dark Theme)
	SetAppTheme(ThemeDark)
	card.Refresh()
	wDark := test.NewWindow(container.NewPadded(card))
	wDark.Resize(fyne.NewSize(700, 220))
	imgDark := wDark.Canvas().Capture()
	if f, err := os.Create("/Users/manuel/.gemini/antigravity/brain/1a2b4731-1e5b-4352-b5d1-c6f5e84045da/scratch/header_card_dark_snapshot.png"); err == nil {
		defer f.Close()
		_ = png.Encode(f, imgDark)
	}
	SetAppTheme(ThemeLight)
}
