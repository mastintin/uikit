//go:build darwin

package uikit

import (
	"fmt"
	"os/exec"
	"strings"
)

// MsgDebug muestra una ventana de alerta nativa del sistema operativo (macOS)
// y BLOQUEA de forma síncrona la ejecución del código Go hasta que se pulse "Aceptar".
// Es ideal para depurar e inspeccionar valores sobre la marcha.
func MsgDebug(values ...any) {
	text := formatDebugMessage(values...)

	// Escapar comillas y barras para AppleScript
	escaped := strings.ReplaceAll(text, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `"`, `\"`)

	script := fmt.Sprintf(`display alert "MsgDebug" message "%s" as informational buttons {"Aceptar"} default button "Aceptar"`, escaped)
	cmd := exec.Command("osascript", "-e", script)
	_ = cmd.Run()
}
