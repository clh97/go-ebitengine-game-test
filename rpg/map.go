package rpg

import (
	"encoding/json"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/math/f64"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"
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

type TiledMap struct {
	Compressionlevel int  `json:"compressionlevel"`
	Height           int  `json:"height"`
	Infinite         bool `json:"infinite"`
	Layers           []struct {
		Data    []int  `json:"data"`
		Height  int    `json:"height"`
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Opacity int    `json:"opacity"`
		Type    string `json:"type"`
		Visible bool   `json:"visible"`
		Width   int    `json:"width"`
		X       int    `json:"x"`
		Y       int    `json:"y"`
	} `json:"layers"`
	Nextlayerid  int    `json:"nextlayerid"`
	Nextobjectid int    `json:"nextobjectid"`
	Orientation  string `json:"orientation"`
	Renderorder  string `json:"renderorder"`
	Tiledversion string `json:"tiledversion"`
	Tileheight   int    `json:"tileheight"`
	Tilesets     []struct {
		Columns     int    `json:"columns"`
		Firstgid    int    `json:"firstgid"`
		Image       string `json:"image"`
		Imageheight int    `json:"imageheight"`
		Imagewidth  int    `json:"imagewidth"`
		Margin      int    `json:"margin"`
		Name        string `json:"name"`
		Spacing     int    `json:"spacing"`
		Tilecount   int    `json:"tilecount"`
		Tileheight  int    `json:"tileheight"`
		Tilewidth   int    `json:"tilewidth"`
	} `json:"tilesets"`
	Tilewidth int    `json:"tilewidth"`
	Type      string `json:"type"`
	Version   string `json:"version"`
	Width     int    `json:"width"`
}

func (gm *GameMap) Init() {
	// w6 h8
	mapImage, _, err := ebitenutil.NewImageFromFile("assets/sprites/tilesets/plains.png")

	if err != nil {
		log.Fatal(err)
	}

	gm.sprite = ebiten.NewImageFromImage(mapImage)

	tiles, width, height := gm.readMap()

	gm.MapSizeX = width
	gm.MapSizeY = height
	gm.Tiles = [][]Tile{
		tiles,
	}

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

func (gm *GameMap) readMap() (tiles []Tile, width, height int) {
	jsonFile, err := os.Open("demo_map.json")

	if err != nil {
		log.Fatal("Unable to load map!")
	}

	defer jsonFile.Close()

	jsonBytes, _ := ioutil.ReadAll(jsonFile)

	var result TiledMap

	json.Unmarshal([]byte(jsonBytes), &result)

	width = result.Width
	height = result.Height

	for _, layer := range result.Layers {
		yIndex := 0
		xIndex := 0
		for i, tileId := range layer.Data {
			xIndex = i % width

			fmt.Println(i, width, i%width == 0)

			if i%width == 0 {
				xIndex = 0
				yIndex++
			}

			tiles = append(tiles, Tile{Id: tileId - 1, Position: Vector2{X: float64(xIndex), Y: float64(yIndex - 1)}})
		}
	}

	return tiles, width, height
}
