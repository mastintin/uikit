package uikit

import (
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

var _ desktop.Hoverable = (*ModernLineChart)(nil)

// ChartAxisY define si una serie pertenece al eje Y izquierdo o derecho.
type ChartAxisY int

const (
	ChartAxisLeft ChartAxisY = iota
	ChartAxisRight
)

// ChartSeries define una curva de datos a graficar.
type ChartSeries struct {
	Name           string                   // Nombre para leyenda y tooltip
	Color          color.Color              // Color representativo
	Axis           ChartAxisY               // Eje Y asignado (Izquierdo o Derecho)
	IsDashed       bool                     // Línea discontinua o continua
	StrokeWidth    float32                  // Grosor del trazo
	Values         []float64                // Puntos Y ordenados respecto a los puntos X
	ValueFormatter func(val float64) string // Formateador opcional de los valores
}

// ChartDetailRow representa una fila de metadatos adicional para la ventana flotante (tooltip).
type ChartDetailRow struct {
	Label string
	Value string
	Color color.Color
	Bold  bool
}

// ChartPoint representa una posición en el eje X con información contextual.
type ChartPoint struct {
	XLabel   string           // Etiqueta mostrada en el eje X (ej: "24m", "10/26")
	Title    string           // Título en la tarjeta flotante (ej: "Anticipo: 24 meses")
	Subtitle string           // Subtítulo en la tarjeta flotante (ej: "Fecha: 26/08/2028")
	Details  []ChartDetailRow // Filas adicionales antes de las series en el tooltip
}

// ChartAxisConfig configura las propiedades de escala y títulos de un eje Y.
type ChartAxisConfig struct {
	Title     string
	Color     color.Color
	Min       *float64
	Max       *float64
	Formatter func(val float64) string
}

// ChartConfig contiene las opciones generales para construir el gráfico.
type ChartConfig struct {
	XAxisTitle string
	LeftAxis   ChartAxisConfig
	RightAxis  *ChartAxisConfig // nil si sólo hay eje izquierdo
	MinSize    fyne.Size
}

// ModernLineChart es un widget vectorial nativo de Fyne interactivo de alto rendimiento.
// Soporta uno o dos ejes Y, líneas continuas y discontinuas (dashed), leyendas adaptativas al tema,
// y hover con línea de referencia vertical y tarjeta flotante (popover card) con valores exactos.
type ModernLineChart struct {
	widget.BaseWidget

	config ChartConfig
	points []ChartPoint
	series []ChartSeries

	// Estado interactivo (hover del cursor)
	isHovered  bool
	hoverIndex int
	cursorPos  fyne.Position
}

// NewModernLineChart crea una nueva instancia de ModernLineChart.
func NewModernLineChart(config ChartConfig) *ModernLineChart {
	if config.MinSize.Width <= 0 || config.MinSize.Height <= 0 {
		config.MinSize = fyne.NewSize(500, 240)
	}

	c := &ModernLineChart{
		config:     config,
		hoverIndex: -1,
	}
	c.ExtendBaseWidget(c)
	return c
}

// SetData asigna los puntos y las series al gráfico y lo redibuja.
func (c *ModernLineChart) SetData(points []ChartPoint, series []ChartSeries) {
	c.points = points
	c.series = series
	c.hoverIndex = -1
	c.isHovered = false
	c.Refresh()
}

// Clear borra los datos del gráfico.
func (c *ModernLineChart) Clear() {
	c.points = nil
	c.series = nil
	c.hoverIndex = -1
	c.isHovered = false
	c.Refresh()
}

// MouseIn captura la entrada del cursor.
func (c *ModernLineChart) MouseIn(e *desktop.MouseEvent) {
	c.handleHover(e.Position)
}

// MouseMoved actualiza la posición del cursor, la línea guía y el tooltip.
func (c *ModernLineChart) MouseMoved(e *desktop.MouseEvent) {
	c.handleHover(e.Position)
}

// MouseOut oculta la línea guía y el tooltip cuando el cursor abandona el gráfico.
func (c *ModernLineChart) MouseOut() {
	if c.isHovered {
		c.isHovered = false
		c.hoverIndex = -1
		c.Refresh()
	}
}

// IsHovered indica si el cursor está sobre el gráfico.
func (c *ModernLineChart) IsHovered() bool {
	return c.isHovered
}

// HoverIndex devuelve el índice del punto activo por el cursor (-1 si inactivo).
func (c *ModernLineChart) HoverIndex() int {
	return c.hoverIndex
}

// getPaddings centraliza los márgenes para asegurar sincronía entre trazado y hover.
func (c *ModernLineChart) getPaddings() (padLeft, padRight, padTop, padBottom float32) {
	padLeft = float32(20)
	padRight = float32(24)
	if c.config.RightAxis != nil {
		padRight = 48
	}
	padTop = float32(46)
	padBottom = float32(38)
	return
}

func (c *ModernLineChart) handleHover(pos fyne.Position) {
	n := len(c.points)
	if n < 2 {
		return
	}

	size := c.Size()
	padLeft, padRight, padTop, padBottom := c.getPaddings()

	plotW := size.Width - padLeft - padRight
	plotH := size.Height - padTop - padBottom

	if plotW <= 0 || plotH <= 0 {
		return
	}

	// Verificar si el cursor se encuentra en la zona del gráfico con tolerancia
	if pos.X < padLeft-15 || pos.X > padLeft+plotW+15 || pos.Y < padTop-20 || pos.Y > padTop+plotH+30 {
		if c.isHovered {
			c.isHovered = false
			c.hoverIndex = -1
			c.Refresh()
		}
		return
	}

	stepX := plotW / float32(n-1)
	idx := int(math.Round(float64(pos.X-padLeft) / float64(stepX)))
	if idx < 0 {
		idx = 0
	}
	if idx >= n {
		idx = n - 1
	}

	if !c.isHovered || c.hoverIndex != idx || math.Abs(float64(c.cursorPos.Y-pos.Y)) > 12 {
		c.isHovered = true
		c.hoverIndex = idx
		c.cursorPos = pos
		c.Refresh()
	}
}

func (c *ModernLineChart) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(ThemeColor(ColorNameInsetPanelBg))
	bg.CornerRadius = 10
	bg.StrokeColor = ThemeColor(ColorNameBorder)
	bg.StrokeWidth = 1

	return &modernLineChartRenderer{
		chart:   c,
		bg:      bg,
		objects: []fyne.CanvasObject{bg},
	}
}

