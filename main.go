package main

import (
	"strconv"

	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
)

func main() {
	b := core.NewBody("7GUIs")

	tasks := []struct {
		name   string
		runner func(*core.Body)
	}{
		{name: "Counter", runner: counter},
		{name: "Temperature Converter", runner: temperatureConverter},
	}

	for _, task := range tasks {
		core.NewButton(b).SetText(task.name).OnClick(func(e events.Event) {
			taskBody := core.NewBody(task.name)
			task.runner(taskBody)
			taskBody.RunWindow()
		})
	}

	b.RunMainWindow()
}

func counter(body *core.Body) {
	i := 0
	textField := core.NewTextField(body).SetText(strconv.Itoa(i))
	textField.SetReadOnly(true)
	core.NewButton(body).SetText("Count").OnClick(func(e events.Event) {
		i++
		textField.SetText(strconv.Itoa(i)).Update()
	})
}

func temperatureConverter(body *core.Body) {
	state := struct {
		unit string
		raw  string
	}{}

	body.Styler(func(s *styles.Style) {
		s.Direction = styles.Row
		s.Align.Items = styles.Center
	})

	celsiusInput := core.NewTextField(body)
	core.NewText(body).SetText("Celsius")
	core.NewText(body).SetText("=")
	fahrenheitInput := core.NewTextField(body)
	core.NewText(body).SetText("Fahrenheit")

	celsiusInput.OnInput(func(e events.Event) {
		state.unit = "c"
		state.raw = celsiusInput.Text()

		if state.raw == "" {
			fahrenheitInput.SetText("").Update()
			return
		}

		parsed, err := strconv.ParseFloat(state.raw, 64)
		if err != nil {
			return
		}
		fahrenheitInput.SetText(strconv.FormatFloat(parsed*9/5+32, 'f', -1, 64)).Update()
	})

	fahrenheitInput.OnInput(func(e events.Event) {
		state.unit = "f"
		state.raw = fahrenheitInput.Text()

		if state.raw == "" {
			celsiusInput.SetText("").Update()
			return
		}

		parsed, err := strconv.ParseFloat(state.raw, 64)
		if err != nil {
			return
		}
		celsiusInput.SetText(strconv.FormatFloat((parsed-32)*5/9, 'f', -1, 64)).Update()
	})
}
