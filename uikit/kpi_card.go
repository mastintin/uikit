package uikit

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// ─────────────────────────────────────────────────────────────────────────────
// KPICardPreset — 4 Variantes estándar y generación al vuelo
// ─────────────────────────────────────────────────────────────────────────────

// KPICardVariant identifica de forma semántica el tipo de métrica a mostrar.
type KPICardVariant string

const (
	KPICardVariantEmerald KPICardVariant = "emerald" // 1. Pensión Neta / Favorable (Verde)
	KPICardVariantBlue    KPICardVariant = "blue"    // 2. Pensión Bruta / Totales (Azul)
	KPICardVariantViolet  KPICardVariant = "violet"  // 3. Base Reguladora / Cálculos (Violeta)
	KPICardVariantAmber   KPICardVariant = "amber"   // 4. % Coeficiente / IRPF / Advertencias (Ámbar)
)

// KPICardPreset encapsula la configuración de una ModernKPICard.
type KPICardPreset struct {
	Variant   KPICardVariant // Si es una de las 4 estándar, lee de la paleta activa
	AccentHex string         // Color principal (si no es estándar, el resto se auto-genera)

	// Anulaciones manuales opcionales
	ValueHex        string
	LightBgHex      string
	DarkBgHex       string
	BorderHex       string
	DarkBorderHex   string
	SubtitleHex     string
	DarkSubtitleHex string
}

// Alias de retrocompatibilidad
type NiceCardPreset = KPICardPreset

// Las 4 variantes estándar de la aplicación:
var (
	KPICardEmerald = KPICardPreset{Variant: KPICardVariantEmerald, AccentHex: "#10b981"}
	KPICardBlue    = KPICardPreset{Variant: KPICardVariantBlue, AccentHex: "#3b82f6"}
	KPICardViolet  = KPICardPreset{Variant: KPICardVariantViolet, AccentHex: "#818cf8"}
	KPICardAmber   = KPICardPreset{Variant: KPICardVariantAmber, AccentHex: "#f59e0b"}

	// Alias heredados para no romper llamadas existentes
	NiceCardEmerald = KPICardEmerald
	NiceCardBlue    = KPICardBlue
	NiceCardViolet  = KPICardViolet
	NiceCardAmber   = KPICardAmber
	NiceCardSlate   = KPICardPreset{AccentHex: "#475569"} // Auto-genera sobre gris
	NiceCardPurple  = KPICardPreset{AccentHex: "#c084fc"} // Auto-genera sobre púrpura
)

// ─────────────────────────────────────────────────────────────────────────────
// ModernKPICard (alias ModernNiceCard)
// ─────────────────────────────────────────────────────────────────────────────

const (
	kpiCardStripW float32 = 5   // anchura de la barra lateral de acento
	kpiCardPadX   float32 = 14  // padding horizontal interior
	kpiCardPadY   float32 = 11  // padding vertical interior
	kpiCardMinH   float32 = 88  // altura mínima
	kpiCardMinW   float32 = 160 // anchura mínima
	kpiCardRadius float32 = 8   // radio de esquinas
	kpiCardGapLV  float32 = 4   // separación label → valor
	kpiCardGapVS  float32 = 5   // separación valor → subtitle
)

// ModernKPICard es la tarjeta de métrica Hero de la aplicación.
type ModernKPICard struct {
	widget.BaseWidget

	preset  KPICardPreset
	labelTx string
	valueTx string
	subTx   string

	// Elementos de canvas
	bg       *canvas.Rectangle
	strip    *canvas.Rectangle
	border   *canvas.Rectangle
	lblLabel *canvas.Text
	lblValue *canvas.Text
	lblSub   *canvas.Text
}

// Alias de retrocompatibilidad
type ModernNiceCard = ModernKPICard

// NewKPICard crea una ModernKPICard con cualquiera de las 4 variantes estándar o personalizada.
func NewKPICard(preset KPICardPreset, label, value, subtitle string) *ModernKPICard {
	c := &ModernKPICard{
		preset:  preset,
		labelTx: strings.ToUpper(label),
		valueTx: value,
		subTx:   subtitle,
	}
	c.build()
	c.ExtendBaseWidget(c)
	return c
}

