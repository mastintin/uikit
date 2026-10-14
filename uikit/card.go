package uikit

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ModernCard implementa un panel o tarjeta flotante estilizada tipo goCuadros
// (.detail-card, .image-card, .detail-top-bar), compartiendo la misma filosofía
// de ModernButton, ModernBadge y ModernTextBox.
//
// Extiende widget.BaseWidget para integrarse como widget nativo en Fyne.
type ModernCard struct {
	widget.BaseWidget

	bg      *canvas.Rectangle
	content fyne.CanvasObject

	cornerRadius  float32
	strokeWidth   float32
	strokeColor   color.Color
	bgColor       color.Color
	paddingLayers int
	hasCustomBg     bool
	hasCustomBorder bool
	isInset         bool

	// Paddings configurables independientes (superior, inferior, lateral)
	paddingTop    float32
	paddingBottom float32
	paddingLeft   float32
	paddingRight  float32

	hasCustomPaddingTop    bool
	hasCustomPaddingBottom bool
	hasCustomPaddingLeft   bool
	hasCustomPaddingRight  bool

	// Soporte de sombras nativas Fyne 2.8 (canvas.Shadow)
	hasShadow       bool
	shadowBlur      float32
	shadowSpread    float32
	shadowOffset    fyne.Position
	shadowColor     color.Color
	hasCustomShadow bool
}

// Alias de tipo para compatibilidad
type CardPanel = ModernCard

type modernCardRenderer struct {
	card *ModernCard
}

func (c *ModernCard) CreateRenderer() fyne.WidgetRenderer {
	return &modernCardRenderer{card: c}
}

func (r *modernCardRenderer) Objects() []fyne.CanvasObject {
	if r.card.content != nil {
		return []fyne.CanvasObject{r.card.bg, r.card.content}
	}
	return []fyne.CanvasObject{r.card.bg}
}

func (r *modernCardRenderer) Destroy() {}

func (r *modernCardRenderer) MinSize() fyne.Size {
	top, right, bottom, left := r.card.effectivePadding()
	var w, h float32
	if r.card.content != nil && r.card.content.Visible() {
		ms := r.card.content.MinSize()
		w = ms.Width
		h = ms.Height
	}
	return fyne.NewSize(w+left+right, h+top+bottom)
}

func (r *modernCardRenderer) Layout(size fyne.Size) {
	r.card.bg.Resize(size)
	r.card.bg.Move(fyne.NewPos(0, 0))

	if r.card.content == nil || !r.card.content.Visible() {
		return
	}

	top, right, bottom, left := r.card.effectivePadding()
	contentW := size.Width - left - right
	contentH := size.Height - top - bottom
	if contentW < 0 {
		contentW = 0
	}
	if contentH < 0 {
		contentH = 0
	}

	r.card.content.Resize(fyne.NewSize(contentW, contentH))
	r.card.content.Move(fyne.NewPos(left, top))
}

func (r *modernCardRenderer) Refresh() {
	r.card.applyStyle()
	canvas.Refresh(r.card.bg)
	if r.card.content != nil {
		canvas.Refresh(r.card.content)
	}
}

// effectivePadding resuelve la combinación de padding según capas base, tema y sobrescrituras específicas.
func (c *ModernCard) effectivePadding() (top, right, bottom, left float32) {
	// 1. Base por defecto calculada según paddingLayers y el tema actual
	var basePad float32 = 0
	if c.paddingLayers > 0 {
		padUnit := float32(8.0)
		if c.paddingLayers == 1 {
			padUnit = 6.0
		} else if app := fyne.CurrentApp(); app != nil && app.Settings() != nil {
			if th := app.Settings().Theme(); th != nil {
				padUnit = th.Size(theme.SizeNamePadding)
				if padUnit < 8.0 {
					padUnit = 8.0
				}
			}
		}
		basePad = float32(c.paddingLayers) * padUnit
	}

	top = basePad
	right = basePad
	bottom = basePad
	left = basePad

	// 2. Si el tema actual define tamaños personalizados para tarjetas estándar (layers >= 2)
	if c.paddingLayers >= 2 {
		if app := fyne.CurrentApp(); app != nil && app.Settings() != nil {
			if th := app.Settings().Theme(); th != nil {
				if appTh, ok := th.(*AppTheme); ok && appTh.customSizes != nil {
					if tVal, ok := appTh.customSizes[SizeNameCardPaddingTop]; ok {
						top = tVal
					}
					if bVal, ok := appTh.customSizes[SizeNameCardPaddingBottom]; ok {
						bottom = bVal
					}
					if lVal, ok := appTh.customSizes[SizeNameCardPaddingLeft]; ok {
						left = lVal
					}
					if rVal, ok := appTh.customSizes[SizeNameCardPaddingRight]; ok {
						right = rVal
					}
				}
			}
		}
	}

	// 3. Sobrescribir con configuraciones explícitas de la propia tarjeta
	if c.hasCustomPaddingTop {
		top = c.paddingTop
	}
	if c.hasCustomPaddingBottom {
		bottom = c.paddingBottom
	}
	if c.hasCustomPaddingLeft {
		left = c.paddingLeft
	}
	if c.hasCustomPaddingRight {
		right = c.paddingRight
	}

	return top, right, bottom, left
}

