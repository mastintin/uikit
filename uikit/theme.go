package uikit

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const (
	ThemeDark  = "dark"
	ThemeLight = "light"
)

// ─────────────────────────────────────────────────────────────────────────────
// Paleta Centralizada de Colores Fundamentales (Tokens Primitivos)
// ─────────────────────────────────────────────────────────────────────────────
var (
	// Pizarras / Grises (Slate)
	ColorSlate950 = color.NRGBA{R: 2, G: 6, B: 23, A: 255}    // #020617 - Fondo ultra oscuro base
	ColorSlate900 = color.NRGBA{R: 15, G: 23, B: 42, A: 255}   // #0f172a - Texto primario claro / Navbar Nice / Tarjetas Dark
	ColorSlate800 = color.NRGBA{R: 30, G: 41, B: 59, A: 255}   // #1e293b - Superficies elevadas Dark / Botón Dark
	ColorSlate700 = color.NRGBA{R: 51, G: 65, B: 85, A: 255}   // #334155 - Bordes y divisores Dark / Texto muted Light
	ColorSlate600 = color.NRGBA{R: 71, G: 85, B: 105, A: 255}  // #475569 - Hover Dark
	ColorSlate500 = color.NRGBA{R: 100, G: 116, B: 139, A: 255} // #64748b - Placeholder / Muted neutro
	ColorSlate400 = color.NRGBA{R: 148, G: 163, B: 184, A: 255} // #94a3b8 - Texto muted Dark / Deshabilitados
	ColorSlate300 = color.NRGBA{R: 203, G: 213, B: 225, A: 255} // #cbd5e1 - Bordes tarjetas e inputs
	ColorSlate200 = color.NRGBA{R: 226, G: 232, B: 240, A: 255} // #e2e8f0 - Bordes generales y separadores Light
	ColorSlate100 = color.NRGBA{R: 241, G: 245, B: 249, A: 255} // #f1f5f9 - Fondo Nice / Divisores / Hover
	ColorSlate50  = color.NRGBA{R: 248, G: 250, B: 252, A: 255} // #f8fafc - Fondo Light / Texto primario Dark
	// ColorWhite defined in typography.go // #ffffff - Blanco puro
	// ColorTransparent

	// Azules (Blue / Sky)
	ColorBlue700 = color.NRGBA{R: 29, G: 78, B: 216, A: 255}   // #1d4ed8 - Azul intenso acento Light
	ColorBlue600 = color.NRGBA{R: 37, G: 99, B: 235, A: 255}   // #2563eb - Azul corporativo estándar
	ColorBlue500 = color.NRGBA{R: 59, G: 130, B: 246, A: 255}  // #3b82f6 - Cursor / Primario focus
	ColorBlue400 = color.NRGBA{R: 96, G: 165, B: 250, A: 255}  // #60a5fa - Azul claro brutos Dark
	ColorBlue200 = color.NRGBA{R: 191, G: 219, B: 254, A: 255} // #bfdbfe - Azul pastel bordes
	ColorBlue100 = color.NRGBA{R: 219, G: 234, B: 254, A: 255} // #dbeafe - Hover primario suave
	ColorBlue50  = color.NRGBA{R: 239, G: 246, B: 255, A: 255} // #eff6ff - Fondo botón primario Light
	ColorSky300  = color.NRGBA{R: 147, G: 197, B: 253, A: 255} // #93c5fd - Sidebar nav inactivo xlsJubila

	// Índigos (Indigo - Acentos Tema Dark)
	ColorIndigo700 = color.NRGBA{R: 67, G: 56, B: 202, A: 255}  // #4338ca
	ColorIndigo600 = color.NRGBA{R: 79, G: 70, B: 229, A: 255}  // #4f46e5 - Botón primario Dark
	ColorIndigo500 = color.NRGBA{R: 99, G: 102, B: 241, A: 255}  // #6366f1 - Bordes botón primario Dark
	ColorIndigo400 = color.NRGBA{R: 129, G: 140, B: 248, A: 255} // #818cf8 - Acento Dark
	ColorIndigo800 = color.NRGBA{R: 55, G: 48, B: 163, A: 255}  // #3730a3 - Press botón primario Dark

	// Cobalto Ejecutivo Nice (xlsJubila)
	ColorCobalt800 = color.NRGBA{R: 30, G: 64, B: 175, A: 255} // #1e40af - Acento cobalto Nice

	// Verdes / Esmeraldas (Green / Emerald)
	ColorGreen800 = color.NRGBA{R: 22, G: 101, B: 52, A: 255}  // #166534 - Hover borde verde Nice
	ColorGreen700 = color.NRGBA{R: 21, G: 128, B: 61, A: 255}  // #15803d - Verde 700 éxito Light
	ColorGreen900 = color.NRGBA{R: 10, G: 77, B: 40, A: 255}   // #0a4d28 - Press verde Nice
	ColorEmerald800 = color.NRGBA{R: 14, G: 106, B: 55, A: 255} // #0e6a37 - Borde primario Nice
	ColorEmerald700 = color.NRGBA{R: 4, G: 120, B: 87, A: 255}  // #047857 - Pensión neta / precio Light
	ColorEmerald600 = color.NRGBA{R: 16, G: 124, B: 65, A: 255} // #107c41 - Verde corporativo xlsJubila
	ColorEmerald400 = color.NRGBA{R: 52, G: 211, B: 153, A: 255} // #34d399 - Éxito / Pensión neta Dark
	ColorEmerald200 = color.NRGBA{R: 167, G: 243, B: 208, A: 255} // #a7f3d0 - Primario muted Nice

	// Tonos Oliva / Salvia (Pastilla ACTIVO y números destacados)
	ColorActiveBg     = color.NRGBA{R: 232, G: 240, B: 230, A: 255} // #E8F0E6 - Oliva blanqueado muy sutil (fondo)
	ColorActiveFg     = color.NRGBA{R: 74, G: 93, B: 78, A: 255}   // #4A5D4E - Verde seco / salvia oscuro (texto y números)
	ColorActiveBorder = color.NRGBA{R: 205, G: 221, B: 202, A: 255} // #CDDDCA - Borde oliva suave

	// Rojos (Red / Danger)
	ColorRed800 = color.NRGBA{R: 153, G: 27, B: 27, A: 255}   // #991b1b - Rojo vino / Texto peligro oscuro
	ColorRed700 = color.NRGBA{R: 185, G: 28, B: 28, A: 255}   // #b91c1c - Rojo intenso press/hover
	ColorRed600 = color.NRGBA{R: 220, G: 38, B: 38, A: 255}   // #dc2626 - Peligro / Error base
	ColorRed400 = color.NRGBA{R: 248, G: 113, B: 113, A: 255} // #f87171 - Peligro Dark
	ColorRed300 = color.NRGBA{R: 252, G: 165, B: 165, A: 255} // #fca5a5 - Peligro atenuado / borde
	ColorRed200 = color.NRGBA{R: 254, G: 202, B: 202, A: 255} // #fecaca - Borde suave / press bg
	ColorRed100 = color.NRGBA{R: 254, G: 226, B: 226, A: 255} // #fee2e2 - Hover fondo peligro Light
	ColorRed50  = color.NRGBA{R: 254, G: 242, B: 242, A: 255} // #fef2f2 - Fondo normal peligro Light

	// Ámbares (Amber / Warning)
	ColorAmber800 = color.NRGBA{R: 146, G: 64, B: 14, A: 255}  // #92400e - Alerta / Warning texto sobre fondo ámbar (#fef3c7)
	ColorAmber700 = color.NRGBA{R: 180, G: 83, B: 9, A: 255}  // #b45309 - Alerta / Warning apagado
	ColorAmber600 = color.NRGBA{R: 217, G: 119, B: 6, A: 255} // #d97706 - Ámbar dorado luminoso y de alto contraste (Light/Nice)
	ColorAmber400 = color.NRGBA{R: 251, G: 191, B: 36, A: 255} // #fbbf24 - Alerta / Warning Dark

	// Tonos Warning / Terracota Suave (Borrado y Advertencias)
	ColorWarningBg          = color.NRGBA{R: 250, G: 219, B: 216, A: 255} // #FADBD8 - Crema/terracota muy suave (fondo botón/badge)
	ColorWarningFg          = color.NRGBA{R: 176, G: 58, B: 46, A: 255}  // #B03A2E - Terracota intenso (texto e icono)
	ColorWarningBorder      = color.NRGBA{R: 245, G: 183, B: 177, A: 255} // #F5B7B1 - Borde suave
	ColorWarningHoverBg     = color.NRGBA{R: 245, G: 203, B: 167, A: 255} // #F5CBA7 - Terracota claro (hover)
	ColorWarningHoverFg     = color.NRGBA{R: 148, G: 49, B: 38, A: 255}  // #943126 - Terracota oscuro (hover texto)
	ColorWarningHoverBorder = color.NRGBA{R: 237, G: 187, B: 153, A: 255} // #EDBB99 - Borde hover
	ColorWarningPressBg     = color.NRGBA{R: 237, G: 187, B: 153, A: 255} // #EDBB99 - Fondo press
	ColorWarningPressFg     = color.NRGBA{R: 120, G: 40, B: 31, A: 255}  // #78281F - Texto press
	ColorWarningPressBorder = color.NRGBA{R: 229, G: 152, B: 102, A: 255} // #E59866 - Borde press

	// Transparencias y Sombras
	ColorBlackAlpha20 = color.NRGBA{R: 0, G: 0, B: 0, A: 20}
	ColorBlackAlpha25 = color.NRGBA{R: 0, G: 0, B: 0, A: 25}
	ColorWhiteAlpha30 = color.NRGBA{R: 255, G: 255, B: 255, A: 30}
)

