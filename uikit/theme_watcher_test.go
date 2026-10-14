package uikit

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func TestThemeWatcher_LifecycleAndDetection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watcher_test_*")
	if err != nil {
		t.Fatalf("Error creando tmpDir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	file1 := filepath.Join(tmpDir, "test.css")
	if err := os.WriteFile(file1, []byte(":root { --accent: #111; }"), 0644); err != nil {
		t.Fatalf("Error escribiendo archivo inicial: %v", err)
	}

	var callCount int32
	var lastChanged string

	watcher := NewThemeWatcher(tmpDir, func(fp string) {
		atomic.AddInt32(&callCount, 1)
		lastChanged = fp
	}).SetInterval(50 * time.Millisecond)

	if watcher.IsActive() {
		t.Errorf("Watcher no debe estar activo antes de Start()")
	}

	watcher.Start()
	if !watcher.IsActive() {
		t.Errorf("Watcher debe estar activo tras Start()")
	}

	// Al iniciar no debe disparar eventos
	time.Sleep(100 * time.Millisecond)
	if atomic.LoadInt32(&callCount) != 0 {
		t.Errorf("No debe haber llamadas iniciales sin cambios, obtenido %d", atomic.LoadInt32(&callCount))
	}

	// Modificar archivo
	time.Sleep(20 * time.Millisecond)
	_ = os.WriteFile(file1, []byte(":root { --accent: #222; }"), 0644)

	// Esperar que el watcher lo detecte
	time.Sleep(150 * time.Millisecond)
	if atomic.LoadInt32(&callCount) == 0 {
		t.Errorf("Se esperaba al menos 1 llamada tras modificar el archivo")
	}
	if filepath.Base(lastChanged) != "test.css" {
		t.Errorf("Esperado test.css, obtenido %s", filepath.Base(lastChanged))
	}

	// Detener watcher
	watcher.Stop()
	if watcher.IsActive() {
		t.Errorf("Watcher no debe estar activo tras Stop()")
	}

	// Modificar tras Stop no debe disparar llamadas adicionales
	countBefore := atomic.LoadInt32(&callCount)
	_ = os.WriteFile(file1, []byte(":root { --accent: #333; }"), 0644)
	time.Sleep(120 * time.Millisecond)

	if atomic.LoadInt32(&callCount) != countBefore {
		t.Errorf("No debe haber llamadas tras Stop(), antes %d ahora %d", countBefore, atomic.LoadInt32(&callCount))
	}
}
