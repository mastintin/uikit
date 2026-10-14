package uikit

import (
	"errors"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestAlerts(t *testing.T) {
	app := test.NewApp()
	app.Settings().SetTheme(NewAppTheme())
	win := app.NewWindow("Test Alerts")

	// 1. Probar formatAlertMessage
	if msg := formatAlertMessage(errors.New("fallo crítico")); msg != "fallo crítico" {
		t.Errorf("expected 'fallo crítico', got '%s'", msg)
	}
	if msg := formatAlertMessage("mensaje texto"); msg != "mensaje texto" {
		t.Errorf("expected 'mensaje texto', got '%s'", msg)
	}
	if msg := formatAlertMessage(nil); msg != "Se ha producido un error imprevisto." {
		t.Errorf("unexpected nil message format: %s", msg)
	}

	// 2. Probar MessageError (asíncrono)
	dErr := MessageError(win, errors.New("error de prueba"), "Fallo de Sistema")
	if dErr == nil {
		t.Fatal("expected dErr dialog not to be nil")
	}
	dErr.Hide()

	// 3. Probar MessageInfo (asíncrono)
	dInfo := MessageInfo(win, "Operación realizada con éxito", "Título Completo Sin Puntos")
	if dInfo == nil {
		t.Fatal("expected dInfo dialog not to be nil")
	}
	dInfo.Hide()

	// 4. Probar MessageYesNo (asíncrono)
	yesTriggered := false
	dYesNo := MessageYesNo(win, "¿Desea guardar los cambios?", func() {
		yesTriggered = true
	})
	if dYesNo == nil {
		t.Fatal("expected dYesNo dialog not to be nil")
	}
	dYesNo.Hide()
	if yesTriggered {
		t.Error("yesTriggered should be false before clicking Yes")
	}

	// 5. Probar MessageConfirm (asíncrono)
	confirmTriggered := false
	dConfirm := MessageConfirm(win, "¿Desea eliminar este registro?", "Eliminar", "Cancelar", true, func() {
		confirmTriggered = true
	}, "Peligro")
	if dConfirm == nil {
		t.Fatal("expected dConfirm dialog not to be nil")
	}
	dConfirm.Hide()
	if confirmTriggered {
		t.Error("confirmTriggered should be false before clicking Confirm")
	}

	// 6. Probar Action y Async (Hilo 2)
	doneAction := make(chan bool, 1)
	wrappedAction := Action(func() {
		doneAction <- true
	})
	wrappedAction()

	select {
	case <-doneAction:
		// OK
	case <-time.After(500 * time.Millisecond):
		t.Error("Action did not execute in background goroutine")
	}

	doneAsync := make(chan bool, 1)
	Async(func() {
		doneAsync <- true
	})
	select {
	case <-doneAsync:
		// OK
	case <-time.After(500 * time.Millisecond):
		t.Error("Async did not execute in background goroutine")
	}

	doneSync := make(chan bool, 1)
	Sync(func() {
		doneSync <- true
	})
	select {
	case <-doneSync:
		// OK
	case <-time.After(500 * time.Millisecond):
		t.Error("Sync did not enqueue into UI loop")
	}

	doneDo := make(chan bool, 1)
	Do(func() {
		doneDo <- true
	})
	select {
	case <-doneDo:
		// OK
	case <-time.After(500 * time.Millisecond):
		t.Error("Do did not enqueue into UI loop")
	}

	// 7. Probar formatDebugMessage
	if dbg := formatDebugMessage("A", 123, true); dbg != "A 123 true" {
		t.Errorf("unexpected debug string: '%s'", dbg)
	}
	if dbg := formatDebugMessage(); dbg != "(sin datos)" {
		t.Errorf("unexpected empty debug string: '%s'", dbg)
	}

	// 8. Probar llamadas síncronas seguras con win == nil
	if MsgYesNo(nil, "pregunta") != false {
		t.Error("expected false with nil window")
	}
	if MsgConfirm(nil, "pregunta", "Sí", "No", false) != false {
		t.Error("expected false with nil window")
	}
	if text, ok := MsgGet(nil, "prompt"); ok || text != "" {
		t.Error("expected empty result with nil window")
	}
	MsgInfo(nil, "info")
	MsgError(nil, "error")
	MsgAlert(nil, "aviso")
	MsgStop(nil, "alto")
}

func TestFormDialog(t *testing.T) {
	app := test.NewApp()
	app.Settings().SetTheme(NewLightTheme())
	win := app.NewWindow("Test Form Dialog")

	entry := widget.NewEntry()
	items := []*widget.FormItem{
		widget.NewFormItem("Campo", entry),
	}

	confirmed := false
	d := NewFormDialog("Prueba Formulario", "Guardar", "Cancelar", items, func(ok bool) {
		confirmed = ok
	}, win)

	if d == nil {
		t.Fatal("expected non-nil FormDialog")
	}

	d.Show()
	d.Hide()
	if confirmed {
		t.Error("confirmed should remain false when dialog was hidden without clicking save")
	}

	d2 := ShowFormDialog("Prueba Show", "Guardar", "Cancelar", items, func(ok bool) {}, win)
	if d2 == nil {
		t.Fatal("expected non-nil ShowFormDialog")
	}
	d2.Hide()
}

func TestFormDialog_ConfirmExecution(t *testing.T) {
	app := test.NewApp()
	app.Settings().SetTheme(NewLightTheme())
	win := app.NewWindow("Test Confirm")

	entry := widget.NewEntry()
	items := []*widget.FormItem{
		widget.NewFormItem("Campo", entry),
	}

	confirmed := false
	d := NewFormDialog("Prueba", "Guardar", "Cancelar", items, func(ok bool) {
		confirmed = ok
	}, win)
	d.Show()

	overlay := win.Canvas().Overlays().Top()
	if overlay == nil {
		t.Fatal("expected dialog overlay")
	}

	var confirmBtn *ModernButton
	var findBtn func(co fyne.CanvasObject)
	findBtn = func(co fyne.CanvasObject) {
		if co == nil {
			return
		}
		if btn, ok := co.(*ModernButton); ok && btn.Text == "Guardar" {
			confirmBtn = btn
			return
		}
		if pop, ok := co.(*widget.PopUp); ok {
			findBtn(pop.Content)
		}
		if c, ok := co.(*fyne.Container); ok {
			for _, child := range c.Objects {
				findBtn(child)
			}
		}
		if mc, ok := co.(*ModernCard); ok {
			findBtn(mc.Content())
		}
		if w, ok := co.(fyne.Widget); ok {
			r := test.WidgetRenderer(w)
			for _, child := range r.Objects() {
				findBtn(child)
			}
		}
	}
	findBtn(overlay)

	if confirmBtn == nil {
		t.Fatal("could not find confirm button in dialog")
	}

	test.Tap(confirmBtn)

	if !confirmed {
		t.Error("expected confirmed to be true when clicking Guardar")
	}
}


