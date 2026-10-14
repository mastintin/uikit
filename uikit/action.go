package uikit

import "fyne.io/fyne/v2"

// Action envuelve cualquier función de negocio para que se ejecute en una
// goroutine secundaria (Hilo 2 / Worker Thread).
// Es el enlace natural entre eventos de interfaz (botones, menús) y funciones
// procedurales secuenciales con diálogos síncronos (MsgYesNo, MsgInfo, etc.).
func Action(fn func()) func() {
	return func() {
		go fn()
	}
}

// Async ejecuta directamente una función en segundo plano (Hilo 2).
func Async(fn func()) {
	go fn()
}

// NewMenuItem crea un elemento de menú de Fyne cuya acción corre en Hilo 2
// de forma automática y transparente (estilo FiveWin).
func NewMenuItem(label string, action func()) *fyne.MenuItem {
	return fyne.NewMenuItem(label, func() {
		go action()
	})
}

// NewMenuItemSeparator crea un separador estándar de menú.
func NewMenuItemSeparator() *fyne.MenuItem {
	return fyne.NewMenuItemSeparator()
}

// Sync encola la ejecución de una función en el hilo principal de la interfaz gráfica (UI thread).
// Es el complemento seguro de Async/Action para modificar widgets o invocar Refresh().
func Sync(fn func()) {
	fyne.Do(fn)
}

// Do es un alias directo de fyne.Do dentro del ecosistema de uikit.
func Do(fn func()) {
	fyne.Do(fn)
}

