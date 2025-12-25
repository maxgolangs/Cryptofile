package ui

import (
	"fmt"
	"os"

	"cryptor/pkg/batch"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowDecryptBatchWindow(parent fyne.Window) {
	window := fyne.CurrentApp().NewWindow("Пакетная расшифровка")
	window.Resize(fyne.NewSize(600, 350))
	window.CenterOnScreen()

	var selectedDir string
	var password string
	var removeOriginal bool

	dirLabel := widget.NewLabel("Директория не выбрана\n💡 Или перетащите папку сюда")
	dirLabel.Wrapping = fyne.TextWrapWord

	handleDrop := func(uri fyne.URI) {
		if uri == nil {
			return
		}
		selectedDir = uri.Path()
		if info, err := os.Stat(selectedDir); err == nil && info.IsDir() {
			dirLabel.SetText(fmt.Sprintf("Выбрано: %s", selectedDir))
		} else {
			dirLabel.SetText("❌ Выберите директорию, а не файл")
			selectedDir = ""
		}
	}

	selectBtn := widget.NewButton("Выбрать директорию", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			selectedDir = uri.Path()
			dirLabel.SetText(fmt.Sprintf("Выбрано: %s", selectedDir))
		}, window).Show()
	})

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Введите пароль")

	removeCheck := widget.NewCheck("Удалить зашифрованные файлы после расшифровки", func(checked bool) {
		removeOriginal = checked
	})
	removeCheck.SetChecked(true)

	statusLabel := widget.NewLabel("")
	statusLabel.Wrapping = fyne.TextWrapWord

	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	var decryptBtn *widget.Button
	decryptBtn = widget.NewButton("Расшифровать все файлы", func() {
		if selectedDir == "" {
			statusLabel.SetText("❌ Пожалуйста, выберите директорию")
			return
		}

		password = passwordEntry.Text
		if password == "" {
			statusLabel.SetText("❌ Введите пароль")
			return
		}

		if _, err := os.Stat(selectedDir); err != nil {
			statusLabel.SetText(fmt.Sprintf("❌ Ошибка: %v", err))
			return
		}

		decryptBtn.Disable()
		progressBar.Show()
		statusLabel.SetText("⏳ Выполняется пакетная расшифровка...")

		go func() {
			processed, errors, total := batch.ProcessDirectoryParallel(selectedDir, password, "decrypt", removeOriginal)

			progressBar.Hide()
			if total == 0 {
				statusLabel.SetText("⚠️ Подходящих файлов для расшифровки не найдено")
			} else {
				statusLabel.SetText(fmt.Sprintf("✅ Обработано: %d из %d файлов. Ошибок: %d", processed, total, errors))
			}
			decryptBtn.Enable()
		}()
	})

	window.SetOnDropped(func(pos fyne.Position, uris []fyne.URI) {
		if len(uris) > 0 {
			handleDrop(uris[0])
		}
	})

	content := container.NewVBox(
		widget.NewCard("", "Выбор директории", container.NewVBox(
			dirLabel,
			selectBtn,
		)),
		widget.NewCard("", "Пароль", container.NewVBox(
			passwordEntry,
		)),
		widget.NewCard("", "Параметры", container.NewVBox(
			removeCheck,
		)),
		decryptBtn,
		progressBar,
		statusLabel,
	)

	window.SetContent(container.NewPadded(content))
	window.Show()
}

