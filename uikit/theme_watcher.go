package uikit

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ThemeWatcher supervisa una carpeta en busca de modificaciones en archivos .css
// de forma ligera (inspección de ModTime a intervalos regulares), 100% multiplataforma,
// sin dependencias de terceros y con consumo prácticamente nulo de recursos.
type ThemeWatcher struct {
	dir      string
	interval time.Duration
	onChange func(filePath string)

	mu       sync.Mutex
	active   bool
	stopCh   chan struct{}
	lastMods map[string]time.Time
}

// NewThemeWatcher crea una instancia del vigía de temas.
func NewThemeWatcher(dir string, onChange func(filePath string)) *ThemeWatcher {
	return &ThemeWatcher{
		dir:      dir,
		interval: 1000 * time.Millisecond,
		onChange: onChange,
		lastMods: make(map[string]time.Time),
	}
}

// SetInterval permite ajustar el intervalo de sondeo (por defecto 1s).
func (w *ThemeWatcher) SetInterval(d time.Duration) *ThemeWatcher {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.interval = d
	return w
}

// SetDir actualiza el directorio vigilado y reinicia el mapa de marcas de tiempo.
func (w *ThemeWatcher) SetDir(dir string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.dir = dir
	w.lastMods = make(map[string]time.Time)
}

// IsActive indica si el vigía está actualmente en funcionamiento.
func (w *ThemeWatcher) IsActive() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.active
}

// Start inicia la vigilancia en segundo plano.
func (w *ThemeWatcher) Start() {
	w.mu.Lock()
	if w.active {
		w.mu.Unlock()
		return
	}
	w.active = true
	w.stopCh = make(chan struct{})

	// Registrar marcas de tiempo actuales para no disparar eventos falsos al arrancar
	w.scanInitialMods()
	w.mu.Unlock()

	go w.watchLoop()
}

func (w *ThemeWatcher) scanInitialMods() {
	if w.dir == "" {
		return
	}
	entries, err := os.ReadDir(w.dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() || strings.ToLower(filepath.Ext(entry.Name())) != ".css" {
			continue
		}
		if info, err := entry.Info(); err == nil {
			w.lastMods[entry.Name()] = info.ModTime()
		}
	}
}

// Stop detiene el vigía de forma inmediata y libera la goroutine.
func (w *ThemeWatcher) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.active {
		return
	}
	w.active = false
	if w.stopCh != nil {
		close(w.stopCh)
		w.stopCh = nil
	}
}

func (w *ThemeWatcher) watchLoop() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-w.stopCh:
			return
		case <-ticker.C:
			w.checkChanges()
		}
	}
}

func (w *ThemeWatcher) checkChanges() {
	w.mu.Lock()
	if !w.active || w.dir == "" {
		w.mu.Unlock()
		return
	}
	dir := w.dir
	callback := w.onChange
	w.mu.Unlock()

	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	var changedFile string

	w.mu.Lock()
	currentFiles := make(map[string]bool)
	for _, entry := range entries {
		if entry.IsDir() || strings.ToLower(filepath.Ext(entry.Name())) != ".css" {
			continue
		}
		name := entry.Name()
		currentFiles[name] = true

		info, err := entry.Info()
		if err != nil {
			continue
		}

		lastMod, exists := w.lastMods[name]
		if !exists || info.ModTime().After(lastMod) {
			w.lastMods[name] = info.ModTime()
			changedFile = filepath.Join(dir, name)
		}
	}

	// Limpiar del registro archivos que hayan sido eliminados
	for recorded := range w.lastMods {
		if !currentFiles[recorded] {
			delete(w.lastMods, recorded)
		}
	}
	w.mu.Unlock()

	if changedFile != "" && callback != nil {
		callback(changedFile)
	}
}
