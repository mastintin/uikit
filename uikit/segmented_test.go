package uikit

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestModernSegmentedControl_Basic(t *testing.T) {
	items := []SegmentItem{
		{ID: "nice", Label: "Nice", Icon: "palette"},
		{ID: "dark", Label: "Oscuro", Icon: "moon"},
		{ID: "light", Label: "Claro", Icon: "sun"},
	}

	var triggeredIdx int = -99
	var triggeredID string

	seg := NewSegmentedControl(items, 0, func(idx int, item SegmentItem) {
		triggeredIdx = idx
		triggeredID = item.ID
	})

	if seg.SelectedIndex() != 0 {
		t.Errorf("Esperado selectedIdx 0, obtenido %d", seg.SelectedIndex())
	}
	if seg.SelectedID() != "nice" {
		t.Errorf("Esperado selectedID 'nice', obtenido '%s'", seg.SelectedID())
	}

	// Cambiar a Oscuro
	seg.SetSelected(1)
	if triggeredIdx != 1 || triggeredID != "dark" {
		t.Errorf("Esperado callback disparado con idx=1 y id='dark', obtenido %d y '%s'", triggeredIdx, triggeredID)
	}
	if seg.SelectedID() != "dark" {
		t.Errorf("Esperado selectedID 'dark', obtenido '%s'", seg.SelectedID())
	}

	// Cambiar a Claro mediante ID
	seg.SetSelectedByID("light")
	if triggeredIdx != 2 || triggeredID != "light" {
		t.Errorf("Esperado callback disparado con idx=2 y id='light', obtenido %d y '%s'", triggeredIdx, triggeredID)
	}

	// Modo silencioso (sin disparar callback)
	triggeredIdx = -99
	seg.SetSelectedSilent(0)
	if triggeredIdx != -99 {
		t.Errorf("SetSelectedSilent no debería disparar callback, pero se disparó con idx=%d", triggeredIdx)
	}
	if seg.SelectedID() != "nice" {
		t.Errorf("Esperado 'nice' tras SetSelectedSilent, obtenido '%s'", seg.SelectedID())
	}

	// ClearSelection
	seg.ClearSelection()
	if seg.SelectedIndex() != -1 || seg.SelectedID() != "" {
		t.Errorf("Esperado ningún elemento seleccionado tras ClearSelection, obtenido idx=%d id='%s'", seg.SelectedIndex(), seg.SelectedID())
	}
}

func TestModernSegmentedControl_TapAndHover(t *testing.T) {
	items := []SegmentItem{
		{ID: "opt1", Label: "Uno"},
		{ID: "opt2", Label: "Dos"},
	}

	var selectedItem string
	seg := NewSegmentedControl(items, 0, func(idx int, item SegmentItem) {
		selectedItem = item.ID
	})

	w := test.NewWindow(seg)
	defer w.Close()

	// Tap en el segundo segmento
	test.Tap(seg.segments[1])
	if selectedItem != "opt2" {
		t.Errorf("Esperado tap en segmento 1 para seleccionar 'opt2', obtenido '%s'", selectedItem)
	}

	// Hover en el primer segmento
	seg.segments[0].MouseIn(nil)
	if !seg.segments[0].isHovered {
		t.Errorf("Esperado isHovered=true tras MouseIn")
	}
	seg.segments[0].MouseOut()
	if seg.segments[0].isHovered {
		t.Errorf("Esperado isHovered=false tras MouseOut")
	}
}
