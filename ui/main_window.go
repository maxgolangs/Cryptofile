package ui

import (
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MainWindow struct {
	app    fyne.App
	window fyne.Window
}

func NewMainWindow() *MainWindow {
	application := app.NewWithID("cryptofile.app")
	
	if runtime.GOOS != "windows" {
		application.Settings().SetTheme(theme.DarkTheme())
	}
	
	window := application.NewWindow("CryptoFile by @MaxGolang")
	window.Resize(fyne.NewSize(700, 550))
	window.CenterOnScreen()

	return &MainWindow{
		app:    application,
		window: window,
	}
}

func (mw *MainWindow) Show() {
	mw.buildUI()
	mw.window.ShowAndRun()
}

func (mw *MainWindow) buildUI() {
	title := widget.NewLabel("🔐 CryptoFile")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := widget.NewLabel("by @MaxGolang - Безопасное шифрование файлов и директорий")
	subtitle.Alignment = fyne.TextAlignCenter
	subtitle.Importance = widget.LowImportance

	encryptFileBtn := widget.NewButton("Зашифровать файл или директорию", func() {
		ShowEncryptSingleWindow(mw.window)
	})
	encryptFileBtn.Importance = widget.HighImportance

	decryptFileBtn := widget.NewButton("Расшифровать файл", func() {
		ShowDecryptSingleWindow(mw.window)
	})

	encryptBatchBtn := widget.NewButton("Зашифровать директорию (рекурсивно)", func() {
		ShowEncryptBatchWindow(mw.window)
	})

	decryptBatchBtn := widget.NewButton("Расшифровать директорию (рекурсивно)", func() {
		ShowDecryptBatchWindow(mw.window)
	})

	headerCard := widget.NewCard("", "", container.NewVBox(
		container.NewCenter(title),
		container.NewCenter(subtitle),
	))

	buttonsCard := widget.NewCard("Операции", "", container.NewVBox(
		encryptFileBtn,
		decryptFileBtn,
		widget.NewSeparator(),
		encryptBatchBtn,
		decryptBatchBtn,
	))

	content := container.NewVBox(
		headerCard,
		widget.NewSeparator(),
		buttonsCard,
	)

	mw.window.SetContent(container.NewPadded(content))
}


