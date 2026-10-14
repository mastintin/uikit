package uikit

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var spanishMonths = []string{
	"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
	"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre",
}

var spanishWeekdays = []string{"Lu", "Ma", "Mi", "Ju", "Vi", "Sá", "Do"}

// calendarDayCell representa una celda interactiva para un día específico del mes.
type calendarDayCell struct {
	widget.BaseWidget

	date       time.Time
	isCurrent  bool
	isSelected bool
	isToday    bool
	onTap      func(time.Time)

	bg        *canvas.Rectangle
	text      *canvas.Text
	root      *fyne.Container
	isHovered bool
}

func newCalendarDayCell(date time.Time, isCurrent, isSelected, isToday bool, onTap func(time.Time)) *calendarDayCell {
	c := &calendarDayCell{
		date:       date,
		isCurrent:  isCurrent,
		isSelected: isSelected,
		isToday:    isToday,
		onTap:      onTap,
	}

	c.bg = canvas.NewRectangle(color.Transparent)
	c.bg.CornerRadius = 6

	dayNum := fmt.Sprintf("%d", date.Day())
	c.text = canvas.NewText(dayNum, ThemeColor(ColorNameTextPrimary))
	c.text.TextSize = 12.5

	c.applyStyle()

	c.root = container.NewStack(c.bg, container.NewCenter(c.text))
	c.ExtendBaseWidget(c)
	return c
}

func (c *calendarDayCell) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.root)
}

func (c *calendarDayCell) MinSize() fyne.Size {
	return fyne.NewSize(34, 30)
}

func (c *calendarDayCell) applyStyle() {
	accent := ThemeColorAccent()

	if c.isSelected {
		c.bg.FillColor = accent
		c.bg.StrokeColor = color.Transparent
		c.bg.StrokeWidth = 0
		c.text.Color = color.White
		c.text.TextStyle = fyne.TextStyle{Bold: true}
	} else if c.isToday {
		c.bg.FillColor = color.Transparent
		c.bg.StrokeColor = accent
		c.bg.StrokeWidth = 1.5
		c.text.Color = accent
		c.text.TextStyle = fyne.TextStyle{Bold: true}
	} else if c.isHovered {
		c.bg.FillColor = ThemeColor(theme.ColorNameHover)
		c.bg.StrokeColor = color.Transparent
		c.bg.StrokeWidth = 0
		c.text.Color = ThemeColor(ColorNameTextPrimary)
		c.text.TextStyle = fyne.TextStyle{Bold: false}
	} else {
		c.bg.FillColor = color.Transparent
		c.bg.StrokeColor = color.Transparent
		c.bg.StrokeWidth = 0
		c.text.TextStyle = fyne.TextStyle{Bold: false}

		if c.isCurrent {
			c.text.Color = ThemeColor(ColorNameTextPrimary)
		} else {
			c.text.Color = ThemeColor(ColorNameTextMuted)
		}
	}
}

func (c *calendarDayCell) Tapped(_ *fyne.PointEvent) {
	if c.onTap != nil {
		c.onTap(c.date)
	}
}

func (c *calendarDayCell) MouseIn(_ *desktop.MouseEvent) {
	c.isHovered = true
	c.applyStyle()
	c.Refresh()
}

func (c *calendarDayCell) MouseOut() {
	c.isHovered = false
	c.applyStyle()
	c.Refresh()
}

func (c *calendarDayCell) MouseMoved(_ *desktop.MouseEvent) {}

func (c *calendarDayCell) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

// ModernCalendar es un componente visual completo de calendario estilo desplegable/almanaque
// inspirado en xlsJubila.
type ModernCalendar struct {
	widget.BaseWidget

	selectedDate time.Time
	viewMonth    time.Time
	onSelect     func(time.Time)
	onClose      func()

	lblMonthYear *canvas.Text
	gridDays     *fyne.Container
	card         *ModernCard
	popup        *widget.PopUp
}

// NewModernCalendar construye una instancia del calendario centrada en la fecha dada.
func NewModernCalendar(selected time.Time, onSelect func(time.Time), onClose func()) *ModernCalendar {
	if selected.IsZero() {
		selected = time.Now()
	}

	mc := &ModernCalendar{
		selectedDate: selected,
		viewMonth:    time.Date(selected.Year(), selected.Month(), 1, 0, 0, 0, 0, time.Local),
		onSelect:     onSelect,
		onClose:      onClose,
	}

	mc.build()
	mc.ExtendBaseWidget(mc)
	return mc
}

