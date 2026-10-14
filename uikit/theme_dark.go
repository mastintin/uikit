package uikit

import (
	"image/color"
)

func init() {
	RegisterTheme(NewDarkPalette())
}

// NewDarkPalette construye la definición de colores de componentes para el tema "Dark".
// Paleta slate oscura de alta profundidad con acentos índigo/sky y contrastes nítidos.
func NewDarkPalette() *ThemePalette {
	return &ThemePalette{
		Name:        ThemeDark,
		DisplayName: "Oscuro (Moderno)",
		IsDark:      true,

		// Superficies y Fondos
		AppBackground:   ColorSlate950, // #020617
		NavbarBg:        ColorSlate900, // #0f172a
		CardBg:          ColorSlate900, // #0f172a
		CardBorder:      ColorSlate800, // #1e293b
		SidebarBg:       ColorSlate900,
		SidebarShadedBg: ColorSlate950,
		FilterDrawerBg:  ColorSlate900,
		InsetPanelBg:    ColorSlate800, // #1e293b

		// Bordes y Divisores
		Border:  ColorSlate700, // #334155
		Divider: ColorSlate700,

		// Tipografía
		TextPrimary:   ColorSlate50,  // #f8fafc
		TextMuted:     ColorSlate400, // #94a3b8
		NavInactiveFg: ColorSky300,   // #93c5fd

		// Estados Semánticos
		Primary: ColorIndigo400,  // #818cf8 - Acento visible para cursor y foco en modo oscuro
		Accent:  ColorIndigo400,  // #818cf8
		Success: ColorEmerald400, // #34d399
		Warning: ColorAmber400,   // #fbbf24
		Danger:  ColorRed400,     // #f87171

		// Controles e Inputs
		InputBackground: ColorSlate800,
		InputBorder:     ColorSlate700,
		Selection:       ColorIndigo800,
		Hover:           ColorSlate700,
		Disabled:        ColorSlate600,
		DisabledButton:  ColorSlate800,

		// ── Componente: ModernHeaderCard ─────────────────────────────────────
		HeaderCard: HeaderCardThemeColors{
			BannerTop:     color.NRGBA{R: 28, G: 82, B: 57, A: 255}, // #1c5239
			BannerBottom:  color.NRGBA{R: 15, G: 56, B: 38, A: 255}, // #0f3826
			BannerDivider: HexColor("#0b291c"),
			BodyBg:        ColorSlate900,
			Border:        ColorSlate800,
			TitleColor:    ColorWhite,
			SubtitleColor: color.NRGBA{R: 167, G: 243, B: 208, A: 220},
		},

		// ── Componente: ModernBadge ──────────────────────────────────────────
		Badges: BadgeThemeColors{
			Success: BadgeVariantColors{
				Bg:     color.NRGBA{R: 6, G: 78, B: 59, A: 200},
				Text:   color.NRGBA{R: 134, G: 239, B: 172, A: 255},
				Border: color.NRGBA{R: 16, G: 185, B: 129, A: 255},
			},
			Danger: BadgeVariantColors{
				Bg:     color.NRGBA{R: 127, G: 29, B: 29, A: 200},
				Text:   ColorRed300,
				Border: color.NRGBA{R: 239, G: 68, B: 68, A: 255},
			},
			Info: BadgeVariantColors{
				Bg:     color.NRGBA{R: 30, G: 58, B: 138, A: 200},
				Text:   ColorSky300,
				Border: ColorBlue500,
			},
			Warning: BadgeVariantColors{
				Bg:     color.NRGBA{R: 120, G: 53, B: 15, A: 200},
				Text:   color.NRGBA{R: 253, G: 230, B: 138, A: 255},
				Border: color.NRGBA{R: 245, G: 158, B: 11, A: 255},
			},
			Purple: BadgeVariantColors{
				Bg:     color.NRGBA{R: 88, G: 28, B: 135, A: 200},
				Text:   color.NRGBA{R: 216, G: 180, B: 254, A: 255},
				Border: color.NRGBA{R: 168, G: 85, B: 247, A: 255},
			},
			Neutral: BadgeVariantColors{
				Bg:     ColorSlate800,
				Text:   ColorSlate300,
				Border: ColorSlate700,
			},
		},

		// ── Componente: ModernButton ─────────────────────────────────────────
		Buttons: ButtonThemeColors{
			Primary: ButtonStateColors{
				Bg:             ColorIndigo600, // #4f46e5
				Fg:             ColorWhite,
				Border:         ColorIndigo500,
				HoverBg:        ColorIndigo700,
				HoverFg:        ColorWhite,
				HoverBorder:    ColorIndigo400,
				PressBg:        ColorIndigo800,
				PressFg:        ColorWhite,
				PressBorder:    ColorIndigo500,
				DisabledBg:     ColorSlate800,
				DisabledFg:     ColorSlate400,
				DisabledBorder: ColorSlate700,
			},
			Default: ButtonStateColors{
				Bg:             ColorSlate800,
				Fg:             ColorSlate100,
				Border:         ColorSlate700,
				HoverBg:        ColorSlate700,
				HoverFg:        ColorWhite,
				HoverBorder:    ColorSlate600,
				PressBg:        ColorSlate900,
				PressFg:        ColorWhite,
				PressBorder:    ColorSlate600,
				DisabledBg:     ColorSlate900,
				DisabledFg:     ColorSlate600,
				DisabledBorder: ColorSlate800,
			},
			Danger: ButtonStateColors{
				Bg:             ColorSlate800,
				Fg:             ColorRed400,
				Border:         ColorRed400,
				HoverBg:        ColorSlate700,
				HoverFg:        ColorWhite,
				HoverBorder:    ColorRed300,
				PressBg:        ColorSlate900,
				PressFg:        ColorWhite,
				PressBorder:    ColorRed400,
				DisabledBg:     ColorSlate900,
				DisabledFg:     ColorSlate600,
				DisabledBorder: ColorSlate800,
			},
		},

		// ── Componente: ModernButton (Modo Cristal / Glass) ──────────────────
		GlassButtons: ButtonThemeColors{
			Primary: ButtonStateColors{
				Bg:             color.NRGBA{R: 79, G: 70, B: 229, A: 45},
				Fg:             ColorIndigo400,
				Border:         color.NRGBA{R: 129, G: 140, B: 248, A: 130},
				HoverBg:        color.NRGBA{R: 99, G: 102, B: 241, A: 75},
				HoverFg:        ColorWhite,
				HoverBorder:    color.NRGBA{R: 129, G: 140, B: 248, A: 200},
				PressBg:        color.NRGBA{R: 67, G: 56, B: 202, A: 110},
				PressFg:        ColorWhite,
				PressBorder:    color.NRGBA{R: 99, G: 102, B: 241, A: 220},
				DisabledBg:     color.NRGBA{R: 79, G: 70, B: 229, A: 20},
				DisabledFg:     ColorSlate400,
				DisabledBorder: color.NRGBA{R: 99, G: 102, B: 241, A: 50},
			},
			Default: ButtonStateColors{
				Bg:             color.NRGBA{R: 30, G: 41, B: 59, A: 140},
				Fg:             ColorSlate100,
				Border:         color.NRGBA{R: 71, G: 85, B: 105, A: 160},
				HoverBg:        color.NRGBA{R: 51, G: 65, B: 85, A: 180},
				HoverFg:        ColorWhite,
				HoverBorder:    ColorSlate600,
				PressBg:        color.NRGBA{R: 15, G: 23, B: 42, A: 210},
				PressFg:        ColorWhite,
				PressBorder:    ColorSlate700,
				DisabledBg:     color.Transparent,
				DisabledFg:     ColorSlate600,
				DisabledBorder: color.Transparent,
			},
			Danger: ButtonStateColors{
				Bg:             color.NRGBA{R: 220, G: 38, B: 38, A: 38},
				Fg:             ColorRed400,
				Border:         color.NRGBA{R: 220, G: 38, B: 38, A: 120},
				HoverBg:        color.NRGBA{R: 220, G: 38, B: 38, A: 68},
				HoverFg:        ColorWhite,
				HoverBorder:    ColorRed300,
				PressBg:        color.NRGBA{R: 185, G: 28, B: 28, A: 100},
				PressFg:        ColorWhite,
				PressBorder:    ColorRed400,
				DisabledBg:     color.NRGBA{R: 220, G: 38, B: 38, A: 15},
				DisabledFg:     ColorSlate600,
				DisabledBorder: color.NRGBA{R: 220, G: 38, B: 38, A: 40},
			},
		},

		// ── Componente: ModernBrowse / Tablas ────────────────────────────────
		Table: TableThemeColors{
			HeaderBg:      HexColor("#1e293b"),
			HeaderFg:      HexColor("#94a3b8"),
			HeaderHoverFg: HexColor("#f8fafc"),
			SortIndicator: HexColor("#38bdf8"),
			RowEvenBg:     HexColor("#0f172a"),
			RowOddBg:      HexColor("#1e293b"),
			RowHoverBg:    HexColor("#334155"),
			RowSelectedBg: HexColor("#1e3a8a"),
			Border:        ColorSlate700,
		},

		// ── Componente: ModernLineChart / Gráficos ───────────────────────────
		Chart: ChartThemeColors{
			TooltipBg:       HexColor("#0f172a"),
			TooltipBorder:   HexColor("#38bdf8"),
			TooltipTitle:    ColorWhite,
			TooltipSubtitle: HexColor("#94a3b8"),
			TooltipDivider:  HexColor("#334155"),
			TooltipText:     HexColor("#f8fafc"),
			TooltipMuted:    HexColor("#94a3b8"),
			GuideLine:       HexColor("#0284c7"),
			BadgeBg:         HexColor("#0f172a"),
			BadgeText:       ColorWhite,
		},

		// ── Componente: Avatar ───────────────────────────────────────────────
		Avatar: AvatarThemeColors{
			Bg: HexColor("#1e293b"),
			Fg: HexColor("#60a5fa"),
		},

		// ── Componente: KPI Cards (Hero metric cards con franja) ─────────────
		KPICards: KPICardThemeColors{
			Emerald: KPICardVariantColors{
				Accent:   HexColor("#10b981"),
				Value:    HexColor("#10b981"),
				Bg:       HexColor("#064e3b"),
				Border:   HexColor("#065f46"),
				Subtitle: HexColor("#6ee7b7"),
			},
			Blue: KPICardVariantColors{
				Accent:   HexColor("#3b82f6"),
				Value:    HexColor("#3b82f6"),
				Bg:       HexColor("#1e3a8a"),
				Border:   HexColor("#1d4ed8"),
				Subtitle: HexColor("#93c5fd"),
			},
			Violet: KPICardVariantColors{
				Accent:   HexColor("#818cf8"),
				Value:    HexColor("#818cf8"),
				Bg:       HexColor("#1e1b4b"),
				Border:   HexColor("#4338ca"),
				Subtitle: HexColor("#a5b4fc"),
			},
			Amber: KPICardVariantColors{
				Accent:   HexColor("#f59e0b"),
				Value:    HexColor("#f59e0b"),
				Bg:       HexColor("#451a03"),
				Border:   HexColor("#78350f"),
				Subtitle: HexColor("#fcd34d"),
			},
		},
	}
}
