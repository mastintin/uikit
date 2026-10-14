package uikit

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SidebarStyle int

const (
	SidebarStyleDefault SidebarStyle = iota
	SidebarStyleShaded
)

// ModernSidebar es un control contenedor genérico para barras laterales (sidebars),
// diseñado bajo la filosofía común de ModernButton, ModernCard, ModernBadge y ModernTextBox:
// - Extiende widget.BaseWidget para integrarse como ciudadano nativo en el árbol de Fyne.
// - Ancho fijo configurable (por defecto 320 px) que impide redimensionamientos indeseados.
// - Fondo estilizado (#ffffff por defecto) y línea divisoria perimetral a la derecha (#e2e8f0).
// - Distribución modular tripartita: Cabecera (Header), Cuerpo (Content) y Pie (Footer).
// - Interfaz fluida (Fluent API) y constructores semánticos estilo xBase / Harbour.
type ModernSidebar struct {
	widget.BaseWidget

	width           float32
	bgColor         color.Color
	borderColor     color.Color
	borderWidth     float32
	showRightBorder bool
	isCustom        bool
	style           SidebarStyle

	header  fyne.CanvasObject
	content fyne.CanvasObject
	footer  fyne.CanvasObject

	widthGuide  *canvas.Rectangle
	bg          *canvas.Rectangle
	rightBorder *canvas.Rectangle
	rootStack   *fyne.Container
}

// Alias de tipo para máxima claridad arquitectónica
type SidebarContainer = ModernSidebar

// NewSidebarContainer es el constructor unificado estilo Harbour / xBase:
// Permite definir la barra lateral completa en una sola llamada.
// - width: anchura fija en píxeles (si es <= 0, usa 320)
// - header: widget para la parte superior (búsqueda, filtros, etc.) o nil
// - content: widget principal del cuerpo (lista, árbol, menú) o nil
// - footer: widget para la base (estado, botones, paginación) o nil
func NewSidebarContainer(width float32, header, content, footer fyne.CanvasObject) *ModernSidebar {
	if width <= 0 {
		width = 320
	}

	sb := &ModernSidebar{
		width:           width,
		bgColor:         ThemeColor(ColorNameSidebarBg),
		borderColor:     ThemeColor(ColorNameBorder),
		borderWidth:     1,
		showRightBorder: true,
		header:          header,
		content:         content,
		footer:          footer,
	}

	sb.widthGuide = canvas.NewRectangle(color.Transparent)
	sb.widthGuide.SetMinSize(fyne.NewSize(sb.width, 1))

	sb.bg = canvas.NewRectangle(sb.bgColor)

	sb.rightBorder = canvas.NewRectangle(sb.borderColor)
	sb.rightBorder.SetMinSize(fyne.NewSize(sb.borderWidth, 1))

	sb.rebuild()
	sb.ExtendBaseWidget(sb)
	return sb
}

// NewModernSidebar es el constructor semántico por defecto con cuerpo principal.
func NewModernSidebar(width float32, content fyne.CanvasObject) *ModernSidebar {
	return NewSidebarContainer(width, nil, content, nil)
}

// NewModernSidebarWithHeader crea un sidebar con cabecera y cuerpo.
func NewModernSidebarWithHeader(width float32, header, content fyne.CanvasObject) *ModernSidebar {
	return NewSidebarContainer(width, header, content, nil)
}

type modernSidebarRenderer struct {
	sidebar *ModernSidebar
}

func (s *ModernSidebar) CreateRenderer() fyne.WidgetRenderer {
	return &modernSidebarRenderer{sidebar: s}
}

func (r *modernSidebarRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.sidebar.rootStack}
}

func (r *modernSidebarRenderer) Destroy() {}

func (r *modernSidebarRenderer) Layout(size fyne.Size) {
	r.sidebar.rootStack.Resize(size)
	r.sidebar.rootStack.Move(fyne.NewPos(0, 0))
}

func (r *modernSidebarRenderer) MinSize() fyne.Size {
	return r.sidebar.MinSize()
}

func (r *modernSidebarRenderer) Refresh() {
	r.sidebar.applyStyle()
	r.sidebar.rootStack.Refresh()
}

