package rpg

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"log"
)

const (
	GridSize = 16
)

type Tile struct {
	Id       int
	Position Vector2
}

type GameMap struct {
	Tiles    [][]Tile
	MapSizeX int
	MapSizeY int

	sprite *ebiten.Image
}

func (gm *GameMap) Init() {
	mapImage, _, err := ebitenutil.NewImageFromFile("assets/sprites/tilesets/grass.png")

	if err != nil {
		log.Fatal(err)
	}

	gm.sprite = ebiten.NewImageFromImage(mapImage)

	gm.Tiles = make([][]Tile, gm.MapSizeY)

	for i := range gm.Tiles {
		gm.Tiles[i] = make([]Tile, gm.MapSizeX)
	}

	for y := 0; y < gm.MapSizeY; y++ {
		for x := 0; x < gm.MapSizeX; x++ {
			gm.Tiles[y][x] = Tile{
				Id: 1,
				Position: Vector2{
					X: float64(x),
					Y: float64(y),
				},
			}
		}
	}
}

func (gm *GameMap) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0xff})
	for _, row := range gm.Tiles {
		for _, tile := range row {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(tile.Position.X*16, tile.Position.Y*16)
			screen.DrawImage(gm.sprite, op)
		}
	}
}
