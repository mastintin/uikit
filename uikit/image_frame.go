package uikit

import (
	"image"
	"image/color"
	"math"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// ModernImageFrame implementa un marco interactivo para imágenes con esquinas curvas,
// escalado responsive (FillContain), efectos dinámicos de hover y cursor, y evento al hacer clic.
//
// Comparte la misma filosofía arquitectónica de ModernButton, ModernBadge, ModernTextBox y ModernCard:
// - Extiende widget.BaseWidget como ciudadano nativo de Fyne.
// - Constructores unificados estilo Harbour / xBase y constructores semánticos ("Sugar").
// - Interfaz fluida (Fluent API) para parametrizar bordes, radios, colores y tamaños.
// - Métodos de alto nivel (.SetFile, .SetImage, .FilePath, .Clear, .SetOnTap).
type ModernImageFrame struct {
	widget.BaseWidget

	bg        *canvas.Rectangle
	img       *canvas.Image
	innerWrap *fyne.Container
	container *fyne.Container

	onTap func()

	cornerRadius      float32
	strokeWidth       float32
	normalBorderColor color.Color
	hoverBorderColor  color.Color
	bgColor           color.Color
	minSize           fyne.Size
	isHovered         bool
	hasPadding        bool
}

// Alias de tipo para compatibilidad total hacia atrás
type ImagePreviewFrame = ModernImageFrame

// Asegurar que implementa interfaces interactivas de Fyne
var (
	_ fyne.Tappable      = (*ModernImageFrame)(nil)
	_ desktop.Hoverable  = (*ModernImageFrame)(nil)
	_ desktop.Cursorable = (*ModernImageFrame)(nil)
)

// NewImageFrame es el constructor unificado estilo Harbour / xBase para crear un marco de imagen.
// - fileOrPath: ruta al archivo de imagen o "" para placeholder
// - cornerRadius: radio de las esquinas en px (ej. 8)
// - bgHex: color de fondo hexadecimal (ej. "#0f172a")
// - borderHex: color del borde perimetral (ej. "#334155")
// - onTap: callback ejecutado al pulsar la imagen (ej. abrir en Vista Previa)
func NewImageFrame(fileOrPath string, cornerRadius float32, bgHex, borderHex string, onTap func()) *ModernImageFrame {
	f := &ModernImageFrame{
		onTap:             onTap,
		cornerRadius:      cornerRadius,
		strokeWidth:       1.5,
		normalBorderColor: HexColor(borderHex),
		hoverBorderColor:  ThemeColorAccent(),
		bgColor:           HexColor(bgHex),
		minSize:           fyne.NewSize(180, 180),
		hasPadding:        true,
	}

	// 1. Fondo estilizado con esquinas redondeadas
	f.bg = canvas.NewRectangle(f.bgColor)
	f.bg.CornerRadius = f.cornerRadius
	f.bg.StrokeWidth = f.strokeWidth
	f.bg.StrokeColor = f.normalBorderColor

	// 2. Imagen interna centrada con FillContain para no deformarse
	f.img = canvas.NewImageFromImage(PlaceholderThumbnail(340, 340))
	f.img.FillMode = canvas.ImageFillContain

	// 3. Montaje del contenedor
	f.rebuildContainer()

	f.ExtendBaseWidget(f)

	if fileOrPath != "" {
		f.SetFile(fileOrPath)
	}

	return f
}

// NewImagePreviewFrame es el constructor semántico por defecto ("Sugar Constructor"):
// crea el marco oscuro pizarra (#0f172a) con bordes curvos (8px) y borde #334155
// característico de goCuadros (.image-preview-wrapper).
func NewImagePreviewFrame(onTap func()) *ModernImageFrame {
	return NewImageFrame("", 8, "#0f172a", "#334155", onTap)
}

func (f *ModernImageFrame) rebuildContainer() {
	var content fyne.CanvasObject = f.img
	if f.hasPadding {
		content = container.NewPadded(f.img)
	}
	f.innerWrap = container.NewStack(f.bg, content)
	if f.container == nil {
		f.container = container.NewStack(f.innerWrap)
	} else {
		f.container.Objects = []fyne.CanvasObject{f.innerWrap}
		f.container.Refresh()
	}
}

func (f *ModernImageFrame) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(f.container)
}

