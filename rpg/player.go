package rpg

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
)

const (
	playerSize = 16

	// Player animation constants
	frameOX     = 0
	frameOY     = 0
	frameWidth  = 48
	frameHeight = 48
	frameCount  = 6
)

type Player struct {
	Position Vector2

	sprite *ebiten.Image
}

func (p *Player) Init() {
	playerImage, _, err := ebitenutil.NewImageFromFile("assets/sprites/characters/player.png")

	if err != nil {
		log.Fatal(err)
	}

	p.sprite = ebiten.NewImageFromImage(playerImage)
}

func (p *Player) String() string {
	return fmt.Sprintf(
		"P: %.1f, Y: %.1f",
		p.Position.X, p.Position.Y,
	)
}

func (p *Player) Draw(screen *ebiten.Image, timer int) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(p.Position.X)*GridSize-(playerSize/2), float64(p.Position.Y)*GridSize-playerSize)

	i := (timer / 6) % frameCount
	sx, sy := frameOX+i*frameWidth, frameOY

	screen.DrawImage(p.sprite.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}
