package uikit

import (
	"image/color"
	"math"
	"strconv"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Widget = (*ModernSlider)(nil)
var _ fyne.Disableable = (*ModernSlider)(nil)

// ModernSlider es un slider interactivo refinado con soporte de pastilla (badge) de valor en tiempo real,
// formateo de unidades (prefijo/sufijo/precisión), anti-jitter, título superior opcional y reactividad fluida.
type ModernSlider struct {
	widget.DisableableWidget

	mu            sync.RWMutex
	min           float64
	max           float64
	step          float64
	value         float64
	title         string
	prefix        string
	suffix        string
	decimals      int // -1 para recorte dinámico de ceros, o número exacto de decimales (ej. 0, 1, 2)
	formatFunc    func(float64) string
	showValue     bool
	showBadge     bool
	onChanged     func(float64)
	onChangeEnded func(float64)
}

// NewModernSlider crea una nueva instancia de ModernSlider con rango mínimo y máximo.
func NewModernSlider(min, max float64) *ModernSlider {
	if max < min {
		min, max = max, min
	}
	s := &ModernSlider{
		min:       min,
		max:       max,
		step:      1,
		value:     min,
		decimals:  -1,
		showValue: true,
		showBadge: true,
	}
	s.ExtendBaseWidget(s)
	return s
}

// NewSlider es un alias conciso de NewModernSlider.
func NewSlider(min, max float64) *ModernSlider {
	return NewModernSlider(min, max)
}

// NewAgeSlider crea un slider especializado en edades de jubilación (por defecto 60 a 70 años, paso 1).
func NewAgeSlider(min, max float64) *ModernSlider {
	if min == 0 && max == 0 {
		min, max = 60, 70
	}
	return NewModernSlider(min, max).
		WithStep(1).
		WithPrecision(0).
		WithSuffix(" años")
}

// NewPercentageSlider crea un slider especializado en porcentajes (ej. jubilación parcial/relevo 0 a 100%).
func NewPercentageSlider(min, max float64) *ModernSlider {
	if min == 0 && max == 0 {
		min, max = 0, 100
	}
	return NewModernSlider(min, max).
		WithStep(1).
		WithPrecision(0).
		WithSuffix(" %")
}

// NewCurrencySlider crea un slider monetario con sufijo de euros y 2 decimales.
func NewCurrencySlider(min, max, step float64) *ModernSlider {
	if step <= 0 {
		step = 50
	}
	return NewModernSlider(min, max).
		WithStep(step).
		WithPrecision(2).
		WithSuffix(" €")
}

// WithStep ajusta el incremento del deslizador.
func (s *ModernSlider) WithStep(step float64) *ModernSlider {
	s.mu.Lock()
	if step > 0 {
		s.step = step
	}
	s.mu.Unlock()
	s.Refresh()
	return s
}

// SetStep ajusta el incremento del deslizador.
func (s *ModernSlider) SetStep(step float64) {
	s.WithStep(step)
}

// Value obtiene el valor actual del slider de manera thread-safe.
func (s *ModernSlider) Value() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

// SetValue define el valor actual, limitándolo al rango [min, max] y alineándolo al step.
func (s *ModernSlider) SetValue(v float64) {
	s.mu.Lock()
	if v < s.min {
		v = s.min
	} else if v > s.max {
		v = s.max
	}
	if s.step > 0 {
		steps := math.Round((v - s.min) / s.step)
		v = s.min + (steps * s.step)
		if v > s.max {
			v = s.max
		}
	}
	s.value = v
	s.mu.Unlock()
	s.Refresh()
}

// WithValue inicializa fluidamente el valor del deslizador.
func (s *ModernSlider) WithValue(v float64) *ModernSlider {
	s.SetValue(v)
	return s
}

// WithTitle añade una etiqueta de título superior (estilo tarjeta de ajustes de macOS/iOS).
func (s *ModernSlider) WithTitle(title string) *ModernSlider {
	s.mu.Lock()
	s.title = title
	s.mu.Unlock()
	s.Refresh()
	return s
}

// SetTitle define el título superior.
func (s *ModernSlider) SetTitle(title string) {
	s.WithTitle(title)
}

// WithSuffix añade un sufijo de texto al valor numérico (ej. " %", " años", " €", " meses").
func (s *ModernSlider) WithSuffix(suffix string) *ModernSlider {
	s.mu.Lock()
	s.suffix = suffix
	s.mu.Unlock()
	s.Refresh()
	return s
}

// WithPrefix añade un prefijo al valor numérico (ej. "€ ", "$ ").
func (s *ModernSlider) WithPrefix(prefix string) *ModernSlider {
	s.mu.Lock()
	s.prefix = prefix
	s.mu.Unlock()
	s.Refresh()
	return s
}

// WithPrecision fija el número de decimales a mostrar (0 para enteros, 1, 2, etc.).
func (s *ModernSlider) WithPrecision(decimals int) *ModernSlider {
	s.mu.Lock()
	s.decimals = decimals
	s.mu.Unlock()
	s.Refresh()
	return s
}

// WithFormat asigna una función personalizada para representar el valor.
func (s *ModernSlider) WithFormat(fn func(float64) string) *ModernSlider {
	s.mu.Lock()
	s.formatFunc = fn
	s.mu.Unlock()
	s.Refresh()
	return s
}

// WithBadge activa o desactiva la pastilla/badge alrededor del valor.
func (s *ModernSlider) WithBadge(show bool) *ModernSlider {
	s.mu.Lock()
	s.showBadge = show
	s.mu.Unlock()
	s.Refresh()
	return s
}

// WithValueVisible controla si se muestra el valor actual numérico.
func (s *ModernSlider) WithValueVisible(show bool) *ModernSlider {
	s.mu.Lock()
	s.showValue = show
	s.mu.Unlock()
	s.Refresh()
	return s
}

// OnChanged registra un callback invocado continuamente mientras se arrastra el slider.
func (s *ModernSlider) OnChanged(fn func(float64)) *ModernSlider {
	s.mu.Lock()
	s.onChanged = fn
	s.mu.Unlock()
	return s
}

// OnChangeEnded registra un callback invocado al soltar el deslizador.
func (s *ModernSlider) OnChangeEnded(fn func(float64)) *ModernSlider {
	s.mu.Lock()
	s.onChangeEnded = fn
	s.mu.Unlock()
	return s
}

// formatValue genera la representación textual formateada de un número según la configuración del slider.
func (s *ModernSlider) formatValue(v float64) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.formatFunc != nil {
		return s.formatFunc(v)
	}

	var numStr string
	if s.decimals >= 0 {
		numStr = strconv.FormatFloat(v, 'f', s.decimals, 64)
	} else {
		numStr = strconv.FormatFloat(v, 'f', 4, 64)
		numStr = strings.TrimRight(numStr, "0")
		numStr = strings.TrimSuffix(numStr, ".")
		if numStr == "" || numStr == "-0" {
			numStr = "0"
		}
	}

	return s.prefix + numStr + s.suffix
}

