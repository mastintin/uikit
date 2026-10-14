package uikit

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// BadgeVariant define las variantes visuales predefinidas para badges y pastillas.
type BadgeVariant int

const (
	BadgeNeutral BadgeVariant = iota
	BadgeSuccess
	BadgeDanger
	BadgeInfo
	BadgeWarning
	BadgePurple
	BadgePrimary = BadgeInfo
)

// Alias de tipo para compatibilidad
type Badge = ModernBadge

// NewBadge es el constructor unificado rápido estilo xBase / Harbour:
// Soporta tanto NewBadge(BadgeSuccess, "Texto") como NewBadge("Texto", BadgeSuccess),
// además de compatibilidad directa con widget.Importance de Fyne.
func NewBadge(arg1 interface{}, arg2 interface{}, isPill ...bool) *ModernBadge {
	b := NewModernBadge()

	var v BadgeVariant = BadgeNeutral
	var text string

	switch a1 := arg1.(type) {
	case BadgeVariant:
		v = a1
		if s, ok := arg2.(string); ok {
			text = s
		}
	case widget.Importance:
		v = ImportanceToBadgeVariant(a1)
		if s, ok := arg2.(string); ok {
			text = s
		}
	case string:
		switch a2 := arg2.(type) {
		case BadgeVariant:
			v = a2
			text = a1
		case widget.Importance:
			v = ImportanceToBadgeVariant(a2)
			text = a1
		case string:
			p1 := ParseBadgeVariant(a1)
			p2 := ParseBadgeVariant(a2)
			if p1 != BadgeNeutral {
				v = p1
				text = a2
			} else {
				v = p2
				text = a1
			}
		default:
			text = a1
		}
	default:
		v = ParseBadgeVariant(arg1)
		if s, ok := arg2.(string); ok {
			text = s
		}
	}

	b.variant = v
	b.text = text
	if len(isPill) > 0 && isPill[0] {
		b.isPill = true
		b.height = 24
		b.paddingX = 10
		b.label.TextSize = 12.5
	}
	b.applyVariantColors()
	b.Refresh()
	return b
}


// ImportanceToBadgeVariant convierte un widget.Importance nativo de Fyne a BadgeVariant.
func ImportanceToBadgeVariant(imp widget.Importance) BadgeVariant {
	switch imp {
	case widget.DangerImportance:
		return BadgeDanger
	case widget.SuccessImportance:
		return BadgeSuccess
	case widget.WarningImportance:
		return BadgeWarning
	case widget.HighImportance:
		return BadgeInfo
	default:
		return BadgeNeutral
	}
}

// ParseBadgeVariant convierte nombres comunes de variantes o tipos Fyne a BadgeVariant.
func ParseBadgeVariant(v interface{}) BadgeVariant {
	switch val := v.(type) {
	case BadgeVariant:
		return val
	case widget.Importance:
		return ImportanceToBadgeVariant(val)
	case string:
		switch strings.ToLower(strings.TrimSpace(val)) {
		case "success", "verde", "green", "ok":
			return BadgeSuccess
		case "danger", "rojo", "red", "error":
			return BadgeDanger
		case "info", "primary", "azul", "blue":
			return BadgeInfo
		case "warning", "ambar", "amber", "alerta", "yellow":
			return BadgeWarning
		case "purple", "violeta", "morado":
			return BadgePurple
		default:
			return BadgeNeutral
		}
	default:
		return BadgeNeutral
	}
}

// ModernBadge es un widget genérico de alta precisión visual para etiquetas,
// pastillas de estado (badges) y contadores redondeados (pills) estilo goCuadros.
type ModernBadge struct {
	widget.BaseWidget

	bg    *canvas.Rectangle
	label *canvas.Text

	text          string
	variant       BadgeVariant
	isCustom      bool
	bgFill        color.Color
	textCol       color.Color
	strokeCol     color.Color
	strokeWidth   float32
	isPill        bool
	cornerRadius  float32
	height        float32
	minWidth      float32
	paddingX      float32
	hideWhenEmpty bool
	alignment     *fyne.TextAlign
}

