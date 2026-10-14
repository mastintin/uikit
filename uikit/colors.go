package uikit

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ─────────────────────────────────────────────────────────────────────────────
// 1. Normalización y Conversión de Formatos de Color
// ─────────────────────────────────────────────────────────────────────────────

// ColorToNRGBA convierte cualquier color.Color de Go a un color.NRGBA de 8 bits.
// Des-premultiplica con precisión matemática el canal alfa si el color proviene
// de formatos premultiplicados (como color.RGBA o la interfaz genérica c.RGBA()),
// asegurando máxima fidelidad de color incluso con transparencias intermedias.
func ColorToNRGBA(c color.Color) color.NRGBA {
	if c == nil {
		return color.NRGBA{}
	}

	switch col := c.(type) {
	case color.NRGBA:
		return col
	case *color.NRGBA:
		if col == nil {
			return color.NRGBA{}
		}
		return *col
	case color.NRGBA64:
		return color.NRGBA{
			R: uint8(col.R >> 8),
			G: uint8(col.G >> 8),
			B: uint8(col.B >> 8),
			A: uint8(col.A >> 8),
		}
	case *color.NRGBA64:
		if col == nil {
			return color.NRGBA{}
		}
		return color.NRGBA{
			R: uint8(col.R >> 8),
			G: uint8(col.G >> 8),
			B: uint8(col.B >> 8),
			A: uint8(col.A >> 8),
		}
	case color.Gray:
		return color.NRGBA{R: col.Y, G: col.Y, B: col.Y, A: 255}
	case *color.Gray:
		if col == nil {
			return color.NRGBA{}
		}
		return color.NRGBA{R: col.Y, G: col.Y, B: col.Y, A: 255}
	case color.Gray16:
		y := uint8(col.Y >> 8)
		return color.NRGBA{R: y, G: y, B: y, A: 255}
	case *color.Gray16:
		if col == nil {
			return color.NRGBA{}
		}
		y := uint8(col.Y >> 8)
		return color.NRGBA{R: y, G: y, B: y, A: 255}
	case color.Alpha:
		return color.NRGBA{R: 255, G: 255, B: 255, A: col.A}
	case *color.Alpha:
		if col == nil {
			return color.NRGBA{}
		}
		return color.NRGBA{R: 255, G: 255, B: 255, A: col.A}
	case color.Alpha16:
		return color.NRGBA{R: 255, G: 255, B: 255, A: uint8(col.A >> 8)}
	case *color.Alpha16:
		if col == nil {
			return color.NRGBA{}
		}
		return color.NRGBA{R: 255, G: 255, B: 255, A: uint8(col.A >> 8)}
	default:
		r, g, b, a := c.RGBA()
		if a == 0 {
			return color.NRGBA{0, 0, 0, 0}
		}
		if a == 0xffff {
			return color.NRGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: 255,
			}
		}
		// Des-premultiplicación precisa: (canal * 0xffff) / a
		return color.NRGBA{
			R: uint8(((r * 0xffff) / a) >> 8),
			G: uint8(((g * 0xffff) / a) >> 8),
			B: uint8(((b * 0xffff) / a) >> 8),
			A: uint8(a >> 8),
		}
	}
}

// HexColor convierte un código hexadecimal CSS (#RGB, #RGBA, #RRGGBB o #RRGGBBAA)
// a color.NRGBA. Si el código no es válido, devuelve color.Transparent.
func HexColor(hex string) color.Color {
	hex = strings.TrimPrefix(strings.TrimSpace(hex), "#")
	var r, g, b, a uint8 = 0, 0, 0, 255
	switch len(hex) {
	case 3: // rgb
		n, err := fmt.Sscanf(hex, "%1x%1x%1x", &r, &g, &b)
		if err != nil || n != 3 {
			return color.Transparent
		}
		r = r*16 + r
		g = g*16 + g
		b = b*16 + b
	case 4: // rgba
		n, err := fmt.Sscanf(hex, "%1x%1x%1x%1x", &r, &g, &b, &a)
		if err != nil || n != 4 {
			return color.Transparent
		}
		r = r*16 + r
		g = g*16 + g
		b = b*16 + b
		a = a*16 + a
	case 6: // rrggbb
		var rgb uint32
		n, err := fmt.Sscanf(hex, "%06x", &rgb)
		if err != nil || n != 1 {
			return color.Transparent
		}
		r = uint8(rgb >> 16)
		g = uint8((rgb >> 8) & 0xFF)
		b = uint8(rgb & 0xFF)
	case 8: // rrggbbaa
		var rgba uint64
		n, err := fmt.Sscanf(hex, "%08x", &rgba)
		if err != nil || n != 1 {
			return color.Transparent
		}
		r = uint8(rgba >> 24)
		g = uint8((rgba >> 16) & 0xFF)
		b = uint8((rgba >> 8) & 0xFF)
		a = uint8(rgba & 0xFF)
	default:
		return color.Transparent
	}
	return color.NRGBA{R: r, G: g, B: b, A: a}
}

