package uikit

import (
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// ParseThemeCSS interpreta el contenido de un archivo .css con design tokens y genera un *ThemePalette.
// Soporta variables CSS con o sin prefijo '--', comentarios multilínea /* */ y de línea //,
// y deduce automáticamente colores secundarios si no se declaran todos los tokens.
func ParseThemeCSS(content string, fallbackName string) (*ThemePalette, error) {
	pal, _, err := ParseThemeCSSWithReport(content, fallbackName)
	return pal, err
}

// ParseThemeCSSWithReport analiza el CSS y devuelve la paleta, advertencias de sintaxis detectadas y cualquier error fatal.
func ParseThemeCSSWithReport(content string, fallbackName string) (*ThemePalette, []string, error) {
	// 1. Eliminar comentarios multilínea /* ... */
	reBlockComment := regexp.MustCompile(`(?s)/\*.*?\*/`)
	content = reBlockComment.ReplaceAllString(content, "")

	var warnings []string
	lines := strings.Split(content, "\n")
	props := make(map[string]string)

	for i, line := range lines {
		lineNum := i + 1
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") {
			continue
		}
		// Ignorar bloques envolventes :root { o }
		if strings.HasPrefix(line, ":root") || line == "{" || line == "}" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			warnings = append(warnings, fmt.Sprintf("línea %d: sintaxis inválida en '%s' (falta ':')", lineNum, line))
			continue
		}

		key := strings.TrimSpace(parts[0])
		key = strings.TrimPrefix(key, "--")
		key = strings.ToLower(key)

		val := strings.TrimSpace(parts[1])
		val = strings.TrimSuffix(val, ";")
		val = strings.TrimSpace(val)

		// Quitar comentarios al final de la línea
		if idx := strings.Index(val, "//"); idx >= 0 {
			val = strings.TrimSpace(val[:idx])
		}
		val = strings.TrimSuffix(val, ";")
		val = strings.TrimSpace(val)
		val = strings.Trim(val, `"'`)

		props[key] = val
	}

	if len(props) == 0 {
		return nil, warnings, fmt.Errorf("el archivo CSS no contiene definiciones de estilo válidas")
	}

	// 2. Resolución de variables CSS cruzadas: var(--nombre) o var(--nombre, fallback)
	reVar := regexp.MustCompile(`var\(\s*--?([a-zA-Z0-9_-]+)\s*(?:,\s*([^)]+))?\s*\)`)
	for iter := 0; iter < 5; iter++ {
		changed := false
		for k, v := range props {
			if strings.Contains(v, "var(") {
				newV := reVar.ReplaceAllStringFunc(v, func(match string) string {
					sub := reVar.FindStringSubmatch(match)
					if len(sub) >= 2 {
						varName := strings.ToLower(sub[1])
						if replacement, ok := props[varName]; ok && replacement != "" {
							return replacement
						}
						if len(sub) >= 3 && sub[2] != "" {
							return strings.TrimSpace(sub[2])
						}
					}
					return match
				})
				if newV != v {
					props[k] = newV
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}

	for k, v := range props {
		if strings.Contains(v, "var(") {
			warnings = append(warnings, fmt.Sprintf("--%s: contiene variable no definida '%s'", k, v))
		}
	}

	name := fallbackName
	if val, ok := props["name"]; ok && val != "" {
		name = val
	}
	name = strings.TrimSuffix(name, filepath.Ext(name))
	name = strings.ToLower(name)
	// Si el nombre colisiona con los temas nativos integrados ("dark", "light"),
	// añadir el sufijo "-css" para que se registre como tema externo independiente y aparezca en el desplegable.
	if name == ThemeDark || name == ThemeLight {
		name = name + "-css"
	}

	displayName := name
	if val, ok := props["display-name"]; ok && val != "" {
		displayName = val
	} else if val, ok := props["title"]; ok && val != "" {
		displayName = val
	}

	isDark := false
	if val, ok := props["is-dark"]; ok {
		isDark = parseBool(val)
	} else if val, ok := props["dark"]; ok {
		isDark = parseBool(val)
	} else if val, ok := props["mode"]; ok && strings.ToLower(val) == "dark" {
		isDark = true
	} else if val, ok := props["type"]; ok && strings.ToLower(val) == "dark" {
		isDark = true
	}

	// Base palette sobre la que aplicar overrides
	var base *ThemePalette
	if isDark {
		base = NewDarkPalette()
	} else {
		base = NewLightPalette()
	}

	pal := *base
	pal.Name = strings.ToLower(name)
	pal.DisplayName = displayName
	pal.IsDark = isDark

	// Parser de colores helper
	parseCSSColor := func(val string) color.Color {
		val = strings.TrimSpace(val)
		if val == "" {
			return nil
		}
		if strings.HasPrefix(val, "#") {
			c := HexColor(val)
			if c == color.Transparent && strings.ToLower(val) != "#0000" && strings.ToLower(val) != "#00000000" && strings.ToLower(val) != "#000" {
				return nil
			}
			return c
		}
		lower := strings.ToLower(val)
		if strings.HasPrefix(lower, "rgba(") || strings.HasPrefix(lower, "rgb(") {
			idxOpen := strings.Index(val, "(")
			idxClose := strings.LastIndex(val, ")")
			if idxOpen >= 0 && idxClose > idxOpen {
				parts := strings.Split(val[idxOpen+1:idxClose], ",")
				if len(parts) == 3 {
					r, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
					g, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
					b, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
					return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
				} else if len(parts) == 4 {
					r, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
					g, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
					b, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
					alphaStr := strings.TrimSpace(parts[3])
					var a uint8 = 255
					if strings.HasSuffix(alphaStr, "%") {
						if pct, err := strconv.ParseFloat(strings.TrimSuffix(alphaStr, "%"), 64); err == nil {
							if pct > 100 {
								pct = 100
							} else if pct < 0 {
								pct = 0
							}
							a = uint8((pct / 100.0) * 255.0)
						}
					} else if f, err := strconv.ParseFloat(alphaStr, 64); err == nil {
						if f <= 1.0 && f >= 0.0 {
							a = uint8(f * 255.0)
						} else if f > 1.0 && f <= 255.0 {
							a = uint8(f)
						} else if f > 255.0 {
							a = 255
						}
					}
					return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: a}
				}
			}
		}
		return nil
	}

	parseColor := func(key string, fallback color.Color) color.Color {
		if val, ok := props[key]; ok && val != "" {
			if c := parseCSSColor(val); c != nil {
				return c
			}
			warnings = append(warnings, fmt.Sprintf("--%s: valor '%s' no es un color válido", key, val))
		}
		return fallback
	}

	// Superficies y Fondos
	pal.AppBackground = parseColor("app-background", parseColor("app-bg", pal.AppBackground))
	pal.NavbarBg = parseColor("navbar-bg", pal.NavbarBg)
	pal.CardBg = parseColor("card-bg", pal.CardBg)
	pal.CardBorder = parseColor("card-border", pal.CardBorder)
	pal.SidebarBg = parseColor("sidebar-bg", pal.SidebarBg)
	pal.SidebarShadedBg = parseColor("sidebar-shaded-bg", pal.SidebarShadedBg)
	pal.FilterDrawerBg = parseColor("filter-drawer-bg", pal.FilterDrawerBg)
	pal.InsetPanelBg = parseColor("inset-panel-bg", pal.InsetPanelBg)

	// Bordes y Divisores
	pal.Border = parseColor("border", pal.Border)
	pal.Divider = parseColor("divider", pal.Divider)

	// Tipografía
	pal.TextPrimary = parseColor("text-primary", pal.TextPrimary)
	pal.TextMuted = parseColor("text-muted", pal.TextMuted)
	pal.NavInactiveFg = parseColor("nav-inactive-fg", pal.NavInactiveFg)

	// Semánticos de Estado
	pal.Primary = parseColor("primary", parseColor("primary-color", parseColor("btn-primary-bg", pal.Primary)))
	pal.Accent = parseColor("accent", pal.Accent)
	pal.Success = parseColor("success", pal.Success)
	pal.Warning = parseColor("warning", pal.Warning)
	pal.Danger = parseColor("danger", pal.Danger)

	// Inputs y Controles
	pal.InputBackground = parseColor("input-bg", parseColor("input-background", pal.InputBackground))
	pal.InputBorder = parseColor("input-border", pal.InputBorder)
	pal.Hover = parseColor("hover", pal.Hover)
	pal.Selection = parseColor("selection", pal.Selection)
	pal.Disabled = parseColor("disabled", pal.Disabled)
	pal.DisabledButton = parseColor("disabled-button", pal.DisabledButton)

	// HeaderCard (Tarjeta Ejecutiva / Banner)
	if c := parseColor("header-banner-top", nil); c != nil {
		pal.HeaderCard.BannerTop = ColorToNRGBA(c)
	}
	if c := parseColor("header-banner-bottom", nil); c != nil {
		pal.HeaderCard.BannerBottom = ColorToNRGBA(c)
	}
	pal.HeaderCard.BannerDivider = parseColor("header-banner-divider", parseColor("header-divider", pal.HeaderCard.BannerDivider))
	pal.HeaderCard.BodyBg = parseColor("header-body-bg", pal.HeaderCard.BodyBg)
	pal.HeaderCard.Border = parseColor("header-border", pal.HeaderCard.Border)
	pal.HeaderCard.TitleColor = parseColor("header-title", pal.HeaderCard.TitleColor)
	pal.HeaderCard.SubtitleColor = parseColor("header-subtitle", pal.HeaderCard.SubtitleColor)

	// Badges
	pal.Badges.Success.Bg = parseColor("badge-success-bg", pal.Badges.Success.Bg)
	pal.Badges.Success.Text = parseColor("badge-success-text", pal.Badges.Success.Text)
	pal.Badges.Success.Border = parseColor("badge-success-border", pal.Badges.Success.Border)

	pal.Badges.Danger.Bg = parseColor("badge-danger-bg", pal.Badges.Danger.Bg)
	pal.Badges.Danger.Text = parseColor("badge-danger-text", pal.Badges.Danger.Text)
	pal.Badges.Danger.Border = parseColor("badge-danger-border", pal.Badges.Danger.Border)

	pal.Badges.Info.Bg = parseColor("badge-info-bg", pal.Badges.Info.Bg)
	pal.Badges.Info.Text = parseColor("badge-info-text", pal.Badges.Info.Text)
	pal.Badges.Info.Border = parseColor("badge-info-border", pal.Badges.Info.Border)

	pal.Badges.Warning.Bg = parseColor("badge-warning-bg", pal.Badges.Warning.Bg)
	pal.Badges.Warning.Text = parseColor("badge-warning-text", pal.Badges.Warning.Text)
	pal.Badges.Warning.Border = parseColor("badge-warning-border", pal.Badges.Warning.Border)

	pal.Badges.Purple.Bg = parseColor("badge-purple-bg", pal.Badges.Purple.Bg)
	pal.Badges.Purple.Text = parseColor("badge-purple-text", pal.Badges.Purple.Text)
	pal.Badges.Purple.Border = parseColor("badge-purple-border", pal.Badges.Purple.Border)

	pal.Badges.Neutral.Bg = parseColor("badge-neutral-bg", pal.Badges.Neutral.Bg)
	pal.Badges.Neutral.Text = parseColor("badge-neutral-text", pal.Badges.Neutral.Text)
	pal.Badges.Neutral.Border = parseColor("badge-neutral-border", pal.Badges.Neutral.Border)

	// Botones
	parseButtonState := func(prefix string, s *ButtonStateColors) {
		s.Bg = parseColor(prefix+"-bg", s.Bg)
		s.Fg = parseColor(prefix+"-fg", s.Fg)
		s.Border = parseColor(prefix+"-border", s.Border)

		s.HoverBg = parseColor(prefix+"-hover-bg", s.HoverBg)
		s.HoverFg = parseColor(prefix+"-hover-fg", s.HoverFg)
		s.HoverBorder = parseColor(prefix+"-hover-border", s.HoverBorder)

		s.PressBg = parseColor(prefix+"-press-bg", s.PressBg)
		s.PressFg = parseColor(prefix+"-press-fg", s.PressFg)
		s.PressBorder = parseColor(prefix+"-press-border", s.PressBorder)

		s.DisabledBg = parseColor(prefix+"-disabled-bg", s.DisabledBg)
		s.DisabledFg = parseColor(prefix+"-disabled-fg", s.DisabledFg)
		s.DisabledBorder = parseColor(prefix+"-disabled-border", s.DisabledBorder)
	}

	parseButtonState("btn-primary", &pal.Buttons.Primary)
	parseButtonState("btn-default", &pal.Buttons.Default)
	parseButtonState("btn-danger", &pal.Buttons.Danger)

	parseButtonState("btn-glass-primary", &pal.GlassButtons.Primary)
	parseButtonState("btn-glass-default", &pal.GlassButtons.Default)
	parseButtonState("btn-glass-danger", &pal.GlassButtons.Danger)

	// Tablas (ModernBrowse)
	pal.Table.HeaderBg = parseColor("table-header-bg", pal.Table.HeaderBg)
	pal.Table.HeaderFg = parseColor("table-header-fg", pal.Table.HeaderFg)
	pal.Table.HeaderHoverFg = parseColor("table-header-hover-fg", pal.Table.HeaderHoverFg)
	pal.Table.SortIndicator = parseColor("table-sort-indicator", pal.Table.SortIndicator)
	pal.Table.RowEvenBg = parseColor("table-row-even-bg", parseColor("table-row-even", pal.Table.RowEvenBg))
	pal.Table.RowOddBg = parseColor("table-row-odd-bg", parseColor("table-row-odd", pal.Table.RowOddBg))
	pal.Table.RowHoverBg = parseColor("table-row-hover-bg", parseColor("table-row-hover", pal.Table.RowHoverBg))
	pal.Table.RowSelectedBg = parseColor("table-row-selected-bg", parseColor("table-row-selected", pal.Table.RowSelectedBg))
	pal.Table.Border = parseColor("table-border", pal.Table.Border)

	// Gráficos (ModernLineChart)
	pal.Chart.TooltipBg = parseColor("chart-tooltip-bg", pal.Chart.TooltipBg)
	pal.Chart.TooltipBorder = parseColor("chart-tooltip-border", pal.Chart.TooltipBorder)
	pal.Chart.TooltipTitle = parseColor("chart-tooltip-title", pal.Chart.TooltipTitle)
	pal.Chart.TooltipSubtitle = parseColor("chart-tooltip-subtitle", pal.Chart.TooltipSubtitle)
	pal.Chart.TooltipDivider = parseColor("chart-tooltip-divider", pal.Chart.TooltipDivider)
	pal.Chart.TooltipText = parseColor("chart-tooltip-text", pal.Chart.TooltipText)
	pal.Chart.TooltipMuted = parseColor("chart-tooltip-muted", pal.Chart.TooltipMuted)
	pal.Chart.GuideLine = parseColor("chart-guideline", pal.Chart.GuideLine)
	pal.Chart.BadgeBg = parseColor("chart-badge-bg", pal.Chart.BadgeBg)
	pal.Chart.BadgeText = parseColor("chart-badge-text", pal.Chart.BadgeText)

	// Avatar
	pal.Avatar.Bg = parseColor("avatar-bg", pal.Avatar.Bg)
	pal.Avatar.Fg = parseColor("avatar-fg", pal.Avatar.Fg)

	// KPI Cards (4 variantes estándar: emerald, blue, violet, amber)
	parseKPI := func(prefix string, variant *KPICardVariantColors) {
		variant.Accent = parseColor(prefix+"-accent", variant.Accent)
		variant.Value = parseColor(prefix+"-value", variant.Value)
		variant.Bg = parseColor(prefix+"-bg", variant.Bg)
		variant.Border = parseColor(prefix+"-border", variant.Border)
		variant.Subtitle = parseColor(prefix+"-subtitle", variant.Subtitle)
	}
	parseKPI("kpi-emerald", &pal.KPICards.Emerald)
	parseKPI("kpi-blue", &pal.KPICards.Blue)
	parseKPI("kpi-violet", &pal.KPICards.Violet)
	parseKPI("kpi-amber", &pal.KPICards.Amber)

	return &pal, warnings, nil
}

// LoadThemeCSSFileWithReport carga un archivo CSS individual y genera su ThemePalette devolviendo advertencias.
func LoadThemeCSSFileWithReport(path string) (*ThemePalette, []string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("error leyendo archivo %s: %w", filepath.Base(path), err)
	}

	filename := filepath.Base(path)
	return ParseThemeCSSWithReport(string(data), filename)
}

