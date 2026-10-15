package uikit

import (
	"image/color"
	"math"
	"strconv"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Widget = (*ModernSlider)(nil)
var _ fyne.Disableable = (*ModernSlider)(nil)
var _ fyne.Scrollable = (*ModernSlider)(nil)

// ModernSlider es un slider interactivo contemporáneo y de alta precisión visual,
// con pastilla en tiempo real, micro-botones (+/-), marcas Min/Max, soporte de rueda de ratón,
// edición manual al hacer clic en el valor, anti-jitter y contraste de pista impecable.
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
	color         color.Color
	trackColor    color.Color
	decimals      int // -1 para recorte dinámico de ceros, o número exacto de decimales (ej. 0, 1, 2)
	formatFunc    func(float64) string
	showValue     bool
	showBadge     bool
	showSteppers  bool
	showMinMax    bool
	minLabel      string
	maxLabel      string
	editableBadge bool
	scrollDelta   float32
	onChanged     func(float64)
	onChangeEnded func(float64)
}

// NewModernSlider crea una nueva instancia de ModernSlider con rango mínimo y máximo.
func NewModernSlider(min, max float64) *ModernSlider {
	if max < min {
		min, max = max, min
	}
	s := &ModernSlider{
		min:           min,
		max:           max,
		step:          1,
		value:         min,
		decimals:      -1,
		showValue:     true,
		showBadge:     true,
		editableBadge: true,
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
		WithSuffix(" años").
		WithSteppers(true).
		WithMinMax(true)
}

// NewPercentageSlider crea un slider especializado en porcentajes (ej. jubilación parcial/relevo 0 a 100%).
func NewPercentageSlider(min, max float64) *ModernSlider {
	if min == 0 && max == 0 {
		min, max = 0, 100
	}
	return NewModernSlider(min, max).
		WithStep(1).
		WithPrecision(0).
		WithSuffix(" %").
		WithSteppers(true).
		WithMinMax(true)
}

// NewCurrencySlider crea un slider monetario con sufijo de euros y 2 decimales.
func NewCurrencySlider(min, max, step float64) *ModernSlider {
	if step <= 0 {
		step = 50
	}
	return NewModernSlider(min, max).
		WithStep(step).
		WithPrecision(2).
		WithSuffix(" €").
		WithSteppers(true).
		WithMinMax(true)
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

// StepUp incrementa el valor en un paso (step) y dispara los callbacks correspondientes.
func (s *ModernSlider) StepUp() {
	s.mu.RLock()
	cur := s.value
	step := s.step
	s.mu.RUnlock()

	s.SetValue(cur + step)

	s.mu.RLock()
	cb := s.onChanged
	cbEnd := s.onChangeEnded
	v := s.value
	s.mu.RUnlock()

	if cb != nil {
		cb(v)
	}
	if cbEnd != nil {
		cbEnd(v)
	}
}

// StepDown decrementa el valor en un paso (step) y dispara los callbacks correspondientes.
func (s *ModernSlider) StepDown() {
	s.mu.RLock()
	cur := s.value
	step := s.step
	s.mu.RUnlock()

	s.SetValue(cur - step)

	s.mu.RLock()
	cb := s.onChanged
	cbEnd := s.onChangeEnded
	v := s.value
	s.mu.RUnlock()

	if cb != nil {
		cb(v)
	}
	if cbEnd != nil {
		cbEnd(v)
	}
}

// Scrolled permite ajustar el deslizador suavemente mediante la rueda del ratón o trackpad.
func (s *ModernSlider) Scrolled(e *fyne.ScrollEvent) {
	if s.Disabled() {
		return
	}

	dx := e.Scrolled.DX
	dy := e.Scrolled.DY

	// Seleccionar el delta del eje predominante para evitar que la micro-deriva
	// horizontal cancele o invierta el sentido de la rueda vertical
	var delta float32
	if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
		delta = dx
	} else {
		delta = dy
	}

	if delta == 0 {
		return
	}

	s.mu.Lock()
	s.scrollDelta += delta
	// En GLFW, 1 muesca de rueda física envía |delta| >= 8 (10 en macOS, 25 en Linux/Windows).
	// Cualquier evento con |delta| >= 8 representa 1 muesca de rueda -> 1 paso exacto.
	// Para trackpads con micro-desplazamientos continuos (|delta| < 8), acumulamos hasta 8.
	stepsToTake := 0
	if s.scrollDelta >= 8.0 {
		stepsToTake = 1
		s.scrollDelta = 0
	} else if s.scrollDelta <= -8.0 {
		stepsToTake = -1
		s.scrollDelta = 0
	}
	s.mu.Unlock()

	if stepsToTake > 0 {
		for i := 0; i < stepsToTake; i++ {
			s.StepUp()
		}
	} else if stepsToTake < 0 {
		for i := 0; i < -stepsToTake; i++ {
			s.StepDown()
		}
	}
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

// WithSteppers activa o desactiva los micro-botones de paso fino [-] y [+].
func (s *ModernSlider) WithSteppers(show bool) *ModernSlider {
	s.mu.Lock()
	s.showSteppers = show
	s.mu.Unlock()
	s.Refresh()
	return s
}

// WithMinMax activa la visualización de etiquetas en los extremos de rango inferior y superior.
func (s *ModernSlider) WithMinMax(show bool) *ModernSlider {
	s.mu.Lock()
	s.showMinMax = show
	s.mu.Unlock()
	s.Refresh()
	return s
}

// WithRangeLabels permite personalizar las etiquetas de los extremos de rango.
func (s *ModernSlider) WithRangeLabels(minStr, maxStr string) *ModernSlider {
	s.mu.Lock()
	s.showMinMax = true
	s.minLabel = minStr
	s.maxLabel = maxStr
	s.mu.Unlock()
	s.Refresh()
	return s
}

// WithEditableBadge permite que el usuario haga clic en la pastilla para escribir el valor numérico.
func (s *ModernSlider) WithEditableBadge(editable bool) *ModernSlider {
	s.mu.Lock()
	s.editableBadge = editable
	s.mu.Unlock()
	return s
}

// WithColor define el color del tirador (thumb) y del segmento activo del slider.
func (s *ModernSlider) WithColor(c color.Color) *ModernSlider {
	s.mu.Lock()
	s.color = c
	s.mu.Unlock()
	s.Refresh()
	return s
}

// WithHexColor define el color principal mediante código hexadecimal (ej. "#2563eb").
func (s *ModernSlider) WithHexColor(hex string) *ModernSlider {
	return s.WithColor(HexColor(hex))
}

// WithTrackColor define manualmente el color de la pista inactiva del slider.
func (s *ModernSlider) WithTrackColor(c color.Color) *ModernSlider {
	s.mu.Lock()
	s.trackColor = c
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

	// 2. Control Slider de Fyne con tema adaptativo (Track nítido y Thumb estilizado)
	innerSlider := newScrollableSlider(minVal, maxVal, s)
	innerSlider.Step = stepVal
	innerSlider.Value = currentVal

	sliderTheme := &modernSliderTheme{slider: s}
	themedSlider := container.NewThemeOverride(innerSlider, sliderTheme)

	// 3. Pastilla (Badge) de valor y micro-editor de entrada
	badgeBg := canvas.NewRectangle(ThemeColor(ColorNameBtnPrimaryBg))
	badgeBg.CornerRadius = 5
	badgeBg.StrokeWidth = 1
	badgeBg.StrokeColor = ThemeColor(ColorNameBtnPrimaryBorder)

	valueText := canvas.NewText("", ThemeColor(ColorNameBtnPrimaryFg))
	valueText.TextSize = 12.5
	valueText.TextStyle = fyne.TextStyle{Bold: true}
	valueText.Alignment = fyne.TextAlignCenter

	valueEntry := widget.NewEntry()
	valueEntry.TextStyle = fyne.TextStyle{Bold: true}
	valueEntry.Hide()

	// 4. Botones de Micro-Paso [-] y [+]
	stepperMinus := newSliderStepper("-", s, func() {
		s.StepDown()
	})
	stepperPlus := newSliderStepper("+", s, func() {
		s.StepUp()
	})

	// 5. Etiquetas de rango Min / Max
	minText := canvas.NewText("", ThemeColor(ColorNameTextMuted))
	minText.TextSize = 10.5
	minText.Alignment = fyne.TextAlignLeading

	maxText := canvas.NewText("", ThemeColor(ColorNameTextMuted))
	maxText.TextSize = 10.5
	maxText.Alignment = fyne.TextAlignTrailing

	r := &modernSliderRenderer{
		widget:       s,
		titleText:    titleText,
		innerSlider:  innerSlider,
		themedSlider: themedSlider,
		badgeBg:      badgeBg,
		valueText:    valueText,
		valueEntry:   valueEntry,
		stepperMinus: stepperMinus,
		stepperPlus:  stepperPlus,
		minText:      minText,
		maxText:      maxText,
	}

	badgeTap := newBadgeTappable(s, func() {
		r.startEditing()
	})
	r.badgeTap = badgeTap

	valueEntry.OnSubmitted = func(text string) {
		val, err := strconv.ParseFloat(strings.TrimSpace(text), 64)
		if err == nil {
			s.SetValue(val)
			s.mu.RLock()
			cb := s.onChanged
			cbEnd := s.onChangeEnded
			v := s.value
			s.mu.RUnlock()
			if cb != nil {
				cb(v)
			}
			if cbEnd != nil {
				cbEnd(v)
			}
		}
		r.stopEditing()
	}

	innerSlider.OnChanged = func(v float64) {
		s.mu.Lock()
		s.value = v
		callback := s.onChanged
		s.mu.Unlock()

		r.updateValueLabel(v)
		r.updateSteppersState()
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
		r.updateSteppersState()
		if callback != nil {
			callback(v)
		}
	}

	r.applyStyles()
	r.calculateBadgeWidth()
	r.updateValueLabel(currentVal)
	r.updateSteppersState()

	return r
}

type modernSliderRenderer struct {
	widget       *ModernSlider
	titleText    *canvas.Text
	innerSlider  *scrollableSlider
	themedSlider *container.ThemeOverride
	badgeBg      *canvas.Rectangle
	valueText    *canvas.Text
	valueEntry   *widget.Entry
	badgeTap     *badgeTappable
	stepperMinus *sliderStepper
	stepperPlus  *sliderStepper
	minText      *canvas.Text
	maxText      *canvas.Text

	cachedBadgeW float32
	isEditing    bool
}

func (r *modernSliderRenderer) Objects() []fyne.CanvasObject {
	r.widget.mu.RLock()
	hasTitle := r.widget.title != ""
	showVal := r.widget.showValue
	showBdg := r.widget.showBadge
	showStp := r.widget.showSteppers
	showMM := r.widget.showMinMax
	r.widget.mu.RUnlock()

	objs := make([]fyne.CanvasObject, 0, 10)
	if hasTitle {
		objs = append(objs, r.titleText)
	}
	if showStp {
		objs = append(objs, r.stepperMinus, r.stepperPlus)
	}
	objs = append(objs, r.themedSlider)

	if showMM {
		objs = append(objs, r.minText, r.maxText)
	}

	if showVal {
		if showBdg {
			objs = append(objs, r.badgeBg, r.badgeTap)
		}
		if r.isEditing {
			objs = append(objs, r.valueEntry)
		} else {
			objs = append(objs, r.valueText)
		}
	}
	return objs
}

func (r *modernSliderRenderer) Destroy() {}

func (r *modernSliderRenderer) startEditing() {
	if r.widget.Disabled() || !r.widget.editableBadge {
		return
	}
	r.isEditing = true
	r.valueEntry.SetText(strconv.FormatFloat(r.widget.Value(), 'f', -1, 64))
	r.valueEntry.Show()
	r.valueText.Hide()
	r.badgeBg.StrokeColor = ThemeColor(ColorNameAccent)
	r.badgeBg.Refresh()
	r.Layout(r.widget.Size())
	canvas.Refresh(r.widget)

	if c := fyne.CurrentApp().Driver().CanvasForObject(r.widget); c != nil {
		c.Focus(r.valueEntry)
	}
}

func (r *modernSliderRenderer) stopEditing() {
	r.isEditing = false
	r.valueEntry.Hide()
	r.valueText.Show()
	r.applyStyles()
	r.updateValueLabel(r.widget.Value())
	r.Layout(r.widget.Size())
	canvas.Refresh(r.widget)
}

// calculateBadgeWidth evalúa el valor más ancho en los extremos y pasos clave para evitar jittering.
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
		r.cachedBadgeW = fyne.Max(maxTextW+22, 52) // Acolchado amplio para pastilla interactiva
	} else {
		r.cachedBadgeW = maxTextW + 4
	}
}

