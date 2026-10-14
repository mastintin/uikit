package uikit

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

func TestLightTheme_Colors(t *testing.T) {
	th := NewLightTheme()
	if th.Name() != ThemeLight {
		t.Errorf("Esperado nombre '%s', obtenido '%s'", ThemeLight, th.Name())
	}
	if th.IsDark() {
		t.Errorf("LightTheme no debería ser oscuro")
	}

	// 1. Fondo de la app (Slate 50)
	appBg := th.Color(ColorNameAppBackground, theme.VariantLight)
	if appBg != (color.NRGBA{R: 248, G: 250, B: 252, A: 255}) {
		t.Errorf("ColorNameAppBackground de LightTheme incorrecto: %v", appBg)
	}

	// 2. Tarjeta en blanco puro
	cardBg := th.Color(ColorNameCardBg, theme.VariantLight)
	if cardBg != ColorWhite {
		t.Errorf("ColorNameCardBg de LightTheme incorrecto: %v", cardBg)
	}

	// 3. Texto principal en slate 900 (#0f172a)
	textPrimary := th.Color(ColorNameTextPrimary, theme.VariantLight)
	if textPrimary != (color.NRGBA{R: 15, G: 23, B: 42, A: 255}) {
		t.Errorf("ColorNameTextPrimary de LightTheme incorrecto: %v", textPrimary)
	}

	// 4. Cursor y foco primario con el verde corporativo de Nice (#107c41)
	primary := th.Color(theme.ColorNamePrimary, theme.VariantLight)
	expectedGreen := HexColor("#107c41")
	if primary != expectedGreen {
		t.Errorf("theme.ColorNamePrimary en LightTheme debe ser el verde de Nice %v, obtenido: %v", expectedGreen, primary)
	}
}

func TestTheme_SwitchingHelpers(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()

	// Available themes
	themes := AvailableThemes()
	if len(themes) != 2 {
		t.Errorf("Esperados 2 temas nativos, obtenidos %d", len(themes))
	}

	// Cycle
	if GetNextTheme(ThemeLight) != ThemeDark {
		t.Errorf("Next after Light debería ser Dark")
	}
	if GetNextTheme(ThemeDark) != ThemeLight {
		t.Errorf("Next after Dark debería ser Light")
	}

	// SetAppTheme
	th := SetAppTheme(ThemeLight)
	if th.Name() != ThemeLight {
		t.Errorf("SetAppTheme(light) devolvió %s", th.Name())
	}
	if GetCurrentThemeName() != ThemeLight {
		t.Errorf("GetCurrentThemeName devolvió %s", GetCurrentThemeName())
	}

	thDark := SetAppTheme(ThemeDark)
	if thDark.Name() != ThemeDark {
		t.Errorf("SetAppTheme(dark) devolvió %s", thDark.Name())
	}
	if !thDark.IsDark() {
		t.Errorf("thDark debería ser oscuro")
	}
	_ = app
}