// NewCard crea un panel estilo tarjeta blanca estándar (.detail-card / .image-card):
// - Fondo blanco (#ffffff)
// - Esquinas redondeadas (10px)
// - Borde perimetral fino de 1px (#e2e8f0)
// - Acolchado interno amplio (~12-16px)
func NewCard(content fyne.CanvasObject) *ModernCard {
	c := NewPanel(content, "", "", 10, 2)
	c.hasCustomBg = false
	c.hasCustomBorder = false
	c.WithElevation(1)
	c.applyStyle()
	return c
}

// NewTintedCard crea una tarjeta con colores de fondo y borde personalizados.
func NewTintedCard(content fyne.CanvasObject, bgHex, borderHex string) *ModernCard {
	c := NewPanel(content, bgHex, borderHex, 12, 2)
	return c
}

// NewPillCard crea una tarjeta estilo cápsula o píldora con acolchado compacto (1 capa)
// y bordes redondeados, ideal para barras de navegación o conmutadores de pestañas.
func NewPillCard(content fyne.CanvasObject, bgHex, borderHex string) *ModernCard {
	c := NewPanel(content, bgHex, borderHex, 10, 1)
	return c
}

// NewCompactCard crea una tarjeta con acolchado simple (~6-8px),
// ideal para barras superiores o de herramientas (.detail-top-bar).
func NewCompactCard(content fyne.CanvasObject) *ModernCard {
	c := NewPanel(content, "", "", 10, 1)
	c.hasCustomBg = false
	c.hasCustomBorder = false
	c.applyStyle()
	return c
}

// NewInsetCard crea una tarjeta con fondo sombreado interior semántico (ColorNameInsetPanelBg),
// ideal para agrupar controles como casillas de verificación dentro de una tarjeta principal.
func NewInsetCard(content fyne.CanvasObject) *ModernCard {
	c := NewPanel(content, "", "", 6, 1)
	c.isInset = true
	c.hasCustomBg = false
	c.hasCustomBorder = false
	c.applyStyle()
	return c
}

// NewPanel es el constructor unificado estilo Harbour / xBase para paneles y tarjetas.
// content: contenido a colocar dentro
// bgHex: color de fondo en formato hexadecimal (ej. "#ffffff")
// borderHex: color del borde perimetral (ej. "#e2e8f0")
// cornerRadius: radio de curvatura de las esquinas en píxeles (ej. 10)
// paddingLayers: capas de padding interno (1 = compacto, 2 = espacioso)
func NewPanel(content fyne.CanvasObject, bgHex, borderHex string, cornerRadius float32, paddingLayers int) *ModernCard {
	if content == nil {
		content = container.NewVBox()
	}
	card := &ModernCard{
		bg:            canvas.NewRectangle(color.Transparent),
		content:       content,
		cornerRadius:  cornerRadius,
		strokeWidth:   1,
		paddingLayers: paddingLayers,
	}

	if bgHex != "" {
		card.bgColor = HexColor(bgHex)
		card.hasCustomBg = true
	}
	if borderHex != "" {
		card.strokeColor = HexColor(borderHex)
		card.hasCustomBorder = true
	}

	if card.paddingLayers < 0 {
		card.paddingLayers = 0
	}

	card.applyStyle()
	card.ExtendBaseWidget(card)
	return card
}

