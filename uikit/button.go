package uikit

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ModernButtonVariant int

const (
	ButtonDefault ModernButtonVariant = iota // .btn (blanco con borde fino gris y texto pizarra)
	ButtonPrimary                            // .btn-primary (azul suave translúcido #eff6ff con borde #bfdbfe y texto azul)
	ButtonDanger                             // .btn-danger (rojo suave pastel #fef2f2 con borde y texto rojo)
)

type ModernButtonSize int

const (
	ButtonSizeRegular   ModernButtonSize = iota // 32px altura, fuente 13px
	ButtonSizeSmall                             // 26px altura, fuente 12px
	ButtonSizeIcon                              // 32x32px cuadrado (solo icono)
	ButtonSizeSmallIcon                         // 24x24px cuadrado (solo icono compacto)
)

// ResourceSlidersHorizontal es un recurso vectorial idéntico a SlidersHorizontal de Lucide en goCuadros.
var ResourceSlidersHorizontal = fyne.NewStaticResource("sliders_horizontal.svg", []byte(`
<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
  <rect x="3" y="4" width="7" height="2" rx="1" fill="#ffffff"/>
  <rect x="14" y="4" width="7" height="2" rx="1" fill="#ffffff"/>
  <rect x="10" y="2" width="4" height="6" rx="2" fill="#ffffff"/>

  <rect x="3" y="11" width="11" height="2" rx="1" fill="#ffffff"/>
  <rect x="18" y="11" width="3" height="2" rx="1" fill="#ffffff"/>
  <rect x="14" y="9" width="4" height="6" rx="2" fill="#ffffff"/>

  <rect x="3" y="18" width="5" height="2" rx="1" fill="#ffffff"/>
  <rect x="12" y="18" width="9" height="2" rx="1" fill="#ffffff"/>
  <rect x="8" y="16" width="4" height="6" rx="2" fill="#ffffff"/>
</svg>
`))

type ModernButtonColors struct {
	NormalBg     color.Color
	NormalFg     color.Color
	NormalBorder color.Color
	HoverBg      color.Color
	HoverFg      color.Color
	HoverBorder  color.Color
	PressBg      color.Color
	PressFg      color.Color
	PressBorder  color.Color
}

// Control Global de Efecto Cristal (Glassmorphism Fyne 2.8)
var globalGlassButtons = false

// SetGlobalGlassButtons activa o desactiva el efecto cristal para todos los ModernButton de la app.
func SetGlobalGlassButtons(active bool) {
	globalGlassButtons = active
	RefreshAllDatePickers()
}

// SetGlobalGlass es un alias conciso de SetGlobalGlassButtons.
func SetGlobalGlass(active bool) {
	globalGlassButtons = active
	RefreshAllDatePickers()
}

// IsGlobalGlassButtons indica si el efecto cristal global está activado.
func IsGlobalGlassButtons() bool {
	return globalGlassButtons
}

// IsGlobalGlass es un alias conciso de IsGlobalGlassButtons.
func IsGlobalGlass() bool {
	return globalGlassButtons
}

// ModernButton es un botón moderno con bordes refinados, esquinas redondeadas (6px),
// transiciones de hover/activo y estados visuales idénticos a goCuadros.
type ModernButton struct {
	widget.DisableableWidget

	Text         string
	Icon         fyne.Resource
	Variant      ModernButtonVariant
	SizeKind     ModernButtonSize
	CustomColors *ModernButtonColors
	CustomHeight float32
	CustomWidth  float32
	CustomPadX   float32
	Borderless   bool

	// Control de efecto Glass (estilo FiveWin lGlass):
	// Si hasExplicitGlass es false, el botón hereda la variable global globalGlassButtons.
	// Si hasExplicitGlass es true, se rige estrictamente por lGlass.
	lGlass           bool
	hasExplicitGlass bool

	OnTapped func()

	hovered bool
	pressed bool
	loading bool
}

// SetLoading activa o desactiva el estado de carga del botón.
// Cuando está en modo carga, muestra un ModernSpinner giratorio y bloquea eventos de clic.
func (b *ModernButton) SetLoading(loading bool) *ModernButton {
	if b.loading == loading {
		return b
	}
	b.loading = loading
	b.Refresh()
	return b
}

// WithLoading configura fluidamente el estado de carga del botón.
func (b *ModernButton) WithLoading(loading bool) *ModernButton {
	return b.SetLoading(loading)
}