const (
	ColorNameBtnDanger       fyne.ThemeColorName = "modernBtnDanger"
	ColorNameBtnDangerMuted  fyne.ThemeColorName = "modernBtnDangerMuted"

	// Colores de ModernButton Danger
	ColorNameBtnDangerBg          fyne.ThemeColorName = "modernBtnDangerBg"
	ColorNameBtnDangerFg          fyne.ThemeColorName = "modernBtnDangerFg"
	ColorNameBtnDangerBorder      fyne.ThemeColorName = "modernBtnDangerBorder"
	ColorNameBtnDangerHoverBg     fyne.ThemeColorName = "modernBtnDangerHoverBg"
	ColorNameBtnDangerHoverFg     fyne.ThemeColorName = "modernBtnDangerHoverFg"
	ColorNameBtnDangerHoverBorder fyne.ThemeColorName = "modernBtnDangerHoverBorder"
	ColorNameBtnDangerPressBg     fyne.ThemeColorName = "modernBtnDangerPressBg"
	ColorNameBtnDangerPressFg     fyne.ThemeColorName = "modernBtnDangerPressFg"
	ColorNameBtnDangerPressBorder fyne.ThemeColorName = "modernBtnDangerPressBorder"
	ColorNameBtnMuted        fyne.ThemeColorName = "modernBtnMuted"
	ColorNameBtnPrimaryMuted fyne.ThemeColorName = "modernBtnPrimaryMuted"

	// Colores de ModernButton Primary
	ColorNameBtnPrimaryBg          fyne.ThemeColorName = "modernBtnPrimaryBg"
	ColorNameBtnPrimaryFg          fyne.ThemeColorName = "modernBtnPrimaryFg"
	ColorNameBtnPrimaryBorder      fyne.ThemeColorName = "modernBtnPrimaryBorder"
	ColorNameBtnPrimaryHoverBg     fyne.ThemeColorName = "modernBtnPrimaryHoverBg"
	ColorNameBtnPrimaryHoverFg     fyne.ThemeColorName = "modernBtnPrimaryHoverFg"
	ColorNameBtnPrimaryHoverBorder fyne.ThemeColorName = "modernBtnPrimaryHoverBorder"
	ColorNameBtnPrimaryPressBg     fyne.ThemeColorName = "modernBtnPrimaryPressBg"
	ColorNameBtnPrimaryPressFg     fyne.ThemeColorName = "modernBtnPrimaryPressFg"
	ColorNameBtnPrimaryPressBorder fyne.ThemeColorName = "modernBtnPrimaryPressBorder"

	// Colores de ModernButton Default
	ColorNameBtnDefaultBg          fyne.ThemeColorName = "modernBtnDefaultBg"
	ColorNameBtnDefaultFg          fyne.ThemeColorName = "modernBtnDefaultFg"
	ColorNameBtnDefaultBorder      fyne.ThemeColorName = "modernBtnDefaultBorder"
	ColorNameBtnDefaultHoverBg     fyne.ThemeColorName = "modernBtnDefaultHoverBg"
	ColorNameBtnDefaultHoverFg     fyne.ThemeColorName = "modernBtnDefaultHoverFg"
	ColorNameBtnDefaultHoverBorder fyne.ThemeColorName = "modernBtnDefaultHoverBorder"
	ColorNameBtnDefaultPressBg     fyne.ThemeColorName = "modernBtnDefaultPressBg"
	ColorNameBtnDefaultPressFg     fyne.ThemeColorName = "modernBtnDefaultPressFg"
	ColorNameBtnDefaultPressBorder fyne.ThemeColorName = "modernBtnDefaultPressBorder"

	// --- Tokens Semánticos de Superficies y Fondos ---
	ColorNameAppBackground   fyne.ThemeColorName = "appBackground"   // Fondo de la app/detalle
	ColorNameNavbarBg        fyne.ThemeColorName = "navbarBg"        // Barra superior
	ColorNameCardBg          fyne.ThemeColorName = "cardBg"          // Tarjetas y paneles
	ColorNameSidebarBg       fyne.ThemeColorName = "sidebarBg"       // Sidebar principal
	ColorNameSidebarShadedBg fyne.ThemeColorName = "sidebarShadedBg" // Sidebar sombreado
	ColorNameFilterDrawerBg  fyne.ThemeColorName = "filterDrawerBg"  // Cajón de filtros
	ColorNameInsetPanelBg    fyne.ThemeColorName = "insetPanelBg"    // Paneles interiores

	// --- Tokens Semánticos de Bordes y Separadores ---
	ColorNameBorder     fyne.ThemeColorName = "modernBorder" // Borde general 1px
	ColorNameCardBorder fyne.ThemeColorName = "cardBorder"   // Borde de tarjetas
	ColorNameDivider    fyne.ThemeColorName = "divider"      // Separador sutil de lista

	// --- Tokens Semánticos de Textos ---
	ColorNameTextPrimary   fyne.ThemeColorName = "textPrimary"   // Título principal
	ColorNameTextMuted     fyne.ThemeColorName = "textMuted"     // Subtítulos y secciones
	ColorNameNavInactiveFg fyne.ThemeColorName = "navInactiveFg" // Texto inactivo de navbar

	// --- Tokens Semánticos de Estados y Acentos Centralizados ---
	ColorNameAccent  fyne.ThemeColorName = "modernAccent"  // Acento principal
	ColorNameSuccess fyne.ThemeColorName = "modernSuccess" // Verde corporativo / éxito
	ColorNameWarning fyne.ThemeColorName = "modernWarning" // Ámbar alerta / retención
	ColorNameDanger  fyne.ThemeColorName = "modernDanger"  // Rojo peligro

	// --- Tokens Semánticos de Elementos de Lista (Obras / Cuadros) ---
	ColorNameItemID          fyne.ThemeColorName = "itemID"          // Etiqueta ID (#2563eb / #60a5fa)
	ColorNameItemPrice       fyne.ThemeColorName = "itemPrice"       // Precio (#047857 / #34d399)
	ColorNameItemThumbBg     fyne.ThemeColorName = "itemThumbBg"     // Fondo miniatura (#e2e8f0 / #1e293b)
	ColorNameItemThumbBorder fyne.ThemeColorName = "itemThumbBorder" // Borde miniatura


	// Tamaños de ModernCard configurables en el tema
	SizeNameCardPaddingTop    fyne.ThemeSizeName = "modernCardPaddingTop"
	SizeNameCardPaddingBottom fyne.ThemeSizeName = "modernCardPaddingBottom"
	SizeNameCardPaddingLeft   fyne.ThemeSizeName = "modernCardPaddingLeft"
	SizeNameCardPaddingRight  fyne.ThemeSizeName = "modernCardPaddingRight"
)

