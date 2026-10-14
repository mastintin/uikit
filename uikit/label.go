package uikit

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// Alias de tipo para máxima concisión estilo xBase / Harbour
type Label = ModernLabel

// ModernLabel es un widget tipográfico de alta fidelidad para interfaces modernas.
// Combina control cromático exacto, tamaños en puntos, negritas, alineación,
// soporte multilínea con salto de palabra (wrap), truncado y eventos de click.
type ModernLabel struct {
	widget.BaseWidget

	text       string
	textSize   float32
	textColor  color.Color
	textStyle  fyne.TextStyle
	alignment  fyne.TextAlign
	wrapping   fyne.TextWrap
	truncation fyne.TextTruncation
	onTapped   func()

	// Control interno para hover si es interactivo
	isHovered      bool
	hasCustomColor bool
	isMuted        bool
}

// NewModernLabel crea un nuevo ModernLabel con texto y estilo tipográfico opcional.
func NewModernLabel(text string, style ...TypographyStyle) *ModernLabel {
	l := &ModernLabel{
		text:           text,
		textSize:       StyleBody.Size,
		textColor:      ColorLight,
		hasCustomColor: false,
		alignment:      fyne.TextAlignLeading,
		wrapping:       fyne.TextWrapOff,
		truncation:     fyne.TextTruncateOff,
	}

	if len(style) > 0 {
		l.textSize = style[0].Size
		l.textStyle = fyne.TextStyle{
			Bold:      style[0].Bold,
			Monospace: style[0].Monospace,
		}
	}

	l.ExtendBaseWidget(l)
	return l
}

// NewModernLabelWithColor crea un ModernLabel especificando texto, estilo y color.
func NewModernLabelWithColor(text string, style TypographyStyle, col color.Color) *ModernLabel {
	l := NewModernLabel(text, style)
	if col != nil {
		l.textColor = col
		l.hasCustomColor = true
	}
	return l
}

// NewLabel es el constructor ultrarrápido y conciso de una sola línea:
// Soporta NewLabel("Texto") y NewLabel("Texto", uikit.ColorMuted)
func NewLabel(text string, col ...color.Color) *ModernLabel {
	l := NewModernLabel(text, StyleBody)
	if len(col) > 0 && col[0] != nil {
		l.textColor = col[0]
		l.hasCustomColor = true
	}
	return l
}

// --- Constructores Semánticos Listos para Usar ---

// NewHeroLabel para grandes importes monetarios o KPIs (26pt, Bold).
func NewHeroLabel(text string, col ...color.Color) *ModernLabel {
	if len(col) > 0 && col[0] != nil {
		return NewModernLabelWithColor(text, StyleHero, col[0])
	}
	l := NewModernLabel(text, StyleHero)
	l.textColor = ColorWhite
	l.hasCustomColor = false
	return l
}

// NewTitleLabel para títulos principales de vistas o tarjetas (18pt - 22pt, Bold).
func NewTitleLabel(text string, col ...color.Color) *ModernLabel {
	if len(col) > 0 && col[0] != nil {
		return NewModernLabelWithColor(text, StyleTitle1, col[0])
	}
	l := NewModernLabel(text, StyleTitle1)
	l.textColor = ColorWhite
	l.hasCustomColor = false
	return l
}

// NewSubtitleLabel para subtítulos destacados (18pt, Bold).
func NewSubtitleLabel(text string, col ...color.Color) *ModernLabel {
	if len(col) > 0 && col[0] != nil {
		return NewModernLabelWithColor(text, StyleTitle2, col[0])
	}
	l := NewModernLabel(text, StyleTitle2)
	l.textColor = ColorWhite
	l.hasCustomColor = false
	return l
}

// NewSectionLabel para cabeceras de sección o nombres de columna (14.5pt, Bold).
func NewSectionLabel(text string, col ...color.Color) *ModernLabel {
	if len(col) > 0 && col[0] != nil {
		return NewModernLabelWithColor(text, StyleSection, col[0])
	}
	l := NewModernLabel(text, StyleSection)
	l.textColor = ColorWhite
	l.hasCustomColor = false
	return l
}