// MinSize asegura que el sidebar reporte siempre con precisión su ancho fijo asignado.
func (s *ModernSidebar) MinSize() fyne.Size {
	minH := float32(0)
	if s.header != nil {
		minH += s.header.MinSize().Height
	}
	if s.content != nil {
		minH += s.content.MinSize().Height
	}
	if s.footer != nil {
		minH += s.footer.MinSize().Height
	}
	return fyne.NewSize(s.width, minH)
}

func (s *ModernSidebar) applyStyle() {
	if s.isCustom {
		s.bg.FillColor = s.bgColor
		s.rightBorder.FillColor = s.borderColor
	} else {
		if s.style == SidebarStyleShaded {
			s.bg.FillColor = ThemeColor(ColorNameSidebarShadedBg)
		} else {
			s.bg.FillColor = ThemeColor(ColorNameSidebarBg)
		}
		s.rightBorder.FillColor = ThemeColor(ColorNameBorder)
	}
	s.bg.Refresh()
	s.rightBorder.Refresh()
}

func (s *ModernSidebar) rebuild() {
	s.widthGuide.SetMinSize(fyne.NewSize(s.width, 1))
	s.applyStyle()
	s.rightBorder.SetMinSize(fyne.NewSize(s.borderWidth, 1))

	var rightObj fyne.CanvasObject
	if s.showRightBorder {
		rightObj = s.rightBorder
	}

	var centerObj fyne.CanvasObject = s.content
	if centerObj == nil {
		centerObj = container.NewVBox()
	}

	innerBorder := container.NewBorder(
		s.header,
		s.footer,
		nil,
		rightObj,
		centerObj,
	)

	newObjects := []fyne.CanvasObject{
		s.widthGuide,
		s.bg,
		innerBorder,
	}

	if s.rootStack == nil {
		s.rootStack = container.NewStack(newObjects...)
	} else {
		s.rootStack.Objects = newObjects
		s.rootStack.Refresh()
		s.Refresh()
	}
}

// --- Métodos de Contenido ---

func (s *ModernSidebar) SetHeader(header fyne.CanvasObject) {
	s.header = header
	s.rebuild()
}

func (s *ModernSidebar) SetContent(content fyne.CanvasObject) {
	s.content = content
	s.rebuild()
}

func (s *ModernSidebar) SetFooter(footer fyne.CanvasObject) {
	s.footer = footer
	s.rebuild()
}

func (s *ModernSidebar) SetWidth(width float32) {
	if width > 0 {
		s.width = width
		s.rebuild()
	}
}

// --- Interfaz Fluida (Fluent API) ---

func (s *ModernSidebar) WithWidth(width float32) *ModernSidebar {
	s.SetWidth(width)
	return s
}

func (s *ModernSidebar) WithBackground(bgHex string) *ModernSidebar {
	s.bgColor = HexColor(bgHex)
	s.isCustom = true
	s.rebuild()
	return s
}

func (s *ModernSidebar) WithRightBorder(borderHex string, width float32) *ModernSidebar {
	s.borderColor = HexColor(borderHex)
	s.isCustom = true
	if width > 0 {
		s.borderWidth = width
		s.showRightBorder = true
	} else {
		s.showRightBorder = false
	}
	s.rebuild()
	return s
}

func (s *ModernSidebar) WithHeader(header fyne.CanvasObject) *ModernSidebar {
	s.SetHeader(header)
	return s
}

func (s *ModernSidebar) WithContent(content fyne.CanvasObject) *ModernSidebar {
	s.SetContent(content)
	return s
}

func (s *ModernSidebar) WithFooter(footer fyne.CanvasObject) *ModernSidebar {
	s.SetFooter(footer)
	return s
}

// WithShadedBackground configura el sidebar con fondo ligeramente sombreado (#f8fafc en claro / #1e293b en oscuro)
// y borde divisor reactivo al tema, ideal para barras de accesos rápidos y navegación secundaria.
func (s *ModernSidebar) WithShadedBackground() *ModernSidebar {
	s.style = SidebarStyleShaded
	s.isCustom = false
	s.rebuild()
	return s
}
