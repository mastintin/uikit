package uikit

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// SwitchVariant define las variantes semánticas de color para el interruptor cuando está activo.
type SwitchVariant int

const (
	SwitchPrimary SwitchVariant = iota // Azul corporativo (por defecto)
	SwitchSuccess                      // Verde esmeralda (activo / OK)
	SwitchDanger                       // Rojo (alerta / peligro)
	SwitchWarning                      // Ámbar (precaución)
)

// SwitchLabelPosition define la posición relativa de la etiqueta respecto al control deslizante.
type SwitchLabelPosition int

const (
	LabelTrailing SwitchLabelPosition = iota // [Switch] [Texto] (por defecto)
	LabelLeading                             // [Texto] [Switch]
)

// Alias para mantener la misma filosofía de tipos unificados del UIKit.
type Switch = ModernSwitch

const (
	defaultTrackWidth  float32 = 38
	defaultTrackHeight float32 = 22
	defaultThumbPad    float32 = 2
	defaultLabelGap    float32 = 8
)

// ModernSwitch es un interruptor digital animado de alta precisión estética (estilo iOS / Tailwind),
// con soporte completo para transiciones fluidas, hover con cursor puntero y halo reactivo,
// sombra tridimensional en el botón, delimitación estricta de hitbox, accesibilidad por teclado
// (espacio y enter), etiqueta de texto integrada opcional (a izquierda o derecha) y variantes semánticas.
type ModernSwitch struct {
	widget.DisableableWidget

	Text          string
	On            bool
	Variant       SwitchVariant
	LabelPosition SwitchLabelPosition
	OnChanged     func(on bool)

	// Dimensiones personalizadas (opcionales)
	customTrackWidth  float32
	customTrackHeight float32

	// Colores personalizados (opcionales)
	isCustomColors bool
	customOnBg     color.Color
	customOffBg    color.Color
	customThumbCol color.Color

	// Estados internos de interacción
	isHovered bool
	isFocused bool

	// Hitbox delimitada del control para evitar clics accidentales en márgenes vacíos
	hitboxPos  fyne.Position
	hitboxSize fyne.Size

	// Objetos canvas del renderizado
	track       *canvas.Rectangle
	thumbHalo   *canvas.Circle // Halo reactivo en hover y foco (ripple sutil)
	thumbShadow *canvas.Circle // Sombra suave proyectada para relieve tridimensional
	thumb       *canvas.Circle // Botón circular deslizante
	label       *canvas.Text
}

var _ desktop.Hoverable = (*ModernSwitch)(nil)
var _ fyne.Disableable = (*ModernSwitch)(nil)
var _ fyne.Focusable = (*ModernSwitch)(nil)
var _ fyne.Tappable = (*ModernSwitch)(nil)
var _ fyne.Widget = (*ModernSwitch)(nil)

// NewModernSwitch crea un nuevo interruptor con los valores predeterminados del UIKit.
func NewModernSwitch() *ModernSwitch {
	sw := &ModernSwitch{
		Variant:       SwitchPrimary,
		LabelPosition: LabelTrailing,
		track:         canvas.NewRectangle(color.Transparent),
		thumbHalo:     canvas.NewCircle(color.Transparent),
		thumbShadow:   canvas.NewCircle(color.Transparent),
		thumb:         canvas.NewCircle(color.Transparent),
		label:         canvas.NewText("", color.Transparent),
	}

	sw.track.CornerRadius = defaultTrackHeight / 2
	sw.thumb.StrokeWidth = 0.5
	sw.thumb.StrokeColor = color.NRGBA{R: 0, G: 0, B: 0, A: 25}

	sw.thumbHalo.Hide()
	sw.label.TextSize = 13
	sw.label.Alignment = fyne.TextAlignLeading

	sw.ExtendBaseWidget(sw)
	return sw
}

// NewSwitch es el constructor flexible y polimórfico estilo xBase / Harbour de UIKit.
// Acepta argumentos en orden libre:
//
//	NewSwitch()
//	NewSwitch(true)
//	NewSwitch(func(on bool){ ... })
//	NewSwitch("Efecto Cristal", true, func(on bool){ ... })
//	NewSwitch("Activar modo", func(on bool){ ... })
//	NewSwitch(SwitchSuccess, "Notificaciones", true, func(on bool){ ... })
//	NewSwitch("Modo", LabelLeading, true, func(on bool){ ... })
func NewSwitch(args ...interface{}) *ModernSwitch {
	sw := NewModernSwitch()

	for _, arg := range args {
		switch val := arg.(type) {
		case bool:
			sw.On = val
		case string:
			sw.Text = val
		case SwitchVariant:
			sw.Variant = val
		case SwitchLabelPosition:
			sw.LabelPosition = val
		case func(bool):
			sw.OnChanged = val
		}
	}

	sw.Refresh()
	return sw
}

