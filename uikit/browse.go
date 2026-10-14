package uikit

import (
	"image/color"
	"sort"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ─────────────────────────────────────────────────────────────────────────────
// ModernBrowse — Control de tabla / browse de auditoría tabular estilo xlsJubila
// ─────────────────────────────────────────────────────────────────────────────

// BrowseAlign define la alineación horizontal de una celda o columna.
type BrowseAlign int

const (
	BrowseAlignLeading BrowseAlign = iota
	BrowseAlignCenter
	BrowseAlignTrailing
)

// BrowseRowType define el estado semántico de una fila en el browse.
type BrowseRowType int

const (
	BrowseRowNormal BrowseRowType = iota
	BrowseRowDiscarded // Fila descartada (fondo rojo suave estilo xlsJubila)
	BrowseRowLaguna    // Laguna de cotización (fondo ámbar suave estilo xlsJubila)
	BrowseRowSuccess   // Fila destacada en positivo / éxito
)

// SortDirection define la dirección de ordenación de una columna.
type SortDirection int

const (
	SortNone SortDirection = iota
	SortAscending
	SortDescending
)

// BrowseColumn define una columna de la tabla.
type BrowseColumn struct {
	Title      string
	Align      BrowseAlign
	Sortable   bool // Si es false (por defecto cuando no se indica), se permite ordenar a menos que se inhabilite explícitamente
	NoSort     bool // Si es true, inhabilita la ordenación en esta columna específica (ej. botones de acciones)
	Comparator func(a, b BrowseRow) int // Comparador personalizado opcional
}

// BrowseCell representa el contenido de una celda individual.
type BrowseCell struct {
	Text      string
	Bold      bool
	Color     color.Color
	CustomObj fyne.CanvasObject
	Align     BrowseAlign
	TextSize  float32
	SortValue any // Valor opcional numérico o tipado para ordenar con máxima precisión
}

// CellText crea una celda de texto estándar.
func CellText(text string, align ...BrowseAlign) BrowseCell {
	al := BrowseAlignLeading
	if len(align) > 0 {
		al = align[0]
	}
	return BrowseCell{Text: text, Align: al}
}

// CellBold crea una celda de texto en negrita.
func CellBold(text string, col color.Color, align ...BrowseAlign) BrowseCell {
	al := BrowseAlignLeading
	if len(align) > 0 {
		al = align[0]
	}
	return BrowseCell{Text: text, Bold: true, Color: col, Align: al}
}

// CellObject crea una celda con un CanvasObject arbitrario (ej. un Badge).
func CellObject(obj fyne.CanvasObject, align ...BrowseAlign) BrowseCell {
	al := BrowseAlignCenter
	if len(align) > 0 {
		al = align[0]
	}
	return BrowseCell{CustomObj: obj, Align: al}
}

// BrowseRow representa una fila completa de datos.
type BrowseRow struct {
	Type  BrowseRowType
	Cells []BrowseCell
}

// browseContinuousLayout apila filas verticalmente de manera 100% continua
// sin espacios ni gaps de theme.Padding(), permitiendo que los fondos cubran
// el alto entero y las líneas divisorias de 1px se posicionen con precisión.
type browseContinuousLayout struct{}

func (b *browseContinuousLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	y := float32(0)
	for _, obj := range objects {
		if !obj.Visible() {
			continue
		}
		h := obj.MinSize().Height
		obj.Resize(fyne.NewSize(size.Width, h))
		obj.Move(fyne.NewPos(0, y))
		y += h
	}
}

func (b *browseContinuousLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w := float32(0)
	h := float32(0)
	for _, obj := range objects {
		if !obj.Visible() {
			continue
		}
		min := obj.MinSize()
		if min.Width > w {
			w = min.Width
		}
		h += min.Height
	}
	return fyne.NewSize(w, h)
}

// ModernBrowse es un control tabular compacto y limpio inspirado en las tablas
// de auditoría de xlsJubila:
//   - Cabecera ejecutiva sólida con fondo corporativo y texto nítido.
//   - Filas limpias continuas que ajustan al alto entero sin cortes laterales ni huecos.
//   - Fondo semántico por fila (normal cebra, descartada en rojo suave, laguna en ámbar).
//   - Adaptación dinámica completa a los temas Nice (xlsJubila), Claro y Oscuro.
type ModernBrowse struct {
	widget.BaseWidget

	columns         []BrowseColumn
	rows            []BrowseRow
	container       *fyne.Container
	rowsBox         *fyne.Container
	scrollRows      *container.Scroll
	headerContainer *fyne.Container
	defaultTextSize float32
	maxHeight       float32

	sortColIdx int           // Columna actualmente ordenada (-1 si ninguna)
	sortDir    SortDirection // Dirección actual (SortNone, SortAscending, SortDescending)

	onSortChanged func(colIdx int, dir SortDirection)
}

// browseTableLayout mantiene la cabecera fija arriba (sticky header)
// y permite que el área de filas (scrollRows) haga scroll verticalmente.
type browseTableLayout struct {
	mb *ModernBrowse
}

func (l *browseTableLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) < 2 {
		return
	}
	header := objects[0]
	scroll := objects[1]

	hHeader := header.MinSize().Height
	header.Resize(fyne.NewSize(size.Width, hHeader))
	header.Move(fyne.NewPos(0, 0))

	hScroll := size.Height - hHeader
	if hScroll < 0 {
		hScroll = 0
	}
	scroll.Resize(fyne.NewSize(size.Width, hScroll))
	scroll.Move(fyne.NewPos(0, hHeader))
}

