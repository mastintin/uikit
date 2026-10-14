package uikit

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// SegmentItem representa una pestaña individual dentro de un ModernSegmentedControl.
type SegmentItem struct {
	ID    string // Identificador único (ej. "nice", "dark", "light")
	Label string // Texto visible corto (ej. "Nice", "Oscuro", "Claro")
	Icon  string // Icono opcional (ej. "palette", "moon", "sun")
}

// ModernSegmentedControl es un control visual interactivo estilo píldora continua
// (Segmented Control nativo de macOS / iOS).
type ModernSegmentedControl struct {
	widget.BaseWidget

	items       []SegmentItem
	selectedIdx int
	onChanged   func(index int, item SegmentItem)

	segments  []*modernSegmentBtn
	container *fyne.Container
	trackBg   *canvas.Rectangle
}

var _ fyne.Widget = (*ModernSegmentedControl)(nil)

// NewSegmentedControl crea un ModernSegmentedControl a partir de una lista de SegmentItems.
func NewSegmentedControl(items []SegmentItem, defaultIdx int, onChanged func(index int, item SegmentItem)) *ModernSegmentedControl {
	if defaultIdx < -1 || defaultIdx >= len(items) {
		defaultIdx = 0
	}

	sc := &ModernSegmentedControl{
		items:       items,
		selectedIdx: defaultIdx,
		onChanged:   onChanged,
		segments:    make([]*modernSegmentBtn, len(items)),
	}

	sc.trackBg = canvas.NewRectangle(color.Transparent)
	sc.trackBg.CornerRadius = 7
	sc.trackBg.StrokeWidth = 1

	segmentObjects := make([]fyne.CanvasObject, len(items))
	for i, item := range items {
		btn := newSegmentBtn(sc, item, i)
		sc.segments[i] = btn
		segmentObjects[i] = btn
	}

	innerBox := container.NewHBox(segmentObjects...)
	paddedBox := container.NewBorder(
		browseSpacer(1, 2),
		browseSpacer(1, 2),
		browseSpacer(2, 1),
		browseSpacer(2, 1),
		innerBox,
	)

	sc.container = container.NewStack(
		sc.trackBg,
		paddedBox,
	)

	sc.ExtendBaseWidget(sc)
	sc.Refresh()
	return sc
}

func (sc *ModernSegmentedControl) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(sc.container)
}

// Refresh actualiza los colores del control y de todos sus segmentos según el tema activo.
func (sc *ModernSegmentedControl) Refresh() {
	sc.trackBg.FillColor = ThemeColor(ColorNameInsetPanelBg)
	sc.trackBg.StrokeColor = ThemeColor(ColorNameBorder)
	sc.trackBg.Refresh()

	for _, seg := range sc.segments {
		if seg != nil {
			seg.Refresh()
		}
	}

	if sc.container != nil {
		sc.container.Refresh()
	}
	sc.BaseWidget.Refresh()
}

// SetSelected cambia la pestaña activa por índice disparando el callback onChanged.
func (sc *ModernSegmentedControl) SetSelected(index int) *ModernSegmentedControl {
	if index < -1 || index >= len(sc.items) {
		return sc
	}
	if sc.selectedIdx == index {
		return sc
	}

	sc.selectedIdx = index
	sc.Refresh()

	if sc.onChanged != nil && index >= 0 && index < len(sc.items) {
		sc.onChanged(index, sc.items[index])
	}
	return sc
}

// SetSelectedByID cambia la pestaña activa por su ID disparando el callback onChanged.
func (sc *ModernSegmentedControl) SetSelectedByID(id string) *ModernSegmentedControl {
	for i, it := range sc.items {
		if strings.EqualFold(it.ID, id) {
			return sc.SetSelected(i)
		}
	}
	return sc.ClearSelection()
}

// SetSelectedSilent actualiza la pestaña activa sin disparar el callback.
func (sc *ModernSegmentedControl) SetSelectedSilent(index int) *ModernSegmentedControl {
	if index < -1 || index >= len(sc.items) {
		return sc
	}
	sc.selectedIdx = index
	sc.Refresh()
	return sc
}

// SetSelectedSilentByID busca y activa la pestaña por ID sin disparar el callback.
func (sc *ModernSegmentedControl) SetSelectedSilentByID(id string) *ModernSegmentedControl {
	for i, it := range sc.items {
		if strings.EqualFold(it.ID, id) {
			return sc.SetSelectedSilent(i)
		}
	}
	return sc.SetSelectedSilent(-1)
}