// IsLoading devuelve true si el botón se encuentra actualmente en estado de carga.
func (b *ModernButton) IsLoading() bool {
	return b.loading
}

func (b *ModernButton) SetBorderless(borderless bool) *ModernButton {
	b.Borderless = borderless
	b.Refresh()
	return b
}

func (b *ModernButton) WithBorderless() *ModernButton {
	return b.SetBorderless(true)
}

// SetGlass configura explícitamente la propiedad lGlass para este botón individual.
func (b *ModernButton) SetGlass(glass bool) *ModernButton {
	b.lGlass = glass
	b.hasExplicitGlass = true
	b.Refresh()
	return b
}

// WithGlass activa explícitamente el efecto cristal en el botón (Fluent API).
func (b *ModernButton) WithGlass() *ModernButton {
	return b.SetGlass(true)
}

// ResetGlass devuelve el botón a su comportamiento predeterminado, volviendo a heredar
// la configuración global de la aplicación.
func (b *ModernButton) ResetGlass() *ModernButton {
	b.hasExplicitGlass = false
	b.Refresh()
	return b
}

// IsGlass determina si el efecto cristal está actualmente activo para este botón
// (bien por sobreescritura local explícita o por herencia global).
func (b *ModernButton) IsGlass() bool {
	if b.hasExplicitGlass {
		return b.lGlass
	}
	return globalGlassButtons
}

func (b *ModernButton) SetCustomColors(c *ModernButtonColors) *ModernButton {
	b.CustomColors = c
	b.Refresh()
	return b
}

func (b *ModernButton) SetHeight(h float32) *ModernButton {
	b.CustomHeight = h
	b.Refresh()
	return b
}

func (b *ModernButton) SetMinWidth(w float32) *ModernButton {
	b.CustomWidth = w
	b.Refresh()
	return b
}

func (b *ModernButton) SetPaddingX(padX float32) *ModernButton {
	b.CustomPadX = padX
	b.Refresh()
	return b
}

func (b *ModernButton) SetMinSize(s fyne.Size) *ModernButton {
	b.CustomWidth = s.Width
	b.CustomHeight = s.Height
	b.Refresh()
	return b
}

var (
	_ fyne.Widget        = (*ModernButton)(nil)
	_ fyne.Tappable      = (*ModernButton)(nil)
	_ desktop.Hoverable  = (*ModernButton)(nil)
	_ desktop.Mouseable  = (*ModernButton)(nil)
	_ desktop.Cursorable = (*ModernButton)(nil)
)

func NewModernButton(text string, icon fyne.Resource, variant ModernButtonVariant, tapped func()) *ModernButton {
	b := &ModernButton{
		Text:     text,
		Icon:     icon,
		Variant:  variant,
		SizeKind: ButtonSizeRegular,
		OnTapped: tapped,
	}
	b.ExtendBaseWidget(b)
	return b
}

func NewModernButtonSmall(text string, icon fyne.Resource, variant ModernButtonVariant, tapped func()) *ModernButton {
	b := &ModernButton{
		Text:     text,
		Icon:     icon,
		Variant:  variant,
		SizeKind: ButtonSizeSmall,
		OnTapped: tapped,
	}
	b.ExtendBaseWidget(b)
	return b
}

func NewModernIconButton(icon fyne.Resource, variant ModernButtonVariant, tapped func()) *ModernButton {
	b := &ModernButton{
		Text:     "",
		Icon:     icon,
		Variant:  variant,
		SizeKind: ButtonSizeIcon,
		OnTapped: tapped,
	}
	b.ExtendBaseWidget(b)
	return b
}

// --- Constructores Semánticos Especializados ("Sugar Constructors") ---

// NewPrimaryButton crea un botón con estilo primario (.btn-primary azul translúcido suave con borde y texto azul)
func NewPrimaryButton(text string, icon fyne.Resource, tapped func()) *ModernButton {
	return NewModernButton(text, icon, ButtonPrimary, tapped)
}

// NewDangerButton crea un botón con estilo de peligro (.btn-danger rojo pastel #fef2f2 con texto #dc2626)
func NewDangerButton(text string, icon fyne.Resource, tapped func()) *ModernButton {
	return NewModernButton(text, icon, ButtonDanger, tapped)
}

// NewDefaultButton crea un botón estándar (.btn blanco con borde fino gris y texto pizarra)
func NewDefaultButton(text string, icon fyne.Resource, tapped func()) *ModernButton {
	return NewModernButton(text, icon, ButtonDefault, tapped)
}

