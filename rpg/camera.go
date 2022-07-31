package rpg

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type Camera struct {
	ViewPort   Vector2
	Position   Vector2
	ZoomFactor int
	Rotation   int
}

func (c *Camera) String() string {
	return fmt.Sprintf(
		"T: %.1f, R: %d, S: %d",
		c.Position, c.Rotation, c.ZoomFactor,
	)
}

func (c *Camera) viewportCenter() Vector2 {
	return Vector2{
		X: c.ViewPort.X * 0.5,
		Y: c.ViewPort.Y * 0.5,
	}
}

func (c *Camera) worldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.Position.X, -c.Position.Y)
	// We want to scale and rotate around center of image / screen
	m.Translate(-c.viewportCenter().X, -c.viewportCenter().Y)
	m.Scale(
		math.Pow(1.01, float64(c.ZoomFactor)),
		math.Pow(1.01, float64(c.ZoomFactor)),
	)
	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)
	m.Translate(c.viewportCenter().X, c.viewportCenter().Y)
	return m
}

func (c *Camera) Render(world, screen *ebiten.Image) {
	screen.DrawImage(world, &ebiten.DrawImageOptions{
		GeoM: c.worldMatrix(),
	})
}

func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
	inverseMatrix := c.worldMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		return inverseMatrix.Apply(float64(posX), float64(posY))
	} else {
		return math.NaN(), math.NaN()
	}
}

func (c *Camera) Reset() {
	c.Position.X = 0
	c.Position.Y = 0
	c.Rotation = 0
	c.ZoomFactor = 0
}
