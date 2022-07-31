package rpg

import "fmt"

type Player struct {
	Position Vector2
}

func (p *Player) String() string {
	return fmt.Sprintf(
		"P: %.1f, Y: %.1f",
		p.Position.X, p.Position.Y,
	)
}
