package uikit

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ModernIcon implementa un control vectorial genérico de iconos, diseñado bajo
// la misma filosofía que ModernButton, ModernBadge, ModernCard, ModernTextBox, ModernSidebar y ModernListBox:
// - Extiende widget.BaseWidget como ciudadano nativo del árbol de Fyne.
// - Soporta iconos estándar vectoriales de alta definición (Lucide icons y Fyne theme icons).
// - Tamaño configurable con precisión en píxeles (sin iconos "enanos" ni deformaciones).
// - Colores personalizados vía código hexadecimal (ej. "#2563eb", "#dc2626", "#64748b").
// - Soporte para fondo circular o redondeado tipo badge / pastilla con borde perimetral.
// - Interacción opcional: Tappable (clic), Hoverable (cambio de color al pasar el cursor) y Cursorable (puntero de mano).
// - Constructores unificados estilo Harbour / xBase y constructores semánticos ("Sugar").
// - Interfaz fluida (Fluent API) encadenable.
type ModernIcon struct {
	widget.BaseWidget

	rawSource    any
	resource     fyne.Resource
	size         float32
	colorHex     string
	hoverHex     string
	themeColor   fyne.ThemeColorName

	// Marco de fondo opcional (badge / botón icono plano)
	hasBg        bool
	badgeSize    float32
	bgColor      color.Color
	bgHoverColor color.Color
	borderColor  color.Color
	borderWidth  float32
	cornerRadius float32

	// Interactividad
	onTap        func()
	isHovered    bool
	disabled     bool

	// Objetos gráficos internos
	img          *canvas.Image
	bgRect       *canvas.Rectangle
	root         *fyne.Container
}

// Asegurar que implementa interfaces interactivas de Fyne
var (
	_ fyne.Tappable      = (*ModernIcon)(nil)
	_ desktop.Hoverable  = (*ModernIcon)(nil)
	_ desktop.Cursorable = (*ModernIcon)(nil)
)

// NewIcon es el constructor unificado estilo Harbour / xBase:
// Permite definir cualquier icono en una sola llamada.
// - nameOrRes: nombre ("search", "close", "palette", "filter", etc.) o fyne.Resource
// - size: tamaño en píxeles (si es <= 0, usa 20 px por defecto)
// - colorHex: color hexadecimal (ej. "#2563eb", "#64748b") o "" para color por defecto
// - onClick: callback opcional al hacer clic, o nil para icono puramente estático
func NewIcon(nameOrRes any, size float32, colorHex string, onClick func()) *ModernIcon {
	if size <= 0 {
		size = 20
	}

	ic := &ModernIcon{
		rawSource: nameOrRes,
		size:      size,
		colorHex:  colorHex,
		onTap:     onClick,
	}

	ic.resource = resolveModernIcon(nameOrRes, colorHex)
	ic.img = canvas.NewImageFromResource(ic.resource)
	ic.img.FillMode = canvas.ImageFillContain
	ic.img.SetMinSize(fyne.NewSize(size, size))

	ic.rebuild()
	ic.ExtendBaseWidget(ic)
	return ic
}

// --- Constructores Semánticos ("Sugar Constructors") ---

// NewStandardIcon crea un icono estándar del tamaño especificado con color de primer plano.
func NewStandardIcon(nameOrRes any, size float32) *ModernIcon {
	return NewIcon(nameOrRes, size, "", nil)
}

// NewColoredIcon crea un icono con un color hexadecimal específico (ej. "#2563eb").
func NewColoredIcon(nameOrRes any, size float32, colorHex string) *ModernIcon {
	return NewIcon(nameOrRes, size, colorHex, nil)
}

// NewTappableIcon crea un icono interactivo que reacciona al clic y muestra cursor de mano.
func NewTappableIcon(nameOrRes any, size float32, onTap func()) *ModernIcon {
	return NewIcon(nameOrRes, size, "", onTap)
}

// NewActionIcon crea un icono de acción interactivo (ideal para botones 'X' de limpiar, cerrar, etc.)
// con transición de color normal a color hover.
func NewActionIcon(nameOrRes any, size float32, normalHex, hoverHex string, onTap func()) *ModernIcon {
	ic := NewIcon(nameOrRes, size, normalHex, onTap)
	ic.hoverHex = hoverHex
	return ic
}

