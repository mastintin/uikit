package uikit

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
)

func TestModernIcon(t *testing.T) {
	app := test.NewApp()
	app.Settings().SetTheme(NewAppTheme())

	// 1. Constructor unificado xBase
	ic := NewIcon("search", 24, "#2563eb", nil)
	if ic == nil {
		t.Fatal("expected icon to not be nil")
	}
	if ic.MinSize().Width != 24 || ic.MinSize().Height != 24 {
		t.Errorf("expected MinSize 24x24, got %v", ic.MinSize())
	}
	if ic.Cursor() != desktop.DefaultCursor {
		t.Errorf("expected default cursor for static icon, got %v", ic.Cursor())
	}

	// 2. Icono interactivo / de acción
	clicked := false
	actionIc := NewActionIcon("close", 18, "#94a3b8", "#0f172a", func() {
		clicked = true
	})
	if actionIc.Cursor() != desktop.PointerCursor {
		t.Errorf("expected pointer cursor for tappable icon, got %v", actionIc.Cursor())
	}

	actionIc.Tapped(&fyne.PointEvent{})
	if !clicked {
		t.Error("expected onTap callback to have fired")
	}

	// Probar Hover
	actionIc.MouseIn(&desktop.MouseEvent{})
	if !actionIc.isHovered {
		t.Error("expected isHovered to be true on MouseIn")
	}
	actionIc.MouseOut()
	if actionIc.isHovered {
		t.Error("expected isHovered to be false on MouseOut")
	}

	// 3. Probar Badge Icon
	badgeIc := NewBadgeIcon("palette", 20, 38, "#eff6ff", "#bfdbfe", "#2563eb")
	if badgeIc.MinSize().Width != 38 || badgeIc.MinSize().Height != 38 {
		t.Errorf("expected badge icon size 38x38, got %v", badgeIc.MinSize())
	}
	if !badgeIc.hasBg {
		t.Error("expected badge icon to have background")
	}

	// 4. Probar Fluent API
	ic.WithSize(32).
		WithColor("#dc2626").
		WithHoverColor("#991b1b").
		WithBackground("#fef2f2", 6).
		WithBorder("#fecaca", 1)

	if ic.size != 32 {
		t.Errorf("expected size 32, got %v", ic.size)
	}
	if ic.colorHex != "#dc2626" {
		t.Errorf("expected color #dc2626, got %v", ic.colorHex)
	}

	// 5. Probar Icono Coloreado
	icPalette := NewColoredIcon("palette", 22, "#2563eb")
	if icPalette.size != 22 {
		t.Errorf("expected colored icon size 22, got %v", icPalette.size)
	}
}