// NewBodyLabel para textos descriptivos y párrafos (12.5pt, Regular).
func NewBodyLabel(text string, col ...color.Color) *ModernLabel {
	if len(col) > 0 && col[0] != nil {
		return NewModernLabelWithColor(text, StyleBody, col[0])
	}
	l := NewModernLabel(text, StyleBody)
	l.textColor = ColorLight
	l.hasCustomColor = false
	return l
}

// NewCaptionLabel para subtítulos secundarios y etiquetas silenciadas (12pt, ColorMuted).
func NewCaptionLabel(text string, col ...color.Color) *ModernLabel {
	if len(col) > 0 && col[0] != nil {
		return NewModernLabelWithColor(text, StyleCaption, col[0])
	}
	l := NewModernLabel(text, StyleCaption)
	l.textColor = ColorMuted
	l.hasCustomColor = false
	l.isMuted = true
	return l
}

// NewValueLabel para valores numéricos o monetarios de celdas y resúmenes.
func NewValueLabel(text string, bold bool, col ...color.Color) *ModernLabel {
	st := StyleValue
	if bold {
		st = StyleValueBold
	}
	if len(col) > 0 && col[0] != nil {
		return NewModernLabelWithColor(text, st, col[0])
	}
	l := NewModernLabel(text, st)
	l.textColor = ColorWhite
	l.hasCustomColor = false
	return l
}

// NewMicroLabel para subtítulos técnicos o marcas diminutas (10.5pt, Bold, ColorMuted).
func NewMicroLabel(text string, col ...color.Color) *ModernLabel {
	if len(col) > 0 && col[0] != nil {
		return NewModernLabelWithColor(text, StyleMicro, col[0])
	}
	l := NewModernLabel(text, StyleMicro)
	l.textColor = ColorMuted
	l.hasCustomColor = false
	l.isMuted = true
	return l
}

// --- Métodos Fluent / Modificadores en Cadena ---

// SetText actualiza el contenido textual del label y refresca la vista.
func (l *ModernLabel) SetText(text string) *ModernLabel {
	l.text = text
	l.Refresh()
	return l
}

// Text devuelve el texto actual.
func (l *ModernLabel) Text() string {
	return l.text
}

// SetSize define el tamaño en puntos de la fuente tipográfica.
func (l *ModernLabel) SetSize(size float32) *ModernLabel {
	l.textSize = size
	l.Refresh()
	return l
}

// SetColor asigna directamente un color.Color.
func (l *ModernLabel) SetColor(col color.Color) *ModernLabel {
	l.textColor = col
	l.hasCustomColor = (col != nil)
	l.Refresh()
	return l
}

// Color devuelve el color actual efectivo del texto.
func (l *ModernLabel) Color() color.Color {
	return l.effectiveColor()
}

func (l *ModernLabel) effectiveColor() color.Color {
	if l.hasCustomColor && l.textColor != nil {
		return l.textColor
	}
	if l.isMuted {
		return ThemeColor(ColorNameTextMuted)
	}
	if l.textColor != nil && l.textColor != ColorWhite && l.textColor != ColorLight && l.textColor != ColorMuted {
		return l.textColor
	}
	return ThemeColor(ColorNameTextPrimary)
}

// TextSize devuelve el tamaño tipográfico actual en puntos.
func (l *ModernLabel) TextSize() float32 {
	return l.textSize
}

// SetHexColor asigna un color a partir de una cadena CSS (#fff, #818cf8, etc.).
func (l *ModernLabel) SetHexColor(hex string) *ModernLabel {
	l.textColor = HexColor(hex)
	l.hasCustomColor = true
	l.Refresh()
	return l
}

// SetBold activa o desactiva la negrita.
func (l *ModernLabel) SetBold(bold bool) *ModernLabel {
	l.textStyle.Bold = bold
	l.Refresh()
	return l
}

