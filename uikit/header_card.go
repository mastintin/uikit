package uikit

import (
	"image"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ModernHeaderCard implementa una tarjeta ejecutiva seccionada en dos partes ("partida en dos"),
// inspirada en los paneles ejecutivos y cabeceras de reporte de xlsJubila:
// 1. Parte Superior (Banner de Cabecera): Fondo con degradado corporativo verde esmeralda
//    con esquinas superiores redondeadas, sin juntas ni solapes, con título, subtítulo e icono.
// 2. Parte Inferior (Cuerpo): Superficie estándar de tarjeta con fondo limpio para controles
//    de captura de datos, botones y formularios.
type ModernHeaderCard struct {
	widget.BaseWidget

	header fyne.CanvasObject
	body   fyne.CanvasObject

	// Objetos gráficos internos del renderizado
	bgOuter         *canvas.Rectangle
	bgBanner        *canvas.Raster
	divLine         *canvas.Rectangle
	outerBorder     *canvas.Rectangle

	// Capas Glassmorphism y Acabado Ejecutivo
	hasGlass        bool
	topRimHighlight *canvas.Rectangle
	iconGlassPillBg *canvas.Rectangle

	// Configuración de estilo
	cornerRadius      float32
	bannerH           float32
	bannerBgColor     color.Color
	bannerTopColor    color.NRGBA
	bannerBottomColor color.NRGBA
	hasCustomBanner   bool
	bodyBgColor       color.Color
	hasCustomBody     bool
	borderColor       color.Color
	hasCustomBorder   bool

	bannerPadTop    float32
	bannerPadBottom float32
	bannerPadLeft   float32
	bannerPadRight  float32

	bodyPadTop    float32
	bodyPadBottom float32
	bodyPadLeft   float32
	bodyPadRight  float32

	// Soporte de sombras nativas Fyne 2.8 (canvas.Shadow)
	hasShadow       bool
	shadowBlur      float32
	shadowSpread    float32
	shadowOffset    fyne.Position
	shadowColor     color.Color
	hasCustomShadow bool

	// Campos de conveniencia para títulos y textos dinámicos
	TitleText    *canvas.Text
	SubtitleText *canvas.Text
	IconWidget   *ModernIcon
}

// Alias semánticos
type HeaderCard = ModernHeaderCard
type BannerCard = ModernHeaderCard
type SplitCard = ModernHeaderCard

type modernHeaderCardRenderer struct {
	card *ModernHeaderCard
}

func (c *ModernHeaderCard) CreateRenderer() fyne.WidgetRenderer {
	return &modernHeaderCardRenderer{card: c}
}

func (r *modernHeaderCardRenderer) Destroy() {}

func (r *modernHeaderCardRenderer) Objects() []fyne.CanvasObject {
	objs := []fyne.CanvasObject{
		r.card.bgOuter,
		r.card.bgBanner,
	}
	if r.card.hasGlass && r.card.topRimHighlight != nil {
		objs = append(objs, r.card.topRimHighlight)
	}
	objs = append(objs, r.card.divLine)
	if r.card.header != nil {
		objs = append(objs, r.card.header)
	}
	if r.card.body != nil {
		objs = append(objs, r.card.body)
	}
	objs = append(objs, r.card.outerBorder)
	return objs
}

func (r *modernHeaderCardRenderer) MinSize() fyne.Size {
	var headerW, headerH float32
	if r.card.header != nil && r.card.header.Visible() {
		ms := r.card.header.MinSize()
		headerW = ms.Width
		headerH = ms.Height
	}

	var bodyW, bodyH float32
	if r.card.body != nil && r.card.body.Visible() {
		ms := r.card.body.MinSize()
		bodyW = ms.Width
		bodyH = ms.Height
	}

	bannerW := headerW + r.card.bannerPadLeft + r.card.bannerPadRight
	bannerH := headerH + r.card.bannerPadTop + r.card.bannerPadBottom
	if bannerH < 36 {
		bannerH = 36
	}

	bodyTotalW := bodyW + r.card.bodyPadLeft + r.card.bodyPadRight
	bodyTotalH := bodyH + r.card.bodyPadTop + r.card.bodyPadBottom

	totalW := bannerW
	if bodyTotalW > totalW {
		totalW = bodyTotalW
	}

	totalH := bannerH + bodyTotalH + 1 // +1px por la línea divisoria
	return fyne.NewSize(totalW, totalH)
}

func (r *modernHeaderCardRenderer) Layout(size fyne.Size) {
	c := r.card

	var headerH float32
	if c.header != nil && c.header.Visible() {
		headerH = c.header.MinSize().Height
	}
	bannerH := headerH + c.bannerPadTop + c.bannerPadBottom
	if bannerH < 36 {
		bannerH = 36
	}
	c.bannerH = bannerH

	// 1. Capa base exterior y borde perimetral continuo
	c.bgOuter.Resize(size)
	c.bgOuter.Move(fyne.NewPos(0, 0))

	c.outerBorder.Resize(size)
	c.outerBorder.Move(fyne.NewPos(0, 0))

	// 2. Banner superior con degradado y esquinas superiores redondeadas
	c.bgBanner.Resize(fyne.NewSize(size.Width, bannerH))
	c.bgBanner.Move(fyne.NewPos(0, 0))

	if c.hasGlass && c.topRimHighlight != nil {
		c.topRimHighlight.Resize(fyne.NewSize(size.Width, 1.5))
		c.topRimHighlight.Move(fyne.NewPos(0, 0))
	}

	// 3. Línea divisoria nítida horizontal
	c.divLine.Resize(fyne.NewSize(size.Width, 1))
	c.divLine.Move(fyne.NewPos(0, bannerH))

	// 4. Posicionamiento del contenido de la cabecera
	if c.header != nil && c.header.Visible() {
		hw := size.Width - c.bannerPadLeft - c.bannerPadRight
		hh := bannerH - c.bannerPadTop - c.bannerPadBottom
		if hw < 0 {
			hw = 0
		}
		if hh < 0 {
			hh = 0
		}
		c.header.Resize(fyne.NewSize(hw, hh))
		c.header.Move(fyne.NewPos(c.bannerPadLeft, c.bannerPadTop))
	}

	// 5. Posicionamiento del cuerpo inferior
	if c.body != nil && c.body.Visible() {
		bodyY := bannerH + 1 + c.bodyPadTop
		bw := size.Width - c.bodyPadLeft - c.bodyPadRight
		bh := size.Height - (bannerH + 1) - c.bodyPadTop - c.bodyPadBottom
		if bw < 0 {
			bw = 0
		}
		if bh < 0 {
			bh = 0
		}
		c.body.Resize(fyne.NewSize(bw, bh))
		c.body.Move(fyne.NewPos(c.bodyPadLeft, bodyY))
	}
}

func (r *modernHeaderCardRenderer) Refresh() {
	r.card.applyStyle()
	canvas.Refresh(r.card.bgOuter)
	canvas.Refresh(r.card.bgBanner)
	if r.card.hasGlass && r.card.topRimHighlight != nil {
		canvas.Refresh(r.card.topRimHighlight)
	}
	canvas.Refresh(r.card.divLine)
	canvas.Refresh(r.card.outerBorder)
	if r.card.header != nil {
		canvas.Refresh(r.card.header)
	}
	if r.card.body != nil {
		canvas.Refresh(r.card.body)
	}
}

// drawBannerRaster genera el degradado vertical continuo del banner con esquinas superiores redondeadas
// y borde inferior completamente plano, eliminando cualquier costura, solape o banda extraña.
func (c *ModernHeaderCard) drawBannerRaster(w, h int) image.Image {
	if w <= 0 || h <= 0 {
		return image.NewNRGBA(image.Rect(0, 0, 1, 1))
	}

	img := image.NewNRGBA(image.Rect(0, 0, w, h))

	topCol := c.bannerTopColor
	botCol := c.bannerBottomColor

	tr, tg, tb, ta := float64(topCol.R), float64(topCol.G), float64(topCol.B), float64(topCol.A)
	br, bg, bb, ba := float64(botCol.R), float64(botCol.G), float64(botCol.B), float64(botCol.A)

	rf := float64(c.cornerRadius)
	if h > 0 && c.bannerH > 0 {
		rf *= float64(h) / float64(c.bannerH)
	}

	for y := 0; y < h; y++ {
		t := float64(0)
		if h > 1 {
			t = float64(y) / float64(h-1)
		}
		curR := uint8(tr + t*(br-tr) + 0.5)
		curG := uint8(tg + t*(bg-tg) + 0.5)
		curB := uint8(tb + t*(bb-tb) + 0.5)
		baseA := ta + t*(ba-ta)

		for x := 0; x < w; x++ {
			alphaFactor := 1.0

			// Esquina superior izquierda
			if float64(x) < rf && float64(y) < rf {
				dx := rf - float64(x) - 0.5
				dy := rf - float64(y) - 0.5
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist > rf {
					alphaFactor = 0.0
				} else if dist > rf-1.0 {
					alphaFactor = rf - dist
				}
			} else if float64(x) >= float64(w)-rf && float64(y) < rf { // Esquina superior derecha
				dx := float64(x) - (float64(w) - rf) + 0.5
				dy := rf - float64(y) - 0.5
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist > rf {
					alphaFactor = 0.0
				} else if dist > rf-1.0 {
					alphaFactor = rf - dist
				}
			}

			finalA := uint8(baseA * alphaFactor + 0.5)
			if finalA == 0 {
				img.SetNRGBA(x, y, color.NRGBA{0, 0, 0, 0})
			} else {
				img.SetNRGBA(x, y, color.NRGBA{R: curR, G: curG, B: curB, A: finalA})
			}
		}
	}
	return img
}

// NewHeaderCard crea una tarjeta partida en dos con cabecera personalizada y cuerpo de contenido.
func NewHeaderCard(header, body fyne.CanvasObject) *ModernHeaderCard {
	c := &ModernHeaderCard{
		header:          header,
		body:            body,
		cornerRadius:    10,
		bannerPadTop:    14,
		bannerPadBottom: 14,
		bannerPadLeft:   18,
		bannerPadRight:  18,
		bodyPadTop:      14,
		bodyPadBottom:   14,
		bodyPadLeft:     18,
		bodyPadRight:    18,
		hasShadow:       true,
		shadowBlur:      6,
		shadowOffset:    fyne.NewPos(0, 2),
		hasGlass:        true,
	}

	c.bgOuter = canvas.NewRectangle(color.Transparent)
	c.bgBanner = canvas.NewRaster(func(w, h int) image.Image {
		return c.drawBannerRaster(w, h)
	})
	c.topRimHighlight = canvas.NewRectangle(color.Transparent)
	c.divLine = canvas.NewRectangle(color.Transparent)
	c.outerBorder = canvas.NewRectangle(color.Transparent)

	c.applyStyle()
	c.ExtendBaseWidget(c)
	return c
}

// NewBannerCard es un alias idéntico de NewHeaderCard.
func NewBannerCard(header, body fyne.CanvasObject) *ModernHeaderCard {
	return NewHeaderCard(header, body)
}

// NewHeaderCardWithTitle construye directamente la tarjeta partida en dos con título, subtítulo e icono
// en la parte superior, y el contenedor de controles en la parte inferior.
func NewHeaderCardWithTitle(
	iconName string,
	title string,
	subtitle string,
	body fyne.CanvasObject,
) *ModernHeaderCard {
	c := NewHeaderCard(nil, body)

	c.TitleText = canvas.NewText(title, color.White)
	c.TitleText.TextStyle = fyne.TextStyle{Bold: true}
	c.TitleText.TextSize = 18

	c.SubtitleText = canvas.NewText(subtitle, ThemeColor(ColorNameTextMuted))
	c.SubtitleText.TextSize = 12.5

	titleBox := container.NewVBox(
		c.TitleText,
		c.SubtitleText,
	)

	var headerContent fyne.CanvasObject
	if iconName != "" {
		c.IconWidget = NewColoredIcon(iconName, 22, "#ffffff")

		pillBg := canvas.NewRectangle(color.NRGBA{R: 255, G: 255, B: 255, A: 32})
		pillBg.CornerRadius = 10
		pillBg.StrokeColor = color.NRGBA{R: 255, G: 255, B: 255, A: 65}
		pillBg.StrokeWidth = 1.0

		c.iconGlassPillBg = pillBg

		pillBlur := canvas.NewBlur(8)
		pillBlur.CornerRadius = 10

		iconPill := container.NewStack(
			container.NewGridWrap(fyne.NewSize(42, 42), pillBlur),
			container.NewGridWrap(fyne.NewSize(42, 42), pillBg),
			container.NewCenter(c.IconWidget),
		)

		headerContent = container.NewHBox(
			container.NewCenter(iconPill),
			headerCardSpacer(14, 1),
			container.NewCenter(titleBox),
		)
	} else {
		headerContent = container.NewCenter(titleBox)
	}

	c.header = headerContent
	c.applyStyle()
	return c
}

func (c *ModernHeaderCard) applyStyle() {
	pal := GetCurrentThemePalette()

	c.bgOuter.CornerRadius = c.cornerRadius
	c.outerBorder.CornerRadius = c.cornerRadius
	c.outerBorder.StrokeWidth = 1

	// 1. Color de fondo del Banner Superior con degradado
	if !c.hasCustomBanner {
		if pal != nil && pal.HeaderCard.BannerTop.A > 0 {
			c.bannerTopColor = pal.HeaderCard.BannerTop
			c.bannerBottomColor = pal.HeaderCard.BannerBottom
			c.divLine.FillColor = pal.HeaderCard.BannerDivider
		} else {
			topCol := ColorToNRGBA(ThemeColor(ColorNameNavbarBg))
			botCol := ColorToNRGBA(ThemeColor(ColorNameAccent))
			c.bannerTopColor = topCol
			c.bannerBottomColor = botCol
			c.divLine.FillColor = ThemeColor(ColorNameBorder)
		}
	} else {
		c.divLine.FillColor = c.bannerBottomColor
	}

	// 1b. Acabado especular sutil en el borde superior
	if c.hasGlass && c.topRimHighlight != nil {
		c.topRimHighlight.CornerRadius = c.cornerRadius
		c.topRimHighlight.FillColor = color.NRGBA{R: 255, G: 255, B: 255, A: 45}
	} else if c.topRimHighlight != nil {
		c.topRimHighlight.FillColor = color.Transparent
	}

	if c.iconGlassPillBg != nil {
		if c.hasGlass {
			c.iconGlassPillBg.FillColor = color.NRGBA{R: 255, G: 255, B: 255, A: 32}
			c.iconGlassPillBg.StrokeColor = color.NRGBA{R: 255, G: 255, B: 255, A: 65}
			c.iconGlassPillBg.StrokeWidth = 1.0
		} else {
			c.iconGlassPillBg.FillColor = color.Transparent
			c.iconGlassPillBg.StrokeColor = color.Transparent
			c.iconGlassPillBg.StrokeWidth = 0
		}
	}

	// 2. Color de fondo del Cuerpo Inferior
	if c.hasCustomBody {
		c.bgOuter.FillColor = c.bodyBgColor
	} else if pal != nil && pal.HeaderCard.BodyBg != nil {
		c.bgOuter.FillColor = pal.HeaderCard.BodyBg
	} else {
		c.bgOuter.FillColor = ThemeColor(ColorNameCardBg)
	}

	// 3. Borde perimetral exterior
	if c.hasCustomBorder {
		c.outerBorder.StrokeColor = c.borderColor
	} else if pal != nil && pal.HeaderCard.Border != nil {
		c.outerBorder.StrokeColor = pal.HeaderCard.Border
	} else {
		c.outerBorder.StrokeColor = ThemeColor(ColorNameCardBorder)
	}

	// 4. Adaptación de textos del banner si se generaron con NewHeaderCardWithTitle
	if c.TitleText != nil {
		if pal != nil && pal.HeaderCard.TitleColor != nil {
			c.TitleText.Color = pal.HeaderCard.TitleColor
		} else {
			c.TitleText.Color = color.White
		}
		c.TitleText.Refresh()
	}
	if c.SubtitleText != nil {
		if pal != nil && pal.HeaderCard.SubtitleColor != nil {
			c.SubtitleText.Color = pal.HeaderCard.SubtitleColor
		} else {
			c.SubtitleText.Color = ThemeColor(ColorNameTextMuted)
		}
		c.SubtitleText.Refresh()
	}
	if c.IconWidget != nil {
		c.IconWidget.SetColor(color.White)
		c.IconWidget.Refresh()
	}
	if c.iconGlassPillBg != nil {
		c.iconGlassPillBg.Refresh()
	}

	// 5. Sombra acelerada por GPU (Fyne 2.8 canvas.Shadow)
	if c.hasShadow {
		shadowCol := c.shadowColor
		if !c.hasCustomShadow || shadowCol == nil {
			isDark := pal != nil && pal.IsDark
			if isDark {
				shadowCol = color.NRGBA{R: 0, G: 0, B: 0, A: 65}
			} else {
				shadowCol = color.NRGBA{R: 0, G: 0, B: 0, A: 25}
			}
		}
		c.bgOuter.Shadow = canvas.Shadow{
			Color:      shadowCol,
			BlurRadius: c.shadowBlur,
			Spread:     c.shadowSpread,
			Offset:     c.shadowOffset,
			Variant:    canvas.DropShadow,
		}
	} else {
		c.bgOuter.Shadow = canvas.Shadow{}
	}
}

// Fluent API

// WithShadow activa una sombra personalizada sobre la tarjeta usando el motor nativo de Fyne 2.8.
func (c *ModernHeaderCard) WithShadow(blurRadius, spread float32, offset fyne.Position, col ...color.Color) *ModernHeaderCard {
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
	c.applyStyle()
	c.Refresh()
	return c
}

// WithElevation aplica un nivel estándar de sombra y elevación (1: sutil/tarjeta, 2: flotante/panel, 3: pronunciada/modal).
func (c *ModernHeaderCard) WithElevation(level int) *ModernHeaderCard {
	if level <= 0 {
		c.hasShadow = false
		c.applyStyle()
		c.Refresh()
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
func (c *ModernHeaderCard) WithoutShadow() *ModernHeaderCard {
	c.hasShadow = false
	c.hasCustomShadow = false
	c.shadowColor = nil
	c.applyStyle()
	c.Refresh()
	return c
}

// WithBannerGradient define un degradado personalizado (color superior e inferior) para la cabecera.
func (c *ModernHeaderCard) WithBannerGradient(top, bottom color.Color) *ModernHeaderCard {
	c.hasCustomBanner = true
	c.bannerTopColor = toNRGBA(top)
	c.bannerBottomColor = toNRGBA(bottom)
	c.bannerBgColor = bottom
	c.applyStyle()
	c.Refresh()
	return c
}

// WithBannerColor asigna un color corporativo base, derivando automáticamente un sutil degradado.
func (c *ModernHeaderCard) WithBannerColor(bg color.Color) *ModernHeaderCard {
	c.hasCustomBanner = true
	c.bannerBgColor = bg
	base := toNRGBA(bg)
	c.bannerTopColor = adjustBrightness(base, 1.06)
	c.bannerBottomColor = adjustBrightness(base, 0.95)
	c.applyStyle()
	c.Refresh()
	return c
}

func (c *ModernHeaderCard) WithBannerHex(hex string) *ModernHeaderCard {
	return c.WithBannerColor(HexColor(hex))
}

func (c *ModernHeaderCard) WithBodyColor(bg color.Color) *ModernHeaderCard {
	c.hasCustomBody = true
	c.bodyBgColor = bg
	c.applyStyle()
	c.Refresh()
	return c
}

func (c *ModernHeaderCard) WithCornerRadius(radius float32) *ModernHeaderCard {
	c.cornerRadius = radius
	c.applyStyle()
	c.Refresh()
	return c
}

func (c *ModernHeaderCard) WithBannerPadding(top, bottom, left, right float32) *ModernHeaderCard {
	c.bannerPadTop = top
	c.bannerPadBottom = bottom
	c.bannerPadLeft = left
	c.bannerPadRight = right
	c.Refresh()
	return c
}

func (c *ModernHeaderCard) WithBodyPadding(top, bottom, left, right float32) *ModernHeaderCard {
	c.bodyPadTop = top
	c.bodyPadBottom = bottom
	c.bodyPadLeft = left
	c.bodyPadRight = right
	c.Refresh()
	return c
}

func (c *ModernHeaderCard) SetHeader(header fyne.CanvasObject) *ModernHeaderCard {
	c.header = header
	c.Refresh()
	return c
}

func (c *ModernHeaderCard) SetBody(body fyne.CanvasObject) *ModernHeaderCard {
	c.body = body
	c.Refresh()
	return c
}

// WithGlass activa o desactiva los efectos de glassmorphism sobre la tarjeta (por defecto true).
func (c *ModernHeaderCard) WithGlass(active bool) *ModernHeaderCard {
	c.hasGlass = active
	c.applyStyle()
	c.Refresh()
	return c
}

// IsGlass indica si el efecto de cristal está activo.
func (c *ModernHeaderCard) IsGlass() bool {
	return c.hasGlass
}

func headerCardSpacer(w, h float32) fyne.CanvasObject {
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(w, h))
	return rect
}

func toNRGBA(col color.Color) color.NRGBA {
	return ColorToNRGBA(col)
}

func adjustBrightness(c color.NRGBA, factor float64) color.NRGBA {
	if factor >= 1.0 {
		return ColorToNRGBA(ColorBrighter(c, float32(factor-1.0)))
	}
	return ColorToNRGBA(ColorDarker(c, float32(1.0-factor)))
}
