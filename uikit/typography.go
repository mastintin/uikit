package uikit

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// Paleta cromática semántica para textos de alto contraste
var (
	ColorWhite   = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	ColorLight   = HexColor("#e2e8f0") // Slate-200
	ColorMuted   = HexColor("#94a3b8") // Slate-400
	ColorDim     = HexColor("#64748b") // Slate-500
	ColorPrimary = HexColor("#818cf8") // Indigo-400
	ColorAccent  = HexColor("#60a5fa") // Blue-400
	ColorSuccess = HexColor("#34d399") // Emerald-400
	ColorWarning = HexColor("#fbbf24") // Amber-400
	ColorDanger  = HexColor("#f87171") // Red-400
)

// TypographyStyle encapsula las propiedades de fuente y escala de la aplicación
type TypographyStyle struct {
	Size      float32
	Bold      bool
	Monospace bool
}

// Estilos tipográficos centralizados del Design System de GoJubila:
var (
	// StyleHero: Grandes cifras monetarias y KPIs principales (26pt, Negrita)
	StyleHero = TypographyStyle{Size: 26, Bold: true}

	// StyleTitle1: Títulos de primer nivel / Pantalla (22pt, Negrita)
	StyleTitle1 = TypographyStyle{Size: 22, Bold: true}

	// StyleTitle2: Subtítulos destacados, método ganador, nombres en fichas (18pt, Negrita)
	StyleTitle2 = TypographyStyle{Size: 18, Bold: true}

	// StyleSection: Títulos de sección, columnas y tarjetas (14.5pt, Negrita)
	StyleSection = TypographyStyle{Size: 14.5, Bold: true}

	// StyleHighlight: Totales netos destacados secundarios (16pt, Negrita)
	StyleHighlight = TypographyStyle{Size: 16, Bold: true}

	// StyleMilestoneDate: Fechas de hitos legales en tarjetas de cabecera (15pt, Negrita)
	StyleMilestoneDate = TypographyStyle{Size: 15, Bold: true}

	// StyleValue: Importes monetarios y valores de tabla/desglose (13pt, Regular)
	StyleValue = TypographyStyle{Size: 13, Bold: false}

	// StyleValueBold: Importes monetarios destacados en tabla/desglose (13pt, Negrita)
	StyleValueBold = TypographyStyle{Size: 13, Bold: true}

	// StyleBody: Párrafos, subtítulos descriptivos y textos explicativos (12.5pt, Regular)
	StyleBody = TypographyStyle{Size: 12.5, Bold: false}

	// StyleBodyBold: Textos descriptivos con énfasis (12.5pt, Negrita)
	StyleBodyBold = TypographyStyle{Size: 12.5, Bold: true}

	// StyleCaption: Etiquetas de concepto en filas de detalle (12pt, Regular)
	StyleCaption = TypographyStyle{Size: 12, Bold: false}

	// StyleCaptionBold: Etiquetas de concepto con énfasis (12pt, Negrita)
	StyleCaptionBold = TypographyStyle{Size: 12, Bold: true}

	// StyleBadge: Textos dentro de pastillas e indicadores (11.5pt, Negrita)
	StyleBadge = TypographyStyle{Size: 11.5, Bold: true}

	// StyleMicro: Subtítulos técnicos y etiquetas compactas (10.5pt, Negrita)
	StyleMicro = TypographyStyle{Size: 10.5, Bold: true}
)

// ApplyStyle aplica un TypographyStyle a un objeto canvas.Text existente
func ApplyStyle(txt *canvas.Text, style TypographyStyle) *canvas.Text {
	if txt == nil {
		return nil
	}
	txt.TextSize = style.Size
	txt.TextStyle = fyne.TextStyle{
		Bold:      style.Bold,
		Monospace: style.Monospace,
	}
	return txt
}

// Funciones modificadoras directas para aplicar a textos existentes:
func ApplyHero(txt *canvas.Text) *canvas.Text          { return ApplyStyle(txt, StyleHero) }
func ApplyTitle1(txt *canvas.Text) *canvas.Text        { return ApplyStyle(txt, StyleTitle1) }
func ApplyTitle2(txt *canvas.Text) *canvas.Text        { return ApplyStyle(txt, StyleTitle2) }
func ApplySection(txt *canvas.Text) *canvas.Text       { return ApplyStyle(txt, StyleSection) }
func ApplyHighlight(txt *canvas.Text) *canvas.Text     { return ApplyStyle(txt, StyleHighlight) }
func ApplyValue(txt *canvas.Text, bold ...bool) *canvas.Text {
	if len(bold) > 0 && bold[0] {
		return ApplyStyle(txt, StyleValueBold)
	}
	return ApplyStyle(txt, StyleValue)
}
func ApplyBody(txt *canvas.Text, bold ...bool) *canvas.Text {
	if len(bold) > 0 && bold[0] {
		return ApplyStyle(txt, StyleBodyBold)
	}
	return ApplyStyle(txt, StyleBody)
}
func ApplyCaption(txt *canvas.Text, bold ...bool) *canvas.Text {
	if len(bold) > 0 && bold[0] {
		return ApplyStyle(txt, StyleCaptionBold)
	}
	return ApplyStyle(txt, StyleCaption)
}
func ApplyBadge(txt *canvas.Text) *canvas.Text { return ApplyStyle(txt, StyleBadge) }
func ApplyMicro(txt *canvas.Text) *canvas.Text { return ApplyStyle(txt, StyleMicro) }

// Constructores de una sola línea que combinan texto, estilo y color opcional:

func NewStyledText(text string, style TypographyStyle, col ...color.Color) *canvas.Text {
	c := color.Color(ColorWhite)
	if len(col) > 0 && col[0] != nil {
		c = col[0]
	}
	txt := canvas.NewText(text, c)
	return ApplyStyle(txt, style)
}

func NewHeroText(text string, col ...color.Color) *canvas.Text {
	return NewStyledText(text, StyleHero, col...)
}

func NewTitle1Text(text string, col ...color.Color) *canvas.Text {
	return NewStyledText(text, StyleTitle1, col...)
}

func NewTitle2Text(text string, col ...color.Color) *canvas.Text {
	return NewStyledText(text, StyleTitle2, col...)
}

func NewSectionText(text string, col ...color.Color) *canvas.Text {
	return NewStyledText(text, StyleSection, col...)
}

func NewHighlightText(text string, col ...color.Color) *canvas.Text {
	return NewStyledText(text, StyleHighlight, col...)
}

func NewValueText(text string, bold bool, col ...color.Color) *canvas.Text {
	if bold {
		return NewStyledText(text, StyleValueBold, col...)
	}
	return NewStyledText(text, StyleValue, col...)
}

func NewBodyText(text string, col ...color.Color) *canvas.Text {
	return NewStyledText(text, StyleBody, col...)
}

func NewCaptionText(text string, col ...color.Color) *canvas.Text {
	c := color.Color(ColorMuted)
	if len(col) > 0 && col[0] != nil {
		c = col[0]
	}
	return NewStyledText(text, StyleCaption, c)
}

func NewMicroText(text string, col ...color.Color) *canvas.Text {
	c := color.Color(ColorMuted)
	if len(col) > 0 && col[0] != nil {
		c = col[0]
	}
	return NewStyledText(text, StyleMicro, c)
}