func (r *modernSliderRenderer) updateValueLabel(v float64) {
	txt := r.widget.formatValue(v)
	r.valueText.Text = txt
	r.valueText.Refresh()
}

func (r *modernSliderRenderer) updateSteppersState() {
	r.widget.mu.RLock()
	val := r.widget.value
	minVal := r.widget.min
	maxVal := r.widget.max
	isDisabled := r.widget.Disabled()
	r.widget.mu.RUnlock()

	r.stepperMinus.SetDisabled(isDisabled || val <= minVal)
	r.stepperPlus.SetDisabled(isDisabled || val >= maxVal)
}

func (r *modernSliderRenderer) MinSize() fyne.Size {
	r.widget.mu.RLock()
	hasTitle := r.widget.title != ""
	showVal := r.widget.showValue
	showStp := r.widget.showSteppers
	showMM := r.widget.showMinMax
	r.widget.mu.RUnlock()

	sliderMin := r.themedSlider.MinSize()
	var totalH float32 = sliderMin.Height
	var totalW float32 = sliderMin.Width

	if showStp {
		totalW += 52 // 22px de cada botón + huecos
	}

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

	if showMM {
		totalH += 14 // Espacio para etiquetas Min/Max
	}

	return fyne.NewSize(fyne.Max(totalW, 180), totalH)
}