// CreateRenderer construye el árbol visual y los enlaces reactivos de ModernSlider.
func (s *ModernSlider) CreateRenderer() fyne.WidgetRenderer {
	s.mu.RLock()
	minVal, maxVal, stepVal, currentVal := s.min, s.max, s.step, s.value
	s.mu.RUnlock()

	// 1. Título superior opcional
	titleText := canvas.NewText("", ThemeColor(ColorNameTextMuted))
	titleText.TextSize = 12.5
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	// 2. Control Slider de Fyne
	innerSlider := widget.NewSlider(minVal, maxVal)
	innerSlider.Step = stepVal
	innerSlider.Value = currentVal

	// 3. Pastilla (Badge) de valor
	badgeBg := canvas.NewRectangle(ThemeColor(ColorNameBtnPrimaryBg))
	badgeBg.CornerRadius = 5
	badgeBg.StrokeWidth = 1
	badgeBg.StrokeColor = ThemeColor(ColorNameBtnPrimaryBorder)

	valueText := canvas.NewText("", ThemeColor(ColorNameBtnPrimaryFg))
	valueText.TextSize = 12.5
	valueText.TextStyle = fyne.TextStyle{Bold: true}
	valueText.Alignment = fyne.TextAlignCenter

	r := &modernSliderRenderer{
		widget:      s,
		titleText:   titleText,
		innerSlider: innerSlider,
		badgeBg:     badgeBg,
		valueText:   valueText,
	}

	innerSlider.OnChanged = func(v float64) {
		s.mu.Lock()
		s.value = v
		callback := s.onChanged
		s.mu.Unlock()

		r.updateValueLabel(v)
		if callback != nil {
			callback(v)
		}
	}

	innerSlider.OnChangeEnded = func(v float64) {
		s.mu.Lock()
		s.value = v
		callback := s.onChangeEnded
		s.mu.Unlock()

		r.updateValueLabel(v)
		if callback != nil {
			callback(v)
		}
	}

	r.applyStyles()
	r.calculateBadgeWidth()
	r.updateValueLabel(currentVal)

	return r
}