// NewBadgeIcon crea un icono contenido en una pastilla redondeada o circular con fondo,
// borde perimetral e icono centrado (reemplaza badges manuales en alertas y cabeceras).
func NewBadgeIcon(nameOrRes any, iconSize, badgeSize float32, bgHex, borderHex, iconHex string) *ModernIcon {
	if badgeSize <= 0 {
		badgeSize = iconSize + 16
	}
	ic := NewIcon(nameOrRes, iconSize, iconHex, nil)
	ic.hasBg = true
	ic.badgeSize = badgeSize
	ic.bgColor = HexColor(bgHex)
	ic.borderColor = HexColor(borderHex)
	ic.borderWidth = 1
	ic.cornerRadius = badgeSize / 2 // circular por defecto
	ic.rebuild()
	return ic
}

// NewCircleIcon crea un icono con fondo circular suave y sin borde.
func NewCircleIcon(nameOrRes any, iconSize, circleSize float32, bgHex, iconHex string) *ModernIcon {
	return NewBadgeIcon(nameOrRes, iconSize, circleSize, bgHex, bgHex, iconHex)
}


// --- Métodos de ciclo de vida y renderizado ---

func (ic *ModernIcon) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(ic.root)
}

func (ic *ModernIcon) MinSize() fyne.Size {
	if ic.hasBg && ic.badgeSize > 0 {
		return fyne.NewSize(ic.badgeSize, ic.badgeSize)
	}
	return fyne.NewSize(ic.size, ic.size)
}

func (ic *ModernIcon) rebuild() {
	ic.img.SetMinSize(fyne.NewSize(ic.size, ic.size))
	ic.img.Resize(fyne.NewSize(ic.size, ic.size))

	var activeHex = ic.colorHex
	if ic.isHovered && ic.hoverHex != "" {
		activeHex = ic.hoverHex
	}
	ic.resource = resolveModernIcon(ic.rawSource, activeHex)
	ic.img.Resource = ic.resource
	ic.img.Refresh()

	if ic.hasBg {
		if ic.bgRect == nil {
			ic.bgRect = canvas.NewRectangle(ic.bgColor)
		}
		ic.bgRect.FillColor = ic.bgColor
		ic.bgRect.StrokeColor = ic.borderColor
		ic.bgRect.StrokeWidth = ic.borderWidth
		ic.bgRect.CornerRadius = ic.cornerRadius
		ic.bgRect.SetMinSize(fyne.NewSize(ic.badgeSize, ic.badgeSize))

		centerWrap := container.NewCenter(ic.img)
		ic.root = container.NewStack(ic.bgRect, centerWrap)
	} else {
		ic.root = container.NewCenter(ic.img)
	}

	ic.Refresh()
}

// --- Interfaz Fluida (Fluent API) ---

func (ic *ModernIcon) WithSize(size float32) *ModernIcon {
	ic.SetSize(size)
	return ic
}

func (ic *ModernIcon) WithColor(colorHex string) *ModernIcon {
	ic.SetHexColor(colorHex)
	return ic
}

func (ic *ModernIcon) WithHoverColor(hoverHex string) *ModernIcon {
	ic.hoverHex = hoverHex
	return ic
}

func (ic *ModernIcon) WithBackground(bgHex string, cornerRadius float32) *ModernIcon {
	ic.hasBg = true
	ic.bgColor = HexColor(bgHex)
	ic.cornerRadius = cornerRadius
	ic.rebuild()
	return ic
}

func (ic *ModernIcon) WithCircleBackground(bgHex string, badgeSize float32) *ModernIcon {
	ic.hasBg = true
	ic.badgeSize = badgeSize
	ic.bgColor = HexColor(bgHex)
	ic.cornerRadius = badgeSize / 2
	ic.rebuild()
	return ic
}

func (ic *ModernIcon) WithBorder(borderHex string, width float32) *ModernIcon {
	ic.borderColor = HexColor(borderHex)
	ic.borderWidth = width
	ic.rebuild()
	return ic
}

