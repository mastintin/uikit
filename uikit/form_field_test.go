package uikit

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

func TestModernTextBoxThemeShading(t *testing.T) {
	test.NewApp()

	// Creamos un textbox desactivado (modo consulta / readonly)
	tb := NewTextBox("Título", "placeholder", "Paisaje al atardecer", true, nil)
	tb.Resize(fyne.NewSize(200, 36))

	if !tb.Disabled() {
		t.Fatalf("Se esperaba que tb estuviera desactivado")
	}

	th := tb.themedInput.Theme
	if th == nil {
		t.Fatalf("Se esperaba que themedInput tuviera un Theme asignado")
	}

	// 1. Verificar color de fondo cuando está desactivado: sombreado #f8fafc
	bgCol := th.Color(theme.ColorNameInputBackground, theme.VariantLight)
	nrgba, ok := bgCol.(color.NRGBA)
	if !ok {
		t.Fatalf("Se esperaba color.NRGBA, obtenido %T", bgCol)
	}
	if nrgba.R != 248 || nrgba.G != 250 || nrgba.B != 252 {
		t.Errorf("Esperado sombreado #f8fafc (248, 250, 252), obtenido (%d, %d, %d)", nrgba.R, nrgba.G, nrgba.B)
	}

	// 2. Verificar que no tenga marco (StrokeWidth / InputBorder == 0)
	borderSize := th.Size(theme.SizeNameInputBorder)
	if borderSize != 0 {
		t.Errorf("Esperado borderSize 0 cuando está desactivado, obtenido %f", borderSize)
	}

	// 3. Verificar color de texto desactivado: #334155 (Slate 700)
	txtCol := th.Color(theme.ColorNameDisabled, theme.VariantLight)
	txtNrgba, ok := txtCol.(color.NRGBA)
	if !ok {
		t.Fatalf("Se esperaba color.NRGBA para texto desactivado, obtenido %T", txtCol)
	}
	if txtNrgba.R != 51 || txtNrgba.G != 65 || txtNrgba.B != 85 {
		t.Errorf("Esperado texto #334155 (51, 65, 85), obtenido (%d, %d, %d)", txtNrgba.R, txtNrgba.G, txtNrgba.B)
	}

	// 4. Cambiar a activo (habilitado / edición)
	tb.Enable()
	if tb.Disabled() {
		t.Fatalf("Se esperaba que tb estuviera activado")
	}

	// Fondo blanco para edición
	bgColActive := th.Color(theme.ColorNameInputBackground, theme.VariantLight)
	activeNrgba := bgColActive.(color.NRGBA)
	if activeNrgba.R != 255 || activeNrgba.G != 255 || activeNrgba.B != 255 {
		t.Errorf("Esperado fondo blanco #ffffff en modo activo, obtenido (%d, %d, %d)", activeNrgba.R, activeNrgba.G, activeNrgba.B)
	}

	// Borde de 1px en modo activo
	activeBorderSize := th.Size(theme.SizeNameInputBorder)
	if activeBorderSize != 1.0 {
		t.Errorf("Esperado borderSize 1.0 en modo activo, obtenido %f", activeBorderSize)
	}

	// 5. Volver a desactivar
	tb.Disable()
	if !tb.Disabled() {
		t.Fatalf("Se esperaba que tb volviera a estar desactivado")
	}
	if th.Size(theme.SizeNameInputBorder) != 0 {
		t.Errorf("Esperado borderSize 0 tras desactivar de nuevo, obtenido %f", th.Size(theme.SizeNameInputBorder))
	}
	bgColDisabledAgain := th.Color(theme.ColorNameInputBackground, theme.VariantLight).(color.NRGBA)
	if bgColDisabledAgain.R != 248 || bgColDisabledAgain.G != 250 || bgColDisabledAgain.B != 252 {
		t.Errorf("Esperado sombreado #f8fafc tras desactivar de nuevo")
	}

	// 6. Verificar los objetos renderizados dentro de Entry
	r := test.WidgetRenderer(tb.Entry)
	objs := r.Objects()
	if len(objs) < 2 {
		t.Fatalf("Se esperaban al menos 2 objetos en renderer de Entry (box y border), obtenidos %d", len(objs))
	}
	boxRect, okBox := objs[0].(*canvas.Rectangle)
	borderRect, okBorder := objs[1].(*canvas.Rectangle)
	if !okBox || !okBorder {
		t.Fatalf("Los objetos 0 y 1 deben ser *canvas.Rectangle")
	}

	// Comprobamos estado desactivado
	if borderRect.StrokeWidth != 0 {
		t.Errorf("Se esperaba border.StrokeWidth == 0 desactivado, obtenido %f", borderRect.StrokeWidth)
	}
	boxFillDisabled := boxRect.FillColor.(color.NRGBA)
	if boxFillDisabled.R != 248 || boxFillDisabled.G != 250 || boxFillDisabled.B != 252 {
		t.Errorf("Se esperaba box.FillColor sombreado desactivado, obtenido %+v", boxFillDisabled)
	}

	// Comprobamos transición a activado
	tb.Enable()
	if borderRect.StrokeWidth != 1.0 {
		t.Errorf("Se esperaba border.StrokeWidth == 1.0 activado, obtenido %f", borderRect.StrokeWidth)
	}
	boxFillActive := boxRect.FillColor.(color.NRGBA)
	if boxFillActive.R != 255 || boxFillActive.G != 255 || boxFillActive.B != 255 {
		t.Errorf("Se esperaba box.FillColor blanco activado, obtenido %+v", boxFillActive)
	}

	// Comprobamos vuelta a desactivado
	tb.Disable()
	if borderRect.StrokeWidth != 0 {
		t.Errorf("Se esperaba border.StrokeWidth == 0 tras desactivar de nuevo, obtenido %f", borderRect.StrokeWidth)
	}
	boxFillDisabled2 := boxRect.FillColor.(color.NRGBA)
	if boxFillDisabled2.R != 248 || boxFillDisabled2.G != 250 || boxFillDisabled2.B != 252 {
		t.Errorf("Se esperaba box.FillColor sombreado desactivado 2, obtenido %+v", boxFillDisabled2)
	}
}