// NewSmallDefaultButton crea un botón estándar compacto (26px de alto)
func NewSmallDefaultButton(text string, icon fyne.Resource, tapped func()) *ModernButton {
	return NewModernButtonSmall(text, icon, ButtonDefault, tapped)
}

// NewSmallPrimaryButton crea un botón primario compacto (26px de alto)
func NewSmallPrimaryButton(text string, icon fyne.Resource, tapped func()) *ModernButton {
	return NewModernButtonSmall(text, icon, ButtonPrimary, tapped)
}

// NewSmallDangerButton crea un botón de peligro compacto (26px de alto)
func NewSmallDangerButton(text string, icon fyne.Resource, tapped func()) *ModernButton {
	return NewModernButtonSmall(text, icon, ButtonDanger, tapped)
}

// NewIconButton crea un botón cuadrado estándar 32x32px solo con icono (.btn-icon)
func NewIconButton(icon fyne.Resource, tapped func()) *ModernButton {
	return NewModernIconButton(icon, ButtonDefault, tapped)
}

// NewSmallIconButton crea un botón cuadrado compacto 24x24px solo con icono
func NewSmallIconButton(icon fyne.Resource, tapped func()) *ModernButton {
	b := &ModernButton{
		Text:     "",
		Icon:     icon,
		Variant:  ButtonDefault,
		SizeKind: ButtonSizeSmallIcon,
		OnTapped: tapped,
	}
	b.ExtendBaseWidget(b)
	return b
}

// WithDisabled permite configurar el estado deshabilitado en cadena
func (b *ModernButton) WithDisabled(disabled bool) *ModernButton {
	if disabled {
		b.Disable()
	} else {
		b.Enable()
	}
	return b
}

// NewButton crea un botón con sintaxis unificada y directa tipo Harbour/xBase:
// btn := ui.NewButton("Primary", "Guardar", "save", lDisable, onGuardar)
// - variant: "primary", "danger", "default", "small", "small-primary", "small-danger", "icon", "small-icon"
// - text: texto visible (o "" para botones solo icono)
// - icon: fyne.Resource o nombre string ("save", "delete", "new", "edit", "print", "cancel", "preview", "folder", "config", "refresh", "image", "upload", "palette", "filter", etc.)
// - disabled: booleano de estado inicial (lDisable)
// - onClick: función callback func() ejecutada al pulsar
func NewButton(variant string, text string, icon any, disabled bool, onClick func()) *ModernButton {
	v := ButtonDefault
	size := ButtonSizeRegular

	switch strings.ToLower(strings.TrimSpace(variant)) {
	case "primary", "btn-primary":
		v = ButtonPrimary
	case "danger", "btn-danger":
		v = ButtonDanger
	case "small", "sm", "btn-sm", "small-default":
		size = ButtonSizeSmall
	case "small-primary", "sm-primary":
		v = ButtonPrimary
		size = ButtonSizeSmall
	case "small-danger", "sm-danger":
		v = ButtonDanger
		size = ButtonSizeSmall
	case "icon", "btn-icon":
		size = ButtonSizeIcon
	case "small-icon", "sm-icon":
		size = ButtonSizeSmallIcon
	default:
		v = ButtonDefault
	}

	var res fyne.Resource
	switch val := icon.(type) {
	case fyne.Resource:
		res = val
	case string:
		res = resolveIconByName(val)
	}

	btn := &ModernButton{
		Text:     text,
		Icon:     res,
		Variant:  v,
		SizeKind: size,
		OnTapped: onClick,
	}
	btn.ExtendBaseWidget(btn)

	if disabled {
		btn.Disable()
	}

	return btn
}

