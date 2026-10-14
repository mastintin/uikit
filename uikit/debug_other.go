//go:build !darwin && !windows && !linux

package uikit

import (
	"fmt"
)

// MsgDebug fallback para plataformas no contempladas directamente.
func MsgDebug(values ...any) {
	text := formatDebugMessage(values...)
	fmt.Printf("\n[MsgDebug] %s\nPresione ENTER para continuar...", text)
	var dummy string
	fmt.Scanln(&dummy)
}