// LoadThemeCSSFile carga un archivo CSS individual y genera su ThemePalette.
func LoadThemeCSSFile(path string) (*ThemePalette, error) {
	pal, _, err := LoadThemeCSSFileWithReport(path)
	return pal, err
}

// ThemeReloadResult agrupa el resultado de registrar temas y los problemas encontrados.
type ThemeReloadResult struct {
	Count    int
	Errors   []string
	Warnings []string
}

// ScanThemesDirDetailed escanea la carpeta y retorna paletas válidas junto con errores y advertencias agrupados.
func ScanThemesDirDetailed(dir string) ([]*ThemePalette, []string, []string) {
	if dir == "" {
		return nil, nil, nil
	}
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return nil, nil, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, []string{fmt.Sprintf("Error leyendo directorio %s: %v", dir, err)}, nil
	}

	var palettes []*ThemePalette
	var allErrors []string
	var allWarnings []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".css" {
			filePath := filepath.Join(dir, entry.Name())
			pal, warns, err := LoadThemeCSSFileWithReport(filePath)
			if err != nil {
				allErrors = append(allErrors, fmt.Sprintf("%s: %v", entry.Name(), err))
			} else if pal != nil {
				palettes = append(palettes, pal)
			}
			for _, w := range warns {
				allWarnings = append(allWarnings, fmt.Sprintf("%s: %s", entry.Name(), w))
			}
		}
	}

	return palettes, allErrors, allWarnings
}