// Constructores de compatibilidad hacia atrás
func NewCardPanel(content fyne.CanvasObject) *ModernCard {
	return NewCard(content)
}

func NewCompactCardPanel(content fyne.CanvasObject) *ModernCard {
	return NewCompactCard(content)
}

// Métodos de configuración y fluidez (Fluent API)

func (c *ModernCard) applyStyle() {
	c.bg.CornerRadius = c.cornerRadius
	c.bg.StrokeWidth = c.strokeWidth

	// Fondo: si no se ha fijado un fondo explícito, usar el semántico
	if c.hasCustomBg {
		c.bg.FillColor = c.bgColor
	} else if c.isInset {
		c.bg.FillColor = ThemeColor(ColorNameInsetPanelBg)
	} else {
		c.bg.FillColor = ThemeColor(ColorNameCardBg)
	}

	// Borde: si no se ha fijado un borde explícito, usar el semántico
	if c.hasCustomBorder {
		c.bg.StrokeColor = c.strokeColor
	} else if c.isInset {
		c.bg.StrokeColor = ThemeColor(ColorNameBorder)
	} else {
		c.bg.StrokeColor = ThemeColor(ColorNameCardBorder)
	}

	// Sombra acelerada por GPU (Fyne 2.8 canvas.Shadow)
	if c.hasShadow && !c.isInset {
		shadowCol := c.shadowColor
		if !c.hasCustomShadow || shadowCol == nil {
			if IsDarkTheme() {
				shadowCol = color.NRGBA{R: 0, G: 0, B: 0, A: 65}
			} else {
				shadowCol = color.NRGBA{R: 0, G: 0, B: 0, A: 22}
			}
		}
		c.bg.Shadow = canvas.Shadow{
			Color:      shadowCol,
			BlurRadius: c.shadowBlur,
			Spread:     c.shadowSpread,
			Offset:     c.shadowOffset,
			Variant:    canvas.DropShadow,
		}
	} else {
		c.bg.Shadow = canvas.Shadow{}
	}
}

// WithShadow activa una sombra personalizada sobre la tarjeta usando el motor nativo de Fyne 2.8.
func (c *ModernCard) WithShadow(blurRadius, spread float32, offset fyne.Position, col ...color.Color) *ModernCard {
	c.hasShadow = true
	c.shadowBlur = blurRadius
	c.shadowSpread = spread
	c.shadowOffset = offset
	if len(col) > 0 && col[0] != nil {
		c.shadowColor = col[0]
		c.hasCustomShadow = true
	} else {
		c.shadowColor = nil
		c.hasCustomShadow = false
	}
	c.rebuild()
	return c
}

// WithElevation aplica un nivel estándar de sombra y elevación (1: sutil/tarjeta, 2: flotante/panel, 3: pronunciada/modal).
func (c *ModernCard) WithElevation(level int) *ModernCard {
	if level <= 0 {
		c.hasShadow = false
		c.rebuild()
		return c
	}
	switch level {
	case 1:
		return c.WithShadow(6, 0, fyne.NewPos(0, 2))
	case 2:
		return c.WithShadow(10, 1, fyne.NewPos(0, 4))
	default:
		return c.WithShadow(16, 2, fyne.NewPos(0, 6))
	}
}

// WithoutShadow desactiva explícitamente cualquier sombra de la tarjeta.
func (c *ModernCard) WithoutShadow() *ModernCard {
	c.hasShadow = false
	c.hasCustomShadow = false
	c.shadowColor = nil
	c.rebuild()
	return c
}

// WithInset configura la tarjeta para usar el color de fondo interior (inset) semántico.
func (c *ModernCard) WithInset() *ModernCard {
	c.isInset = true
	c.hasCustomBg = false
	c.hasCustomBorder = false
	c.hasShadow = false
	c.rebuild()
	return c
}

// WithBorderColor asigna un color de borde explícito a la tarjeta sin alterar el fondo semántico.
func (c *ModernCard) WithBorderColor(col color.Color, width float32) *ModernCard {
	c.strokeColor = col
	c.strokeWidth = width
	c.hasCustomBorder = (col != nil && width > 0)
	c.rebuild()
	return c
}