// AppTheme personaliza los colores y tamaños de la aplicación para una apariencia
// moderna, profesional y con textos nítidos de alto contraste.
type AppTheme struct {
	fyne.Theme
	name          string
	customColors  map[fyne.ThemeColorName]color.Color
	customSizes   map[fyne.ThemeSizeName]float32
	forcedVariant *fyne.ThemeVariant
}

// Name devuelve el nombre identificativo del tema actual.
func (t *AppTheme) Name() string {
	if t.name == "" {
		if t.IsDark() {
			return ThemeDark
		}
		return ThemeLight
	}
	return t.name
}

// NewLightTheme crea la variante de tema claro de la aplicación (superficies limpias, textos oscuros nítidos).
func NewLightTheme() *AppTheme {
	light := theme.VariantLight
	th := &AppTheme{
		Theme:         theme.DefaultTheme(),
		name:          ThemeLight,
		customColors:  make(map[fyne.ThemeColorName]color.Color),
		customSizes:   make(map[fyne.ThemeSizeName]float32),
		forcedVariant: &light,
	}
	th.SetCardPadding(18, 20, 18, 20)
	return th
}

// NewDarkTheme crea la variante de tema oscuro de la aplicación (paleta slate/indigo).
func NewDarkTheme() *AppTheme {
	dark := theme.VariantDark
	th := &AppTheme{
		Theme:         theme.DarkTheme(),
		name:          ThemeDark,
		customColors:  make(map[fyne.ThemeColorName]color.Color),
		customSizes:   make(map[fyne.ThemeSizeName]float32),
		forcedVariant: &dark,
	}
	th.SetCardPadding(18, 20, 18, 20)
	return th
}

