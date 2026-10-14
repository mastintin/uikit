package uikit

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestColorToNRGBA_Unpremultiplication(t *testing.T) {
	// color.RGBA almacena valores alfa-premultiplicados:
	// Si R = 128 con Alpha = 128 (50% de opacidad), el valor no premultiplicado debe ser 255 (rojo puro al 50% de alfa)
	premultiplied := color.RGBA{R: 128, G: 0, B: 0, A: 128}
	unpremultiplied := ColorToNRGBA(premultiplied)

	if unpremultiplied.A != 128 {
		t.Errorf("Esperado A=128, obtenido %d", unpremultiplied.A)
	}
	// Permite tolerancia de +/- 2 por redondeo en enteros de 16 bits
	if unpremultiplied.R < 253 || unpremultiplied.R > 255 {
		t.Errorf("Esperado R des-premultiplicado ~255, obtenido %d", unpremultiplied.R)
	}
	if unpremultiplied.G != 0 || unpremultiplied.B != 0 {
		t.Errorf("Esperado G=0 B=0, obtenido G=%d B=%d", unpremultiplied.G, unpremultiplied.B)
	}

	// Nil debe devolver NRGBA vacío
	nilResult := ColorToNRGBA(nil)
	if nilResult != (color.NRGBA{}) {
		t.Errorf("Esperado NRGBA vacio para nil, obtenido %v", nilResult)
	}

	// Totalmente transparente (Alpha=0)
	transp := color.RGBA{R: 0, G: 0, B: 0, A: 0}
	if ColorToNRGBA(transp).A != 0 {
		t.Errorf("Esperado A=0 para transparente")
	}

	// NRGBA nativo debe devolverse idéntico
	original := color.NRGBA{R: 16, G: 124, B: 65, A: 255}
	if ColorToNRGBA(original) != original {
		t.Errorf("NRGBA nativo debería preservarse sin cambios")
	}
}

func TestHexColor_And_ColorToHex(t *testing.T) {
	// 1. #RRGGBB estándar
	c1 := HexColor("#107c41")
	n1 := ColorToNRGBA(c1)
	if n1.R != 16 || n1.G != 124 || n1.B != 65 || n1.A != 255 {
		t.Errorf("HexColor(#107c41) falló: %v", n1)
	}
	hex1 := ColorToHex(c1)
	if hex1 != "#107c41" {
		t.Errorf("ColorToHex falló: esperado #107c41, obtenido %s", hex1)
	}

	// 2. #RGB corto (#fff)
	c2 := HexColor("#fff")
	n2 := ColorToNRGBA(c2)
	if n2.R != 255 || n2.G != 255 || n2.B != 255 || n2.A != 255 {
		t.Errorf("HexColor(#fff) falló: %v", n2)
	}

	// 3. Hex inválido devuelve Transparent
	cInv := HexColor("invalido")
	if cInv != color.Transparent {
		t.Errorf("Hex inválido debería devolver color.Transparent")
	}

	// 4. Hex con canal alfa #rrggbbaa
	cAlpha := HexColor("#107c4180")
	hexAlpha := ColorToHex(cAlpha)
	if hexAlpha != "#107c4180" {
		t.Errorf("ColorToHex con alfa falló: esperado #107c4180, obtenido %s", hexAlpha)
	}
}

func TestColorWithAlpha(t *testing.T) {
	base := HexColor("#107c41")
	half := ColorWithAlpha(base, 0.5)
	nHalf := ColorToNRGBA(half)

	if nHalf.R != 16 || nHalf.G != 124 || nHalf.B != 65 {
		t.Errorf("Canales RGB no deben alterarse con ColorWithAlpha: %v", nHalf)
	}
	if nHalf.A < 127 || nHalf.A > 128 {
		t.Errorf("Esperado A ~128, obtenido %d", nHalf.A)
	}

	// Clamping
	if ColorToNRGBA(ColorWithAlpha(base, -0.5)).A != 0 {
		t.Errorf("Factor negativo debe clamplear A a 0")
	}
	if ColorToNRGBA(ColorWithAlpha(base, 2.0)).A != 255 {
		t.Errorf("Factor > 1 debe clamplear A a 255")
	}
}