func (c *ModernCard) rebuild() {
	c.applyStyle()
	c.Refresh()
}

func (c *ModernCard) SetContent(content fyne.CanvasObject) *ModernCard {
	if content == nil {
		content = container.NewVBox()
	}
	c.content = content
	c.rebuild()
	return c
}

func (c *ModernCard) Content() fyne.CanvasObject {
	return c.content
}

func (c *ModernCard) WithRadius(radius float32) *ModernCard {
	c.cornerRadius = radius
	c.rebuild()
	return c
}

func (c *ModernCard) WithBorder(borderHex string, width float32) *ModernCard {
	c.strokeColor = HexColor(borderHex)
	c.strokeWidth = width
	c.hasCustomBorder = (borderHex != "" && width > 0)
	c.rebuild()
	return c
}

func (c *ModernCard) WithBackground(bgHex string) *ModernCard {
	c.bgColor = HexColor(bgHex)
	c.hasCustomBg = (bgHex != "")
	c.rebuild()
	return c
}

func (c *ModernCard) WithBgColor(col color.Color) *ModernCard {
	c.bgColor = col
	c.hasCustomBg = (col != nil)
	c.rebuild()
	return c
}

func (c *ModernCard) WithPadding(layers int) *ModernCard {
	c.paddingLayers = layers
	c.hasCustomPaddingTop = false
	c.hasCustomPaddingBottom = false
	c.hasCustomPaddingLeft = false
	c.hasCustomPaddingRight = false
	c.rebuild()
	return c
}

// WithPaddingTop define un padding superior explícito para la tarjeta.
func (c *ModernCard) WithPaddingTop(pt float32) *ModernCard {
	c.paddingTop = pt
	c.hasCustomPaddingTop = true
	c.rebuild()
	return c
}

// WithPaddingBottom define un padding inferior explícito para la tarjeta.
func (c *ModernCard) WithPaddingBottom(pb float32) *ModernCard {
	c.paddingBottom = pb
	c.hasCustomPaddingBottom = true
	c.rebuild()
	return c
}

// WithVerticalPadding define paddings superior e inferior explícitos para la tarjeta.
func (c *ModernCard) WithVerticalPadding(top, bottom float32) *ModernCard {
	c.paddingTop = top
	c.hasCustomPaddingTop = true
	c.paddingBottom = bottom
	c.hasCustomPaddingBottom = true
	c.rebuild()
	return c
}

// WithPaddingV define el mismo padding vertical (superior e inferior) para la tarjeta.
func (c *ModernCard) WithPaddingV(v float32) *ModernCard {
	return c.WithVerticalPadding(v, v)
}

// WithHorizontalPadding define paddings izquierdo y derecho explícitos para la tarjeta.
func (c *ModernCard) WithHorizontalPadding(left, right float32) *ModernCard {
	c.paddingLeft = left
	c.hasCustomPaddingLeft = true
	c.paddingRight = right
	c.hasCustomPaddingRight = true
	c.rebuild()
	return c
}

// WithCustomPadding define los 4 márgenes internos de padding de forma independiente.
func (c *ModernCard) WithCustomPadding(top, right, bottom, left float32) *ModernCard {
	c.paddingTop = top
	c.hasCustomPaddingTop = true
	c.paddingRight = right
	c.hasCustomPaddingRight = true
	c.paddingBottom = bottom
	c.hasCustomPaddingBottom = true
	c.paddingLeft = left
	c.hasCustomPaddingLeft = true
	c.rebuild()
	return c
}

// SetPaddingTop define el padding superior explícito.
func (c *ModernCard) SetPaddingTop(pt float32) *ModernCard {
	return c.WithPaddingTop(pt)
}

// SetPaddingBottom define el padding inferior explícito.
func (c *ModernCard) SetPaddingBottom(pb float32) *ModernCard {
	return c.WithPaddingBottom(pb)
}

// SetVerticalPadding define los paddings superior e inferior explícitos.
func (c *ModernCard) SetVerticalPadding(top, bottom float32) *ModernCard {
	return c.WithVerticalPadding(top, bottom)
}

// SetCustomPadding define los 4 márgenes de padding independientes.
func (c *ModernCard) SetCustomPadding(top, right, bottom, left float32) *ModernCard {
	return c.WithCustomPadding(top, right, bottom, left)
}
