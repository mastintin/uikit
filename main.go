package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"uikit-demo/uikit"
)

// findThemesDir resuelve de forma ultra-robusta la ubicación de la carpeta "themes"
// independientemente de si se ejecuta desde terminal, Finder, IDE o subdirectorios.
func findThemesDir() string {
	// 1. Comprobar ruta directa en el directorio de trabajo actual
	if stat, err := os.Stat("themes"); err == nil && stat.IsDir() {
		if abs, err := filepath.Abs("themes"); err == nil {
			return abs
		}
		return "themes"
	}

	// 2. Comprobar ruta relativa al binario ejecutable
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		cand := filepath.Join(exeDir, "themes")
		if stat, err := os.Stat(cand); err == nil && stat.IsDir() {
			return cand
		}
	}

	// 3. Buscar subiendo por el árbol de directorios padres
	if cwd, err := os.Getwd(); err == nil {
		curr := cwd
		for i := 0; i < 5; i++ {
			cand := filepath.Join(curr, "themes")
			if stat, err := os.Stat(cand); err == nil && stat.IsDir() {
				return cand
			}
			parent := filepath.Dir(curr)
			if parent == curr {
				break
			}
			curr = parent
		}
	}

	// 4. Ruta absoluta de contingencia del proyecto
	hardcoded := "/Users/manuel/antiProyectos/uikit/themes"
	if stat, err := os.Stat(hardcoded); err == nil && stat.IsDir() {
		return hardcoded
	}

	return "themes"
}