// NewPrimarySwitch crea un interruptor con color primario.
func NewPrimarySwitch(text string, on bool, onChanged func(bool)) *ModernSwitch {
	sw := NewSwitch(text, on, onChanged)
	sw.Variant = SwitchPrimary
	return sw
}

// NewSuccessSwitch crea un interruptor semántico en verde esmeralda.
func NewSuccessSwitch(text string, on bool, onChanged func(bool)) *ModernSwitch {
	sw := NewSwitch(text, on, onChanged)
	sw.Variant = SwitchSuccess
	return sw
}

// NewDangerSwitch crea un interruptor semántico en rojo de alerta.
func NewDangerSwitch(text string, on bool, onChanged func(bool)) *ModernSwitch {
	sw := NewSwitch(text, on, onChanged)
	sw.Variant = SwitchDanger
	return sw
}

// SetOn actualiza el estado del interruptor y dispara el callback si hubo cambio.
func (s *ModernSwitch) SetOn(on bool) *ModernSwitch {
	if s.On == on {
		return s
	}
	s.On = on
	if s.OnChanged != nil {
		s.OnChanged(on)
	}
	s.Refresh()
	return s
}

// SetChecked es un alias de conveniencia idéntico a widget.Check.
func (s *ModernSwitch) SetChecked(checked bool) *ModernSwitch {
	return s.SetOn(checked)
}

// Checked devuelve el estado activo del interruptor.
func (s *ModernSwitch) Checked() bool {
	return s.On
}

// IsOn devuelve el estado activo del interruptor.
func (s *ModernSwitch) IsOn() bool {
	return s.On
}

// Toggle invierte el estado actual del switch.
func (s *ModernSwitch) Toggle() *ModernSwitch {
	return s.SetOn(!s.On)
}

// SetText asigna la etiqueta textual del interruptor.
func (s *ModernSwitch) SetText(text string) *ModernSwitch {
	s.Text = text
	s.Refresh()
	return s
}

// GetText devuelve el texto de la etiqueta asociada.
func (s *ModernSwitch) GetText() string {
	return s.Text
}

// SetLabelPosition define la posición de la etiqueta respecto al switch.
func (s *ModernSwitch) SetLabelPosition(pos SwitchLabelPosition) *ModernSwitch {
	s.LabelPosition = pos
	s.Refresh()
	return s
}

// WithLeadingLabel coloca la etiqueta antes del switch: [Texto] [Switch].
func (s *ModernSwitch) WithLeadingLabel() *ModernSwitch {
	return s.SetLabelPosition(LabelLeading)
}

// WithTrailingLabel coloca la etiqueta después del switch: [Switch] [Texto].
func (s *ModernSwitch) WithTrailingLabel() *ModernSwitch {
	return s.SetLabelPosition(LabelTrailing)
}

// SetVariant actualiza la variante semántica de color.
func (s *ModernSwitch) SetVariant(v SwitchVariant) *ModernSwitch {
	s.Variant = v
	s.isCustomColors = false
	s.Refresh()
	return s
}

// SetOnChanged asigna el listener de cambio de estado.
func (s *ModernSwitch) SetOnChanged(fn func(bool)) *ModernSwitch {
	s.OnChanged = fn
	return s
}

// SetCustomColors permite personalizar libremente los colores de pista y botón.
func (s *ModernSwitch) SetCustomColors(onBg, offBg, thumbCol color.Color) *ModernSwitch {
	s.isCustomColors = true
	s.customOnBg = onBg
	s.customOffBg = offBg
	s.customThumbCol = thumbCol
	s.Refresh()
	return s
}

// SetDimensions permite sobreescribir el ancho y alto de la pista deslizante.
func (s *ModernSwitch) SetDimensions(width, height float32) *ModernSwitch {
	s.customTrackWidth = width
	s.customTrackHeight = height
	s.Refresh()
	return s
}

// --- Delimitación de Hitbox y Control de Puntero ---

