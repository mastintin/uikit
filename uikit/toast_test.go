package uikit

import (
	"image/color"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestParseToastVariant(t *testing.T) {
	tests := []struct {
		input    any
		expected ToastVariant
	}{
		{"success", ToastSuccess},
		{"OK", ToastSuccess},
		{"error", ToastDanger},
		{"danger", ToastDanger},
		{"err", ToastDanger},
		{"warning", ToastWarning},
		{"warn", ToastWarning},
		{"alerta", ToastWarning},
		{"info", ToastInfo},
		{"unknown", ToastInfo},
		{ToastSuccess, ToastSuccess},
		{ToastDanger, ToastDanger},
		{ToastWarning, ToastWarning},
		{ToastInfo, ToastInfo},
		{12345, ToastInfo},
	}

	for _, tc := range tests {
		actual := ParseToastVariant(tc.input)
		if actual != tc.expected {
			t.Errorf("ParseToastVariant(%v) = %v, expected %v", tc.input, actual, tc.expected)
		}
	}
}

func TestModernToast_ShowAndHide(t *testing.T) {
	_ = test.NewApp()

	toast := NewModernToast()
	if toast == nil {
		t.Fatal("Expected NewModernToast to return non-nil instance")
	}

	container := toast.Container()
	if container.Visible() {
		t.Errorf("Expected toast container to be initially hidden")
	}

	toast.Show("success", "Operación guardada")
	if !container.Visible() {
		t.Errorf("Expected toast container to be visible after Show")
	}

	if toast.label.Text() != "Operación guardada" {
		t.Errorf("Expected label text 'Operación guardada', got '%s'", toast.label.Text())
	}

	toast.Hide()
	if container.Visible() {
		t.Errorf("Expected toast container to be hidden after Hide")
	}
}

func TestModernToast_CustomColorsAndDuration(t *testing.T) {
	_ = test.NewApp()

	toast := NewModernToast()
	customBorder := color.NRGBA{R: 255, G: 0, B: 255, A: 255}
	customIcon := color.NRGBA{R: 0, G: 255, B: 255, A: 255}
	customText := color.NRGBA{R: 255, G: 255, B: 0, A: 255}

	toast.SetCustomColors(ToastSuccess, customBorder, customIcon, customText)
	toast.SetDuration(100 * time.Millisecond)

	toast.Show(ToastSuccess, "Mensaje personalizado")
	if toast.label.Color() != customText {
		t.Errorf("Expected custom text color, got %v", toast.label.Color())
	}

	if toast.Card() == nil {
		t.Errorf("Expected Card() to return non-nil ModernCard")
	}
}

func TestModernToast_TappedDismissAndQueue(t *testing.T) {
	_ = test.NewApp()

	toast := NewModernToast()
	toast.SetDuration(5 * time.Second)

	// 1. Mostrar primer mensaje y encolar segundo
	toast.Show("info", "Primer Mensaje")
	if !toast.Container().Visible() {
		t.Errorf("Expected toast container to be visible")
	}

	toast.Queue("success", "Segundo Mensaje Encolado")

	// 2. Hacer clic sobre el toast para descartar el primer mensaje inmediatamente
	toast.tapWrapper.Tapped(&fyne.PointEvent{})

	// Debe haber avanzado al segundo mensaje de la cola
	if toast.label.Text() != "Segundo Mensaje Encolado" {
		t.Errorf("Expected 'Segundo Mensaje Encolado', got '%s'", toast.label.Text())
	}

	// 3. Descartar el segundo mensaje
	toast.tapWrapper.Tapped(&fyne.PointEvent{})
	if toast.Container().Visible() {
		t.Errorf("Expected toast container to be hidden after dismissing queue")
	}
}

func TestModernToast_ShowPopUp(t *testing.T) {
	_ = test.NewApp()
	win := test.NewWindow(nil)
	defer win.Close()

	pop := ShowPopUp(win.Canvas(), ToastWarning, "Aviso de Overlay", 100*time.Millisecond)
	if pop == nil {
		t.Fatal("Expected ShowPopUp to return non-nil PopUp")
	}
}

