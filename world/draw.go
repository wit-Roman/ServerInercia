package world

import (
	"math"
	"server/constant"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
)

var op = &ebiten.DrawImageOptions{}

//op.ColorM.Scale(255, 255, 255, 1)

func drawPlayer(ship TShip, screen *ebiten.Image) {
	op.GeoM.Reset()
	angle := ship.Body.Angle()

	// TODO: проверить по модулю
	v0 := cp.Vector{Ship.X, Ship.Y}
	//getShipPartPosition(Ship.X0, Ship.Y0, angle, body.UserData)

	op.GeoM.Translate(v0.X, v0.Y)

	impulseSum := 0.0
	ship.Body.EachConstraint(func(constraint *cp.Constraint) {
		impulseSum += constraint.Class.GetImpulse()
	})

	x1 := v0.X
	x2 := v0.X
	x3 := v0.X
	y1 := v0.Y
	y2 := v0.Y
	y3 := v0.Y

	x1 -= constant.Size
	x3 += constant.Size
	y1 -= constant.Size
	y3 -= constant.Size

	//center s := atan2(P.y-T.y,P.x-T.x) - atan2(A.y-T.y,A.x-T.x)
	sin := math.Sin(angle)
	cos := math.Cos(angle)

	x1r := (x1-v0.X)*cos - (y1-v0.Y)*sin + v0.X
	y1r := (x1-v0.X)*sin + (y1-v0.Y)*cos + v0.Y

	x2r := (x2-v0.X)*cos - (y2-v0.Y)*sin + v0.X
	y2r := (x2-v0.X)*sin + (y2-v0.Y)*cos + v0.Y

	x3r := (x3-v0.X)*cos - (y3-v0.Y)*sin + v0.X
	y3r := (x3-v0.X)*sin + (y3-v0.Y)*cos + v0.Y

	vertices := []ebiten.Vertex{{float32(x1r), float32(y1r), 0, 0, 0, 0, 0, 0}, {float32(x2r), float32(y2r), 0, 0, 0, 0, 0, 0}, {float32(x3r), float32(y3r), 0, 0, 0, 0, 0, 0}}
	indices := []uint16{0, 1, 2}

	screen.DrawTriangles(vertices, indices, &ebiten.Image{}, &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	})
}
