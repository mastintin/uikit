package uikit

import (
	"image/color"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// ToastVariant define las variantes semánticas disponibles para el toast.
type ToastVariant int

const (
	ToastInfo ToastVariant = iota
	ToastSuccess
	ToastWarning
	ToastDanger
	ToastError = ToastDanger
)

// Alias de tipo para compatibilidad y rapidez
type Toast = ModernToast

// ParseToastVariant convierte un valor cualquiera (string o ToastVariant) a ToastVariant.
func ParseToastVariant(v any) ToastVariant {
	switch val := v.(type) {
	case ToastVariant:
		return val
	case string:
		s := strings.ToLower(strings.TrimSpace(val))
		switch s {
		case "success", "ok":
			return ToastSuccess
		case "error", "danger", "err":
			return ToastDanger
		case "warning", "warn", "alerta":
			return ToastWarning
		default:
			return ToastInfo
		}
	default:
		return ToastInfo
	}
}

// ToastVariantStyle permite definir o sobrescribir iconos y colores para una variante.
type ToastVariantStyle struct {
	IconName    string
	BorderColor color.Color
	IconColor   color.Color
	TextColor   color.Color
}

type toastQueuedMsg struct {
	variant any
	message string
}

// ModernToast es una notificación flotante estilizada, no invasiva y con temporizador de auto-cierre,
// animación elástica de entrada, cierre instantáneo al hacer clic (dismiss) y soporte para cola de mensajes.
// Es 100% independiente del dominio de negocio y reutilizable en cualquier aplicación Fyne.
type ModernToast struct {
	card       *ModernCard
	wrapper    *fyne.Container
	tapWrapper *tappableToast
	icon       *ModernIcon
	label      *ModernLabel
	duration   time.Duration
	timer      *time.Timer
	showAnim   *fyne.Animation
	isVisible  bool
	queue      []toastQueuedMsg
	mu         sync.Mutex

	customStyles map[ToastVariant]*ToastVariantStyle
}

// tappableToast envuelve la tarjeta para permitir descartar la notificación al hacer clic.
type tappableToast struct {
	widget.BaseWidget
	content  fyne.CanvasObject
	onTapped func()
}

var _ desktop.Hoverable = (*tappableToast)(nil)
var _ fyne.Tappable = (*tappableToast)(nil)

func newTappableToast(content fyne.CanvasObject, onTapped func()) *tappableToast {
	t := &tappableToast{
		content:  content,
		onTapped: onTapped,
	}
	t.ExtendBaseWidget(t)
	return t
}

func (t *tappableToast) Tapped(_ *fyne.PointEvent) {
	if t.onTapped != nil {
		t.onTapped()
	}
}

func (t *tappableToast) TappedSecondary(_ *fyne.PointEvent) {}

func (t *tappableToast) MouseIn(_ *desktop.MouseEvent)    {}
func (t *tappableToast) MouseMoved(_ *desktop.MouseEvent) {}
func (t *tappableToast) MouseOut()                        {}

func (t *tappableToast) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (t *tappableToast) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.content)
}

// NewModernToast crea un nuevo componente toast flotante listo para añadir al layout.
func NewModernToast() *ModernToast {
	lbl := NewModernLabel("").SetSize(12.5).SetBold(true).SetColor(ThemeColor(ColorNameTextPrimary))
	ic := NewColoredIcon("info", 16, "#38bdf8")

	content := container.NewHBox(
		container.NewCenter(ic),
		lbl,
	)

	card := NewCard(content).WithBorderColor(ColorPrimary, 1)

	var mt *ModernToast
	tapWrap := newTappableToast(card, func() {
		if mt != nil {
			mt.Hide()
		}
	})

	wrapper := container.NewCenter(tapWrap)
	wrapper.Hide()

	mt = &ModernToast{
		card:         card,
		wrapper:      wrapper,
		tapWrapper:   tapWrap,
		icon:         ic,
		label:        lbl,
		duration:     3500 * time.Millisecond,
		customStyles: make(map[ToastVariant]*ToastVariantStyle),
	}
	return mt
}

// NewToast es un constructor alternativo rápido de ModernToast.
func NewToast() *ModernToast {
	return NewModernToast()
}

