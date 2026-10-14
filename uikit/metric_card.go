package uikit

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// Alias de tipo para máxima concisión estilo xBase / Harbour
type MetricCard = ModernMetricCard

// ModernMetricCard es una tarjeta KPI moderna para dashboards de métricas y analítica.
// Presenta un título superior, un icono en caja redondeada con color distintivo,
// una cifra destacada de gran tamaño (hero) y dos etiquetas informativas al pie.
type ModernMetricCard struct {
	Card       *ModernCard
	lblTitle   *ModernLabel
	iconBox    *fyne.Container
	lblValue   *ModernLabel
	lblFootL   *ModernLabel
	lblFootR   *ModernLabel
	valueColor color.Color
}

// NewMetricCard crea una nueva tarjeta de métrica KPI lista para incrustar en layouts.
func NewMetricCard(
	title string,
	iconName string,
	iconHex string,
	iconBgHex string,
	valueColor color.Color,
) *ModernMetricCard {
	mc := &ModernMetricCard{
		valueColor: valueColor,
		lblTitle:   NewMicroLabel(title),
		lblValue:   NewHeroLabel("--", valueColor),
		lblFootL:   NewCaptionLabel(""),
		lblFootR:   NewCaptionLabel("").SetBold(true),
	}

	// Icono en caja pequeña redondeada
	iBg := canvas.NewRectangle(HexColor(iconBgHex))
	iBg.CornerRadius = 8
	iBg.Resize(fyne.NewSize(28, 28))
	icon := NewColoredIcon(iconName, 15, iconHex)
	mc.iconBox = container.NewStack(iBg, container.NewCenter(icon))

	topRow := container.NewBorder(
		nil, nil,
		mc.lblTitle,
		mc.iconBox,
	)

	footRow := container.NewBorder(
		nil, nil,
		mc.lblFootL,
		mc.lblFootR,
	)

	body := container.NewVBox(
		topRow,
		container.NewPadded(mc.lblValue),
		footRow,
	)

	mc.Card = NewCard(body).WithVerticalPadding(18, 18)
	return mc
}

// SetValue actualiza el valor principal de la tarjeta de métrica.
func (mc *ModernMetricCard) SetValue(v string) *ModernMetricCard {
	mc.lblValue.SetText(v)
	return mc
}

// SetFooters actualiza las dos etiquetas informativas inferiores con el color temático de alto contraste.
func (mc *ModernMetricCard) SetFooters(left, right string, rightColor ...color.Color) *ModernMetricCard {
	mc.lblFootL.SetText(left)
	mc.lblFootL.SetColor(ThemeColor(ColorNameTextMuted))
	mc.lblFootR.SetText(right)
	if len(rightColor) > 0 && rightColor[0] != nil {
		mc.lblFootR.SetColor(rightColor[0])
	} else {
		mc.lblFootR.SetColor(ThemeColor(ColorNameTextMuted))
	}
	return mc
}

// SetTitle actualiza el título superior de la métrica.
func (mc *ModernMetricCard) SetTitle(title string) *ModernMetricCard {
	mc.lblTitle.SetText(title)
	mc.lblTitle.SetColor(ThemeColor(ColorNameTextMuted))
	return mc
}

// SetValueColor actualiza el color de la cifra principal.
func (mc *ModernMetricCard) SetValueColor(c color.Color) *ModernMetricCard {
	mc.valueColor = c
	mc.lblValue.SetColor(c)
	return mc
}
