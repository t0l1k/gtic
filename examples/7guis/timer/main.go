package timer

import (
	"etic"
	"etic/examples/7guis/data"
	"fmt"
	"log"
	"time"
)

func Timer(id etic.ElementID, fn etic.SlotFn[*etic.Button]) *etic.Element {
	s := etic.NewElement(id)
	var (
		duration time.Duration = 30 * time.Second
		elapsed  time.Duration
		reset    func()
		running  bool
	)
	lblHelper := etic.NewLabel("scene.timer.label.helper", "A timer with progress, adjustable duration, and reset.")
	lblHelper.Font.Set(etic.FontSystem)
	lblDur := etic.NewLabel("scene.timer.label.duration", "")
	progress := etic.NewProgress("scene.timer.progress", 1, 0, 1, etic.Horizontal)
	slider := etic.NewSlider("scene.timer.slider", 0.3, 0, 1, 0.01, etic.Horizontal, func(f float32) {
		duration = time.Duration(f*100) * time.Second
		if !running {
			reset()
		}
		log.Println("slider value", f)
	})
	btnReset := etic.NewButton("scene.timer.button.reset", "Reset", func(b *etic.Button) {
		running = !running
		if running {
			b.Text.Set("Stop")
			log.Println("Stop")
		} else {
			b.Text.Set("Start")
			log.Println("Start")
		}
	})

	reset = func() {
		elapsed = 0
		progress.Value.Set(0)
		lblDur.Text.Set(fmt.Sprintf("%0.1f", duration.Seconds()))
	}

	s.OnUpdate = func(t *etic.Console) {
		if !running {
			return
		}
		elapsed += t.Tick()
		if elapsed >= duration {
			elapsed -= duration
			running = false
			btnReset.Text.Set("Start")
			log.Println("Done")
		}
		progress.Value.Set(float32(elapsed.Seconds() / duration.Seconds()))
		lblDur.Text.Set(fmt.Sprintf("%.1f", (duration - elapsed).Seconds()))
	}

	s.Layout().Set(etic.NewVerticalFlexLayout([]float32{0.1, 0.9}, 0.05))
	s.Add(etic.TopBar("scene.timer.topbar", data.Timer, data.QuitDemo, fn))

	contBtns := etic.NewElement("scene.timer.buttons.container")
	contBtns.Layout().Set(etic.NewVerticalBoxLayout(0.05))
	contBtns.Add(lblHelper)
	contBtns.Add(progress)
	contBtns.Add(slider)
	contBtns.Add(lblDur)
	contBtns.Add(btnReset)
	s.Add(contBtns)
	return s
}
