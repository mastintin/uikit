package uikit

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
)

func TestModernLabel_Constructors(t *testing.T) {
	// Base
	l1 := NewModernLabel("Prueba")
	if l1.Text() != "Prueba" {
		t.Errorf("Esperado 'Prueba', obtenido '%s'", l1.Text())
	}
	if l1.textSize != StyleBody.Size {
		t.Errorf("Esperado tamaño %.1f, obtenido %.1f", StyleBody.Size, l1.textSize)
	}

	// Con Style
	l2 := NewModernLabel("Título", StyleHero)
	if l2.textSize != StyleHero.Size || !l2.textStyle.Bold {
		t.Errorf("Estilo Hero no aplicado correctamente")
	}

	// NewLabel (rápido)
	l3 := NewLabel("Texto Simple", ColorMuted)
	if l3.Text() != "Texto Simple" || l3.textColor != ColorMuted {
		t.Errorf("NewLabel con color no asignó correctamente los valores")
	}

	// Sugar constructors
	hero := NewHeroLabel("100.000 €")
	if hero.textSize != StyleHero.Size || !hero.textStyle.Bold {
		t.Errorf("NewHeroLabel falló")
	}

	title := NewTitleLabel("Cabecera")
	if title.textSize != StyleTitle1.Size {
		t.Errorf("NewTitleLabel falló")
	}

	sec := NewSectionLabel("Sección")
	if sec.textSize != StyleSection.Size {
		t.Errorf("NewSectionLabel falló")
	}

	caption := NewCaptionLabel("Detalle")
	if caption.textSize != StyleCaption.Size || caption.textColor != ColorMuted {
		t.Errorf("NewCaptionLabel falló")
	}

	val := NewValueLabel("1.234 €", true, ColorSuccess)
	if val.textSize != StyleValueBold.Size || !val.textStyle.Bold || val.textColor != ColorSuccess {
		t.Errorf("NewValueLabel falló")
	}

	micro := NewMicroLabel("ID-001")
	if micro.textSize != StyleMicro.Size {
		t.Errorf("NewMicroLabel falló")
	}
}

func TestModernLabel_FluentMethods(t *testing.T) {
	lbl := NewModernLabel("Inicial")
	customCol := color.NRGBA{R: 255, G: 0, B: 0, A: 255}

	lbl.SetText("Modificado").
		SetSize(15.5).
		SetColor(customCol).
		SetBold(true).
		SetItalic(true).
		SetMonospace(true).
		AlignCenter().
		WithWrap(true)

	if lbl.Text() != "Modificado" {
		t.Errorf("SetText falló")
	}
	if lbl.textSize != 15.5 {
		t.Errorf("SetSize falló")
	}
	if lbl.textColor != customCol {
		t.Errorf("SetColor falló")
	}
	if !lbl.textStyle.Bold || !lbl.textStyle.Italic || !lbl.textStyle.Monospace {
		t.Errorf("Estilos de texto fallaron")
	}
	if lbl.alignment != fyne.TextAlignCenter {
		t.Errorf("AlignCenter falló")
	}
	if lbl.wrapping != fyne.TextWrapWord {
		t.Errorf("WithWrap falló")
	}

	lbl.AlignLeft()
	if lbl.alignment != fyne.TextAlignLeading {
		t.Errorf("AlignLeft falló")
	}

	lbl.AlignRight()
	if lbl.alignment != fyne.TextAlignTrailing {
		t.Errorf("AlignRight falló")
	}

	lbl.SetHexColor("#34d399")
	expectedCol := HexColor("#34d399")
	if lbl.textColor != expectedCol {
		t.Errorf("SetHexColor falló")
	}
}

func TestModernLabel_Renderer(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	lbl := NewModernLabel("Línea 1\nLínea 2").SetSize(14)
	r := lbl.CreateRenderer()

	objs := r.Objects()
	if len(objs) != 2 {
		t.Fatalf("Esperados 2 objetos de línea, obtenidos %d", len(objs))
	}

	minSize := r.MinSize()
	if minSize.Width <= 0 || minSize.Height <= 0 {
		t.Errorf("MinSize inválido: %+v", minSize)
	}

	r.Layout(fyne.NewSize(200, 50))
	r.Refresh()
}

func TestModernLabel_Tapped(t *testing.T) {
	clicked := false
	lbl := NewModernLabel("Click me").SetOnTapped(func() {
		clicked = true
	})

	if lbl.Cursor() != desktop.PointerCursor {
		t.Errorf("Esperado desktop.PointerCursor, obtenido %v", lbl.Cursor())
	}

	lbl.Tapped(&fyne.PointEvent{})
	if !clicked {
		t.Errorf("Callback OnTapped no fue ejecutado")
	}
}
