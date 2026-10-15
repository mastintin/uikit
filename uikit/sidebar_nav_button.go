package uikit

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	_ fyne.Widget        = (*ModernSidebarNavButton)(nil)
	_ fyne.Tappable      = (*ModernSidebarNavButton)(nil)
	_ desktop.Hoverable  = (*ModernSidebarNavButton)(nil)
	_ desktop.Cursorable = (*ModernSidebarNavButton)(nil)
)

// ModernSidebarNavButton implementa un botón de navegación lateral con el estilo elegante y ejecutivo de gofyneJubila:
// - Barra lateral de acento de 3px a la izquierda cuando está activo
// - Icono vectorial semántico alineado
// - Tipografía adaptable al tema (negrita y color de acento cuando está activo, gris/muted cuando está inactivo)
// - Fondo suave de selección (theme.ColorNameSelection)
// - Cursor de mano interactivo y hover feedback
type ModernSidebarNavButton struct {
	widget.BaseWidget

	label     string
	iconName  string
	active    bool
	isHovered bool
	onTap     func()

	bg    *canvas.Rectangle
	strip *canvas.Rectangle
	icon  *ModernIcon
	lbl   *canvas.Text
}

// NewSidebarNavButton crea un nuevo botón de navegación lateral estilo gofyneJubila.
func NewSidebarNavButton(label, iconName string, onTap func()) *ModernSidebarNavButton {
	b := &ModernSidebarNavButton{
		label:    label,
		iconName: iconName,
		onTap:    onTap,
	}

	b.bg = canvas.NewRectangle(color.Transparent)
	b.bg.CornerRadius = 6

	b.strip = canvas.NewRectangle(color.Transparent)
	b.strip.CornerRadius = 0

	colMuted := ThemeColor(ColorNameTextMuted)
	b.icon = NewColoredIcon(iconName, 17, ColorToHex(colMuted))
	b.lbl = canvas.NewText(label, colMuted)
	b.lbl.TextSize = 13.5

	b.applyStyle()
	b.ExtendBaseWidget(b)
	return b
}

// SetActive establece el estado seleccionado del botón.
func (b *ModernSidebarNavButton) SetActive(active bool) {
	b.active = active
	b.applyStyle()
	b.Refresh()
}

// IsActive retorna si el botón está actualmente activo.
func (b *ModernSidebarNavButton) IsActive() bool {
	return b.active
}

// SetLabel actualiza el texto del botón.
func (b *ModernSidebarNavButton) SetLabel(label string) {
	b.label = label
	b.lbl.Text = label
	b.Refresh()
}

// SetIconName actualiza el icono del botón.
func (b *ModernSidebarNavButton) SetIconName(iconName string) {
	b.iconName = iconName
	b.icon.SetIcon(iconName)
	b.applyStyle()
	b.Refresh()
}

func (b *ModernSidebarNavButton) applyStyle() {
	isDark := IsDarkTheme()

	if b.active {
		b.bg.CornerRadius = 6
		b.lbl.TextStyle = fyne.TextStyle{Bold: true}

		b.bg.FillColor = ThemeColor(theme.ColorNameSelection)
		b.strip.FillColor = ThemeColor(ColorNameAccent)
		b.icon.SetColor(ThemeColor(ColorNameAccent))
		if isDark {
			b.lbl.Color = ThemeColor(ColorNameTextPrimary)
		} else {
			b.lbl.Color = ThemeColor(ColorNameAccent)
		}
	} else if b.isHovered {
		b.bg.CornerRadius = 6
		b.lbl.TextStyle = fyne.TextStyle{}

		if isDark {
			b.bg.FillColor = HexColor("#1e293b")
		} else {
			b.bg.FillColor = HexColor("#f1f5f9")
		}
		b.strip.FillColor = color.Transparent
		b.icon.SetColor(ThemeColor(ColorNameTextPrimary))
		b.lbl.Color = ThemeColor(ColorNameTextPrimary)
	} else {
		b.bg.FillColor = color.Transparent
		b.bg.CornerRadius = 6
		b.strip.FillColor = color.Transparent
		b.lbl.TextStyle = fyne.TextStyle{}

		colMuted := ThemeColor(ColorNameTextMuted)
		b.icon.SetColor(colMuted)
		b.lbl.Color = colMuted
	}

	b.bg.Refresh()
	b.strip.Refresh()
	b.lbl.Refresh()
}

func (b *ModernSidebarNavButton) Tapped(_ *fyne.PointEvent) {
	if b.onTap != nil {
		b.onTap()
	}
}

func (b *ModernSidebarNavButton) TappedSecondary(_ *fyne.PointEvent) {}

func (b *ModernSidebarNavButton) MouseIn(_ *desktop.MouseEvent) {
	b.isHovered = true
	b.applyStyle()
	b.Refresh()
}

func (b *ModernSidebarNavButton) MouseMoved(_ *desktop.MouseEvent) {}

func (b *ModernSidebarNavButton) MouseOut() {
	b.isHovered = false
	b.applyStyle()
	b.Refresh()
}

func (b *ModernSidebarNavButton) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (b *ModernSidebarNavButton) CreateRenderer() fyne.WidgetRenderer {
	return &modernSidebarNavButtonRenderer{btn: b}
}

type modernSidebarNavButtonRenderer struct {
	btn *ModernSidebarNavButton
}

func (r *modernSidebarNavButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.btn.bg, r.btn.strip, r.btn.icon, r.btn.lbl}
}

func (r *modernSidebarNavButtonRenderer) Destroy() {}

func (r *modernSidebarNavButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(160, 42)
}

func (r *modernSidebarNavButtonRenderer) Layout(size fyne.Size) {
	r.btn.bg.Resize(size)
	r.btn.bg.Move(fyne.NewPos(0, 0))

	// Barra lateral de acento (3px) a la izquierda
	r.btn.strip.Resize(fyne.NewSize(3, size.Height))
	r.btn.strip.Move(fyne.NewPos(0, 0))

	// Icono centrado verticalmente
	iconSize := float32(18)
	iconX := float32(16)
	iconY := (size.Height - iconSize) / 2
	r.btn.icon.Resize(fyne.NewSize(iconSize, iconSize))
	r.btn.icon.Move(fyne.NewPos(iconX, iconY))

	// Texto centrado verticalmente
	lblH := r.btn.lbl.MinSize().Height
	lblY := (size.Height - lblH) / 2
	r.btn.lbl.Resize(fyne.NewSize(size.Width-iconX-iconSize-10, lblH))
	r.btn.lbl.Move(fyne.NewPos(iconX+iconSize+10, lblY))
}

func (r *modernSidebarNavButtonRenderer) Refresh() {
	r.btn.applyStyle()
	canvas.Refresh(r.btn)
}
