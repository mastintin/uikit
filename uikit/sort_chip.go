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

// SortChipOption representa una columna u opción de ordenación seleccionable.
type SortChipOption struct {
	ID    string // Identificador o clave
	Label string // Texto visible en el chip y en el menú
	Index int    // Índice numérico asociado (por ejemplo en ModernBrowse)
}

// ModernSortChip es un control visual interactivo estilo pastilla (chip) que muestra
// el criterio actual de ordenación y permite cambiar columna y dirección mediante un menú emergente.
type ModernSortChip struct {
	widget.BaseWidget

	options    []SortChipOption
	currentIdx int
	currentDir SortDirection
	defaultIdx int
	defaultDir SortDirection
	onChanged  func(opt SortChipOption, dir SortDirection)

	bg          *canvas.Rectangle
	iconSort    *canvas.Text
	label       *canvas.Text
	iconChevron *canvas.Text
	content     *fyne.Container

	isHovered bool
}

var _ desktop.Hoverable = (*ModernSortChip)(nil)
var _ fyne.Tappable = (*ModernSortChip)(nil)
var _ fyne.Widget = (*ModernSortChip)(nil)

// NewSortChip crea un ModernSortChip a partir de nombres de columnas y valores por defecto.
func NewSortChip(columnTitles []string, defaultColumn string, defaultDir SortDirection, onChanged func(col string, dir SortDirection)) *ModernSortChip {
	options := make([]SortChipOption, len(columnTitles))
	defaultIdx := 0

	for i, title := range columnTitles {
		options[i] = SortChipOption{
			ID:    title,
			Label: title,
			Index: i,
		}
		if strings.EqualFold(title, defaultColumn) {
			defaultIdx = i
		}
	}

	var cb func(opt SortChipOption, dir SortDirection)
	if onChanged != nil {
		cb = func(opt SortChipOption, dir SortDirection) {
			onChanged(opt.Label, dir)
		}
	}

	return NewSortChipOptions(options, defaultIdx, defaultDir, cb)
}

// NewSortChipOptions crea un ModernSortChip con opciones estructuradas completas.
func NewSortChipOptions(options []SortChipOption, defaultIdx int, defaultDir SortDirection, onChanged func(opt SortChipOption, dir SortDirection)) *ModernSortChip {
	if defaultIdx < 0 || defaultIdx >= len(options) {
		defaultIdx = 0
	}
	if defaultDir != SortAscending && defaultDir != SortDescending {
		defaultDir = SortAscending
	}

	sc := &ModernSortChip{
		options:    options,
		currentIdx: defaultIdx,
		currentDir: defaultDir,
		defaultIdx: defaultIdx,
		defaultDir: defaultDir,
		onChanged:  onChanged,
	}

	sc.bg = canvas.NewRectangle(color.Transparent)
	sc.bg.CornerRadius = 13
	sc.bg.StrokeWidth = 1

	sc.iconSort = canvas.NewText("⇅", color.Transparent)
	sc.iconSort.TextSize = 11.5
	sc.iconSort.TextStyle = fyne.TextStyle{Bold: true}

	sc.label = canvas.NewText("", color.Transparent)
	sc.label.TextSize = 12.0
	sc.label.TextStyle = fyne.TextStyle{Bold: true}

	sc.iconChevron = canvas.NewText("▾", color.Transparent)
	sc.iconChevron.TextSize = 10.0

	sc.content = container.NewHBox(
		sc.iconSort,
		sc.label,
		sc.iconChevron,
	)

	sc.ExtendBaseWidget(sc)
	sc.Refresh()
	return sc
}

// NewBrowseSortChip crea un ModernSortChip vinculado y bidireccionalmente sincronizado con un ModernBrowse.
func NewBrowseSortChip(mb *ModernBrowse) *ModernSortChip {
	if mb == nil {
		return NewSortChipOptions(nil, 0, SortAscending, nil)
	}

	var options []SortChipOption
	for i, col := range mb.columns {
		if !col.NoSort && strings.TrimSpace(col.Title) != "" {
			options = append(options, SortChipOption{
				ID:    col.Title,
				Label: col.Title,
				Index: i,
			})
		}
	}

	defaultIdx := 0
	defaultDir := SortAscending

	if mb.sortColIdx >= 0 {
		for i, opt := range options {
			if opt.Index == mb.sortColIdx {
				defaultIdx = i
				if mb.sortDir != SortNone {
					defaultDir = mb.sortDir
				}
				break
			}
		}
	}

	chip := NewSortChipOptions(options, defaultIdx, defaultDir, func(opt SortChipOption, dir SortDirection) {
		mb.SortBy(opt.Index, dir)
	})

	// Sincronización inversa cuando el usuario hace clic en una cabecera de columna de ModernBrowse
	prevOnSort := mb.onSortChanged
	mb.SetOnSortChanged(func(colIdx int, dir SortDirection) {
		for i, opt := range chip.options {
			if opt.Index == colIdx {
				chip.SetSortSilent(i, dir)
				break
			}
		}
		if prevOnSort != nil {
			prevOnSort(colIdx, dir)
		}
	})

	return chip
}