func (ic *ModernIcon) WithOnTap(onTap func()) *ModernIcon {
	ic.SetOnTap(onTap)
	return ic
}

func (ic *ModernIcon) WithDisabled(disabled bool) *ModernIcon {
	ic.disabled = disabled
	ic.Refresh()
	return ic
}

// --- Métodos de Alto Nivel ---

func (ic *ModernIcon) SetIcon(nameOrRes any) {
	ic.rawSource = nameOrRes
	ic.rebuild()
}

func (ic *ModernIcon) SetSize(size float32) {
	if size > 0 {
		ic.size = size
		ic.rebuild()
	}
}

func (ic *ModernIcon) SetHexColor(hex string) {
	ic.colorHex = hex
	ic.rebuild()
}

// SetColor permite asignar directamente un color.Color al icono vectorial.
func (ic *ModernIcon) SetColor(c color.Color) {
	if c == nil {
		ic.SetHexColor("")
		return
	}
	ic.SetHexColor(ColorToHex(c))
}


func (ic *ModernIcon) SetOnTap(onTap func()) {
	ic.onTap = onTap
}

func (ic *ModernIcon) Enable() {
	ic.disabled = false
	ic.Refresh()
}

func (ic *ModernIcon) Disable() {
	ic.disabled = true
	ic.Refresh()
}

// --- Eventos Interactivos ---

func (ic *ModernIcon) Tapped(_ *fyne.PointEvent) {
	if ic.disabled || ic.onTap == nil {
		return
	}
	ic.onTap()
}

func (ic *ModernIcon) MouseIn(_ *desktop.MouseEvent) {
	if ic.disabled || ic.onTap == nil {
		return
	}
	ic.isHovered = true
	if ic.hoverHex != "" {
		ic.resource = resolveModernIcon(ic.rawSource, ic.hoverHex)
		ic.img.Resource = ic.resource
		ic.img.Refresh()
	}
}

func (ic *ModernIcon) MouseMoved(_ *desktop.MouseEvent) {}

func (ic *ModernIcon) MouseOut() {
	if ic.disabled || ic.onTap == nil {
		return
	}
	ic.isHovered = false
	if ic.hoverHex != "" {
		ic.resource = resolveModernIcon(ic.rawSource, ic.colorHex)
		ic.img.Resource = ic.resource
		ic.img.Refresh()
	}
}