// NewAppTheme devuelve el tema claro estándar de la aplicación.
func NewAppTheme() *AppTheme {
	return NewLightTheme()
}

// ForcedVariant devuelve la variante forzada configurada o nil si sigue la del sistema.
func (t *AppTheme) ForcedVariant() *fyne.ThemeVariant {
	return t.forcedVariant
}

// SetForcedVariant permite forzar una variante de tema específica.
func (t *AppTheme) SetForcedVariant(v fyne.ThemeVariant) *AppTheme {
	t.forcedVariant = &v
	if v == theme.VariantDark {
		t.Theme = theme.DarkTheme()
		t.name = ThemeDark
	} else {
		t.Theme = theme.DefaultTheme()
		if t.name == ThemeDark {
			t.name = ThemeLight
		}
	}
	return t
}

// IsDark indica si el tema actual está en modo oscuro.
func (t *AppTheme) IsDark() bool {
	if t.name == ThemeDark {
		return true
	}
	if t.forcedVariant != nil {
		return *t.forcedVariant == theme.VariantDark
	}
	return false
}

// IsDarkTheme comprueba de forma global si la aplicación o el tema activo está en modo oscuro.
func IsDarkTheme() bool {
	if app := fyne.CurrentApp(); app != nil && app.Settings() != nil {
		th := app.Settings().Theme()
		if appTh, ok := th.(*AppTheme); ok {
			return appTh.IsDark()
		}
		return app.Settings().ThemeVariant() == theme.VariantDark
	}
	return false
}