// CreateSortChip genera un ModernSortChip directamente desde la instancia de ModernBrowse.
func (mb *ModernBrowse) CreateSortChip() *ModernSortChip {
	return NewBrowseSortChip(mb)
}

// IsAtDefault indica si el chip está actualmente en su criterio y dirección por defecto.
func (sc *ModernSortChip) IsAtDefault() bool {
	return sc.currentIdx == sc.defaultIdx && sc.currentDir == sc.defaultDir
}

// CurrentSort devuelve la opción y dirección actualmente activas.
func (sc *ModernSortChip) CurrentSort() (SortChipOption, SortDirection) {
	if len(sc.options) == 0 {
		return SortChipOption{}, SortNone
	}
	return sc.options[sc.currentIdx], sc.currentDir
}

// SetSort establece la columna y dirección de ordenación disparando el callback.
func (sc *ModernSortChip) SetSort(optIdx int, dir SortDirection) *ModernSortChip {
	if optIdx < 0 || optIdx >= len(sc.options) {
		return sc
	}
	if dir != SortAscending && dir != SortDescending {
		dir = SortAscending
	}

	changed := sc.currentIdx != optIdx || sc.currentDir != dir
	sc.currentIdx = optIdx
	sc.currentDir = dir
	sc.Refresh()

	if changed && sc.onChanged != nil {
		sc.onChanged(sc.options[optIdx], dir)
	}
	return sc
}

// SetSortSilent actualiza el estado visual del chip sin disparar el listener de cambio.
func (sc *ModernSortChip) SetSortSilent(optIdx int, dir SortDirection) *ModernSortChip {
	if optIdx < 0 || optIdx >= len(sc.options) {
		return sc
	}
	if dir != SortAscending && dir != SortDescending {
		dir = SortAscending
	}

	sc.currentIdx = optIdx
	sc.currentDir = dir
	sc.Refresh()
	return sc
}

// Reset restablece el chip a su ordenación por defecto llamando al callback.
func (sc *ModernSortChip) Reset() *ModernSortChip {
	return sc.SetSort(sc.defaultIdx, sc.defaultDir)
}

// ShowMenu despliega el menú contextual de ordenación justo debajo del chip.
func (sc *ModernSortChip) ShowMenu() {
	if len(sc.options) == 0 {
		return
	}

	currentDir := sc.currentDir

	var items []*fyne.MenuItem

	// 1. Título de Columnas
	titleCols := fyne.NewMenuItem("ORDENAR POR", nil)
	titleCols.Disabled = true
	items = append(items, titleCols)

	// Lista de columnas disponibles
	for i, opt := range sc.options {
		optIdx := i
		item := fyne.NewMenuItem(opt.Label, func() {
			sc.SetSort(optIdx, currentDir)
		})
		if optIdx == sc.currentIdx {
			item.Icon = theme.ConfirmIcon()
		}
		items = append(items, item)
	}

	items = append(items, fyne.NewMenuItemSeparator())

	// 2. Título de Dirección
	titleDir := fyne.NewMenuItem("DIRECCIÓN", nil)
	titleDir.Disabled = true
	items = append(items, titleDir)

	itemAsc := fyne.NewMenuItem("Ascendente (A-Z, 0-9)", func() {
		sc.SetSort(sc.currentIdx, SortAscending)
	})
	if currentDir == SortAscending {
		itemAsc.Icon = theme.ConfirmIcon()
	}
	items = append(items, itemAsc)

	itemDesc := fyne.NewMenuItem("Descendente (Z-A, 9-0)", func() {
		sc.SetSort(sc.currentIdx, SortDescending)
	})
	if currentDir == SortDescending {
		itemDesc.Icon = theme.ConfirmIcon()
	}
	items = append(items, itemDesc)

	// 3. Separador y Restablecer
	items = append(items, fyne.NewMenuItemSeparator())
	itemReset := fyne.NewMenuItem("Restablecer orden por defecto", func() {
		sc.Reset()
	})
	itemReset.Disabled = sc.IsAtDefault()
	items = append(items, itemReset)

	menu := fyne.NewMenu("", items...)

	// Posicionar y mostrar menú relativo al widget
	driver := fyne.CurrentApp().Driver()
	canvasObj := driver.CanvasForObject(sc)
	if canvasObj != nil {
		pos := fyne.NewPos(0, sc.Size().Height+4)
		widget.ShowPopUpMenuAtRelativePosition(menu, canvasObj, pos, sc)
	}
}