func (ic *ModernIcon) Cursor() desktop.Cursor {
	if !ic.disabled && ic.onTap != nil {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

// --- Catálogo y Resolución de Iconos Vectoriales (Lucide SVG & Theme) ---

// IconResource resuelve y genera un recurso vectorial fyne.Resource a partir de un nombre
// estándar ("search", "export", "filter", etc.) con el color hexadecimal solicitado opcional.
func IconResource(src any, colorHex ...string) fyne.Resource {
	hex := ""
	if len(colorHex) > 0 {
		hex = colorHex[0]
	}
	return resolveModernIcon(src, hex)
}

// resolveModernIcon convierte un nombre string o fyne.Resource en un recurso SVG con el color solicitado
func resolveModernIcon(src any, colorHex string) fyne.Resource {
	var baseRes fyne.Resource

	switch val := src.(type) {
	case fyne.Resource:
		baseRes = val
	case string:
		baseRes = getVectorIconByName(val)
	default:
		baseRes = theme.InfoIcon()
	}

	if baseRes == nil {
		baseRes = theme.InfoIcon()
	}

	if colorHex == "" {
		return baseRes
	}

	// Si es un SVG, inyectamos o reemplazamos el color directamente en el vector
	content := string(baseRes.Content())
	if strings.Contains(content, "<svg") {
		// Reemplazar currentColor o fill/stroke existente
		replaced := strings.ReplaceAll(content, "currentColor", colorHex)
		cleanName := strings.TrimSuffix(baseRes.Name(), ".svg")
		return fyne.NewStaticResource(cleanName+"_"+colorHex+".svg", []byte(replaced))
	}

	return theme.NewColoredResource(baseRes, theme.ColorNamePrimary)
}

func getVectorIconByName(name string) fyne.Resource {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "search", "buscar", "lupa":
		return iconSvgSearch
	case "close", "cancel", "cancelar", "x", "limpiar":
		return iconSvgClose
	case "palette", "paleta", "color", "photoshop", "art":
		return iconSvgPalette
	case "filter", "filtros", "sliders":
		return iconSvgSliders
	case "new", "nuevo", "plus", "add", "crear":
		return iconSvgPlus
	case "delete", "borrar", "trash", "eliminar":
		return iconSvgTrash
	case "print", "imprimir", "ficha":
		return iconSvgPrinter
	case "settings", "config", "ajustes", "gear":
		return iconSvgSettings
	case "info", "informacion":
		return iconSvgInfo
	case "error", "alerta", "alert":
		return iconSvgError
	case "question", "pregunta", "ayuda":
		return iconSvgQuestion
	case "check", "confirmar", "ok":
		return iconSvgCheck
	case "export", "exportar", "csv", "excel":
		return iconSvgFileSpreadsheet
	case "dashboard", "home", "layout-dashboard":
		return iconSvgDashboard
	case "users", "group":
		return iconSvgUsers
	case "user", "account", "profile":
		return iconSvgUser
	case "user-plus", "user_plus":
		return iconSvgUserPlus
	case "coins", "currency":
		return iconSvgCoins
	case "calculator":
		return iconSvgCalculator
	case "briefcase", "portfolio":
		return iconSvgBriefcase
	case "calendar", "date":
		return iconSvgCalendar
	case "trending-up", "chart-line":
		return iconSvgTrendingUp
	case "shield-check", "security":
		return iconSvgShieldCheck
	case "file-text", "document":
		return iconSvgFileText
	case "arrow-right", "siguiente":
		return iconSvgArrowRight
	case "chevron-left", "anterior":
		return iconSvgChevronLeft
	case "chevron-right":
		return iconSvgChevronRight
	case "eye", "ver", "ojo":
		return iconSvgEye
	case "wifi", "red":
		return iconSvgWifi
	case "sparkles", "destacado", "estrella":
		return iconSvgSparkles
	case "alert-triangle", "advertencia", "warning":
		return iconSvgAlertTriangle
	case "download", "descargar":
		return iconSvgDownload
	case "wallet", "cartera":
		return iconSvgWallet
	case "scale", "balanza":
		return iconSvgScale
	case "award", "premio", "trofeo":
		return iconSvgAward
	case "piggy-bank", "piggy", "piggybank", "hucha", "savings":
		return iconSvgPiggyBank
	case "save", "guardar", "grabar":
		return iconSvgSave
	case "edit", "editar", "modificar", "pencil":
		return iconSvgEdit
	case "upload", "subir", "import", "importar":
		return iconSvgUpload
	case "refresh", "refrescar", "sync", "sincronizar":
		return iconSvgRefresh
	case "database", "db", "basedatos", "sql":
		return iconSvgDatabase
	case "power", "apagar", "salir", "power-off", "shutdown":
		return iconSvgPower
	case "sun", "sol", "light", "tema":
		return iconSvgSun
	case "moon", "luna", "dark":
		return iconSvgMoon
	default:
		return nil
	}
}

// SVGs oficiales estándar Lucide (viewBox 0 0 24 24, stroke 2, vector nítido a cualquier escala)

var iconSvgSearch = fyne.NewStaticResource("lucide_search.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <circle cx="11" cy="11" r="8"/>
  <path d="m21 21-4.3-4.3"/>
</svg>`))

var iconSvgClose = fyne.NewStaticResource("lucide_close.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M18 6 6 18"/>
  <path d="m6 6 12 12"/>
</svg>`))

