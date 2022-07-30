package rpg

import "fmt"

type Tile struct {
	Id   int
	PosX int
	PosY int
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
				Id:   1,
				PosX: x,
				PosY: y,
			}
		}
	}

	fmt.Println(gm.Tiles)
}
