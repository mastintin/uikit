package uikit

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ColorTextMuted es el color fallback #64748b (Slate 500)
var ColorTextMuted = ColorSlate500

// NewFieldLabel genera una etiqueta estilizada idéntica a .form-field label de goCuadros:
// - Fuente: 12.5px
// - Peso: Negrita / Semi-bold
// - Color: dinámico según tema (ThemeColor(ColorNameTextMuted))
func NewFieldLabel(text string) *ModernLabel {
	return NewModernLabel(text, StyleBodyBold).
		SetSize(12.5).
		SetColor(ThemeColor(ColorNameTextMuted))
}

// NewSectionTitle genera un título de sección idéntico a .section-title de goCuadros:
// - Fuente: 13.5px
// - Peso: Negrita
// - Color: dinámico según tema (ThemeColor(ColorNameTextMuted))
// - Texto en mayúsculas
func NewSectionTitle(text string) *ModernLabel {
	return NewModernLabel(strings.ToUpper(text), StyleSection).
		SetSize(13.5).
		SetColor(ThemeColor(ColorNameTextMuted))
}

// FormField empaqueta una etiqueta estilizada y un widget de entrada en disposición vertical,
// idéntico a .form-field de goCuadros (label arriba, input abajo, gap 4px).
func FormField(labelText string, input fyne.CanvasObject) *fyne.Container {
	lbl := NewFieldLabel(labelText)
	return container.NewVBox(lbl, input)
}

// TextBoxKind identifica el tipo de control interno de ModernTextBox.
type TextBoxKind int

const (
	KindStandard TextBoxKind = iota
	KindMultiLine
	KindSelectEntry
	KindSelect
)

// ModernTextBox encapsula la etiqueta superior, el control de entrada (Entry, SelectEntry o Select),
// el tema visual de alto contraste y la lógica de estado bajo una filosofía idéntica a ModernButton.
//
// Extiende widget.BaseWidget para integrarse de forma nativa en el motor de renderizado de Fyne.
type ModernTextBox struct {
	widget.BaseWidget

	Label       *ModernLabel
	Entry       *widget.Entry
	SelectEntry *widget.SelectEntry
	Select      *widget.Select
	Container   *fyne.Container
	themedInput *container.ThemeOverride

	disabled  bool
	kind      TextBoxKind
	isUpper   bool
	onChanged func(string)
}

func (m *ModernTextBox) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(m.Container)
}

// ModernTextBoxTheme proporciona la apariencia visual exacta de los campos de entrada de goCuadros:
// - Modo inactivo (readonly): fondo sombreado suave (#f8fafc), sin marco (StrokeWidth: 0), texto nítido (#334155)
// - Modo activo (edición): fondo blanco limpio (#ffffff), borde sutil (#cbd5e1), foco azul (#2563eb)
type ModernTextBoxTheme struct {
	box *ModernTextBox
}

func NewModernTextBoxTheme(box *ModernTextBox) *ModernTextBoxTheme {
	return &ModernTextBoxTheme{
		box: box,
	}
}

func (t *ModernTextBoxTheme) activeTheme() fyne.Theme {
	if app := fyne.CurrentApp(); app != nil && app.Settings() != nil {
		if th := app.Settings().Theme(); th != nil {
			return th
		}
	}
	return theme.DefaultTheme()
}

func (t *ModernTextBoxTheme) isDark(variant fyne.ThemeVariant) bool {
	if IsDarkTheme() || variant == theme.VariantDark {
		return true
	}
	if app := fyne.CurrentApp(); app != nil && app.Settings() != nil {
		if appTh, ok := app.Settings().Theme().(*AppTheme); ok && appTh.ForcedVariant() != nil {
			return *appTh.ForcedVariant() == theme.VariantDark
		}
		return app.Settings().ThemeVariant() == theme.VariantDark
	}
	return false
}

