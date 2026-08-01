package game

import (
	"etic"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

const Title = "Увернись от Крипов!"

type GameState int

func (s GameState) String() string {
	return []string{"StateMenu", "StateStarting", "StatePlaying", "StateGameOver"}[s]
}

const (
	StateMenu GameState = iota
	StateStarting
	StatePlaying
	StateGameOver
)

var scale = 0.4

type Game struct {
	etic.Element
	ready                            bool
	state                            etic.Property[GameState]
	score                            etic.Property[int]
	lblScore, lblMessage             *etic.Label
	btnStart                         *etic.Button
	timerStart, timerScore, timerMob *etic.Timer
	player                           *player
	creeps                           map[int]*creep
	creepId                          int
}

func New() *Game {
	g := &Game{}
	g.lblMessage = etic.NewLabel("label.message", "Увернись \nот Крипов!\nPress SPACE\n to start")
	g.lblMessage.Font.Set(etic.FontTitle)
	g.lblScore = etic.NewLabel("label.score", "0")
	g.lblScore.Font.Set(etic.FontTitle)
	g.btnStart = etic.NewButton("button.start", "Start", func(b *etic.Button) {
		g.state.Set(StateStarting)
	})
	g.btnStart.Font.Set(etic.FontTitle)
	g.score = *etic.NewProperty(0)
	g.score.OnChange.Connect(func(value int) {
		g.lblScore.Text.Set(strconv.Itoa(value))
	})
	g.state = *etic.NewProperty(StateMenu)
	g.state.OnChange.Connect(func(value GameState) {
		switch value {
		case StateMenu:
			g.lblMessage.Text.Set("Увернись \nот Крипов!\nPress SPACE\n to start")
			g.lblMessage.Show()
			g.btnStart.Show()
			g.score.Set(0)
			g.player.Hide()
			log.Println(value)
		case StateStarting:
			for k, creep := range g.creeps {
				delete(g.creeps, k)
				g.Remove(creep)
			}
			g.creepId = 0
			g.player.ClearReady()
			g.timerStart.Start()
			g.lblMessage.Text.Set("Get Ready!")
			g.btnStart.Hide()
			log.Println(value)
		case StatePlaying:
			g.lblMessage.Hide()
			g.timerScore.Start()
			g.timerMob.Start()
			log.Println(value)
		case StateGameOver:
			g.player.Hide()
			g.timerStart.Start()
			g.timerMob.Stop()
			g.timerScore.Stop()
			g.lblMessage.Show()
			g.lblMessage.Text.Set(fmt.Sprintf("Game Over!\nScore: %d", g.score.Get()))
			log.Println(value)
		}
	})
	g.timerStart = etic.NewTimer("timer.start")
	g.timerStart.Setup(2*time.Second, false, true)
	g.timerStart.Timeout.Connect(func(t *etic.Timer) {
		switch g.state.Get() {
		case StateStarting:
			g.state.Set(StatePlaying)
		case StateGameOver:
			g.state.Set(StateMenu)
		}
	})

	g.timerScore = etic.NewTimer("timer.score")
	g.timerScore.Setup(1*time.Second, false, false)
	g.timerScore.Timeout.Connect(func(t *etic.Timer) {
		g.score.Set(g.score.Get() + 1)
	})

	g.timerMob = etic.NewTimer("timer.mob")
	g.timerMob.Setup(500*time.Millisecond, false, false)
	g.timerMob.Timeout.Connect(func(t *etic.Timer) {
		if !(g.state.Get() == StatePlaying) {
			return
		}
		id := g.creepId
		g.creepId++
		if g.creeps == nil {
			g.creeps = make(map[int]*creep)
		}

		mobKind := [][]etic.SpriteID{{enemyFly1, enemyFly2}, {enemySwim1, enemySwim2}, {enemyWalk1, enemyWalk2}}

		creep := NewCreep(etic.ElementID(fmt.Sprintf("creep %v", strconv.Itoa(id)))).Setup(mobKind[rand.Intn(3)])
		g.creeps[id] = creep
		g.Add(creep)
		// log.Println("new creep", id, len(g.creeps), creep.Bounds().Get())
	})

	g.player = NewPlayer("player")
	g.player.Setup([]etic.SpriteID{playerSprWalk1, playerSprWalk2}, []etic.SpriteID{playerSprUp1, playerSprUp2})
	g.player.hit.Connect(func(value bool) {
		if value {
			g.state.Set(StateGameOver)
		}
	})
	g.state.Set(StateMenu)
	g.Add(g.player)
	g.Add(g.lblMessage)
	g.Add(g.lblScore)
	g.Add(g.btnStart)
	g.Add(g.timerStart)
	g.Add(g.timerScore)
	g.Add(g.timerMob)
	return g
}
func (g *Game) Arrange(w, h float32) {
	marginW := w * 0.1
	marginH := h * 0.1
	g.lblScore.Bounds().Set(etic.Rect(0, 0, w, marginH*2))
	g.lblMessage.Bounds().Set(etic.Rect(marginW, marginH*2, w-marginW, marginH*7))
	g.btnStart.Bounds().Set(etic.Rect(marginW, marginH*7, w-marginW, marginH*8))
}
func (g *Game) Init(t *etic.Console) {
	g.Arrange(t.Width, t.Height)
	g.creeps = nil
	g.ready = true
}
func (g *Game) Update(t *etic.Console) {
	g.Element.Update(t)
	if t.Btn(ebiten.KeySpace) && g.state.Get() == StateMenu {
		g.state.Set(StateStarting)
	}
	for id, creep := range g.creeps {
		creep.IsColliding = false
		g.player.IsColliding = false
		if g.state.Get() == StatePlaying {
			creepBounds := creep.Bounds().Get()
			playerBounds := g.player.Bounds().Get()
			if creepBounds.Overlaps(playerBounds) {
				g.player.hit.Emit(true)
				creep.IsColliding = true
				g.player.IsColliding = true
			}
		}
		if len(g.creeps) > 0 {
			r := etic.Rect(0, 0, t.Width, t.Height).Inset(-100)
			if !r.Overlaps(creep.Bounds().Get()) && creep.IsReady() {
				g.Remove(creep)
				delete(g.creeps, id)
			}
		}
	}
}
func (g *Game) TIC(t *etic.Console) {
	if !g.ready {
		g.Init(t)
	}
	t.Cls(colornames.Teal)
	g.Update(t)
	g.Draw(t)
}
