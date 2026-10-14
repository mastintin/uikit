package uikit

import (
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

func TestParseThemeCSS_Basic(t *testing.T) {
	css := `
/* Tema Nord de Prueba */
name: Nord
display-name: Nord Frost
is-dark: true

--app-background: #2e3440;
--card-bg:        #3b4252;
--card-border:    #4c566a;

--accent:         #88c0d0;
--success:        #a3be8c;
--warning:        #ebcb8b;
--danger:         #bf616a;

--text-primary:   #eceff4;
--text-muted:     #d8dee9;

--btn-primary-bg: #88c0d0;
--btn-primary-fg: #2e3440;
`

	pal, err := ParseThemeCSS(css, "nord.css")
	if err != nil {
		t.Fatalf("ParseThemeCSS falló: %v", err)
	}

	if pal.Name != "nord" {
		t.Errorf("Esperado Name 'nord', obtenido '%s'", pal.Name)
	}
	if pal.DisplayName != "Nord Frost" {
		t.Errorf("Esperado DisplayName 'Nord Frost', obtenido '%s'", pal.DisplayName)
	}
	if !pal.IsDark {
		t.Errorf("Esperado IsDark true")
	}

	// Verificar colores parseados
	if ColorToHex(pal.AppBackground) != "#2e3440" {
		t.Errorf("Esperado AppBackground #2e3440, obtenido %s", ColorToHex(pal.AppBackground))
	}
	if ColorToHex(pal.CardBg) != "#3b4252" {
		t.Errorf("Esperado CardBg #3b4252, obtenido %s", ColorToHex(pal.CardBg))
	}
	if ColorToHex(pal.Accent) != "#88c0d0" {
		t.Errorf("Esperado Accent #88c0d0, obtenido %s", ColorToHex(pal.Accent))
	}
	if ColorToHex(pal.Buttons.Primary.Bg) != "#88c0d0" {
		t.Errorf("Esperado Buttons.Primary.Bg #88c0d0, obtenido %s", ColorToHex(pal.Buttons.Primary.Bg))
	}
}

func TestScanThemesDir_And_Register(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "jubila_themes_test_*")
	if err != nil {
		t.Fatalf("Error creando tmpDir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	defer UnregisterTheme("dracula")

	file1 := filepath.Join(tmpDir, "dracula.css")
	css1 := `
name: Dracula
display-name: Dracula Dark
is-dark: true
app-bg: #282a36
card-bg: #44475a
accent: #bd93f9
`
	if err := os.WriteFile(file1, []byte(css1), 0644); err != nil {
		t.Fatalf("Error escribiendo archivo: %v", err)
	}

	count := LoadAndRegisterExternalThemes(tmpDir)
	if count != 1 {
		t.Fatalf("Esperado 1 tema registrado, obtenido %d", count)
	}

	pal := GetThemePalette("dracula")
	if pal == nil {
		t.Fatalf("Esperado tema 'dracula' registrado en catálogo global")
	}
	if pal.DisplayName != "Dracula Dark" {
		t.Errorf("Esperado 'Dracula Dark', obtenido '%s'", pal.DisplayName)
	}
	if ColorToHex(pal.Accent) != "#bd93f9" {
		t.Errorf("Esperado Accent #bd93f9, obtenido %s", ColorToHex(pal.Accent))
	}
}

func TestParseThemeCSS_DarkMasterTemplate(t *testing.T) {
	darkPath := filepath.Join("..", "themes", "dark.css")
	pal, err := LoadThemeCSSFile(darkPath)
	if err != nil {
		t.Fatalf("Error al cargar themes/dark.css: %v", err)
	}

	if pal.Name != "dark-css" {
		t.Errorf("Esperado Name 'dark-css', obtenido '%s'", pal.Name)
	}
	if pal.DisplayName != "Oscuro Moderno (CSS)" {
		t.Errorf("Esperado DisplayName 'Oscuro Moderno (CSS)', obtenido '%s'", pal.DisplayName)
	}
	if !pal.IsDark {
		t.Errorf("Esperado IsDark true")
	}

	// Superficies
	if ColorToHex(pal.AppBackground) != "#020617" {
		t.Errorf("AppBackground erróneo: %s", ColorToHex(pal.AppBackground))
	}
	if ColorToHex(pal.CardBg) != "#0f172a" {
		t.Errorf("CardBg erróneo: %s", ColorToHex(pal.CardBg))
	}

	// Badges
	if ColorToHex(pal.Badges.Success.Text) != "#86efac" {
		t.Errorf("Badge Success Text erróneo: %s", ColorToHex(pal.Badges.Success.Text))
	}
	if ColorToHex(pal.Badges.Danger.Border) != "#ef4444" {
		t.Errorf("Badge Danger Border erróneo: %s", ColorToHex(pal.Badges.Danger.Border))
	}
	if ColorToHex(pal.Badges.Warning.Text) != "#fde68a" {
		t.Errorf("Badge Warning Text erróneo: %s", ColorToHex(pal.Badges.Warning.Text))
	}
	if ColorToHex(pal.Badges.Purple.Border) != "#a855f7" {
		t.Errorf("Badge Purple Border erróneo: %s", ColorToHex(pal.Badges.Purple.Border))
	}

	// Botones
	if ColorToHex(pal.Buttons.Primary.Bg) != "#4f46e5" {
		t.Errorf("Btn Primary Bg erróneo: %s", ColorToHex(pal.Buttons.Primary.Bg))
	}
	if ColorToHex(pal.Buttons.Primary.HoverBg) != "#4338ca" {
		t.Errorf("Btn Primary HoverBg erróneo: %s", ColorToHex(pal.Buttons.Primary.HoverBg))
	}
	if ColorToHex(pal.Buttons.Primary.PressBg) != "#3730a3" {
		t.Errorf("Btn Primary PressBg erróneo: %s", ColorToHex(pal.Buttons.Primary.PressBg))
	}
	if ColorToHex(pal.Buttons.Danger.Fg) != "#f87171" {
		t.Errorf("Btn Danger Fg erróneo: %s", ColorToHex(pal.Buttons.Danger.Fg))
	}

	// Tablas
	if ColorToHex(pal.Table.HeaderBg) != "#1e293b" {
		t.Errorf("Table HeaderBg erróneo: %s", ColorToHex(pal.Table.HeaderBg))
	}
	if ColorToHex(pal.Table.SortIndicator) != "#38bdf8" {
		t.Errorf("Table SortIndicator erróneo: %s", ColorToHex(pal.Table.SortIndicator))
	}
	if ColorToHex(pal.Table.RowSelectedBg) != "#1e3a8a" {
		t.Errorf("Table RowSelectedBg erróneo: %s", ColorToHex(pal.Table.RowSelectedBg))
	}

	// Gráficos y Avatares
	if ColorToHex(pal.Chart.TooltipBorder) != "#38bdf8" {
		t.Errorf("Chart TooltipBorder erróneo: %s", ColorToHex(pal.Chart.TooltipBorder))
	}
	if ColorToHex(pal.Avatar.Fg) != "#60a5fa" {
		t.Errorf("Avatar Fg erróneo: %s", ColorToHex(pal.Avatar.Fg))
	}
}

func TestParseThemeCSS_RGBA(t *testing.T) {
	css := `
name: CustomRGBA
is-dark: true
--app-background: rgba(10, 20, 30, 1.0)
--card-bg: rgba(40, 50, 60, 0.5)
--badge-success-bg: rgba(6, 78, 59, 200)
`
	pal, err := ParseThemeCSS(css, "rgba.css")
	if err != nil {
		t.Fatalf("ParseThemeCSS falló: %v", err)
	}

	if ColorToHex(pal.AppBackground) != "#0a141e" {
		t.Errorf("Esperado #0a141e, obtenido %s", ColorToHex(pal.AppBackground))
	}

	nrgbaBg := ColorToNRGBA(pal.Badges.Success.Bg)
	if nrgbaBg.R != 6 || nrgbaBg.G != 78 || nrgbaBg.B != 59 || nrgbaBg.A != 200 {
		t.Errorf("Esperado RGBA(6, 78, 59, 200), obtenido RGBA(%d, %d, %d, %d)",
			nrgbaBg.R, nrgbaBg.G, nrgbaBg.B, nrgbaBg.A)
	}
}

func TestParseThemeCSS_SampleThemes(t *testing.T) {
	// 1. Probar themes/emerald.css
	emeraldPath := filepath.Join("..", "themes", "emerald.css")
	palEmerald, err := LoadThemeCSSFile(emeraldPath)
	if err != nil {
		t.Fatalf("Error al cargar themes/emerald.css: %v", err)
	}
	if palEmerald.Name != "emerald" {
		t.Errorf("Esperado Name 'emerald', obtenido '%s'", palEmerald.Name)
	}
	if palEmerald.DisplayName != "Esmeralda Pro" {
		t.Errorf("Esperado DisplayName 'Esmeralda Pro', obtenido '%s'", palEmerald.DisplayName)
	}
	if palEmerald.IsDark {
		t.Errorf("Esperado IsDark false para emerald")
	}
	if ColorToHex(palEmerald.AppBackground) != "#f0fdf4" {
		t.Errorf("Emerald AppBackground erróneo: %s", ColorToHex(palEmerald.AppBackground))
	}
	if ColorToHex(palEmerald.Buttons.Primary.Bg) != "#059669" {
		t.Errorf("Emerald Btn Primary Bg erróneo: %s", ColorToHex(palEmerald.Buttons.Primary.Bg))
	}

	// 2. Probar themes/nord.css
	nordPath := filepath.Join("..", "themes", "nord.css")
	palNord, err := LoadThemeCSSFile(nordPath)
	if err != nil {
		t.Fatalf("Error al cargar themes/nord.css: %v", err)
	}
	if palNord.Name != "nord" {
		t.Errorf("Esperado Name 'nord', obtenido '%s'", palNord.Name)
	}
	if palNord.DisplayName != "Nord Ártico" {
		t.Errorf("Esperado DisplayName 'Nord Ártico', obtenido '%s'", palNord.DisplayName)
	}
	if !palNord.IsDark {
		t.Errorf("Esperado IsDark true para nord")
	}
	if ColorToHex(palNord.AppBackground) != "#2e3440" {
		t.Errorf("Nord AppBackground erróneo: %s", ColorToHex(palNord.AppBackground))
	}
	if ColorToHex(palNord.Buttons.Primary.Bg) != "#88c0d0" {
		t.Errorf("Nord Btn Primary Bg erróneo: %s", ColorToHex(palNord.Buttons.Primary.Bg))
	}

	// 3. Probar themes/nice.css
	nicePath := filepath.Join("..", "themes", "nice.css")
	palNice, err := LoadThemeCSSFile(nicePath)
	if err != nil {
		t.Fatalf("Error al cargar themes/nice.css: %v", err)
	}
	if palNice.Name != "nice" {
		t.Errorf("Esperado Name 'nice', obtenido '%s'", palNice.Name)
	}
	if palNice.DisplayName != "Nice (xlsJubila)" {
		t.Errorf("Esperado DisplayName 'Nice (xlsJubila)', obtenido '%s'", palNice.DisplayName)
	}
	if palNice.IsDark {
		t.Errorf("Esperado IsDark false para nice")
	}
	if ColorToHex(palNice.AppBackground) != "#f1f5f9" {
		t.Errorf("Nice AppBackground erróneo: %s", ColorToHex(palNice.AppBackground))
	}
	if ColorToHex(palNice.Buttons.Primary.Bg) != "#107c41" {
		t.Errorf("Nice Btn Primary Bg erróneo: %s", ColorToHex(palNice.Buttons.Primary.Bg))
	}
	if ColorToHex(palNice.NavbarBg) != "#0f172a" {
		t.Errorf("Nice NavbarBg erróneo: %s", ColorToHex(palNice.NavbarBg))
	}
}

func TestParseThemeCSS_VariableReferencesAndRGB(t *testing.T) {
	css := `
:root {
  --name: test-variables;
  --display-name: "Variables Dinámicas";
  --is-dark: true;

  --danger-rgb: 239, 68, 68; /* #ef4444 en RGB */
  
  /* Ahora puedes jugar con las opacidades manteniendo la consistencia */
  --badge-danger-bg:     rgba(var(--danger-rgb), 0.15); 
  --badge-danger-border: rgba(var(--danger-rgb), 0.4);
  --badge-danger-text:   #fca5a5;
}
`

	pal, err := ParseThemeCSS(css, "vars.css")
	if err != nil {
		t.Fatalf("ParseThemeCSS falló con variables dinámicas: %v", err)
	}

	// 1. Verificar Badge Danger Background: rgba(239, 68, 68, 0.15)
	bgNRGBA, ok := pal.Badges.Danger.Bg.(color.NRGBA)
	if !ok {
		t.Fatalf("Esperado pal.Badges.Danger.Bg de tipo color.NRGBA, obtenido %T", pal.Badges.Danger.Bg)
	}
	if bgNRGBA.R != 239 || bgNRGBA.G != 68 || bgNRGBA.B != 68 {
		t.Errorf("Esperado RGB (239, 68, 68), obtenido (%d, %d, %d)", bgNRGBA.R, bgNRGBA.G, bgNRGBA.B)
	}
	// 0.15 * 255 = 38.25 -> 38
	if bgNRGBA.A < 37 || bgNRGBA.A > 39 {
		t.Errorf("Esperado Alpha aprox 38 (0.15 * 255), obtenido %d", bgNRGBA.A)
	}

	// 2. Verificar Badge Danger Border: rgba(239, 68, 68, 0.4)
	borderNRGBA, ok := pal.Badges.Danger.Border.(color.NRGBA)
	if !ok {
		t.Fatalf("Esperado pal.Badges.Danger.Border de tipo color.NRGBA, obtenido %T", pal.Badges.Danger.Border)
	}
	if borderNRGBA.R != 239 || borderNRGBA.G != 68 || borderNRGBA.B != 68 {
		t.Errorf("Esperado RGB (239, 68, 68), obtenido (%d, %d, %d)", borderNRGBA.R, borderNRGBA.G, borderNRGBA.B)
	}
	// 0.4 * 255 = 102
	if borderNRGBA.A < 101 || borderNRGBA.A > 103 {
		t.Errorf("Esperado Alpha aprox 102 (0.4 * 255), obtenido %d", borderNRGBA.A)
	}

	// 3. Verificar Badge Danger Text: #fca5a5
	if ColorToHex(pal.Badges.Danger.Text) != "#fca5a5" {
		t.Errorf("Esperado Text #fca5a5, obtenido %s", ColorToHex(pal.Badges.Danger.Text))
	}
}

func TestCreateThemeFromTemplate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "jubila_template_test_*")
	if err != nil {
		t.Fatalf("Error creando tmpDir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 1. Crear tema oscuro desde plantilla
	darkPath, err := CreateThemeFromTemplate(tmpDir, "mi-tema-oscuro", "Mi Tema Oscuro", true)
	if err != nil {
		t.Fatalf("CreateThemeFromTemplate (dark) falló: %v", err)
	}
	palDark, err := LoadThemeCSSFile(darkPath)
	if err != nil {
		t.Fatalf("Error cargando plantilla oscura generada: %v", err)
	}
	if !palDark.IsDark {
		t.Errorf("Esperado tema oscuro IsDark=true")
	}
	if palDark.Name != "mi-tema-oscuro" {
		t.Errorf("Esperado Name 'mi-tema-oscuro', obtenido '%s'", palDark.Name)
	}

	// 2. Crear tema claro desde plantilla
	lightPath, err := CreateThemeFromTemplate(tmpDir, "mi-tema-claro", "Mi Tema Claro", false)
	if err != nil {
		t.Fatalf("CreateThemeFromTemplate (light) falló: %v", err)
	}
	palLight, err := LoadThemeCSSFile(lightPath)
	if err != nil {
		t.Fatalf("Error cargando plantilla clara generada: %v", err)
	}
	if palLight.IsDark {
		t.Errorf("Esperado tema claro IsDark=false")
	}

	// 3. Intentar crear duplicado debe fallar con mensaje claro
	_, errDupl := CreateThemeFromTemplate(tmpDir, "mi-tema-oscuro", "Duplicado", true)
	if errDupl == nil {
		t.Errorf("Esperado error al crear tema duplicado con mismo ID")
	}
}

func TestParseThemeCSS_SyntaxWarningsAndErrors(t *testing.T) {
	// 1. Archivo con avisos de sintaxis (línea sin colon, color inválido, variable no resuelta)
	cssWithWarnings := `
:root {
  --name: tema-con-errores;
  --display-name: "Tema con Errores";
  linea_sin_dos_puntos_invalida
  --card-bg: color_inventado_invalido;
  --accent: #12345;
  --btn-primary-bg: var(--variable_que_no_existe);
}
`

	pal, warnings, err := ParseThemeCSSWithReport(cssWithWarnings, "errores.css")
	if err != nil {
		t.Fatalf("ParseThemeCSSWithReport no debe abortar en advertencias recuperables: %v", err)
	}
	if pal == nil {
		t.Fatalf("Esperada paleta con fallbacks seguros generada")
	}

	if len(warnings) < 3 {
		t.Errorf("Esperadas al menos 3 advertencias de sintaxis, obtenidas %d: %v", len(warnings), warnings)
	}

	// 2. Archivo completamente vacío o sin propiedades debe retornar error
	_, _, errEmpty := ParseThemeCSSWithReport("   \n // solo comentarios \n /* nada */ ", "vacio.css")
	if errEmpty == nil {
		t.Errorf("Esperado error al procesar CSS vacío")
	}
}