// resolveIconByName mapea nombres comunes a recursos de icono oficiales de Fyne o vectoriales de Lucide
func resolveIconByName(name string) fyne.Resource {
	if res := getVectorIconByName(name); res != nil {
		return res
	}
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "save", "grabar", "guardar":
		return theme.DocumentSaveIcon()
	case "delete", "borrar", "eliminar":
		return theme.DeleteIcon()
	case "new", "nuevo", "add", "crear":
		return theme.ContentAddIcon()
	case "edit", "editar", "modificar":
		return theme.DocumentCreateIcon()
	case "print", "imprimir":
		return theme.DocumentPrintIcon()
	case "cancel", "cancelar", "close":
		return theme.CancelIcon()
	case "preview", "vista-previa", "ojo":
		return theme.VisibilityIcon()
	case "folder", "finder", "abrir":
		return theme.FolderOpenIcon()
	case "config", "settings", "ajustes":
		return theme.SettingsIcon()
	case "refresh", "refrescar", "limpiar":
		return theme.ViewRefreshIcon()
	case "image", "foto", "imagen":
		return theme.FileImageIcon()
	case "upload", "subir":
		return theme.UploadIcon()
	case "palette", "photoshop", "color":
		return theme.ColorPaletteIcon()
	case "filter", "filtros":
		return ResourceSlidersHorizontal
	case "export", "exportar", "download", "descargar":
		return theme.DownloadIcon()
	case "storage", "db", "database":
		return theme.StorageIcon()
	case "check", "seleccionar", "aceptar", "confirm":
		return theme.ConfirmIcon()
	case "up", "subir-nivel", "arriba":
		return theme.MoveUpIcon()
	case "home", "inicio":
		return theme.HomeIcon()
	case "desktop", "escritorio":
		return theme.DesktopIcon()
	case "document", "documento", "documentos":
		return theme.DocumentIcon()
	case "new-folder", "carpeta-nueva":
		return theme.FolderNewIcon()
	default:
		return nil
	}
}



func (b *ModernButton) SetText(text string) {
	b.Text = text
	b.Refresh()
}

func (b *ModernButton) SetIcon(icon fyne.Resource) {
	b.Icon = icon
	b.Refresh()
}

func (b *ModernButton) SetIconByName(name string) *ModernButton {
	b.Icon = IconResource(name)
	b.Refresh()
	return b
}

func (b *ModernButton) SetVariant(v ModernButtonVariant) {
	b.Variant = v
	b.Refresh()
}