// MinSize define un tamaño mínimo que permite al marco encogerse
// o expandirse fluidamente según el espacio del contenedor padre.
func (f *ModernImageFrame) MinSize() fyne.Size {
	if f.minSize.Width > 0 && f.minSize.Height > 0 {
		return f.minSize
	}
	return fyne.NewSize(180, 180)
}

// --- Métodos de Alto Nivel y Estado ---

// SetFile carga una imagen desde la ruta de archivo indicada.
func (f *ModernImageFrame) SetFile(path string) {
	if path != "" {
		if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
			f.img.File = path
			f.img.Image = nil
			f.img.Refresh()
			return
		}
	}
	f.SetImage(PlaceholderThumbnail(340, 340))
}

// SetImage asigna una imagen en memoria (ej. placeholder o raster).
func (f *ModernImageFrame) SetImage(img image.Image) {
	f.img.File = ""
	f.img.Image = img
	f.img.Refresh()
}

// FilePath retorna la ruta de archivo cargada actualmente (o cadena vacía).
func (f *ModernImageFrame) FilePath() string {
	return f.img.File
}

// Clear restaura el marco a su imagen de marcador de posición (placeholder).
func (f *ModernImageFrame) Clear() {
	f.SetImage(PlaceholderThumbnail(340, 340))
}

// SetOnTap asigna o actualiza la acción al pulsar la imagen.
func (f *ModernImageFrame) SetOnTap(onTap func()) {
	f.onTap = onTap
}

// --- Interfaz Fluida (Fluent API) ---

// WithRadius ajusta el radio de curvatura de las esquinas.
func (f *ModernImageFrame) WithRadius(radius float32) *ModernImageFrame {
	f.cornerRadius = radius
	f.bg.CornerRadius = radius
	f.bg.Refresh()
	return f
}

// WithBackground cambia el color de fondo mediante código hexadecimal.
func (f *ModernImageFrame) WithBackground(bgHex string) *ModernImageFrame {
	f.bgColor = HexColor(bgHex)
	f.bg.FillColor = f.bgColor
	f.bg.Refresh()
	return f
}

// WithBorder configura los colores y grosor de borde en estado normal y hover.
func (f *ModernImageFrame) WithBorder(normalHex, hoverHex string, width float32) *ModernImageFrame {
	f.normalBorderColor = HexColor(normalHex)
	f.hoverBorderColor = HexColor(hoverHex)
	f.strokeWidth = width
	f.bg.StrokeWidth = width
	if f.isHovered {
		f.bg.StrokeColor = f.hoverBorderColor
	} else {
		f.bg.StrokeColor = f.normalBorderColor
	}
	f.bg.Refresh()
	return f
}

// WithMinSize personaliza el tamaño mínimo que puede adoptar el marco.
func (f *ModernImageFrame) WithMinSize(size fyne.Size) *ModernImageFrame {
	f.minSize = size
	f.Refresh()
	return f
}

// WithPadding habilita o deshabilita el margen de seguridad interno.
func (f *ModernImageFrame) WithPadding(padding bool) *ModernImageFrame {
	f.hasPadding = padding
	f.rebuildContainer()
	f.Refresh()
	return f
}

// WithFillMode permite configurar el modo de ajuste (Contain, Stretch, Original).
func (f *ModernImageFrame) WithFillMode(mode canvas.ImageFill) *ModernImageFrame {
	f.img.FillMode = mode
	f.img.Refresh()
	return f
}

// --- Eventos Interactivos ---

// Tapped responde al clic en cualquier parte del marco o imagen.
func (f *ModernImageFrame) Tapped(_ *fyne.PointEvent) {
	if f.onTap != nil {
		f.onTap()
	}
}

// Cursor devuelve PointerCursor para indicar que es interactivo.
func (f *ModernImageFrame) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

