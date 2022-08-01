package rpg

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/math/f64"
	"image"
	"image/color"
	"log"
)

const (
	GridSize = 16
	TileSize = 16
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
	// w6 h8
	mapImage, _, err := ebitenutil.NewImageFromFile("assets/sprites/tilesets/plains.png")

	if err != nil {
		log.Fatal(err)
	}

	gm.sprite = ebiten.NewImageFromImage(mapImage)

	//gm.Tiles = make([][]Tile, gm.MapSizeY)

	//for i := range gm.Tiles {
	//	gm.Tiles[i] = make([]Tile, gm.MapSizeX)
	//}

	//for y := 0; y < gm.MapSizeY; y++ {
	//	for x := 0; x < gm.MapSizeX; x++ {
	//		gm.Tiles[y][x] = Tile{
	//			Id: 1,
	//			Position: Vector2{
	//				X: float64(x),
	//				Y: float64(y),
	//			},
	//		}
	//	}
	//}
}

func (gm *GameMap) getSpriteCoordsById(id int) (coords f64.Vec4) {
	w, _ := gm.sprite.Size()

	x1 := float64(id % (w / TileSize) * TileSize)
	y1 := float64(id / (w / TileSize) * TileSize)
	x2 := x1 + TileSize
	y2 := y1 + TileSize
	return f64.Vec4{
		x1,
		y1,
		x2,
		y2,
	}

}

func (gm *GameMap) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	for _, row := range gm.Tiles {
		for _, tile := range row {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(tile.Position.X*16, tile.Position.Y*16)
			spriteCoords := gm.getSpriteCoordsById(tile.Id)

			screen.DrawImage(
				gm.sprite.SubImage(
					image.Rect(
						int(spriteCoords[0]),
						int(spriteCoords[1]),
						int(spriteCoords[2]),
						int(spriteCoords[3]),
					),
				).(*ebiten.Image),
				op,
			)
		}
	}
}
