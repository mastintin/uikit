package uikit

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Dimensiones intermedias, cómodas y perfectamente equilibradas para el diálogo de archivos
const (
	ModernFileDialogDefaultWidth  float32 = 820
	ModernFileDialogDefaultHeight float32 = 510
)

// calculateDialogSize determina un tamaño espacioso y proporcional al tamaño actual de la ventana
func calculateDialogSize(win fyne.Window) fyne.Size {
	targetW := ModernFileDialogDefaultWidth
	targetH := ModernFileDialogDefaultHeight

	if win != nil && win.Canvas() != nil {
		winSize := win.Canvas().Size()
		if winSize.Width > 200 {
			w := winSize.Width * 0.85
			if w < targetW {
				targetW = w
			}
			if targetW > winSize.Width-30 {
				targetW = winSize.Width - 30
			}
		}
		if winSize.Height > 200 {
			h := winSize.Height * 0.85
			if h < targetH {
				targetH = h
			}
			if targetH > winSize.Height-30 {
				targetH = winSize.Height - 30
			}
		}
	}

	return fyne.NewSize(targetW, targetH)
}

type fileEntryInfo struct {
	name    string
	path    string
	isDir   bool
	size    int64
	modTime time.Time
}

func formatBytes(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	units := []string{"KB", "MB", "GB", "TB"}
	val := float64(bytes)
	for _, u := range units {
		val /= 1024
		if val < 1024 {
			return fmt.Sprintf("%.1f %s", val, u)
		}
	}
	return fmt.Sprintf("%.1f TB", val)
}

type locationShortcut struct {
	name string
	path string
	icon fyne.Resource
}

type fileDialogKind int

const (
	dialogKindSave fileDialogKind = iota
	dialogKindOpenFile
	dialogKindOpenFolder
)

// ShowSaveFileDialog proyecta un diálogo de archivos 100% uikit amplio y moderno para guardar archivos.
// Utiliza exclusivamente los botones y estética del sistema uikit (ModernButton con azul translúcido suave).
func ShowSaveFileDialog(
	win fyne.Window,
	title string,
	defaultFileName string,
	ext string, // ej: "csv", "pdf"
	onSelected func(filePath string),
) {
	showModernFileDialog(win, dialogKindSave, title, defaultFileName, ext, nil, "", onSelected)
}

// ShowOpenFileDialog proyecta un diálogo de archivos 100% uikit amplio y moderno para seleccionar/abrir archivos.
// Utiliza exclusivamente los botones y estética del sistema uikit (ModernButton con azul translúcido suave).
func ShowOpenFileDialog(
	win fyne.Window,
	title string,
	extensions []string,
	onSelected func(filePath string),
) {
	showModernFileDialog(win, dialogKindOpenFile, title, "", "", extensions, "", onSelected)
}

// ShowOpenFileDialogWithInitial proyecta un diálogo de archivos con ruta inicial preferente.
func ShowOpenFileDialogWithInitial(
	win fyne.Window,
	title string,
	extensions []string,
	initialPath string,
	onSelected func(filePath string),
) {
	showModernFileDialog(win, dialogKindOpenFile, title, "", "", extensions, initialPath, onSelected)
}

// ShowFolderDialog proyecta un diálogo 100% uikit amplio y moderno para seleccionar carpetas/directorios.
// Utiliza exclusivamente los botones y estética del sistema uikit (ModernButton con azul translúcido suave).
func ShowFolderDialog(
	win fyne.Window,
	title string,
	onSelected func(folderPath string),
) {
	showModernFileDialog(win, dialogKindOpenFolder, title, "", "", nil, "", onSelected)
}

// ShowFolderDialogWithInitial proyecta un diálogo de selección de carpetas con ruta inicial preferente.
func ShowFolderDialogWithInitial(
	win fyne.Window,
	title string,
	initialPath string,
	onSelected func(folderPath string),
) {
	showModernFileDialog(win, dialogKindOpenFolder, title, "", "", nil, initialPath, onSelected)
}

