# UIKit para Go & Fyne

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Fyne Version](https://img.shields.io/badge/Fyne-v2.8+-2563EB?style=flat)](https://fyne.io)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Biblioteca visual y catálogo de componentes modernos, responsivos y fluidos para el desarrollo de aplicaciones de escritorio multiplataforma (macOS, Windows, Linux) utilizando **Go** y **[Fyne v2](https://fyne.io)**.

---

## 🎯 Finalidad y Filosofía

Fyne proporciona una base sólida y multiplataforma para interfaces gráficas en Go, pero los proyectos empresariales y modernos a menudo requieren una estética más contemporánea, patrones de diseño avanzados y componentes listos para usar sin tener que diseñarlos desde cero.

**UIKit** amplía el ecosistema de Fyne ofreciendo:
- **Estética Contemporánea**: Bordes redondeados, paletas de colores armónicas, efectos de cristal (*glassmorphism*) y tipografía refinada.
- **Componentes Reactivos de Alto Nivel**: Deslizadores fluidos con soporte para divisas, porcentajes y rangos con formato en tiempo real, selectores de fecha con calendario emergente, tarjetas de métricas (KPIs), gráficas vectoriales y botones con spinners de carga asíncronos integrados.
- **Sistema de Temas Dinámico & CSS**: Soporte completo para temas Claro, Oscuro y personalización avanzada basada en variables de estilo CSS (`nord`, `emerald`, `dark`, `nice`), con capacidad de recarga en tiempo real (*hot-reload*).

---

## ✨ Componentes Incluidos

| Componente | Archivo | Descripción |
| :--- | :--- | :--- |
| **ModernButton** | [`uikit/button.go`](file:///Users/manuel/antiProyectos/uikit/uikit/button.go) | Botones primarios, de peligro, secundarios, compactos e iconos con soporte de estados de carga (`SetLoading`) y modo *glass*. |
| **KPICard / MetricCard** | [`uikit/kpi_card.go`](file:///Users/manuel/antiProyectos/uikit/uikit/kpi_card.go), [`uikit/metric_card.go`](file:///Users/manuel/antiProyectos/uikit/uikit/metric_card.go) | Tarjetas de métricas con indicadores visuales de tendencia (subida/bajada), variantes cromáticas (Emerald, Blue, Violet, Amber). |
| **ModernSlider** | [`uikit/slider.go`](file:///Users/manuel/antiProyectos/uikit/uikit/slider.go) | Deslizadores reactivos con API fluida (`NewAgeSlider`, `NewPercentageSlider`, `NewCurrencySlider`) con formateo automático, pasos y precisión. |
| **ModernSpinner** | [`uikit/spinner.go`](file:///Users/manuel/antiProyectos/uikit/uikit/spinner.go) | Indicador circular de progreso/actividad con 4 tamaños (`Small`, `Medium`, `Large`, `Hero`) y variantes semánticas. |
| **ModernDatePicker** | [`uikit/datepicker.go`](file:///Users/manuel/antiProyectos/uikit/uikit/datepicker.go) | Selector de fechas intuitivo con ventana emergente de calendario interactivo ([`calendar_popup.go`](file:///Users/manuel/antiProyectos/uikit/uikit/calendar_popup.go)). |
| **ModernBadge & Chips** | [`uikit/badge.go`](file:///Users/manuel/antiProyectos/uikit/uikit/badge.go), [`uikit/sort_chip.go`](file:///Users/manuel/antiProyectos/uikit/uikit/sort_chip.go) | Insignias de estado (Success, Warning, Danger, Primary, Neutral) y chips interactivos para ordenación y filtrado. |
| **ModernSwitch** | [`uikit/switch.go`](file:///Users/manuel/antiProyectos/uikit/uikit/switch.go) | Interruptor estilo toggle animado y personalizable. |
| **ModernLineChart** | [`uikit/modern_line_chart.go`](file:///Users/manuel/antiProyectos/uikit/uikit/modern_line_chart.go) | Gráfico de líneas y series temporales trazado vectorialmente. |
| **Toasts, Alerts & Popups** | [`uikit/toast.go`](file:///Users/manuel/antiProyectos/uikit/uikit/toast.go), [`uikit/alerts.go`](file:///Users/manuel/antiProyectos/uikit/uikit/alerts.go) | Mensajes flotantes temporizados (Toasts) y diálogos de confirmación / aviso estilizados. |
| **HeaderCard & Cards** | [`uikit/header_card.go`](file:///Users/manuel/antiProyectos/uikit/uikit/header_card.go), [`uikit/card.go`](file:///Users/manuel/antiProyectos/uikit/uikit/card.go) | Contenedores visuales con sombras suaves, degradados de título y divisores estructurales. |
| **Segmented & Listbox** | [`uikit/segmented.go`](file:///Users/manuel/antiProyectos/uikit/uikit/segmented.go), [`uikit/listbox.go`](file:///Users/manuel/antiProyectos/uikit/uikit/listbox.go) | Selector por segmentos estilo píldora y listas seleccionables modernas. |
| **Form Fields & Browse** | [`uikit/form_field.go`](file:///Users/manuel/antiProyectos/uikit/uikit/form_field.go), [`uikit/browse.go`](file:///Users/manuel/antiProyectos/uikit/uikit/browse.go) | Campos de formulario con validación y selectores de archivos con diseño integrado. |

---

## 🎨 Sistema de Temas y CSS

La biblioteca incluye temas preconfigurados y un cargador de temas desde archivos de especificación CSS:
- **Tema Claro** (`uikit.NewLightTheme()`)
- **Tema Oscuro** (`uikit.NewDarkTheme()`)
- **Temas CSS incluidos en `themes/`**:
  - `nord.css`: Paleta nórdica ártica oficial (Polar Night, Snow Storm, Frost).
  - `emerald.css`: Acentos esmeralda y tonos frescos para paneles corporativos.
  - `dark.css`: Modo oscuro puro de alto contraste.
  - `nice.css`: Estilo equilibrado y luminoso.
- Soporte de observación de cambios de archivo en tiempo de ejecución para desarrollo interactivo (`theme_watcher.go`).

---

## 🚀 Inicio Rápido

### Prerrequisitos
- **Go** 1.21 o superior
- Dependencias de Fyne instaladas en tu sistema (OpenGL/compilador C de tu plataforma). Consulta la [guía de requisitos de Fyne](https://docs.fyne.io/started/#prerequisites).

### Ejecutar la Galería de Demostración

El repositorio incluye un catálogo interactivo completo en [`main.go`](file:///Users/manuel/antiProyectos/uikit/main.go) que muestra todos los componentes en acción:

```bash
git clone https://github.com/mastintin/uikit.git
cd uikit
go run main.go
```

### Ejemplo Básico de Uso

```go
package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"

	"uikit-demo/uikit"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Ejemplo UIKit")

	// Aplicar tema moderno
	myApp.Settings().SetTheme(uikit.NewAppTheme())

	// Botón con estado de carga
	var btn *uikit.ModernButton
	btn = uikit.NewPrimaryButton("Guardar Datos", theme.DocumentSaveIcon(), func() {
		btn.SetLoading(true)
		go func() {
			time.Sleep(2 * time.Second)
			btn.SetLoading(false)
			uikit.ShowPopUp(myWindow.Canvas(), uikit.ToastSuccess, "¡Guardado con éxito!")
		}()
	})

	// Slider reactivo de divisa
	slider := uikit.NewCurrencySlider(500, 5000, 50).
		WithTitle("Presupuesto Objetivo").
		WithValue(1500)

	// Tarjeta de métrica
	kpi := uikit.NewKPICard(uikit.KPICardEmerald, "Rendimiento", "98.4%", "+2.3% vs mes anterior")

	myWindow.SetContent(container.NewVBox(
		kpi,
		slider,
		btn,
	))

	myWindow.Resize(fyne.NewSize(500, 400))
	myWindow.ShowAndRun()
}
```

---

## 📂 Estructura del Proyecto

```text
uikit/
├── main.go               # Aplicación de catálogo interactivo / demo
├── go.mod                # Declaración del módulo y dependencias
├── go.sum                # Checksums de dependencias Go
├── themes/               # Archivos CSS de temas (Nord, Emerald, Dark, Nice)
│   ├── dark.css
│   ├── emerald.css
│   ├── nice.css
│   └── nord.css
└── uikit/                # Paquete con los componentes visuales
    ├── button.go         # Botones modernos y variantes
    ├── kpi_card.go       # Tarjetas KPI y métricas
    ├── slider.go         # Deslizadores numéricos reactivos
    ├── spinner.go        # Indicadores de carga circulares
    ├── datepicker.go     # Selector de fechas y calendario
    ├── theme*.go         # Gestión de temas, paletas y CSS
    └── ...               # Más de 30 controles y pruebas unitarias
```

---

## 📄 Licencia

Este proyecto se distribuye bajo los términos de la licencia de código abierto. Consulta el archivo de licencia para más detalles.
