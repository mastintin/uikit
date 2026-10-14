package uikit

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

type testProduct struct {
	ID    string
	Name  string
	Price float64
}

func TestModernListBox(t *testing.T) {
	app := test.NewApp()
	app.Settings().SetTheme(NewAppTheme())

	// 1. Probar NewStringListBox con cadenas simples
	var selectedString string
	strList := NewStringListBox([]string{"Oleo", "Acuarela", "Grabado"}, func(val string) {
		selectedString = val
	})

	if strList.Length() != 3 {
		t.Errorf("expected length 3, got %d", strList.Length())
	}

	strList.Select(1)
	if selectedString != "Acuarela" {
		t.Errorf("expected 'Acuarela', got '%s'", selectedString)
	}

	item, ok := strList.SelectedItem()
	if !ok || item != "Acuarela" {
		t.Errorf("expected SelectedItem 'Acuarela', got '%s'", item)
	}
	if strList.SelectedIndex() != 1 {
		t.Errorf("expected SelectedIndex 1, got %d", strList.SelectedIndex())
	}

	// 2. Probar NewListBox con estructura fuertemente tipada y genéricos
	products := []testProduct{
		{ID: "P1", Name: "Lienzo 50x50", Price: 15.50},
		{ID: "P2", Name: "Pintura Azul Ultramar", Price: 8.20},
		{ID: "P3", Name: "Caballete de Madera", Price: 45.00},
	}

	var selectedProduct testProduct
	prodList := NewListBox[testProduct](
		products,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(p testProduct, o fyne.CanvasObject) {
			if lbl, ok := o.(*widget.Label); ok {
				lbl.SetText(p.Name)
			}
		},
		func(p testProduct) {
			selectedProduct = p
		},
	)

	// 3. Probar SelectBy con predicado
	found := prodList.SelectBy(func(p testProduct) bool {
		return p.ID == "P3"
	})
	if !found {
		t.Error("expected SelectBy to find product P3")
	}
	if selectedProduct.Name != "Caballete de Madera" {
		t.Errorf("expected 'Caballete de Madera', got '%s'", selectedProduct.Name)
	}
	if prodList.SelectedIndex() != 2 {
		t.Errorf("expected SelectedIndex 2, got %d", prodList.SelectedIndex())
	}

	// 4. Probar límites seguros (índice fuera de rango)
	prodList.Select(999)
	if prodList.SelectedIndex() != 2 {
		t.Errorf("out of range select should not change index, got %d", prodList.SelectedIndex())
	}

	// 5. Probar Clear
	prodList.Clear()
	if prodList.Length() != 0 {
		t.Errorf("expected length 0 after Clear, got %d", prodList.Length())
	}
	_, hasItem := prodList.SelectedItem()
	if hasItem {
		t.Error("expected no selected item after Clear")
	}

	// 6. Probar Fluent API
	prodList.WithItems(products)
	if prodList.Length() != 3 {
		t.Errorf("expected length 3 after WithItems, got %d", prodList.Length())
	}
}