func (b *ModernButton) Tapped(_ *fyne.PointEvent) {
	if b.Disabled() || b.loading {
		return
	}
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

func (b *ModernButton) MouseIn(_ *desktop.MouseEvent) {
	if b.Disabled() || b.loading {
		return
	}
	b.hovered = true
	b.Refresh()
}

func (b *ModernButton) MouseMoved(_ *desktop.MouseEvent) {}

func (b *ModernButton) MouseOut() {
	if b.Disabled() || b.loading {
		return
	}
	b.hovered = false
	b.pressed = false
	b.Refresh()
}

func (b *ModernButton) MouseDown(_ *desktop.MouseEvent) {
	if b.Disabled() || b.loading {
		return
	}
	b.pressed = true
	b.Refresh()
}

func (b *ModernButton) MouseUp(_ *desktop.MouseEvent) {
	if b.Disabled() || b.loading {
		return
	}
	b.pressed = false
	b.Refresh()
}

func (b *ModernButton) Cursor() desktop.Cursor {
	if b.Disabled() || b.loading {
		return desktop.DefaultCursor
	}
	return desktop.PointerCursor
}

type modernButtonRenderer struct {
	button  *ModernButton
	blur    *canvas.Blur
	bg      *canvas.Rectangle
	icon    *canvas.Image
	label   *canvas.Text
	spinner *ModernSpinner
}

func (b *ModernButton) CreateRenderer() fyne.WidgetRenderer {
	blur := canvas.NewBlur(10)
	blur.CornerRadius = 6

	bg := canvas.NewRectangle(color.Transparent)
	bg.CornerRadius = 6
	bg.StrokeWidth = 1

	icon := canvas.NewImageFromResource(nil)
	icon.FillMode = canvas.ImageFillContain

	label := canvas.NewText("", ThemeColor(ColorNameTextPrimary))
	label.TextStyle = fyne.TextStyle{Bold: true}
	label.Alignment = fyne.TextAlignCenter

	spinner := NewSpinnerSmall()
	spinner.Hide()

	r := &modernButtonRenderer{
		button:  b,
		blur:    blur,
		bg:      bg,
		icon:    icon,
		label:   label,
		spinner: spinner,
	}
	r.applyStyle()
	return r
}

func (r *modernButtonRenderer) Objects() []fyne.CanvasObject {
	if r.button.IsGlass() {
		return []fyne.CanvasObject{r.blur, r.bg, r.icon, r.label, r.spinner}
	}
	return []fyne.CanvasObject{r.bg, r.icon, r.label, r.spinner}
}

func (r *modernButtonRenderer) Destroy() {
	r.spinner.Stop()
}

func (r *modernButtonRenderer) MinSize() fyne.Size {
	b := r.button
	var minH float32 = 32
	var padX float32 = 24 // 12px a cada lado
	var iconSize float32 = 15
	var gap float32 = 6

	if b.CustomPadX > 0 {
		padX = b.CustomPadX
	}

	if b.CustomHeight > 0 {
		minH = b.CustomHeight
	} else if b.SizeKind == ButtonSizeSmall {
		minH = 26
		if b.CustomPadX <= 0 {
			padX = 16 // 8px a cada lado
		}
		iconSize = 13
		gap = 5
	} else if b.SizeKind == ButtonSizeIcon {
		h := float32(32)
		if b.CustomHeight > 0 {
			h = b.CustomHeight
		}
		w := float32(32)
		if b.CustomWidth > 0 {
			w = b.CustomWidth
		}
		return fyne.NewSize(w, h)
	} else if b.SizeKind == ButtonSizeSmallIcon {
		h := float32(24)
		if b.CustomHeight > 0 {
			h = b.CustomHeight
		}
		w := float32(24)
		if b.CustomWidth > 0 {
			w = b.CustomWidth
		}
		return fyne.NewSize(w, h)
	}

	if b.CustomHeight >= 42 {
		iconSize = 17
		gap = 8
		if b.CustomPadX <= 0 {
			padX = 28
		}
	} else if b.CustomHeight >= 36 {
		iconSize = 16
		gap = 7
	}

	isLoading := b.loading
	hasVisual := (b.Icon != nil && !isLoading) || isLoading
	hasText := b.Text != ""

	if !hasText && hasVisual {
		w := minH
		if b.CustomWidth > 0 {
			w = b.CustomWidth
		}
		return fyne.NewSize(w, minH)
	}

	var minW float32 = padX
	if hasText {
		minW += r.label.MinSize().Width
	}
	if hasVisual {
		minW += iconSize
		if hasText {
			minW += gap
		}
	}

	if b.CustomWidth > 0 && minW < b.CustomWidth {
		minW = b.CustomWidth
	}

	return fyne.NewSize(minW, minH)
}

func (r *modernButtonRenderer) Layout(size fyne.Size) {
	if r.button.IsGlass() {
		r.blur.Resize(size)
		r.blur.Move(fyne.NewPos(0, 0))
		r.blur.Show()
	} else {
		r.blur.Hide()
	}
	r.bg.Resize(size)
	r.bg.Move(fyne.NewPos(0, 0))

	b := r.button
	isLoading := b.loading
	hasVisual := (b.Icon != nil && !isLoading) || isLoading
	hasText := b.Text != ""

	var iconSize float32 = 15
	var gap float32 = 6
	if b.CustomHeight >= 42 {
		iconSize = 17
		gap = 8
	} else if b.CustomHeight >= 36 {
		iconSize = 16
		gap = 7
	} else if b.SizeKind == ButtonSizeSmall {
		iconSize = 13
		gap = 5
	} else if b.SizeKind == ButtonSizeIcon {
		iconSize = 16
	} else if b.SizeKind == ButtonSizeSmallIcon {
		iconSize = 14
	}

	if hasVisual && hasText {
		labelSize := r.label.MinSize()
		totalW := iconSize + gap + labelSize.Width
		startX := (size.Width - totalW) / 2
		if startX < 0 {
			startX = 0
		}

		visualY := (size.Height - iconSize) / 2
		if isLoading {
			r.icon.Hide()
			r.spinner.Resize(fyne.NewSize(iconSize, iconSize))
			r.spinner.Move(fyne.NewPos(startX, visualY))
			r.spinner.Show()
		} else {
			r.spinner.Hide()
			r.icon.Resize(fyne.NewSize(iconSize, iconSize))
			r.icon.Move(fyne.NewPos(startX, visualY))
			r.icon.Show()
		}

		labelY := (size.Height - labelSize.Height) / 2
		r.label.Resize(labelSize)
		r.label.Move(fyne.NewPos(startX+iconSize+gap, labelY))
		r.label.Show()
	} else if hasVisual && !hasText {
		visualY := (size.Height - iconSize) / 2
		visualX := (size.Width - iconSize) / 2
		if isLoading {
			r.icon.Hide()
			r.spinner.Resize(fyne.NewSize(iconSize, iconSize))
			r.spinner.Move(fyne.NewPos(visualX, visualY))
			r.spinner.Show()
		} else {
			r.spinner.Hide()
			r.icon.Resize(fyne.NewSize(iconSize, iconSize))
			r.icon.Move(fyne.NewPos(visualX, visualY))
			r.icon.Show()
		}
		r.label.Resize(fyne.NewSize(0, 0))
		r.label.Move(fyne.NewPos(0, 0))
		r.label.Hide()
	} else if !hasVisual && hasText {
		r.icon.Hide()
		r.spinner.Hide()
		labelSize := r.label.MinSize()
		labelX := (size.Width - labelSize.Width) / 2
		labelY := (size.Height - labelSize.Height) / 2
		r.label.Resize(labelSize)
		r.label.Move(fyne.NewPos(labelX, labelY))
		r.label.Show()
	}
}

func (r *modernButtonRenderer) Refresh() {
	r.applyStyle()
	r.Layout(r.button.Size())
	canvas.Refresh(r.button)
}

func (r *modernButtonRenderer) applyStyle() {
	b := r.button

	// 1. Text Size
	if b.SizeKind == ButtonSizeSmall {
		r.label.TextSize = 12.5
	} else if b.CustomHeight >= 42 {
		r.label.TextSize = 14.5
	} else {
		r.label.TextSize = 13.5
	}
	r.label.Text = b.Text

	// 2. Colores y estado interactivo
	var themeColor fyne.ThemeColorName
	_ = themeColor

	if b.IsGlass() {
		r.blur.Show()
		r.blur.CornerRadius = r.bg.CornerRadius
		r.bg.Shadow = canvas.Shadow{
			Color:      color.NRGBA{R: 0, G: 0, B: 0, A: 20},
			BlurRadius: 4,
			Spread:     0,
			Offset:     fyne.NewPos(0, 1.5),
			Variant:    canvas.DropShadow,
		}

		if b.CustomColors != nil {
			themeColor = theme.ColorNameForeground
			if b.Disabled() {
				r.label.Color = ThemeColor(ColorNameBtnMuted)
				r.bg.FillColor = color.Transparent
				r.bg.StrokeColor = color.Transparent
			} else if b.pressed {
				if b.CustomColors.PressBg != nil { r.bg.FillColor = b.CustomColors.PressBg }
				if b.CustomColors.PressFg != nil { r.label.Color = b.CustomColors.PressFg }
				if b.CustomColors.PressBorder != nil { r.bg.StrokeColor = b.CustomColors.PressBorder }
			} else if b.hovered {
				if b.CustomColors.HoverBg != nil { r.bg.FillColor = b.CustomColors.HoverBg }
				if b.CustomColors.HoverFg != nil { r.label.Color = b.CustomColors.HoverFg }
				if b.CustomColors.HoverBorder != nil { r.bg.StrokeColor = b.CustomColors.HoverBorder }
			} else {
				if b.CustomColors.NormalBg != nil { r.bg.FillColor = b.CustomColors.NormalBg }
				if b.CustomColors.NormalFg != nil { r.label.Color = b.CustomColors.NormalFg }
				if b.CustomColors.NormalBorder != nil { r.bg.StrokeColor = b.CustomColors.NormalBorder }
			}
		} else {
			pal := GetCurrentThemePalette()
			var btnColors ButtonStateColors
			if pal != nil {
				switch b.Variant {
				case ButtonPrimary:
					btnColors = pal.GlassButtons.Primary
				case ButtonDanger:
					btnColors = pal.GlassButtons.Danger
				default:
					btnColors = pal.GlassButtons.Default
				}
			}

			if b.Disabled() {
				themeColor = ColorNameBtnMuted
				if b.Variant == ButtonPrimary {
					themeColor = ColorNameBtnPrimaryMuted
				} else if b.Variant == ButtonDanger {
					themeColor = ColorNameBtnDangerMuted
				}
				if btnColors.DisabledBg != nil {
					r.bg.FillColor = btnColors.DisabledBg
					r.bg.StrokeColor = btnColors.DisabledBorder
					r.label.Color = btnColors.DisabledFg
				} else {
					r.bg.FillColor = color.Transparent
					r.bg.StrokeColor = color.Transparent
					r.label.Color = ThemeColor(themeColor)
				}
			} else if b.pressed {
				if b.Variant == ButtonPrimary {
					themeColor = ColorNameBtnPrimaryPressFg
				} else if b.Variant == ButtonDanger {
					themeColor = ColorNameBtnDangerPressFg
				} else {
					themeColor = ColorNameBtnDefaultPressFg
				}
				if btnColors.PressBg != nil {
					r.bg.FillColor = btnColors.PressBg
					r.bg.StrokeColor = btnColors.PressBorder
					r.label.Color = btnColors.PressFg
				}
			} else if b.hovered {
				if b.Variant == ButtonPrimary {
					themeColor = ColorNameBtnPrimaryHoverFg
				} else if b.Variant == ButtonDanger {
					themeColor = ColorNameBtnDangerHoverFg
				} else {
					themeColor = ColorNameBtnDefaultHoverFg
				}
				if btnColors.HoverBg != nil {
					r.bg.FillColor = btnColors.HoverBg
					r.bg.StrokeColor = btnColors.HoverBorder
					r.label.Color = btnColors.HoverFg
				}
			} else {
				if b.Variant == ButtonPrimary {
					themeColor = ColorNameBtnPrimaryFg
				} else if b.Variant == ButtonDanger {
					themeColor = ColorNameBtnDangerFg
				} else {
					themeColor = ColorNameBtnDefaultFg
				}
				if btnColors.Bg != nil {
					r.bg.FillColor = btnColors.Bg
					r.bg.StrokeColor = btnColors.Border
					r.label.Color = btnColors.Fg
				}
			}
		}
	} else {
		r.blur.Hide()
		r.bg.Shadow = canvas.Shadow{}

		if b.CustomColors != nil {
			themeColor = theme.ColorNameForeground
			if b.Disabled() {
				r.label.Color = ThemeColor(ColorNameBtnMuted)
				r.bg.FillColor = color.Transparent
				r.bg.StrokeColor = color.Transparent
			} else if b.pressed {
				if b.CustomColors.PressBg != nil { r.bg.FillColor = b.CustomColors.PressBg }
				if b.CustomColors.PressFg != nil { r.label.Color = b.CustomColors.PressFg }
				if b.CustomColors.PressBorder != nil { r.bg.StrokeColor = b.CustomColors.PressBorder }
			} else if b.hovered {
				if b.CustomColors.HoverBg != nil { r.bg.FillColor = b.CustomColors.HoverBg }
				if b.CustomColors.HoverFg != nil { r.label.Color = b.CustomColors.HoverFg }
				if b.CustomColors.HoverBorder != nil { r.bg.StrokeColor = b.CustomColors.HoverBorder }
			} else {
				if b.CustomColors.NormalBg != nil { r.bg.FillColor = b.CustomColors.NormalBg }
				if b.CustomColors.NormalFg != nil { r.label.Color = b.CustomColors.NormalFg }
				if b.CustomColors.NormalBorder != nil { r.bg.StrokeColor = b.CustomColors.NormalBorder }
			}
		} else {
			switch b.Variant {
			case ButtonPrimary:
				if b.Disabled() {
					themeColor = ColorNameBtnPrimaryMuted
					r.label.Color = ThemeColor(ColorNameBtnPrimaryMuted)
					r.bg.FillColor = ThemeColor(ColorNameBtnPrimaryBg)
					r.bg.StrokeColor = ThemeColor(ColorNameBtnPrimaryBorder)
				} else if b.pressed {
					themeColor = ColorNameBtnPrimaryPressFg
					r.bg.FillColor = ThemeColor(ColorNameBtnPrimaryPressBg)
					r.label.Color = ThemeColor(ColorNameBtnPrimaryPressFg)
					r.bg.StrokeColor = ThemeColor(ColorNameBtnPrimaryPressBorder)
				} else if b.hovered {
					themeColor = ColorNameBtnPrimaryHoverFg
					r.bg.FillColor = ThemeColor(ColorNameBtnPrimaryHoverBg)
					r.label.Color = ThemeColor(ColorNameBtnPrimaryHoverFg)
					r.bg.StrokeColor = ThemeColor(ColorNameBtnPrimaryHoverBorder)
				} else {
					themeColor = ColorNameBtnPrimaryFg
					r.bg.FillColor = ThemeColor(ColorNameBtnPrimaryBg)
					r.label.Color = ThemeColor(ColorNameBtnPrimaryFg)
					r.bg.StrokeColor = ThemeColor(ColorNameBtnPrimaryBorder)
				}

			case ButtonDanger:
				if b.Disabled() {
					themeColor = ColorNameBtnDangerMuted
					r.label.Color = ThemeColor(ColorNameBtnDangerMuted)
					r.bg.FillColor = ThemeColor(ColorNameBtnDangerBg)
					r.bg.StrokeColor = ThemeColor(ColorNameBtnDangerBorder)
				} else if b.pressed {
					themeColor = ColorNameBtnDangerPressFg
					r.bg.FillColor = ThemeColor(ColorNameBtnDangerPressBg)
					r.label.Color = ThemeColor(ColorNameBtnDangerPressFg)
					r.bg.StrokeColor = ThemeColor(ColorNameBtnDangerPressBorder)
				} else if b.hovered {
					themeColor = ColorNameBtnDangerHoverFg
					r.bg.FillColor = ThemeColor(ColorNameBtnDangerHoverBg)
					r.label.Color = ThemeColor(ColorNameBtnDangerHoverFg)
					r.bg.StrokeColor = ThemeColor(ColorNameBtnDangerHoverBorder)
				} else {
					themeColor = ColorNameBtnDangerFg
					r.bg.FillColor = ThemeColor(ColorNameBtnDangerBg)
					r.label.Color = ThemeColor(ColorNameBtnDangerFg)
					r.bg.StrokeColor = ThemeColor(ColorNameBtnDangerBorder)
				}

			default: // ButtonDefault
				if b.Disabled() {
					themeColor = ColorNameBtnMuted
					r.label.Color = ThemeColor(ColorNameBtnMuted)
					r.bg.FillColor = ThemeColor(ColorNameBtnDefaultBg)
					r.bg.StrokeColor = ThemeColor(ColorNameBtnDefaultBorder)
				} else if b.pressed {
					themeColor = ColorNameBtnDefaultPressFg
					r.bg.FillColor = ThemeColor(ColorNameBtnDefaultPressBg)
					r.label.Color = ThemeColor(ColorNameBtnDefaultPressFg)
					r.bg.StrokeColor = ThemeColor(ColorNameBtnDefaultPressBorder)
				} else if b.hovered {
					themeColor = ColorNameBtnDefaultHoverFg
					r.bg.FillColor = ThemeColor(ColorNameBtnDefaultHoverBg)
					r.label.Color = ThemeColor(ColorNameBtnDefaultHoverFg)
					r.bg.StrokeColor = ThemeColor(ColorNameBtnDefaultHoverBorder)
				} else {
					themeColor = ColorNameBtnDefaultFg
					r.bg.FillColor = ThemeColor(ColorNameBtnDefaultBg)
					r.label.Color = ThemeColor(ColorNameBtnDefaultFg)
					r.bg.StrokeColor = ThemeColor(ColorNameBtnDefaultBorder)
				}
			}
		}
	}

	if b.Borderless {
		r.blur.Hide()
		r.bg.Shadow = canvas.Shadow{}
		r.bg.StrokeColor = color.Transparent
		if !b.hovered && !b.pressed {
			r.bg.FillColor = color.Transparent
		}
	}

	// 3. Icon
	if b.Icon != nil && !b.loading {
		content := string(b.Icon.Content())
		if strings.Contains(content, "<svg") && (strings.Contains(content, "currentColor") || strings.Contains(content, "stroke=")) {
			hex := ColorToHex(r.label.Color)
			r.icon.Resource = IconResource(b.Icon, hex)
		} else {
			standardColor := theme.ColorNameForeground
			if b.Variant == ButtonPrimary {
				standardColor = theme.ColorNamePrimary
			} else if b.Variant == ButtonDanger {
				standardColor = theme.ColorNameError
			}
			r.icon.Resource = theme.NewColoredResource(b.Icon, standardColor)
		}
		r.icon.Image = nil
		r.icon.Refresh()
		r.icon.Show()
	} else {
		r.icon.Resource = nil
		r.icon.Image = nil
		r.icon.Hide()
	}

	// 4. Spinner
	var spinnerSize float32 = 15
	if b.CustomHeight >= 42 {
		spinnerSize = 17
	} else if b.CustomHeight >= 36 {
		spinnerSize = 16
	} else if b.SizeKind == ButtonSizeSmall {
		spinnerSize = 13
	} else if b.SizeKind == ButtonSizeIcon {
		spinnerSize = 16
	} else if b.SizeKind == ButtonSizeSmallIcon {
		spinnerSize = 14
	}

	if b.loading {
		r.spinner.WithSize(spinnerSize)
		r.spinner.WithColor(r.label.Color)
		if !r.spinner.IsRunning() {
			r.spinner.Start()
		}
		r.spinner.Show()
	} else {
		if r.spinner.IsRunning() {
			r.spinner.Stop()
		}
		r.spinner.Hide()
	}

	if b.Text != "" {
		r.label.Show()
	} else {
		r.label.Hide()
	}

	r.label.Refresh()
	r.bg.Refresh()
}
