package uikit

import (
	"image/color"
)

func init() {
	RegisterTheme(NewLightPalette())
}

// NewLightPalette construye la definición de colores de componentes para el tema "Light".
// Tema claro estándar con superficies limpias y acentos azul corporativo Tailwind.
func NewLightPalette() *ThemePalette {
	return &ThemePalette{
		Name:        ThemeLight,
		DisplayName: "Claro (Estándar)",
		IsDark:      false,

		// Superficies y Fondos
		AppBackground:   ColorSlate50,  // #f8fafc
		NavbarBg:        ColorWhite,    // #ffffff
		CardBg:          ColorWhite,    // #ffffff
		CardBorder:      ColorSlate300, // #cbd5e1
		SidebarBg:       ColorWhite,
		SidebarShadedBg: ColorSlate50,
		FilterDrawerBg:  ColorSlate100,
		InsetPanelBg:    ColorSlate50,

		// Bordes y Divisores
		Border:  ColorSlate200, // #e2e8f0
		Divider: ColorSlate200,

		// Tipografía
		TextPrimary:   ColorSlate900, // #0f172a
		TextMuted:     ColorSlate700, // #334155
		NavInactiveFg: ColorBlue600,  // #2563eb

		// Estados Semánticos
		Primary: HexColor("#107c41"), // Verde corporativo Nice: cursor y foco de alta visibilidad y contraste
		Accent:  ColorBlue700,         // #1d4ed8
		Success: ColorGreen700,        // #15803d
		Warning: ColorAmber400,        // #fbbf24
		Danger:  ColorWarningFg,

		// Controles e Inputs
		InputBackground: ColorWhite,
		InputBorder:     ColorSlate300,
		Selection:       ColorBlue100, // #dbeafe
		Hover:           ColorSlate100,
		Disabled:        ColorSlate400,
		DisabledButton:  ColorSlate100,

		// ── Componente: ModernHeaderCard ─────────────────────────────────────
		HeaderCard: HeaderCardThemeColors{
			BannerTop:     color.NRGBA{R: 10, G: 98, B: 75, A: 255}, // #0a624b
			BannerBottom:  color.NRGBA{R: 6, G: 78, B: 59, A: 255},  // #064e3b
			BannerDivider: HexColor("#022c22"),
			BodyBg:        ColorWhite,
			Border:        ColorSlate300,
			TitleColor:    ColorWhite,
			SubtitleColor: color.NRGBA{R: 167, G: 243, B: 208, A: 220},
		},

		// ── Componente: ModernBadge ──────────────────────────────────────────
		Badges: BadgeThemeColors{
			Success: BadgeVariantColors{
				Bg:     ColorActiveBg,
				Text:   ColorActiveFg,
				Border: ColorActiveBorder,
			},
			Danger: BadgeVariantColors{
				Bg:     ColorWarningBg,
				Text:   ColorWarningFg,
				Border: ColorWarningBorder,
			},
			Info: BadgeVariantColors{
				Bg:     ColorBlue50,
				Text:   ColorBlue600,
				Border: ColorBlue200,
			},
			Warning: BadgeVariantColors{
				Bg:     HexColor("#FEF3C7"),
				Text:   HexColor("#92400E"),
				Border: HexColor("#FDE68A"),
			},
			Purple: BadgeVariantColors{
				Bg:     HexColor("#f3e8ff"),
				Text:   HexColor("#7e22ce"),
				Border: HexColor("#e9d5ff"),
			},
			Neutral: BadgeVariantColors{
				Bg:     ColorSlate100,
				Text:   ColorSlate600,
				Border: ColorSlate200,
			},
		},

		// ── Componente: ModernButton ─────────────────────────────────────────
		Buttons: ButtonThemeColors{
			Primary: ButtonStateColors{
				Bg:             ColorBlue50, // #eff6ff
				Fg:             ColorBlue600,
				Border:         ColorBlue200,
				HoverBg:        ColorBlue100,
				HoverFg:        ColorBlue700,
				HoverBorder:    ColorBlue500,
				PressBg:        ColorBlue200,
				PressFg:        ColorBlue700,
				PressBorder:    ColorBlue600,
				DisabledBg:     color.NRGBA{R: 37, G: 99, B: 235, A: 15},
				DisabledFg:     ColorBlue200,
				DisabledBorder: color.NRGBA{R: 37, G: 99, B: 235, A: 40},
			},
			Default: ButtonStateColors{
				Bg:             ColorWhite,
				Fg:             ColorSlate900,
				Border:         ColorSlate300,
				HoverBg:        ColorSlate100,
				HoverFg:        ColorSlate900,
				HoverBorder:    ColorSlate400,
				PressBg:        ColorSlate200,
				PressFg:        ColorSlate900,
				PressBorder:    ColorSlate400,
				DisabledBg:     ColorSlate100,
				DisabledFg:     ColorSlate400,
				DisabledBorder: ColorSlate200,
			},
			Danger: ButtonStateColors{
				Bg:             ColorWarningBg,
				Fg:             ColorWarningFg,
				Border:         ColorWarningBorder,
				HoverBg:        ColorWarningHoverBg,
				HoverFg:        ColorWarningHoverFg,
				HoverBorder:    ColorWarningHoverBorder,
				PressBg:        ColorWarningPressBg,
				PressFg:        ColorWarningPressFg,
				PressBorder:    ColorWarningPressBorder,
				DisabledBg:     ColorSlate100,
				DisabledFg:     ColorSlate400,
				DisabledBorder: ColorSlate200,
			},
		},

		// ── Componente: ModernButton (Modo Cristal / Glass) ──────────────────
		GlassButtons: ButtonThemeColors{
			Primary: ButtonStateColors{
				Bg:             color.NRGBA{R: 37, G: 99, B: 235, A: 38},
				Fg:             ColorBlue600,
				Border:         color.NRGBA{R: 37, G: 99, B: 235, A: 120},
				HoverBg:        color.NRGBA{R: 37, G: 99, B: 235, A: 68},
				HoverFg:        ColorBlue700,
				HoverBorder:    color.NRGBA{R: 29, G: 78, B: 216, A: 190},
				PressBg:        color.NRGBA{R: 29, G: 78, B: 216, A: 100},
				PressFg:        ColorBlue700,
				PressBorder:    color.NRGBA{R: 29, G: 78, B: 216, A: 220},
				DisabledBg:     color.NRGBA{R: 37, G: 99, B: 235, A: 15},
				DisabledFg:     ColorBlue200,
				DisabledBorder: color.NRGBA{R: 37, G: 99, B: 235, A: 40},
			},
			Default: ButtonStateColors{
				Bg:             color.NRGBA{R: 255, G: 255, B: 255, A: 140},
				Fg:             ColorSlate900,
				Border:         color.NRGBA{R: 203, G: 213, B: 225, A: 160},
				HoverBg:        color.NRGBA{R: 255, G: 255, B: 255, A: 200},
				HoverFg:        ColorSlate900,
				HoverBorder:    ColorSlate300,
				PressBg:        color.NRGBA{R: 241, G: 245, B: 249, A: 220},
				PressFg:        ColorSlate900,
				PressBorder:    ColorSlate300,
				DisabledBg:     color.Transparent,
				DisabledFg:     ColorSlate400,
				DisabledBorder: color.Transparent,
			},
			Danger: ButtonStateColors{
				Bg:             color.NRGBA{R: 245, G: 203, B: 167, A: 70},  // #F5CBA7 translúcido suave
				Fg:             HexColor("#B03A2E"),
				Border:         color.NRGBA{R: 208, G: 90, B: 63, A: 130},  // #D05A3F
				HoverBg:        color.NRGBA{R: 250, G: 219, B: 216, A: 120}, // #FADBD8
				HoverFg:        HexColor("#943126"),
				HoverBorder:    color.NRGBA{R: 176, G: 58, B: 46, A: 180},  // #B03A2E
				PressBg:        color.NRGBA{R: 237, G: 187, B: 153, A: 160},
				PressFg:        HexColor("#78281F"),
				PressBorder:    color.NRGBA{R: 148, G: 49, B: 38, A: 210},
				DisabledBg:     color.NRGBA{R: 245, G: 203, B: 167, A: 25},
				DisabledFg:     HexColor("#D05A3F"),
				DisabledBorder: color.NRGBA{R: 208, G: 90, B: 63, A: 40},
			},
		},

		// ── Componente: ModernBrowse / Tablas ────────────────────────────────
		Table: TableThemeColors{
			HeaderBg:      HexColor("#334155"),
			HeaderFg:      ColorWhite,
			HeaderHoverFg: HexColor("#bae6fd"),
			SortIndicator: HexColor("#bae6fd"),
			RowEvenBg:     ColorWhite,
			RowOddBg:      HexColor("#f8fafc"),
			RowHoverBg:    HexColor("#f1f5f9"),
			RowSelectedBg: HexColor("#e2e8f0"),
			Border:        ColorSlate200,
		},

		// ── Componente: ModernLineChart / Gráficos ───────────────────────────
		Chart: ChartThemeColors{
			TooltipBg:       ColorWhite,
			TooltipBorder:   ColorBlue600,
			TooltipTitle:    ColorSlate900,
			TooltipSubtitle: ColorSlate500,
			TooltipDivider:  ColorSlate200,
			TooltipText:     ColorSlate900,
			TooltipMuted:    ColorSlate500,
			GuideLine:       ColorBlue600,
			BadgeBg:         ColorBlue600,
			BadgeText:       ColorWhite,
		},

		// ── Componente: Avatar ───────────────────────────────────────────────
		Avatar: AvatarThemeColors{
			Bg: HexColor("#dbeafe"),
			Fg: HexColor("#1d4ed8"),
		},

		// ── Componente: KPI Cards (Hero metric cards con franja) ─────────────
		KPICards: KPICardThemeColors{
			Emerald: KPICardVariantColors{
				Accent:   HexColor("#10b981"),
				Value:    HexColor("#047857"),
				Bg:       HexColor("#ecfdf5"),
				Border:   HexColor("#a7f3d0"),
				Subtitle: HexColor("#065f46"),
			},
			Blue: KPICardVariantColors{
				Accent:   HexColor("#3b82f6"),
				Value:    HexColor("#1e3a8a"),
				Bg:       HexColor("#eff6ff"),
				Border:   HexColor("#bfdbfe"),
				Subtitle: HexColor("#1e40af"),
			},
			Violet: KPICardVariantColors{
				Accent:   HexColor("#818cf8"),
				Value:    HexColor("#4338ca"),
				Bg:       HexColor("#eef2ff"),
				Border:   HexColor("#c7d2fe"),
				Subtitle: HexColor("#4338ca"),
			},
			Amber: KPICardVariantColors{
				Accent:   HexColor("#f59e0b"),
				Value:    HexColor("#92400e"),
				Bg:       HexColor("#fffbeb"),
				Border:   HexColor("#fde68a"),
				Subtitle: HexColor("#78350f"),
			},
		},
	}
}