// NewNiceCard es el alias retrocompatible para NewKPICard.
func NewNiceCard(preset NiceCardPreset, label, value, subtitle string) *ModernKPICard {
	return NewKPICard(preset, label, value, subtitle)
}

// NewKPICardHex crea una ModernKPICard a partir de CUALQUIER color hexadecimal.
// Los tonos de fondo (claro/oscuro), borde y valor se generan automáticamente al vuelo.
func NewKPICardHex(accentHex, label, value, subtitle string) *ModernKPICard {
	preset := KPICardPreset{AccentHex: accentHex}
	return NewKPICard(preset, label, value, subtitle)
}

// NewNiceCardHex es el alias retrocompatible para NewKPICardHex.
func NewNiceCardHex(accentHex, label, value, subtitle string) *ModernKPICard {
	return NewKPICardHex(accentHex, label, value, subtitle)
}

func (c *ModernKPICard) build() {
	isDark := IsDarkTheme()

	// Fondo y sombra
	c.bg = canvas.NewRectangle(c.resolveBg(isDark))
	c.bg.CornerRadius = kpiCardRadius
	c.applyShadow(isDark)

	// Borde exterior
	c.border = canvas.NewRectangle(color.Transparent)
	c.border.CornerRadius = kpiCardRadius
	c.border.StrokeWidth = 1
	c.border.StrokeColor = c.resolveBorder(isDark)

	// Barra lateral de acento
	c.strip = canvas.NewRectangle(c.resolveAccent())
	c.strip.CornerRadius = 0

	// Textos
	c.lblLabel = canvas.NewText(c.labelTx, c.resolveAccent())
	c.lblLabel.TextSize = 10.5
	c.lblLabel.TextStyle = fyne.TextStyle{Bold: true}

	c.lblValue = canvas.NewText(c.valueTx, c.resolveValue(isDark))
	c.lblValue.TextSize = 24
	c.lblValue.TextStyle = fyne.TextStyle{Bold: true}

	c.lblSub = canvas.NewText(c.subTx, c.resolveSubtitle(isDark))
	c.lblSub.TextSize = 11
}

func (c *ModernKPICard) CreateRenderer() fyne.WidgetRenderer {
	return &kpiCardRenderer{card: c}
}

// ── API fluida ───────────────────────────────────────────────────────────────

func (c *ModernKPICard) SetLabel(text string) *ModernKPICard {
	c.labelTx = strings.ToUpper(text)
	c.lblLabel.Text = c.labelTx
	c.lblLabel.Refresh()
	return c
}

func (c *ModernKPICard) SetValue(text string) *ModernKPICard {
	c.valueTx = text
	c.lblValue.Text = text
	c.lblValue.Refresh()
	return c
}

func (c *ModernKPICard) SetSubtitle(text string) *ModernKPICard {
	c.subTx = text
	c.lblSub.Text = text
	c.lblSub.Refresh()
	return c
}

func (c *ModernKPICard) SetPreset(preset KPICardPreset) *ModernKPICard {
	c.preset = preset
	c.Refresh()
	return c
}

func (c *ModernKPICard) SetValueSize(size float32) *ModernKPICard {
	c.lblValue.TextSize = size
	c.Refresh()
	return c
}

func (c *ModernKPICard) applyThemeStyle() {
	isDark := IsDarkTheme()

	accentCol := c.resolveAccent()
	c.bg.FillColor = c.resolveBg(isDark)
	c.applyShadow(isDark)
	c.bg.Refresh()

	c.border.FillColor = color.Transparent
	c.border.StrokeColor = c.resolveBorder(isDark)
	c.border.StrokeWidth = 1
	c.border.CornerRadius = kpiCardRadius
	c.border.Refresh()

	c.strip.FillColor = accentCol
	c.strip.Refresh()

	c.lblLabel.Color = accentCol
	c.lblLabel.Text = c.labelTx
	c.lblLabel.Refresh()

	c.lblValue.Color = c.resolveValue(isDark)
	c.lblValue.Text = c.valueTx
	c.lblValue.Refresh()

	c.lblSub.Color = c.resolveSubtitle(isDark)
	c.lblSub.Text = c.subTx
	c.lblSub.Refresh()
}