func showModernFileDialog(
	win fyne.Window,
	kind fileDialogKind,
	title string,
	defaultFileName string,
	ext string,
	openExts []string,
	initialPath string,
	onSelected func(filePath string),
) {
	if onSelected == nil {
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil || homeDir == "" {
		homeDir = "."
	}

	// 1. Directorio inicial preferente
	initialDir := homeDir
	if initialPath != "" {
		cleanInit := initialPath
		if !filepath.IsAbs(cleanInit) {
			if abs, err := filepath.Abs(cleanInit); err == nil {
				cleanInit = abs
			}
		}
		if info, err := os.Stat(cleanInit); err == nil {
			if info.IsDir() {
				initialDir = cleanInit
			} else {
				initialDir = filepath.Dir(cleanInit)
			}
		} else {
			dirPart := filepath.Dir(cleanInit)
			if infoDir, err := os.Stat(dirPart); err == nil && infoDir.IsDir() {
				initialDir = dirPart
			}
		}
	} else if kind == dialogKindSave {
		docs := filepath.Join(homeDir, "Documents")
		if _, err := os.Stat(docs); err == nil {
			initialDir = docs
		}
	} else {
		pics := filepath.Join(homeDir, "Pictures")
		if _, err := os.Stat(pics); err == nil {
			initialDir = pics
		} else {
			docs := filepath.Join(homeDir, "Documents")
			if _, err := os.Stat(docs); err == nil {
				initialDir = docs
			}
		}
	}

	// Extensiones permitidas
	var allowedExtensions []string
	cleanExt := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(ext)), ".")
	if kind == dialogKindSave {
		if cleanExt != "" {
			allowedExtensions = []string{"." + cleanExt}
			if !strings.HasSuffix(strings.ToLower(defaultFileName), "."+cleanExt) {
				defaultFileName += "." + cleanExt
			}
		}
	} else if kind == dialogKindOpenFile {
		for _, e := range openExts {
			clean := strings.ToLower(strings.TrimSpace(e))
			if !strings.HasPrefix(clean, ".") {
				clean = "." + clean
			}
			allowedExtensions = append(allowedExtensions, clean)
		}
	}

	loadEntries := func(dir string) []fileEntryInfo {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil
		}

		var dirs []fileEntryInfo
		var files []fileEntryInfo

		for _, entry := range entries {
			name := entry.Name()
			if strings.HasPrefix(name, ".") {
				continue // Ocultar archivos ocultos
			}

			info, err := entry.Info()
			if err != nil {
				continue
			}

			fullPath := filepath.Join(dir, name)
			isAppBundle := strings.HasSuffix(strings.ToLower(name), ".app")

			if entry.IsDir() && !isAppBundle {
				dirs = append(dirs, fileEntryInfo{
					name:    name,
					path:    fullPath,
					isDir:   true,
					modTime: info.ModTime(),
				})
			} else if kind != dialogKindOpenFolder {
				if len(allowedExtensions) > 0 {
					fileExt := strings.ToLower(filepath.Ext(name))
					matched := false
					for _, aExt := range allowedExtensions {
						if fileExt == aExt {
							matched = true
							break
						}
					}
					if !matched {
						continue
					}
				}

				var size int64
				if !entry.IsDir() {
					size = info.Size()
				}

				files = append(files, fileEntryInfo{
					name:    name,
					path:    fullPath,
					isDir:   false,
					size:    size,
					modTime: info.ModTime(),
				})
			}
		}

		sort.Slice(dirs, func(i, j int) bool {
			return strings.ToLower(dirs[i].name) < strings.ToLower(dirs[j].name)
		})
		sort.Slice(files, func(i, j int) bool {
			return strings.ToLower(files[i].name) < strings.ToLower(files[j].name)
		})

		return append(dirs, files...)
	}

	currentDir := initialDir
	currentItems := loadEntries(currentDir)
	var selectedFilePath string

	if title == "" {
		switch kind {
		case dialogKindSave:
			title = "Guardar Archivo"
		case dialogKindOpenFolder:
			title = "Seleccionar Carpeta"
		default:
			title = "Seleccionar Archivo"
		}
	}

	var customDialog dialog.Dialog
	var fileList *widget.List
	var nameEntry *widget.Entry
	var selectedFileLabel *widget.Label
	var selectedFolderLabel *widget.Label
	var btnConfirm *ModernButton
	var btnUp *ModernButton
	var sidebarList *ModernListBox[locationShortcut]

	var shortcuts []locationShortcut
	addShortcut := func(name, path string, icon fyne.Resource) {
		if _, err := os.Stat(path); err == nil {
			shortcuts = append(shortcuts, locationShortcut{
				name: name,
				path: path,
				icon: icon,
			})
		}
	}

	if runtime.GOOS == "darwin" {
		addShortcut("Aplicaciones", "/Applications", theme.StorageIcon())
	}
	addShortcut("Documentos", filepath.Join(homeDir, "Documents"), theme.DocumentIcon())
	addShortcut("Imágenes", filepath.Join(homeDir, "Pictures"), theme.FileImageIcon())
	addShortcut("Escritorio", filepath.Join(homeDir, "Desktop"), theme.DesktopIcon())
	addShortcut("Descargas", filepath.Join(homeDir, "Downloads"), theme.DownloadIcon())
	addShortcut("Inicio", homeDir, theme.HomeIcon())

	pathLabel := widget.NewLabel(currentDir)
	pathLabel.TextStyle = fyne.TextStyle{Bold: true}
	pathLabel.Truncation = fyne.TextTruncateClip

	updateDirectory := func(newDir string) {
		cleanNew := filepath.Clean(newDir)
		info, err := os.Stat(cleanNew)
		if err != nil || !info.IsDir() {
			return
		}

		currentDir = cleanNew
		currentItems = loadEntries(currentDir)
		pathLabel.SetText(currentDir)
		if selectedFolderLabel != nil {
			selectedFolderLabel.SetText(currentDir)
		}

		parentDir := filepath.Dir(currentDir)
		if btnUp != nil {
			if parentDir == currentDir {
				btnUp.Disable()
			} else {
				btnUp.Enable()
			}
		}

		if sidebarList != nil {
			found := -1
			for i, sc := range shortcuts {
				if sc.path == currentDir {
					found = i
					break
				}
			}
			if found >= 0 {
				if sidebarList.selectedIndex != found {
					sidebarList.Select(found)
				}
			} else {
				sidebarList.UnselectAll()
			}
		}

		if fileList != nil {
			fileList.UnselectAll()
			fileList.Refresh()
		}
	}

	doConfirm := func() {
		var finalPath string
		switch kind {
		case dialogKindSave:
			fileName := ""
			if nameEntry != nil {
				fileName = strings.TrimSpace(nameEntry.Text)
			}
			if fileName == "" {
				MsgError(win, "Por favor introduzca un nombre para el archivo", "Validación")
				return
			}
			if cleanExt != "" && !strings.HasSuffix(strings.ToLower(fileName), "."+cleanExt) {
				fileName += "." + cleanExt
			}
			finalPath = filepath.Join(currentDir, fileName)
		case dialogKindOpenFolder:
			finalPath = currentDir
		case dialogKindOpenFile:
			if selectedFilePath == "" {
				MsgError(win, "Por favor seleccione un archivo de la lista", "Validación")
				return
			}
			finalPath = selectedFilePath
		}

		if customDialog != nil {
			customDialog.Hide()
		}
		if onSelected != nil {
			onSelected(finalPath)
		}
	}

	// 2. Botones de acción principales uikit (ModernButton con estilo Primary translúcido)
	confirmLabel := "Guardar"
	confirmIcon := "save"
	confirmDisabled := false

	switch kind {
	case dialogKindSave:
		confirmLabel = "Guardar"
		confirmIcon = "save"
		confirmDisabled = false
	case dialogKindOpenFolder:
		confirmLabel = "Seleccionar Carpeta"
		confirmIcon = "check"
		confirmDisabled = false
	case dialogKindOpenFile:
		confirmLabel = "Seleccionar"
		confirmIcon = "check"
		confirmDisabled = true
	}

	btnConfirm = NewButton("Primary", confirmLabel, confirmIcon, confirmDisabled, func() {
		doConfirm()
	})

	btnCancel := NewButton("Default", "Cancelar", "cancel", false, func() {
		if customDialog != nil {
			customDialog.Hide()
		}
	})

	btnUp = NewButton("Default", "Subir", "up", false, func() {
		parent := filepath.Dir(currentDir)
		if parent != currentDir {
			updateDirectory(parent)
		}
	})

	// 3. Lista de archivos y carpetas
	var lastSelectedID = -1
	var lastSelectedTime time.Time

	fileList = widget.NewList(
		func() int {
			return len(currentItems)
		},
		func() fyne.CanvasObject {
			icon := widget.NewIcon(theme.FolderIcon())
			name := widget.NewLabel("sample_file_or_directory_name")
			details := widget.NewLabel("124 KB  •  01/01/2026")
			details.TextStyle = fyne.TextStyle{Italic: true}
			details.Alignment = fyne.TextAlignTrailing

			return container.NewBorder(
				nil, nil,
				container.NewHBox(icon, name),
				details,
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id < 0 || id >= len(currentItems) {
				return
			}
			entry := currentItems[id]
			row := item.(*fyne.Container)
			leftBox := row.Objects[0].(*fyne.Container)
			iconWidget := leftBox.Objects[0].(*widget.Icon)
			nameLabel := leftBox.Objects[1].(*widget.Label)
			detailsLabel := row.Objects[1].(*widget.Label)

			nameLabel.SetText(entry.name)
			if entry.isDir {
				iconWidget.SetResource(theme.FolderIcon())
				nameLabel.TextStyle = fyne.TextStyle{Bold: true}
				detailsLabel.SetText("Carpeta  •  " + entry.modTime.Format("02/01/2006"))
			} else {
				fExt := strings.ToLower(filepath.Ext(entry.name))
				if fExt == ".jpg" || fExt == ".jpeg" || fExt == ".png" || fExt == ".psd" {
					iconWidget.SetResource(theme.FileImageIcon())
				} else if fExt == ".csv" || fExt == ".txt" {
					iconWidget.SetResource(theme.FileTextIcon())
				} else if fExt == ".db" || fExt == ".sqlite" || fExt == ".sqlite3" {
					iconWidget.SetResource(theme.StorageIcon())
				} else if fExt == ".app" || fExt == ".exe" {
					iconWidget.SetResource(theme.ColorPaletteIcon())
				} else {
					iconWidget.SetResource(theme.DocumentIcon())
				}
				nameLabel.TextStyle = fyne.TextStyle{Bold: false}
				if fExt == ".app" {
					detailsLabel.SetText("Aplicación  •  " + entry.modTime.Format("02/01/2006"))
				} else {
					detailsLabel.SetText(fmt.Sprintf("%s  •  %s", formatBytes(entry.size), entry.modTime.Format("02/01/2006 15:04")))
				}
			}
			nameLabel.Refresh()
			detailsLabel.Refresh()
		},
	)

	fileList.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(currentItems) {
			return
		}
		entry := currentItems[id]
		now := time.Now()
		isDouble := (lastSelectedID == id && now.Sub(lastSelectedTime) < 450*time.Millisecond)
		lastSelectedID = id
		lastSelectedTime = now

		if entry.isDir {
			updateDirectory(entry.path)
			return
		}

		// Archivo seleccionado
		selectedFilePath = entry.path
		if kind == dialogKindSave {
			if nameEntry != nil {
				nameEntry.SetText(entry.name)
			}
		} else if kind == dialogKindOpenFile {
			if selectedFileLabel != nil {
				selectedFileLabel.SetText(entry.name)
			}
			if btnConfirm != nil {
				btnConfirm.Enable()
			}
		}

		if isDouble {
			doConfirm()
		}
	}

	// 4. Barra lateral con Accesos Rápidos usando ModernSidebar y ModernListBox
	sidebarList = NewListBox[locationShortcut](
		shortcuts,
		func() fyne.CanvasObject {
			icon := widget.NewIcon(theme.FolderIcon())
			label := widget.NewLabelWithStyle("Ubicación", fyne.TextAlignLeading, fyne.TextStyle{})
			return container.NewHBox(icon, label)
		},
		func(item locationShortcut, obj fyne.CanvasObject) {
			box := obj.(*fyne.Container)
			box.Objects[0].(*widget.Icon).SetResource(item.icon)
			box.Objects[1].(*widget.Label).SetText(item.name)
		},
		func(item locationShortcut) {
			if item.path != currentDir {
				updateDirectory(item.path)
			}
		},
	)

	// Seleccionar la carpeta inicial en la lista si coincide
	for i, sc := range shortcuts {
		if sc.path == currentDir {
			sidebarList.Select(i)
			break
		}
	}

	lblShortcuts := container.NewPadded(widget.NewLabelWithStyle("UBICACIONES", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	sidebar := NewSidebarContainer(170, lblShortcuts, sidebarList, nil).WithShadedBackground()

	// 5. Cabecera con título, ruta actual y botón de subir nivel / nueva carpeta
	var modeBadgeIcon fyne.Resource
	switch kind {
	case dialogKindSave:
		modeBadgeIcon = theme.DocumentSaveIcon()
	case dialogKindOpenFolder:
		modeBadgeIcon = theme.FolderOpenIcon()
	default:
		modeBadgeIcon = theme.DocumentIcon()
	}
	iconBadge := createAlertIconBadge("#eff6ff", "#bfdbfe", theme.ColorNamePrimary, modeBadgeIcon)
	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	titleBox := container.NewHBox(iconBadge, container.NewCenter(titleLabel))

	btnNewFolder := NewButton("Default", "Nueva Carpeta", "new-folder", false, func() {
		folderEntry := widget.NewEntry()
		folderEntry.SetPlaceHolder("Nombre de la nueva carpeta")

		var promptDialog dialog.Dialog
		btnCreateFolder := NewButton("Primary", "Crear Carpeta", "check", false, func() {
			newFolderName := strings.TrimSpace(folderEntry.Text)
			if newFolderName == "" {
				return
			}
			newPath := filepath.Join(currentDir, newFolderName)
			if err := os.MkdirAll(newPath, 0755); err != nil {
				MsgError(win, fmt.Errorf("error al crear carpeta: %w", err), "Error")
				return
			}
			if promptDialog != nil {
				promptDialog.Hide()
			}
			updateDirectory(newPath)
		})
		btnCancelFolder := NewButton("Default", "Cancelar", "cancel", false, func() {
			if promptDialog != nil {
				promptDialog.Hide()
			}
		})

		promptContent := container.NewVBox(
			widget.NewLabel("Introduzca el nombre de la carpeta:"),
			folderEntry,
			widget.NewSeparator(),
			container.NewBorder(nil, nil, nil, container.NewHBox(btnCancelFolder, btnCreateFolder)),
		)
		promptCard := NewPanel(promptContent, "", "", 12, 1)
		promptDialog = dialog.NewCustomWithoutButtons("Nueva Carpeta", promptCard, win)
		promptDialog.Resize(fyne.NewSize(380, 160))
		promptDialog.Show()
	})

	headerTop := container.NewBorder(
		nil, nil,
		titleBox,
		container.NewHBox(btnUp, btnNewFolder),
	)

	pathBar := container.NewHBox(
		widget.NewIcon(theme.FolderIcon()),
		pathLabel,
	)

	header := container.NewVBox(
		headerTop,
		pathBar,
		widget.NewSeparator(),
	)

	// 6. Pie de diálogo con input de nombre y botones Guardar/Seleccionar y Cancelar
	var footerContent fyne.CanvasObject
	actionsBar := container.NewHBox(btnCancel, btnConfirm)

	if kind == dialogKindSave {
		nameEntry = widget.NewEntry()
		nameEntry.SetText(defaultFileName)
		nameEntry.OnSubmitted = func(string) {
			doConfirm()
		}

		lbl := widget.NewLabelWithStyle("Nombre de archivo:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		footerContent = container.NewVBox(
			widget.NewSeparator(),
			container.NewBorder(nil, nil, lbl, actionsBar, container.NewPadded(nameEntry)),
		)
	} else if kind == dialogKindOpenFolder {
		selectedFolderLabel = widget.NewLabel(currentDir)
		selectedFolderLabel.TextStyle = fyne.TextStyle{Bold: true}
		selectedFolderLabel.Truncation = fyne.TextTruncateClip

		lbl := widget.NewLabelWithStyle("Carpeta seleccionada:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		footerContent = container.NewVBox(
			widget.NewSeparator(),
			container.NewBorder(nil, nil, lbl, actionsBar, container.NewPadded(selectedFolderLabel)),
		)
	} else {
		selectedFileLabel = widget.NewLabel("Haga doble clic en un archivo o selecciónelo y pulse Seleccionar")
		selectedFileLabel.TextStyle = fyne.TextStyle{Italic: true}
		selectedFileLabel.Truncation = fyne.TextTruncateClip

		lbl := widget.NewLabelWithStyle("Archivo seleccionado:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		footerContent = container.NewVBox(
			widget.NewSeparator(),
			container.NewBorder(nil, nil, lbl, actionsBar, container.NewPadded(selectedFileLabel)),
		)
	}

	// 7. Estructura y proyección del diálogo
	centerBody := container.NewBorder(
		nil, nil,
		sidebar,
		nil,
		container.NewPadded(fileList),
	)

	mainLayout := container.NewBorder(
		header,
		footerContent,
		nil, nil,
		centerBody,
	)

	card := NewPanel(mainLayout, "", "", 12, 1)

	if win == nil {
		if kind == dialogKindSave && defaultFileName != "" {
			onSelected(filepath.Join(currentDir, defaultFileName))
		} else if kind == dialogKindOpenFolder {
			onSelected(currentDir)
		}
		return
	}

	customDialog = dialog.NewCustomWithoutButtons("", card, win)
	if customDialog == nil {
		return
	}
	customDialog.Resize(calculateDialogSize(win))
	customDialog.Show()
}

// ShowSavedSuccessModal muestra una tarjeta modal moderna de confirmación tras guardar un archivo,
// permitiendo al usuario ver el resumen de la operación y abrir la carpeta contenedora con 1 clic.
func ShowSavedSuccessModal(
	win fyne.Window,
	title string,
	message string,
	savedFilePath string,
	onRevealInFinder func(path string),
) {
	if win == nil {
		return
	}

	var d dialog.Dialog

	btnAccept := NewButton("Primary", "Aceptar", "check", false, func() {
		if d != nil {
			d.Hide()
		}
	})

	var buttons []fyne.CanvasObject
	if savedFilePath != "" && onRevealInFinder != nil {
		revealLabel := "Mostrar en Finder"
		if runtime.GOOS == "windows" {
			revealLabel = "Mostrar en Explorador"
		} else if runtime.GOOS != "darwin" {
			revealLabel = "Abrir Carpeta"
		}

		btnReveal := NewButton("Default", revealLabel, "folder", false, func() {
			if d != nil {
				d.Hide()
			}
			onRevealInFinder(savedFilePath)
		})
		buttons = append(buttons, btnReveal)
	}
	buttons = append(buttons, btnAccept)

	buttonBar := container.NewHBox(buttons...)

	iconBadge := createAlertIconBadge("#eff6ff", "#bfdbfe", theme.ColorNamePrimary, theme.DocumentSaveIcon())
	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	header := container.NewHBox(iconBadge, container.NewCenter(titleLabel))

	msgLabel := widget.NewLabel(message)
	msgLabel.Wrapping = fyne.TextWrapWord

	pathCard := container.NewPadded(container.NewVBox(
		widget.NewLabelWithStyle("Ruta:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(savedFilePath),
	))

	contentBox := container.NewVBox(
		header,
		widget.NewSeparator(),
		msgLabel,
		pathCard,
		widget.NewSeparator(),
		container.NewBorder(nil, nil, nil, buttonBar),
	)

	card := NewPanel(contentBox, "", "", 12, 1)
	d = dialog.NewCustomWithoutButtons("", card, win)
	if d == nil {
		return
	}
	d.Resize(fyne.NewSize(500, 230))
	d.Show()
}