type modernLineChartRenderer struct {
	chart   *ModernLineChart
	bg      *canvas.Rectangle
	objects []fyne.CanvasObject
}

func (r *modernLineChartRenderer) Destroy() {}

func (r *modernLineChartRenderer) MinSize() fyne.Size {
	return r.chart.config.MinSize
}

func (r *modernLineChartRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *modernLineChartRenderer) Refresh() {
	r.bg.FillColor = ThemeColor(ColorNameInsetPanelBg)
	r.bg.StrokeColor = ThemeColor(ColorNameBorder)
	r.bg.StrokeWidth = 1
	r.bg.Refresh()
	r.Layout(r.chart.Size())
}

func (r *modernLineChartRenderer) Layout(size fyne.Size) {
	r.bg.FillColor = ThemeColor(ColorNameInsetPanelBg)
	r.bg.StrokeColor = ThemeColor(ColorNameBorder)
	r.bg.StrokeWidth = 1
	r.bg.Resize(size)
	r.bg.Refresh()

	points := r.chart.points
	series := r.chart.series
	n := len(points)

	objs := []fyne.CanvasObject{r.bg}

	if n < 2 || len(series) == 0 || size.Width < 200 || size.Height < 120 {
		emptyLbl := canvas.NewText("Sin datos de proyección temporal", ThemeColor(ColorNameTextMuted))
		emptyLbl.TextSize = 12
		emptyLbl.Move(fyne.NewPos(size.Width/2-100, size.Height/2-8))
		objs = append(objs, emptyLbl)
		r.objects = objs
		canvas.Refresh(r.chart)
		return
	}

	colTextMuted := ThemeColor(ColorNameTextMuted)
	colTextPrimary := ThemeColor(ColorNameTextPrimary)
	colBorder := ThemeColor(ColorNameBorder)

	// ── 1. Leyenda Superior ────────────────────────────────────────────────
	legY := float32(14)
	currX := float32(24)

	for _, s := range series {
		if s.IsDashed {
			line1 := canvas.NewLine(s.Color)
			line1.StrokeWidth = 2.5
			line1.Position1 = fyne.NewPos(currX, legY+6)
			line1.Position2 = fyne.NewPos(currX+6, legY+6)

			line2 := canvas.NewLine(s.Color)
			line2.StrokeWidth = 2.5
			line2.Position1 = fyne.NewPos(currX+9, legY+6)
			line2.Position2 = fyne.NewPos(currX+15, legY+6)
			objs = append(objs, line1, line2)
		} else {
			sq := canvas.NewRectangle(s.Color)
			sq.Resize(fyne.NewSize(12, 8))
			sq.Move(fyne.NewPos(currX, legY+2))
			sq.CornerRadius = 2
			objs = append(objs, sq)
		}

		txt := canvas.NewText(s.Name, colTextPrimary)
		txt.TextSize = 11.5
		txt.TextStyle = fyne.TextStyle{Bold: true}
		txt.Move(fyne.NewPos(currX+18, legY))
		objs = append(objs, txt)

		currX += float32(len(s.Name)*7) + 36
	}

	// ── 2. Dimensiones del Área de Gráfica ─────────────────────────────────
	padLeft, padRight, padTop, padBottom := r.chart.getPaddings()

	plotW := size.Width - padLeft - padRight
	plotH := size.Height - padTop - padBottom

	if plotW <= 0 || plotH <= 0 {
		r.objects = objs
		return
	}

	// ── 3. Cálculo de Escalas de Ejes Y ────────────────────────────────────
	// Eje Izquierdo
	minLeft := math.MaxFloat64
	maxLeft := -math.MaxFloat64
	hasLeft := false

	// Eje Derecho
	minRight := math.MaxFloat64
	maxRight := -math.MaxFloat64
	hasRight := false

	for _, s := range series {
		for _, v := range s.Values {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				continue
			}
			if s.Axis == ChartAxisRight {
				hasRight = true
				if v < minRight {
					minRight = v
				}
				if v > maxRight {
					maxRight = v
				}
			} else {
				hasLeft = true
				if v < minLeft {
					minLeft = v
				}
				if v > maxLeft {
					maxLeft = v
				}
			}
		}
	}

	// Ajuste Eje Izquierdo
	if !hasLeft {
		minLeft = 0
		maxLeft = 100
	} else {
		span := maxLeft - minLeft
		if span <= 0 {
			span = 100
		}

		if r.chart.config.LeftAxis.Min != nil {
			minLeft = *r.chart.config.LeftAxis.Min
		} else {
			// Auto-escalado dinámico alrededor de los datos (10% de holgura inferior, sin bajar de 0)
			minLeft = math.Max(0, minLeft-span*0.10)
		}

		if r.chart.config.LeftAxis.Max != nil {
			maxLeft = *r.chart.config.LeftAxis.Max
		} else {
			// 10% de holgura superior
			maxLeft = maxLeft + span*0.10
		}
	}

	// Ajuste Eje Derecho
	if r.chart.config.RightAxis != nil {
		if !hasRight {
			minRight = 60
			maxRight = 80
		} else {
			if r.chart.config.RightAxis.Min != nil {
				minRight = *r.chart.config.RightAxis.Min
			} else {
				span := maxRight - minRight
				if span < 2 {
					span = 2
				}
				minRight = math.Floor(minRight - span*0.08)
			}
			if r.chart.config.RightAxis.Max != nil {
				maxRight = *r.chart.config.RightAxis.Max
			} else {
				span := maxRight - minRight
				if span < 2 {
					span = 2
				}
				maxRight = math.Ceil(maxRight + span*0.08)
			}
		}
	}

	// ── 4. Títulos de Ejes Y y Cuadrícula Horizontal ───────────────────────
	if r.chart.config.LeftAxis.Title != "" {
		colLeft := r.chart.config.LeftAxis.Color
		if colLeft == nil {
			colLeft = colTextPrimary
		}
		lblLeftTitle := canvas.NewText(r.chart.config.LeftAxis.Title, colLeft)
		lblLeftTitle.TextSize = 10.5
		lblLeftTitle.TextStyle = fyne.TextStyle{Bold: true}
		lblLeftTitle.Move(fyne.NewPos(padLeft+6, padTop-16))
		objs = append(objs, lblLeftTitle)
	}

	if r.chart.config.RightAxis != nil && r.chart.config.RightAxis.Title != "" {
		colRight := r.chart.config.RightAxis.Color
		if colRight == nil {
			colRight = colTextPrimary
		}
		lblRightTitle := canvas.NewText(r.chart.config.RightAxis.Title, colRight)
		lblRightTitle.TextSize = 10.5
		lblRightTitle.TextStyle = fyne.TextStyle{Bold: true}
		lblRightTitle.Move(fyne.NewPos(padLeft+plotW-70, padTop-16))
		objs = append(objs, lblRightTitle)
	}

	gridSteps := 4
	for i := 0; i <= gridSteps; i++ {
		ratio := float32(i) / float32(gridSteps)
		yPos := padTop + plotH - (ratio * plotH)

		// Línea de cuadrícula horizontal
		gridLine := canvas.NewLine(colBorder)
		gridLine.StrokeWidth = 1
		gridLine.Position1 = fyne.NewPos(padLeft, yPos)
		gridLine.Position2 = fyne.NewPos(padLeft+plotW, yPos)
		objs = append(objs, gridLine)

		// Etiqueta Eje Izquierdo — DENTRO del área del gráfico descansando sobre la línea de cuadrícula
		valLeft := minLeft + (float64(ratio) * (maxLeft - minLeft))
		leftStr := fmt.Sprintf("%.0f", valLeft)
		if r.chart.config.LeftAxis.Formatter != nil {
			leftStr = r.chart.config.LeftAxis.Formatter(valLeft)
		}
		lblL := canvas.NewText(leftStr, colTextMuted)
		lblL.TextSize = 10.5
		lblL.TextStyle = fyne.TextStyle{Bold: true}
		lblL.Move(fyne.NewPos(padLeft+6, yPos-13))
		objs = append(objs, lblL)

		// Etiqueta Eje Derecho (si está habilitado)
		if r.chart.config.RightAxis != nil {
			valRight := minRight + (float64(ratio) * (maxRight - minRight))
			rightStr := fmt.Sprintf("%.0f", valRight)
			if r.chart.config.RightAxis.Formatter != nil {
				rightStr = r.chart.config.RightAxis.Formatter(valRight)
			}
			lblR := canvas.NewText(rightStr, colTextMuted)
			lblR.TextSize = 10.5
			lblR.Move(fyne.NewPos(padLeft+plotW+8, yPos-6))
			objs = append(objs, lblR)
		}
	}

	// ── 5. Eje X y Marcas ──────────────────────────────────────────────────
	stepX := plotW / float32(n-1)

	for idx, pt := range points {
		x := padLeft + float32(idx)*stepX

		showLabel := false
		if n <= 12 {
			showLabel = true
		} else if n <= 24 {
			showLabel = (idx%2 == 0) || idx == n-1
		} else {
			showLabel = (idx%4 == 0) || idx == n-1
		}

		if showLabel && pt.XLabel != "" {
			xLbl := canvas.NewText(pt.XLabel, colTextMuted)
			xLbl.TextSize = 10
			xLbl.Move(fyne.NewPos(x-12, padTop+plotH+6))
			objs = append(objs, xLbl)
		}
	}

	if r.chart.config.XAxisTitle != "" {
		lblAxisBottom := canvas.NewText(r.chart.config.XAxisTitle, colTextPrimary)
		lblAxisBottom.TextSize = 11
		lblAxisBottom.TextStyle = fyne.TextStyle{Bold: true}
		lblAxisBottom.Move(fyne.NewPos(padLeft+(plotW/2)-45, padTop+plotH+22))
		objs = append(objs, lblAxisBottom)
	}

	// ── 6. Trazado de Curvas y Puntos ──────────────────────────────────────
	drawContinuousLine := func(pts []fyne.Position, col color.Color, width float32) {
		for i := 0; i < len(pts)-1; i++ {
			if pts[i].Y < 0 || pts[i+1].Y < 0 {
				continue
			}
			line := canvas.NewLine(col)
			line.StrokeWidth = width
			line.Position1 = pts[i]
			line.Position2 = pts[i+1]
			objs = append(objs, line)
		}
	}

	drawDashedLine := func(pts []fyne.Position, col color.Color, width float32) {
		dashLen := float64(5.0)
		gapLen := float64(4.0)
		pattern := dashLen + gapLen

		for i := 0; i < len(pts)-1; i++ {
			p1 := pts[i]
			p2 := pts[i+1]
			if p1.Y < 0 || p2.Y < 0 {
				continue
			}
			dx := float64(p2.X - p1.X)
			dy := float64(p2.Y - p1.Y)
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= 0 {
				continue
			}

			for t := 0.0; t < dist; t += pattern {
				tEnd := math.Min(t+dashLen, dist)
				pos1 := fyne.NewPos(
					p1.X+float32((t/dist)*dx),
					p1.Y+float32((t/dist)*dy),
				)
				pos2 := fyne.NewPos(
					p1.X+float32((tEnd/dist)*dx),
					p1.Y+float32((tEnd/dist)*dy),
				)
				line := canvas.NewLine(col)
				line.StrokeWidth = width
				line.Position1 = pos1
				line.Position2 = pos2
				objs = append(objs, line)
			}
		}
	}

	// Almacenar posiciones calculadas para el hover
	seriesPoints := make([][]fyne.Position, len(series))

	for sIdx, s := range series {
		var sPts []fyne.Position
		minV := minLeft
		maxV := maxLeft
		if s.Axis == ChartAxisRight {
			minV = minRight
			maxV = maxRight
		}

		spanV := maxV - minV
		if spanV <= 0 {
			spanV = 1
		}

		for idx := 0; idx < n; idx++ {
			x := padLeft + float32(idx)*stepX
			val := 0.0
			hasVal := false
			if idx < len(s.Values) {
				val = s.Values[idx]
				if !math.IsNaN(val) && !math.IsInf(val, 0) {
					hasVal = true
				}
			}

			if hasVal {
				// Acotar a rango para que no se salga del canvas
				clampedVal := math.Min(math.Max(val, minV), maxV)
				ratio := float32((clampedVal - minV) / spanV)
				y := padTop + plotH - (ratio * plotH)
				sPts = append(sPts, fyne.NewPos(x, y))
			} else {
				sPts = append(sPts, fyne.NewPos(x, -9999))
			}
		}
		seriesPoints[sIdx] = sPts

		// Dibujar trazo
		width := s.StrokeWidth
		if width <= 0 {
			width = 2.2
		}

		if s.IsDashed {
			drawDashedLine(sPts, s.Color, width)
		} else {
			drawContinuousLine(sPts, s.Color, width)
		}

		// Dibujar vértices discretos
		dotRadius := float32(2.5)
		if width > 2.5 {
			dotRadius = 3.0
		}
		for _, p := range sPts {
			if p.Y < 0 {
				continue
			}
			dot := canvas.NewCircle(s.Color)
			dot.Resize(fyne.NewSize(dotRadius*2, dotRadius*2))
			dot.Move(fyne.NewPos(p.X-dotRadius, p.Y-dotRadius))
			objs = append(objs, dot)
		}
	}

	// ── 7. HOVER: Línea de Referencia y Tarjeta Flotante (Tooltip Popover) ─
	if r.chart.isHovered && r.chart.hoverIndex >= 0 && r.chart.hoverIndex < n {
		hIdx := r.chart.hoverIndex
		pt := points[hIdx]
		refX := padLeft + float32(hIdx)*stepX

		// Paleta adaptativa del Tooltip y elementos guía
		pal := GetCurrentThemePalette()

		var tooltipBgColor color.Color = ThemeColor(ColorNameCardBg)
		var tooltipBorderColor color.Color = ThemeColorAccent()
		var tooltipTitleColor color.Color = ThemeColor(ColorNameTextPrimary)
		var tooltipSubtitleColor color.Color = ThemeColor(ColorNameTextMuted)
		var tooltipSepColor color.Color = ThemeColor(ColorNameDivider)
		var tooltipTextColor color.Color = ThemeColor(ColorNameTextPrimary)
		var tooltipMutedColor color.Color = ThemeColor(ColorNameTextMuted)
		var guideLineColor color.Color = ThemeColorAccent()
		var badgeBgColor color.Color = ThemeColorAccent()
		var badgeTextColor color.Color = ThemeColor(ColorNameTextPrimary)

		if pal != nil {
			tooltipBgColor = pal.Chart.TooltipBg
			tooltipBorderColor = pal.Chart.TooltipBorder
			tooltipTitleColor = pal.Chart.TooltipTitle
			tooltipSubtitleColor = pal.Chart.TooltipSubtitle
			tooltipSepColor = pal.Chart.TooltipDivider
			tooltipTextColor = pal.Chart.TooltipText
			tooltipMutedColor = pal.Chart.TooltipMuted
			guideLineColor = pal.Chart.GuideLine
			badgeBgColor = pal.Chart.BadgeBg
			badgeTextColor = pal.Chart.BadgeText
		}

		// 7.1. Línea de referencia vertical
		guideLine := canvas.NewLine(guideLineColor)
		guideLine.StrokeWidth = 1.6
		guideLine.Position1 = fyne.NewPos(refX, padTop)
		guideLine.Position2 = fyne.NewPos(refX, padTop+plotH)
		objs = append(objs, guideLine)

		// 7.2. Badge indicador en el eje X
		badgeLabel := pt.XLabel
		badgeWidth := float32(math.Max(38, float64(len(badgeLabel)*8+14)))
		badgeBg := canvas.NewRectangle(badgeBgColor)
		badgeBg.CornerRadius = 4
		badgeBg.Resize(fyne.NewSize(badgeWidth, 18))
		badgeBg.Move(fyne.NewPos(refX-(badgeWidth/2), padTop+plotH+3))

		badgeTxt := canvas.NewText(badgeLabel, badgeTextColor)
		badgeTxt.TextSize = 10
		badgeTxt.TextStyle = fyne.TextStyle{Bold: true}
		badgeTxt.Alignment = fyne.TextAlignCenter
		badgeTxt.Move(fyne.NewPos(refX-(badgeWidth/2)+2, padTop+plotH+4))
		objs = append(objs, badgeBg, badgeTxt)

		// 7.3. Halos circulares concéntricos sobre cada serie
		for sIdx, s := range series {
			if sIdx < len(seriesPoints) && hIdx < len(seriesPoints[sIdx]) {
				p := seriesPoints[sIdx][hIdx]
				if p.Y < 0 {
					continue
				}
				halo := canvas.NewCircle(s.Color)
				halo.Resize(fyne.NewSize(12, 12))
				halo.Move(fyne.NewPos(p.X-6, p.Y-6))

				inner := canvas.NewCircle(ColorWhite)
				inner.Resize(fyne.NewSize(6, 6))
				inner.Move(fyne.NewPos(p.X-3, p.Y-3))

				objs = append(objs, halo, inner)
			}
		}

		// 7.4. Ventana-Card Tooltip Popover Flotante
		cardW := float32(236)
		totalRows := len(pt.Details) + len(series)
		cardH := float32(48 + totalRows*18)
		if cardH < 110 {
			cardH = 110
		}

		// Orientación horizontal automática para no salirse de la ventana
		cardX := refX + 16
		if cardX+cardW > size.Width-10 {
			cardX = refX - cardW - 16
		}

		// Posicionamiento vertical automático
		cardY := r.chart.cursorPos.Y - (cardH / 2)
		if cardY < padTop-6 {
			cardY = padTop - 6
		}
		if cardY+cardH > size.Height-8 {
			cardY = size.Height - 8 - cardH
		}

		cardBg := canvas.NewRectangle(tooltipBgColor)
		cardBg.CornerRadius = 8
		cardBg.StrokeColor = tooltipBorderColor
		cardBg.StrokeWidth = 1.6
		cardBg.Resize(fyne.NewSize(cardW, cardH))
		cardBg.Move(fyne.NewPos(cardX, cardY))

		// Sombra nativa Fyne 2.8 acelerada por GPU
		var shadowAlpha uint8 = 30
		if pal != nil && pal.IsDark {
			shadowAlpha = 70
		}
		cardBg.Shadow = canvas.Shadow{
			Color:      color.NRGBA{R: 0, G: 0, B: 0, A: shadowAlpha},
			BlurRadius: 8,
			Spread:     0,
			Offset:     fyne.NewPos(0, 3),
			Variant:    canvas.DropShadow,
		}
		objs = append(objs, cardBg)

		ty := cardY + 8

		// Título del Tooltip
		titleText := pt.Title
		if titleText == "" {
			titleText = pt.XLabel
		}
		tHeader := canvas.NewText(titleText, tooltipTitleColor)
		tHeader.TextSize = 12
		tHeader.TextStyle = fyne.TextStyle{Bold: true}
		tHeader.Move(fyne.NewPos(cardX+10, ty))
		objs = append(objs, tHeader)
		ty += 16

		// Subtítulo (opcional)
		if pt.Subtitle != "" {
			tSub := canvas.NewText(pt.Subtitle, tooltipSubtitleColor)
			tSub.TextSize = 10.5
			tSub.Move(fyne.NewPos(cardX+10, ty))
			objs = append(objs, tSub)
			ty += 16
		}

		// Separador fino
		sep := canvas.NewLine(tooltipSepColor)
		sep.StrokeWidth = 1
		sep.Position1 = fyne.NewPos(cardX+10, ty)
		sep.Position2 = fyne.NewPos(cardX+cardW-10, ty)
		objs = append(objs, sep)
		ty += 6

		// Filas de Detalles adicionales (si existen)
		for _, det := range pt.Details {
			colVal := det.Color
			if colVal == nil {
				if det.Bold && pal != nil {
					colVal = pal.Success
				} else {
					colVal = tooltipTextColor
				}
			}

			tDetLabel := canvas.NewText(det.Label+":", tooltipMutedColor)
			tDetLabel.TextSize = 10.5
			tDetLabel.Move(fyne.NewPos(cardX+10, ty))

			tDetVal := canvas.NewText(det.Value, colVal)
			tDetVal.TextSize = 10.5
			if det.Bold {
				tDetVal.TextStyle = fyne.TextStyle{Bold: true}
			}
			tDetVal.Move(fyne.NewPos(cardX+float32(len(det.Label)*6)+18, ty))
			objs = append(objs, tDetLabel, tDetVal)
			ty += 18
		}

		// Filas de cada Serie
		for _, s := range series {
			dot := canvas.NewRectangle(s.Color)
			if s.IsDashed {
				dot.Resize(fyne.NewSize(8, 3))
				dot.Move(fyne.NewPos(cardX+10, ty+5))
			} else {
				dot.Resize(fyne.NewSize(8, 8))
				dot.CornerRadius = 2
				dot.Move(fyne.NewPos(cardX+10, ty+3))
			}
			objs = append(objs, dot)

			val := 0.0
			hasVal := false
			if hIdx < len(s.Values) {
				val = s.Values[hIdx]
				if !math.IsNaN(val) && !math.IsInf(val, 0) {
					hasVal = true
				}
			}
			valStr := "N/D"
			if hasVal {
				valStr = fmt.Sprintf("%.2f", val)
				if s.ValueFormatter != nil {
					valStr = s.ValueFormatter(val)
				}
			} else if s.ValueFormatter != nil {
				valStr = s.ValueFormatter(math.NaN())
			}

			tSeriesName := canvas.NewText(s.Name+":", tooltipMutedColor)
			tSeriesName.TextSize = 10.5
			tSeriesName.Move(fyne.NewPos(cardX+24, ty))

			tSeriesVal := canvas.NewText(valStr, tooltipTextColor)
			tSeriesVal.TextSize = 10.5
			tSeriesVal.TextStyle = fyne.TextStyle{Bold: true}
			tSeriesVal.Move(fyne.NewPos(cardX+float32(len(s.Name)*6)+32, ty))

			objs = append(objs, tSeriesName, tSeriesVal)
			ty += 18
		}
	}

	r.objects = objs
	canvas.Refresh(r.chart)
}