// ScanThemesDir escanea una carpeta en busca de archivos .css y devuelve las paletas cargadas.
func ScanThemesDir(dir string) ([]*ThemePalette, error) {
	palettes, _, _ := ScanThemesDirDetailed(dir)
	return palettes, nil
}

// LoadAndRegisterExternalThemesDetailed escanea y registra temas devolviendo reporte completo de errores y advertencias.
func LoadAndRegisterExternalThemesDetailed(dir string) ThemeReloadResult {
	palettes, errs, warns := ScanThemesDirDetailed(dir)
	registered := 0
	for _, pal := range palettes {
		RegisterTheme(pal)
		registered++
	}
	return ThemeReloadResult{
		Count:    registered,
		Errors:   errs,
		Warnings: warns,
	}
}

// LoadAndRegisterExternalThemes escanea una carpeta y registra todos los temas CSS en el catálogo global de uikit.
// Retorna la cantidad de temas externos registrados con éxito.
func LoadAndRegisterExternalThemes(dir string) int {
	res := LoadAndRegisterExternalThemesDetailed(dir)
	return res.Count
}

func parseBool(s string) bool {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "true" || s == "1" || s == "yes" || s == "si" || s == "dark" {
		return true
	}
	b, _ := strconv.ParseBool(s)
	return b
}


