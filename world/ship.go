package world

import (
	"image/color"
	"math"
	"server/constant"

	"github.com/jakecoffman/cp"
)

type TShip struct {
	Weight float64    // нужно привязать к реальному
	Color  color.RGBA // пока не используется
	X      float64    // реальная координата на экране
	Y      float64    // реальная координата на экране
	X0     float64    // координата для отрисовки на экране
	Y0     float64    // координата для отрисовки на экране
	// Граничные индексы видимых тайлов
	Path              [4]int
	Angle             int
	AngleVelocity     float64
	AngleAcceleration float64
	VelocityX         float64
	VelocityY         float64
	// Скорость движения по экрану Velocity*Multiplier
	VelocityMultiplier float64
	VelocityX0         float64
	VelocityY0         float64
	ForceX             float64
	ForceY             float64
	Body               *cp.Body
	//PrevUpdateTime time.Time
}

var Ship = TShip{
	Weight: 40,
	Color:  color.RGBA{255, 0, 55, 255},
	// TODO случайная позиция
	X:                 200,
	Y:                 200,
	X0:                200,
	Y0:                200,
	Path:              [4]int{0, 0, 0, 0},
	Angle:             180,
	AngleVelocity:     0,
	AngleAcceleration: 0,
	//PrevUpdateTime:     time.Now(),
	VelocityX:          0,
	VelocityY:          0,
	VelocityMultiplier: 1.1,
	VelocityX0:         0,
	VelocityY0:         0,
	ForceX:             0,
	ForceY:             0,
}

func MoveShip() {
	// Ускорение убывает от веса
	Ship.ForceX *= math.Pow(constant.Friction, Ship.Weight)
	Ship.ForceY *= math.Pow(constant.Friction, Ship.Weight)
	// Торможение для поворота
	Ship.AngleAcceleration *= math.Pow(constant.Friction, Ship.Weight)
	Ship.AngleVelocity += Ship.AngleAcceleration
	weightMult := 1 + Ship.Weight/100
	Ship.AngleVelocity *= math.Pow(constant.Friction, (weightMult / (1 + Ship.AngleVelocity)))

	moveShipOnReal()
}

func moveShipOnReal() {
	angle := int(180 * Ship.Body.Angle() / math.Pi)
	Ship.Angle = angle
	pos := Ship.Body.Position()
	Ship.X = pos.X
	Ship.Y = pos.Y
	vel := Ship.Body.Velocity()
	Ship.VelocityX = vel.X
	Ship.VelocityY = vel.Y

	Ship.Body.SetAngularVelocity(Ship.AngleVelocity)
	Ship.Body.SetForce(cp.Vector{Ship.ForceX, Ship.ForceY})
}

func OnPressGas() {
	radian := float64(Ship.Angle+90) * math.Pi / 180
	Ship.ForceX += math.Cos(radian) * 1000
	Ship.ForceY += math.Sin(radian) * 1000

	/*if ship.Ship.Weight > 0 {
		radian := (float64(ship.Ship.Angle) - 90) * math.Pi / 180
		ship.Ship.AccelerationX = math.Cos(radian) * ship.Ship.AccelerationMultipler
		ship.Ship.AccelerationY = math.Sin(radian) * ship.Ship.AccelerationMultipler
	}*/
}

func OnPressLeft() {
	Ship.AngleAcceleration -= 0.1
}

func OnPressRight() {
	Ship.AngleAcceleration += 0.1
}
