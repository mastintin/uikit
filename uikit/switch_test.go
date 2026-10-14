package uikit

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
)

func TestModernSwitch_Basic(t *testing.T) {
	test.NewApp()

	// 1. Constructor por defecto
	sw := NewModernSwitch()
	if sw.On {
		t.Errorf("Expected initial state to be false (OFF)")
	}
	if sw.Variant != SwitchPrimary {
		t.Errorf("Expected default variant to be SwitchPrimary")
	}

	// 2. Toggle y OnChanged
	var changedVal bool
	var callbackCalls int
	sw.SetOnChanged(func(on bool) {
		changedVal = on
		callbackCalls++
	})

	sw.Toggle()
	if !sw.On || !sw.Checked() || !sw.IsOn() {
		t.Errorf("Expected switch to be ON after Toggle()")
	}
	if callbackCalls != 1 || !changedVal {
		t.Errorf("Expected callback called once with true, got calls=%d, val=%v", callbackCalls, changedVal)
	}

	sw.SetChecked(false)
	if sw.On {
		t.Errorf("Expected switch to be OFF after SetChecked(false)")
	}
	if callbackCalls != 2 || changedVal {
		t.Errorf("Expected callback called twice with false, got calls=%d, val=%v", callbackCalls, changedVal)
	}
}

func TestModernSwitch_Constructors(t *testing.T) {
	test.NewApp()

	// Constructor xBase polimórfico
	var triggered bool
	sw1 := NewSwitch("Efecto Cristal", true, func(on bool) {
		triggered = on
	})
	if sw1.GetText() != "Efecto Cristal" {
		t.Errorf("Expected text 'Efecto Cristal', got '%s'", sw1.GetText())
	}
	if !sw1.On {
		t.Errorf("Expected switch sw1 to be ON")
	}
	sw1.Toggle()
	if triggered != false {
		t.Errorf("Expected triggered to be false after toggle")
	}

	// Constructores semánticos
	swSucc := NewSuccessSwitch("Aprobado", false, nil)
	if swSucc.Variant != SwitchSuccess {
		t.Errorf("Expected SwitchSuccess variant")
	}

	swDang := NewDangerSwitch("Eliminar", true, nil)
	if swDang.Variant != SwitchDanger {
		t.Errorf("Expected SwitchDanger variant")
	}

	swPrim := NewPrimarySwitch("Activar", true, nil)
	if swPrim.Variant != SwitchPrimary {
		t.Errorf("Expected SwitchPrimary variant")
	}
}

func TestModernSwitch_HitboxAndLabelPositions(t *testing.T) {
	test.NewApp()

	// 1. Probar hitbox delimitada
	sw := NewSwitch("Guardar cambios", false, nil)
	w := test.NewWindow(sw)
	defer w.Close()

	sw.Resize(fyne.NewSize(400, 30))
	sw.Refresh()

	// Clic dentro de la hitbox (en la posición 10, 10, que es sobre el switch)
	sw.Tapped(&fyne.PointEvent{Position: fyne.NewPos(10, 10)})
	if !sw.On {
		t.Errorf("Expected ON when clicking inside switch hitbox")
	}

	// Clic fuera de la hitbox (en 350, 10, margen derecho vacío)
	sw.Tapped(&fyne.PointEvent{Position: fyne.NewPos(350, 10)})
	// Debe mantenerse en ON (no conmutar)
	if !sw.On {
		t.Errorf("Expected state to stay ON when clicking outside hitbox at (350, 10)")
	}

	// 2. Probar LabelLeading y LabelTrailing
	sw.WithLeadingLabel()
	if sw.LabelPosition != LabelLeading {
		t.Errorf("Expected LabelLeading")
	}
	sw.Resize(fyne.NewSize(200, 30))
	sw.Refresh()

	sw.WithTrailingLabel()
	if sw.LabelPosition != LabelTrailing {
		t.Errorf("Expected LabelTrailing")
	}
}

func TestModernSwitch_InteractionShadowAndHalo(t *testing.T) {
	test.NewApp()

	sw := NewSwitch("Accesible", false, nil)
	w := test.NewWindow(sw)
	defer w.Close()
	sw.Resize(fyne.NewSize(120, 30))

	// Cursor dentro de hitbox
	sw.MouseIn(&desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(10, 10)}})
	if sw.Cursor() != desktop.PointerCursor {
		t.Errorf("Expected PointerCursor when hovered inside hitbox, got %v", sw.Cursor())
	}

	// Hover fuera de hitbox
	sw.MouseMoved(&desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(300, 10)}})
	if sw.isHovered {
		t.Errorf("Expected isHovered false when mouse is outside hitbox")
	}
	if sw.Cursor() != desktop.DefaultCursor {
		t.Errorf("Expected DefaultCursor when outside hitbox")
	}

	// Teclado: Espacio
	sw.TypedRune(' ')
	if !sw.On {
		t.Errorf("Expected ON after TypedRune(' ')")
	}

	// Teclado: Enter
	sw.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	if sw.On {
		t.Errorf("Expected OFF after Enter key")
	}

	// Foco y Halo
	sw.FocusGained()
	if !sw.isFocused {
		t.Errorf("Expected isFocused true after FocusGained")
	}
	if !sw.thumbHalo.Visible() {
		t.Errorf("Expected thumbHalo to be visible when focused")
	}

	sw.FocusLost()
	if sw.isFocused {
		t.Errorf("Expected isFocused false after FocusLost")
	}
	if sw.thumbHalo.Visible() {
		t.Errorf("Expected thumbHalo to be hidden after FocusLost")
	}

	// Deshabilitado
	sw.Disable()
	sw.Tapped(&fyne.PointEvent{Position: fyne.NewPos(10, 10)})
	if sw.On {
		t.Errorf("Expected state to remain unchanged when disabled")
	}
}

func TestModernSwitch_RendererObjects(t *testing.T) {
	test.NewApp()

	sw := NewSwitch("Test Switch", true, nil)
	w := test.NewWindow(sw)
	defer w.Close()

	r := test.WidgetRenderer(sw)
	objs := r.Objects()

	// Deben estar: track, thumbHalo, thumbShadow, thumb, label (5 objetos)
	if len(objs) != 5 {
		t.Errorf("Expected 5 renderer objects for labeled switch, got %d", len(objs))
	}

	// Sin texto: deben ser 4 objetos (track, thumbHalo, thumbShadow, thumb)
	swNoText := NewSwitch(false)
	rNoText := test.WidgetRenderer(swNoText)
	if len(rNoText.Objects()) != 4 {
		t.Errorf("Expected 4 renderer objects for switch without text, got %d", len(rNoText.Objects()))
	}
}