// HexToNRGBA convierte una cadena hexadecimal a color.NRGBA.
func HexToNRGBA(hex string) color.NRGBA {
	return ColorToNRGBA(HexColor(hex))
}

// ColorToHex convierte un color.Color a formato string hexadecimal "#rrggbb" (o "#rrggbbaa" si tiene canal alfa).
func ColorToHex(c color.Color) string {
	if c == nil {
		return ""
	}
	nrgba := ColorToNRGBA(c)
	if nrgba.A == 255 {
		return fmt.Sprintf("#%02x%02x%02x", nrgba.R, nrgba.G, nrgba.B)
	}
	return fmt.Sprintf("#%02x%02x%02x%02x", nrgba.R, nrgba.G, nrgba.B, nrgba.A)
}

// ─────────────────────────────────────────────────────────────────────────────
// 2. Manipulación de Opacidad e Interpolación
// ─────────────────────────────────────────────────────────────────────────────

// ColorWithAlpha devuelve una copia del color indicado ajustando su opacidad
// mediante un factor escalar (0.0 = totalmente transparente, 1.0 = opacidad original intacta).
func ColorWithAlpha(c color.Color, factor float32) color.Color {
	if c == nil {
		return color.Transparent
	}
	if factor < 0 {
		factor = 0
	} else if factor > 1 {
		factor = 1
	}
	n := ColorToNRGBA(c)
	n.A = uint8(float32(n.A)*factor + 0.5)
	return n
}

// BlendColor interpola linealmente entre dos colores según un factor t (0.0 = c1, 1.0 = c2).
func BlendColor(c1, c2 color.Color, t float32) color.Color {
	if t <= 0 {
		return c1
	}
	if t >= 1 {
		return c2
	}
	n1 := ColorToNRGBA(c1)
	n2 := ColorToNRGBA(c2)

	invT := 1.0 - t
	r := uint8(float32(n1.R)*invT + float32(n2.R)*t + 0.5)
	g := uint8(float32(n1.G)*invT + float32(n2.G)*t + 0.5)
	b := uint8(float32(n1.B)*invT + float32(n2.B)*t + 0.5)
	a := uint8(float32(n1.A)*invT + float32(n2.A)*t + 0.5)

	return color.NRGBA{R: r, G: g, B: b, A: a}
}

// ─────────────────────────────────────────────────────────────────────────────
// 3. Brillo y Oscurecimiento con Conservación de Tono
// ─────────────────────────────────────────────────────────────────────────────

// ColorBrighter aclara un color en la proporción especificada (entre 0.0 y 1.0).
// Utiliza un algoritmo de redistribución armónica de luminosidad que evita
// virajes cromáticos indeseados cuando los canales primarios saturan a 255.
func ColorBrighter(c color.Color, factor float32) color.Color {
	if c == nil {
		return nil
	}
	if factor <= 0 {
		return c
	}
	if factor > 1 {
		factor = 1
	}

	n := ColorToNRGBA(c)
	multiplier := 1.0 + factor

	rf := (float32(n.R) / 255.0) * multiplier
	gf := (float32(n.G) / 255.0) * multiplier
	bf := (float32(n.B) / 255.0) * multiplier

	rf, gf, bf = balanceLuminanceRGB(rf, gf, bf)

	return color.NRGBA{
		R: clampChannel(rf * 255.0),
		G: clampChannel(gf * 255.0),
		B: clampChannel(bf * 255.0),
		A: n.A,
	}
}