func (c *ModernKPICard) Refresh() {
	c.applyThemeStyle()
	c.BaseWidget.Refresh()
}

func (c *ModernKPICard) applyShadow(isDark bool) {
	var shadowCol color.Color
	if isDark {
		shadowCol = color.NRGBA{R: 0, G: 0, B: 0, A: 65}
	} else {
		shadowCol = color.NRGBA{R: 0, G: 0, B: 0, A: 22}
	}
	c.bg.Shadow = canvas.Shadow{
		Color:      shadowCol,
		BlurRadius: 6,
		Spread:     0,
		Offset:     fyne.NewPos(0, 2),
		Variant:    canvas.DropShadow,
	}
}

// ── Resolución de colores: de la paleta activa o generados al vuelo ──────────

func (c *ModernKPICard) getPaletteColors() *KPICardVariantColors {
	pal := GetCurrentThemePalette()
	if pal == nil {
		return nil
	}
	switch c.preset.Variant {
	case KPICardVariantEmerald:
		return &pal.KPICards.Emerald
	case KPICardVariantBlue:
		return &pal.KPICards.Blue
	case KPICardVariantViolet:
		return &pal.KPICards.Violet
	case KPICardVariantAmber:
		return &pal.KPICards.Amber
	default:
		return nil
	}
}

func (c *ModernKPICard) resolveAccent() color.Color {
	if pCols := c.getPaletteColors(); pCols != nil && pCols.Accent != nil {
		return pCols.Accent
	}
	if c.preset.AccentHex != "" {
		return HexColor(c.preset.AccentHex)
	}
	return ThemeColorAccent()
}

func (c *ModernKPICard) resolveBg(dark bool) color.Color {
	if pCols := c.getPaletteColors(); pCols != nil && pCols.Bg != nil {
		return pCols.Bg
	}
	if dark {
		if c.preset.DarkBgHex != "" {
			return HexColor(c.preset.DarkBgHex)
		}
		// Generación al vuelo en oscuro: mezcla sutil del acento sobre fondo oscuro
		if c.preset.AccentHex != "" {
			return blendColor(HexColor(c.preset.AccentHex), color.NRGBA{R: 15, G: 23, B: 42, A: 255}, 0.20)
		}
		return ThemeColor(ColorNameCardBg)
	}
	if c.preset.LightBgHex != "" {
		return HexColor(c.preset.LightBgHex)
	}
	// Generación al vuelo en claro: tono pastel suave mezclando el acento con blanco
	if c.preset.AccentHex != "" {
		return blendColor(HexColor(c.preset.AccentHex), color.NRGBA{R: 255, G: 255, B: 255, A: 255}, 0.08)
	}
	return ThemeColor(ColorNameCardBg)
}

func (c *ModernKPICard) resolveBorder(dark bool) color.Color {
	if pCols := c.getPaletteColors(); pCols != nil && pCols.Border != nil {
		return pCols.Border
	}
	if dark {
		if c.preset.DarkBorderHex != "" {
			return HexColor(c.preset.DarkBorderHex)
		}
		if c.preset.AccentHex != "" {
			return blendColor(HexColor(c.preset.AccentHex), color.NRGBA{R: 30, G: 41, B: 59, A: 255}, 0.40)
		}
		return c.resolveAccent()
	}
	if c.preset.BorderHex != "" {
		return HexColor(c.preset.BorderHex)
	}
	if c.preset.AccentHex != "" {
		return blendColor(HexColor(c.preset.AccentHex), color.NRGBA{R: 255, G: 255, B: 255, A: 255}, 0.25)
	}
	return c.resolveAccent()
}

func (c *ModernKPICard) resolveValue(dark bool) color.Color {
	if pCols := c.getPaletteColors(); pCols != nil && pCols.Value != nil {
		return pCols.Value
	}
	if dark {
		return c.resolveAccent()
	}
	if c.preset.ValueHex != "" {
		return HexColor(c.preset.ValueHex)
	}
	// En modo claro, oscurece el acento un 25% para máxima legibilidad
	if c.preset.AccentHex != "" {
		return blendColor(HexColor(c.preset.AccentHex), color.NRGBA{R: 0, G: 0, B: 0, A: 255}, 0.70)
	}
	return c.resolveAccent()
}