var defaultResolverTheme = NewLightTheme()

// ThemeColor devuelve el color del token semántico solicitado según el tema activo y la variante actual.
func ThemeColor(name fyne.ThemeColorName) color.Color {
	variant := theme.VariantLight

	if app := fyne.CurrentApp(); app != nil && app.Settings() != nil {
		th := app.Settings().Theme()
		if appTheme, ok := th.(*AppTheme); ok {
			if appTheme.ForcedVariant() != nil {
				variant = *appTheme.ForcedVariant()
			} else {
				variant = app.Settings().ThemeVariant()
			}
			return appTheme.Color(name, variant)
		}
		variant = app.Settings().ThemeVariant()
	}
	return defaultResolverTheme.Color(name, variant)
}

// SetCustomSize permite registrar o sobrescribir un tamaño temático específico.
func (t *AppTheme) SetCustomSize(name fyne.ThemeSizeName, s float32) {
	if t.customSizes == nil {
		t.customSizes = make(map[fyne.ThemeSizeName]float32)
	}
	t.customSizes[name] = s
}

// SetCardPaddingTop configura el padding superior por defecto para las tarjetas en el tema.
func (t *AppTheme) SetCardPaddingTop(pt float32) {
	t.SetCustomSize(SizeNameCardPaddingTop, pt)
}

// SetCardPaddingBottom configura el padding inferior por defecto para las tarjetas en el tema.
func (t *AppTheme) SetCardPaddingBottom(pb float32) {
	t.SetCustomSize(SizeNameCardPaddingBottom, pb)
}

// SetCardVerticalPadding configura el padding superior e inferior por defecto para las tarjetas en el tema.
func (t *AppTheme) SetCardVerticalPadding(top, bottom float32) {
	t.SetCustomSize(SizeNameCardPaddingTop, top)
	t.SetCustomSize(SizeNameCardPaddingBottom, bottom)
}

