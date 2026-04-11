package main

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
)

func main() {
	b := core.NewBody("7GUIs")

	tasks := []struct {
		name   string
		runner func(*core.Body)
	}{
		{name: "Counter", runner: counter},
		{name: "Temperature Converter", runner: temperatureConverter},
		{name: "Flight Booker", runner: flightBooker},
		{name: "Timer", runner: timer},
		{name: "CRUD", runner: crud},
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

func flightBooker(body *core.Body) {
	const (
		oneWayFlight = "one-way flight"
		returnFlight = "return flight"
	)

	option := oneWayFlight
	startRaw := time.Now().Format(time.DateOnly)
	returnRaw := startRaw

	optionChooser := core.Bind(&option, core.NewChooser(body).SetStrings(oneWayFlight, returnFlight))
	startInput := core.Bind(&startRaw, core.NewTextField(body))
	returnInput := core.Bind(&returnRaw, core.NewTextField(body))
	bookButton := core.NewButton(body).SetText("Book")

	returnInput.SetEnabled(false)

	setInputStyle := func(input *core.TextField, valid bool) {
		input.Styler(func(s *styles.Style) {
			if valid {
				s.Background = colors.Scheme.SurfaceContainer
			} else {
				s.Background = colors.Scheme.Error.Container
			}
		})
	}

	validate := func() {
		startDate, err := time.Parse(time.DateOnly, startRaw)
		setInputStyle(startInput, err == nil)

		ok := err == nil
		if option == returnFlight {
			returnDate, err := time.Parse(time.DateOnly, returnRaw)
			setInputStyle(returnInput, err == nil)
			ok = ok && err == nil && !startDate.After(returnDate)
		}

		bookButton.SetEnabled(ok)
		bookButton.Update()
	}

	optionChooser.OnChange(func(e events.Event) {
		returnInput.SetEnabled(option == returnFlight)
		returnInput.Update()
		validate()
	})

	startInput.OnChange(func(e events.Event) { validate() })
	returnInput.OnChange(func(e events.Event) { validate() })

	bookButton.OnClick(func(e events.Event) {
		switch option {
		case oneWayFlight:
			core.MessageSnackbar(body, fmt.Sprintf("You have booked a one-way flight on %s.", startRaw))
		case returnFlight:
			core.MessageSnackbar(body, fmt.Sprintf("You have booked a return flight on %s and %s.", startRaw, returnRaw))
		}
	})
}

func timer(body *core.Body) {
	ticks := 0
	target := 0.0
	elapsed := 0.0

	progressFrame := core.NewFrame(body)
	progressFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Row
		s.Align.Items = styles.Center
		s.Grow.X = 1
	})
	core.NewText(progressFrame).SetText("Elapsed Time:")
	progressBar := core.Bind(&elapsed, core.NewMeter(progressFrame).SetMax(60))

	elapsedDisplay := core.NewText(body).SetText("0.0 s")

	controlFrame := core.NewFrame(body)
	controlFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Row
		s.Align.Items = styles.Center
		s.Grow.X = 1
	})
	core.NewText(controlFrame).SetText("Duration:")
	core.Bind(&target, core.NewSlider(controlFrame).SetMax(60).SetEnforceStep(true))

	core.NewButton(body).SetText("Reset").AsButton().OnClick(func(e events.Event) {
		ticks = 0
		elapsed = 0
		progressBar.Update()
		elapsedDisplay.SetText("0.0 s").Update()
	})

	go func() {
		for range time.Tick(100 * time.Millisecond) {
			body.AsyncLock()
			if ticks < int(math.Round(target*10)) {
				ticks++
				elapsed = float64(ticks) / 10
				progressBar.Update()
				elapsedDisplay.SetText(fmt.Sprintf("%.1f s", elapsed)).Update()
			}
			body.AsyncUnlock()
		}
	}()
}

func crud(body *core.Body) {
	type Person struct {
		Name    string
		Surname string
	}

	people := []Person{}

	body.Styler(func(s *styles.Style) { s.Gap.Y.Dp(16) })

	type Search struct {
		FilterPrefix string
	}

	search := &Search{}
	filter := func() []Person {
		if search.FilterPrefix == "" {
			return people
		}
		filtered := []Person{}
		for _, person := range people {
			if strings.HasPrefix(person.Surname, search.FilterPrefix) {
				filtered = append(filtered, person)
			}
		}
		return filtered
	}
	searchForm := core.NewForm(body).SetStruct(search)

	formFrame := core.NewFrame(body)
	formFrame.Styler(func(s *styles.Style) {
		s.Grow.SetScalar(1)
		s.Gap.X.Dp(16)
	})

	items := &[]string{}
	selected := -1
	personList := core.NewList(formFrame).SetSlice(items).BindSelect(&selected)
	personList.Styler(func(s *styles.Style) {
		s.Max.X.Pw(16.67)
		s.Border.Width.Set(units.Dp(1))
	})
	personList.Updater(func() {
		filtered := filter()
		*items = make([]string, 0, len(filtered))
		for _, person := range filtered {
			*items = append(*items, fmt.Sprintf("%v, %v", person.Name, person.Surname))
		}
	})

	person := &Person{}
	personForm := core.NewForm(formFrame).SetStruct(person)
	personForm.Styler(func(s *styles.Style) { s.Grow.X = 1 })

	opFrame := core.NewFrame(body)
	core.NewButton(opFrame).SetText("Create").OnClick(func(e events.Event) {
		if person.Name != "" && person.Surname != "" {
			people = append(people, *person)
			personList.Update()

			*person = Person{}
			personForm.Update()
		}
	})

	updateButton := core.NewButton(opFrame).SetText("Update").SetEnabled(false)
	updateButton.OnClick(func(e events.Event) {
		if person.Name != "" && person.Surname != "" {
			people[selected] = *person
			personList.SetSelectedIndex(-1)
			personList.Send(events.Select)
			personList.Update()

			*person = Person{}
			personForm.Update()
		}
	})

	deleteButton := core.NewButton(opFrame).SetText("Delete").SetEnabled(false)
	deleteButton.OnClick(func(e events.Event) {
		people = slices.Delete(people, selected, selected+1)
		personList.SetSelectedIndex(-1)
		personList.Send(events.Select)
		personList.Update()
	})

	personList.OnChange(func(e events.Event) {
		enabled := selected != -1
		updateButton.SetEnabled(enabled).Update()
		deleteButton.SetEnabled(enabled).Update()
	})

	searchForm.OnChange(func(e events.Event) { personList.Update() })
}
