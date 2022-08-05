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
	Position Position

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
		"P: %d, Y: %d",
		p.Position.X, p.Position.Y,
	)
}

func (p *Player) Move(x, y int) {
	nextPosX := p.Position.X + x
	nextPosY := p.Position.Y + y

	p.Position.X = nextPosX
	p.Position.Y = nextPosY
}

func (p *Player) Draw(screen *ebiten.Image, timer int) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(p.Position.X)*GridSize, float64(p.Position.Y)*GridSize)

	i := (timer / 6) % frameCount
	sx, sy := frameOX+i*frameWidth, frameOY

	screen.DrawImage(p.sprite.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}