func (t *ModernTextBoxTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	dark := t.isDark(variant)
	effectiveVariant := theme.VariantLight
	if dark {
		effectiveVariant = theme.VariantDark
	}

	if t.box != nil && t.box.Disabled() {
		switch name {
		case theme.ColorNameInputBackground, theme.ColorNameDisabledButton:
			if dark {
				return ColorSlate800
			}
			// Sombreado de fondo (#f8fafc / Slate 50)
			return ColorSlate50

		case theme.ColorNameDisabled, theme.ColorNameForeground:
			if dark {
				return ColorSlate300
			}
			// Texto oscuro y nítido (#334155 / Slate 700)
			return ColorSlate700
		}
	} else {
		// Modo activo / edición
		switch name {
		case theme.ColorNameInputBackground, theme.ColorNameButton:
			if dark {
				return ColorSlate900
			}
			return ColorWhite

		case theme.ColorNameInputBorder:
			if dark {
				return ColorSlate600
			}
			return ColorSlate300

		case theme.ColorNameForeground:
			if dark {
				return ColorSlate50
			}
			return ColorSlate900

		case theme.ColorNamePlaceHolder:
			if dark {
				return ColorSlate500
			}
			return ColorSlate400

		case theme.ColorNamePrimary:
			if dark {
				return ColorBlue500
			}
			return ColorBlue600
		}
	}

	return t.activeTheme().Color(name, effectiveVariant)
}

func (t *ModernTextBoxTheme) Font(style fyne.TextStyle) fyne.Resource {
	return t.activeTheme().Font(style)
}

func (t *ModernTextBoxTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.activeTheme().Icon(name)
}

func (t *ModernTextBoxTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameInputBorder:
		if t.box != nil && t.box.Disabled() {
			// Sin marco cuando no está activo: el sombreado delimita el campo
			return 0
		}
		return 1.0
	case theme.SizeNamePadding:
		// Padding compacto idéntico a inputs web/goCuadros (4px) para evitar desbordamiento vertical
		return 4.0
	case theme.SizeNameInputRadius:
		return 6.0
	}
	return t.activeTheme().Size(name)
}

// NewTextBox crea un campo de texto estándar de una línea encapsulado con su etiqueta.
// label: texto de la etiqueta superior (ej. "Ubicación física")
// placeholder: sugerencia visual atenuada (ej. "ej: Pasillo 2")
// initialValue: valor inicial del campo
// disabled: true para iniciar en modo solo lectura
// onChanged: callback ejecutado en cada modificación
func NewTextBox(label, placeholder, initialValue string, disabled bool, onChanged func(string)) *ModernTextBox {
	e := widget.NewEntry()
	e.SetPlaceHolder(placeholder)
	e.SetText(initialValue)

	m := &ModernTextBox{
		Label:     NewFieldLabel(label),
		Entry:     e,
		kind:      KindStandard,
		disabled:  disabled,
		onChanged: onChanged,
	}

	m.setupEntryHandler()
	m.themedInput = container.NewThemeOverride(e, NewModernTextBoxTheme(m))
	m.Container = container.NewVBox(m.Label, m.themedInput)
	m.ExtendBaseWidget(m)

	if disabled {
		m.Disable()
	}
	return m
}

type inlineInputWrapper struct {
	widget.BaseWidget
	content fyne.CanvasObject
	width   float32
}

func (w *inlineInputWrapper) CreateRenderer() fyne.WidgetRenderer {
	return &inlineInputWrapperRenderer{wrapper: w}
}

type inlineInputWrapperRenderer struct {
	wrapper *inlineInputWrapper
}

func (r *inlineInputWrapperRenderer) Layout(size fyne.Size) {
	if r.wrapper.content != nil {
		r.wrapper.content.Resize(size)
		r.wrapper.content.Move(fyne.NewPos(0, 0))
	}
}