func main() {
	myApp := app.NewWithID("com.uikit.showcase")
	myWindow := myApp.NewWindow("UIKit Component Gallery & CSS Themes - Fyne 2.8")
	myWindow.Resize(fyne.NewSize(1280, 880))

	// -------------------------------------------------------------------------
	// 0. Carga e Inicialización de Temas CSS Externos
	// -------------------------------------------------------------------------
	themesDir := findThemesDir()
	initReport := uikit.LoadAndRegisterExternalThemesDetailed(themesDir)
	uikit.SetAppTheme(uikit.ThemeLight)

	// Fondo de aplicación reactivo al tema
	appBg := canvas.NewRectangle(uikit.ThemeColor(uikit.ColorNameAppBackground))

	// Helper para obtener la lista de temas y sus nombres descriptivos
	getThemesList := func() ([]string, map[string]string, string) {
		var displayList []string
		idMap := make(map[string]string)
		currName := uikit.GetCurrentThemeName()
		currentDisplay := "Claro (Light)"

		for _, id := range uikit.AvailableThemes() {
			name := id
			if p := uikit.GetThemePalette(id); p != nil && p.DisplayName != "" {
				name = p.DisplayName
			} else if id == uikit.ThemeLight {
				name = "Claro (Light)"
			} else if id == uikit.ThemeDark {
				name = "Oscuro (Dark)"
			}
			displayList = append(displayList, name)
			idMap[name] = id
			if id == currName {
				currentDisplay = name
			}
		}
		return displayList, idMap, currentDisplay
	}

	themeDisplays, themeMap, currentDisplay := getThemesList()

	// -------------------------------------------------------------------------
	// 1. Cabecera Principal (HeaderCard con título degradado)
	// -------------------------------------------------------------------------
	header := uikit.NewHeaderCardWithTitle(
		"sliders",
		"CATÁLOGO DE COMPONENTES UIKIT",
		"Biblioteca visual moderna para aplicaciones Go + Fyne (Controles contemporáneos, responsivos y fluidos)",
		nil,
	)

	// -------------------------------------------------------------------------
	// 2. Panel Superior Destacado: Selector y Gestor de Temas CSS
	// -------------------------------------------------------------------------
	lblThemesInfo := uikit.NewCaptionLabel(
		fmt.Sprintf("Carpeta: %s • %d temas CSS cargados", themesDir, initReport.Count),
		uikit.ThemeColor(uikit.ColorNameTextMuted),
	)

	var selectTopTheme *widget.Select
	var selectSidebarTheme *widget.Select
	var chipThemeTop *uikit.ModernChipFilter
	var lblBrand *canvas.Text
	var sidebar *uikit.ModernSidebar

	applyTheme := func(themeID string, showFeedback bool) {
		th := uikit.SetAppTheme(themeID)

		// 1. Fondo de ventana sincronizado
		appBg.FillColor = uikit.ThemeColor(uikit.ColorNameAppBackground)
		appBg.Refresh()

		// 2. Elementos gráficos personalizados del sidebar
		if lblBrand != nil {
			lblBrand.Color = uikit.ThemeColor(uikit.ColorNameTextPrimary)
			lblBrand.Refresh()
		}

		// 3. Refrescar recursivamente todo el árbol de controles
		uikit.RefreshTree(myWindow.Content())
		myWindow.Content().Refresh()

		// 4. Sincronizar selectores
		_, _, disp := getThemesList()
		if selectTopTheme != nil && selectTopTheme.Selected != disp {
			selectTopTheme.Selected = disp
			selectTopTheme.Refresh()
		}
		if selectSidebarTheme != nil && selectSidebarTheme.Selected != disp {
			selectSidebarTheme.Selected = disp
			selectSidebarTheme.Refresh()
		}
		if chipThemeTop != nil && chipThemeTop.Selected() != disp {
			chipThemeTop.SetSelectedSilent(disp)
		}

		if showFeedback {
			displayName := th.Name()
			if pal := uikit.GetThemePalette(themeID); pal != nil && pal.DisplayName != "" {
				displayName = pal.DisplayName
			}
			uikit.ShowPopUp(myWindow.Canvas(), uikit.ToastSuccess, fmt.Sprintf("Tema aplicado: %s", displayName))
		}
	}

	selectTopTheme = widget.NewSelect(themeDisplays, func(selected string) {
		if id, ok := themeMap[selected]; ok {
			applyTheme(id, true)
		}
	})
	selectTopTheme.Selected = currentDisplay

	btnReloadTop := uikit.NewButton("Primary", "Recargar Temas CSS", theme.ViewRefreshIcon(), false, func() {
		rep := uikit.LoadAndRegisterExternalThemesDetailed(themesDir)
		td, tm, cur := getThemesList()
		themeDisplays = td
		themeMap = tm

		selectTopTheme.Options = themeDisplays
		selectTopTheme.Selected = cur
		selectTopTheme.Refresh()

		if selectSidebarTheme != nil {
			selectSidebarTheme.Options = themeDisplays
			selectSidebarTheme.Selected = cur
			selectSidebarTheme.Refresh()
		}

		if chipThemeTop != nil {
			chipThemeTop.SetOptions(themeDisplays)
			chipThemeTop.SetSelectedSilent(cur)
		}

		lblThemesInfo.SetText(fmt.Sprintf("Carpeta: %s • %d temas CSS cargados", themesDir, rep.Count))

		applyTheme(uikit.GetCurrentThemeName(), false)

		if len(rep.Errors) > 0 {
			uikit.ShowPopUp(myWindow.Canvas(), uikit.ToastDanger, fmt.Sprintf("Error CSS: %s", rep.Errors[0]))
		} else {
			uikit.ShowPopUp(myWindow.Canvas(), uikit.ToastSuccess, fmt.Sprintf("%d temas CSS recargados en caliente", rep.Count))
		}
	})

	btnOpenFolder := uikit.NewButton("Default", "Abrir Carpeta Temas", theme.FolderOpenIcon(), false, func() {
		_ = uikit.OpenFolderInOS(themesDir)
	})

	switchGlass := uikit.NewSwitch("Efecto Cristal (Glass)", false, func(active bool) {
		uikit.SetGlobalGlassButtons(active)
		uikit.RefreshTree(myWindow.Content())
		myWindow.Content().Refresh()
	})

	themeControlsRow := container.NewHBox(
		container.NewCenter(widget.NewLabelWithStyle("Tema Activo:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		selectTopTheme,
		uiSpacer(8, 1),
		btnReloadTop,
		btnOpenFolder,
		uiSpacer(8, 1),
		switchGlass,
	)

	cardThemeManager := uikit.NewCard(container.NewVBox(
		container.NewHBox(
			container.NewCenter(uikit.NewColoredIcon("sliders", 20, "#107c41")),
			uikit.NewSectionTitle("Gestor y Selector de Temas Visuales (CSS en Caliente)"),
			layout.NewSpacer(),
			lblThemesInfo,
		),
		widget.NewSeparator(),
		themeControlsRow,
	))

	// -------------------------------------------------------------------------
	// 3. Tarjetas KPI / Métricas (KPICard)
	// -------------------------------------------------------------------------
	kpi1 := uikit.NewKPICard(uikit.KPICardEmerald, "Usuarios Activos", "14.280", "+12.5% este mes")
	kpi2 := uikit.NewKPICard(uikit.KPICardBlue, "Tiempo de Respuesta", "120 ms", "-18% latencia")
	kpi3 := uikit.NewKPICard(uikit.KPICardViolet, "Disponibilidad", "99.98 %", "Alta disponibilidad")
	kpi4 := uikit.NewKPICard(uikit.KPICardAmber, "Incidencias", "0 críticas", "Estable")

	kpiGrid := container.NewGridWithColumns(4, kpi1, kpi2, kpi3, kpi4)
	cardKPIs := uikit.NewCard(container.NewVBox(
		uikit.NewSectionTitle("Métricas y Tarjetas KPI (KPICard)"),
		kpiGrid,
	))

	// -------------------------------------------------------------------------
	// 4. Sección de Botones Modernos (ModernButton con Loading & Glass)
	// -------------------------------------------------------------------------
	var btnLoading *uikit.ModernButton
	btnLoading = uikit.NewPrimaryButton("Procesar Cálculo", theme.MediaPlayIcon(), func() {
		btnLoading.SetLoading(true)
		go func() {
			time.Sleep(2500 * time.Millisecond)
			btnLoading.SetLoading(false)
			uikit.ShowPopUp(myWindow.Canvas(), uikit.ToastSuccess, "Cálculo completado correctamente")
		}()
	})

	btnPrimary := uikit.NewPrimaryButton("Guardar Registro", theme.DocumentSaveIcon(), func() {
		uikit.ShowPopUp(myWindow.Canvas(), uikit.ToastSuccess, "Registro guardado correctamente")
	})

	btnDanger := uikit.NewDangerButton("Eliminar Elemento", theme.DeleteIcon(), func() {
		uikit.MessageYesNo(myWindow, "¿Estás seguro de eliminar el registro seleccionado?", func() {
			uikit.ShowPopUp(myWindow.Canvas(), uikit.ToastWarning, "Elemento eliminado")
		}, "Atención")
	})

	btnDefault := uikit.NewDefaultButton("Cancelar Operación", theme.CancelIcon(), func() {
		uikit.ShowPopUp(myWindow.Canvas(), uikit.ToastInfo, "Operación cancelada")
	})

	btnSmall := uikit.NewSmallPrimaryButton("Acción Compacta", theme.ContentAddIcon(), func() {})
	btnIcon := uikit.NewIconButton(theme.SettingsIcon(), func() {})

	buttonsRow := container.NewHBox(
		btnLoading,
		btnPrimary,
		btnDanger,
		btnDefault,
		btnSmall,
		btnIcon,
	)

	cardButtons := uikit.NewCard(container.NewVBox(
		uikit.NewSectionTitle("Botones & Estados de Carga (ModernButton)"),
		buttonsRow,
	))

	// -------------------------------------------------------------------------
	// 5. Sección de Deslizadores Modernos (ModernSlider)
	// -------------------------------------------------------------------------
	lblSliderOutput := uikit.NewModernLabel("Mueve los deslizadores para ver la reactividad en tiempo real...", uikit.StyleBody)

	sliderAge := uikit.NewAgeSlider(60, 70).
		WithTitle("Edad Deseada de Jubilación").
		WithValue(65).
		OnChanged(func(v float64) {
			lblSliderOutput.SetText(fmt.Sprintf("Edad seleccionada: %.0f años", v))
		})

	sliderPct := uikit.NewPercentageSlider(0, 100).
		WithTitle("Porcentaje de Jubilación Parcial / Relevo").
		WithValue(50).
		WithStep(5).
		OnChanged(func(v float64) {
			lblSliderOutput.SetText(fmt.Sprintf("Porcentaje parcial: %.0f %%", v))
		})

	sliderCurr := uikit.NewCurrencySlider(1000, 4500, 50).
		WithTitle("Base de Cotización Estimada Mensual").
		WithValue(2650).
		OnChanged(func(v float64) {
			lblSliderOutput.SetText(fmt.Sprintf("Base estimada: %.2f €", v))
		})

	sliderInline := uikit.NewSlider(0, 10).
		WithPrefix("Factor x").
		WithPrecision(1).
		WithValue(3.5)

	slidersBox := container.NewVBox(
		sliderAge,
		sliderPct,
		sliderCurr,
		container.NewHBox(widget.NewLabel("Slider Inline sin título:"), sliderInline),
		lblSliderOutput,
	)

	cardSliders := uikit.NewCard(container.NewVBox(
		uikit.NewSectionTitle("Deslizadores Numéricos Reactivos (ModernSlider)"),
		slidersBox,
	))

	// -------------------------------------------------------------------------
	// 6. Sección de Indicadores de Carga (ModernSpinner)
	// -------------------------------------------------------------------------
	spinSmall := uikit.NewPrimarySpinner(uikit.SpinnerSmall).Start()
	spinMed := uikit.NewSuccessSpinner(uikit.SpinnerMedium).Start()
	spinLarge := uikit.NewDangerSpinner(uikit.SpinnerLarge).Start()
	spinHero := uikit.NewAccentSpinner(uikit.SpinnerHero).Start()

	btnToggleSpinners := uikit.NewDefaultButton("Iniciar / Pausar Spinners", theme.MediaPlayIcon(), func() {
		spinSmall.Toggle()
		spinMed.Toggle()
		spinLarge.Toggle()
		spinHero.Toggle()
	})

	spinnersRow := container.NewHBox(
		spinSmall, widget.NewLabel("16px (Small)"),
		spinMed, widget.NewLabel("24px (Medium)"),
		spinLarge, widget.NewLabel("36px (Large)"),
		spinHero, widget.NewLabel("48px (Hero)"),
		btnToggleSpinners,
	)

	cardSpinners := uikit.NewCard(container.NewVBox(
		uikit.NewSectionTitle("Indicadores Circulares de Actividad (ModernSpinner)"),
		spinnersRow,
	))

	// -------------------------------------------------------------------------
	// 7. Badges, Pastillas e Interruptores (ModernBadge, Switch)
	// -------------------------------------------------------------------------
	badge1 := uikit.NewBadge(uikit.BadgeSuccess, "Aprobado")
	badge2 := uikit.NewBadge(uikit.BadgePrimary, "En Revisión")
	badge3 := uikit.NewBadge(uikit.BadgeWarning, "Pendiente")
	badge4 := uikit.NewBadge(uikit.BadgeDanger, "Rechazado")
	badge5 := uikit.NewBadge(uikit.BadgeNeutral, "Borrador")

	switchFeature := uikit.NewSwitch("Ahorro de Energía", true, func(b bool) {})
	switchNotify := uikit.NewSwitch("Notificaciones Push", false, func(b bool) {})

	badgesRow := container.NewHBox(
		badge1, badge2, badge3, badge4, badge5,
		widget.NewSeparator(),
		switchFeature,
		switchNotify,
	)

	// -------------------------------------------------------------------------
	// 8. Filtros Desplegables Tipo Pastilla (ModernChipFilter)
	// -------------------------------------------------------------------------
	chipYear := uikit.NewChipFilter("Año:", []string{"Todos los años", "2026", "2025", "2024", "2023"}, "Todos los años", func(sel string) {
		uikit.ShowPopUp(myWindow.Canvas(), uikit.ToastInfo, "Filtro activo: "+sel)
	})

	chipStatus := uikit.NewChipFilter("Estado:", []string{"Todos", "Activo", "Pendiente", "Archivado"}, "Todos", func(sel string) {
		uikit.ShowPopUp(myWindow.Canvas(), uikit.ToastInfo, "Estado filtrado: "+sel)
	})

	chipThemeTop = uikit.NewChipFilter("Tema CSS:", themeDisplays, currentDisplay, func(selected string) {
		if id, ok := themeMap[selected]; ok {
			applyTheme(id, true)
		}
	})

	chipFiltersRow := container.NewHBox(
		widget.NewLabel("Filtros Interactivos:"),
		chipYear,
		chipStatus,
		chipThemeTop,
	)

	// -------------------------------------------------------------------------
	// 9. Selector de Fechas (ModernDatePicker)
	// -------------------------------------------------------------------------
	datePicker := uikit.NewModernDatePicker(time.Now(), func(t time.Time) {
		uikit.ShowPopUp(myWindow.Canvas(), uikit.ToastInfo, "Fecha seleccionada: "+t.Format("02/01/2006"))
	})
	datePickerRow := container.NewHBox(
		widget.NewLabel("Fecha de Referencia:"),
		datePicker.Container,
	)

	cardBadgesAndFilters := uikit.NewCard(container.NewVBox(
		uikit.NewSectionTitle("Insignias, Filtros, Interruptores & Fechas"),
		container.NewVBox(badgesRow, chipFiltersRow, datePickerRow),
	))

	// -------------------------------------------------------------------------
	// 10. Contenedor Dinámico Principal (Centro)
	// -------------------------------------------------------------------------
	allContent := container.NewVBox(
		header,
		cardThemeManager,
		cardKPIs,
		cardButtons,
		cardSliders,
		cardSpinners,
		cardBadgesAndFilters,
	)

	scrollArea := container.NewVScroll(allContent)

	// Función para conmutar vista central desde el Sidebar
	setView := func(target fyne.CanvasObject) {
		scrollArea.Content = target
		scrollArea.ScrollToTop()
		scrollArea.Refresh()
	}

	// -------------------------------------------------------------------------
	// 11. Construcción de ModernSidebar (Barra Lateral Izquierda)
	// -------------------------------------------------------------------------
	brandIcon := uikit.NewColoredIcon("sliders", 22, "#107c41")
	lblBrand = canvas.NewText("UIKit Gallery", uikit.ThemeColor(uikit.ColorNameTextPrimary))
	lblBrand.TextStyle = fyne.TextStyle{Bold: true}
	lblBrand.TextSize = 15
	badgeVersion := uikit.NewBadge(uikit.BadgeSuccess, "v2.8")

	selectSidebarTheme = widget.NewSelect(themeDisplays, func(selected string) {
		if id, ok := themeMap[selected]; ok {
			applyTheme(id, true)
		}
	})
	selectSidebarTheme.Selected = currentDisplay

	sidebarHeader := container.NewPadded(container.NewVBox(
		container.NewHBox(
			brandIcon,
			lblBrand,
			layout.NewSpacer(),
			badgeVersion,
		),
		uikit.NewCaptionLabel("Catálogo Oficial de Componentes", uikit.ThemeColor(uikit.ColorNameTextMuted)),
		uiSpacer(1, 4),
		uikit.NewCaptionLabel("TEMA ACTIVO:", uikit.ThemeColor(uikit.ColorNameTextMuted)),
		selectSidebarTheme,
		widget.NewSeparator(),
	))

	// Botones de Navegación del Sidebar con el estilo ejecutivo de gofyneJubila
	var navBtns []*uikit.ModernSidebarNavButton
	addNav := func(title, iconName string, action func()) *uikit.ModernSidebarNavButton {
		var btn *uikit.ModernSidebarNavButton
		btn = uikit.NewSidebarNavButton(title, iconName, func() {
			for _, b := range navBtns {
				b.SetActive(false)
			}
			btn.SetActive(true)
			action()
		})
		navBtns = append(navBtns, btn)
		return btn
	}

	btnAll := addNav("Vista Completa", "home", func() {
		setView(allContent)
	})
	btnAll.SetActive(true) // Seleccionado por defecto con barra de acento

	btnThemesView := addNav("Temas & Estilos CSS", "sliders", func() {
		setView(container.NewVBox(header, cardThemeManager))
	})

	btnKPI := addNav("Métricas & KPIs", "trending-up", func() {
		setView(container.NewVBox(header, cardThemeManager, cardKPIs))
	})

	btnBtns := addNav("Botones & Carga", "play", func() {
		setView(container.NewVBox(header, cardThemeManager, cardButtons))
	})

	btnSliders := addNav("Deslizadores (Sliders)", "sliders", func() {
		setView(container.NewVBox(header, cardThemeManager, cardSliders))
	})

	btnSpin := addNav("Indicadores (Spinners)", "refresh", func() {
		setView(container.NewVBox(header, cardThemeManager, cardSpinners))
	})

	btnChips := addNav("Filtros & Insignias", "filter", func() {
		setView(container.NewVBox(header, cardThemeManager, cardBadgesAndFilters))
	})

	sidebarNavList := container.NewVBox(
		btnAll,
		btnThemesView,
		btnKPI,
		btnBtns,
		btnSliders,
		btnSpin,
		btnChips,
	)

	sidebarFooter := container.NewPadded(container.NewVBox(
		widget.NewSeparator(),
		container.NewCenter(uikit.NewCaptionLabel("Fyne 2.8 • Modern UI Toolkit", uikit.ThemeColor(uikit.ColorNameTextMuted))),
	))

	// Sidebar contenedor nativo UIKit
	sidebar = uikit.NewSidebarContainer(
		260,
		sidebarHeader,
		container.NewPadded(sidebarNavList),
		sidebarFooter,
	)

	// Layout con distribución Border: Sidebar fijo a la izquierda + Scroll central
	mainLayout := container.NewBorder(
		nil, nil,
		sidebar, nil,
		container.NewPadded(scrollArea),
	)

	// Envolver con fondo nativo reactivo
	rootStack := container.NewStack(appBg, mainLayout)
	myWindow.SetContent(rootStack)

	// Iniciar vigía de temas CSS en segundo plano (ThemeWatcher)
	watcher := uikit.NewThemeWatcher(themesDir, func(filePath string) {
		uikit.LoadAndRegisterExternalThemesDetailed(themesDir)
		applyTheme(uikit.GetCurrentThemeName(), false)
		uikit.ShowPopUp(myWindow.Canvas(), uikit.ToastInfo, fmt.Sprintf("CSS recargado en vivo: %s", filepath.Base(filePath)))
	})
	watcher.Start()

	myWindow.SetOnClosed(func() {
		watcher.Stop()
	})

	myWindow.ShowAndRun()
}

func uiSpacer(w, h float32) fyne.CanvasObject {
	rect := canvas.NewRectangle(nil)
	rect.SetMinSize(fyne.NewSize(w, h))
	return rect
}