type modernSliderRenderer struct {
	widget      *ModernSlider
	titleText   *canvas.Text
	innerSlider *widget.Slider
	badgeBg     *canvas.Rectangle
	valueText   *canvas.Text

	cachedBadgeW float32
}

func (r *modernSliderRenderer) Objects() []fyne.CanvasObject {
	r.widget.mu.RLock()
	hasTitle := r.widget.title != ""
	showVal := r.widget.showValue
	showBdg := r.widget.showBadge
	r.widget.mu.RUnlock()

	objs := make([]fyne.CanvasObject, 0, 4)
	if hasTitle {
		objs = append(objs, r.titleText)
	}
	objs = append(objs, r.innerSlider)
	if showVal {
		if showBdg {
			objs = append(objs, r.badgeBg)
		}
		objs = append(objs, r.valueText)
	}
	return objs
}

func (r *modernSliderRenderer) Destroy() {}

// calculateBadgeWidth evalúa el valor más ancho en los extremos y pasos clave para evitar
// oscilaciones y jitter horizontal mientras se arrastra el slider.
func (r *modernSliderRenderer) calculateBadgeWidth() {
	r.widget.mu.RLock()
	showVal := r.widget.showValue
	showBdg := r.widget.showBadge
	minVal, maxVal, stepVal := r.widget.min, r.widget.max, r.widget.step
	r.widget.mu.RUnlock()

	if !showVal {
		r.cachedBadgeW = 0
		return
	}

	candidates := []float64{
		minVal,
		maxVal,
		minVal + stepVal,
		maxVal - stepVal,
		0,
		(minVal + maxVal) / 2,
	}

	dummyText := canvas.NewText("", ThemeColor(ColorNameTextPrimary))
	dummyText.TextSize = 12.5
	dummyText.TextStyle = fyne.TextStyle{Bold: true}

	var maxTextW float32 = 0
	for _, val := range candidates {
		dummyText.Text = r.widget.formatValue(val)
		w := dummyText.MinSize().Width
		if w > maxTextW {
			maxTextW = w
		}
	}

	if showBdg {
		r.cachedBadgeW = fyne.Max(maxTextW+18, 48) // Margen de 9px a cada lado, min 48px
	} else {
		r.cachedBadgeW = maxTextW + 4
	}
}

func (r *modernSliderRenderer) updateValueLabel(v float64) {
	txt := r.widget.formatValue(v)
	r.valueText.Text = txt
	r.valueText.Refresh()
}

func (r *modernSliderRenderer) MinSize() fyne.Size {
	r.widget.mu.RLock()
	hasTitle := r.widget.title != ""
	showVal := r.widget.showValue
	r.widget.mu.RUnlock()

	sliderMin := r.innerSlider.MinSize()
	var totalH float32 = sliderMin.Height
	var totalW float32 = sliderMin.Width

	if showVal {
		totalW += r.cachedBadgeW + theme.Padding()
	}

	if hasTitle {
		titleH := r.titleText.MinSize().Height + 4
		totalH += titleH
		titleW := r.titleText.MinSize().Width
		if showVal {
			titleW += r.cachedBadgeW + theme.Padding()
		}
		if titleW > totalW {
			totalW = titleW
		}
	}

	return fyne.NewSize(fyne.Max(totalW, 160), totalH)
}