func (l *browseTableLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) < 2 {
		return fyne.NewSize(0, 0)
	}
	header := objects[0]

	wHeader := header.MinSize().Width
	wRows := l.mb.rowsBox.MinSize().Width
	w := wHeader
	if wRows > w {
		w = wRows
	}

	hHeader := header.MinSize().Height
	hRows := l.mb.rowsBox.MinSize().Height

	if l.mb.maxHeight > 0 && hRows > l.mb.maxHeight {
		hRows = l.mb.maxHeight
	}

	return fyne.NewSize(w, hHeader+hRows)
}

// NewModernBrowse crea un nuevo browse configurado con las columnas especificadas.
func NewModernBrowse(columns []BrowseColumn) *ModernBrowse {
	mb := &ModernBrowse{
		columns:         columns,
		defaultTextSize: 13.5,
		maxHeight:       440,
		rowsBox:         container.New(&browseContinuousLayout{}),
		headerContainer: container.NewStack(),
		sortColIdx:      -1,
		sortDir:         SortNone,
	}

	mb.build()
	mb.ExtendBaseWidget(mb)
	return mb
}

// SetMaxHeight configura la altura máxima visible de las filas con scroll interno.
// Si maxHeight <= 0, no se limita la altura (todas las filas se despliegan continuas sin scroll).
func (mb *ModernBrowse) SetMaxHeight(h float32) *ModernBrowse {
	mb.maxHeight = h
	if mb.container != nil {
		mb.container.Refresh()
	}
	return mb
}

// SetTextSize ajusta el tamaño de fuente por defecto de las celdas de datos.
func (mb *ModernBrowse) SetTextSize(size float32) *ModernBrowse {
	mb.defaultTextSize = size
	mb.rebuildRows()
	return mb
}

// CreateRenderer implementa la interfaz fyne.Widget usando NewSimpleRenderer.
func (mb *ModernBrowse) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(mb.container)
}

func browseSpacer(w, h float32) fyne.CanvasObject {
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(w, h))
	return rect
}

func (mb *ModernBrowse) build() {
	mb.headerContainer.Objects = []fyne.CanvasObject{mb.buildHeader()}

	mb.scrollRows = container.NewVScroll(mb.rowsBox)

	inner := container.New(&browseTableLayout{mb: mb},
		mb.headerContainer,
		mb.scrollRows,
	)

	// Contenedor principal con tarjeta única envolvente para toda la tabla
	card := NewCard(inner).WithVerticalPadding(0, 0).WithHorizontalPadding(0, 0)
	mb.container = container.NewStack(card)
}

// browseHeaderCell es un widget interactivo para cada cabecera con hover, click y puntero.
type browseHeaderCell struct {
	widget.BaseWidget

	title     string
	align     BrowseAlign
	sortDir   SortDirection
	canSort   bool
	onTapped  func()
	isHovered bool
}

func newBrowseHeaderCell(title string, align BrowseAlign, sortDir SortDirection, canSort bool, onTapped func()) *browseHeaderCell {
	c := &browseHeaderCell{
		title:    title,
		align:    align,
		sortDir:  sortDir,
		canSort:  canSort,
		onTapped: onTapped,
	}
	c.ExtendBaseWidget(c)
	return c
}

