package temper

import (
	"etic"
	"etic/examples/7guis/data"
	"fmt"
	"strconv"
)

func toF(c float64) float64 { return c*(9.0/5.0) + 32.0 }
func toC(f float64) float64 { return (f - 32.0) * (5.0 / 9.0) }

func TemeratureConv(id etic.ElementID, fn etic.SlotFn[*etic.Button]) *etic.Element {
	tempConvCelsius := etic.NewInputLine("tempconv.c", "Celsius", func(il *etic.InputLine) {})
	tempConvFahrenheit := etic.NewInputLine("tempconv.f", "Fahrenheit", func(il *etic.InputLine) {})

	tempConvCelsius.OnChanged.Connect(func(s string) {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return
		}
		result := toF(v)
		tempConvFahrenheit.Text.Set("")
		tempConvFahrenheit.PlaceHolder.Set(fmt.Sprintf("%.1f", result))
	})
	tempConvFahrenheit.OnChanged.Connect(func(s string) {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return
		}
		result := toC(v)
		tempConvCelsius.Text.Set("")
		tempConvCelsius.PlaceHolder.Set(fmt.Sprintf("%.1f", result))
	})

	t := etic.NewElement(id)
	t.Layout().Set(etic.NewVerticalFlexLayout([]float32{0.1, 0.1, 0.8}, 0.05))
	topbar := etic.TopBar("scene.temperature.converter.topbar", data.TempConv, data.QuitDemo, fn)
	helper := etic.NewLabel("scene.temperature.converter.helper", "Bidirectional data flow with user-provided text input.")
	helper.Font.Set(etic.FontSystem)
	cont := etic.NewElement("scene.temperature.converter.container")
	cont.Layout().Set(etic.NewHorizontalBoxLayout(0.05))
	cont.Add(tempConvCelsius)
	cont.Add(etic.NewLabel("scene.temperature.converter.label", "="))
	cont.Add(tempConvFahrenheit)
	t.Add(topbar)
	t.Add(helper)
	t.Add(cont)
	return t
}