func (r *inlineInputWrapperRenderer) MinSize() fyne.Size {
	var h float32 = 36
	if r.wrapper.content != nil {
		ms := r.wrapper.content.MinSize()
		h = ms.Height
	}
	targetW := r.wrapper.width
	if targetW <= 0 {
		targetW = 80
	}
	return fyne.NewSize(targetW, h)
}

func (r *inlineInputWrapperRenderer) Refresh() {
	if r.wrapper.content != nil {
		r.wrapper.content.Refresh()
	}
}

func (r *inlineInputWrapperRenderer) Objects() []fyne.CanvasObject {
	if r.wrapper.content != nil {
		return []fyne.CanvasObject{r.wrapper.content}
	}
	return nil
}

func (r *inlineInputWrapperRenderer) Destroy() {}

// NewInlineInput crea un campo de entrada horizontal estilizado con ModernTextBoxTheme
// con ancho explícito configurable y altura nativa perfecta para que no se recorte ni deslice.
// Devuelve el widget.Entry junto con su objeto CanvasObject envolvente.
func NewInlineInput(placeholder, initialValue string, width float32, onChanged func(string)) (*widget.Entry, fyne.CanvasObject) {
	e := widget.NewEntry()
	e.SetPlaceHolder(placeholder)
	e.SetText(initialValue)
	if onChanged != nil {
		e.OnChanged = onChanged
	}

	m := &ModernTextBox{
		Entry:     e,
		kind:      KindStandard,
		disabled:  false,
		onChanged: onChanged,
	}
	m.setupEntryHandler()

	thOverride := container.NewThemeOverride(e, NewModernTextBoxTheme(m))
	wrapper := &inlineInputWrapper{
		content: thOverride,
		width:   width,
	}
	wrapper.ExtendBaseWidget(wrapper)
	return e, wrapper
}

// NewUpperTextBox crea un campo de texto que fuerza automáticamente mayúsculas al teclear.
func NewUpperTextBox(label, placeholder, initialValue string, disabled bool, onChanged func(string)) *ModernTextBox {
	m := NewTextBox(label, placeholder, initialValue, disabled, onChanged)
	m.WithUpper()
	return m
}

// NewTextBoxWithAction crea un campo de texto con un botón de acción lateral (ej. selector de archivos/carpetas).
func NewTextBoxWithAction(
	label string,
	placeholder string,
	initialValue string,
	disabled bool,
	actionIcon string,
	onAction func(),
) *ModernTextBox {
	e := widget.NewEntry()
	e.SetPlaceHolder(placeholder)
	e.SetText(initialValue)

	m := &ModernTextBox{
		Label:    NewFieldLabel(label),
		Entry:    e,
		kind:     KindStandard,
		disabled: disabled,
	}

	m.setupEntryHandler()
	m.themedInput = container.NewThemeOverride(e, NewModernTextBoxTheme(m))

	btnAction := NewButton("Icon", "", actionIcon, disabled, onAction)
	inputRow := container.NewBorder(nil, nil, nil, btnAction, m.themedInput)

	m.Container = container.NewVBox(m.Label, inputRow)
	m.ExtendBaseWidget(m)

	if disabled {
		m.Disable()
	}
	return m
}

// NewMultiLineTextBox crea un área de texto multilínea con altura configurada en filas visibles.
func NewMultiLineTextBox(label, placeholder, initialValue string, rows int, disabled bool, onChanged func(string)) *ModernTextBox {
	e := widget.NewMultiLineEntry()
	e.SetPlaceHolder(placeholder)
	e.SetText(initialValue)
	if rows > 1 {
		e.SetMinRowsVisible(rows)
	}

	m := &ModernTextBox{
		Label:     NewFieldLabel(label),
		Entry:     e,
		kind:      KindMultiLine,
		disabled:  disabled,
		onChanged: onChanged,
	}

	m.setupEntryHandler()
	m.themedInput = container.NewThemeOverride(e, NewModernTextBoxTheme(m))
	m.Container = container.NewVBox(m.Label, m.themedInput)
	m.ExtendBaseWidget(m)

	if disabled {
		m.Disable()
	}
	return m
}