// Alias de tipos para mantener compatibilidad con código existente
type CountBadge = ModernBadge

// NewModernBadge crea un nuevo badge totalmente configurable.
func NewModernBadge() *ModernBadge {
	b := &ModernBadge{
		bg:            canvas.NewRectangle(color.Transparent),
		label:         canvas.NewText("", color.Transparent),
		variant:       BadgeNeutral,
		strokeWidth:   1,
		cornerRadius:  4,
		height:        18,
		minWidth:      0,
		paddingX:      7,
		hideWhenEmpty: false,
	}

	b.label.TextSize = 11.5
	b.label.TextStyle = fyne.TextStyle{Bold: true}
	b.label.Alignment = fyne.TextAlignCenter

	b.applyVariantColors()
	b.ExtendBaseWidget(b)
	return b
}

// NewCustomBadge permite crear un badge con colores hexadecimales directos y medidas libres.
func NewCustomBadge(text, bgHex, textHex, borderHex string, cornerRadius, height float32) *ModernBadge {
	b := NewModernBadge()
	b.text = text
	b.isCustom = true
	b.bgFill = HexColor(bgHex)
	b.textCol = HexColor(textHex)
	b.strokeCol = HexColor(borderHex)
	b.cornerRadius = cornerRadius
	b.height = height
	b.Refresh()
	return b
}

// Sugar Constructors (constructores semánticos listos para usar)
func NewSuccessBadge(text string) *ModernBadge {
	return NewBadge(BadgeSuccess, text, false)
}

func NewDangerBadge(text string) *ModernBadge {
	return NewBadge(BadgeDanger, text, false)
}

func NewInfoBadge(text string) *ModernBadge {
	return NewBadge(BadgeInfo, text, false)
}

func NewWarningBadge(text string) *ModernBadge {
	return NewBadge(BadgeWarning, text, false)
}

func NewNeutralBadge(text string) *ModernBadge {
	return NewBadge(BadgeNeutral, text, false)
}

func NewPillBadge(variant interface{}, text string) *ModernBadge {
	return NewBadge(variant, text, true)
}


// NewCountBadge crea la pastilla del contador del Navbar con adaptación dinámica de tamaño,
// bordes redondeados y estilo moderno (#eff6ff / #2563eb / borde #bfdbfe).
func NewCountBadge() *ModernBadge {
	b := NewBadge(BadgeInfo, "", true)
	b.height = 24
	b.minWidth = 80
	b.paddingX = 10
	b.label.TextSize = 12.5
	return b
}

// NewIDBadge crea una etiqueta para códigos de identificación (.detail-id-tag):
// - Fuente monospace negrita
// - Adaptación automática a modo claro/oscuro
func NewIDBadge(id string) *ModernBadge {
	b := NewModernBadge()
	b.text = id
	b.height = 26
	b.paddingX = 10
	b.cornerRadius = 6
	b.variant = BadgeNeutral
	b.strokeWidth = 1
	b.label.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	b.label.TextSize = 13
	b.applyVariantColors()
	b.Refresh()
	return b
}

// Métodos de configuración y fluidez

func (b *ModernBadge) SetText(text string) *ModernBadge {
	b.text = text
	b.Refresh()
	return b
}

func (b *ModernBadge) Text() string {
	return b.text
}

func (b *ModernBadge) SetVariant(variant BadgeVariant) *ModernBadge {
	b.variant = variant
	b.isCustom = false
	b.applyVariantColors()
	b.Refresh()
	return b
}

func (b *ModernBadge) SetHexColors(bgHex, textHex, borderHex string) *ModernBadge {
	b.isCustom = true
	b.bgFill = HexColor(bgHex)
	b.textCol = HexColor(textHex)
	b.strokeCol = HexColor(borderHex)
	b.Refresh()
	return b
}

func (b *ModernBadge) SetColors(bg, text, border color.Color) *ModernBadge {
	b.isCustom = true
	b.bgFill = bg
	b.textCol = text
	b.strokeCol = border
	b.Refresh()
	return b
}