func (mc *ModernCalendar) build() {
	mc.lblMonthYear = canvas.NewText("", ThemeColor(ColorNameTextPrimary))
	mc.lblMonthYear.TextStyle = fyne.TextStyle{Bold: true}
	mc.lblMonthYear.TextSize = 14.5

	btnPrevYear := NewButton("Default", "«", nil, false, func() { mc.navigate(-12) }).SetHeight(28)
	btnPrevMonth := NewButton("Default", "‹", nil, false, func() { mc.navigate(-1) }).SetHeight(28)
	btnNextMonth := NewButton("Default", "›", nil, false, func() { mc.navigate(1) }).SetHeight(28)
	btnNextYear := NewButton("Default", "»", nil, false, func() { mc.navigate(12) }).SetHeight(28)

	headerNav := container.NewBorder(
		nil, nil,
		container.NewHBox(btnPrevYear, btnPrevMonth),
		container.NewHBox(btnNextMonth, btnNextYear),
		container.NewCenter(mc.lblMonthYear),
	)

	// Cabecera de días de la semana (L M X J V S D)
	weekdayRow := container.NewGridWithColumns(7)
	for _, wd := range spanishWeekdays {
		t := canvas.NewText(wd, ThemeColor(ColorNameTextMuted))
		t.TextSize = 11.5
		t.TextStyle = fyne.TextStyle{Bold: true}
		weekdayRow.Add(container.NewCenter(t))
	}

	mc.gridDays = container.NewGridWithColumns(7)
	mc.rebuildGrid()

	btnToday := NewButton("Primary", "Hoy", "calendar", false, func() {
		mc.selectDate(time.Now())
	}).SetHeight(30)

	btnClose := NewButton("Default", "Cerrar", "close", false, func() {
		mc.close()
	}).SetHeight(30)

	footer := container.NewBorder(
		nil, nil,
		btnToday,
		btnClose,
		nil,
	)

	content := container.NewVBox(
		headerNav,
		widget.NewSeparator(),
		weekdayRow,
		mc.gridDays,
		widget.NewSeparator(),
		footer,
	)

	mc.card = NewCard(content).WithElevation(2).WithVerticalPadding(10, 10).WithHorizontalPadding(12, 12)
}

func (mc *ModernCalendar) navigate(deltaMonths int) {
	mc.viewMonth = mc.viewMonth.AddDate(0, deltaMonths, 0)
	mc.rebuildGrid()
}

func (mc *ModernCalendar) rebuildGrid() {
	year := mc.viewMonth.Year()
	month := mc.viewMonth.Month()

	monthIdx := int(month) - 1
	if monthIdx >= 0 && monthIdx < len(spanishMonths) {
		mc.lblMonthYear.Text = fmt.Sprintf("%s %d", spanishMonths[monthIdx], year)
	} else {
		mc.lblMonthYear.Text = fmt.Sprintf("%s %d", month.String(), year)
	}
	mc.lblMonthYear.Color = ThemeColor(ColorNameTextPrimary)
	mc.lblMonthYear.Refresh()

	mc.gridDays.Objects = nil

	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()
	daysInPrevMonth := time.Date(year, month, 0, 0, 0, 0, 0, time.Local).Day()

	// ISO weekday donde Lunes = 0 ... Domingo = 6
	startWeekday := (int(firstOfMonth.Weekday()) + 6) % 7

	today := time.Now()

	// Cuadrícula fija de 42 celdas (6 semanas completas de 7 días)
	for i := 0; i < 42; i++ {
		var cellDate time.Time
		isCurrent := false

		if i < startWeekday {
			// Días del mes anterior
			d := daysInPrevMonth - startWeekday + 1 + i
			cellDate = time.Date(year, month-1, d, 0, 0, 0, 0, time.Local)
		} else if i < startWeekday+daysInMonth {
			// Días del mes actual
			d := i - startWeekday + 1
			cellDate = time.Date(year, month, d, 0, 0, 0, 0, time.Local)
			isCurrent = true
		} else {
			// Días del mes siguiente
			d := i - (startWeekday + daysInMonth) + 1
			cellDate = time.Date(year, month+1, d, 0, 0, 0, 0, time.Local)
		}

		isSelected := isSameDay(cellDate, mc.selectedDate)
		isToday := isSameDay(cellDate, today)

		cell := newCalendarDayCell(cellDate, isCurrent, isSelected, isToday, func(d time.Time) {
			mc.selectDate(d)
		})
		mc.gridDays.Add(cell)
	}

	mc.gridDays.Refresh()
}

func (mc *ModernCalendar) selectDate(d time.Time) {
	mc.selectedDate = d
	if mc.onSelect != nil {
		mc.onSelect(d)
	}
	mc.close()
}

func (mc *ModernCalendar) close() {
	if mc.popup != nil {
		mc.popup.Hide()
	}
	if mc.onClose != nil {
		mc.onClose()
	}
}

func (mc *ModernCalendar) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(mc.card)
}

func isSameDay(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}

// ShowCalendarPopup abre el calendario en un PopUp flotante en el lienzo del canvas dado
// o relativo al objeto objetivo (target).
func ShowCalendarPopup(canvas fyne.Canvas, target fyne.CanvasObject, initialDate time.Time, onSelect func(time.Time)) *widget.PopUp {
	if canvas == nil {
		return nil
	}

	var p *widget.PopUp
	cal := NewModernCalendar(initialDate, func(selected time.Time) {
		if onSelect != nil {
			onSelect(selected)
		}
	}, func() {
		if p != nil {
			p.Hide()
		}
	})

	p = widget.NewPopUp(cal, canvas)
	cal.popup = p

	if target != nil {
		p.ShowAtRelativePosition(fyne.NewPos(0, target.Size().Height+4), target)
	} else {
		p.Show()
	}

	return p
}
