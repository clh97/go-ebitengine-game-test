package rpg

import (
	"encoding/json"
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
	Position Position
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
			op.GeoM.Translate(float64(tile.Position.X+1)*16, float64(tile.Position.Y+1)*16)

			if tile.Id == 48 {
				emptyImage := ebiten.NewImage(16, 16)
				op := &ebiten.DrawImageOptions{}
				op.ColorM.ScaleWithColor(color.Black)
				// Filter must be 'nearest' filter (default).
				// Linear filtering would make edges blurred.
				screen.DrawImage(emptyImage.SubImage(image.Rect(tile.Position.X, tile.Position.Y, tile.Position.X*16, tile.Position.Y*16)).(*ebiten.Image), op)
				return
			}

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

func (gm *GameMap) TilesAt(x, y int) (tiles []Tile) {
	for _, layer := range gm.Tiles {
		for _, tile := range layer {
			if tile.Position.X == x && tile.Position.Y == y {
				tiles = append(tiles, tile)
			}
		}
	}
	return tiles
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
		for i, tileId := range layer.Data {
			xIndex := i % width

			if i%width == 0 {
				xIndex = 0
				yIndex++
			}

			tiles = append(tiles, Tile{Id: tileId - 1, Position: Position{X: xIndex, Y: yIndex}})
		}
	}

	return tiles, width, height
}
