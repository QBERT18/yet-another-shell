package gui

import (
	"errors"
	"os"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/QBERT18/yet-another-shell/shell"
)

type outputWriter struct {
	buffer strings.Builder
	output *widget.Label
	scroll *container.Scroll
	mu     sync.Mutex
}

func (w *outputWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buffer.Write(p)
	w.output.SetText(w.buffer.String())
	w.scroll.ScrollToBottom()
	return len(p), nil
}

func getPrompt() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "> "
	}
	return cwd + "> "
}

func Run() {
	myApp := app.New()
	myApp.Settings().SetTheme(&terminalTheme{})
	window := myApp.NewWindow("Yet Another Shell")

	logoResource := fyne.NewStaticResource("logo.png", logoBytes)
	logoImage := canvas.NewImageFromResource(logoResource)
	logoImage.FillMode = canvas.ImageFillContain
	logoImage.SetMinSize(fyne.NewSize(40, 40))

	titleLabel := widget.NewLabel("Yet Another Shell")
	titleLabel.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

	topBar := container.NewHBox(logoImage, titleLabel)

	outputLabel := widget.NewLabel("")
	outputLabel.Wrapping = fyne.TextWrapBreak
	outputLabel.TextStyle = fyne.TextStyle{Monospace: true}
	outputLabel.Selectable = true
	scrollArea := container.NewScroll(outputLabel)

	writer := &outputWriter{
		output: outputLabel,
		scroll: scrollArea,
	}

	engine := shell.NewEngine(writer, writer)

	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("Enter command...")
	inputEntry.TextStyle = fyne.TextStyle{Monospace: true}
	inputEntry.OnSubmitted = func(cmd string) {
		prompt := getPrompt()

		writer.mu.Lock()
		writer.buffer.WriteString(prompt + cmd + "\n")
		writer.output.SetText(writer.buffer.String())
		writer.scroll.ScrollToBottom()
		writer.mu.Unlock()

		inputEntry.SetText("")
		inputEntry.Disable()

		go func() {
			err := engine.Execute(cmd)
			if err != nil {
				if errors.Is(err, shell.ErrExit) {
					window.Close()
					return
				}
				writer.Write([]byte("error: " + err.Error() + "\n"))
			}

			inputEntry.Enable()
			window.Canvas().Focus(inputEntry)
		}()
	}

	content := container.NewBorder(
		topBar,     // top
		inputEntry, // bottom
		nil,        // left
		nil,        // right
		scrollArea, // center
	)

	window.SetContent(content)
	window.Resize(fyne.NewSize(800, 500))
	window.ShowAndRun()
}
