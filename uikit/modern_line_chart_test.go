package uikit

import (
	"fmt"
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
)

func TestModernLineChart_CreationAndRender(t *testing.T) {
	test.NewApp()

	minV := 0.0
	cfg := ChartConfig{
		XAxisTitle: "Meses de Anticipo",
		LeftAxis: ChartAxisConfig{
			Title: "Duración (Años)",
			Min:   &minV,
			Formatter: func(v float64) string {
				return fmt.Sprintf("%.1f a", v)
			},
		},
		RightAxis: &ChartAxisConfig{
			Title: "Edad (Años)",
			Formatter: func(v float64) string {
				return fmt.Sprintf("%.0f a", v)
			},
		},
		MinSize: fyne.NewSize(600, 300),
	}

	chart := NewModernLineChart(cfg)
	if chart == nil {
		t.Fatalf("Expected NewModernLineChart to return non-nil instance")
	}

	points := []ChartPoint{
		{XLabel: "24m", Title: "Anticipo: 24 meses", Subtitle: "Fecha: 01/01/2026"},
		{XLabel: "18m", Title: "Anticipo: 18 meses", Subtitle: "Fecha: 01/07/2026"},
		{XLabel: "12m", Title: "Anticipo: 12 meses", Subtitle: "Fecha: 01/01/2027"},
		{XLabel: "6m", Title: "Anticipo: 6 meses", Subtitle: "Fecha: 01/07/2027"},
		{XLabel: "1m", Title: "Anticipo: 1 mes", Subtitle: "Fecha: 01/12/2027"},
	}

	series := []ChartSeries{
		{
			Name:        "Duración Bruto",
			Color:       HexColor("#3b82f6"),
			Axis:        ChartAxisLeft,
			StrokeWidth: 2.2,
			Values:      []float64{10.5, 9.2, 7.8, 5.4, 2.1},
		},
		{
			Name:        "Edad Bruto",
			Color:       HexColor("#60a5fa"),
			Axis:        ChartAxisRight,
			IsDashed:    true,
			StrokeWidth: 2.0,
			Values:      []float64{75.5, 74.2, 72.8, 70.4, 67.1},
		},
	}

	chart.SetData(points, series)
	chart.Resize(fyne.NewSize(650, 320))
	chart.Refresh()

	// Probar interacción con el ratón
	chart.MouseIn(&desktop.MouseEvent{
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(250, 160)},
	})
	if !chart.isHovered {
		t.Errorf("Expected isHovered to be true after MouseIn")
	}
	if chart.hoverIndex < 0 {
		t.Errorf("Expected valid hoverIndex after MouseIn")
	}

	chart.MouseMoved(&desktop.MouseEvent{
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(450, 160)},
	})
	if !chart.isHovered {
		t.Errorf("Expected isHovered to be true after MouseMoved")
	}

	chart.MouseOut()
	if chart.isHovered {
		t.Errorf("Expected isHovered to be false after MouseOut")
	}
	if chart.hoverIndex != -1 {
		t.Errorf("Expected hoverIndex to be -1 after MouseOut")
	}

	// Probar vaciado de datos
	chart.Clear()
	if len(chart.points) != 0 || len(chart.series) != 0 {
		t.Errorf("Expected Clear to empty points and series")
	}
}

func TestModernLineChart_SingleAxis(t *testing.T) {
	test.NewApp()

	cfg := ChartConfig{
		XAxisTitle: "Fecha de Jubilación",
		LeftAxis: ChartAxisConfig{
			Title: "Pensión (€/mes)",
			Formatter: func(v float64) string {
				return fmt.Sprintf("%.0f €", v)
			},
		},
	}

	chart := NewModernLineChart(cfg)
	points := []ChartPoint{
		{XLabel: "01/26", Title: "Fecha: 01/2026"},
		{XLabel: "06/26", Title: "Fecha: 06/2026"},
		{XLabel: "12/26", Title: "Fecha: 12/2026"},
	}
	series := []ChartSeries{
		{
			Name:   "Pensión Bruta",
			Color:  color.RGBA{R: 200, G: 50, B: 50, A: 255},
			Axis:   ChartAxisLeft,
			Values: []float64{2100.0, 2150.0, 2200.0},
		},
	}

	chart.SetData(points, series)
	chart.Resize(fyne.NewSize(500, 250))
	chart.Refresh()

	if len(chart.points) != 3 {
		t.Errorf("Expected 3 points, got %d", len(chart.points))
	}
}

func TestModernLineChart_DynamicAutoScaling(t *testing.T) {
	test.NewApp()

	// Sin forzar Min (Min == nil)
	cfg := ChartConfig{
		XAxisTitle: "Fecha",
		LeftAxis: ChartAxisConfig{
			Title: "Pensión (€/mes)",
		},
	}

	chart := NewModernLineChart(cfg)
	points := []ChartPoint{
		{XLabel: "01/26"},
		{XLabel: "06/26"},
		{XLabel: "12/26"},
	}
	series := []ChartSeries{
		{
			Name:   "Pensión Bruta",
			Color:  color.RGBA{R: 50, G: 150, B: 250, A: 255},
			Axis:   ChartAxisLeft,
			Values: []float64{2000.0, 2200.0, 2400.0},
		},
	}

	chart.SetData(points, series)
	chart.Resize(fyne.NewSize(600, 300))
	chart.Refresh()

	// Comprobar que el renderer se ejecutó y generó los objetos
	renderer := test.WidgetRenderer(chart)
	if len(renderer.Objects()) < 5 {
		t.Errorf("Expected renderer to generate objects including grid and series lines")
	}
}

