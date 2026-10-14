package uikit

import (
	"fmt"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type alertKind int

const (
	alertKindError alertKind = iota
	alertKindInfo
	alertKindYesNo
	alertKindWarning
)

// =============================================================================
// FAMILIA ASÍNCRONA (Bajo nivel Fyne, no bloqueante, retorno de dialog.Dialog)
// =============================================================================

// MessageError muestra un diálogo modal de error y retorna inmediatamente dialog.Dialog.
func MessageError(win fyne.Window, message any, title ...string) dialog.Dialog {
	t := "Error"
	if len(title) > 0 && strings.TrimSpace(title[0]) != "" {
		t = strings.TrimSpace(title[0])
	}
	msgStr := formatAlertMessage(message)
	return showModernAlertInternal(win, alertKindError, t, msgStr, "Aceptar", "", false, nil, nil)
}

// MessageInfo muestra un diálogo modal informativo y retorna inmediatamente dialog.Dialog.
func MessageInfo(win fyne.Window, message string, title ...string) dialog.Dialog {
	t := "Información"
	if len(title) > 0 && strings.TrimSpace(title[0]) != "" {
		t = strings.TrimSpace(title[0])
	}
	return showModernAlertInternal(win, alertKindInfo, t, message, "Aceptar", "", false, nil, nil)
}

// MessageYesNo muestra un diálogo de confirmación asíncrono con callback onYes.
func MessageYesNo(win fyne.Window, message string, onYes func(), title ...string) dialog.Dialog {
	t := "Confirmación"
	if len(title) > 0 && strings.TrimSpace(title[0]) != "" {
		t = strings.TrimSpace(title[0])
	}
	return showModernAlertInternal(win, alertKindYesNo, t, message, "Sí", "No", false, onYes, nil)
}

// MessageConfirm permite personalizar etiquetas de botones con callback onConfirm.
func MessageConfirm(win fyne.Window, message string, confirmLabel, cancelLabel string, isDanger bool, onConfirm func(), title ...string) dialog.Dialog {
	t := "Confirmación"
	if len(title) > 0 && strings.TrimSpace(title[0]) != "" {
		t = strings.TrimSpace(title[0])
	}
	return showModernAlertInternal(win, alertKindYesNo, t, message, confirmLabel, cancelLabel, isDanger, onConfirm, nil)
}

// =============================================================================
// FAMILIA SÍNCRONA FIVEWIN (Para Hilo 2, bloqueante mediante canales Go)
// =============================================================================

// MsgInfo muestra un diálogo informativo estilo FiveWin sin bloquear el hilo principal.
func MsgInfo(win fyne.Window, message string, title ...string) {
	if win == nil {
		return
	}
	t := "Información"
	if len(title) > 0 && strings.TrimSpace(title[0]) != "" {
		t = strings.TrimSpace(title[0])
	}

	showModernAlertInternal(win, alertKindInfo, t, message, "Aceptar", "", false, nil, nil)
}

// MsgError muestra un diálogo de error estilo FiveWin sin bloquear el hilo principal.
func MsgError(win fyne.Window, message any, title ...string) {
	if win == nil {
		return
	}
	t := "Error"
	if len(title) > 0 && strings.TrimSpace(title[0]) != "" {
		t = strings.TrimSpace(title[0])
	}
	msgStr := formatAlertMessage(message)

	showModernAlertInternal(win, alertKindError, t, msgStr, "Aceptar", "", false, nil, nil)
}

// MsgAlert muestra un diálogo de advertencia estilo FiveWin sin bloquear el hilo principal.
func MsgAlert(win fyne.Window, message string, title ...string) {
	if win == nil {
		return
	}
	t := "Aviso"
	if len(title) > 0 && strings.TrimSpace(title[0]) != "" {
		t = strings.TrimSpace(title[0])
	}

	showModernAlertInternal(win, alertKindWarning, t, message, "Aceptar", "", false, nil, nil)
}

// MsgStop muestra un diálogo de parada/alto crítico sin bloquear el hilo principal.
func MsgStop(win fyne.Window, message string, title ...string) {
	if win == nil {
		return
	}
	t := "Alto"
	if len(title) > 0 && strings.TrimSpace(title[0]) != "" {
		t = strings.TrimSpace(title[0])
	}

	showModernAlertInternal(win, alertKindError, t, message, "Aceptar", "", true, nil, nil)
}

// MsgYesNo muestra un diálogo de confirmación estilo FiveWin y BLOQUEA la goroutine (Hilo 2).
// Retorna true si el usuario pulsó "Sí" o false si pulsó "No" / cerró el diálogo.
func MsgYesNo(win fyne.Window, message string, title ...string) bool {
	t := "Confirmación"
	if len(title) > 0 && strings.TrimSpace(title[0]) != "" {
		t = strings.TrimSpace(title[0])
	}
	return MsgConfirm(win, message, "Sí", "No", false, t)
}

// MsgConfirm muestra un diálogo de confirmación con etiquetas personalizadas y retorno booleano directo.
func MsgConfirm(win fyne.Window, message string, confirmLabel, cancelLabel string, isDanger bool, title ...string) bool {
	if win == nil {
		return false
	}
	t := "Confirmación"
	if len(title) > 0 && strings.TrimSpace(title[0]) != "" {
		t = strings.TrimSpace(title[0])
	}

	ch := make(chan bool, 1)
	var once sync.Once
	sendResult := func(val bool) {
		once.Do(func() {
			ch <- val
		})
	}

	d := showModernAlertInternal(win, alertKindYesNo, t, message, confirmLabel, cancelLabel, isDanger,
		func() { sendResult(true) },
		func() { sendResult(false) },
	)
	if d == nil {
		return false
	}
	d.SetOnClosed(func() {
		sendResult(false)
	})

	return <-ch
}

// MsgGet muestra un cuadro modal para introducir un texto estilo FiveWin y BLOQUEA la goroutine.
// Parámetros opcionales en title:
//   - title[0]: Título de la ventana
//   - title[1]: Texto inicial por defecto
//
// Retorna el texto ingresado y true si el usuario confirmó con Aceptar, o false si Canceló.
func MsgGet(win fyne.Window, prompt string, title ...string) (string, bool) {
	if win == nil {
		return "", false
	}
	type result struct {
		text string
		ok   bool
	}
	ch := make(chan result, 1)
	var once sync.Once
	sendResult := func(res result) {
		once.Do(func() {
			ch <- res
		})
	}

	t := "Entrada de datos"
	if len(title) > 0 && strings.TrimSpace(title[0]) != "" {
		t = strings.TrimSpace(title[0])
	}
	initialText := ""
	if len(title) > 1 {
		initialText = title[1]
	}

	entry := widget.NewEntry()
	entry.SetPlaceHolder(prompt)
	if initialText != "" {
		entry.SetText(initialText)
	}

	var d dialog.Dialog

	btnCancel := NewButton("Default", "Cancelar", "cancel", false, func() {
		sendResult(result{text: "", ok: false})
		if d != nil {
			d.Hide()
		}
	})

	btnConfirm := NewButton("Primary", "Aceptar", "check", false, func() {
		sendResult(result{text: entry.Text, ok: true})
		if d != nil {
			d.Hide()
		}
	})

	entry.OnSubmitted = func(_ string) {
		sendResult(result{text: entry.Text, ok: true})
		if d != nil {
			d.Hide()
		}
	}

	iconBadge := createAlertIconBadge("#eff6ff", "#bfdbfe", theme.ColorNamePrimary, theme.InfoIcon())
	titleLabel := widget.NewLabelWithStyle(t, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	header := container.NewHBox(iconBadge, container.NewCenter(titleLabel))

	msgPrompt := widget.NewLabel(prompt)
	msgPrompt.Wrapping = fyne.TextWrapWord

	buttonBar := container.NewHBox(btnCancel, btnConfirm)

	contentBox := container.NewVBox(
		header,
		widget.NewSeparator(),
		container.NewPadded(container.NewVBox(
			msgPrompt,
			entry,
		)),
		widget.NewSeparator(),
		container.NewBorder(nil, nil, nil, buttonBar),
	)

	card := NewPanel(contentBox, "", "", 12, 1)
	d = dialog.NewCustomWithoutButtons("", card, win)
	if d == nil {
		return "", false
	}
	d.SetOnClosed(func() {
		sendResult(result{text: "", ok: false})
	})
	d.Resize(fyne.NewSize(460, 220))
	fyne.Do(func() {
		d.Show()
		if win != nil && win.Canvas() != nil {
			win.Canvas().Focus(entry)
		}
	})

	res := <-ch
	return res.text, res.ok
}

// formatAlertMessage convierte cualquier tipo de mensaje en string formateado.
func formatAlertMessage(message any) string {
	switch m := message.(type) {
	case nil:
		return "Se ha producido un error imprevisto."
	case error:
		return m.Error()
	case string:
		return m
	case fmt.Stringer:
		return m.String()
	default:
		return fmt.Sprintf("%v", m)
	}
}

// formatDebugMessage concatena argumentos para MsgDebug.
func formatDebugMessage(values ...any) string {
	if len(values) == 0 {
		return "(sin datos)"
	}
	var parts []string
	for _, v := range values {
		parts = append(parts, fmt.Sprint(v))
	}
	return strings.Join(parts, " ")
}

// showModernAlertInternal construye y proyecta la tarjeta de diálogo modal estilizada.
func showModernAlertInternal(
	win fyne.Window,
	kind alertKind,
	title string,
	message string,
	confirmText string,
	cancelText string,
	isDanger bool,
	onConfirm func(),
	onCancel func(),
) dialog.Dialog {
	if win == nil {
		return nil
	}

	var d dialog.Dialog

	// 1. Configurar icono semántico y colores según el tipo de alerta
	var iconRes fyne.Resource
	var bgHex, borderHex string
	var iconColorName fyne.ThemeColorName

	switch kind {
	case alertKindError:
		iconRes = theme.ErrorIcon()
		bgHex = ColorToHex(ColorWarningBg)
		borderHex = ColorToHex(ColorWarningBorder)
		iconColorName = theme.ColorNameError
	case alertKindInfo:
		iconRes = theme.InfoIcon()
		bgHex = "#eff6ff"     // Pastel suave azul
		borderHex = "#bfdbfe" // Borde azul suave
		iconColorName = theme.ColorNamePrimary
	case alertKindWarning:
		iconRes = theme.WarningIcon()
		bgHex = ColorToHex(ColorWarningBg)
		borderHex = ColorToHex(ColorWarningBorder)
		iconColorName = theme.ColorNameWarning
	case alertKindYesNo:
		if isDanger {
			iconRes = theme.DeleteIcon()
			bgHex = ColorToHex(ColorWarningBg)
			borderHex = ColorToHex(ColorWarningBorder)
			iconColorName = theme.ColorNameError
		} else {
			iconRes = theme.QuestionIcon()
			bgHex = ColorToHex(ColorWarningBg)
			borderHex = ColorToHex(ColorWarningBorder)
			iconColorName = theme.ColorNameWarning
		}
	}

	iconBadge := createAlertIconBadge(bgHex, borderHex, iconColorName, iconRes)

	// 2. Cabecera con título
	titleLabel := NewSectionLabel(title)

	header := container.NewHBox(
		iconBadge,
		container.NewCenter(titleLabel),
	)

	// 3. Mensaje multilínea con salto de línea automático
	msgLabel := NewBodyLabel(message).WithWrap(true)

	// 4. Barra de botones
	var buttons []fyne.CanvasObject

	if cancelText != "" {
		btnCancel := NewButton("Default", cancelText, "cancel", false, func() {
			if onCancel != nil {
				onCancel()
			}
			if d != nil {
				d.Hide()
			}
		})
		buttons = append(buttons, btnCancel)
	}

	confirmVariant := "Primary"
	confirmIcon := "check"
	if isDanger {
		confirmVariant = "Danger"
		confirmIcon = "delete"
	}

	btnConfirm := NewButton(confirmVariant, confirmText, confirmIcon, false, func() {
		if onConfirm != nil {
			onConfirm()
		}
		if d != nil {
			d.Hide()
		}
	})
	buttons = append(buttons, btnConfirm)

	buttonBar := container.NewHBox(buttons...)
	bottomRow := container.NewBorder(nil, nil, nil, buttonBar)

	// 5. Montar tarjeta estilizada
	contentBox := container.NewVBox(
		header,
		widget.NewSeparator(),
		container.NewPadded(msgLabel),
		bottomRow,
	)

	card := NewCard(contentBox)

	d = dialog.NewCustomWithoutButtons("", card, win)
	if d == nil {
		return nil
	}
	d.Resize(fyne.NewSize(440, 190))
	fyne.Do(func() {
		d.Show()
	})

	return d
}

// createAlertIconBadge crea una pastilla circular de 38x38 px con fondo pastel y su icono centrado.
func createAlertIconBadge(bgHex, borderHex string, iconColor fyne.ThemeColorName, iconRes fyne.Resource) fyne.CanvasObject {
	var iconHex string
	if IsDarkTheme() {
		switch iconColor {
		case theme.ColorNameError:
			bgHex, borderHex = "#450a0a", "#7f1d1d"
			iconHex = ThemeDangerHex()
		case theme.ColorNamePrimary:
			bgHex, borderHex = "#172554", "#1e3a8a"
			iconHex = ThemeAccentHex()
		case theme.ColorNameWarning:
			bgHex, borderHex = "#451a03", "#78350f"
			iconHex = ThemeWarningHex()
		default:
			bgHex, borderHex = "#172554", "#1e3a8a"
			iconHex = ThemeAccentHex()
		}
	} else {
		switch iconColor {
		case theme.ColorNameError:
			iconHex = ThemeDangerHex()
		case theme.ColorNamePrimary:
			iconHex = ThemeAccentHex()
		case theme.ColorNameWarning:
			iconHex = ThemeWarningHex()
		default:
			iconHex = ThemeAccentHex()
		}
	}
	return NewBadgeIcon(iconRes, 20, 38, bgHex, borderHex, iconHex)
}

// NewFormDialog construye un cuadro de diálogo modal moderno para formularios con botones UIKit
// (ModernButton con soporte de efecto cristal, sombra GPU y estados interactivos de alto contraste).
func NewFormDialog(
	title string,
	confirmText string,
	cancelText string,
	items []*widget.FormItem,
	onConfirm func(bool),
	win fyne.Window,
) dialog.Dialog {
	if win == nil {
		return nil
	}
	if confirmText == "" {
		confirmText = "Guardar"
	}
	if cancelText == "" {
		cancelText = "Cancelar"
	}

	var d dialog.Dialog
	called := false
	invokeCallback := func(ok bool) {
		if !called {
			called = true
			if onConfirm != nil {
				onConfirm(ok)
			}
		}
	}

	btnCancel := NewButton("Default", cancelText, "cancel", false, func() {
		invokeCallback(false)
		if d != nil {
			d.Hide()
		}
	})

	btnConfirm := NewButton("Primary", confirmText, "check", false, func() {
		invokeCallback(true)
		if d != nil {
			d.Hide()
		}
	})

	buttonBar := container.NewHBox(btnCancel, btnConfirm)
	bottomRow := container.NewBorder(nil, nil, nil, buttonBar)

	iconBadge := createAlertIconBadge("#eff6ff", "#bfdbfe", theme.ColorNamePrimary, theme.DocumentCreateIcon())
	titleLabel := NewSectionLabel(title)
	header := container.NewHBox(
		iconBadge,
		container.NewCenter(titleLabel),
	)

	form := widget.NewForm(items...)

	contentBox := container.NewVBox(
		header,
		widget.NewSeparator(),
		container.NewPadded(form),
		widget.NewSeparator(),
		bottomRow,
	)

	card := NewCard(contentBox).WithElevation(3)

	d = dialog.NewCustomWithoutButtons("", card, win)
	if d == nil {
		return nil
	}
	d.SetOnClosed(func() {
		invokeCallback(false)
	})

	initialHeight := float32(180 + len(items)*44)
	if initialHeight < 240 {
		initialHeight = 240
	}
	if initialHeight > 640 {
		initialHeight = 640
	}
	d.Resize(fyne.NewSize(500, initialHeight))

	return d
}

// ShowFormDialog construye y proyecta de inmediato un diálogo modal de formulario UIKit.
func ShowFormDialog(
	title string,
	confirmText string,
	cancelText string,
	items []*widget.FormItem,
	onConfirm func(bool),
	win fyne.Window,
) dialog.Dialog {
	d := NewFormDialog(title, confirmText, cancelText, items, onConfirm, win)
	if d != nil {
		fyne.Do(func() {
			d.Show()
		})
	}
	return d
}