func (c *ModernKPICard) resolveSubtitle(dark bool) color.Color {
	if pCols := c.getPaletteColors(); pCols != nil && pCols.Subtitle != nil {
		return pCols.Subtitle
	}
	if dark {
		if c.preset.DarkSubtitleHex != "" {
			return HexColor(c.preset.DarkSubtitleHex)
		}
		return ThemeColor(ColorNameTextMuted)
	}
	if c.preset.SubtitleHex != "" {
		return HexColor(c.preset.SubtitleHex)
	}
	return ThemeColor(ColorNameTextMuted)
}

// blendColor interpola suavemente entre dos colores con un factor t (0.0 = colB, 1.0 = colA)
func blendColor(colA, colB color.Color, factor float32) color.Color {
	return BlendColor(colB, colA, factor)
}

// ─────────────────────────────────────────────────────────────────────────────
// kpiCardRenderer
// ─────────────────────────────────────────────────────────────────────────────

type kpiCardRenderer struct {
	card *ModernKPICard
}

func (r *kpiCardRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		r.card.bg,
		r.card.border,
		r.card.strip,
		r.card.lblLabel,
		r.card.lblValue,
		r.card.lblSub,
	}
}

func (r *kpiCardRenderer) Destroy() {}

func (r *kpiCardRenderer) MinSize() fyne.Size {
	c := r.card
	labelW := c.lblLabel.MinSize().Width
	valueW := c.lblValue.MinSize().Width
	subW := c.lblSub.MinSize().Width

	contentW := labelW
	if valueW > contentW {
		contentW = valueW
	}
	if subW > contentW {
		contentW = subW
	}
	minW := kpiCardStripW + kpiCardPadX + contentW + kpiCardPadX
	if minW < kpiCardMinW {
		minW = kpiCardMinW
	}

	labelH := c.lblLabel.MinSize().Height
	valueH := c.lblValue.MinSize().Height
	subH := c.lblSub.MinSize().Height
	minH := kpiCardPadY + labelH + kpiCardGapLV + valueH + kpiCardGapVS + subH + kpiCardPadY
	if minH < kpiCardMinH {
		minH = kpiCardMinH
	}

	return fyne.NewSize(minW, minH)
}

func (r *kpiCardRenderer) Layout(size fyne.Size) {
	c := r.card

	// Fondo y borde
	c.bg.Resize(size)
	c.bg.Move(fyne.NewPos(0, 0))
	c.border.Resize(size)
	c.border.Move(fyne.NewPos(0, 0))

	// Barra lateral de acento (izquierda)
	c.strip.Resize(fyne.NewSize(kpiCardStripW, size.Height))
	c.strip.Move(fyne.NewPos(0, 0))

	// Área de contenido
	contentX := kpiCardStripW + kpiCardPadX
	contentW := size.Width - contentX - kpiCardPadX

	labelH := c.lblLabel.MinSize().Height
	valueH := c.lblValue.MinSize().Height
	subH := c.lblSub.MinSize().Height
	totalH := labelH + kpiCardGapLV + valueH + kpiCardGapVS + subH

	startY := (size.Height - totalH) / 2
	if startY < kpiCardPadY {
		startY = kpiCardPadY
	}

	// Label
	c.lblLabel.Resize(fyne.NewSize(contentW, labelH))
	c.lblLabel.Move(fyne.NewPos(contentX, startY))

	// Valor
	valueY := startY + labelH + kpiCardGapLV
	c.lblValue.Resize(fyne.NewSize(contentW, valueH))
	c.lblValue.Move(fyne.NewPos(contentX, valueY))

	// Subtítulo
	subY := valueY + valueH + kpiCardGapVS
	c.lblSub.Resize(fyne.NewSize(contentW, subH))
	c.lblSub.Move(fyne.NewPos(contentX, subY))
}

func (r *kpiCardRenderer) Refresh() {
	r.card.applyThemeStyle()
}