// OpenFolderInOS abre la carpeta especificada en el explorador de archivos del sistema operativo (Finder, Explorer o xdg-open).
func OpenFolderInOS(path string) error {
	if path == "" {
		return fmt.Errorf("ruta vacía")
	}
	_ = os.MkdirAll(path, 0755)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("explorer", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}

// GenerateThemeCSSTemplate genera el código CSS completo y formateado para una nueva plantilla de tema.
func GenerateThemeCSSTemplate(id, displayName string, isDark bool) string {
	if isDark {
		return fmt.Sprintf(`/* ==========================================================================
   GoJubila PRO - Tema: %[1]s
   Generado desde Plantilla Oscura
   ========================================================================== */

:root {
  /* Metadatos */
  --name: %[2]s;
  --display-name: "%[1]s";
  --is-dark: true;

  /* 1. Superficies y Fondos Base */
  --app-background:      #020617; /* Fondo principal de la ventana (Slate 950) */
  --navbar-bg:           #0f172a; /* Barra superior de navegación */
  --card-bg:             #0f172a; /* Fondo de tarjetas y paneles */
  --card-border:         #1e293b; /* Borde de tarjetas (Slate 800) */
  --sidebar-bg:          #0f172a; /* Menú lateral */
  --sidebar-shaded-bg:   #020617; /* Fondo secundario de la barra lateral */
  --filter-drawer-bg:    #0f172a; /* Cajón de filtros */
  --inset-panel-bg:      #1e293b; /* Paneles interiores / cajas */

  /* 2. Bordes y Divisores */
  --border:              #334155; /* Bordes estructurales (Slate 700) */
  --divider:             #334155; /* Divisores horizontales */

  /* 3. Tipografía */
  --text-primary:        #f8fafc; /* Texto principal (Slate 50) */
  --text-muted:          #94a3b8; /* Texto secundario / etiquetas (Slate 400) */
  --nav-inactive-fg:     #93c5fd; /* Botones o enlaces inactivos */

  /* 4. Estados Semánticos y Tripletas RGB */
  --accent:              #818cf8; /* Acento primario (Indigo 400) */
  --accent-rgb:          129, 140, 248;
  --success:             #34d399; /* Éxito / Pensión favorable */
  --success-rgb:         52, 211, 153;
  --warning:             #fbbf24; /* Advertencia / Alerta */
  --warning-rgb:         251, 191, 36;
  --danger:              #f87171; /* Penalización / Error */
  --danger-rgb:          248, 113, 113;

  /* 5. Controles e Inputs */
  --input-bg:            #1e293b;
  --input-border:        #334155;
  --hover:               #334155;
  --selection:           #3730a3;
  --disabled:            #475569;
  --disabled-button:     #1e293b;

  /* 6. Cabecera Ejecutiva (ModernHeaderCard) */
  --header-banner-top:      #1c5239;
  --header-banner-bottom:   #0f3826;
  --header-banner-divider:  #062316;
  --header-body-bg:         #0f172a;
  --header-border:          #1e293b;
  --header-title:           #ffffff;
  --header-subtitle:        rgba(167, 243, 208, 0.86);

  /* 7. Insignias / Badges (utilizando variables con opacidad) */
  --badge-success-bg:       rgba(var(--success-rgb), 0.18);
  --badge-success-text:     #86efac;
  --badge-success-border:   var(--success);

  --badge-danger-bg:        rgba(var(--danger-rgb), 0.18);
  --badge-danger-text:      #fca5a5;
  --badge-danger-border:    var(--danger);

  --badge-info-bg:          rgba(var(--accent-rgb), 0.18);
  --badge-info-text:        #93c5fd;
  --badge-info-border:      var(--accent);

  --badge-warning-bg:       rgba(var(--warning-rgb), 0.18);
  --badge-warning-text:     #fde68a;
  --badge-warning-border:   var(--warning);

  --badge-purple-bg:        rgba(168, 85, 247, 0.18);
  --badge-purple-text:      #d8b4fe;
  --badge-purple-border:    #a855f7;

  --badge-neutral-bg:       #1e293b;
  --badge-neutral-text:     #cbd5e1;
  --badge-neutral-border:   #334155;

  /* 8. Botones de Acción (ModernButton) */
  --btn-primary-bg:              #4f46e5;
  --btn-primary-fg:              #ffffff;
  --btn-primary-border:          #6366f1;
  --btn-primary-hover-bg:        #4338ca;
  --btn-primary-hover-fg:        #ffffff;
  --btn-primary-hover-border:    #818cf8;
  --btn-primary-press-bg:        #3730a3;
  --btn-primary-press-fg:        #ffffff;
  --btn-primary-press-border:    #6366f1;
  --btn-primary-disabled-bg:     #1e293b;
  --btn-primary-disabled-fg:     #94a3b8;
  --btn-primary-disabled-border: #334155;

  --btn-default-bg:              #1e293b;
  --btn-default-fg:              #f1f5f9;
  --btn-default-border:          #334155;
  --btn-default-hover-bg:        #334155;
  --btn-default-hover-fg:        #ffffff;
  --btn-default-hover-border:    #475569;
  --btn-default-press-bg:        #0f172a;
  --btn-default-press-fg:        #ffffff;
  --btn-default-press-border:    #475569;
  --btn-default-disabled-bg:     #0f172a;
  --btn-default-disabled-fg:     #475569;
  --btn-default-disabled-border: #1e293b;

  --btn-danger-bg:               #1e293b;
  --btn-danger-fg:               #f87171;
  --btn-danger-border:           #f87171;
  --btn-danger-hover-bg:         #334155;
  --btn-danger-hover-fg:         #ffffff;
  --btn-danger-hover-border:     #fca5a5;
  --btn-danger-press-bg:         #0f172a;
  --btn-danger-press-fg:         #ffffff;
  --btn-danger-press-border:     #f87171;
  --btn-danger-disabled-bg:      #0f172a;
  --btn-danger-disabled-fg:      #475569;
  --btn-danger-disabled-border:  #1e293b;

  /* 9. Tablas y Rejillas de Datos (ModernBrowse) */
  --table-header-bg:        #1e293b;
  --table-header-fg:        #94a3b8;
  --table-header-hover-fg:  #f8fafc;
  --table-sort-indicator:   #38bdf8;
  --table-row-even:         #0f172a;
  --table-row-odd:          #1e293b;
  --table-row-hover:        #334155;
  --table-row-selected:     #1e3a8a;
  --table-border:           #334155;

  /* 10. Gráficos y Tooltips (ModernLineChart) */
  --chart-tooltip-bg:       #0f172a;
  --chart-tooltip-border:   #38bdf8;
  --chart-tooltip-title:    #ffffff;
  --chart-tooltip-subtitle: #94a3b8;
  --chart-tooltip-divider:  #334155;
  --chart-tooltip-text:     #f8fafc;
  --chart-tooltip-muted:    #94a3b8;
  --chart-guideline:        #0284c7;
  --chart-badge-bg:         #0f172a;
  --chart-badge-text:       #ffffff;

  /* 11. Identidad y Avatares */
  --avatar-bg:              #1e293b;
  --avatar-fg:              #60a5fa;
}
`, displayName, id)
	}

	return fmt.Sprintf(`/* ==========================================================================
   GoJubila PRO - Tema: %[1]s
   Generado desde Plantilla Clara
   ========================================================================== */

:root {
  /* Metadatos */
  --name: %[2]s;
  --display-name: "%[1]s";
  --is-dark: false;

  /* 1. Superficies y Fondos Base */
  --app-background:      #f8fafc; /* Fondo principal de la ventana (Slate 50) */
  --navbar-bg:           #ffffff; /* Barra superior de navegación */
  --card-bg:             #ffffff; /* Fondo de tarjetas y paneles */
  --card-border:         #e2e8f0; /* Borde perimetral de tarjetas (Slate 200) */
  --sidebar-bg:          #ffffff; /* Menú lateral */
  --sidebar-shaded-bg:   #f1f5f9; /* Fondo secundario de la barra lateral */
  --filter-drawer-bg:    #f8fafc; /* Cajón de filtros */
  --inset-panel-bg:      #f1f5f9; /* Paneles interiores / cajas */

  /* 2. Bordes y Divisores */
  --border:              #e2e8f0; /* Bordes estructurales (Slate 200) */
  --divider:             #e2e8f0; /* Divisores horizontales */

  /* 3. Tipografía */
  --text-primary:        #0f172a; /* Texto principal de alto contraste (Slate 900) */
  --text-muted:          #64748b; /* Texto secundario / etiquetas (Slate 500) */
  --nav-inactive-fg:     #64748b; /* Botones o enlaces inactivos */

  /* 4. Estados Semánticos y Tripletas RGB */
  --accent:              #2563eb; /* Acento primario (Blue 600) */
  --accent-rgb:          37, 99, 235;
  --success:             #16a34a; /* Éxito / Pensión favorable (Green 600) */
  --success-rgb:         22, 163, 74;
  --warning:             #d97706; /* Advertencia / Alerta (Amber 600) */
  --warning-rgb:         217, 119, 6;
  --danger:              #dc2626; /* Penalización / Error (Red 600) */
  --danger-rgb:          220, 38, 38;

  /* 5. Controles e Inputs */
  --input-bg:            #ffffff;
  --input-border:        #cbd5e1;
  --hover:               #f1f5f9;
  --selection:           #dbeafe;
  --disabled:            #94a3b8;
  --disabled-button:     #f1f5f9;

  /* 6. Cabecera Ejecutiva (ModernHeaderCard) */
  --header-banner-top:      #15803d;
  --header-banner-bottom:   #166534;
  --header-banner-divider:  #14532d;
  --header-body-bg:         #ffffff;
  --header-border:          #cbd5e1;
  --header-title:           #ffffff;
  --header-subtitle:        rgba(255, 255, 255, 0.85);

  /* 7. Insignias / Badges (utilizando variables con opacidad) */
  --badge-success-bg:       rgba(var(--success-rgb), 0.12);
  --badge-success-text:     #15803d;
  --badge-success-border:   #86efac;

  --badge-danger-bg:        rgba(var(--danger-rgb), 0.12);
  --badge-danger-text:      #b91c1c;
  --badge-danger-border:    #fca5a5;

  --badge-info-bg:          rgba(var(--accent-rgb), 0.12);
  --badge-info-text:        #1d4ed8;
  --badge-info-border:      #bfdbfe;

  --badge-warning-bg:       rgba(var(--warning-rgb), 0.12);
  --badge-warning-text:     #b45309;
  --badge-warning-border:   #fde68a;

  --badge-purple-bg:        rgba(168, 85, 247, 0.12);
  --badge-purple-text:      #7e22ce;
  --badge-purple-border:    #e9d5ff;

  --badge-neutral-bg:       #f1f5f9;
  --badge-neutral-text:     #475569;
  --badge-neutral-border:   #cbd5e1;

  /* 8. Botones de Acción (ModernButton) */
  --btn-primary-bg:              #2563eb;
  --btn-primary-fg:              #ffffff;
  --btn-primary-border:          #1d4ed8;
  --btn-primary-hover-bg:        #1d4ed8;
  --btn-primary-hover-fg:        #ffffff;
  --btn-primary-hover-border:    #1e40af;
  --btn-primary-press-bg:        #1e40af;
  --btn-primary-press-fg:        #ffffff;
  --btn-primary-press-border:    #172554;
  --btn-primary-disabled-bg:     #f1f5f9;
  --btn-primary-disabled-fg:     #94a3b8;
  --btn-primary-disabled-border: #e2e8f0;

  --btn-default-bg:              #ffffff;
  --btn-default-fg:              #0f172a;
  --btn-default-border:          #cbd5e1;
  --btn-default-hover-bg:        #f1f5f9;
  --btn-default-hover-fg:        #0f172a;
  --btn-default-hover-border:    #94a3b8;
  --btn-default-press-bg:        #e2e8f0;
  --btn-default-press-fg:        #0f172a;
  --btn-default-press-border:    #64748b;
  --btn-default-disabled-bg:     #f8fafc;
  --btn-default-disabled-fg:     #94a3b8;
  --btn-default-disabled-border: #e2e8f0;

  --btn-danger-bg:               #ffffff;
  --btn-danger-fg:               #dc2626;
  --btn-danger-border:           #fca5a5;
  --btn-danger-hover-bg:         #fef2f2;
  --btn-danger-hover-fg:         #b91c1c;
  --btn-danger-hover-border:     #f87171;
  --btn-danger-press-bg:         #fee2e2;
  --btn-danger-press-fg:         #991b1b;
  --btn-danger-press-border:     #ef4444;
  --btn-danger-disabled-bg:      #f8fafc;
  --btn-danger-disabled-fg:      #94a3b8;
  --btn-danger-disabled-border:  #e2e8f0;

  /* 9. Tablas y Rejillas de Datos (ModernBrowse) */
  --table-header-bg:        #f8fafc;
  --table-header-fg:        #64748b;
  --table-header-hover-fg:  #0f172a;
  --table-sort-indicator:   #2563eb;
  --table-row-even:         #ffffff;
  --table-row-odd:          #f8fafc;
  --table-row-hover:        #f1f5f9;
  --table-row-selected:     #eff6ff;
  --table-border:           #e2e8f0;

  /* 10. Gráficos y Tooltips (ModernLineChart) */
  --chart-tooltip-bg:       #ffffff;
  --chart-tooltip-border:   #2563eb;
  --chart-tooltip-title:    #0f172a;
  --chart-tooltip-subtitle: #64748b;
  --chart-tooltip-divider:  #e2e8f0;
  --chart-tooltip-text:     #0f172a;
  --chart-tooltip-muted:    #64748b;
  --chart-guideline:        #2563eb;
  --chart-badge-bg:         #ffffff;
  --chart-badge-text:       #0f172a;

  /* 11. Identidad y Avatares */
  --avatar-bg:              #f1f5f9;
  --avatar-fg:              #2563eb;
}
`, displayName, id)
}

// CreateThemeFromTemplate crea un nuevo archivo .css con la plantilla seleccionada dentro del directorio de temas.
// Devuelve la ruta absoluta del archivo creado o un error descriptivo si ya existe o no es válido.
func CreateThemeFromTemplate(dir, rawID, displayName string, isDark bool) (string, error) {
	if dir == "" {
		return "", fmt.Errorf("directorio de temas no especificado")
	}
	_ = os.MkdirAll(dir, 0755)

	cleanID := strings.TrimSpace(rawID)
	cleanID = strings.TrimSuffix(cleanID, ".css")
	cleanID = strings.ToLower(cleanID)
	cleanID = strings.ReplaceAll(cleanID, " ", "-")
	cleanID = strings.ReplaceAll(cleanID, "_", "-")
	reClean := regexp.MustCompile(`[^a-z0-9\-]`)
	cleanID = reClean.ReplaceAllString(cleanID, "")

	if cleanID == "" {
		return "", fmt.Errorf("debes indicar un identificador válido para el tema")
	}

	if cleanID == ThemeDark || cleanID == ThemeLight {
		cleanID = cleanID + "-css"
	}

	if strings.TrimSpace(displayName) == "" {
		displayName = strings.Title(strings.ReplaceAll(cleanID, "-", " "))
	}

	fileName := cleanID + ".css"
	targetPath := filepath.Join(dir, fileName)

	if _, err := os.Stat(targetPath); err == nil {
		return "", fmt.Errorf("ya existe un tema con el archivo '%s'", fileName)
	}

	content := GenerateThemeCSSTemplate(cleanID, displayName, isDark)
	err := os.WriteFile(targetPath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("error al escribir plantilla de tema: %w", err)
	}

	return targetPath, nil
}


