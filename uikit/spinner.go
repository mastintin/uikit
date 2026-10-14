package uikit

import (
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Constantes de dimensiones estándar para ModernSpinner
const (
	SpinnerSmall  float32 = 16 // Adecuado para botones, badges y chips
	SpinnerMedium float32 = 24 // Tamaño estándar para tarjetas y cabeceras
	SpinnerLarge  float32 = 36 // Paneles de carga y modales
	SpinnerHero   float32 = 48 // Pantallas completas de carga o splash

	defaultCutoutRatio = float32(0.72) // Relación radio interior/exterior (grosor de trazo óptimo)
	defaultAnimDuration = 1100 * time.Millisecond
	spinnerMinSweep     = float32(24)
	spinnerMaxSweep     = float32(280)
	spinnerGrowRange    = spinnerMaxSweep - spinnerMinSweep
	spinnerCycleWrap    = float32(16)
)

var _ fyne.Widget = (*ModernSpinner)(nil)

// ModernSpinner es un indicador circular animado moderno con anillo base sutil (track)
// y arco rotativo dinámico con curva suave inchworm.
type ModernSpinner struct {
	widget.BaseWidget

	mu          sync.Mutex
	size        float32
	color       color.Color
	trackColor  color.Color
	showTrack   bool
	strokeRatio float32
	started     bool
}

// NewModernSpinner crea un nuevo spinner circular moderno con tamaño configurable.
// Si no se especifica tamaño, adopta SpinnerMedium (24px).
func NewModernSpinner(size ...float32) *ModernSpinner {
	sz := SpinnerMedium
	if len(size) > 0 && size[0] > 0 {
		sz = size[0]
	}

	s := &ModernSpinner{
		size:        sz,
		showTrack:   true,
		strokeRatio: defaultCutoutRatio,
	}
	s.ExtendBaseWidget(s)
	return s
}

// NewSpinner es un alias conciso de NewModernSpinner.
func NewSpinner(size ...float32) *ModernSpinner {
	return NewModernSpinner(size...)
}

// NewSpinnerSmall crea un spinner de 16px optimizado para botones y líneas compactas.
func NewSpinnerSmall() *ModernSpinner {
	return NewModernSpinner(SpinnerSmall)
}

// NewSpinnerMedium crea un spinner estándar de 24px.
func NewSpinnerMedium() *ModernSpinner {
	return NewModernSpinner(SpinnerMedium)
}

// NewSpinnerLarge crea un spinner prominente de 36px.
func NewSpinnerLarge() *ModernSpinner {
	return NewModernSpinner(SpinnerLarge)
}

// NewPrimarySpinner crea un spinner estilizado con el color primario de la aplicación (#2563eb).
func NewPrimarySpinner(size ...float32) *ModernSpinner {
	return NewModernSpinner(size...).WithColor(HexColor("#2563eb"))
}

// NewSuccessSpinner crea un spinner estilizado con el color verde esmeralda (#10b981).
func NewSuccessSpinner(size ...float32) *ModernSpinner {
	return NewModernSpinner(size...).WithColor(HexColor("#10b981"))
}

// NewDangerSpinner crea un spinner estilizado con el color rojo peligro (#ef4444).
func NewDangerSpinner(size ...float32) *ModernSpinner {
	return NewModernSpinner(size...).WithColor(HexColor("#ef4444"))
}

// NewAccentSpinner crea un spinner estilizado con el color ámbar/acento (#f59e0b).
func NewAccentSpinner(size ...float32) *ModernSpinner {
	return NewModernSpinner(size...).WithColor(HexColor("#f59e0b"))
}

// WithSize ajusta fluidamente el tamaño exterior en píxeles.
func (s *ModernSpinner) WithSize(size float32) *ModernSpinner {
	s.mu.Lock()
	if size > 0 {
		s.size = size
	}
	s.mu.Unlock()
	s.Refresh()
	return s
}

// SetSize actualiza el tamaño exterior del spinner.
func (s *ModernSpinner) SetSize(size float32) *ModernSpinner {
	return s.WithSize(size)
}

// WithColor define el color principal del segmento móvil.
func (s *ModernSpinner) WithColor(c color.Color) *ModernSpinner {
	s.mu.Lock()
	s.color = c
	s.mu.Unlock()
	s.Refresh()
	return s
}

// SetColor define el color principal del segmento móvil.
func (s *ModernSpinner) SetColor(c color.Color) *ModernSpinner {
	return s.WithColor(c)
}

// WithHexColor define el color principal a partir de una cadena hexadecimal (ej. "#2563eb").
func (s *ModernSpinner) WithHexColor(hex string) *ModernSpinner {
	return s.WithColor(HexColor(hex))
}

// WithTrackColor define manualmente el color del anillo de fondo (track ring).
func (s *ModernSpinner) WithTrackColor(c color.Color) *ModernSpinner {
	s.mu.Lock()
	s.trackColor = c
	s.mu.Unlock()
	s.Refresh()
	return s
}

// WithTrack activa o desactiva la visibilidad del anillo de fondo tenue.
func (s *ModernSpinner) WithTrack(show bool) *ModernSpinner {
	s.mu.Lock()
	s.showTrack = show
	s.mu.Unlock()
	s.Refresh()
	return s
}

// WithStrokeRatio ajusta el grosor relativo del trazo (entre 0.0 y 0.95).
func (s *ModernSpinner) WithStrokeRatio(ratio float32) *ModernSpinner {
	s.mu.Lock()
	if ratio >= 0.1 && ratio <= 0.95 {
		s.strokeRatio = ratio
	}
	s.mu.Unlock()
	s.Refresh()
	return s
}

// Start inicia la animación de giro. Es seguro invocarlo concurrentemente.
func (s *ModernSpinner) Start() *ModernSpinner {
	s.mu.Lock()
	if s.started {
		s.mu.Unlock()
		return s
	}
	s.started = true
	s.mu.Unlock()

	s.Refresh()
	return s
}

// Stop detiene la animación de giro.
func (s *ModernSpinner) Stop() *ModernSpinner {
	s.mu.Lock()
	if !s.started {
		s.mu.Unlock()
		return s
	}
	s.started = false
	s.mu.Unlock()

	s.Refresh()
	return s
}

// Toggle alterna el estado activo/detenido del spinner.
func (s *ModernSpinner) Toggle() *ModernSpinner {
	s.mu.Lock()
	active := s.started
	s.mu.Unlock()

	if active {
		return s.Stop()
	}
	return s.Start()
}

// IsRunning indica si el spinner está actualmente animándose.
func (s *ModernSpinner) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.started
}

