package world

import (
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"

	"server/constant"
)

var Players = make(map[string]TPlayer)

type TPlayer struct {
	Ship     TShip
	InputCh  chan int8
	OutputCh chan interface{}
}

type Game struct {
	isStarted bool
	space     *cp.Space
	//ctx         *context.Context
	//viewport    Viewport
}

var g *Game

func (g *Game) Update() error {
	for _, player := range Players {
		switch <-player.InputCh {
		case 0:
			OnPressGas()
		case 1:
			OnPressLeft()
		case 2:
			OnPressRight()
		}
	}

	g.space.Step(0.015)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, player := range Players {
		drawPlayer(player.Ship, screen)
	}
}

// Layout sets window size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 1200
}

func Create() {
	rand.Seed(time.Now().UnixNano())

	ebiten.SetWindowSize(constant.ScreenWidth, constant.ScreenHeight)
	ebiten.SetWindowTitle("Server Test")

	g = NewGame()

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

func NewGame() *Game {
	space := cp.NewSpace()
	space.Iterations = 1
	space.SetDamping(0.975)

	return &Game{isStarted: true, space: space}
}

func CreateNewPlayer(token string) TPlayer {
	currLen := float64(len(Players)-20) * 500.0
	x := currLen * 2
	y := currLen

	body := createNewPlayerBody(x, y)

	Players[token] = TPlayer{TShip{X: x, Y: y, Body: body}, make(chan int8), make(chan interface{})}

	command := make(map[string]interface{})
	command["startPosition"] = struct {
		X float64
		Y float64
	}{X: Players[token].Ship.X, Y: Players[token].Ship.Y}
	go func() {
		Players[token].OutputCh <- command
	}()

	return Players[token]
}

func createNewPlayerBody(posX, posY float64) *cp.Body {
	angle := 0

	var vertices = []cp.Vector{
		{-constant.Size, -constant.Size},
		{0, 0},
		{constant.Size, -constant.Size},
	}
	var triangleMoment = cp.MomentForPoly(Ship.Weight, 3, vertices, cp.Vector{}, 0)
	body := g.space.AddBody(cp.NewBody(Ship.Weight, triangleMoment))
	body.SetAngle(float64(angle) * math.Pi / 180)
	body.SetPosition(cp.Vector{})
	shape := g.space.AddShape(cp.NewPolyShape(body, 3, vertices, cp.NewTransformIdentity(), 0))
	shape.SetElasticity(0.4)
	shape.SetFriction(0.8)

	return body
}