func (b *ModernBadge) SetPill(isPill bool) *ModernBadge {
	b.isPill = isPill
	if isPill && b.height < 24 {
		b.height = 24
	}
	b.Refresh()
	return b
}

func (b *ModernBadge) SetCornerRadius(radius float32) *ModernBadge {
	b.cornerRadius = radius
	b.Refresh()
	return b
}

func (b *ModernBadge) SetHeight(height float32) *ModernBadge {
	b.height = height
	b.Refresh()
	return b
}

func (b *ModernBadge) SetMinWidth(minWidth float32) *ModernBadge {
	b.minWidth = minWidth
	b.Refresh()
	return b
}

func (b *ModernBadge) SetPaddingX(paddingX float32) *ModernBadge {
	b.paddingX = paddingX
	b.Refresh()
	return b
}

func (b *ModernBadge) SetHideWhenEmpty(hide bool) *ModernBadge {
	b.hideWhenEmpty = hide
	b.Refresh()
	return b
}

// SetAlignment permite alinear horizontalmente la pastilla dentro de su espacio asignado
// (por ejemplo: fyne.TextAlignLeading, fyne.TextAlignCenter, fyne.TextAlignTrailing).
func (b *ModernBadge) SetAlignment(align fyne.TextAlign) *ModernBadge {
	b.alignment = &align
	b.Refresh()
	return b
}

// Alignment devuelve la alineación configurada o (TextAlignCenter, false) si no se configuró explícitamente.
func (b *ModernBadge) Alignment() (fyne.TextAlign, bool) {
	if b.alignment == nil {
		return fyne.TextAlignCenter, false
	}
	return *b.alignment, true
}

// SetImportance mapea y asigna una importancia estándar de Fyne a la variante visual del badge.
func (b *ModernBadge) SetImportance(imp widget.Importance) *ModernBadge {
	return b.SetVariant(ImportanceToBadgeVariant(imp))
}

func (b *ModernBadge) applyVariantColors() {
	if b.isCustom {
		return
	}
	pal := GetCurrentThemePalette()
	if pal == nil {
		// Fallback resiliente al tema nativo de Fyne si no hay paleta cargada
		th := fyne.CurrentApp().Settings().Theme()
		v := fyne.CurrentApp().Settings().ThemeVariant()
		switch b.variant {
		case BadgeSuccess:
			b.bgFill = th.Color(theme.ColorNameSuccess, v)
			b.textCol = th.Color(theme.ColorNameForegroundOnSuccess, v)
		case BadgeDanger:
			b.bgFill = th.Color(theme.ColorNameError, v)
			b.textCol = th.Color(theme.ColorNameForegroundOnError, v)
		case BadgeWarning:
			b.bgFill = th.Color(theme.ColorNameWarning, v)
			b.textCol = th.Color(theme.ColorNameForegroundOnWarning, v)
		case BadgeInfo, BadgePurple:
			b.bgFill = th.Color(theme.ColorNamePrimary, v)
			b.textCol = th.Color(theme.ColorNameForegroundOnPrimary, v)
		default:
			b.bgFill = th.Color(theme.ColorNameDisabled, v)
			b.textCol = th.Color(theme.ColorNameForeground, v)
		}
		b.strokeCol = color.Transparent
		return
	}

	var colors BadgeVariantColors
	switch b.variant {
	case BadgeSuccess:
		colors = pal.Badges.Success
	case BadgeDanger:
		colors = pal.Badges.Danger
	case BadgeInfo:
		colors = pal.Badges.Info
	case BadgeWarning:
		colors = pal.Badges.Warning
	case BadgePurple:
		colors = pal.Badges.Purple
	case BadgeNeutral:
		fallthrough
	default:
		colors = pal.Badges.Neutral
	}

	b.bgFill = colors.Bg
	b.textCol = colors.Text
	b.strokeCol = colors.Border
}

// Implementación del Renderer
type modernBadgeRenderer struct {
	badge *ModernBadge
}

func (b *ModernBadge) CreateRenderer() fyne.WidgetRenderer {
	return &modernBadgeRenderer{badge: b}
}