func TestModernTextBoxThemeDarkMode(t *testing.T) {
	test.NewApp()

	tb := NewTextBox("Título", "placeholder", "Paisaje al atardecer", false, nil)
	th := tb.themedInput.Theme

	// 1. Modo edición activo en tema oscuro: la letra debe ser blanca (#f8fafc) y el fondo oscuro (#0f172a)
	fgDark := th.Color(theme.ColorNameForeground, theme.VariantDark).(color.NRGBA)
	if fgDark.R != 248 || fgDark.G != 250 || fgDark.B != 252 {
		t.Errorf("Esperado texto blanco #f8fafc (248, 250, 252) en modo activo oscuro, obtenido (%d, %d, %d)", fgDark.R, fgDark.G, fgDark.B)
	}

	bgDark := th.Color(theme.ColorNameInputBackground, theme.VariantDark).(color.NRGBA)
	if bgDark.R != 15 || bgDark.G != 23 || bgDark.B != 42 {
		t.Errorf("Esperado fondo oscuro #0f172a en modo activo oscuro, obtenido (%d, %d, %d)", bgDark.R, bgDark.G, bgDark.B)
	}

	// 2. Modo edición activo en tema claro: la letra debe ser negra pizarra (#0f172a) y el fondo blanco
	fgLight := th.Color(theme.ColorNameForeground, theme.VariantLight).(color.NRGBA)
	if fgLight.R != 15 || fgLight.G != 23 || fgLight.B != 42 {
		t.Errorf("Esperado texto oscuro #0f172a en modo activo claro, obtenido (%d, %d, %d)", fgLight.R, fgLight.G, fgLight.B)
	}

	// 3. Modo solo lectura (desactivado) en tema oscuro
	tb.Disable()
	bgDisabledDark := th.Color(theme.ColorNameInputBackground, theme.VariantDark).(color.NRGBA)
	if bgDisabledDark.R != 30 || bgDisabledDark.G != 41 || bgDisabledDark.B != 59 {
		t.Errorf("Esperado fondo #1e293b en modo solo lectura oscuro, obtenido (%d, %d, %d)", bgDisabledDark.R, bgDisabledDark.G, bgDisabledDark.B)
	}

	fgDisabledDark := th.Color(theme.ColorNameForeground, theme.VariantDark).(color.NRGBA)
	if fgDisabledDark.R != 203 || fgDisabledDark.G != 213 || fgDisabledDark.B != 225 {
		t.Errorf("Esperado texto #cbd5e1 en modo solo lectura oscuro, obtenido (%d, %d, %d)", fgDisabledDark.R, fgDisabledDark.G, fgDisabledDark.B)
	}
}

