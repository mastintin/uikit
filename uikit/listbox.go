package uikit

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ModernListBox es un control de lista genérico y fuertemente tipado (utilizando genéricos [T any]),
// diseñado bajo la filosofía común de ModernButton, ModernCard, ModernSidebar y ModernTextBox:
// - Extiende widget.BaseWidget como ciudadano nativo del árbol de Fyne.
// - Aprovecha el motor de reciclado virtual de alta velocidad de widget.List.
// - Elimina la fontanería repetitiva (manejo manual de índices, conversiones interface{} y comprobaciones de límites).
// - Proporciona métodos de alto nivel: .SetItems(), .SelectBy(), .SelectedItem(), .SelectedIndex(), etc.
type ModernListBox[T any] struct {
	widget.BaseWidget

	items         []T
	selectedIndex int

	createItem func() fyne.CanvasObject
	updateItem func(item T, obj fyne.CanvasObject)
	onSelected func(item T)

	list      *widget.List
	container *fyne.Container
}

// Alias de tipo para mayor comodidad
type ListBox[T any] = ModernListBox[T]

// NewListBox es el constructor unificado estilo xBase / Harbour para cualquier tipo de datos.
// - items: colección inicial de elementos tipados
// - createItem: función constructora de plantilla visual para cada fila
// - updateItem: función actualizadora que inyecta los datos de tipo T en la fila reciclada
// - onSelected: callback opcional que recibe directamente el elemento seleccionado de tipo T
func NewListBox[T any](
	items []T,
	createItem func() fyne.CanvasObject,
	updateItem func(item T, obj fyne.CanvasObject),
	onSelected func(item T),
) *ModernListBox[T] {
	lb := &ModernListBox[T]{
		items:         items,
		selectedIndex: -1,
		createItem:    createItem,
		updateItem:    updateItem,
		onSelected:    onSelected,
	}

	lb.list = widget.NewList(
		func() int {
			return len(lb.items)
		},
		func() fyne.CanvasObject {
			if lb.createItem != nil {
				return lb.createItem()
			}
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			if id < 0 || id >= len(lb.items) {
				return
			}
			if lb.updateItem != nil {
				lb.updateItem(lb.items[id], o)
			}
		},
	)

	lb.list.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(lb.items) {
			lb.selectedIndex = id
			if lb.onSelected != nil {
				lb.onSelected(lb.items[id])
			}
		}
	}

	lb.list.OnUnselected = func(id widget.ListItemID) {
		if lb.selectedIndex == id {
			lb.selectedIndex = -1
		}
	}

	lb.list.HideSeparators = true
	lb.container = container.NewStack(lb.list)
	lb.ExtendBaseWidget(lb)
	return lb
}

// NewCustomListBox crea una lista tipada especificando únicamente sus funciones de plantilla.
func NewCustomListBox[T any](
	createItem func() fyne.CanvasObject,
	updateItem func(item T, obj fyne.CanvasObject),
) *ModernListBox[T] {
	return NewListBox[T](nil, createItem, updateItem, nil)
}

// NewStringListBox es un constructor semántico ("Sugar Constructor") para listas de cadenas de texto.
func NewStringListBox(items []string, onSelected func(selected string)) *ModernListBox[string] {
	createFn := func() fyne.CanvasObject {
		lbl := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{})
		lbl.Truncation = fyne.TextTruncateEllipsis
		return container.NewPadded(lbl)
	}

	updateFn := func(item string, o fyne.CanvasObject) {
		if padded, ok := o.(*fyne.Container); ok && len(padded.Objects) > 0 {
			if lbl, ok := padded.Objects[0].(*widget.Label); ok {
				lbl.SetText(item)
			}
		}
	}

	lb := NewListBox[string](items, createFn, updateFn, onSelected)
	lb.list.HideSeparators = false
	return lb
}

func (l *ModernListBox[T]) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(l.container)
}

