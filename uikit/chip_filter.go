package uikit

import (
	"image/color"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ModernChipFilter es un control visual contemporáneo tipo pastilla (chip) con selector desplegable,
// diseñado para filtrar datos por criterios específicos (ej. "Año: 2026", "Tipo: Real", etc.).
type ModernChipFilter struct {
	widget.BaseWidget

	mu            sync.RWMutex
	prefix        string
	options       []string
	selected      string
	defaultOption string
	onChanged     func(selected string)

	bg          *canvas.Rectangle
	iconText    *canvas.Text
	lblPrefix   *canvas.Text
	lblValue    *canvas.Text
	iconChevron *canvas.Text
	btnClear    *canvas.Text
	content     *fyne.Container

	isHovered bool
}

var _ desktop.Hoverable = (*ModernChipFilter)(nil)
var _ fyne.Tappable = (*ModernChipFilter)(nil)
var _ fyne.Widget = (*ModernChipFilter)(nil)

// NewChipFilter crea un nuevo ModernChipFilter.
func NewChipFilter(prefix string, options []string, defaultOption string, onChanged func(selected string)) *ModernChipFilter {
	if defaultOption == "" {
		if len(options) > 0 {
			defaultOption = options[0]
		} else {
			defaultOption = "Todos"
		}
	}

	cf := &ModernChipFilter{
		prefix:        prefix,
		options:       options,
		selected:      defaultOption,
		defaultOption: defaultOption,
		onChanged:     onChanged,
	}

	cf.bg = canvas.NewRectangle(color.Transparent)
	cf.bg.CornerRadius = 13
	cf.bg.StrokeWidth = 1

	// Icono sutil de filtro / embudo
	cf.iconText = canvas.NewText("⏳", color.Transparent)
	cf.iconText.TextSize = 11.5
	cf.iconText.TextStyle = fyne.TextStyle{Bold: true}

	cf.lblPrefix = canvas.NewText(prefix, color.Transparent)
	cf.lblPrefix.TextSize = 11.5
	cf.lblPrefix.TextStyle = fyne.TextStyle{Bold: false}

	cf.lblValue = canvas.NewText("", color.Transparent)
	cf.lblValue.TextSize = 11.5
	cf.lblValue.TextStyle = fyne.TextStyle{Bold: true}

	cf.iconChevron = canvas.NewText("▾", color.Transparent)
	cf.iconChevron.TextSize = 10.0

	cf.btnClear = canvas.NewText("✕", color.Transparent)
	cf.btnClear.TextSize = 10.0
	cf.btnClear.TextStyle = fyne.TextStyle{Bold: true}
	cf.btnClear.Hide()

	cf.content = container.NewHBox(
		cf.iconText,
		cf.lblPrefix,
		cf.lblValue,
		cf.iconChevron,
		cf.btnClear,
	)

	cf.ExtendBaseWidget(cf)
	cf.Refresh()
	return cf
}

// WithPrefix define o actualiza el texto prefijo del chip (ej. "Año:").
func (cf *ModernChipFilter) WithPrefix(prefix string) *ModernChipFilter {
	cf.mu.Lock()
	cf.prefix = prefix
	cf.mu.Unlock()
	cf.Refresh()
	return cf
}

// SetOptions actualiza dinámicamente la lista de opciones disponibles en el selector.
func (cf *ModernChipFilter) SetOptions(options []string) *ModernChipFilter {
	cf.mu.Lock()
	cf.options = options
	// Si la opción seleccionada ya no existe y no es la por defecto, restablecer
	found := false
	for _, opt := range options {
		if opt == cf.selected {
			found = true
			break
		}
	}
	if !found && len(options) > 0 {
		cf.selected = cf.defaultOption
	}
	cf.mu.Unlock()
	cf.Refresh()
	return cf
}

// Selected obtiene la opción actualmente seleccionada.
func (cf *ModernChipFilter) Selected() string {
	cf.mu.RLock()
	defer cf.mu.RUnlock()
	return cf.selected
}

// SetSelected establece la opción seleccionada y dispara el callback si hubo cambios.
func (cf *ModernChipFilter) SetSelected(option string) *ModernChipFilter {
	cf.mu.Lock()
	changed := cf.selected != option
	cf.selected = option
	cb := cf.onChanged
	cf.mu.Unlock()

	cf.Refresh()

	if changed && cb != nil {
		cb(option)
	}
	return cf
}

// SetSelectedSilent establece la opción seleccionada sin disparar el callback.
func (cf *ModernChipFilter) SetSelectedSilent(option string) *ModernChipFilter {
	cf.mu.Lock()
	cf.selected = option
	cf.mu.Unlock()
	cf.Refresh()
	return cf
}

// IsAtDefault indica si el filtro está actualmente en su opción por defecto.
func (cf *ModernChipFilter) IsAtDefault() bool {
	cf.mu.RLock()
	defer cf.mu.RUnlock()
	return cf.selected == cf.defaultOption
}

// Reset restablece el chip a su opción por defecto.
func (cf *ModernChipFilter) Reset() *ModernChipFilter {
	cf.mu.RLock()
	def := cf.defaultOption
	cf.mu.RUnlock()
	return cf.SetSelected(def)
}

// ShowMenu despliega el menú contextual de selección de opciones justo debajo del chip.
func (cf *ModernChipFilter) ShowMenu() {
	cf.mu.RLock()
	opts := cf.options
	sel := cf.selected
	def := cf.defaultOption
	prefix := cf.prefix
	cf.mu.RUnlock()

	if len(opts) == 0 {
		return
	}

	var items []*fyne.MenuItem

	// Título informativo del menú
	titleText := "FILTRAR"
	if prefix != "" {
		titleText = "FILTRAR POR " + strings.ToUpper(strings.TrimSuffix(prefix, ":"))
	}
	titleItem := fyne.NewMenuItem(titleText, nil)
	titleItem.Disabled = true
	items = append(items, titleItem)

	for _, opt := range opts {
		captured := opt
		item := fyne.NewMenuItem(captured, func() {
			cf.SetSelected(captured)
		})
		if captured == sel {
			item.Icon = theme.ConfirmIcon()
		}
		items = append(items, item)
	}

	if sel != def {
		items = append(items, fyne.NewMenuItemSeparator())
		resetItem := fyne.NewMenuItem("Quitar filtro ("+def+")", func() {
			cf.Reset()
		})
		items = append(items, resetItem)
	}

	menu := fyne.NewMenu("", items...)

	driver := fyne.CurrentApp().Driver()
	canvasObj := driver.CanvasForObject(cf)
	if canvasObj != nil {
		pos := fyne.NewPos(0, cf.Size().Height+4)
		widget.ShowPopUpMenuAtRelativePosition(menu, canvasObj, pos, cf)
	}
}

// Cursor define cursor de mano interactivo.
func (cf *ModernChipFilter) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (cf *ModernChipFilter) MouseIn(_ *desktop.MouseEvent) {
	cf.isHovered = true
	cf.Refresh()
}

func (cf *ModernChipFilter) MouseMoved(_ *desktop.MouseEvent) {}

func (cf *ModernChipFilter) MouseOut() {
	cf.isHovered = false
	cf.Refresh()
}

func (cf *ModernChipFilter) Tapped(e *fyne.PointEvent) {
	// Si está filtrado y el usuario hace clic en el área de la '✕' derecha, restablecer inmediatamente
	if !cf.IsAtDefault() && e.Position.X > cf.Size().Width-28 {
		cf.Reset()
		return
	}
	cf.ShowMenu()
}

type modernChipFilterRenderer struct {
	cf   *ModernChipFilter
	root *fyne.Container
}

func (cf *ModernChipFilter) CreateRenderer() fyne.WidgetRenderer {
	padX := float32(10)
	padded := container.NewBorder(
		nil, nil,
		browseSpacer(padX, 1),
		browseSpacer(padX, 1),
		cf.content,
	)
	root := container.NewStack(cf.bg, padded)

	r := &modernChipFilterRenderer{
		cf:   cf,
		root: root,
	}
	r.Refresh()
	return r
}

func (r *modernChipFilterRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.root}
}