func (r *modernSliderRenderer) Layout(size fyne.Size) {
	r.widget.mu.RLock()
	hasTitle := r.widget.title != ""
	showVal := r.widget.showValue
	showBdg := r.widget.showBadge
	showStp := r.widget.showSteppers
	showMM := r.widget.showMinMax
	minVal := r.widget.min
	maxVal := r.widget.max
	minLbl := r.widget.minLabel
	maxLbl := r.widget.maxLabel
	r.widget.mu.RUnlock()

	pad := theme.Padding()
	badgeW := r.cachedBadgeW
	badgeH := float32(24)

	var sliderStartX float32 = 0
	var sliderWidth float32 = size.Width

	// 1. Fila de Cabecera (Título y Badge)
	var sliderTopY float32 = 0
	if hasTitle {
		titleMin := r.titleText.MinSize()
		row1H := fyne.Max(titleMin.Height, badgeH)

		r.titleText.Resize(fyne.NewSize(size.Width-badgeW-pad, titleMin.Height))
		r.titleText.Move(fyne.NewPos(0, (row1H-titleMin.Height)/2))

		if showVal {
			badgeX := size.Width - badgeW
			badgeY := (row1H - badgeH) / 2
			r.layoutBadge(badgeX, badgeY, badgeW, badgeH, showBdg)
		}
		sliderTopY = row1H + 4
	}

	// 2. Fila del Slider y Steppers
	var sliderRowH float32 = size.Height - sliderTopY
	if showMM {
		sliderRowH -= 14
	}
	if sliderRowH < r.themedSlider.MinSize().Height {
		sliderRowH = r.themedSlider.MinSize().Height
	}

	if !hasTitle && showVal {
		sliderWidth = size.Width - badgeW - pad
		if sliderWidth < 80 {
			sliderWidth = 80
		}
		badgeX := size.Width - badgeW
		badgeY := (size.Height - badgeH) / 2
		r.layoutBadge(badgeX, badgeY, badgeW, badgeH, showBdg)
	}

	stepperSize := float32(22)
	stepperGap := float32(6)

	if showStp {
		stepperY := sliderTopY + (sliderRowH-stepperSize)/2
		r.stepperMinus.Resize(fyne.NewSize(stepperSize, stepperSize))
		r.stepperMinus.Move(fyne.NewPos(sliderStartX, stepperY))
		r.stepperMinus.Show()

		sliderStartX += stepperSize + stepperGap
		sliderWidth -= (stepperSize + stepperGap) * 2

		r.stepperPlus.Resize(fyne.NewSize(stepperSize, stepperSize))
		r.stepperPlus.Move(fyne.NewPos(sliderStartX+sliderWidth+stepperGap, stepperY))
		r.stepperPlus.Show()
	} else {
		r.stepperMinus.Hide()
		r.stepperPlus.Hide()
	}

	r.themedSlider.Resize(fyne.NewSize(sliderWidth, sliderRowH))
	r.themedSlider.Move(fyne.NewPos(sliderStartX, sliderTopY))

	// 3. Etiquetas de rango inferior Min / Max
	if showMM {
		if minLbl == "" {
			minLbl = r.widget.formatValue(minVal)
		}
		if maxLbl == "" {
			maxLbl = r.widget.formatValue(maxVal)
		}
		r.minText.Text = minLbl
		r.maxText.Text = maxLbl

		minTextH := r.minText.MinSize().Height
		mmY := sliderTopY + sliderRowH - 2

		r.minText.Resize(fyne.NewSize(sliderWidth/2, minTextH))
		r.minText.Move(fyne.NewPos(sliderStartX, mmY))
		r.minText.Show()

		r.maxText.Resize(fyne.NewSize(sliderWidth/2, minTextH))
		r.maxText.Move(fyne.NewPos(sliderStartX+sliderWidth/2, mmY))
		r.maxText.Show()
	} else {
		r.minText.Hide()
		r.maxText.Hide()
	}
}

