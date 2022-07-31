package rpg

type Tile struct {
	Id       int
	Position Vector2
}

type GameMap struct {
	Tiles    [][]Tile
	MapSizeX int
	MapSizeY int
}

func (gm *GameMap) Init() {
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