// SetItalic activa o desactiva la cursiva.
func (l *ModernLabel) SetItalic(italic bool) *ModernLabel {
	l.textStyle.Italic = italic
	l.Refresh()
	return l
}

// SetMonospace activa o desactiva la fuente monoespaciada.
func (l *ModernLabel) SetMonospace(mono bool) *ModernLabel {
	l.textStyle.Monospace = mono
	l.Refresh()
	return l
}

// SetStyle aplica un TypographyStyle existente (tamaño, bold, monospace).
func (l *ModernLabel) SetStyle(style TypographyStyle) *ModernLabel {
	l.textSize = style.Size
	l.textStyle = fyne.TextStyle{
		Bold:      style.Bold,
		Monospace: style.Monospace,
	}
	l.Refresh()
	return l
}

// SetAlign establece la alineación horizontal del texto.
func (l *ModernLabel) SetAlign(align fyne.TextAlign) *ModernLabel {
	l.alignment = align
	l.Refresh()
	return l
}

// AlignLeft alinea el texto a la izquierda.
func (l *ModernLabel) AlignLeft() *ModernLabel {
	return l.SetAlign(fyne.TextAlignLeading)
}

// AlignCenter centra horizontalmente el texto.
func (l *ModernLabel) AlignCenter() *ModernLabel {
	return l.SetAlign(fyne.TextAlignCenter)
}

// AlignRight alinea el texto a la derecha.
func (l *ModernLabel) AlignRight() *ModernLabel {
	return l.SetAlign(fyne.TextAlignTrailing)
}

// SetWrap configura el modo de ajuste y salto de línea.
func (l *ModernLabel) SetWrap(wrap fyne.TextWrap) *ModernLabel {
	l.wrapping = wrap
	l.Refresh()
	return l
}

// WithWrap es un atajo booleano para activar ajuste por palabras (TextWrapWord).
func (l *ModernLabel) WithWrap(wrap bool) *ModernLabel {
	if wrap {
		l.wrapping = fyne.TextWrapWord
	} else {
		l.wrapping = fyne.TextWrapOff
	}
	l.Refresh()
	return l
}

// SetTruncate configura el modo de truncado de texto.
func (l *ModernLabel) SetTruncate(trunc fyne.TextTruncation) *ModernLabel {
	l.truncation = trunc
	l.Refresh()
	return l
}

// SetOnTapped asigna un callback de click (hace que el label se comporte como link o botón de texto).
func (l *ModernLabel) SetOnTapped(fn func()) *ModernLabel {
	l.onTapped = fn
	l.Refresh()
	return l
}

// --- Soporte de Eventos de Ratón e Interacción ---

func (l *ModernLabel) Tapped(e *fyne.PointEvent) {
	if l.onTapped != nil {
		l.onTapped()
	}
}