// NewSelectEntryBox crea una caja de texto editable con autocompletado y sugerencias.
func NewSelectEntryBox(label string, options []string, initialValue string, disabled bool, onChanged func(string)) *ModernTextBox {
	se := widget.NewSelectEntry(options)
	se.SetText(initialValue)

	m := &ModernTextBox{
		Label:       NewFieldLabel(label),
		SelectEntry: se,
		kind:        KindSelectEntry,
		disabled:    disabled,
		onChanged:   onChanged,
	}

	se.OnChanged = func(s string) {
		if m.isUpper {
			upper := strings.ToUpper(s)
			if upper != s {
				se.SetText(upper)
				return
			}
		}
		if m.onChanged != nil {
			m.onChanged(s)
		}
	}

	m.themedInput = container.NewThemeOverride(se, NewModernTextBoxTheme(m))
	m.Container = container.NewVBox(m.Label, m.themedInput)
	m.ExtendBaseWidget(m)

	if disabled {
		m.Disable()
	}
	return m
}

// NewSelectBox crea un desplegable de selección fija.
func NewSelectBox(label string, options []string, selected string, disabled bool, onChanged func(string)) *ModernTextBox {
	s := widget.NewSelect(options, onChanged)
	if selected != "" {
		s.SetSelected(selected)
	}

	m := &ModernTextBox{
		Label:     NewFieldLabel(label),
		Select:    s,
		kind:      KindSelect,
		disabled:  disabled,
		onChanged: onChanged,
	}

	m.themedInput = container.NewThemeOverride(s, NewModernTextBoxTheme(m))
	m.Container = container.NewVBox(m.Label, m.themedInput)
	m.ExtendBaseWidget(m)

	if disabled {
		m.Disable()
	}
	return m
}

func (m *ModernTextBox) setupEntryHandler() {
	if m.Entry == nil {
		return
	}
	m.Entry.OnChanged = func(s string) {
		if m.isUpper {
			upper := strings.ToUpper(s)
			if upper != s {
				m.Entry.SetText(upper)
				return
			}
		}
		if m.onChanged != nil {
			m.onChanged(s)
		}
	}
}

// Métodos fluidos y de configuración

func (m *ModernTextBox) WithUpper() *ModernTextBox {
	m.isUpper = true
	if m.Entry != nil && m.Entry.Text != "" {
		m.Entry.SetText(strings.ToUpper(m.Entry.Text))
	} else if m.SelectEntry != nil && m.SelectEntry.Text != "" {
		m.SelectEntry.SetText(strings.ToUpper(m.SelectEntry.Text))
	}
	return m
}

func (m *ModernTextBox) WithPlaceholder(ph string) *ModernTextBox {
	if m.Entry != nil {
		m.Entry.SetPlaceHolder(ph)
	} else if m.SelectEntry != nil {
		m.SelectEntry.SetPlaceHolder(ph)
	}
	return m
}

func (m *ModernTextBox) WithDisabled(disabled bool) *ModernTextBox {
	m.SetEnabled(!disabled)
	return m
}

// Métodos de lectura y asignación de valor

func (m *ModernTextBox) Text() string {
	if m.Entry != nil {
		return m.Entry.Text
	}
	if m.SelectEntry != nil {
		return m.SelectEntry.Text
	}
	if m.Select != nil {
		return m.Select.Selected
	}
	return ""
}

func (m *ModernTextBox) SetText(val string) *ModernTextBox {
	if m.Entry != nil {
		m.Entry.SetText(val)
	} else if m.SelectEntry != nil {
		m.SelectEntry.SetText(val)
	} else if m.Select != nil {
		m.Select.SetSelected(val)
	}
	return m
}

