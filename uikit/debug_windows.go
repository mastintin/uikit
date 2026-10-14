//go:build windows

package uikit

import (
	"syscall"
	"unsafe"
)

var (
	user32          = syscall.NewLazyDLL("user32.dll")
	procMessageBoxW = user32.NewProc("MessageBoxW")
)

const (
	mbOK              = 0x00000000
	mbIconInformation = 0x00000040
)

// MsgDebug muestra un cuadro de diálogo nativo de Windows (MessageBoxW)
// y BLOQUEA de forma síncrona la ejecución del código Go hasta que se pulse "Aceptar".
func MsgDebug(values ...any) {
	text := formatDebugMessage(values...)
	title := "MsgDebug"

	lpText, _ := syscall.UTF16PtrFromString(text)
	lpCaption, _ := syscall.UTF16PtrFromString(title)

	procMessageBoxW.Call(
		0,
		uintptr(unsafe.Pointer(lpText)),
		uintptr(unsafe.Pointer(lpCaption)),
		uintptr(mbOK|mbIconInformation),
	)
}
