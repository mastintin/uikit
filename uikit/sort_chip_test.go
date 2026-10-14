package uikit

import (
	"testing"

	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
)

func TestModernSortChip_Standalone(t *testing.T) {
	test.NewApp()

	columns := []string{"Periodo", "Base", "Estado"}
	var lastCol string
	var lastDir SortDirection
	calls := 0

	chip := NewSortChip(columns, "Periodo", SortAscending, func(col string, dir SortDirection) {
		lastCol = col
		lastDir = dir
		calls++
	})

	if !chip.IsAtDefault() {
		t.Errorf("Expected chip to start at default")
	}

	opt, dir := chip.CurrentSort()
	if opt.Label != "Periodo" || dir != SortAscending {
		t.Errorf("Expected Periodo Ascending, got %s %v", opt.Label, dir)
	}

	// 1. Cambiar a Base Descendente
	chip.SetSort(1, SortDescending)
	if chip.IsAtDefault() {
		t.Errorf("Expected chip NOT to be at default after change")
	}
	if calls != 1 || lastCol != "Base" || lastDir != SortDescending {
		t.Errorf("Expected callback with Base Descending, got %s %v (calls=%d)", lastCol, lastDir, calls)
	}

	// 2. Restablecer orden
	chip.Reset()
	if !chip.IsAtDefault() {
		t.Errorf("Expected chip to return to default after Reset()")
	}
	if calls != 2 || lastCol != "Periodo" || lastDir != SortAscending {
		t.Errorf("Expected callback reset to Periodo Ascending, got %s %v", lastCol, lastDir)
	}
}

func TestModernSortChip_BrowseIntegration(t *testing.T) {
	test.NewApp()

	columns := []BrowseColumn{
		{Title: "AÑO", Align: BrowseAlignCenter},
		{Title: "IMPORTE", Align: BrowseAlignTrailing},
		{Title: "ACCIONES", Align: BrowseAlignCenter, NoSort: true},
	}

	browse := NewModernBrowse(columns)
	chip := browse.CreateSortChip()

	if chip == nil {
		t.Fatal("Expected CreateSortChip to return a valid ModernSortChip")
	}

	// Debe excluir la columna NoSort ("ACCIONES")
	if len(chip.options) != 2 {
		t.Fatalf("Expected 2 sortable options, got %d", len(chip.options))
	}
	if chip.options[0].Label != "AÑO" || chip.options[1].Label != "IMPORTE" {
		t.Errorf("Unexpected chip options: %+v", chip.options)
	}

	// 1. Cambiar orden desde el chip -> debe actualizar el browse
	chip.SetSort(1, SortDescending) // IMPORTE Descendente
	bCol, bDir := browse.CurrentSort()
	if bCol != 1 || bDir != SortDescending {
		t.Errorf("Expected browse to sort by col 1 Descending, got col=%d dir=%v", bCol, bDir)
	}

	// 2. Sincronización inversa: hacer ToggleSort en el browse -> debe actualizar el chip
	browse.ToggleSort(0) // AÑO Ascendente
	chipOpt, chipDir := chip.CurrentSort()
	if chipOpt.Index != 0 || chipDir != SortAscending {
		t.Errorf("Expected chip to reflect browse ToggleSort(0), got opt=%+v dir=%v", chipOpt, chipDir)
	}
}

func TestModernSortChip_InteractionAndDesktop(t *testing.T) {
	test.NewApp()

	chip := NewSortChip([]string{"Col1", "Col2"}, "Col1", SortAscending, nil)

	// Cursor
	if chip.Cursor() != desktop.PointerCursor {
		t.Errorf("Expected PointerCursor, got %v", chip.Cursor())
	}

	// Hover
	chip.MouseIn(&desktop.MouseEvent{})
	if !chip.isHovered {
		t.Errorf("Expected isHovered true after MouseIn")
	}

	chip.MouseOut()
	if chip.isHovered {
		t.Errorf("Expected isHovered false after MouseOut")
	}

	// Renderer
	w := test.NewWindow(chip)
	defer w.Close()

	minSize := chip.MinSize()
	if minSize.Height < 28 {
		t.Errorf("Expected minSize.Height >= 28, got %f", minSize.Height)
	}
}