// MinSize calcula el tamaño mínimo cuadrado requerido por el spinner.
func (s *ModernSpinner) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	s.mu.Lock()
	sz := s.size
	if sz <= 0 {
		sz = SpinnerMedium
	}
	s.mu.Unlock()
	return fyne.NewSquareSize(sz)
}

// CreateRenderer inicializa el renderer del widget y su bucle de animación.
func (s *ModernSpinner) CreateRenderer() fyne.WidgetRenderer {
	s.mu.Lock()
	ratio := s.strokeRatio
	s.mu.Unlock()

	// 1. Anillo base estático (track background)
	trackArc := canvas.NewArc(0, 360, ratio, color.Transparent)

	// 2. Arco móvil animado
	spinArc := canvas.NewArc(0, 0, ratio, color.Transparent)

	r := &modernSpinnerRenderer{
		spinner:  s,
		trackArc: trackArc,
		spinArc:  spinArc,
	}

	r.anim = &fyne.Animation{
		Duration:    defaultAnimDuration,
		RepeatCount: fyne.AnimationRepeatForever,
		Curve:       fyne.AnimationLinear,
		Tick:        r.animate,
	}

	r.updateColors()

	if s.IsRunning() {
		r.start()
	}

	return r
}

type modernSpinnerRenderer struct {
	spinner  *ModernSpinner
	trackArc *canvas.Arc
	spinArc  *canvas.Arc
	anim     *fyne.Animation

	wasStarted      bool
	completedCycles float32
	lastDone        float32
}

func (r *modernSpinnerRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.trackArc, r.spinArc}
}

func (r *modernSpinnerRenderer) MinSize() fyne.Size {
	return r.spinner.MinSize()
}

