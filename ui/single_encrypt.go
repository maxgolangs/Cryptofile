package ui

import (
	"fmt"
	"os"

	"cryptor/pkg/encrypt"
	"cryptor/constants"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowEncryptSingleWindow(parent fyne.Window) {
	window := fyne.CurrentApp().NewWindow("Шифрование файла или директории")
	window.Resize(fyne.NewSize(600, 400))
	window.CenterOnScreen()

	var selectedPath string
	var password string
	var removeOriginal bool

	pathLabel := widget.NewLabel("Файл не выбран\n💡 Или перетащите файл/папку сюда")
	pathLabel.Wrapping = fyne.TextWrapWord

	handleDrop := func(uri fyne.URI) {
		if uri == nil {
			return
		}
		selectedPath = uri.Path()
		pathLabel.SetText(fmt.Sprintf("Выбрано: %s", selectedPath))
	}

	selectFileBtn := widget.NewButton("Выбрать файл", func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()
			
			uri := reader.URI()
			if uri != nil {
				selectedPath = uri.Path()
				pathLabel.SetText(fmt.Sprintf("Выбрано: %s", selectedPath))
			}
		}, window)
		fileDialog.Show()
	})

	selectDirBtn := widget.NewButton("Выбрать директорию", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			selectedPath = uri.Path()
			pathLabel.SetText(fmt.Sprintf("Выбрано: %s", selectedPath))
		}, window).Show()
	})

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Введите пароль (минимум 6 символов)")

	confirmPasswordEntry := widget.NewPasswordEntry()
	confirmPasswordEntry.SetPlaceHolder("Подтвердите пароль")

	removeCheck := widget.NewCheck("Удалить оригинал после шифрования", func(checked bool) {
		removeOriginal = checked
	})
	removeCheck.SetChecked(true)

	statusLabel := widget.NewLabel("")
	statusLabel.Wrapping = fyne.TextWrapWord

	var encryptBtn *widget.Button
	encryptBtn = widget.NewButton("Зашифровать", func() {
		if selectedPath == "" {
			statusLabel.SetText("❌ Пожалуйста, выберите файл или директорию")
			return
		}

		password = passwordEntry.Text
		if len(password) < constants.MinPasswordLength {
			statusLabel.SetText(fmt.Sprintf("❌ Пароль должен содержать минимум %d символов", constants.MinPasswordLength))
			return
		}

		if password != confirmPasswordEntry.Text {
			statusLabel.SetText("❌ Пароли не совпадают")
			return
		}

		if _, err := os.Stat(selectedPath); err != nil {
			statusLabel.SetText(fmt.Sprintf("❌ Ошибка: %v", err))
			return
		}

		encryptBtn.Disable()
		statusLabel.SetText("⏳ Выполняется шифрование...")

		go func() {
			result, err := encrypt.EncryptPath(selectedPath, password, removeOriginal)
			if err != nil {
				statusLabel.SetText(fmt.Sprintf("❌ Ошибка: %v", err))
				encryptBtn.Enable()
				return
			}

			statusLabel.SetText(fmt.Sprintf("✅ Успешно! Зашифрованный файл: %s", result))
			encryptBtn.Enable()
		}()
	})

	window.SetOnDropped(func(pos fyne.Position, uris []fyne.URI) {
		if len(uris) > 0 {
			handleDrop(uris[0])
		}
	})

	content := container.NewVBox(
		widget.NewCard("", "Выбор файла или директории", container.NewVBox(
			pathLabel,
			container.NewHBox(selectFileBtn, selectDirBtn),
		)),
		widget.NewCard("", "Пароль", container.NewVBox(
			passwordEntry,
			confirmPasswordEntry,
		)),
		widget.NewCard("", "Параметры", container.NewVBox(
			removeCheck,
		)),
		encryptBtn,
		statusLabel,
	)

	window.SetContent(container.NewPadded(content))
	window.Show()
}