// Container devuelve el contenedor centrador del toast para incrustarlo en la ventana o stack de la app.
func (t *ModernToast) Container() *fyne.Container {
	return t.wrapper
}

// Card devuelve la tarjeta interna ModernCard por si se desea ajustar radios o padding adicional.
func (t *ModernToast) Card() *ModernCard {
	return t.card
}

// SetDuration configura el tiempo que permanece visible el toast antes de desvanecerse (por defecto 3.5s).
func (t *ModernToast) SetDuration(d time.Duration) *ModernToast {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.duration = d
	return t
}

// SetCustomColors permite a una aplicación sobrescribir los colores de una variante específica.
func (t *ModernToast) SetCustomColors(variant ToastVariant, border, icon, text color.Color) *ModernToast {
	t.mu.Lock()
	defer t.mu.Unlock()

	style, exists := t.customStyles[variant]
	if !exists {
		style = &ToastVariantStyle{}
		t.customStyles[variant] = style
	}
	if border != nil {
		style.BorderColor = border
	}
	if icon != nil {
		style.IconColor = icon
	}
	if text != nil {
		style.TextColor = text
	}
	return t
}

// SetVariantStyle configura el estilo completo (icono y colores) de una variante.
func (t *ModernToast) SetVariantStyle(variant ToastVariant, style ToastVariantStyle) *ModernToast {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.customStyles[variant] = &style
	return t
}

// Queue añade un mensaje a la cola para mostrarse secuencialmente después de que finalice el actual.
func (t *ModernToast) Queue(variant any, message string) {
	t.mu.Lock()
	if !t.isVisible {
		t.mu.Unlock()
		t.Show(variant, message)
		return
	}
	t.queue = append(t.queue, toastQueuedMsg{variant: variant, message: message})
	t.mu.Unlock()
}

// Show muestra la notificación con el mensaje y la variante semántica especificada.
// Soporta tanto Show(ToastSuccess, "Texto") como Show("success", "Texto").
func (t *ModernToast) Show(variant any, message string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.timer != nil {
		t.timer.Stop()
		t.timer = nil
	}

	t.isVisible = true
	v := ParseToastVariant(variant)

	// Resolver estilos por defecto según la variante semántica.
	var (
		iconName  string
		borderCol color.Color
		iconCol   color.Color
		textCol   color.Color = ThemeColor(ColorNameTextPrimary)
	)

	switch v {
	case ToastSuccess:
		iconName = "check"
		borderCol = ThemeColor(ColorNameSuccess)
		if borderCol == nil {
			borderCol = ColorSuccess
		}
		iconCol = borderCol
	case ToastDanger:
		iconName = "error"
		borderCol = ThemeColor(ColorNameDanger)
		if borderCol == nil {
			borderCol = ColorDanger
		}
		iconCol = borderCol
	case ToastWarning:
		iconName = "warning"
		borderCol = ThemeColor(ColorNameWarning)
		if borderCol == nil {
			borderCol = ColorWarning
		}
		iconCol = borderCol
	default: // ToastInfo
		iconName = "info"
		borderCol = ThemeColor(ColorNameAccent)
		if borderCol == nil {
			borderCol = ColorPrimary
		}
		iconCol = borderCol
	}

	// Aplicar sobrescrituras personalizadas si existen
	if custom, ok := t.customStyles[v]; ok && custom != nil {
		if custom.IconName != "" {
			iconName = custom.IconName
		}
		if custom.BorderColor != nil {
			borderCol = custom.BorderColor
		}
		if custom.IconColor != nil {
			iconCol = custom.IconColor
		}
		if custom.TextColor != nil {
			textCol = custom.TextColor
		}
	}

	// Configurar widgets
	t.label.SetText(message).SetColor(textCol)
	t.icon.SetIcon(iconName)
	if iconCol != nil {
		t.icon.SetColor(iconCol)
	}
	t.card.WithBorderColor(borderCol, 1)

	t.card.Refresh()
	t.wrapper.Show()
	t.wrapper.Refresh()

	// Animación suave de entrada elástica (160ms) si las animaciones están activas
	if app := fyne.CurrentApp(); app != nil && app.Settings() != nil && app.Settings().ShowAnimations() {
		if t.showAnim != nil {
			t.showAnim.Stop()
		}
		const defaultPadH float32 = 18
		t.card.WithHorizontalPadding(defaultPadH*0.7, defaultPadH*0.7)
		t.card.Refresh()

		t.showAnim = fyne.NewAnimation(160*time.Millisecond, func(done float32) {
			curPad := defaultPadH * (0.7 + 0.3*fyne.AnimationEaseOut(done))
			t.card.WithHorizontalPadding(curPad, curPad)
			t.card.Refresh()
		})
		t.showAnim.Curve = fyne.AnimationEaseOut
		t.showAnim.Start()
	}

	// Programar auto-ocultación segura con fyne.Do
	duration := t.duration
	if duration <= 0 {
		duration = 3500 * time.Millisecond
	}
	t.timer = time.AfterFunc(duration, func() {
		fyne.Do(func() { t.hideUI() })
	})
}