func (r *modernSpinnerRenderer) Layout(size fyne.Size) {
	rad := fyne.Min(size.Width, size.Height) / 2
	if rad <= 0 {
		return
	}
	mid := fyne.NewPos(size.Width/2, size.Height/2)
	arcSize := fyne.NewSquareSize(rad * 2)
	arcPos := mid.Subtract(fyne.NewSquareOffsetPos(rad))

	r.trackArc.Resize(arcSize)
	r.trackArc.Move(arcPos)

	r.spinArc.Resize(arcSize)
	r.spinArc.Move(arcPos)

	if r.spinner.IsRunning() && !fyne.CurrentApp().Settings().ShowAnimations() {
		r.drawStaticArc()
	}
}

func (r *modernSpinnerRenderer) Refresh() {
	r.spinner.mu.Lock()
	running := r.spinner.started
	ratio := r.spinner.strokeRatio
	r.spinner.mu.Unlock()

	r.trackArc.CutoutRatio = ratio
	r.spinArc.CutoutRatio = ratio

	if running {
		if !r.wasStarted {
			r.start()
		}
	} else if r.wasStarted {
		r.stop()
	}

	r.updateColors()
	canvas.Refresh(r.trackArc)
	canvas.Refresh(r.spinArc)
}

func (r *modernSpinnerRenderer) Destroy() {
	r.spinner.mu.Lock()
	r.spinner.started = false
	r.spinner.mu.Unlock()
	r.stop()
}

func (r *modernSpinnerRenderer) updateColors() {
	r.spinner.mu.Lock()
	mainCol := r.spinner.color
	trackCol := r.spinner.trackColor
	showTrack := r.spinner.showTrack
	r.spinner.mu.Unlock()

	if mainCol == nil {
		if curTheme := fyne.CurrentApp().Settings().Theme(); curTheme != nil {
			v := fyne.CurrentApp().Settings().ThemeVariant()
			mainCol = curTheme.Color(theme.ColorNamePrimary, v)
		}
	}
	if mainCol == nil {
		mainCol = color.NRGBA{R: 37, G: 99, B: 235, A: 255}
	}

	if showTrack {
		if trackCol != nil {
			r.trackArc.FillColor = trackCol
		} else {
			r.trackArc.FillColor = ColorWithAlpha(mainCol, 0.16)
		}
	} else {
		r.trackArc.FillColor = color.Transparent
	}

	r.spinArc.FillColor = mainCol
}

func (r *modernSpinnerRenderer) animate(done float32) {
	if done < r.lastDone {
		r.completedCycles++
		if r.completedCycles >= spinnerCycleWrap {
			r.completedCycles -= spinnerCycleWrap
		}
	}
	r.lastDone = done

	cycleAdvance := r.completedCycles * spinnerGrowRange
	spin := (r.completedCycles + done) * 360
	base := cycleAdvance + spin

	var tailFrac, headFrac float32
	if done < 0.5 {
		headFrac = easeInOut(done * 2)
	} else {
		tailFrac = easeInOut((done - 0.5) * 2)
		headFrac = 1
	}

	r.spinArc.StartAngle = base + tailFrac*spinnerGrowRange
	r.spinArc.EndAngle = base + spinnerMinSweep + headFrac*spinnerGrowRange
	r.spinArc.Refresh()
}

// easeInOut implementa un suavizado cúbico monótono hermite (3t^2 - 2t^3)
func easeInOut(t float32) float32 {
	if t <= 0 {
		return 0
	}
	if t >= 1 {
		return 1
	}
	return t * t * (3 - 2*t)
}

func (r *modernSpinnerRenderer) start() {
	r.wasStarted = true
	r.completedCycles = 0
	r.lastDone = 0

	if !fyne.CurrentApp().Settings().ShowAnimations() {
		r.drawStaticArc()
		return
	}
	r.anim.Start()
}

func (r *modernSpinnerRenderer) stop() {
	r.wasStarted = false
	r.anim.Stop()
	r.spinArc.StartAngle = 0
	r.spinArc.EndAngle = 0
	r.spinArc.Refresh()
}

func (r *modernSpinnerRenderer) drawStaticArc() {
	r.spinArc.StartAngle = 0
	r.spinArc.EndAngle = spinnerMaxSweep
	r.spinArc.Refresh()
}
