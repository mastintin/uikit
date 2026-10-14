package uikit

import (
	"image/color"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ModernDatePicker es un control avanzado de selección de fecha que combina:
// 1. Un campo de texto formateado con icono interactivo de almanaque (ActionItem)
//    que despliega un calendario completo estilo xlsJubila al hacer clic.
// 2. Cuatro botones opcionales de navegación rápida (−1a, −1m, +1m, +1a),
//    controlados mediante la variable ShowNavButtons (por defecto activa: true).
//
// # Callbacks disponibles
//
//   - onChanged   → se dispara cuando la fecha cambia por botón, Enter o selección en el calendario.
//   - onChangedLive → se dispara también mientras el usuario teclea, en cuanto el texto
//     alcanza una longitud de 10 caracteres y resulta parseable.
//
// Uso básico:
//
//	dp := uikit.NewModernDatePicker(time.Now(), func(t time.Time) {
//	    fmt.Println("Nueva fecha:", t)
//	})
//	dp.SetShowNavButtons(false) // oculta los botones + - si solo se desea el textbox con almanaque
type ModernDatePicker struct {
	widget.BaseWidget

	// Container es el objeto raíz que se debe incluir en la UI.
	Container *fyne.Container

	// ShowNavButtons determina si se muestran los botones de navegación (−1a, −1m, +1m, +1a).
	// Por defecto es true.
	ShowNavButtons bool

	// onChanged se invoca al confirmar la fecha (botón de navegación, Enter o calendario).
	onChanged func(time.Time)

	// onChangedLive se invoca también al teclear, en cuanto el texto es parseable
	// con longitud exacta de 10 caracteres. Útil para recalcular en vivo.
	onChangedLive func(time.Time)

	currentDate time.Time
	entryDate   *widget.Entry
	height      float32 // altura compartida para botones y entry (default 38)
	entryWidth  float32 // ancho del campo de texto de fecha (default 130)
	window      fyne.Window

	btnPrevYear   *ModernButton
	btnPrevMonth  *ModernButton
	btnNextMonth  *ModernButton
	btnNextYear   *ModernButton
	btnCalendar   *calendarActionItem
	calendarPopup *widget.PopUp

	// isUpdating bloquea los callbacks de OnChanged/OnChangedLive mientras el widget
	// escribe programáticamente en el campo de texto (SetDate, etc.).
	isUpdating bool
}

var (
	activeDatePickersMu sync.Mutex
	activeDatePickers   []*ModernDatePicker
)

func registerDatePicker(dp *ModernDatePicker) {
	activeDatePickersMu.Lock()
	defer activeDatePickersMu.Unlock()
	activeDatePickers = append(activeDatePickers, dp)
}

// RefreshAllDatePickers actualiza inmediatamente el tema visual de todos los selectores de fecha de la app.
func RefreshAllDatePickers() {
	activeDatePickersMu.Lock()
	defer activeDatePickersMu.Unlock()
	for _, dp := range activeDatePickers {
		if dp != nil {
			dp.Refresh()
		}
	}
}

// NewModernDatePicker crea un ModernDatePicker inicializado con la fecha dada.
// Si initial es zero-value, se usa time.Now().
// onChanged puede ser nil.
func NewModernDatePicker(initial time.Time, onChanged func(time.Time)) *ModernDatePicker {
	if initial.IsZero() {
		initial = time.Now()
	}

	dp := &ModernDatePicker{
		onChanged:      onChanged,
		currentDate:    initial,
		height:         38,
		entryWidth:     130,
		ShowNavButtons: true,
	}
	dp.ExtendBaseWidget(dp)
	dp.build()
	registerDatePicker(dp)
	return dp
}

// build construye el layout interno (llamado una sola vez en el constructor).
func (dp *ModernDatePicker) build() {
	h := dp.height

	dp.btnPrevYear = NewButton("Default", "−1a", nil, false, func() { dp.adjustMonths(-12) }).SetHeight(h)
	dp.btnPrevMonth = NewButton("Default", "−1m", "chevron-left", false, func() { dp.adjustMonths(-1) }).SetHeight(h).SetPaddingX(20)
	dp.btnNextMonth = NewButton("Default", "+1m", "chevron-right", false, func() { dp.adjustMonths(1) }).SetHeight(h).SetPaddingX(20)
	dp.btnNextYear = NewButton("Default", "+1a", nil, false, func() { dp.adjustMonths(12) }).SetHeight(h)

	dp.entryDate = widget.NewEntry()
	dp.entryDate.TextStyle = fyne.TextStyle{Bold: true}
	dp.entryDate.PlaceHolder = "DD/MM/AAAA"

	// Icono de almanaque incorporado dentro del textbox con margen de seguridad anti-solapamiento
	dp.btnCalendar = newCalendarActionItem(func() {
		dp.ShowCalendar()
	})
	dp.entryDate.ActionItem = dp.btnCalendar

	dp.setEntryText(formatDateSpanishDP(dp.currentDate.Format("2006-01-02")))

	// OnChanged: recalculo en vivo cuando el texto es una fecha válida (8-10 chars).
	dp.entryDate.OnChanged = func(s string) {
		if dp.isUpdating {
			return
		}
		trimmed := strings.TrimSpace(s)
		if len(trimmed) >= 8 && len(trimmed) <= 10 {
			if parsed, ok := parseDate(trimmed); ok {
				newISO := parsed.Format("2006-01-02")
				if newISO != dp.currentDate.Format("2006-01-02") {
					dp.currentDate = parsed
					dp.notifyLive()
				}
			}
		}
	}

	// OnSubmitted: confirmar con Enter normaliza el texto y dispara onChanged.
	dp.entryDate.OnSubmitted = func(text string) {
		text = strings.TrimSpace(text)
		if text == "" {
			return
		}
		parsed, ok := parseDate(text)
		if !ok {
			return
		}
		dp.currentDate = parsed
		dp.setEntryText(formatDateSpanishDP(dp.currentDate.Format("2006-01-02")))
		dp.notify()
	}

	dp.Container = container.NewHBox()
	dp.rebuildContainer()
}

func (dp *ModernDatePicker) rebuildContainer() {
	if dp.Container == nil {
		dp.Container = container.NewHBox()
	}
	h := dp.height
	ew := dp.entryWidth
	if ew <= 0 {
		ew = 130
	}
	entryWrap := container.NewGridWrap(fyne.NewSize(ew, h), dp.entryDate)

	if dp.ShowNavButtons {
		dp.Container.Objects = []fyne.CanvasObject{
			dp.btnPrevYear,
			dp.btnPrevMonth,
			entryWrap,
			dp.btnNextMonth,
			dp.btnNextYear,
		}
	} else {
		dp.Container.Objects = []fyne.CanvasObject{
			entryWrap,
		}
	}
	dp.Container.Refresh()
}

type modernDatePickerRenderer struct {
	dp *ModernDatePicker
}

func (r *modernDatePickerRenderer) Destroy() {}
func (r *modernDatePickerRenderer) Objects() []fyne.CanvasObject {
	if r.dp.Container != nil {
		return []fyne.CanvasObject{r.dp.Container}
	}
	return nil
}
func (r *modernDatePickerRenderer) MinSize() fyne.Size {
	if r.dp.Container != nil {
		return r.dp.Container.MinSize()
	}
	return fyne.NewSize(130, 38)
}
func (r *modernDatePickerRenderer) Layout(size fyne.Size) {
	if r.dp.Container != nil {
		r.dp.Container.Resize(size)
		r.dp.Container.Move(fyne.NewPos(0, 0))
	}
}
func (r *modernDatePickerRenderer) Refresh() {
	if r.dp != nil {
		r.dp.refreshButtons()
	}
}

func (dp *ModernDatePicker) CreateRenderer() fyne.WidgetRenderer {
	return &modernDatePickerRenderer{dp: dp}
}

// ── API pública ───────────────────────────────────────────────────────────────

// SetShowNavButtons permite activar o desactivar los botones de navegación de mes y año (−1a, −1m, +1m, +1a).
// Si es false, únicamente se visualiza el campo de texto con su icono de almanaque. Por defecto es true.
func (dp *ModernDatePicker) SetShowNavButtons(show bool) *ModernDatePicker {
	dp.ShowNavButtons = show
	dp.rebuildContainer()
	return dp
}

// WithNavButtons es un alias semántico de SetShowNavButtons para encadenamiento fluido.
func (dp *ModernDatePicker) WithNavButtons(show bool) *ModernDatePicker {
	return dp.SetShowNavButtons(show)
}

// SetWindow asigna la ventana padre para el despliegue del popup del calendario.
func (dp *ModernDatePicker) SetWindow(w fyne.Window) *ModernDatePicker {
	dp.window = w
	return dp
}

// ShowCalendar despliega el calendario interactivo completo estilo xlsJubila.
func (dp *ModernDatePicker) ShowCalendar() {
	c := dp.getCanvas()
	if c == nil {
		return
	}

	if dp.calendarPopup != nil {
		dp.calendarPopup.Hide()
		dp.calendarPopup = nil
	}

	cur := dp.Date()
	dp.calendarPopup = ShowCalendarPopup(c, dp.entryDate, cur, func(selected time.Time) {
		dp.currentDate = selected
		dp.setEntryText(formatDateSpanishDP(selected.Format("2006-01-02")))
		dp.notify()
		dp.calendarPopup = nil
	})
}

func (dp *ModernDatePicker) getCanvas() fyne.Canvas {
	if dp.window != nil {
		return dp.window.Canvas()
	}
	if app := fyne.CurrentApp(); app != nil && app.Driver() != nil {
		if c := app.Driver().CanvasForObject(dp.entryDate); c != nil {
			return c
		}
		if c := app.Driver().CanvasForObject(dp.Container); c != nil {
			return c
		}
		wins := app.Driver().AllWindows()
		if len(wins) > 0 {
			return wins[0].Canvas()
		}
	}
	return nil
}

// Date devuelve la fecha actualmente seleccionada.
// Si el usuario ha editado directamente el campo de texto, lee el valor parseado del campo.
func (dp *ModernDatePicker) Date() time.Time {
	if dp.entryDate != nil {
		txt := strings.TrimSpace(dp.entryDate.Text)
		if txt != "" {
			if parsed, ok := parseDate(txt); ok {
				dp.currentDate = parsed
			}
		}
	}
	return dp.currentDate
}

// DateString devuelve la fecha en formato ISO "2006-01-02".
func (dp *ModernDatePicker) DateString() string {
	return dp.Date().Format("2006-01-02")
}

// SetDate establece la fecha programáticamente sin disparar ningún callback.
func (dp *ModernDatePicker) SetDate(t time.Time) *ModernDatePicker {
	if t.IsZero() {
		return dp
	}
	dp.currentDate = t
	dp.setEntryText(formatDateSpanishDP(dp.currentDate.Format("2006-01-02")))
	return dp
}

// SetDateString acepta un string en cualquiera de los tres formatos soportados
// ("DD/MM/AAAA", "DD-MM-AAAA", "AAAA-MM-DD") sin disparar ningún callback.
// Devuelve false si el formato no es reconocido.
func (dp *ModernDatePicker) SetDateString(s string) bool {
	t, ok := parseDate(strings.TrimSpace(s))
	if !ok {
		return false
	}
	dp.SetDate(t)
	return true
}

// SetHeight permite ajustar la altura compartida de botones y entry en cadena.
// Debe llamarse antes de añadir el Container a la UI para que tenga efecto visual.
func (dp *ModernDatePicker) SetHeight(h float32) *ModernDatePicker {
	dp.height = h
	dp.rebuildContainer()
	return dp
}

// SetEntryWidth permite ajustar el ancho del campo de texto de fecha (por defecto 130px).
func (dp *ModernDatePicker) SetEntryWidth(w float32) *ModernDatePicker {
	dp.entryWidth = w
	dp.rebuildContainer()
	return dp
}

// SetOnChanged reemplaza el callback que se dispara al confirmar la fecha.
func (dp *ModernDatePicker) SetOnChanged(fn func(time.Time)) *ModernDatePicker {
	dp.onChanged = fn
	return dp
}

// SetOnChangedLive establece un callback adicional que se dispara mientras el usuario
// teclea, en cuanto el campo contiene exactamente 10 caracteres y forman una fecha
// válida. Ideal para recálculos inmediatos sin necesitar pulsar Enter.
func (dp *ModernDatePicker) SetOnChangedLive(fn func(time.Time)) *ModernDatePicker {
	dp.onChangedLive = fn
	return dp
}

func (dp *ModernDatePicker) refreshButtons() {
	if dp.btnPrevYear != nil {
		dp.btnPrevYear.Refresh()
	}
	if dp.btnPrevMonth != nil {
		dp.btnPrevMonth.Refresh()
	}
	if dp.btnNextMonth != nil {
		dp.btnNextMonth.Refresh()
	}
	if dp.btnNextYear != nil {
		dp.btnNextYear.Refresh()
	}
	if dp.entryDate != nil {
		dp.entryDate.Refresh()
	}
	if dp.btnCalendar != nil {
		dp.btnCalendar.Refresh()
	}
	if dp.Container != nil {
		dp.Container.Refresh()
	}
}

func (dp *ModernDatePicker) Refresh() {
	dp.refreshButtons()
	dp.BaseWidget.Refresh()
}

// ── Navegación interna ────────────────────────────────────────────────────────

func (dp *ModernDatePicker) adjustMonths(delta int) {
	// Si el usuario ha editado el texto manualmente, leerlo antes de ajustar.
	if text := strings.TrimSpace(dp.entryDate.Text); text != "" {
		if t, ok := parseDate(text); ok {
			dp.currentDate = t
		}
	}
	dp.currentDate = dp.currentDate.AddDate(0, delta, 0)
	dp.setEntryText(formatDateSpanishDP(dp.currentDate.Format("2006-01-02")))
	dp.notify()
}

// setEntryText escribe en el entry con la bandera isUpdating activa para suprimir
// los callbacks de OnChanged durante actualizaciones programáticas.
func (dp *ModernDatePicker) setEntryText(text string) {
	dp.isUpdating = true
	dp.entryDate.SetText(text)
	dp.isUpdating = false
}

func (dp *ModernDatePicker) notify() {
	if dp.onChanged != nil {
		dp.onChanged(dp.currentDate)
	}
	// Los botones de navegación también disparan el live callback.
	if dp.onChangedLive != nil {
		dp.onChangedLive(dp.currentDate)
	}
}

func (dp *ModernDatePicker) notifyLive() {
	if dp.onChangedLive != nil {
		dp.onChangedLive(dp.currentDate)
	}
}

// ── Helpers privados ──────────────────────────────────────────────────────────

// parseDate admite DD/MM/AAAA, D/M/AAAA, DD-MM-AAAA, D-M-AAAA y AAAA-MM-DD.
func parseDate(s string) (time.Time, bool) {
	for _, layout := range []string{"02/01/2006", "2/1/2006", "02-01-2006", "2-1-2006", "2006-01-02", "2006/01/02"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// formatDateSpanishDP convierte "AAAA-MM-DD" → "DD/MM/AAAA".
// Es una copia local para evitar dependencia circular con el paquete ui.
func formatDateSpanishDP(iso string) string {
	t, err := time.Parse("2006-01-02", iso)
	if err != nil {
		return iso
	}
	return t.Format("02/01/2006")
}

// calendarActionItem es un control interactivo estilizado y compacto para el ActionItem
// del textbox, con márgenes de seguridad para no tocar ni superponerse sobre el marco.
type calendarActionItem struct {
	widget.BaseWidget

	img       *canvas.Image
	bgHover   *canvas.Rectangle
	isHovered bool
	onTap     func()
}

func newCalendarActionItem(onTap func()) *calendarActionItem {
	item := &calendarActionItem{
		onTap: onTap,
	}

	calColor := ColorToHex(ThemeColorSuccess())
	if calColor == "" {
		calColor = "#107c41"
	}

	item.img = canvas.NewImageFromResource(IconResource("calendar", calColor))
	item.img.FillMode = canvas.ImageFillContain

	item.bgHover = canvas.NewRectangle(color.Transparent)
	item.bgHover.CornerRadius = 4

	item.ExtendBaseWidget(item)
	return item
}

func (a *calendarActionItem) CreateRenderer() fyne.WidgetRenderer {
	return &calendarActionItemRenderer{
		item:    a,
		objects: []fyne.CanvasObject{a.bgHover, a.img},
	}
}

func (a *calendarActionItem) MinSize() fyne.Size {
	return fyne.NewSize(24, 24)
}

func (a *calendarActionItem) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (a *calendarActionItem) Tapped(_ *fyne.PointEvent) {
	if a.onTap != nil {
		a.onTap()
	}
}

func (a *calendarActionItem) MouseIn(_ *desktop.MouseEvent) {
	a.isHovered = true
	a.applyStyle()
	a.Refresh()
}

func (a *calendarActionItem) MouseOut() {
	a.isHovered = false
	a.applyStyle()
	a.Refresh()
}

func (a *calendarActionItem) MouseMoved(_ *desktop.MouseEvent) {}

func (a *calendarActionItem) applyStyle() {
	if a.isHovered {
		a.bgHover.FillColor = ThemeColor(theme.ColorNameHover)
	} else {
		a.bgHover.FillColor = color.Transparent
	}

	calColor := ColorToHex(ThemeColorSuccess())
	if calColor == "" {
		calColor = "#107c41"
	}
	a.img.Resource = IconResource("calendar", calColor)
	a.img.Refresh()
}

type calendarActionItemRenderer struct {
	item    *calendarActionItem
	objects []fyne.CanvasObject
}

func (r *calendarActionItemRenderer) Destroy() {}
func (r *calendarActionItemRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *calendarActionItemRenderer) MinSize() fyne.Size           { return r.item.MinSize() }

func (r *calendarActionItemRenderer) Layout(size fyne.Size) {
	// Deja un margen de 2px en horizontal y 3px en vertical para que la pastilla de hover
	// quede completamente dentro del marco sin tocar nunca el borde del entry.
	const padX float32 = 2
	const padY float32 = 3
	bgW := size.Width - padX*2
	bgH := size.Height - padY*2
	if bgW < 14 {
		bgW = 14
	}
	if bgH < 14 {
		bgH = 14
	}

	r.item.bgHover.Resize(fyne.NewSize(bgW, bgH))
	r.item.bgHover.Move(fyne.NewPos((size.Width-bgW)/2, (size.Height-bgH)/2))

	const iconSize float32 = 16
	r.item.img.Resize(fyne.NewSquareSize(iconSize))
	r.item.img.Move(fyne.NewPos((size.Width-iconSize)/2, (size.Height-iconSize)/2))
}

func (r *calendarActionItemRenderer) Refresh() {
	r.item.applyStyle()
	canvas.Refresh(r.item.bgHover)
	canvas.Refresh(r.item.img)
}

