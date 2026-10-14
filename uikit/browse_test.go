package uikit

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestModernBrowse_Sorting(t *testing.T) {
	test.NewApp()

	cols := []BrowseColumn{
		{Title: "MESES", Align: BrowseAlignCenter},
		{Title: "PENSIÓN", Align: BrowseAlignTrailing},
		{Title: "ACCIONES", Align: BrowseAlignCenter, NoSort: true},
	}

	browse := NewModernBrowse(cols)

	row1 := BrowseRow{
		Cells: []BrowseCell{
			CellText("24 meses"),
			CellText("1.500,00 €"),
			CellText("Editar"),
		},
	}
	row2 := BrowseRow{
		Cells: []BrowseCell{
			CellText("6 meses"),
			CellText("2.200,50 €"),
			CellText("Editar"),
		},
	}
	row3 := BrowseRow{
		Cells: []BrowseCell{
			CellText("12 meses"),
			CellText("1.800,00 €"),
			CellText("Editar"),
		},
	}

	browse.SetRows([]BrowseRow{row1, row2, row3})

	// Inicialmente sin ordenación
	initialRows := browse.getDisplayRows()
	if initialRows[0].Cells[0].Text != "24 meses" || initialRows[1].Cells[0].Text != "6 meses" {
		t.Fatalf("expected initial order preserved")
	}

	// 1. Click en columna 0 (MESES): Orden Ascendente (numérico: 6, 12, 24)
	browse.ToggleSort(0)
	if browse.sortDir != SortAscending || browse.sortColIdx != 0 {
		t.Fatalf("expected col 0 ascending, got col=%d dir=%v", browse.sortColIdx, browse.sortDir)
	}
	ascRows := browse.getDisplayRows()
	if ascRows[0].Cells[0].Text != "6 meses" || ascRows[1].Cells[0].Text != "12 meses" || ascRows[2].Cells[0].Text != "24 meses" {
		t.Errorf("expected numeric sort asc [6, 12, 24], got [%s, %s, %s]",
			ascRows[0].Cells[0].Text, ascRows[1].Cells[0].Text, ascRows[2].Cells[0].Text)
	}

	// 2. Segundo click en columna 0 (MESES): Orden Descendente (24, 12, 6)
	browse.ToggleSort(0)
	if browse.sortDir != SortDescending || browse.sortColIdx != 0 {
		t.Fatalf("expected col 0 descending")
	}
	descRows := browse.getDisplayRows()
	if descRows[0].Cells[0].Text != "24 meses" || descRows[1].Cells[0].Text != "12 meses" || descRows[2].Cells[0].Text != "6 meses" {
		t.Errorf("expected numeric sort desc [24, 12, 6], got [%s, %s, %s]",
			descRows[0].Cells[0].Text, descRows[1].Cells[0].Text, descRows[2].Cells[0].Text)
	}

	// 3. Click en columna 1 (PENSIÓN monetaria en €): Orden Ascendente (1.500, 1.800, 2.200,50)
	browse.ToggleSort(1)
	if browse.sortDir != SortAscending || browse.sortColIdx != 1 {
		t.Fatalf("expected col 1 ascending")
	}
	pensionAsc := browse.getDisplayRows()
	if pensionAsc[0].Cells[1].Text != "1.500,00 €" || pensionAsc[1].Cells[1].Text != "1.800,00 €" || pensionAsc[2].Cells[1].Text != "2.200,50 €" {
		t.Errorf("expected currency sort asc [1500, 1800, 2200], got [%s, %s, %s]",
			pensionAsc[0].Cells[1].Text, pensionAsc[1].Cells[1].Text, pensionAsc[2].Cells[1].Text)
	}

	// 4. Click en columna 2 (ACCIONES con NoSort: true) no debe alterar la ordenación
	browse.ToggleSort(2)
	if browse.sortColIdx != 1 {
		t.Errorf("NoSort column should not change sort index, got %d", browse.sortColIdx)
	}

	// 5. Test interacción directa con widget de cabecera browseHeaderCell
	cellWidget := newBrowseHeaderCell("Columna", BrowseAlignCenter, SortNone, true, func() {
		browse.ToggleSort(0)
	})
	test.Tap(cellWidget)
	if browse.sortColIdx != 0 {
		t.Errorf("expected tap on browseHeaderCell to trigger sorting")
	}

	var _ fyne.CanvasObject = browse.Container()
}

func TestModernBrowse_StickyHeaderAndScroll(t *testing.T) {
	test.NewApp()

	cols := []BrowseColumn{
		{Title: "COL1", Align: BrowseAlignCenter},
		{Title: "COL2", Align: BrowseAlignTrailing},
	}

	browse := NewModernBrowse(cols)

	// Añadir 30 filas
	for i := 1; i <= 30; i++ {
		browse.AddRow(CellText("A"), CellText("B"))
	}

	// Por defecto maxHeight = 440
	minSizeDefault := browse.Container().MinSize()
	headerH := browse.headerContainer.MinSize().Height

	// La altura total debe estar acotada por maxHeight (440 + headerH)
	if minSizeDefault.Height > 440+headerH+10 {
		t.Errorf("expected default height bounded around 440+headerH, got %f", minSizeDefault.Height)
	}

	// Ajustar maxHeight a 300
	browse.SetMaxHeight(300)
	minSize300 := browse.Container().MinSize()
	if minSize300.Height > 300+headerH+10 {
		t.Errorf("expected height bounded around 300+headerH, got %f", minSize300.Height)
	}

	if browse.scrollRows == nil {
		t.Fatalf("expected scrollRows to be non-nil")
	}
}