// ColorDarker oscurece un color en la proporción especificada (entre 0.0 y 1.0).
func ColorDarker(c color.Color, factor float32) color.Color {
	if c == nil {
		return nil
	}
	if factor <= 0 {
		return c
	}
	if factor > 1 {
		factor = 1
	}

	n := ColorToNRGBA(c)
	divisor := 1.0 + factor

	rf := (float32(n.R) / 255.0) / divisor
	gf := (float32(n.G) / 255.0) / divisor
	bf := (float32(n.B) / 255.0) / divisor

	return color.NRGBA{
		R: clampChannel(rf * 255.0),
		G: clampChannel(gf * 255.0),
		B: clampChannel(bf * 255.0),
		A: n.A,
	}
}

// balanceLuminanceRGB redistribuye el exceso lumínico cuando algún canal sobrepasa 1.0,
// evitando la pérdida de pureza y la deriva hacia blanco o amarillo en colores saturados.
func balanceLuminanceRGB(r, g, b float32) (float32, float32, float32) {
	maxVal := r
	if g > maxVal {
		maxVal = g
	}
	if b > maxVal {
		maxVal = b
	}

	if maxVal <= 1.0 {
		return r, g, b
	}

	sum := r + g + b
	if sum >= 3.0 {
		return 1.0, 1.0, 1.0
	}

	ratio := (3.0 - sum) / (3.0*maxVal - sum)
	neutral := 1.0 - ratio*maxVal

	return neutral + ratio*r, neutral + ratio*g, neutral + ratio*b
}

func clampChannel(v float32) uint8 {
	if v <= 0 {
		return 0
	}
	if v >= 255 {
		return 255
	}
	return uint8(v + 0.5)
}

// ─────────────────────────────────────────────────────────────────────────────
// 4. Percepción de Luminancia y Contraste WCAG
// ─────────────────────────────────────────────────────────────────────────────

// ColorLuminance calcula la luminancia relativa perceptual según ITU-R BT.709
// con corrección gamma sRGB a espacio lineal (rango 0.0 a 1.0).
func ColorLuminance(c color.Color) float32 {
	if c == nil {
		return 0
	}
	n := ColorToNRGBA(c)
	r := float64(n.R) / 255.0
	g := float64(n.G) / 255.0
	b := float64(n.B) / 255.0

	toLinear := func(val float64) float64 {
		if val <= 0.04045 {
			return val / 12.92
		}
		return math.Pow((val+0.055)/1.055, 2.4)
	}

	lum := 0.2126*toLinear(r) + 0.7152*toLinear(g) + 0.0722*toLinear(b)
	return float32(lum)
}

// IsColorDark indica si un color es oscuro (luminancia perceptual < 0.45).
func IsColorDark(c color.Color) bool {
	return ColorLuminance(c) < 0.45
}

// ContrastTextColor selecciona automáticamente blanco puro (#ffffff) para fondos oscuros
// o pizarra oscura (#0f172a) para fondos claros, garantizando óptima legibilidad.
func ContrastTextColor(bg color.Color) color.Color {
	if IsColorDark(bg) {
		return ColorWhite
	}
	return ColorSlate900
}

// ─────────────────────────────────────────────────────────────────────────────
// 5. Mapeo Semántico de Importancia
// ─────────────────────────────────────────────────────────────────────────────

// ButtonForegroundFor devuelve el ThemeColorName que corresponde al nivel de
// importancia de un botón de Fyne (High, Danger, Success, Warning o Normal).
func ButtonForegroundFor(importance widget.Importance) fyne.ThemeColorName {
	switch importance {
	case widget.HighImportance:
		return theme.ColorNameForegroundOnPrimary
	case widget.DangerImportance:
		return theme.ColorNameForegroundOnError
	case widget.SuccessImportance:
		return theme.ColorNameForegroundOnSuccess
	case widget.WarningImportance:
		return theme.ColorNameForegroundOnWarning
	default:
		return theme.ColorNameForeground
	}
}
