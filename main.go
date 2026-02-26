package main

import (
	"strconv"

	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
)

func main() {
	b := core.NewBody()
	counter(b)
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