func (r *modernBadgeRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.badge.bg, r.badge.label}
}

func (r *modernBadgeRenderer) Destroy() {}

func (r *modernBadgeRenderer) MinSize() fyne.Size {
	b := r.badge
	if b.text == "" && b.hideWhenEmpty {
		return fyne.NewSize(0, 0)
	}

	labelSize := fyne.MeasureText(b.text, b.label.TextSize, b.label.TextStyle)
	width := labelSize.Width + (b.paddingX * 2)
	if width < b.minWidth {
		width = b.minWidth
	}
	return fyne.NewSize(width, b.height)
}

func (r *modernBadgeRenderer) Layout(size fyne.Size) {
	b := r.badge
	if b.text == "" && b.hideWhenEmpty {
		b.bg.Resize(fyne.NewSize(0, 0))
		b.label.Resize(fyne.NewSize(0, 0))
		return
	}

	labelSize := fyne.MeasureText(b.text, b.label.TextSize, b.label.TextStyle)

	// Si no hay alineación explícita configurada, mantenemos el comportamiento clásico (rellenar size)
	if b.alignment == nil {
		b.bg.Resize(size)
		b.bg.Move(fyne.NewPos(0, 0))

		if b.isPill {
			b.bg.CornerRadius = size.Height / 2
		} else {
			b.bg.CornerRadius = b.cornerRadius
		}

		labelX := (size.Width - labelSize.Width) / 2
		labelY := (size.Height - labelSize.Height) / 2
		b.label.Resize(labelSize)
		b.label.Move(fyne.NewPos(labelX, labelY))
		return
	}

	// Comportamiento con alineación explícita (estilo fyne-kx)
	badgeWidth := labelSize.Width + (b.paddingX * 2)
	if badgeWidth < b.minWidth {
		badgeWidth = b.minWidth
	}
	if badgeWidth > size.Width && size.Width > 0 {
		badgeWidth = size.Width
	}

	badgeHeight := b.height
	if badgeHeight > size.Height && size.Height > 0 {
		badgeHeight = size.Height
	} else if badgeHeight <= 0 {
		badgeHeight = size.Height
	}

	var xPos float32
	switch *b.alignment {
	case fyne.TextAlignTrailing:
		xPos = size.Width - badgeWidth
	case fyne.TextAlignCenter:
		xPos = (size.Width - badgeWidth) / 2
	case fyne.TextAlignLeading:
		fallthrough
	default:
		xPos = 0
	}
	if xPos < 0 {
		xPos = 0
	}

	yPos := (size.Height - badgeHeight) / 2
	if yPos < 0 {
		yPos = 0
	}

	b.bg.Move(fyne.NewPos(xPos, yPos))
	b.bg.Resize(fyne.NewSize(badgeWidth, badgeHeight))

	if b.isPill {
		b.bg.CornerRadius = badgeHeight / 2
	} else {
		b.bg.CornerRadius = b.cornerRadius
	}

	labelX := xPos + (badgeWidth - labelSize.Width) / 2
	labelY := yPos + (badgeHeight - labelSize.Height) / 2
	b.label.Resize(labelSize)
	b.label.Move(fyne.NewPos(labelX, labelY))
}

func (r *modernBadgeRenderer) Refresh() {
	b := r.badge
	b.applyVariantColors()
	if b.text == "" && b.hideWhenEmpty {
		b.bg.FillColor = color.Transparent
		b.bg.StrokeColor = color.Transparent
		b.bg.StrokeWidth = 0
		b.label.Text = ""
		b.label.Color = color.Transparent
		b.bg.Hide()
		b.label.Hide()
	} else {
		b.bg.FillColor = b.bgFill
		b.bg.StrokeColor = b.strokeCol
		b.bg.StrokeWidth = b.strokeWidth
		b.label.Text = b.text
		b.label.Color = b.textCol
		b.bg.Show()
		b.label.Show()
	}

	b.label.Refresh()
	b.bg.Refresh()
	r.Layout(b.Size())
	canvas.Refresh(b)
}