// ClearSelection desmarca cualquier segmento activo.
func (sc *ModernSegmentedControl) ClearSelection() *ModernSegmentedControl {
	sc.selectedIdx = -1
	sc.Refresh()
	return sc
}

// SelectedIndex devuelve el índice de la opción seleccionada (-1 si ninguna).
func (sc *ModernSegmentedControl) SelectedIndex() int {
	return sc.selectedIdx
}

// SelectedID devuelve el ID de la opción seleccionada ("" si ninguna).
func (sc *ModernSegmentedControl) SelectedID() string {
	if sc.selectedIdx >= 0 && sc.selectedIdx < len(sc.items) {
		return sc.items[sc.selectedIdx].ID
	}
	return ""
}

// ─────────────────────────────────────────────────────────────────────────────
// Botón Interno de Segmento
// ─────────────────────────────────────────────────────────────────────────────

type modernSegmentBtn struct {
	widget.BaseWidget

	parent    *ModernSegmentedControl
	item      SegmentItem
	index     int
	isHovered bool

	bg      *canvas.Rectangle
	icon    *ModernIcon
	txt     *canvas.Text
	content *fyne.Container
}

var _ desktop.Hoverable = (*modernSegmentBtn)(nil)
var _ fyne.Tappable = (*modernSegmentBtn)(nil)
var _ fyne.Widget = (*modernSegmentBtn)(nil)

func newSegmentBtn(parent *ModernSegmentedControl, item SegmentItem, index int) *modernSegmentBtn {
	btn := &modernSegmentBtn{
		parent: parent,
		item:   item,
		index:  index,
	}

	btn.bg = canvas.NewRectangle(color.Transparent)
	btn.bg.CornerRadius = 5
	btn.bg.StrokeWidth = 1

	btn.txt = canvas.NewText(item.Label, color.Transparent)
	btn.txt.TextSize = 12.0

	var objs []fyne.CanvasObject
	if item.Icon != "" {
		btn.icon = NewColoredIcon(item.Icon, 13, "#64748b")
		objs = append(objs, container.NewCenter(btn.icon))
		objs = append(objs, browseSpacer(4, 1))
	}
	objs = append(objs, container.NewCenter(btn.txt))

	innerHBox := container.NewHBox(objs...)
	btn.content = container.NewStack(
		btn.bg,
		container.NewBorder(
			browseSpacer(1, 4),
			browseSpacer(1, 4),
			browseSpacer(10, 1),
			browseSpacer(10, 1),
			innerHBox,
		),
	)

	btn.ExtendBaseWidget(btn)
	return btn
}

func (b *modernSegmentBtn) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(b.content)
}

func (b *modernSegmentBtn) Refresh() {
	isSelected := b.parent.selectedIdx == b.index

	if isSelected {
		b.bg.FillColor = ThemeColor(ColorNameCardBg)
		b.bg.StrokeColor = ThemeColor(ColorNameBorder)
		b.bg.StrokeWidth = 1
		b.txt.Color = ThemeColor(ColorNameTextPrimary)
		b.txt.TextStyle = fyne.TextStyle{Bold: true}
		if b.icon != nil {
			b.icon.SetColor(ThemeColorAccent())
		}
	} else {
		b.bg.StrokeWidth = 0
		if b.isHovered {
			b.bg.FillColor = ThemeColor(theme.ColorNameHover)
		} else {
			b.bg.FillColor = color.Transparent
		}

		b.txt.Color = ThemeColor(ColorNameTextMuted)
		b.txt.TextStyle = fyne.TextStyle{Bold: false}
		if b.icon != nil {
			b.icon.SetColor(ThemeColor(ColorNameTextMuted))
		}
	}

	b.bg.Refresh()
	b.txt.Refresh()
	if b.icon != nil {
		b.icon.Refresh()
	}
	if b.content != nil {
		b.content.Refresh()
	}
	b.BaseWidget.Refresh()
}

func (b *modernSegmentBtn) Tapped(_ *fyne.PointEvent) {
	b.parent.SetSelected(b.index)
}

func (b *modernSegmentBtn) MouseIn(_ *desktop.MouseEvent) {
	b.isHovered = true
	b.Refresh()
}

func (b *modernSegmentBtn) MouseMoved(_ *desktop.MouseEvent) {}

func (b *modernSegmentBtn) MouseOut() {
	b.isHovered = false
	b.Refresh()
}
