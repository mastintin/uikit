package uikit

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestModernSidebarNavButton(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	clicked := false
	btn := NewSidebarNavButton("Inicio", "home", func() {
		clicked = true
	})

	if btn == nil {
		t.Fatal("NewSidebarNavButton returned nil")
	}

	if btn.IsActive() {
		t.Errorf("Expected button to be inactive by default")
	}

	btn.SetActive(true)
	if !btn.IsActive() {
		t.Errorf("Expected button to be active after SetActive(true)")
	}

	// Probar interacción de clic
	test.Tap(btn)
	if !clicked {
		t.Errorf("Expected callback to be triggered on tap")
	}

	// Probar dimensiones mínimas
	size := btn.MinSize()
	if size.Width < 160 || size.Height < 40 {
		t.Errorf("Expected MinSize around {160, 42}, got %v", size)
	}

	// Layout y renderizado
	w := app.NewWindow("Test")
	w.SetContent(btn)
	w.Resize(fyne.NewSize(240, 50))
	btn.Refresh()
}