func (c *browseHeaderCell) Tapped(_ *fyne.PointEvent) {
	if c.canSort && c.onTapped != nil {
		c.onTapped()
	}
}

func (c *browseHeaderCell) MouseIn(_ *desktop.MouseEvent) {
	if c.canSort {
		c.isHovered = true
		c.Refresh()
	}
}

func (c *browseHeaderCell) MouseOut() {
	if c.canSort {
		c.isHovered = false
		c.Refresh()
	}
}

func (c *browseHeaderCell) MouseMoved(_ *desktop.MouseEvent) {}

func (c *browseHeaderCell) Cursor() desktop.Cursor {
	if c.canSort {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

func (c *browseHeaderCell) CreateRenderer() fyne.WidgetRenderer {
	txt := canvas.NewText(c.title, ColorWhite)
	txt.TextSize = 12.5
	txt.TextStyle = fyne.TextStyle{Bold: true}

	var indicator *canvas.Text
	if c.sortDir == SortAscending {
		indicator = canvas.NewText(" ▲", ThemeColorAccent())
		indicator.TextSize = 10.5
		indicator.TextStyle = fyne.TextStyle{Bold: true}
	} else if c.sortDir == SortDescending {
		indicator = canvas.NewText(" ▼", ThemeColorAccent())
		indicator.TextSize = 10.5
		indicator.TextStyle = fyne.TextStyle{Bold: true}
	}

	hoverBg := canvas.NewRectangle(color.Transparent)
	hoverBg.CornerRadius = 3

	var content fyne.CanvasObject
	if indicator != nil {
		content = container.NewHBox(txt, indicator)
	} else {
		content = txt
	}

	var aligned fyne.CanvasObject
	switch c.align {
	case BrowseAlignCenter:
		aligned = container.NewCenter(content)
	case BrowseAlignTrailing:
		aligned = container.NewBorder(nil, nil, nil, content, browseSpacer(1, 1))
	default:
		aligned = container.NewBorder(nil, nil, content, nil, browseSpacer(1, 1))
	}

	root := container.NewStack(hoverBg, aligned)

	r := &browseHeaderCellRenderer{
		cell:      c,
		txt:       txt,
		indicator: indicator,
		hoverBg:   hoverBg,
		root:      root,
	}
	r.applyStyle()
	return r
}

type browseHeaderCellRenderer struct {
	cell      *browseHeaderCell
	txt       *canvas.Text
	indicator *canvas.Text
	hoverBg   *canvas.Rectangle
	root      *fyne.Container
}

func (r *browseHeaderCellRenderer) applyStyle() {
	pal := GetCurrentThemePalette()

	var textCol color.Color = ThemeColor(ColorNameTextPrimary)
	var indicatorCol color.Color = ThemeColorAccent()
	if pal != nil {
		if pal.Table.HeaderFg != nil {
			textCol = pal.Table.HeaderFg
		}
		if (r.cell.isHovered || r.cell.sortDir != SortNone) && pal.Table.HeaderHoverFg != nil {
			textCol = pal.Table.HeaderHoverFg
		}
		if pal.Table.SortIndicator != nil {
			indicatorCol = pal.Table.SortIndicator
		}
	}

	r.txt.Color = textCol
	if r.indicator != nil {
		r.indicator.Color = indicatorCol
	}

	if r.cell.isHovered {
		r.hoverBg.FillColor = ThemeColor(theme.ColorNameHover)
	} else {
		r.hoverBg.FillColor = color.Transparent
	}
	r.hoverBg.Refresh()
	r.txt.Refresh()
	if r.indicator != nil {
		r.indicator.Refresh()
	}
}

func (r *browseHeaderCellRenderer) Layout(size fyne.Size) {
	r.root.Resize(size)
}

func (r *browseHeaderCellRenderer) MinSize() fyne.Size {
	return r.root.MinSize()
}

func (r *browseHeaderCellRenderer) Refresh() {
	r.applyStyle()
}

func (r *browseHeaderCellRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.root}
}

func (r *browseHeaderCellRenderer) Destroy() {}

func (mb *ModernBrowse) buildHeader() fyne.CanvasObject {
	numCols := len(mb.columns)
	if numCols == 0 {
		return container.NewHBox()
	}

	pal := GetCurrentThemePalette()
	var bgCol color.Color = ThemeColor(ColorNameNavbarBg)
	if pal != nil && pal.Table.HeaderBg != nil {
		bgCol = pal.Table.HeaderBg
	}

	headerBg := canvas.NewRectangle(bgCol)
	headerBg.CornerRadius = 0 // Plano continuo

	headers := make([]fyne.CanvasObject, numCols)
	for i, col := range mb.columns {
		colIdx := i
		canSort := !col.NoSort

		dir := SortNone
		if mb.sortColIdx == colIdx {
			dir = mb.sortDir
		}

		cell := newBrowseHeaderCell(col.Title, col.Align, dir, canSort, func() {
			mb.ToggleSort(colIdx)
		})
		headers[i] = cell
	}

	grid := container.NewGridWithColumns(numCols, headers...)
	headerContent := container.NewBorder(
		browseSpacer(1, 7),
		browseSpacer(1, 7),
		browseSpacer(12, 1),
		browseSpacer(12, 1),
		grid,
	)

	return container.NewStack(headerBg, headerContent)
}

// SetRows reemplaza todas las filas del browse y refresca la visualización.
func (mb *ModernBrowse) SetRows(rows []BrowseRow) *ModernBrowse {
	mb.rows = rows
	mb.rebuildRows()
	return mb
}

// AddRow añade una fila normal al browse.
func (mb *ModernBrowse) AddRow(cells ...BrowseCell) *ModernBrowse {
	mb.rows = append(mb.rows, BrowseRow{Type: BrowseRowNormal, Cells: cells})
	mb.rebuildRows()
	return mb
}

// AddSemanticRow añade una fila con un tipo semántico específico (ej. BrowseRowDiscarded).
func (mb *ModernBrowse) AddSemanticRow(rowType BrowseRowType, cells ...BrowseCell) *ModernBrowse {
	mb.rows = append(mb.rows, BrowseRow{Type: rowType, Cells: cells})
	mb.rebuildRows()
	return mb
}

// Clear vacía todas las filas del browse.
func (mb *ModernBrowse) Clear() *ModernBrowse {
	mb.rows = nil
	mb.rowsBox.Objects = nil
	mb.rowsBox.Refresh()
	if mb.scrollRows != nil {
		mb.scrollRows.Offset = fyne.NewPos(0, 0)
		mb.scrollRows.Refresh()
	}
	if mb.container != nil {
		mb.container.Refresh()
	}
	return mb
}

// Container devuelve el contenedor Fyne listo para incrustar en cualquier vista.
func (mb *ModernBrowse) Container() *fyne.Container {
	return mb.container
}

func (mb *ModernBrowse) rebuildRows() {
	var savedOffset fyne.Position
	hasSavedOffset := false
	if mb.scrollRows != nil && (mb.scrollRows.Offset.X != 0 || mb.scrollRows.Offset.Y != 0) {
		savedOffset = mb.scrollRows.Offset
		hasSavedOffset = true
	}

	mb.rowsBox.Objects = nil
	numCols := len(mb.columns)
	if numCols == 0 || len(mb.rows) == 0 {
		mb.rowsBox.Refresh()
		if mb.scrollRows != nil {
			mb.scrollRows.Refresh()
		}
		if mb.container != nil {
			mb.container.Refresh()
		}
		return
	}

	displayRows := mb.getDisplayRows()


	dividerCol := ThemeColor(ColorNameDivider)
	pal := GetCurrentThemePalette()

	for rowIndex, row := range displayRows {
		// Determinar colores de fondo y texto por fila
		var rowBgColor color.Color
		var defaultTextColor color.Color

		switch row.Type {
		case BrowseRowDiscarded:
			// Descartado: tokens semánticos de peligro / atenuado
			if pal != nil {
				rowBgColor = pal.Badges.Danger.Bg
				defaultTextColor = pal.Badges.Danger.Text
			} else {
				rowBgColor = ThemeColor(theme.ColorNameError)
				defaultTextColor = ColorWhite
			}
		case BrowseRowLaguna:
			// Laguna: tokens semánticos de advertencia
			if pal != nil {
				rowBgColor = pal.Badges.Warning.Bg
				defaultTextColor = pal.Badges.Warning.Text
			} else {
				rowBgColor = ThemeColor(theme.ColorNameWarning)
				defaultTextColor = ColorWhite
			}
		default:
			if pal != nil && pal.Table.RowOddBg != nil {
				if rowIndex%2 == 1 {
					rowBgColor = pal.Table.RowOddBg
				} else {
					rowBgColor = pal.Table.RowEvenBg
				}
				defaultTextColor = pal.TextPrimary
			} else {
				if rowIndex%2 == 1 {
					rowBgColor = ThemeColor(ColorNameInsetPanelBg)
				} else {
					rowBgColor = ThemeColor(ColorNameCardBg)
				}
				defaultTextColor = ThemeColor(ColorNameTextPrimary)
			}
		}


		// Fondo de fila continuo que cubre el 100% de la superficie
		bgRect := canvas.NewRectangle(rowBgColor)
		bgRect.CornerRadius = 0 // Plano y continuo sin muescas ni cortes laterales

		// Divisor tenue inferior a todo lo ancho (1px)
		divLine := canvas.NewRectangle(dividerCol)
		divLine.SetMinSize(fyne.NewSize(1, 1))

		baseTextSize := mb.defaultTextSize
		if baseTextSize <= 0 {
			baseTextSize = 13.5
		}

		// Celdas
		cellObjects := make([]fyne.CanvasObject, numCols)
		for cIdx := 0; cIdx < numCols; cIdx++ {
			if cIdx < len(row.Cells) {
				cell := row.Cells[cIdx]
				if cell.CustomObj != nil {
					switch cell.Align {
					case BrowseAlignCenter:
						cellObjects[cIdx] = container.NewCenter(cell.CustomObj)
					case BrowseAlignTrailing:
						cellObjects[cIdx] = container.NewBorder(nil, nil, nil, cell.CustomObj, browseSpacer(1, 1))
					default:
						cellObjects[cIdx] = container.NewBorder(nil, nil, cell.CustomObj, nil, browseSpacer(1, 1))
					}
				} else {
					txtColor := cell.Color
					if txtColor == nil {
						txtColor = defaultTextColor
					}
					txt := canvas.NewText(cell.Text, txtColor)
					sz := baseTextSize
					if cell.TextSize > 0 {
						sz = cell.TextSize
					}
					txt.TextSize = sz
					txt.TextStyle = fyne.TextStyle{Bold: cell.Bold}

					switch cell.Align {
					case BrowseAlignCenter:
						cellObjects[cIdx] = container.NewCenter(txt)
					case BrowseAlignTrailing:
						cellObjects[cIdx] = container.NewBorder(nil, nil, nil, txt, browseSpacer(1, 1))
					default:
						cellObjects[cIdx] = container.NewBorder(nil, nil, txt, nil, browseSpacer(1, 1))
					}
				}
			} else {
				cellObjects[cIdx] = browseSpacer(1, 1)
			}
		}

		grid := container.NewGridWithColumns(numCols, cellObjects...)
		rowContent := container.NewBorder(
			browseSpacer(1, 8),
			browseSpacer(1, 8),
			browseSpacer(12, 1),
			browseSpacer(12, 1),
			grid,
		)

		rowStack := container.NewStack(
			bgRect,
			container.NewBorder(nil, divLine, nil, nil, rowContent),
		)

		mb.rowsBox.Add(rowStack)
	}

	mb.rowsBox.Refresh()
	if mb.scrollRows != nil {
		if hasSavedOffset {
			mb.scrollRows.ScrollToOffset(savedOffset)
		}
		mb.scrollRows.Refresh()
	}
	if mb.container != nil {
		mb.container.Refresh()
	}
}

// GetScrollOffset devuelve la posición de scroll actual dentro del área de filas.
func (mb *ModernBrowse) GetScrollOffset() fyne.Position {
	if mb.scrollRows != nil {
		return mb.scrollRows.Offset
	}
	return fyne.Position{}
}

// SetScrollOffset establece la posición de desplazamiento dentro del browse.
func (mb *ModernBrowse) SetScrollOffset(pos fyne.Position) {
	if mb.scrollRows != nil {
		mb.scrollRows.ScrollToOffset(pos)
		mb.scrollRows.Refresh()
	}
}

// ScrollToRow desplaza la tabla para asegurar que la fila indicada por rowIndex quede visible.
func (mb *ModernBrowse) ScrollToRow(rowIndex int) {
	if mb.scrollRows == nil || mb.rowsBox == nil {
		return
	}
	objs := mb.rowsBox.Objects
	if rowIndex < 0 || rowIndex >= len(objs) {
		return
	}

	y := float32(0)
	for i := 0; i < rowIndex; i++ {
		y += objs[i].MinSize().Height
	}
	rowHeight := objs[rowIndex].MinSize().Height

	viewportHeight := mb.scrollRows.Size().Height
	if viewportHeight <= 0 && mb.maxHeight > 0 {
		viewportHeight = mb.maxHeight
	}

	currentY := mb.scrollRows.Offset.Y
	// Si la fila ya se encuentra visible en el viewport actual, no desplazar
	if viewportHeight > 0 && y >= currentY && (y+rowHeight) <= (currentY+viewportHeight) {
		return
	}

	// Dejar un margen de 1 fila arriba para dar contexto al usuario
	targetY := y
	if rowIndex > 0 {
		targetY = y - objs[rowIndex-1].MinSize().Height
	}
	if targetY < 0 {
		targetY = 0
	}
	mb.scrollRows.ScrollToOffset(fyne.NewPos(0, targetY))
	mb.scrollRows.Refresh()
}

// ToggleSort conmuta la ordenación para una columna dada:
// de None/Descendente -> Ascendente, de Ascendente -> Descendente.
func (mb *ModernBrowse) ToggleSort(colIdx int) {
	if colIdx < 0 || colIdx >= len(mb.columns) {
		return
	}
	if mb.columns[colIdx].NoSort {
		return
	}

	if mb.sortColIdx == colIdx {
		if mb.sortDir == SortAscending {
			mb.sortDir = SortDescending
		} else {
			mb.sortDir = SortAscending
		}
	} else {
		mb.sortColIdx = colIdx
		mb.sortDir = SortAscending
	}

	mb.SortByWithScrollReset(colIdx, mb.sortDir, true)
}

// SortBy permite ordenar programáticamente por una columna y dirección específica sin resetear el scroll a 0 si no se solicita.
func (mb *ModernBrowse) SortBy(colIdx int, dir SortDirection) {
	mb.SortByWithScrollReset(colIdx, dir, false)
}

// SortByWithScrollReset permite ordenar especificando si se resetea la posición de scroll a la cabecera.
func (mb *ModernBrowse) SortByWithScrollReset(colIdx int, dir SortDirection, resetScroll bool) {
	if colIdx < 0 || colIdx >= len(mb.columns) {
		return
	}
	mb.sortColIdx = colIdx
	mb.sortDir = dir
	if mb.headerContainer != nil {
		mb.headerContainer.Objects = []fyne.CanvasObject{mb.buildHeader()}
		mb.headerContainer.Refresh()
	}
	if resetScroll && mb.scrollRows != nil {
		mb.scrollRows.Offset = fyne.NewPos(0, 0)
	}
	mb.rebuildRows()

	if mb.onSortChanged != nil {
		mb.onSortChanged(mb.sortColIdx, mb.sortDir)
	}
}

// SetOnSortChanged configura un callback que se notifica cuando cambia la ordenación del browse.
func (mb *ModernBrowse) SetOnSortChanged(fn func(colIdx int, dir SortDirection)) *ModernBrowse {
	mb.onSortChanged = fn
	return mb
}

// CurrentSort devuelve el índice de columna y dirección actualmente ordenados.
func (mb *ModernBrowse) CurrentSort() (int, SortDirection) {
	return mb.sortColIdx, mb.sortDir
}

// Columns devuelve la lista de columnas configuradas en el browse.
func (mb *ModernBrowse) Columns() []BrowseColumn {
	return mb.columns
}

func (mb *ModernBrowse) sortAndRefresh() {
	mb.SortByWithScrollReset(mb.sortColIdx, mb.sortDir, true)
}

func (mb *ModernBrowse) getDisplayRows() []BrowseRow {
	if mb.sortColIdx < 0 || mb.sortColIdx >= len(mb.columns) || mb.sortDir == SortNone {
		return mb.rows
	}

	colIdx := mb.sortColIdx
	col := mb.columns[colIdx]

	display := make([]BrowseRow, len(mb.rows))
	copy(display, mb.rows)

	sort.SliceStable(display, func(i, j int) bool {
		rA := display[i]
		rB := display[j]

		var cmp int
		if col.Comparator != nil {
			cmp = col.Comparator(rA, rB)
		} else {
			var cellA, cellB BrowseCell
			if colIdx < len(rA.Cells) {
				cellA = rA.Cells[colIdx]
			}
			if colIdx < len(rB.Cells) {
				cellB = rB.Cells[colIdx]
			}
			cmp = compareBrowseCells(cellA, cellB)
		}

		if mb.sortDir == SortDescending {
			return cmp > 0
		}
		return cmp < 0
	})

	return display
}

func compareBrowseCells(a, b BrowseCell) int {
	// 1. Si tienen SortValue explícito
	if a.SortValue != nil || b.SortValue != nil {
		return compareAny(a.SortValue, b.SortValue)
	}

	// 2. Extraer texto limpio
	strA := strings.TrimSpace(a.Text)
	strB := strings.TrimSpace(b.Text)

	// Fechas estilo DD/MM/AAAA o MM / AAAA
	if tA, ok := parseSpanishDate(strA); ok {
		if tB, ok := parseSpanishDate(strB); ok {
			if tA.Before(tB) {
				return -1
			} else if tA.After(tB) {
				return 1
			}
			return 0
		}
	}

	// Numérico / Monetario / Porcentaje
	numA, isNumA := parseNumericString(strA)
	numB, isNumB := parseNumericString(strB)
	if isNumA && isNumB {
		if numA < numB {
			return -1
		} else if numA > numB {
			return 1
		}
		return 0
	}

	// Comparación léxica estándar
	return strings.Compare(strings.ToLower(strA), strings.ToLower(strB))
}

func compareAny(v1, v2 any) int {
	if v1 == nil && v2 == nil {
		return 0
	}
	if v1 == nil {
		return -1
	}
	if v2 == nil {
		return 1
	}

	switch val1 := v1.(type) {
	case int:
		if val2, ok := v2.(int); ok {
			return val1 - val2
		}
	case int64:
		if val2, ok := v2.(int64); ok {
			if val1 < val2 {
				return -1
			} else if val1 > val2 {
				return 1
			}
			return 0
		}
	case float64:
		if val2, ok := v2.(float64); ok {
			if val1 < val2 {
				return -1
			} else if val1 > val2 {
				return 1
			}
			return 0
		}
	case float32:
		if val2, ok := v2.(float32); ok {
			if val1 < val2 {
				return -1
			} else if val1 > val2 {
				return 1
			}
			return 0
		}
	case time.Time:
		if val2, ok := v2.(time.Time); ok {
			if val1.Before(val2) {
				return -1
			} else if val1.After(val2) {
				return 1
			}
			return 0
		}
	case string:
		if val2, ok := v2.(string); ok {
			return strings.Compare(strings.ToLower(val1), strings.ToLower(val2))
		}
	}
	return 0
}

func parseSpanishDate(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	// DD/MM/AAAA
	if t, err := time.Parse("02/01/2006", s); err == nil {
		return t, true
	}
	// MM / AAAA
	parts := strings.Split(s, "/")
	if len(parts) == 2 {
		mStr := strings.TrimSpace(parts[0])
		yStr := strings.TrimSpace(parts[1])
		if m, err := strconv.Atoi(mStr); err == nil {
			if y, err := strconv.Atoi(yStr); err == nil {
				return time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC), true
			}
		}
	}
	return time.Time{}, false
}

