package uikit

import (
	"image"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestImagePreviewFrame(t *testing.T) {
	tapped := false
	frame := NewImagePreviewFrame(func() {
		tapped = true
	})

	if frame == nil {
		t.Fatal("expected frame to not be nil")
	}

	// 1. Verificar esquinas redondeadas y colores iniciales
	if frame.bg.CornerRadius != 8 {
		t.Errorf("expected CornerRadius 8, got %v", frame.bg.CornerRadius)
	}

	// 2. Verificar Cursor Pointer
	if frame.Cursor() != desktop.PointerCursor {
		t.Errorf("expected desktop.PointerCursor, got %v", frame.Cursor())
	}

	// 3. Simular clic / Tapped
	test.Tap(frame)
	if !tapped {
		t.Error("expected onTap callback to be called on tap")
	}

	// 4. Probar hover (MouseIn / MouseOut)
	frame.MouseIn(&desktop.MouseEvent{})
	if !frame.isHovered {
		t.Error("expected isHovered to be true after MouseIn")
	}
	if frame.bg.StrokeColor != frame.hoverBorderColor {
		t.Error("expected stroke color to change to hoverBorderColor")
	}

	frame.MouseOut()
	if frame.isHovered {
		t.Error("expected isHovered to be false after MouseOut")
	}
	if frame.bg.StrokeColor != frame.normalBorderColor {
		t.Error("expected stroke color to restore to normalBorderColor")
	}

	// 5. Probar SetImage y SetFile (archivo inexistente -> placeholder)
	frame.SetFile("/non/existent/file.jpg")
	if frame.img.Image == nil {
		t.Error("expected placeholder image when file does not exist")
	}

	m := image.NewRGBA(image.Rect(0, 0, 10, 10))
	frame.SetImage(m)
	if frame.img.Image != m {
		t.Error("expected set image to match")
	}

	// 6. Verificar que MinSize devuelve dimensiones válidas
	min := frame.MinSize()
	if min != fyne.NewSize(180, 180) {
		t.Errorf("expected MinSize (180, 180), got %v", min)
	}

	// 7. Probar comportamiento con NewCardPanel, NewBorder y NewVScroll
	bottomControls := container.NewVBox(widget.NewLabel("Buttons"), widget.NewLabel("Path"))
	imageColumn := container.NewBorder(nil, bottomControls, nil, nil, frame)
	card := NewCardPanel(imageColumn)

	card.Resize(fyne.NewSize(400, 600))
	card.Refresh()

	// 8. Probar dentro de NewVScroll
	vscroll := container.NewVScroll(card)
	vscroll.Resize(fyne.NewSize(400, 600))
	vscroll.Refresh()
	t.Logf("Inside VScroll - Frame size: %v", frame.Size())

	// 9. Probar Fluent API y constructor unificado NewImageFrame
	custom := NewImageFrame("", 12, "#1e293b", "#475569", nil).
		WithRadius(14).
		WithBackground("#000000").
		WithBorder("#111111", "#222222", 2).
		WithMinSize(fyne.NewSize(220, 220))

	if custom.cornerRadius != 14 {
		t.Errorf("expected cornerRadius 14, got %v", custom.cornerRadius)
	}
	if custom.MinSize() != fyne.NewSize(220, 220) {
		t.Errorf("expected minSize (220, 220), got %v", custom.MinSize())
	}
}
