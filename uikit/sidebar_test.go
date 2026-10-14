package uikit

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestModernSidebar(t *testing.T) {
	app := test.NewApp()
	app.Settings().SetTheme(NewAppTheme())

	header := widget.NewLabel("Header Area")
	content := widget.NewLabel("Content Area")
	footer := widget.NewLabel("Footer Area")

	// 1. Constructor unificado xBase
	sb := NewSidebarContainer(320, header, content, footer)
	if sb == nil {
		t.Fatal("expected sidebar container to not be nil")
	}

	// 2. Ancho fijo
	if sb.MinSize().Width != 320 {
		t.Errorf("expected MinSize.Width 320, got %v", sb.MinSize().Width)
	}

	// 3. Probar Fluent API
	sb.WithWidth(350).
		WithBackground("#f8fafc").
		WithRightBorder("#cbd5e1", 2).
		WithHeader(widget.NewLabel("New Header"))

	if sb.width != 350 {
		t.Errorf("expected width 350, got %v", sb.width)
	}
	if sb.borderWidth != 2 {
		t.Errorf("expected borderWidth 2, got %v", sb.borderWidth)
	}

	// 4. Probar layout en un contenedor Border (como en la app principal)
	mainArea := container.NewBorder(nil, nil, sb, nil, widget.NewLabel("Detail Area"))
	mainArea.Resize(fyne.NewSize(1200, 700))
	mainArea.Refresh()

	if sb.Size().Width != 350 {
		t.Errorf("expected sidebar allocated width to be 350, got %v", sb.Size().Width)
	}
}