func parseNumericString(s string) (float64, bool) {
	if s == "" {
		return 0, false
	}
	clean := s
	// Quitar símbolos y sufijos comunes: €, %, m, a, etc.
	clean = strings.ReplaceAll(clean, "€", "")
	clean = strings.ReplaceAll(clean, "%", "")
	clean = strings.ReplaceAll(clean, "meses", "")
	clean = strings.ReplaceAll(clean, "mes", "")
	clean = strings.ReplaceAll(clean, " ", "")

	// Quitar "m" o "a" finales como "24 m" o "65.50 a"
	clean = strings.TrimSuffix(clean, "m")
	clean = strings.TrimSuffix(clean, "a")

	// Formato español: 1.250,50 -> 1250.50
	if strings.Contains(clean, ",") {
		clean = strings.ReplaceAll(clean, ".", "")
		clean = strings.ReplaceAll(clean, ",", ".")
	}

	val, err := strconv.ParseFloat(clean, 64)
	if err == nil {
		return val, true
	}
	return 0, false
}

// Refresh actualiza los colores del browse al conmutar el tema de la aplicación.
func (mb *ModernBrowse) Refresh() {
	if mb.headerContainer != nil {
		mb.headerContainer.Objects = []fyne.CanvasObject{mb.buildHeader()}
		mb.headerContainer.Refresh()
	}
	mb.rebuildRows()
	if mb.container != nil {
		mb.container.Refresh()
	}
	mb.BaseWidget.Refresh()
}
