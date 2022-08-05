package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "image/png"
	"log"
	"os"
	rpg "rpg/rpg"
)

const (
	screenWidth  = 480
	screenHeight = 320

	gridSize   = 16
	worldSizeX = screenWidth / gridSize
	worldSizeY = screenHeight / gridSize
)

const (
	dirNone = iota
	dirLeft
	dirRight
	dirDown
	dirUp
)

var (
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
	g.timer %= 240

	g.handleInput()

	if g.needsToMovePlayer() {
		switch g.moveDirection {
		case dirLeft:
			if !g.collidesWithWall() {
				g.player.Move(-1, 0)
				g.camera.Position.X -= gridSize
			}
		case dirRight:
			if !g.collidesWithWall() {
				g.player.Move(1, 0)
				g.camera.Position.X += gridSize
			}
		case dirUp:
			if !g.collidesWithWall() {
				g.player.Move(0, -1)
				g.camera.Position.Y -= gridSize
			}
		case dirDown:
			if !g.collidesWithWall() {
				g.player.Move(0, 1)
				g.camera.Position.Y += gridSize
			}
		}
	}

	return nil
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
	switch g.moveDirection {
	//case dirLeft:
	//
	//case dirRight:

	case dirUp:
		nextTiles := g.gameMap.TilesAt(g.player.Position.X, g.player.Position.Y)
		for _, tile := range nextTiles {
			if tile.Id == 48 {
				return true
			}
		}
		return false
	case dirDown:
		nextTiles := g.gameMap.TilesAt(g.player.Position.X, g.player.Position.Y+1)
		for _, tile := range nextTiles {
			if tile.Id == 48 {
				return true
			}
		}
		return false
	}
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
	g.gameMap.Draw(g.world)
	g.player.Draw(g.world, g.timer)

	g.camera.Render(g.world, screen)

	g.drawDebugInfo(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)

	ebiten.SetWindowTitle("RPG")
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
	ebiten.SetMaxTPS(60)

	player := rpg.Player{
		Position: rpg.Position{
			X: (worldSizeX / 2) - 1,
			Y: (worldSizeY / 2) - 1,
		},
	}

	gameMap := rpg.GameMap{
		//MapSizeX: 10,
		//MapSizeY: 10,
		//Tiles: [][]rpg.Tile{
		//	{
		//		rpg.Tile{Id: 1, Position: rpg.Vector2{X: 1.0, Y: 0.0}}, rpg.Tile{Id: 2, Position: rpg.Vector2{X: 2.0, Y: 0.0}}, rpg.Tile{Id: 3, Position: rpg.Vector2{X: 3.0, Y: 0.0}},
		//		rpg.Tile{Id: 7, Position: rpg.Vector2{X: 1.0, Y: 1.0}}, rpg.Tile{Id: 8, Position: rpg.Vector2{X: 2.0, Y: 1.0}}, rpg.Tile{Id: 9, Position: rpg.Vector2{X: 3.0, Y: 1.0}},
		//		rpg.Tile{Id: 7, Position: rpg.Vector2{X: 1.0, Y: 2.0}}, rpg.Tile{Id: 8, Position: rpg.Vector2{X: 2.0, Y: 2.0}}, rpg.Tile{Id: 9, Position: rpg.Vector2{X: 3.0, Y: 2.0}},
		//		rpg.Tile{Id: 7, Position: rpg.Vector2{X: 1.0, Y: 3.0}}, rpg.Tile{Id: 8, Position: rpg.Vector2{X: 2.0, Y: 3.0}}, rpg.Tile{Id: 9, Position: rpg.Vector2{X: 3.0, Y: 3.0}},
		//		rpg.Tile{Id: 7, Position: rpg.Vector2{X: 1.0, Y: 4.0}}, rpg.Tile{Id: 8, Position: rpg.Vector2{X: 2.0, Y: 4.0}}, rpg.Tile{Id: 9, Position: rpg.Vector2{X: 3.0, Y: 4.0}},
		//		rpg.Tile{Id: 7, Position: rpg.Vector2{X: 1.0, Y: 5.0}}, rpg.Tile{Id: 8, Position: rpg.Vector2{X: 2.0, Y: 5.0}}, rpg.Tile{Id: 9, Position: rpg.Vector2{X: 3.0, Y: 5.0}},
		//		rpg.Tile{Id: 13, Position: rpg.Vector2{X: 1.0, Y: 6.0}}, rpg.Tile{Id: 14, Position: rpg.Vector2{X: 2.0, Y: 6.0}}, rpg.Tile{Id: 15, Position: rpg.Vector2{X: 3.0, Y: 6.0}},
		//	},
		//},
	}

	camera := rpg.Camera{
		ViewPort: rpg.Vector2{X: screenWidth, Y: screenHeight},
	}

	player.Init()
	gameMap.Init()

	game := &Game{
		moveTime: 5,
		player:   player,
		gameMap:  gameMap,
		camera:   camera,
	}

	game.world = ebiten.NewImage(gameMap.MapSizeX*16, gameMap.MapSizeY*16)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