func (m *ModernTextBox) Float() float64 {
	s := strings.TrimSpace(m.Text())
	s = strings.ReplaceAll(s, ",", ".")
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

func (m *ModernTextBox) SetFloat(val float64, decimals int) *ModernTextBox {
	if val <= 0 {
		m.SetText("")
	} else {
		format := fmt.Sprintf("%%.%df", decimals)
		m.SetText(fmt.Sprintf(format, val))
	}
	return m
}

func (m *ModernTextBox) Int() int {
	s := strings.TrimSpace(m.Text())
	val, _ := strconv.Atoi(s)
	return val
}

func (m *ModernTextBox) SetInt(val int) *ModernTextBox {
	if val <= 0 {
		m.SetText("")
	} else {
		m.SetText(strconv.Itoa(val))
	}
	return m
}

func (m *ModernTextBox) SetOptions(options []string) *ModernTextBox {
	if m.SelectEntry != nil {
		m.SelectEntry.SetOptions(options)
	} else if m.Select != nil {
		m.Select.Options = options
		m.Select.Refresh()
	}
	return m
}

func (m *ModernTextBox) SetLabel(text string) *ModernTextBox {
	if m.Label != nil {
		m.Label.SetText(text)
	}
	return m
}

// Métodos de habilitación y estado

func (m *ModernTextBox) Disabled() bool {
	if m.Entry != nil {
		return m.Entry.Disabled()
	}
	if m.SelectEntry != nil {
		return m.SelectEntry.Disabled()
	}
	if m.Select != nil {
		return m.Select.Disabled()
	}
	return m.disabled
}

func triggerWidgetLayout(w fyne.Widget) {
	if w == nil {
		return
	}
	sz := w.Size()
	w.Resize(sz.AddWidthHeight(0, 0.1))
	w.Resize(sz)
}

func (m *ModernTextBox) Enable() *ModernTextBox {
	m.disabled = false
	if m.Entry != nil {
		m.Entry.Enable()
		triggerWidgetLayout(m.Entry)
	}
	if m.SelectEntry != nil {
		m.SelectEntry.Enable()
		triggerWidgetLayout(m.SelectEntry)
	}
	if m.Select != nil {
		m.Select.Enable()
		triggerWidgetLayout(m.Select)
	}
	if m.themedInput != nil {
		m.themedInput.Refresh()
	}
	m.Refresh()
	return m
}

func (m *ModernTextBox) Disable() *ModernTextBox {
	m.disabled = true
	if m.Entry != nil {
		m.Entry.Disable()
		triggerWidgetLayout(m.Entry)
	}
	if m.SelectEntry != nil {
		m.SelectEntry.Disable()
		triggerWidgetLayout(m.SelectEntry)
	}
	if m.Select != nil {
		m.Select.Disable()
		triggerWidgetLayout(m.Select)
	}
	if m.themedInput != nil {
		m.themedInput.Refresh()
	}
	m.Refresh()
	return m
}

func (m *ModernTextBox) SetEnabled(enabled bool) *ModernTextBox {
	if enabled {
		return m.Enable()
	}
	return m.Disable()
}

// Refresh actualiza los colores y estilos del campo de texto y su etiqueta según el tema activo.
func (m *ModernTextBox) Refresh() {
	if m.Label != nil {
		m.Label.SetColor(ThemeColor(ColorNameTextMuted))
	}
	if m.themedInput != nil {
		m.themedInput.Refresh()
	}
	m.BaseWidget.Refresh()
}

// FieldGroup representa un conjunto de campos para realizar operaciones por lote
// como habilitar/deshabilitar o limpiar en una sola instrucción.
type FieldGroup []*ModernTextBox

func (g FieldGroup) SetEnabled(enabled bool) {
	for _, f := range g {
		f.SetEnabled(enabled)
	}
}

func (g FieldGroup) Clear() {
	for _, f := range g {
		f.SetText("")
	}
}

func (g FieldGroup) Refresh() {
	for _, f := range g {
		f.Refresh()
	}
}

