package uikit

import (
	"image/color"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// HeaderCardThemeColors define los colores específicos de ModernHeaderCard (tarjeta partida).
type HeaderCardThemeColors struct {
	BannerTop     color.NRGBA // Color superior del degradado del banner
	BannerBottom  color.NRGBA // Color inferior del degradado del banner
	BannerDivider color.Color // Línea divisoria inferior del banner
	BodyBg        color.Color // Fondo del cuerpo inferior
	Border        color.Color // Borde perimetral exterior de la tarjeta
	TitleColor    color.Color // Color del título del banner
	SubtitleColor color.Color // Color del subtítulo del banner
}

// BadgeVariantColors define el triplete visual (fondo, texto, borde) de un estado de badge.
type BadgeVariantColors struct {
	Bg     color.Color
	Text   color.Color
	Border color.Color
}

// BadgeThemeColors agrupa los colores de todas las variantes de badges de la aplicación.
type BadgeThemeColors struct {
	Success BadgeVariantColors
	Danger  BadgeVariantColors
	Info    BadgeVariantColors
	Warning BadgeVariantColors
	Purple  BadgeVariantColors
	Neutral BadgeVariantColors
}

// ButtonStateColors define los estados visuales (normal, hover, press, disabled) de un tipo de botón.
type ButtonStateColors struct {
	Bg             color.Color
	Fg             color.Color
	Border         color.Color
	HoverBg        color.Color
	HoverFg        color.Color
	HoverBorder    color.Color
	PressBg        color.Color
	PressFg        color.Color
	PressBorder    color.Color
	DisabledBg     color.Color
	DisabledFg     color.Color
	DisabledBorder color.Color
}

// ButtonThemeColors agrupa los colores para los tres tipos de botones de ModernButton.
type ButtonThemeColors struct {
	Primary ButtonStateColors
	Default ButtonStateColors
	Danger  ButtonStateColors
}

// TableThemeColors define los colores para tablas y vistas de datos (ModernBrowse).
type TableThemeColors struct {
	HeaderBg      color.Color
	HeaderFg      color.Color
	HeaderHoverFg color.Color
	SortIndicator color.Color
	RowEvenBg     color.Color
	RowOddBg      color.Color
	RowHoverBg    color.Color
	RowSelectedBg color.Color
	Border        color.Color
}

// ChartThemeColors define la paleta para gráficos y sus tooltips (ModernLineChart).
type ChartThemeColors struct {
	TooltipBg       color.Color
	TooltipBorder   color.Color
	TooltipTitle    color.Color
	TooltipSubtitle color.Color
	TooltipDivider  color.Color
	TooltipText     color.Color
	TooltipMuted    color.Color
	GuideLine       color.Color
	BadgeBg         color.Color
	BadgeText       color.Color
}

// AvatarThemeColors define los colores para avatares e iniciales de usuario.
type AvatarThemeColors struct {
	Bg color.Color
	Fg color.Color
}

// KPICardVariantColors define el juego completo de colores para una variante de tarjeta KPI Hero.
type KPICardVariantColors struct {
	Accent   color.Color // Franja lateral y texto del rótulo
	Value    color.Color // Cifra numérica principal
	Bg       color.Color // Fondo suave tintado
	Border   color.Color // Borde perimetral exterior
	Subtitle color.Color // Subtítulo o texto explicativo al pie
}

// KPICardThemeColors agrupa las 4 variantes estándar de tarjetas KPI Hero en GoJubila.
type KPICardThemeColors struct {
	Emerald KPICardVariantColors // 1. Pensión Neta / Favorable (Verde)
	Blue    KPICardVariantColors // 2. Pensión Bruta / Totales (Azul)
	Violet  KPICardVariantColors // 3. Base Reguladora / Cálculo (Violeta)
	Amber   KPICardVariantColors // 4. % Pensión / Coeficiente / Advertencia (Ámbar)
}

// ThemePalette es la especificación centralizada y modular de un tema en GoJubila.
// Cada tema se define en su propio archivo implementando esta estructura con los colores de cada componente.
type ThemePalette struct {
	Name        string
	DisplayName string
	IsDark      bool

	// Superficies y Fondos
	AppBackground   color.Color
	NavbarBg        color.Color
	CardBg          color.Color
	CardBorder      color.Color
	SidebarBg       color.Color
	SidebarShadedBg color.Color
	FilterDrawerBg  color.Color
	InsetPanelBg    color.Color

	// Bordes y Divisores Generales
	Border  color.Color
	Divider color.Color

	// Tipografía General
	TextPrimary   color.Color
	TextMuted     color.Color
	NavInactiveFg color.Color

	// Tokens Semánticos de Estado
	Primary color.Color // Color primario de marca / cursor y foco en inputs
	Accent  color.Color
	Success color.Color
	Warning color.Color
	Danger  color.Color

	// Controles e Inputs
	InputBackground color.Color
	InputBorder     color.Color
	Selection       color.Color
	Hover           color.Color
	Disabled        color.Color
	DisabledButton  color.Color

	// Componentes Específicos
	HeaderCard   HeaderCardThemeColors
	Badges       BadgeThemeColors
	Buttons      ButtonThemeColors
	GlassButtons ButtonThemeColors
	Table        TableThemeColors
	Chart        ChartThemeColors
	Avatar       AvatarThemeColors
	KPICards     KPICardThemeColors
}

// Color resuelve un token semántico de Fyne o de GoJubila a partir de la paleta del tema.
func (p *ThemePalette) Color(name fyne.ThemeColorName) color.Color {
	if p == nil {
		return nil
	}
	switch name {
	// Buttons: Primary
	case ColorNameBtnPrimaryBg:
		return p.Buttons.Primary.Bg
	case ColorNameBtnPrimaryFg:
		return p.Buttons.Primary.Fg
	case ColorNameBtnPrimaryBorder:
		return p.Buttons.Primary.Border
	case ColorNameBtnPrimaryHoverBg:
		return p.Buttons.Primary.HoverBg
	case ColorNameBtnPrimaryHoverFg:
		return p.Buttons.Primary.HoverFg
	case ColorNameBtnPrimaryHoverBorder:
		return p.Buttons.Primary.HoverBorder
	case ColorNameBtnPrimaryPressBg:
		return p.Buttons.Primary.PressBg
	case ColorNameBtnPrimaryPressFg:
		return p.Buttons.Primary.PressFg
	case ColorNameBtnPrimaryPressBorder:
		return p.Buttons.Primary.PressBorder
	case ColorNameBtnPrimaryMuted:
		return p.Buttons.Primary.DisabledFg

	// Buttons: Default
	case ColorNameBtnDefaultBg:
		return p.Buttons.Default.Bg
	case ColorNameBtnDefaultFg:
		return p.Buttons.Default.Fg
	case ColorNameBtnDefaultBorder:
		return p.Buttons.Default.Border
	case ColorNameBtnDefaultHoverBg:
		return p.Buttons.Default.HoverBg
	case ColorNameBtnDefaultHoverFg:
		return p.Buttons.Default.HoverFg
	case ColorNameBtnDefaultHoverBorder:
		return p.Buttons.Default.HoverBorder
	case ColorNameBtnDefaultPressBg:
		return p.Buttons.Default.PressBg
	case ColorNameBtnDefaultPressFg:
		return p.Buttons.Default.PressFg
	case ColorNameBtnDefaultPressBorder:
		return p.Buttons.Default.PressBorder

	// Buttons: Danger
	case ColorNameBtnDanger, ColorNameBtnDangerFg:
		return p.Buttons.Danger.Fg
	case ColorNameBtnDangerBg:
		return p.Buttons.Danger.Bg
	case ColorNameBtnDangerBorder:
		return p.Buttons.Danger.Border
	case ColorNameBtnDangerHoverBg:
		return p.Buttons.Danger.HoverBg
	case ColorNameBtnDangerHoverFg:
		return p.Buttons.Danger.HoverFg
	case ColorNameBtnDangerHoverBorder:
		return p.Buttons.Danger.HoverBorder
	case ColorNameBtnDangerPressBg:
		return p.Buttons.Danger.PressBg
	case ColorNameBtnDangerPressFg:
		return p.Buttons.Danger.PressFg
	case ColorNameBtnDangerPressBorder:
		return p.Buttons.Danger.PressBorder
	case ColorNameBtnDangerMuted:
		return p.Buttons.Danger.DisabledFg

	case ColorNameBtnMuted:
		return p.Disabled

	// Superficies y Fondos
	case ColorNameAppBackground:
		return p.AppBackground
	case ColorNameNavbarBg:
		return p.NavbarBg
	case ColorNameCardBg:
		return p.CardBg
	case ColorNameCardBorder:
		return p.CardBorder
	case ColorNameSidebarBg:
		return p.SidebarBg
	case ColorNameSidebarShadedBg:
		return p.SidebarShadedBg
	case ColorNameFilterDrawerBg:
		return p.FilterDrawerBg
	case ColorNameInsetPanelBg:
		return p.InsetPanelBg

	// Bordes y Divisores
	case ColorNameBorder:
		return p.Border
	case ColorNameDivider:
		return p.Divider

	// Tipografía
	case ColorNameTextPrimary:
		return p.TextPrimary
	case ColorNameTextMuted:
		return p.TextMuted
	case ColorNameNavInactiveFg:
		return p.NavInactiveFg

	// Estados Semánticos
	case ColorNameAccent:
		return p.Accent
	case ColorNameSuccess:
		return p.Success
	case ColorNameWarning:
		return p.Warning
	case ColorNameDanger:
		return p.Danger

	// Elementos de Lista (Obras / Cuadros)
	case ColorNameItemID:
		if p.IsDark {
			return ColorBlue400
		}
		return ColorBlue600
	case ColorNameItemPrice:
		if p.IsDark {
			return ColorEmerald400
		}
		return ColorEmerald700
	case ColorNameItemThumbBg:
		if p.IsDark {
			return ColorSlate800
		}
		return ColorSlate200
	case ColorNameItemThumbBorder:
		if p.IsDark {
			return ColorWhiteAlpha30
		}
		return ColorBlackAlpha20

	// Tokens Estándar de Fyne
	case theme.ColorNamePrimary:
		if p.Primary != nil {
			return p.Primary
		}
		if p.Buttons.Primary.Bg != nil {
			return p.Buttons.Primary.Bg
		}
		return ColorGreen700
	case theme.ColorNameFocus:
		if p.Primary != nil {
			return p.Primary
		}
		return p.Accent
	case theme.ColorNameForegroundOnPrimary, theme.ColorNameForegroundOnError:
		return ColorWhite
	case theme.ColorNameBackground:
		return p.AppBackground
	case theme.ColorNameMenuBackground, theme.ColorNameOverlayBackground:
		return p.CardBg
	case theme.ColorNameError:
		return p.Danger
	case theme.ColorNameButton:
		return p.Buttons.Default.Bg
	case theme.ColorNameHover:
		return p.Hover
	case theme.ColorNameSelection:
		return p.Selection
	case theme.ColorNameDisabled:
		return p.Disabled
	case theme.ColorNameForeground:
		return p.TextPrimary
	case theme.ColorNameInputBorder:
		return p.InputBorder
	case theme.ColorNameDisabledButton:
		return p.DisabledButton
	case theme.ColorNameInputBackground:
		return p.InputBackground
	case theme.ColorNameSeparator:
		return p.Divider
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Registro Global y Resolución de Paletas
// ─────────────────────────────────────────────────────────────────────────────

var (
	themeRegistryMu sync.RWMutex
	themeRegistry   = make(map[string]*ThemePalette)
	themeOrder      []string
)

// RegisterTheme registra una nueva paleta de tema en el catálogo global.
func RegisterTheme(p *ThemePalette) {
	if p == nil || p.Name == "" {
		return
	}
	themeRegistryMu.Lock()
	defer themeRegistryMu.Unlock()

	key := strings.ToLower(p.Name)
	if _, exists := themeRegistry[key]; !exists {
		themeOrder = append(themeOrder, p.Name)
	}
	themeRegistry[key] = p
}

// UnregisterTheme elimina una paleta del registro global (útil para pruebas o recarga de temas).
func UnregisterTheme(name string) {
	themeRegistryMu.Lock()
	defer themeRegistryMu.Unlock()

	key := strings.ToLower(name)
	delete(themeRegistry, key)

	for i, n := range themeOrder {
		if strings.EqualFold(n, name) {
			themeOrder = append(themeOrder[:i], themeOrder[i+1:]...)
			break
		}
	}
}

// GetThemePalette devuelve la paleta asociada al identificador de tema.
// Si no existe, recurre a ThemeLight como alternativa segura por defecto.
func GetThemePalette(name string) *ThemePalette {
	themeRegistryMu.RLock()
	defer themeRegistryMu.RUnlock()

	key := strings.ToLower(strings.TrimSpace(name))
	if p, ok := themeRegistry[key]; ok {
		return p
	}
	if p, ok := themeRegistry[ThemeLight]; ok {
		return p
	}
	// Si todavía no hay nada registrado, crear un fallback seguro
	return nil
}

// GetCurrentThemePalette devuelve la paleta del tema activo en la aplicación.
func GetCurrentThemePalette() *ThemePalette {
	name := GetCurrentThemeName()
	p := GetThemePalette(name)
	if p == nil {
		// Fallback inicial seguro antes de que inicialicen los paquetes
		p = GetThemePalette(ThemeLight)
	}
	return p
}

// AvailableThemes devuelve la lista ordenada de temas registrados.
func AvailableThemes() []string {
	themeRegistryMu.RLock()
	defer themeRegistryMu.RUnlock()

	if len(themeOrder) == 0 {
		return []string{ThemeLight, ThemeDark}
	}
	res := make([]string, len(themeOrder))
	copy(res, themeOrder)
	return res
}

// GetNextTheme devuelve el siguiente tema en el ciclo de rotación nativo (Claro <-> Oscuro).
func GetNextTheme(current string) string {
	if strings.EqualFold(current, ThemeLight) {
		return ThemeDark
	}
	return ThemeLight
}