// MouseIn resalta el borde al pasar el ratón.
func (f *ModernImageFrame) MouseIn(_ *desktop.MouseEvent) {
	f.isHovered = true
	f.bg.StrokeColor = f.hoverBorderColor
	f.bg.Refresh()
}

// MouseMoved cumple con la interfaz desktop.Hoverable.
func (f *ModernImageFrame) MouseMoved(_ *desktop.MouseEvent) {}

// MouseOut restaura el borde normal al salir el ratón.
func (f *ModernImageFrame) MouseOut() {
	f.isHovered = false
	f.bg.StrokeColor = f.normalBorderColor
	f.bg.Refresh()
}

// ApplyRoundedCorners recorta las 4 esquinas de la imagen aplicando una máscara circular con antialiasing.
func ApplyRoundedCorners(img *image.RGBA, radius float64) {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	if radius <= 0 || w <= 0 || h <= 0 {
		return
	}

	r := radius
	if r > float64(w)/2 {
		r = float64(w) / 2
	}
	if r > float64(h)/2 {
		r = float64(h) / 2
	}

	cornerSize := int(math.Ceil(r + 1.0))
	if cornerSize > w/2 {
		cornerSize = w / 2
	}
	if cornerSize > h/2 {
		cornerSize = h / 2
	}

	minX, minY := bounds.Min.X, bounds.Min.Y
	maxX, maxY := bounds.Max.X, bounds.Max.Y

	checkPixel := func(x, y int, cx, cy float64) {
		dx := (float64(x) + 0.5) - cx
		dy := (float64(y) + 0.5) - cy
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist >= r+0.5 {
			img.SetRGBA(x, y, color.RGBA{0, 0, 0, 0})
		} else if dist > r-0.5 {
			factor := (r + 0.5 - dist)
			c := img.RGBAAt(x, y)
			newA := uint8(math.Round(float64(c.A) * factor))
			newR := uint8(math.Round(float64(c.R) * factor))
			newG := uint8(math.Round(float64(c.G) * factor))
			newB := uint8(math.Round(float64(c.B) * factor))
			img.SetRGBA(x, y, color.RGBA{R: newR, G: newG, B: newB, A: newA})
		}
	}

	// Esquina Superior Izquierda
	cxTL := float64(minX) + r
	cyTL := float64(minY) + r
	for y := minY; y < minY+cornerSize; y++ {
		for x := minX; x < minX+cornerSize; x++ {
			checkPixel(x, y, cxTL, cyTL)
		}
	}

	// Esquina Superior Derecha
	cxTR := float64(maxX) - r
	cyTR := float64(minY) + r
	for y := minY; y < minY+cornerSize; y++ {
		for x := maxX - cornerSize; x < maxX; x++ {
			checkPixel(x, y, cxTR, cyTR)
		}
	}

	// Esquina Inferior Izquierda
	cxBL := float64(minX) + r
	cyBL := float64(maxY) - r
	for y := maxY - cornerSize; y < maxY; y++ {
		for x := minX; x < minX+cornerSize; x++ {
			checkPixel(x, y, cxBL, cyBL)
		}
	}

	// Esquina Inferior Derecha
	cxBR := float64(maxX) - r
	cyBR := float64(maxY) - r
	for y := maxY - cornerSize; y < maxY; y++ {
		for x := maxX - cornerSize; x < maxX; x++ {
			checkPixel(x, y, cxBR, cyBR)
		}
	}
}

// PlaceholderThumbnail genera una miniatura cuadrada gris suave (#e2e8f0) con esquinas redondeadas.
func PlaceholderThumbnail(w, h int) image.Image {
	size := w
	if h > size {
		size = h
	}
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	bg := color.RGBA{R: 226, G: 232, B: 240, A: 255}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			img.Set(x, y, bg)
		}
	}
	radius := float64(size) * (8.0 / 60.0)
	if radius > 16 {
		radius = 16
	}
	ApplyRoundedCorners(img, radius)
	return img
}
