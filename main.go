package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"
	rpg "rpg/rpg"
)

const (
	screenWidth  = 480
	screenHeight = 320

	gridSize   = 20
	worldSizeX = screenWidth / gridSize
	worldSizeY = screenHeight / gridSize

	frameOX     = 0
	frameOY     = 0
	frameWidth  = 48
	frameHeight = 48
	frameCount  = 6
)

const (
	dirNone = iota
	dirLeft
	dirRight
	dirDown
	dirUp
)

var (
	playerSprite  *ebiten.Image
	gameMapSprite *ebiten.Image
)

type Game struct {
	timer         int
	moveTime      int
	gridSize      int
	moveDirection int

	keys  []ebiten.Key
	world *ebiten.Image

	player  rpg.Player
	gameMap rpg.GameMap
	camera  rpg.Camera
}

func (g *Game) Update() error {
	g.timer++

	g.handleInput()

	if g.needsToMovePlayer() {
		if g.collidesWithWall() {
			fmt.Println("Colliding with wall")
		}

		switch g.moveDirection {
		case dirLeft:
			g.player.Position.X--
			g.camera.Position.X -= gridSize
		case dirRight:
			g.player.Position.X++
			g.camera.Position.X += gridSize
		case dirUp:
			g.player.Position.Y--
			g.camera.Position.Y -= gridSize
		case dirDown:
			g.player.Position.Y++
			g.camera.Position.Y += gridSize
		}
	}

	return nil
}

func (g *Game) drawMap(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0xff})
	for _, row := range g.gameMap.Tiles {
		for _, tile := range row {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(tile.Position.X*16, tile.Position.Y*16)
			screen.DrawImage(gameMapSprite, op)
		}
	}
}

func (g *Game) drawPlayer(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(g.player.Position.X)*gridSize, float64(g.player.Position.Y)*gridSize)

	i := (g.timer / 6) % frameCount
	sx, sy := frameOX+i*frameWidth, frameOY

	screen.DrawImage(playerSprite.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}

func (g *Game) drawDebugInfo(screen *ebiten.Image) {
	worldX, worldY := g.camera.ScreenToWorld(ebiten.CursorPosition())
	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("FPS: %0.2f\nMove (WASD/Arrows)\nReset (Space)", ebiten.CurrentFPS()),
	)

	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("%s\n%s\nCursor World Pos: %.2f,%.2f",
			g.camera.String(),
			g.player.String(),
			worldX, worldY),
		0, screenHeight-48,
	)
}

func (g *Game) needsToMovePlayer() bool {
	return g.timer%g.moveTime == 0
}

func (g *Game) collidesWithWall() bool {
	return false
}

func (g *Game) handleInput() {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		if g.moveDirection != dirRight {
			g.moveDirection = dirLeft
		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeyA) {
		g.moveDirection = dirNone
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if g.moveDirection != dirLeft {
			g.moveDirection = dirRight
		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeyD) {
		g.moveDirection = dirNone
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if g.moveDirection != dirUp {
			g.moveDirection = dirDown

		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeyS) {
		g.moveDirection = dirNone
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		if g.moveDirection != dirDown {
			g.moveDirection = dirUp
		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeyW) {
		g.moveDirection = dirNone
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.exit()
	}
}

func (g *Game) exit() {
	os.Exit(0)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawMap(g.world)
	g.drawPlayer(g.world)

	g.camera.Render(g.world, screen)

	g.drawDebugInfo(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	playerImage, _, err := ebitenutil.NewImageFromFile("assets/sprites/characters/player.png")
	mapImage, _, err := ebitenutil.NewImageFromFile("assets/sprites/tilesets/grass.png")

	if err != nil {
		log.Fatal(err)
	}

	playerSprite = ebiten.NewImageFromImage(playerImage)
	gameMapSprite = ebiten.NewImageFromImage(mapImage)

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)

	ebiten.SetWindowTitle("RPG")

	player := rpg.Player{
		Position: rpg.Vector2{
			X: 10,
			Y: 7,
		},
	}

	gameMap := rpg.GameMap{
		MapSizeX: 60,
		MapSizeY: 40,
		//Tiles: [][]rpg.Tile{
		//	{
		//		rpg.Tile{Id: 1}, rpg.Tile{Id: 1}, rpg.Tile{Id: 1},
		//		rpg.Tile{Id: 1}, rpg.Tile{Id: 1}, rpg.Tile{Id: 1},
		//		rpg.Tile{Id: 1}, rpg.Tile{Id: 1}, rpg.Tile{Id: 1},
		//	},
		//},
	}

	gameMap.Init()

	game := &Game{
		moveTime: 5,
		player:   player,
		gameMap:  gameMap,
	}

	game.world = ebiten.NewImage(gameMap.MapSizeX*16, gameMap.MapSizeY*16)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