// hideUI realiza las operaciones de UI de ocultación. Debe llamarse siempre desde
// el hilo de Fyne (directamente o envuelto en fyne.Do).
func (t *ModernToast) hideUI() {
	t.mu.Lock()
	if t.timer != nil {
		t.timer.Stop()
		t.timer = nil
	}
	if t.showAnim != nil {
		t.showAnim.Stop()
	}
	wasVisible := t.isVisible
	t.isVisible = false

	// Si hay mensajes encolados, despachar el siguiente
	if len(t.queue) > 0 {
		next := t.queue[0]
		t.queue = t.queue[1:]
		t.mu.Unlock()
		t.Show(next.variant, next.message)
		return
	}
	t.mu.Unlock()

	if wasVisible && t.wrapper != nil {
		t.wrapper.Hide()
	}
}

// Hide oculta inmediatamente el toast y cancela cualquier temporizador activo.
// Es seguro llamarlo desde cualquier goroutine.
func (t *ModernToast) Hide() {
	fyne.Do(func() { t.hideUI() })
}

// ShowPopUp muestra un toast flotante en la capa de Overlay sobre el canvas de una ventana,
// asegurando que flote por encima de cualquier otro elemento o modal.
func ShowPopUp(c fyne.Canvas, variant any, message string, duration ...time.Duration) *widget.PopUp {
	v := ParseToastVariant(variant)
	var (
		iconName  = "info"
		borderCol = ThemeColor(ColorNameAccent)
		iconCol   color.Color
		textCol   = ThemeColor(ColorNameTextPrimary)
	)
	if borderCol == nil {
		borderCol = ColorPrimary
	}
	iconCol = borderCol

	switch v {
	case ToastSuccess:
		iconName = "check"
		borderCol = ThemeColor(ColorNameSuccess)
		if borderCol == nil {
			borderCol = ColorSuccess
		}
		iconCol = borderCol
	case ToastDanger:
		iconName = "error"
		borderCol = ThemeColor(ColorNameDanger)
		if borderCol == nil {
			borderCol = ColorDanger
		}
		iconCol = borderCol
	case ToastWarning:
		iconName = "warning"
		borderCol = ThemeColor(ColorNameWarning)
		if borderCol == nil {
			borderCol = ColorWarning
		}
		iconCol = borderCol
	}

	lbl := NewModernLabel(message).SetSize(12.5).SetBold(true).SetColor(textCol)
	ic := NewIcon(iconName, 16, "", nil)
	if iconCol != nil {
		ic.SetColor(iconCol)
	}

	content := container.NewHBox(
		container.NewCenter(ic),
		lbl,
	)
	card := NewCard(content).WithBorderColor(borderCol, 1).WithRadius(10)

	var pop *widget.PopUp
	tapWrapper := newTappableToast(card, func() {
		if pop != nil {
			pop.Hide()
		}
	})

	pop = widget.NewPopUp(tapWrapper, c)
	dur := 3500 * time.Millisecond
	if len(duration) > 0 && duration[0] > 0 {
		dur = duration[0]
	}

	time.AfterFunc(dur, func() {
		fyne.Do(func() {
			pop.Hide()
		})
	})

	pop.Show()
	return pop
}
