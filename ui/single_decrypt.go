package ui

import (
	"fmt"
	"os"

	"cryptor/pkg/decrypt"
	"cryptor/internal/file"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowDecryptSingleWindow(parent fyne.Window) {
	window := fyne.CurrentApp().NewWindow("Расшифровка файла")
	window.Resize(fyne.NewSize(600, 350))
	window.CenterOnScreen()

	var selectedPath string
	var password string
	var removeOriginal bool

	pathLabel := widget.NewLabel("Файл не выбран\n💡 Или перетащите файл сюда")
	pathLabel.Wrapping = fyne.TextWrapWord

	handleDrop := func(uri fyne.URI) {
		if uri == nil {
			return
		}
		selectedPath = uri.Path()
		if file.IsEncryptedCandidate(selectedPath) {
			pathLabel.SetText(fmt.Sprintf("Выбрано: %s", selectedPath))
		} else {
			pathLabel.SetText(fmt.Sprintf("❌ Файл не является зашифрованным файлом CryptoFile: %s", selectedPath))
			selectedPath = ""
		}
	}

	selectBtn := widget.NewButton("Выбрать зашифрованный файл (.encrypted)", func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()
			
			uri := reader.URI()
			if uri == nil {
				return
			}
			selectedPath = uri.Path()
			
			if file.IsEncryptedCandidate(selectedPath) {
				pathLabel.SetText(fmt.Sprintf("Выбрано: %s", selectedPath))
			} else {
				pathLabel.SetText(fmt.Sprintf("❌ Файл не является зашифрованным файлом CryptoFile: %s", selectedPath))
				selectedPath = ""
			}
		}, window)
		fileDialog.Show()
	})

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Введите пароль")

	removeCheck := widget.NewCheck("Удалить зашифрованный файл после расшифровки", func(checked bool) {
		removeOriginal = checked
	})
	removeCheck.SetChecked(true)

	statusLabel := widget.NewLabel("")
	statusLabel.Wrapping = fyne.TextWrapWord

	var decryptBtn *widget.Button
	decryptBtn = widget.NewButton("Расшифровать", func() {
		if selectedPath == "" {
			statusLabel.SetText("❌ Пожалуйста, выберите зашифрованный файл")
			return
		}

		if !file.IsEncryptedCandidate(selectedPath) {
			statusLabel.SetText("❌ Выбранный файл не является зашифрованным файлом CryptoFile")
			return
		}

		password = passwordEntry.Text
		if password == "" {
			statusLabel.SetText("❌ Введите пароль")
			return
		}

		if _, err := os.Stat(selectedPath); err != nil {
			statusLabel.SetText(fmt.Sprintf("❌ Ошибка: %v", err))
			return
		}

		decryptBtn.Disable()
		statusLabel.SetText("⏳ Выполняется расшифровка...")

		go func() {
			result, err := decrypt.DecryptFile(selectedPath, password, removeOriginal)
			if err != nil {
				statusLabel.SetText(fmt.Sprintf("❌ Ошибка: %v", err))
				decryptBtn.Enable()
				return
			}

			statusLabel.SetText(fmt.Sprintf("✅ Успешно! Расшифрованный файл: %s", result))
			decryptBtn.Enable()
		}()
	})

	window.SetOnDropped(func(pos fyne.Position, uris []fyne.URI) {
		if len(uris) > 0 {
			handleDrop(uris[0])
		}
	})

	content := container.NewVBox(
		widget.NewCard("", "Выбор зашифрованного файла", container.NewVBox(
			pathLabel,
			selectBtn,
		)),
		widget.NewCard("", "Пароль", container.NewVBox(
			passwordEntry,
		)),
		widget.NewCard("", "Параметры", container.NewVBox(
			removeCheck,
		)),
		decryptBtn,
		statusLabel,
	)

	window.SetContent(container.NewPadded(content))
	window.Show()
}