func (l *ModernListBox[T]) MinSize() fyne.Size {
	if l.list != nil {
		return l.list.MinSize()
	}
	return fyne.NewSize(100, 100)
}

// --- Métodos de Alto Nivel y Datos ---

// SetItems reemplaza la lista completa de elementos y refresca el control de forma segura.
func (l *ModernListBox[T]) SetItems(items []T) {
	l.items = items
	if l.selectedIndex >= len(items) {
		l.selectedIndex = -1
	}
	l.list.Refresh()
}

// Items devuelve la colección actual de elementos.
func (l *ModernListBox[T]) Items() []T {
	return l.items
}

// Length devuelve la cantidad de elementos en la lista.
func (l *ModernListBox[T]) Length() int {
	return len(l.items)
}

// Select selecciona el elemento en la posición dada (0-indexed) de forma segura.
func (l *ModernListBox[T]) Select(index int) {
	if index >= 0 && index < len(l.items) {
		l.selectedIndex = index
		l.list.Select(index)
		if index == 0 {
			l.list.ScrollToTop()
		}
	}
}

// ScrollToTop desplaza la lista al inicio absoluto (offset 0).
func (l *ModernListBox[T]) ScrollToTop() {
	if l.list != nil {
		l.list.ScrollToTop()
	}
}

// Unselect deselecciona el elemento en la posición dada.
func (l *ModernListBox[T]) Unselect(index int) {
	if index >= 0 && index < len(l.items) {
		l.list.Unselect(index)
		if l.selectedIndex == index {
			l.selectedIndex = -1
		}
	}
}

// UnselectAll deselecciona cualquier elemento seleccionado actualmente.
func (l *ModernListBox[T]) UnselectAll() {
	if l.selectedIndex >= 0 {
		l.list.UnselectAll()
		l.selectedIndex = -1
	}
}

// SelectBy busca el primer elemento que cumpla la condición del predicado y lo selecciona.
// Devuelve true si encontró y seleccionó el elemento, false en caso contrario.
func (l *ModernListBox[T]) SelectBy(predicate func(item T) bool) bool {
	if predicate == nil {
		return false
	}
	for i, item := range l.items {
		if predicate(item) {
			l.Select(i)
			return true
		}
	}
	return false
}

// SelectedIndex devuelve el índice del elemento seleccionado actualmente (-1 si no hay ninguno).
func (l *ModernListBox[T]) SelectedIndex() int {
	return l.selectedIndex
}

// SelectedItem devuelve el elemento seleccionado actualmente y un booleano indicando si existe.
func (l *ModernListBox[T]) SelectedItem() (T, bool) {
	if l.selectedIndex >= 0 && l.selectedIndex < len(l.items) {
		return l.items[l.selectedIndex], true
	}
	var zero T
	return zero, false
}

// Clear vacía la lista de elementos.
func (l *ModernListBox[T]) Clear() {
	l.SetItems(nil)
}

// Refresh solicita el repintado de la lista.
func (l *ModernListBox[T]) Refresh() {
	if l.list != nil {
		l.list.Refresh()
	}
	l.BaseWidget.Refresh()
}

// SetOnSelected actualiza el callback de selección.
func (l *ModernListBox[T]) SetOnSelected(fn func(item T)) {
	l.onSelected = fn
}

// UnderlyingList expone el widget.List interno de Fyne si se requiere acceso de bajo nivel.
func (l *ModernListBox[T]) UnderlyingList() *widget.List {
	return l.list
}

// --- Interfaz Fluida (Fluent API) ---

func (l *ModernListBox[T]) WithItems(items []T) *ModernListBox[T] {
	l.SetItems(items)
	return l
}

func (l *ModernListBox[T]) WithOnSelected(fn func(item T)) *ModernListBox[T] {
	l.SetOnSelected(fn)
	return l
}

func (l *ModernListBox[T]) WithSeparators(show bool) *ModernListBox[T] {
	if l.list != nil {
		l.list.HideSeparators = !show
		l.list.Refresh()
	}
	return l
}
