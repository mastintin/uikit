package uikit

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

func TestModernFileDialogs(t *testing.T) {
	app := test.NewApp()
	app.Settings().SetTheme(NewAppTheme())
	win := app.NewWindow("Test File Dialogs")
	win.Resize(fyne.NewSize(1200, 700))

	// 1. Probar tamaño calculado (intermedio 820x510)
	size := calculateDialogSize(win)
	if size.Width != 820 || size.Height != 510 {
		t.Errorf("expected intermediate dialog size (820x510), got %v", size)
	}

	// 2. Probar creación de ShowSaveFileDialog
	saveTriggered := false
	ShowSaveFileDialog(win, "Guardar Datos", "datos.csv", "csv", func(p string) {
		saveTriggered = true
	})
	if saveTriggered {
		t.Error("save callback should not trigger immediately")
	}

	topOverlay := win.Canvas().Overlays().Top()
	if topOverlay == nil {
		t.Fatal("expected modal overlay for save dialog")
	}

	// 3. Probar creación de ShowOpenFileDialog
	openTriggered := false
	ShowOpenFileDialog(win, "Abrir Imagen", []string{"jpg", "png"}, func(p string) {
		openTriggered = true
	})
	if openTriggered {
		t.Error("open callback should not trigger immediately")
	}

	// 4. Probar creación de ShowFolderDialog
	folderTriggered := false
	ShowFolderDialog(win, "Seleccionar Carpeta", func(p string) {
		folderTriggered = true
	})
	if folderTriggered {
		t.Error("folder callback should not trigger immediately")
	}

	// 5. Probar creación de ShowFolderDialogWithInitial
	initialFolderTriggered := false
	ShowFolderDialogWithInitial(win, "Seleccionar Carpeta Fotos", "/tmp", func(p string) {
		initialFolderTriggered = true
	})
	if initialFolderTriggered {
		t.Error("initial folder callback should not trigger immediately")
	}

	// 6. Probar creación de ShowOpenFileDialogWithInitial
	initialFileTriggered := false
	ShowOpenFileDialogWithInitial(win, "Seleccionar DB", []string{"db"}, "/tmp/cuadros.db", func(p string) {
		initialFileTriggered = true
	})
	if initialFileTriggered {
		t.Error("initial file callback should not trigger immediately")
	}

	// 7. Probar ShowSavedSuccessModal
	revealed := false
	ShowSavedSuccessModal(win, "Guardado con éxito", "Archivo exportado correctamente", "/tmp/test.csv", func(p string) {
		revealed = true
	})
	if revealed {
		t.Error("reveal callback should not trigger immediately")
	}
}

func TestModernFileDialog_DarkTheme(t *testing.T) {
	app := test.NewApp()
	darkTheme := NewAppTheme().SetForcedVariant(theme.VariantDark)
	app.Settings().SetTheme(darkTheme)
	win := app.NewWindow("Test Dark File Dialog")
	win.Resize(fyne.NewSize(1200, 700))

	ShowSaveFileDialog(win, "Exportar Auditoría a Excel", "auditoria.xlsx", "xlsx", func(p string) {})

	topOverlay := win.Canvas().Overlays().Top()
	if topOverlay == nil {
		t.Fatal("expected modal overlay for save dialog in dark theme")
	}

	// Verify that the dialog card does not use hardcoded white (#ffffff)
	cardBg := ThemeColor(ColorNameCardBg)
	nrgba, ok := cardBg.(color.NRGBA)
	if !ok || nrgba.R > 50 || nrgba.G > 50 || nrgba.B > 50 {
		t.Errorf("expected dark theme card background, got %v", cardBg)
	}
}