// SetCardPadding configura todos los paddings por defecto de las tarjetas en el tema (top, right, bottom, left).
func (t *AppTheme) SetCardPadding(top, right, bottom, left float32) {
	t.SetCustomSize(SizeNameCardPaddingTop, top)
	t.SetCustomSize(SizeNameCardPaddingRight, right)
	t.SetCustomSize(SizeNameCardPaddingBottom, bottom)
	t.SetCustomSize(SizeNameCardPaddingLeft, left)
}

// SetCustomColor permite registrar o sobrescribir un color temático específico.
func (t *AppTheme) SetCustomColor(name fyne.ThemeColorName, c color.Color) {
	if t.customColors == nil {
		t.customColors = make(map[fyne.ThemeColorName]color.Color)
	}
	t.customColors[name] = c
}

// SetButtonPrimaryColors permite personalizar los colores del botón primario en el tema.
func (t *AppTheme) SetButtonPrimaryColors(normalBg, normalFg, normalBorder, hoverBg, hoverFg, hoverBorder color.Color) {
	t.SetCustomColor(ColorNameBtnPrimaryBg, normalBg)
	t.SetCustomColor(ColorNameBtnPrimaryFg, normalFg)
	t.SetCustomColor(ColorNameBtnPrimaryBorder, normalBorder)
	t.SetCustomColor(ColorNameBtnPrimaryHoverBg, hoverBg)
	t.SetCustomColor(ColorNameBtnPrimaryHoverFg, hoverFg)
	t.SetCustomColor(ColorNameBtnPrimaryHoverBorder, hoverBorder)
}

// SetButtonDefaultColors permite personalizar los colores del botón por defecto en el tema.
func (t *AppTheme) SetButtonDefaultColors(normalBg, normalFg, normalBorder, hoverBg, hoverFg, hoverBorder color.Color) {
	t.SetCustomColor(ColorNameBtnDefaultBg, normalBg)
	t.SetCustomColor(ColorNameBtnDefaultFg, normalFg)
	t.SetCustomColor(ColorNameBtnDefaultBorder, normalBorder)
	t.SetCustomColor(ColorNameBtnDefaultHoverBg, hoverBg)
	t.SetCustomColor(ColorNameBtnDefaultHoverFg, hoverFg)
	t.SetCustomColor(ColorNameBtnDefaultHoverBorder, hoverBorder)
}

func (t *AppTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if t.forcedVariant != nil {
		variant = *t.forcedVariant
	}
	if t.customColors != nil {
		if c, ok := t.customColors[name]; ok && c != nil {
			return c
		}
	}
	pal := GetThemePalette(t.Name())
	if pal != nil {
		if c := pal.Color(name); c != nil {
			return c
		}
	}
	return t.Theme.Color(name, variant)
}


func (t *AppTheme) Size(name fyne.ThemeSizeName) float32 {
	if t.customSizes != nil {
		if s, ok := t.customSizes[name]; ok {
			return s
		}
	}
	switch name {
	case theme.SizeNameButtonRadius, theme.SizeNameInputRadius, theme.SizeNameSelectionRadius:
		return 6.0 // 6px idéntico a var(--radius-sm): 6px de goCuadros
	case theme.SizeNamePadding:
		return 6.0
	case theme.SizeNameText:
		return 14.5
	case theme.SizeNameHeadingText:
		return 22.0
	case theme.SizeNameSubHeadingText:
		return 16.0
	case theme.SizeNameCaptionText:
		return 11.5
	}
	return t.Theme.Size(name)
}

func (t *AppTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace && style.Bold {
		return t.Theme.Font(fyne.TextStyle{Monospace: true})
	}
	return t.Theme.Font(style)
}



var currentThemeName = ThemeLight

// GetCurrentThemeName devuelve el identificador del tema activo ("light", "dark", o nombre de tema CSS).
func GetCurrentThemeName() string {
	if app := fyne.CurrentApp(); app != nil && app.Settings() != nil {
		if appTh, ok := app.Settings().Theme().(*AppTheme); ok {
			return appTh.Name()
		}
	}
	return currentThemeName
}