var iconSvgPalette = fyne.NewStaticResource("lucide_palette.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <circle cx="13.5" cy="6.5" r=".5" fill="currentColor"/>
  <circle cx="17.5" cy="10.5" r=".5" fill="currentColor"/>
  <circle cx="8.5" cy="7.5" r=".5" fill="currentColor"/>
  <circle cx="6.5" cy="12.5" r=".5" fill="currentColor"/>
  <path d="M12 2C6.5 2 2 6.5 2 12s4.5 10 10 10c.926 0 1.648-.746 1.648-1.688 0-.437-.18-.835-.437-1.125-.29-.289-.438-.652-.438-1.125a1.64 1.64 0 0 1 1.668-1.668h1.996c3.051 0 5.555-2.503 5.555-5.554C21.965 6.012 17.461 2 12 2z"/>
</svg>`))

var iconSvgSliders = fyne.NewStaticResource("lucide_sliders.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <line x1="21" x2="14" y1="4" y2="4"/>
  <line x1="10" x2="3" y1="4" y2="4"/>
  <line x1="21" x2="12" y1="12" y2="12"/>
  <line x1="8" x2="3" y1="12" y2="12"/>
  <line x1="21" x2="16" y1="20" y2="20"/>
  <line x1="12" x2="3" y1="20" y2="20"/>
  <line x1="14" x2="14" y1="2" y2="6"/>
  <line x1="8" x2="8" y1="10" y2="14"/>
  <line x1="16" x2="16" y1="18" y2="22"/>
</svg>`))

var iconSvgPlus = fyne.NewStaticResource("lucide_plus.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M5 12h14"/>
  <path d="M12 5v14"/>
</svg>`))

var iconSvgTrash = fyne.NewStaticResource("lucide_trash.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M3 6h18"/>
  <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/>
  <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/>
  <line x1="10" x2="10" y1="11" y2="17"/>
  <line x1="14" x2="14" y1="11" y2="17"/>
</svg>`))

var iconSvgPrinter = fyne.NewStaticResource("lucide_printer.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/>
  <path d="M6 9V3a1 1 0 0 1 1-1h10a1 1 0 0 1 1 1v6"/>
  <rect width="12" height="8" x="6" y="14" rx="1"/>
</svg>`))

var iconSvgSettings = fyne.NewStaticResource("lucide_settings.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/>
  <circle cx="12" cy="12" r="3"/>
</svg>`))

var iconSvgInfo = fyne.NewStaticResource("lucide_info.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <circle cx="12" cy="12" r="10"/>
  <path d="M12 16v-4"/>
  <path d="M12 8h.01"/>
</svg>`))

var iconSvgError = fyne.NewStaticResource("lucide_error.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <circle cx="12" cy="12" r="10"/>
  <line x1="12" x2="12" y1="8" y2="12"/>
  <line x1="12" x2="12.01" y1="16" y2="16"/>
</svg>`))

var iconSvgQuestion = fyne.NewStaticResource("lucide_question.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <circle cx="12" cy="12" r="10"/>
  <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/>
  <line x1="12" x2="12.01" y1="17" y2="17"/>
</svg>`))

var iconSvgCheck = fyne.NewStaticResource("lucide_check.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <polyline points="20 6 9 17 4 12"/>
</svg>`))

var iconSvgFileSpreadsheet = fyne.NewStaticResource("lucide_file_spreadsheet.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"/>
  <path d="M14 2v4a2 2 0 0 0 2 2h4"/>
  <path d="M8 13h2"/>
  <path d="M14 13h2"/>
  <path d="M8 17h2"/>
  <path d="M14 17h2"/>
</svg>`))

var iconSvgDashboard = fyne.NewStaticResource("lucide_dashboard.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <rect width="7" height="9" x="3" y="3" rx="1"/>
  <rect width="7" height="5" x="14" y="3" rx="1"/>
  <rect width="7" height="9" x="14" y="12" rx="1"/>
  <rect width="7" height="5" x="3" y="16" rx="1"/>
</svg>`))

var iconSvgUsers = fyne.NewStaticResource("lucide_users.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
  <circle cx="9" cy="7" r="4"/>
  <path d="M22 21v-2a4 4 0 0 0-3-3.87"/>
  <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
</svg>`))

var iconSvgUser = fyne.NewStaticResource("lucide_user.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/>
  <circle cx="12" cy="7" r="4"/>
</svg>`))