func (l *ModernLabel) Cursor() desktop.Cursor {
	if l.onTapped != nil {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

func (l *ModernLabel) MouseIn(e *desktop.MouseEvent) {
	if l.onTapped != nil {
		l.isHovered = true
		l.Refresh()
	}
}

func (l *ModernLabel) MouseOut() {
	if l.onTapped != nil && l.isHovered {
		l.isHovered = false
		l.Refresh()
	}
}

func (l *ModernLabel) MouseMoved(e *desktop.MouseEvent) {}

// --- Renderizado del Widget ---

func (l *ModernLabel) CreateRenderer() fyne.WidgetRenderer {
	r := &modernLabelRenderer{
		label:     l,
		linePool:  make([]*canvas.Text, 0, 4),
		lastWidth: -1,
	}
	r.rebuildLines(0)
	return r
}

type modernLabelRenderer struct {
	label     *ModernLabel
	linePool  []*canvas.Text
	lineTexts []string
	lastWidth float32
}

func (r *modernLabelRenderer) Objects() []fyne.CanvasObject {
	objs := make([]fyne.CanvasObject, len(r.lineTexts))
	for i := range r.lineTexts {
		objs[i] = r.linePool[i]
	}
	return objs
}

func (r *modernLabelRenderer) Destroy() {}

func (r *modernLabelRenderer) lineHeight() float32 {
	l := r.label
	sz := fyne.MeasureText("Mg", l.textSize, l.textStyle)
	h := sz.Height
	if h < l.textSize {
		h = l.textSize
	}
	return h + 2
}

func (r *modernLabelRenderer) rebuildLines(availWidth float32) {
	l := r.label
	raw := l.text
	if raw == "" {
		r.lineTexts = nil
		return
	}

	var computed []string
	if l.wrapping == fyne.TextWrapWord && availWidth > 20 {
		// Ajuste automático por palabras respetando saltos de línea existentes
		paras := strings.Split(raw, "\n")
		for _, para := range paras {
			words := strings.Fields(para)
			if len(words) == 0 {
				computed = append(computed, "")
				continue
			}

			curLine := words[0]
			for _, w := range words[1:] {
				testLine := curLine + " " + w
				testSize := fyne.MeasureText(testLine, l.textSize, l.textStyle)
				if testSize.Width <= availWidth {
					curLine = testLine
				} else {
					computed = append(computed, curLine)
					curLine = w
				}
			}
			computed = append(computed, curLine)
		}
	} else {
		// Separación estándar por saltos de línea \n
		computed = strings.Split(raw, "\n")
	}

	r.lineTexts = computed

	// Asegurar que linePool tiene suficientes objetos canvas.Text
	effCol := l.effectiveColor()
	if effCol == nil {
		effCol = ThemeColor(ColorNameTextPrimary)
	}
	for len(r.linePool) < len(r.lineTexts) {
		t := canvas.NewText("", effCol)
		r.linePool = append(r.linePool, t)
	}

	for i, txt := range r.lineTexts {
		obj := r.linePool[i]
		obj.Text = txt
		obj.TextSize = l.textSize
		obj.TextStyle = l.textStyle
		obj.Color = effCol
		obj.Alignment = l.alignment
		obj.Show()
	}

	// Ocultar objetos sobrantes en el pool
	for i := len(r.lineTexts); i < len(r.linePool); i++ {
		r.linePool[i].Hide()
	}
}

func (r *modernLabelRenderer) MinSize() fyne.Size {
	l := r.label
	if l.text == "" {
		return fyne.NewSize(0, r.lineHeight())
	}

	// En MinSize calculamos las líneas según saltos manuales \n
	lines := strings.Split(l.text, "\n")
	var maxW float32 = 0
	for _, line := range lines {
		sz := fyne.MeasureText(line, l.textSize, l.textStyle)
		if sz.Width > maxW {
			maxW = sz.Width
		}
	}

	totalH := r.lineHeight() * float32(len(lines))
	return fyne.NewSize(maxW, totalH)
}

func (r *modernLabelRenderer) Layout(size fyne.Size) {
	l := r.label
	lh := r.lineHeight()

	// Si hay wrap, reconstruir líneas en base al ancho asignado
	if l.wrapping == fyne.TextWrapWord && size.Width != r.lastWidth {
		r.lastWidth = size.Width
		r.rebuildLines(size.Width)
	}

	totalH := lh * float32(len(r.lineTexts))
	var offsetY float32 = 0
	if size.Height > totalH && len(r.lineTexts) == 1 {
		offsetY = (size.Height - totalH) / 2
	}

	for i := range r.lineTexts {
		obj := r.linePool[i]
		obj.Resize(fyne.NewSize(size.Width, lh))
		obj.Move(fyne.NewPos(0, offsetY+float32(i)*lh))
		obj.Refresh()
	}
}

func (r *modernLabelRenderer) Refresh() {
	l := r.label
	r.rebuildLines(r.lastWidth)
	effCol := l.effectiveColor()
	if effCol == nil {
		effCol = ThemeColor(ColorNameTextPrimary)
	}
	for i := range r.lineTexts {
		r.linePool[i].Color = effCol
		r.linePool[i].Refresh()
	}
	canvas.Refresh(l)
}