func (r *modernSliderRenderer) layoutBadge(x, y, w, h float32, showBdg bool) {
	if showBdg {
		r.badgeBg.Resize(fyne.NewSize(w, h))
		r.badgeBg.Move(fyne.NewPos(x, y))
		r.badgeBg.Show()

		r.badgeTap.Resize(fyne.NewSize(w, h))
		r.badgeTap.Move(fyne.NewPos(x, y))
		r.badgeTap.Show()
	} else {
		r.badgeBg.Hide()
		r.badgeTap.Hide()
	}

	if r.isEditing {
		r.valueEntry.Resize(fyne.NewSize(w, h))
		r.valueEntry.Move(fyne.NewPos(x, y))
		r.valueEntry.Show()
		r.valueText.Hide()
	} else {
		r.valueEntry.Hide()
		valMin := r.valueText.MinSize()
		r.valueText.Resize(fyne.NewSize(w, valMin.Height))
		r.valueText.Move(fyne.NewPos(x, y+(h-valMin.Height)/2))
		r.valueText.Show()
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
	r.themedSlider.Refresh()

	r.titleText.Text = titleStr
	r.applyStyles()
	r.calculateBadgeWidth()
	r.updateValueLabel(val)
	r.updateSteppersState()

	canvas.Refresh(r.widget)
}

func (r *modernSliderRenderer) applyStyles() {
	isDisabled := r.widget.Disabled()

	if isDisabled {
		r.titleText.Color = ThemeColor(ColorNameTextMuted)
		r.badgeBg.FillColor = ThemeColor(ColorNameBtnDefaultBg)
		r.badgeBg.StrokeColor = color.Transparent
		r.valueText.Color = ThemeColor(ColorNameBtnMuted)
		r.minText.Color = ThemeColor(ColorNameBtnMuted)
		r.maxText.Color = ThemeColor(ColorNameBtnMuted)
	} else {
		r.titleText.Color = ThemeColor(ColorNameTextPrimary)
		r.badgeBg.FillColor = ThemeColor(ColorNameBtnPrimaryBg)
		r.badgeBg.StrokeColor = ThemeColor(ColorNameBtnPrimaryBorder)
		r.valueText.Color = ThemeColor(ColorNameBtnPrimaryFg)
		r.minText.Color = ThemeColor(ColorNameTextMuted)
		r.maxText.Color = ThemeColor(ColorNameTextMuted)
	}

	r.titleText.Refresh()
	r.badgeBg.Refresh()
	r.valueText.Refresh()
	r.minText.Refresh()
	r.maxText.Refresh()
}

// ─────────────────────────────────────────────────────────────────────────────
// modernSliderTheme: Estilizado del carril, grosor (5px) y halo en hover
// ─────────────────────────────────────────────────────────────────────────────
type modernSliderTheme struct {
	slider *ModernSlider
}

var _ fyne.Theme = (*modernSliderTheme)(nil)

func (t *modernSliderTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameInputBackground:
		// Pista inactiva del slider: contraste nítido y distinguible sobre cualquier fondo
		t.slider.mu.RLock()
		customTrack := t.slider.trackColor
		t.slider.mu.RUnlock()
		if customTrack != nil {
			return customTrack
		}
		if variant == theme.VariantDark {
			return ColorSlate700 // #334155 en modo oscuro
		}
		// En modo claro: ColorSlate300 (#cbd5e1) destaca con perfecta nitidez sobre tarjetas blancas (#ffffff)
		return ColorSlate300

	case theme.ColorNameForeground:
		// Tirador (Thumb) y pista activa: color primario de la marca o personalizado
		t.slider.mu.RLock()
		customCol := t.slider.color
		t.slider.mu.RUnlock()
		if customCol != nil {
			return customCol
		}
		if pal := GetCurrentThemePalette(); pal != nil && pal.Primary != nil {
			return pal.Primary
		}
		return ThemeColor(ColorNameAccent)

	case theme.ColorNameHover:
		// Halo reactivo suave y translúcido en hover alrededor del tirador (estilo iOS / macOS)
		col := ThemeColor(ColorNameAccent)
		if pal := GetCurrentThemePalette(); pal != nil && pal.Primary != nil {
			col = pal.Primary
		}
		return ColorWithAlpha(col, 0.22)

	case theme.ColorNameFocus:
		// Halo en foco activo
		col := ThemeColor(ColorNameAccent)
		if pal := GetCurrentThemePalette(); pal != nil && pal.Primary != nil {
			col = pal.Primary
		}
		return ColorWithAlpha(col, 0.35)

	case theme.ColorNameDisabled:
		return ThemeColor(ColorNameBtnMuted)
	}

	if curTheme := fyne.CurrentApp().Settings().Theme(); curTheme != nil {
		return curTheme.Color(name, variant)
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (t *modernSliderTheme) Font(style fyne.TextStyle) fyne.Resource {
	if curTheme := fyne.CurrentApp().Settings().Theme(); curTheme != nil {
		return curTheme.Font(style)
	}
	return theme.DefaultTheme().Font(style)
}

func (t *modernSliderTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	if curTheme := fyne.CurrentApp().Settings().Theme(); curTheme != nil {
		return curTheme.Icon(name)
	}
	return theme.DefaultTheme().Icon(name)
}

func (t *modernSliderTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameInputBorder {
		return 2.5 // Espesor óptimo de la pista del slider (5px de alto en lugar de 2px)
	}
	if curTheme := fyne.CurrentApp().Settings().Theme(); curTheme != nil {
		return curTheme.Size(name)
	}
	return theme.DefaultTheme().Size(name)
}

// ─────────────────────────────────────────────────────────────────────────────
// scrollableSlider: Extiende widget.Slider para recibir eventos Scrollable
// directamente en cualquier punto de la pista o thumb sin depender de búsqueda de ancestro.
// ─────────────────────────────────────────────────────────────────────────────
type scrollableSlider struct {
	widget.Slider
	parent *ModernSlider
}

var _ fyne.Scrollable = (*scrollableSlider)(nil)

func newScrollableSlider(min, max float64, parent *ModernSlider) *scrollableSlider {
	s := &scrollableSlider{parent: parent}
	s.Min = min
	s.Max = max
	s.ExtendBaseWidget(s)
	return s
}

func (s *scrollableSlider) Scrolled(e *fyne.ScrollEvent) {
	if s.parent != nil {
		s.parent.Scrolled(e)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// sliderStepper: Micro-botones de paso fino [-] y [+]
// ─────────────────────────────────────────────────────────────────────────────
type sliderStepper struct {
	widget.BaseWidget
	symbol   string
	parent   *ModernSlider
	onTap    func()
	bg       *canvas.Rectangle
	label    *canvas.Text
	hovered  bool
	disabled bool
}

var _ fyne.Tappable = (*sliderStepper)(nil)
var _ desktop.Hoverable = (*sliderStepper)(nil)
var _ fyne.Scrollable = (*sliderStepper)(nil)

func (st *sliderStepper) Scrolled(e *fyne.ScrollEvent) {
	if st.parent != nil {
		st.parent.Scrolled(e)
	}
}

func newSliderStepper(symbol string, parent *ModernSlider, onTap func()) *sliderStepper {
	bg := canvas.NewRectangle(ThemeColor(ColorNameBtnDefaultBg))
	bg.CornerRadius = 4
	bg.StrokeWidth = 1
	bg.StrokeColor = ThemeColor(ColorNameBorder)

	lbl := canvas.NewText(symbol, ThemeColor(ColorNameTextPrimary))
	lbl.TextSize = 13
	lbl.TextStyle = fyne.TextStyle{Bold: true}
	lbl.Alignment = fyne.TextAlignCenter

	st := &sliderStepper{
		symbol: symbol,
		parent: parent,
		onTap:  onTap,
		bg:     bg,
		label:  lbl,
	}
	st.ExtendBaseWidget(st)
	return st
}

func (st *sliderStepper) CreateRenderer() fyne.WidgetRenderer {
	return &sliderStepperRenderer{stepper: st}
}

func (st *sliderStepper) Tapped(_ *fyne.PointEvent) {
	if st.disabled {
		return
	}
	if st.onTap != nil {
		st.onTap()
	}
}

func (st *sliderStepper) MouseIn(_ *desktop.MouseEvent) {
	if st.disabled {
		return
	}
	st.hovered = true
	st.Refresh()
}

func (st *sliderStepper) MouseOut() {
	st.hovered = false
	st.Refresh()
}

func (st *sliderStepper) MouseMoved(_ *desktop.MouseEvent) {}

func (st *sliderStepper) Cursor() desktop.Cursor {
	if st.disabled {
		return desktop.DefaultCursor
	}
	return desktop.PointerCursor
}

func (st *sliderStepper) SetDisabled(d bool) {
	if st.disabled != d {
		st.disabled = d
		st.Refresh()
	}
}

type sliderStepperRenderer struct {
	stepper *sliderStepper
}

func (r *sliderStepperRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.stepper.bg, r.stepper.label}
}

func (r *sliderStepperRenderer) Destroy() {}

func (r *sliderStepperRenderer) MinSize() fyne.Size {
	return fyne.NewSize(22, 22)
}

func (r *sliderStepperRenderer) Layout(size fyne.Size) {
	r.stepper.bg.Resize(size)
	r.stepper.bg.Move(fyne.NewPos(0, 0))

	lblMin := r.stepper.label.MinSize()
	r.stepper.label.Resize(fyne.NewSize(size.Width, lblMin.Height))
	r.stepper.label.Move(fyne.NewPos(0, (size.Height-lblMin.Height)/2))
}

func (r *sliderStepperRenderer) Refresh() {
	st := r.stepper
	if st.disabled {
		st.bg.FillColor = color.Transparent
		st.bg.StrokeColor = ThemeColor(ColorNameBorder)
		st.label.Color = ThemeColor(ColorNameBtnMuted)
	} else if st.hovered {
		st.bg.FillColor = ThemeColor(ColorNameBtnDefaultHoverBg)
		st.bg.StrokeColor = ThemeColor(ColorNameAccent)
		st.label.Color = ThemeColor(ColorNameAccent)
	} else {
		st.bg.FillColor = ThemeColor(ColorNameBtnDefaultBg)
		st.bg.StrokeColor = ThemeColor(ColorNameBorder)
		st.label.Color = ThemeColor(ColorNameTextPrimary)
	}
	st.bg.Refresh()
	st.label.Refresh()
}

// ─────────────────────────────────────────────────────────────────────────────
// badgeTappable: Área táctil sobre la pastilla de valor
// ─────────────────────────────────────────────────────────────────────────────
type badgeTappable struct {
	widget.BaseWidget
	parent *ModernSlider
	onTap  func()
}

var _ fyne.Tappable = (*badgeTappable)(nil)
var _ fyne.Scrollable = (*badgeTappable)(nil)

func (b *badgeTappable) Scrolled(e *fyne.ScrollEvent) {
	if b.parent != nil {
		b.parent.Scrolled(e)
	}
}

func newBadgeTappable(parent *ModernSlider, onTap func()) *badgeTappable {
	bt := &badgeTappable{parent: parent, onTap: onTap}
	bt.ExtendBaseWidget(bt)
	return bt
}

func (b *badgeTappable) CreateRenderer() fyne.WidgetRenderer {
	return &badgeTappableRenderer{bt: b}
}

func (b *badgeTappable) Tapped(_ *fyne.PointEvent) {
	if b.onTap != nil {
		b.onTap()
	}
}

func (b *badgeTappable) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

type badgeTappableRenderer struct {
	bt *badgeTappable
}

func (r *badgeTappableRenderer) Objects() []fyne.CanvasObject { return nil }
func (r *badgeTappableRenderer) Destroy()                    {}
func (r *badgeTappableRenderer) MinSize() fyne.Size          { return fyne.NewSize(0, 0) }
func (r *badgeTappableRenderer) Layout(_ fyne.Size)          {}
func (r *badgeTappableRenderer) Refresh()                    {}