var iconSvgUserPlus = fyne.NewStaticResource("lucide_user_plus.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
  <circle cx="9" cy="7" r="4"/>
  <line x1="19" x2="19" y1="8" y2="14"/>
  <line x1="22" x2="16" y1="11" y2="11"/>
</svg>`))

var iconSvgCoins = fyne.NewStaticResource("lucide_coins.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <circle cx="8" cy="8" r="6"/>
  <path d="M18.09 10.37A6 6 0 1 1 10.34 18"/>
  <path d="M7 6h1v4"/>
  <path d="m16.71 13.88.7.71-2.82 2.82"/>
</svg>`))

var iconSvgCalculator = fyne.NewStaticResource("lucide_calculator.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <rect width="16" height="20" x="4" y="2" rx="2"/>
  <line x1="8" x2="16" y1="6" y2="6"/>
  <line x1="16" x2="16" y1="14" y2="18"/>
  <path d="M16 10h.01"/>
  <path d="M12 10h.01"/>
  <path d="M8 10h.01"/>
  <path d="M12 14h.01"/>
  <path d="M8 14h.01"/>
  <path d="M12 18h.01"/>
  <path d="M8 18h.01"/>
</svg>`))

var iconSvgBriefcase = fyne.NewStaticResource("lucide_briefcase.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <rect width="20" height="14" x="2" y="7" rx="2" ry="2"/>
  <path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"/>
</svg>`))

var iconSvgCalendar = fyne.NewStaticResource("lucide_calendar.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <rect width="18" height="18" x="3" y="4" rx="2" ry="2"/>
  <line x1="16" x2="16" y1="2" y2="6"/>
  <line x1="8" x2="8" y1="2" y2="6"/>
  <line x1="3" x2="21" y1="10" y2="10"/>
</svg>`))

var iconSvgTrendingUp = fyne.NewStaticResource("lucide_trending_up.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <polyline points="22 7 13.5 15.5 8.5 10.5 2 17"/>
  <polyline points="16 7 22 7 22 13"/>
</svg>`))

var iconSvgShieldCheck = fyne.NewStaticResource("lucide_shield_check.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10"/>
  <path d="m9 12 2 2 4-4"/>
</svg>`))

var iconSvgFileText = fyne.NewStaticResource("lucide_file_text.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"/>
  <path d="M14 2v4a2 2 0 0 0 2 2h4"/>
  <path d="M10 9H8"/>
  <path d="M16 13H8"/>
  <path d="M16 17H8"/>
</svg>`))

var iconSvgArrowRight = fyne.NewStaticResource("lucide_arrow_right.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M5 12h14"/>
  <path d="m12 5 7 7-7 7"/>
</svg>`))

var iconSvgChevronLeft = fyne.NewStaticResource("lucide_chevron_left.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="m15 18-6-6 6-6"/>
</svg>`))

var iconSvgChevronRight = fyne.NewStaticResource("lucide_chevron_right.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="m9 18 6-6-6-6"/>
</svg>`))

var iconSvgEye = fyne.NewStaticResource("lucide_eye.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z"/>
  <circle cx="12" cy="12" r="3"/>
</svg>`))

var iconSvgWifi = fyne.NewStaticResource("lucide_wifi.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M12 20h.01"/>
  <path d="M2 8.82a15 15 0 0 1 20 0"/>
  <path d="M5 12.86a10 10 0 0 1 14 0"/>
  <path d="M8.5 16.43a5 5 0 0 1 7 0"/>
</svg>`))

var iconSvgSparkles = fyne.NewStaticResource("lucide_sparkles.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="m12 3-1.912 5.813a2 2 0 0 1-1.275 1.275L3 12l5.813 1.912a2 2 0 0 1 1.275 1.275L12 21l1.912-5.813a2 2 0 0 1 1.275-1.275L21 12l-5.813-1.912a2 2 0 0 1-1.275-1.275L12 3Z"/>
</svg>`))

var iconSvgAlertTriangle = fyne.NewStaticResource("lucide_alert_triangle.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/>
  <line x1="12" x2="12" y1="9" y2="13"/>
  <line x1="12" x2="12.01" y1="17" y2="17"/>
</svg>`))

var iconSvgDownload = fyne.NewStaticResource("lucide_download.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
  <polyline points="7 10 12 15 17 10"/>
  <line x1="12" x2="12" y1="15" y2="3"/>
</svg>`))

