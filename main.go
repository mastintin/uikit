package main

import (
	"fmt"
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

func main() {
	myApp := app.NewWithID("com.uikit.showcase")
	myWindow := myApp.NewWindow("UIKit Component Gallery & CSS Themes - Fyne 2.8")
	myWindow.Resize(fyne.NewSize(1260, 860))

	// -------------------------------------------------------------------------
	// 0. Carga e Inicialización de Temas CSS Externos
	// -------------------------------------------------------------------------
	uikit.LoadAndRegisterExternalThemesDetailed("themes")
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
	// 2. Tarjetas KPI / Métricas (KPICard)
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
	// 3. Sección de Botones Modernos (ModernButton con Loading & Glass)
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
	// 4. Sección de Deslizadores Modernos (ModernSlider)
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
	// 5. Sección de Indicadores de Carga (ModernSpinner)
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
	// 6. Badges, Pastillas e Interruptores (ModernBadge, Switch)
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
	// 7. Filtros Desplegables Tipo Pastilla (ModernChipFilter)
	// -------------------------------------------------------------------------
	chipYear := uikit.NewChipFilter("Año:", []string{"Todos los años", "2026", "2025", "2024", "2023"}, "Todos los años", func(sel string) {
		uikit.ShowPopUp(myWindow.Canvas(), uikit.ToastInfo, "Filtro activo: "+sel)
	})

	chipStatus := uikit.NewChipFilter("Estado:", []string{"Todos", "Activo", "Pendiente", "Archivado"}, "Todos", func(sel string) {
		uikit.ShowPopUp(myWindow.Canvas(), uikit.ToastInfo, "Estado filtrado: "+sel)
	})

	var chipThemeTop *uikit.ModernChipFilter

	chipFiltersRow := container.NewHBox(
		widget.NewLabel("Filtros Interactivos:"),
		chipYear,
		chipStatus,
	)

	// -------------------------------------------------------------------------
	// 8. Selector de Fechas (ModernDatePicker)
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
	// 9. Contenedor Dinámico Principal (Centro)
	// -------------------------------------------------------------------------
	allContent := container.NewVBox(
		header,
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
	// 10. Construcción de ModernSidebar (Barra Lateral Izquierda)
	// -------------------------------------------------------------------------
	brandIcon := uikit.NewColoredIcon("sliders", 22, "#107c41")
	lblBrand := canvas.NewText("UIKit Gallery", uikit.ThemeColor(uikit.ColorNameTextPrimary))
	lblBrand.TextStyle = fyne.TextStyle{Bold: true}
	lblBrand.TextSize = 15
	badgeVersion := uikit.NewBadge(uikit.BadgeSuccess, "v2.8")

	sidebarHeader := container.NewPadded(container.NewVBox(
		container.NewHBox(
			brandIcon,
			lblBrand,
			layout.NewSpacer(),
			badgeVersion,
		),
		uikit.NewCaptionLabel("Catálogo Oficial & Gestor CSS", uikit.ThemeColor(uikit.ColorNameTextMuted)),
		widget.NewSeparator(),
	))

	// Botones de Navegación del Sidebar
	var navBtns []*uikit.ModernButton
	addNav := func(title string, ic fyne.Resource, action func()) *uikit.ModernButton {
		var btn *uikit.ModernButton
		btn = uikit.NewButton("Default", title, ic, false, func() {
			for _, b := range navBtns {
				b.SetVariant(uikit.ButtonDefault)
			}
			btn.SetVariant(uikit.ButtonPrimary)
			action()
		})
		navBtns = append(navBtns, btn)
		return btn
	}

	btnAll := addNav("Vista Completa", theme.HomeIcon(), func() {
		setView(allContent)
	})
	btnAll.SetVariant(uikit.ButtonPrimary) // Seleccionado por defecto

	btnKPI := addNav("Métricas & KPIs", theme.ComputerIcon(), func() {
		setView(container.NewVBox(header, cardKPIs))
	})

	btnBtns := addNav("Botones & Carga", theme.MediaPlayIcon(), func() {
		setView(container.NewVBox(header, cardButtons))
	})

	btnSliders := addNav("Deslizadores (Sliders)", theme.MoveDownIcon(), func() {
		setView(container.NewVBox(header, cardSliders))
	})

	btnSpin := addNav("Indicadores (Spinners)", theme.ViewRefreshIcon(), func() {
		setView(container.NewVBox(header, cardSpinners))
	})

	btnChips := addNav("Filtros & Insignias", theme.MenuIcon(), func() {
		setView(container.NewVBox(header, cardBadgesAndFilters))
	})

	sidebarNavList := container.NewVBox(
		btnAll,
		btnKPI,
		btnBtns,
		btnSliders,
		btnSpin,
		btnChips,
	)

	// -------------------------------------------------------------------------
	// 11. Selector y Gestor de Temas CSS en Vivo
	// -------------------------------------------------------------------------
	var sidebarThemeSelect *widget.Select
	var sidebar *uikit.ModernSidebar

	applyTheme := func(themeID string, showFeedback bool) {
		th := uikit.SetAppTheme(themeID)

		// 1. Fondo de la ventana sincronizado
		appBg.FillColor = uikit.ThemeColor(uikit.ColorNameAppBackground)
		appBg.Refresh()

		// 2. Elementos gráficos personalizados del sidebar
		lblBrand.Color = uikit.ThemeColor(uikit.ColorNameTextPrimary)
		lblBrand.Refresh()

		// 3. Refrescar recursivamente todo el árbol de controles
		uikit.RefreshTree(myWindow.Content())
		myWindow.Content().Refresh()

		// 4. Sincronizar selectores
		_, _, disp := getThemesList()
		if sidebarThemeSelect != nil && sidebarThemeSelect.Selected != disp {
			sidebarThemeSelect.Selected = disp
			sidebarThemeSelect.Refresh()
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

	sidebarThemeSelect = widget.NewSelect(themeDisplays, func(selected string) {
		if id, ok := themeMap[selected]; ok {
			applyTheme(id, true)
		}
	})
	sidebarThemeSelect.Selected = currentDisplay

	// Chip Filter en la barra de filtros superior para alternar temas directamente
	chipThemeTop = uikit.NewChipFilter("Tema CSS:", themeDisplays, currentDisplay, func(selected string) {
		if id, ok := themeMap[selected]; ok {
			applyTheme(id, true)
		}
	})
	chipFiltersRow.Add(chipThemeTop)

	btnReloadThemes := uikit.NewButton("Primary", "Recargar Temas CSS", theme.ViewRefreshIcon(), false, func() {
		report := uikit.LoadAndRegisterExternalThemesDetailed("themes")
		td, tm, cur := getThemesList()
		themeDisplays = td
		themeMap = tm

		sidebarThemeSelect.Options = themeDisplays
		sidebarThemeSelect.Selected = cur
		sidebarThemeSelect.Refresh()

		chipThemeTop.SetOptions(themeDisplays)
		chipThemeTop.SetSelectedSilent(cur)

		applyTheme(uikit.GetCurrentThemeName(), false)

		if len(report.Errors) > 0 {
			uikit.ShowPopUp(myWindow.Canvas(), uikit.ToastDanger, fmt.Sprintf("Error CSS: %s", report.Errors[0]))
		} else {
			uikit.ShowPopUp(myWindow.Canvas(), uikit.ToastSuccess, fmt.Sprintf("%d temas CSS registrados en caliente", report.Count))
		}
	})

	switchGlass := uikit.NewSwitch("Efecto Cristal (Glass)", false, func(active bool) {
		uikit.SetGlobalGlassButtons(active)
		uikit.RefreshTree(myWindow.Content())
		myWindow.Content().Refresh()
	})

	sidebarFooter := container.NewPadded(container.NewVBox(
		widget.NewSeparator(),
		uikit.NewCaptionLabel("TEMA VISUAL & CSS:", uikit.ThemeColor(uikit.ColorNameTextMuted)),
		sidebarThemeSelect,
		uiSpacer(1, 4),
		btnReloadThemes,
		switchGlass,
		container.NewCenter(uikit.NewCaptionLabel("Fyne 2.8 • Modern UI Toolkit", uikit.ThemeColor(uikit.ColorNameTextMuted))),
	))

	// Sidebar contenedor nativo UIKit
	sidebar = uikit.NewSidebarContainer(
		255,
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

	// Envolver con fondo nativo
	rootStack := container.NewStack(appBg, mainLayout)
	myWindow.SetContent(rootStack)

	// Iniciar vigía de temas CSS en segundo plano (ThemeWatcher)
	watcher := uikit.NewThemeWatcher("themes", func(filePath string) {
		uikit.LoadAndRegisterExternalThemesDetailed("themes")
		applyTheme(uikit.GetCurrentThemeName(), false)
		uikit.ShowPopUp(myWindow.Canvas(), uikit.ToastInfo, fmt.Sprintf("Tema CSS recargado en vivo (%s)", filepath.Base(filePath)))
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