func (r *modernSliderRenderer) Layout(size fyne.Size) {
	r.widget.mu.RLock()
	hasTitle := r.widget.title != ""
	showVal := r.widget.showValue
	showBdg := r.widget.showBadge
	r.widget.mu.RUnlock()

	pad := theme.Padding()
	badgeW := r.cachedBadgeW
	badgeH := float32(24)

	if hasTitle {
		// Disposición tipo Tarjeta de Ajuste (Settings Row):
		// Fila 1: Título (izq) y Valor (der)
		// Fila 2: Slider abarcando todo el ancho
		titleMin := r.titleText.MinSize()
		row1H := fyne.Max(titleMin.Height, badgeH)

		r.titleText.Resize(fyne.NewSize(size.Width-badgeW-pad, titleMin.Height))
		r.titleText.Move(fyne.NewPos(0, (row1H-titleMin.Height)/2))

		if showVal {
			badgeX := size.Width - badgeW
			badgeY := (row1H - badgeH) / 2
			if showBdg {
				r.badgeBg.Resize(fyne.NewSize(badgeW, badgeH))
				r.badgeBg.Move(fyne.NewPos(badgeX, badgeY))
				r.badgeBg.Show()
			} else {
				r.badgeBg.Hide()
			}

			valMin := r.valueText.MinSize()
			r.valueText.Resize(fyne.NewSize(badgeW, valMin.Height))
			r.valueText.Move(fyne.NewPos(badgeX, badgeY+(badgeH-valMin.Height)/2))
			r.valueText.Show()
		}

		sliderY := row1H + 4
		sliderH := size.Height - sliderY
		if sliderH < r.innerSlider.MinSize().Height {
			sliderH = r.innerSlider.MinSize().Height
		}
		r.innerSlider.Resize(fyne.NewSize(size.Width, sliderH))
		r.innerSlider.Move(fyne.NewPos(0, sliderY))

	} else {
		// Disposición Horizontal Inline:
		// [ ===== Slider ===== ] [ Badge ]
		sliderW := size.Width
		if showVal {
			sliderW = size.Width - badgeW - pad
			if sliderW < 60 {
				sliderW = 60
			}

			badgeX := size.Width - badgeW
			badgeY := (size.Height - badgeH) / 2
			if showBdg {
				r.badgeBg.Resize(fyne.NewSize(badgeW, badgeH))
				r.badgeBg.Move(fyne.NewPos(badgeX, badgeY))
				r.badgeBg.Show()
			} else {
				r.badgeBg.Hide()
			}

			valMin := r.valueText.MinSize()
			r.valueText.Resize(fyne.NewSize(badgeW, valMin.Height))
			r.valueText.Move(fyne.NewPos(badgeX, badgeY+(badgeH-valMin.Height)/2))
			r.valueText.Show()
		}

		sliderMinH := r.innerSlider.MinSize().Height
		r.innerSlider.Resize(fyne.NewSize(sliderW, sliderMinH))
		r.innerSlider.Move(fyne.NewPos(0, (size.Height-sliderMinH)/2))
	}
}

func (r *modernSliderRenderer) Refresh() {
	r.widget.mu.RLock()
	val := r.widget.value
	stepVal := r.widget.step
	titleStr := r.widget.title
	isDisabled := r.widget.Disabled()
	r.widget.mu.RUnlock()

	r.innerSlider.Step = stepVal
	if r.innerSlider.Value != val {
		r.innerSlider.Value = val
		r.innerSlider.Refresh()
	}

	if isDisabled {
		r.innerSlider.Disable()
	} else {
		r.innerSlider.Enable()
	}

	r.titleText.Text = titleStr
	r.applyStyles()
	r.calculateBadgeWidth()
	r.updateValueLabel(val)

	canvas.Refresh(r.widget)
}

func (r *modernSliderRenderer) applyStyles() {
	isDisabled := r.widget.Disabled()

	if isDisabled {
		r.titleText.Color = ThemeColor(ColorNameTextMuted)
		r.badgeBg.FillColor = ThemeColor(ColorNameBtnDefaultBg)
		r.badgeBg.StrokeColor = color.Transparent
		r.valueText.Color = ThemeColor(ColorNameBtnMuted)
	} else {
		r.titleText.Color = ThemeColor(ColorNameTextPrimary)
		r.badgeBg.FillColor = ThemeColor(ColorNameBtnPrimaryBg)
		r.badgeBg.StrokeColor = ThemeColor(ColorNameBtnPrimaryBorder)
		r.valueText.Color = ThemeColor(ColorNameBtnPrimaryFg)
	}

	r.titleText.Refresh()
	r.badgeBg.Refresh()
	r.valueText.Refresh()
}