func TestBlendColor(t *testing.T) {
	cWhite := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	cBlack := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

	// Extremos t=0 y t=1
	if ColorToNRGBA(BlendColor(cBlack, cWhite, 0.0)) != cBlack {
		t.Errorf("Blend t=0 debe devolver c1")
	}
	if ColorToNRGBA(BlendColor(cBlack, cWhite, 1.0)) != cWhite {
		t.Errorf("Blend t=1 debe devolver c2")
	}

	// Punto medio t=0.5 -> gris ~128
	mid := ColorToNRGBA(BlendColor(cBlack, cWhite, 0.5))
	if mid.R < 127 || mid.R > 128 {
		t.Errorf("Esperado canal medio ~128, obtenido %d", mid.R)
	}
}

func TestColorBrighter_And_ColorDarker(t *testing.T) {
	base := HexColor("#107c41") // R:16, G:124, B:65
	nBase := ColorToNRGBA(base)

	// Aclarar 20%
	brighter := ColorBrighter(base, 0.20)
	nBright := ColorToNRGBA(brighter)

	if nBright.G <= nBase.G {
		t.Errorf("ColorBrighter debe aumentar el canal principal: base=%d bright=%d", nBase.G, nBright.G)
	}

	// Oscurecer 20%
	darker := ColorDarker(base, 0.20)
	nDark := ColorToNRGBA(darker)

	if nDark.G >= nBase.G {
		t.Errorf("ColorDarker debe disminuir el canal principal: base=%d dark=%d", nBase.G, nDark.G)
	}

	// Factor 0 no altera el color
	if ColorToNRGBA(ColorBrighter(base, 0)) != nBase {
		t.Errorf("Factor 0 no debe alterar el color")
	}
}

func TestLuminance_And_Contrast(t *testing.T) {
	cWhite := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	cDarkSlate := color.NRGBA{R: 15, G: 23, B: 42, A: 255}

	lumWhite := ColorLuminance(cWhite)
	if lumWhite < 0.99 {
		t.Errorf("Luminancia del blanco debe ser ~1.0, obtenida %f", lumWhite)
	}
	if IsColorDark(cWhite) {
		t.Errorf("Blanco no es un color oscuro")
	}
	// Texto sobre blanco debe ser oscuro
	if ContrastTextColor(cWhite) != ColorSlate900 {
		t.Errorf("Texto sobre blanco debe ser ColorSlate900")
	}

	lumDark := ColorLuminance(cDarkSlate)
	if lumDark > 0.15 {
		t.Errorf("Luminancia de Slate 900 debe ser baja, obtenida %f", lumDark)
	}
	if !IsColorDark(cDarkSlate) {
		t.Errorf("Slate 900 debe clasificarse como oscuro")
	}
	// Texto sobre fondo oscuro debe ser blanco
	if ContrastTextColor(cDarkSlate) != ColorWhite {
		t.Errorf("Texto sobre fondo oscuro debe ser ColorWhite")
	}
}

func TestButtonForegroundFor(t *testing.T) {
	if ButtonForegroundFor(widget.HighImportance) != theme.ColorNameForegroundOnPrimary {
		t.Errorf("HighImportance debe mapear a ForegroundOnPrimary")
	}
	if ButtonForegroundFor(widget.DangerImportance) != theme.ColorNameForegroundOnError {
		t.Errorf("DangerImportance debe mapear a ForegroundOnError")
	}
	if ButtonForegroundFor(widget.SuccessImportance) != theme.ColorNameForegroundOnSuccess {
		t.Errorf("SuccessImportance debe mapear a ForegroundOnSuccess")
	}
	if ButtonForegroundFor(widget.WarningImportance) != theme.ColorNameForegroundOnWarning {
		t.Errorf("WarningImportance debe mapear a ForegroundOnWarning")
	}
	if ButtonForegroundFor(widget.MediumImportance) != theme.ColorNameForeground {
		t.Errorf("MediumImportance debe mapear a Foreground")
	}
}
