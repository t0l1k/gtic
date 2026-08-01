package main

import (
	"etic"
	"math"
	"math/rand"
	"slices"
	"strconv"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

var (
	//go:embed sprites.png
	spritesPNG    []byte
	w, h          float32 = 800, 800
	width, height float32 = 16, 16
	gridSize              = w / width
)

type Dir int

const (
	Up Dir = iota
	Down
	Left
	Right
)

type Snake struct {
	body    []etic.Point[float32]
	dir     etic.Point[float32]
	dirName Dir
	speed   int
	game    *Game
}

func (s *Snake) Reset(g *Game) {
	s.game = g
	s.body = nil
	s.body = append(s.body, etic.Pt[float32](13, 8), etic.Pt[float32](14, 8), etic.Pt[float32](15, 8))
	s.dir = etic.Left
	s.dirName = Left
	s.speed = 10
}
func (s *Snake) change(n Dir, dir etic.Point[float32]) { s.dirName = n; s.dir = dir }
func (s *Snake) update(t *etic.Console) {
	switch {
	case t.BtnP(ebiten.KeyUp) && s.dir.Y == 0:
		s.change(Up, etic.Up)
	case t.BtnP(ebiten.KeyArrowDown) && s.dir.Y == 0:
		s.change(Down, etic.Down)
	case t.BtnP(ebiten.KeyArrowLeft) && s.dir.X == 0:
		s.change(Left, etic.Left)
	case t.BtnP(ebiten.KeyArrowRight) && s.dir.X == 0:
		s.change(Right, etic.Right)
	}
	if t.Frame()%s.speed == 0 {
		newPos := s.body[0].Add(s.dir)
		if !newPos.In(etic.Rect(0, 0, width, height)) {
			s.game.gameOver = true
		}
		for i := 0; i < len(s.body); i++ {
			if s.body[i] == newPos {
				s.game.gameOver = true
			}
		}
		s.body = slices.Insert(s.body, 0, newPos)
		s.body = s.body[:len(s.body)-1]
	}
	if s.isEat(s.game.food.pos) {
		s.body = append(s.body, s.game.food.pos)
		s.game.food.Reset(setNewFood(width, height, s.game.snake.body))
		s.game.score++
		if len(s.body)%5 == 0 {
			s.speed--
		}
	}
}
func (s *Snake) draw(t *etic.Console) {
	for i, v := range s.body {
		sprite := 3
		flip := 0.0
		if i == 0 {
			switch s.dirName {
			case Up:
				sprite = 1
			case Down:
				sprite = 1
				flip = 2
			case Left:
				sprite = 2
				flip = 1
			case Right:
				sprite = 2
			}
		}
		t.Spr(
			etic.SpriteID(strconv.Itoa(sprite)),
			v.X*gridSize,
			v.Y*gridSize,
			etic.SprOpt{
				Scale: gridSize / 8,
				Flip:  etic.Flip(flip),
			})
	}
}
func (s *Snake) isEat(food etic.Point[float32]) bool { return s.body[0].Eq(food) }

type Food struct{ pos etic.Point[float32] }

func (f *Food) Reset(x, y float32) { f.pos.X = x; f.pos.Y = y }
func (f *Food) draw(t *etic.Console) {
	t.Spr(etic.SpriteID(
		strconv.Itoa(0)),
		f.pos.X*gridSize,
		f.pos.Y*gridSize,
		etic.SprOpt{
			Scale: gridSize / 8,
		})
}
func setNewFood(w, h float32, snakeBody []etic.Point[float32]) (x, y float32) {
	x, y = float32(math.Floor(float64(rand.Float32()*w))), float32(math.Floor(float64(rand.Float32()*h)))
	for _, v := range snakeBody {
		if v.X == x && v.Y == y {
			setNewFood(w, h, snakeBody)
		}
	}
	return x, y
}

type Field struct{}

func (f *Field) draw(t *etic.Console) {
	var i float32
	for i = 0; i < t.Height; i++ {
		t.Line(0, i*gridSize, t.Width, i*gridSize, colornames.Darkgreen)
		t.Line(i*gridSize, 0, i*gridSize, t.Height, colornames.Darkgreen)
	}
}

type Game struct {
	score           int
	field           Field
	food            Food
	snake           Snake
	Ready, gameOver bool
}

func NewGame() *Game { return &Game{} }
func (g *Game) Reset(t *etic.Console) {
	g.gameOver = false
	g.score = 0
	g.snake.Reset(g)
	g.food.Reset(8, 8)
}
func (g *Game) update(t *etic.Console) {
	if !g.Ready {
		g.Reset(t)
		g.Ready = true
	}
	if g.gameOver {
		return
	}
	g.snake.update(t)
}
func (g *Game) draw(t *etic.Console) {
	t.Cls(colornames.Darkgray)
	g.field.draw(t)
	g.food.draw(t)
	g.snake.draw(t)
	scale := gridSize / float32(t.Face(etic.FontSystem).Size)
	t.Print("Score:"+strconv.Itoa(g.score), 1, 1, colornames.Aqua, scale)
	t.Print("Speed:"+strconv.Itoa(g.snake.speed), 1, gridSize, colornames.Aqua, scale)
	if g.gameOver {
		t.Print("Game Over!!!\nPress Enter to play again.", 1, gridSize*2, colornames.Red, scale)
	}
}
func StartGame(title string) *etic.Cart {
	cart := etic.NewCart(title, w, h)
	game := NewGame()
	cart.OnBoot = func(t *etic.Console) {
		game.Ready = false
		t.LoadSpriteSheet(etic.ApplyColorKey(etic.LoadImage(spritesPNG), colornames.Black), 8, 8)
	}
	cart.OnTic = func(t *etic.Console) {
		if t.BtnP(ebiten.KeyEscape) {
			t.Exit()
		}
		if t.BtnP(ebiten.KeyEnter) {
			t.Reset()
		}
		game.update(t)
		game.draw(t)
	}
	return cart
}

func main() { etic.Load(StartGame("Snake")).Run() }