// CreateThemeByName construye una instancia de AppTheme a partir del identificador de tema.
func CreateThemeByName(name string) *AppTheme {
	pal := GetThemePalette(name)
	themeName := ThemeLight
	if pal != nil {
		themeName = pal.Name
	}
	switch strings.ToLower(themeName) {
	case ThemeDark:
		return NewDarkTheme()
	case ThemeLight:
		return NewLightTheme()
	default:
		variant := theme.VariantLight
		baseTheme := theme.DefaultTheme()
		if pal != nil && pal.IsDark {
			variant = theme.VariantDark
			baseTheme = theme.DarkTheme()
		}
		th := &AppTheme{
			Theme:         baseTheme,
			name:          themeName,
			customColors:  make(map[fyne.ThemeColorName]color.Color),
			customSizes:   make(map[fyne.ThemeSizeName]float32),
			forcedVariant: &variant,
		}
		th.SetCardPadding(18, 20, 18, 20)
		return th
	}
}

// RefreshTree recorre en profundidad el árbol de CanvasObjects invocando Refresh()
// en cada widget y contenedor, asegurando que todos los componentes reactivos
// actualicen sus paletas de color, bordes, estados y fuentes tras un cambio de tema en vivo.
func RefreshTree(obj fyne.CanvasObject) {
	if obj == nil {
		return
	}
	obj.Refresh()
	if c, ok := obj.(*fyne.Container); ok {
		for _, child := range c.Objects {
			RefreshTree(child)
		}
	}
}

// SetAppTheme aplica el tema globalmente en vivo (en caliente) y persiste la selección en preferencias.
func SetAppTheme(name string) *AppTheme {
	th := CreateThemeByName(name)
	currentThemeName = th.Name()
	if app := fyne.CurrentApp(); app != nil {
		app.Settings().SetTheme(th)
		if prefs := app.Preferences(); prefs != nil {
			prefs.SetString("app_theme", th.Name())
		}
		if drv := app.Driver(); drv != nil {
			for _, win := range drv.AllWindows() {
				if c := win.Content(); c != nil {
					RefreshTree(c)
				}
			}
		}
	}
	RefreshAllDatePickers()
	return th
}

// LoadInitialTheme recupera el tema guardado en preferencias y lo aplica de inmediato. Por defecto es "light".
func LoadInitialTheme() *AppTheme {
	saved := ThemeLight
	if app := fyne.CurrentApp(); app != nil {
		if prefs := app.Preferences(); prefs != nil {
			saved = prefs.StringWithFallback("app_theme", ThemeLight)
			if saved == "nice" {
				saved = ThemeLight
			}
			glass := prefs.BoolWithFallback("app_glass", true)
			SetGlobalGlass(glass)
		}
	}
	return SetAppTheme(saved)
}

// --- Funciones Globales de Acceso Semántico Centralizado ---
// Permiten que cualquier control o vista de la aplicación acceda a la paleta centralizada
// sin hardcodear colores ni duplicar condiciones if/else de temas.

// ThemeColorTextPrimary devuelve el color principal de texto según el tema activo.
func ThemeColorTextPrimary() color.Color {
	return ThemeColor(ColorNameTextPrimary)
}

// ThemeColorTextMuted devuelve el color secundario/atenuado de texto según el tema activo.
func ThemeColorTextMuted() color.Color {
	return ThemeColor(ColorNameTextMuted)
}

// ThemeColorAccent devuelve el color de acento según el tema activo.
func ThemeColorAccent() color.Color {
	return ThemeColor(ColorNameAccent)
}

// ThemeColorSuccess devuelve el color de éxito o datos positivos según el tema activo.
func ThemeColorSuccess() color.Color {
	return ThemeColor(ColorNameSuccess)
}

// ThemeColorWarning devuelve el color de aviso/retención según el tema activo.
func ThemeColorWarning() color.Color {
	return ThemeColor(ColorNameWarning)
}

// ThemeColorDanger devuelve el color de peligro/error según el tema activo.
func ThemeColorDanger() color.Color {
	return ThemeColor(ColorNameDanger)
}



// ThemeAccentHex devuelve la cadena hexadecimal del color de acento del tema activo (para iconos/SVG).
func ThemeAccentHex() string {
	return ColorToHex(ThemeColorAccent())
}

// ThemeSuccessHex devuelve la cadena hexadecimal del color de éxito del tema activo.
func ThemeSuccessHex() string {
	return ColorToHex(ThemeColorSuccess())
}

// ThemeWarningHex devuelve la cadena hexadecimal del color de advertencia del tema activo.
func ThemeWarningHex() string {
	return ColorToHex(ThemeColorWarning())
}

// ThemeDangerHex devuelve la cadena hexadecimal del color de peligro/error del tema activo.
func ThemeDangerHex() string {
	return ColorToHex(ThemeColorDanger())
}