func (r *modernChipFilterRenderer) Destroy() {}

func (r *modernChipFilterRenderer) MinSize() fyne.Size {
	s := r.root.MinSize()
	var minH float32 = 28
	if s.Height < minH {
		s.Height = minH
	}
	return s
}

func (r *modernChipFilterRenderer) Layout(size fyne.Size) {
	r.root.Resize(size)
	r.cf.bg.CornerRadius = size.Height / 2
}

func (r *modernChipFilterRenderer) Refresh() {
	cf := r.cf
	cf.mu.RLock()
	prefix := cf.prefix
	sel := cf.selected
	def := cf.defaultOption
	isHovered := cf.isHovered
	cf.mu.RUnlock()

	cf.lblPrefix.Text = prefix
	cf.lblValue.Text = sel

	isFiltered := sel != def

	if isFiltered {
		cf.btnClear.Show()
		cf.iconChevron.Hide()
	} else {
		cf.btnClear.Hide()
		cf.iconChevron.Show()
	}

	colTextPrimary := ThemeColor(ColorNameTextPrimary)
	colTextMuted := ThemeColor(ColorNameTextMuted)
	colAccent := ThemeColor(ColorNameAccent)
	colBorder := ThemeColor(ColorNameBorder)
	colBtnDefaultBg := ThemeColor(ColorNameBtnDefaultBg)
	colBtnHoverBg := ThemeColor(ColorNameBtnDefaultHoverBg)

	if isFiltered {
		// Estado ACTIVO: resaltado con acento de marca
		cf.bg.StrokeColor = colAccent
		cf.bg.StrokeWidth = 1.5
		cf.iconText.Color = colAccent
		cf.lblPrefix.Color = colAccent
		cf.lblValue.Color = colTextPrimary
		cf.btnClear.Color = colAccent
		if isHovered {
			cf.bg.FillColor = colBtnHoverBg
		} else {
			cf.bg.FillColor = colBtnDefaultBg
		}
	} else {
		// Estado INACTIVO / POR DEFECTO
		cf.bg.StrokeWidth = 1.0
		cf.iconText.Color = colTextMuted
		cf.lblPrefix.Color = colTextMuted
		cf.lblValue.Color = colTextPrimary
		cf.iconChevron.Color = colTextMuted
		if isHovered {
			cf.bg.StrokeColor = colAccent
			cf.bg.FillColor = colBtnHoverBg
		} else {
			cf.bg.StrokeColor = colBorder
			cf.bg.FillColor = colBtnDefaultBg
		}
	}

	cf.bg.Refresh()
	cf.iconText.Refresh()
	cf.lblPrefix.Refresh()
	cf.lblValue.Refresh()
	cf.iconChevron.Refresh()
	cf.btnClear.Refresh()
}