// --- Eventos de Escritorio y Renderer ---

func (sc *ModernSortChip) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (sc *ModernSortChip) MouseIn(_ *desktop.MouseEvent) {
	sc.isHovered = true
	sc.Refresh()
}

func (sc *ModernSortChip) MouseMoved(_ *desktop.MouseEvent) {}

func (sc *ModernSortChip) MouseOut() {
	sc.isHovered = false
	sc.Refresh()
}

func (sc *ModernSortChip) Tapped(_ *fyne.PointEvent) {
	sc.ShowMenu()
}

type modernSortChipRenderer struct {
	sc   *ModernSortChip
	root *fyne.Container
}

func (sc *ModernSortChip) CreateRenderer() fyne.WidgetRenderer {
	padX := float32(11)
	padded := container.NewBorder(
		nil, nil,
		browseSpacer(padX, 1),
		browseSpacer(padX, 1),
		sc.content,
	)
	root := container.NewStack(sc.bg, padded)

	r := &modernSortChipRenderer{
		sc:   sc,
		root: root,
	}
	r.Refresh()
	return r
}

func (r *modernSortChipRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.root}
}

func (r *modernSortChipRenderer) Destroy() {}

func (r *modernSortChipRenderer) MinSize() fyne.Size {
	s := r.root.MinSize()
	var minH float32 = 28
	if s.Height < minH {
		s.Height = minH
	}
	return s
}

func (r *modernSortChipRenderer) Layout(size fyne.Size) {
	r.root.Resize(size)
	r.sc.bg.CornerRadius = size.Height / 2
}

func (r *modernSortChipRenderer) Refresh() {
	sc := r.sc
	if len(sc.options) == 0 {
		sc.label.Text = "Sin ordenar"
		sc.iconSort.Text = "⇅"
	} else {
		opt := sc.options[sc.currentIdx]
		sc.label.Text = opt.Label
		if sc.currentDir == SortAscending {
			sc.iconSort.Text = "↑"
		} else if sc.currentDir == SortDescending {
			sc.iconSort.Text = "↓"
		} else {
			sc.iconSort.Text = "⇅"
		}
	}

	// Estilos dinámicos adaptativos según tema y estado
	atDefault := sc.IsAtDefault()
	pal := GetCurrentThemePalette()

	var bgFill, strokeCol, textCol color.Color

	if atDefault {
		// Modo reposo / neutro
		bgFill = ThemeColor(ColorNameInsetPanelBg)
		strokeCol = ThemeColor(ColorNameBorder)
		textCol = ThemeColor(ColorNameTextMuted)
		if sc.isHovered {
			bgFill = ThemeColor(theme.ColorNameHover)
			textCol = ThemeColor(ColorNameTextPrimary)
		}
	} else {
		// Modo activo / modificado (usa tokens informativos del tema)
		if pal != nil {
			bgFill = pal.Badges.Info.Bg
			strokeCol = pal.Badges.Info.Border
			textCol = pal.Badges.Info.Text
			if sc.isHovered {
				strokeCol = pal.Accent
			}
		} else {
			bgFill = ThemeColor(theme.ColorNameSelection)
			strokeCol = ThemeColor(ColorNameAccent)
			textCol = ThemeColor(ColorNameAccent)
		}
	}

	sc.bg.FillColor = bgFill
	sc.bg.StrokeColor = strokeCol
	sc.iconSort.Color = textCol
	sc.label.Color = textCol
	sc.iconChevron.Color = textCol

	sc.bg.Refresh()
	sc.iconSort.Refresh()
	sc.label.Refresh()
	sc.iconChevron.Refresh()
	r.root.Refresh()
}
