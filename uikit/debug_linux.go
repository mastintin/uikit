//go:build linux

package uikit

import (
	"fmt"
	"os/exec"
)

// MsgDebug muestra un cuadro de diálogo nativo en Linux (usando zenity o kdialog)
// y BLOQUEA de forma síncrona la ejecución del código Go hasta que se pulse "Aceptar".
func MsgDebug(values ...any) {
	text := formatDebugMessage(values...)
	title := "MsgDebug"

	if path, err := exec.LookPath("zenity"); err == nil {
		cmd := exec.Command(path, "--info", "--title="+title, "--text="+text, "--ok-label=Aceptar")
		_ = cmd.Run()
		return
	}

	if path, err := exec.LookPath("kdialog"); err == nil {
		cmd := exec.Command(path, "--msgbox", text, "--title", title)
		_ = cmd.Run()
		return
	}

	fmt.Printf("\n[MsgDebug] %s\nPresione ENTER para continuar...", text)
	var dummy string
	fmt.Scanln(&dummy)
}