func (s *ModernSwitch) hitboxContains(pos fyne.Position) bool {
	if s.hitboxSize.IsZero() {
		sz := s.Size()
		return pos.X >= 0 && pos.X <= sz.Width && pos.Y >= 0 && pos.Y <= sz.Height
	}
	return pos.X >= s.hitboxPos.X && pos.X <= (s.hitboxPos.X+s.hitboxSize.Width) &&
		pos.Y >= s.hitboxPos.Y && pos.Y <= (s.hitboxPos.Y+s.hitboxSize.Height)
}

func (s *ModernSwitch) Cursor() desktop.Cursor {
	if !s.Disabled() && s.isHovered {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

func (s *ModernSwitch) MouseIn(me *desktop.MouseEvent) {
	s.MouseMoved(me)
}

func (s *ModernSwitch) MouseMoved(me *desktop.MouseEvent) {
	if s.Disabled() {
		return
	}
	s.updateHover(&me.Position)
}

func (s *ModernSwitch) MouseOut() {
	s.updateHover(nil)
}

func (s *ModernSwitch) updateHover(pos *fyne.Position) {
	oldHover := s.isHovered
	s.isHovered = pos != nil && s.hitboxContains(*pos)
	if oldHover != s.isHovered {
		s.Refresh()
	}
}

func (s *ModernSwitch) Tapped(pe *fyne.PointEvent) {
	if s.Disabled() {
		return
	}
	if !s.hitboxContains(pe.Position) {
		return
	}
	s.Toggle()
}

func (s *ModernSwitch) TappedSecondary(_ *fyne.PointEvent) {}

func (s *ModernSwitch) FocusGained() {
	if s.Disabled() {
		return
	}
	s.isFocused = true
	s.Refresh()
}

func (s *ModernSwitch) FocusLost() {
	s.isFocused = false
	s.Refresh()
}

func (s *ModernSwitch) TypedRune(r rune) {
	if s.Disabled() {
		return
	}
	if r == ' ' || r == '\n' {
		s.Toggle()
	}
}

func (s *ModernSwitch) TypedKey(key *fyne.KeyEvent) {
	if s.Disabled() {
		return
	}
	if key.Name == fyne.KeySpace || key.Name == fyne.KeyReturn {
		s.Toggle()
	}
}

// --- Lógica de Color del Switch ---

func (s *ModernSwitch) getActiveTrackColor() color.Color {
	if s.isCustomColors && s.customOnBg != nil {
		return s.customOnBg
	}
	pal := GetCurrentThemePalette()
	if pal != nil {
		switch s.Variant {
		case SwitchSuccess:
			return pal.Success
		case SwitchDanger:
			return pal.Danger
		case SwitchWarning:
			return pal.Warning
		case SwitchPrimary:
			fallthrough
		default:
			return pal.Accent
		}
	}

	th := fyne.CurrentApp().Settings().Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	switch s.Variant {
	case SwitchSuccess:
		return th.Color(theme.ColorNameSuccess, v)
	case SwitchDanger:
		return th.Color(theme.ColorNameError, v)
	case SwitchWarning:
		return th.Color(theme.ColorNameWarning, v)
	default:
		return th.Color(theme.ColorNamePrimary, v)
	}
}

func (s *ModernSwitch) getInactiveTrackColor() color.Color {
	if s.isCustomColors && s.customOffBg != nil {
		return s.customOffBg
	}
	if IsDarkTheme() {
		return color.NRGBA{R: 51, G: 65, B: 85, A: 255} // Slate 700
	}
	return color.NRGBA{R: 203, G: 213, B: 225, A: 255} // Slate 300
}

func (s *ModernSwitch) getThumbColor() color.Color {
	if s.isCustomColors && s.customThumbCol != nil {
		return s.customThumbCol
	}
	return color.NRGBA{R: 255, G: 255, B: 255, A: 255} // Blanco puro
}

func (s *ModernSwitch) getLabelColor() color.Color {
	if s.Disabled() {
		return ThemeColor(theme.ColorNameDisabled)
	}
	return ThemeColor(ColorNameTextPrimary)
}

func interpolateColor(c1, c2 color.Color, t float32) color.Color {
	return BlendColor(c1, c2, t)
}

func applyAlpha(c color.Color, factor float32) color.Color {
	return ColorWithAlpha(c, factor)
}

// --- Renderer del Widget ---

type modernSwitchRenderer struct {
	sw       *ModernSwitch
	progress float32 // 0.0 = OFF, 1.0 = ON
	anim     *fyne.Animation
	lastOn   bool
}

func (s *ModernSwitch) CreateRenderer() fyne.WidgetRenderer {
	r := &modernSwitchRenderer{
		sw:     s,
		lastOn: s.On,
	}
	if s.On {
		r.progress = 1.0
	} else {
		r.progress = 0.0
	}
	r.updateColors()
	return r
}

func (r *modernSwitchRenderer) Objects() []fyne.CanvasObject {
	objs := []fyne.CanvasObject{
		r.sw.track,
		r.sw.thumbHalo,
		r.sw.thumbShadow,
		r.sw.thumb,
	}
	if strings.TrimSpace(r.sw.Text) != "" {
		objs = append(objs, r.sw.label)
	}
	return objs
}

func (r *modernSwitchRenderer) Destroy() {
	if r.anim != nil {
		r.anim.Stop()
	}
}

func (r *modernSwitchRenderer) getTrackDimensions() (float32, float32) {
	w := defaultTrackWidth
	h := defaultTrackHeight
	if r.sw.customTrackWidth > 0 {
		w = r.sw.customTrackWidth
	}
	if r.sw.customTrackHeight > 0 {
		h = r.sw.customTrackHeight
	}
	return w, h
}

func (r *modernSwitchRenderer) MinSize() fyne.Size {
	trackW, trackH := r.getTrackDimensions()
	width := trackW
	height := trackH

	if strings.TrimSpace(r.sw.Text) != "" {
		labelSize := fyne.MeasureText(r.sw.Text, r.sw.label.TextSize, r.sw.label.TextStyle)
		width += defaultLabelGap + labelSize.Width
		if labelSize.Height > height {
			height = labelSize.Height
		}
	}

	return fyne.NewSize(width, height)
}

func (r *modernSwitchRenderer) Layout(size fyne.Size) {
	trackW, trackH := r.getTrackDimensions()
	hasLabel := strings.TrimSpace(r.sw.Text) != ""

	var labelSize fyne.Size
	if hasLabel {
		labelSize = fyne.MeasureText(r.sw.Text, r.sw.label.TextSize, r.sw.label.TextStyle)
	}

	trackY := (size.Height - trackH) / 2

	if hasLabel && r.sw.LabelPosition == LabelLeading {
		// [Texto] [Switch]
		labelX := float32(0)
		labelY := (size.Height - labelSize.Height) / 2
		r.sw.label.Move(fyne.NewPos(labelX, labelY))
		r.sw.label.Resize(labelSize)

		trackX := labelSize.Width + defaultLabelGap
		r.sw.track.Move(fyne.NewPos(trackX, trackY))
		r.sw.track.Resize(fyne.NewSize(trackW, trackH))

		maxH := trackH
		if labelSize.Height > maxH {
			maxH = labelSize.Height
		}
		r.sw.hitboxPos = fyne.NewPos(0, (size.Height-maxH)/2)
		r.sw.hitboxSize = fyne.NewSize(trackX+trackW, maxH)
	} else {
		// [Switch] [Texto] (Default)
		trackX := float32(0)
		r.sw.track.Move(fyne.NewPos(trackX, trackY))
		r.sw.track.Resize(fyne.NewSize(trackW, trackH))

		totalW := trackW
		if hasLabel {
			labelX := trackW + defaultLabelGap
			labelY := (size.Height - labelSize.Height) / 2
			r.sw.label.Move(fyne.NewPos(labelX, labelY))
			r.sw.label.Resize(labelSize)
			totalW += defaultLabelGap + labelSize.Width
		}

		maxH := trackH
		if hasLabel && labelSize.Height > maxH {
			maxH = labelSize.Height
		}
		r.sw.hitboxPos = fyne.NewPos(0, (size.Height-maxH)/2)
		r.sw.hitboxSize = fyne.NewSize(totalW, maxH)
	}

	r.sw.track.CornerRadius = trackH / 2
	r.updateThumbPosition()
}

func (r *modernSwitchRenderer) Refresh() {
	r.sw.label.Text = r.sw.Text
	r.sw.label.Color = r.sw.getLabelColor()

	if r.sw.On != r.lastOn {
		r.lastOn = r.sw.On
		r.animateThumb()
	} else {
		r.updateThumbPosition()
		r.updateColors()
	}

	r.sw.track.Refresh()
	r.sw.thumbHalo.Refresh()
	r.sw.thumbShadow.Refresh()
	r.sw.thumb.Refresh()
	r.sw.label.Refresh()
	r.Layout(r.sw.Size())
}

func (r *modernSwitchRenderer) updateThumbPosition() {
	trackW, trackH := r.getTrackDimensions()
	pad := defaultThumbPad
	thumbSize := trackH - (pad * 2)

	size := r.sw.Size()
	trackY := (size.Height - trackH) / 2

	var trackX float32
	if strings.TrimSpace(r.sw.Text) != "" && r.sw.LabelPosition == LabelLeading {
		labelSize := fyne.MeasureText(r.sw.Text, r.sw.label.TextSize, r.sw.label.TextStyle)
		trackX = labelSize.Width + defaultLabelGap
	}

	travelDistance := trackW - thumbSize - (pad * 2)
	thumbX := trackX + pad + (r.progress * travelDistance)
	thumbY := trackY + pad

	// Botón circular principal
	r.sw.thumb.Position1 = fyne.NewPos(thumbX, thumbY)
	r.sw.thumb.Position2 = fyne.NewPos(thumbX+thumbSize, thumbY+thumbSize)

	// Sombra suave proyectada hacia abajo (1px de relieve)
	const shadowOffsetY = float32(1.0)
	r.sw.thumbShadow.Position1 = fyne.NewPos(thumbX, thumbY+shadowOffsetY)
	r.sw.thumbShadow.Position2 = fyne.NewPos(thumbX+thumbSize, thumbY+thumbSize+shadowOffsetY)

	// Halo circular de hover y foco (diámetro = thumbSize + 8px)
	const haloPadding = float32(4.0)
	r.sw.thumbHalo.Position1 = fyne.NewPos(thumbX-haloPadding, thumbY-haloPadding)
	r.sw.thumbHalo.Position2 = fyne.NewPos(thumbX+thumbSize+haloPadding, thumbY+thumbSize+haloPadding)
}

func (r *modernSwitchRenderer) animateThumb() {
	if r.anim != nil {
		r.anim.Stop()
	}

	target := float32(0.0)
	if r.sw.On {
		target = 1.0
	}

	app := fyne.CurrentApp()
	if app != nil && app.Settings() != nil && !app.Settings().ShowAnimations() {
		r.progress = target
		r.updateThumbPosition()
		r.updateColors()
		r.sw.track.Refresh()
		r.sw.thumbHalo.Refresh()
		r.sw.thumbShadow.Refresh()
		r.sw.thumb.Refresh()
		return
	}

	from := r.progress
	delta := target - from
	r.anim = fyne.NewAnimation(canvas.DurationShort, func(done float32) {
		r.progress = from + delta*done
		r.updateThumbPosition()
		r.updateColors()
		canvas.Refresh(r.sw.track)
		canvas.Refresh(r.sw.thumbHalo)
		canvas.Refresh(r.sw.thumbShadow)
		canvas.Refresh(r.sw.thumb)
	})
	r.anim.Curve = fyne.AnimationEaseInOut
	r.anim.Start()
}

func (r *modernSwitchRenderer) updateColors() {
	onColor := r.sw.getActiveTrackColor()
	offColor := r.sw.getInactiveTrackColor()

	// Mezcla suave de color durante la transición
	trackColor := interpolateColor(offColor, onColor, r.progress)

	// Halo y borde reactivo de Foco / Hover
	if r.sw.isFocused {
		r.sw.thumbHalo.FillColor = applyAlpha(ThemeColor(theme.ColorNameFocus), 0.28)
		r.sw.thumbHalo.Show()
		r.sw.track.StrokeColor = ThemeColor(theme.ColorNameFocus)
		r.sw.track.StrokeWidth = 1.5
	} else if r.sw.isHovered && !r.sw.Disabled() {
		r.sw.thumbHalo.FillColor = applyAlpha(onColor, 0.20)
		r.sw.thumbHalo.Show()
		r.sw.track.StrokeColor = applyAlpha(onColor, 0.35)
		r.sw.track.StrokeWidth = 1
	} else {
		r.sw.thumbHalo.Hide()
		r.sw.track.StrokeColor = color.Transparent
		r.sw.track.StrokeWidth = 0
	}

	thumbColor := r.sw.getThumbColor()

	// Sombra física del thumb
	if IsDarkTheme() {
		r.sw.thumbShadow.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 90}
	} else {
		r.sw.thumbShadow.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 35}
	}

	if r.sw.Disabled() {
		trackColor = applyAlpha(trackColor, 0.45)
		thumbColor = applyAlpha(thumbColor, 0.6)
		r.sw.thumbShadow.FillColor = color.Transparent
		r.sw.thumbHalo.Hide()
	}

	r.sw.track.FillColor = trackColor
	r.sw.thumb.FillColor = thumbColor
}
