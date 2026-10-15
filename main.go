package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"uikit-demo/uikit"
)

func main() {
	myApp := app.NewWithID("com.uikit.showcase")
	myWindow := myApp.NewWindow("UIKit Component Gallery - Fyne 2.8")
	myWindow.Resize(fyne.NewSize(1040, 800))

	// Configurar tema claro por defecto
	currentTheme := uikit.NewAppTheme()
	myApp.Settings().SetTheme(currentTheme)

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

	switchGlass := uikit.NewSwitch("Efecto Cristal (Glass)", false, func(active bool) {
		uikit.SetGlobalGlassButtons(active)
	})

	buttonsRow := container.NewHBox(
		btnLoading,
		btnPrimary,
		btnDanger,
		btnDefault,
		btnSmall,
		btnIcon,
		switchGlass,
	)

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

	// -------------------------------------------------------------------------
	// 8. Barra de Control de Tema (Light / Dark)
	// -------------------------------------------------------------------------
	switchTheme := uikit.NewSwitch("Tema Oscuro (Dark Mode)", false, func(dark bool) {
		if dark {
			myApp.Settings().SetTheme(uikit.NewDarkTheme())
		} else {
			myApp.Settings().SetTheme(uikit.NewLightTheme())
		}
	})

	themeBar := container.NewHBox(
		widget.NewLabel("Ajuste de Entorno:"),
		switchTheme,
	)

	// -------------------------------------------------------------------------
	// Ensamblado en Tarjetas de Sección
	// -------------------------------------------------------------------------
	cardButtons := uikit.NewCard(container.NewVBox(
		uikit.NewSectionTitle("Botones & Estados de Carga (ModernButton)"),
		buttonsRow,
	))

	cardSliders := uikit.NewCard(container.NewVBox(
		uikit.NewSectionTitle("Deslizadores Numéricos Reactivos (ModernSlider)"),
		slidersBox,
	))

	cardSpinners := uikit.NewCard(container.NewVBox(
		uikit.NewSectionTitle("Indicadores Circulares de Actividad (ModernSpinner)"),
		spinnersRow,
	))

	cardBadges := uikit.NewCard(container.NewVBox(
		uikit.NewSectionTitle("Insignias, Filtros, Interruptores & Fechas"),
		container.NewVBox(badgesRow, chipFiltersRow, datePickerRow),
	))

	content := container.NewVBox(
		header,
		themeBar,
		kpiGrid,
		cardButtons,
		cardSliders,
		cardSpinners,
		cardBadges,
	)

	// Contenedor con scroll vertical para navegación fluida
	scroll := container.NewVScroll(content)

	myWindow.SetContent(container.NewPadded(scroll))
	myWindow.ShowAndRun()
}