var iconSvgWallet = fyne.NewStaticResource("lucide_wallet.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M19 7V4a1 1 0 0 0-1-1H5a2 2 0 0 0 0 4h15a1 1 0 0 1 1 1v4h-3a2 2 0 0 0 0 4h3a1 1 0 0 0 1-1v-2a1 1 0 0 0-1-1"/>
  <path d="M3 5v14a2 2 0 0 0 2 2h15a1 1 0 0 0 1-1v-4"/>
</svg>`))

var iconSvgScale = fyne.NewStaticResource("lucide_scale.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="m16 16 3-8 3 8c-.87.65-1.92 1-3 1s-2.13-.35-3-1Z"/>
  <path d="m2 16 3-8 3 8c-.87.65-1.92 1-3 1s-2.13-.35-3-1Z"/>
  <path d="M7 21h10"/>
  <path d="M12 3v18"/>
  <path d="M3 7h2c2 0 5-1 7-2 2 1 5 2 7 2h2"/>
</svg>`))

var iconSvgSave = fyne.NewStaticResource("lucide_save.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M15.2 3a2 2 0 0 1 1.4.6l3.8 3.8a2 2 0 0 1 .6 1.4V19a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2z"/>
  <path d="M17 21v-7a1 1 0 0 0-1-1H8a1 1 0 0 0-1 1v7"/>
  <path d="M7 3v4a1 1 0 0 0 1 1h7"/>
</svg>`))

var iconSvgEdit = fyne.NewStaticResource("lucide_edit.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M21.174 6.812a1 1 0 0 0-3.986-3.987L3.842 16.174a2 2 0 0 0-.5.83l-1.321 4.352a.5.5 0 0 0 .623.622l4.353-1.32a2 2 0 0 0 .83-.497z"/>
  <path d="m15 5 4 4"/>
</svg>`))

var iconSvgUpload = fyne.NewStaticResource("lucide_upload.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
  <polyline points="17 8 12 3 7 8"/>
  <line x1="12" x2="12" y1="3" y2="15"/>
</svg>`))

var iconSvgRefresh = fyne.NewStaticResource("lucide_refresh.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/>
  <path d="M21 3v5h-5"/>
  <path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/>
  <path d="M8 16H3v5"/>
</svg>`))

var iconSvgAward = fyne.NewStaticResource("lucide_award.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <circle cx="12" cy="8" r="6"/>
  <path d="M15.477 12.89 17 22l-5-3-5 3 1.523-9.11"/>
</svg>`))

var iconSvgPiggyBank = fyne.NewStaticResource("lucide_piggy_bank.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M19 5c-1.5 0-2.8 1.4-3 2-3.5-1.5-11-.3-11 5 0 1.8 0 3 2 4.5V20h4v-2h3v2h4v-4c1-.5 1.7-1 2-2h2v-4h-2c0-1-.5-1.5-1-2V5z"/>
  <path d="M2 9v1c0 1.1.9 2 2 2h1"/>
  <path d="M16 11h.01"/>
</svg>`))

var iconSvgDatabase = fyne.NewStaticResource("lucide_database.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <ellipse cx="12" cy="5" rx="9" ry="3"/>
  <path d="M3 5v14c0 1.66 4.03 3 9 3s9-1.34 9-3V5"/>
  <path d="M3 12c0 1.66 4.03 3 9 3s9-1.34 9-3"/>
</svg>`))

var iconSvgPower = fyne.NewStaticResource("lucide_power.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M12 2v10"/>
  <path d="M18.4 6.6a9 9 0 1 1-12.77.04"/>
</svg>`))

var iconSvgSun = fyne.NewStaticResource("lucide_sun.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <circle cx="12" cy="12" r="4"/>
  <path d="M12 2v2"/>
  <path d="M12 20v2"/>
  <path d="m4.93 4.93 1.41 1.41"/>
  <path d="m17.66 17.66 1.41 1.41"/>
  <path d="M2 12h2"/>
  <path d="M20 12h2"/>
  <path d="m6.34 17.66-1.41 1.41"/>
  <path d="m19.07 4.93-1.41 1.41"/>
</svg>`))

var iconSvgMoon = fyne.NewStaticResource("lucide_moon.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z"/>
</svg>`))




