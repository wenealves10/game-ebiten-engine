package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/wenealves10/game-ebiten-engine/animations"
	"github.com/wenealves10/game-ebiten-engine/camera"
	"github.com/wenealves10/game-ebiten-engine/constants"
	"github.com/wenealves10/game-ebiten-engine/entities"
	"github.com/wenealves10/game-ebiten-engine/spritesheet"
	"github.com/wenealves10/game-ebiten-engine/tilemap"
	"github.com/wenealves10/game-ebiten-engine/tileset"
)

const (
	screenWidth  = 320
	screenHeight = 240
	gravity      = 800.0
	jumpImpulse  = -300.0
)

type Game struct {
	player            *entities.Player
	playerSpriteSheet *spritesheet.SpriteSheet
	tilemapJSON       *tilemap.TilemapJSON
	tilesets          []tileset.Tileset
	cam               *camera.Camera
}

func NewGame() *Game {
	tilemapJSON, err := tilemap.NewTilemapJSON("assets/maps/maps.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created tilemapJSON")

	tilesets, err := tilemapJSON.GetTilesetPath()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created tilemapJSON and tilesets")

	playerImgIdle, _, err := ebitenutil.NewImageFromFile("assets/ninja/Idle32x32.png")
	if err != nil {
		log.Fatal(err)
	}

	playerSpriteSheet := spritesheet.NewSpriteSheet(11, 0, constants.Tilesize*2)

	return &Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImgIdle,
				X:   100,
				Y:   180,
			},
			Health: 100,
			Animations: map[entities.PlayerState]*animations.Animation{
				entities.Idle: animations.NewAnimation(0, 10, 1, 5.0),
			},
		},
		playerSpriteSheet: playerSpriteSheet,
		tilemapJSON:       tilemapJSON,
		tilesets:          tilesets,
		cam:               camera.NewCamera(0.0, 0.0),
	}
}

func (g *Game) Update() error {

	g.player.Dx = 0.0

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.Dx = -2
		g.player.Flip = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.Dx = 2
		g.player.Flip = false
	}

	g.player.X += g.player.Dx

	activeAnim := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnim != nil {
		activeAnim.Update()
	}

	const dt = 1.0 / 60.0
	const groundY = 180

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.player.Y >= groundY {
		g.player.Dy = jumpImpulse
	}

	g.player.Dy += gravity * dt
	g.player.Y += g.player.Dy * dt

	if g.player.Y > groundY {
		g.player.Y = groundY
		g.player.Dy = 0
	}

	g.cam.FollowTarget(g.player.X+8, g.player.Y+8, 320, 240)
	g.cam.Constrain(
		float64(g.tilemapJSON.Layers[0].Width)*constants.Tilesize,
		float64(g.tilemapJSON.Layers[0].Height)*constants.Tilesize,
		320,
		240,
	)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}

	for _, layer := range g.tilemapJSON.Layers {
		for index, id := range layer.Data {
			for _, tileset := range g.tilesets {
				if id == 0 {
					continue
				}

				x := index % layer.Width
				y := index / layer.Width

				x *= constants.Tilesize
				y *= constants.Tilesize

				img := tileset.Img(id)

				opts.GeoM.Translate(float64(x), float64(y))

				opts.GeoM.Translate(0.0, -float64(img.Bounds().Dy())+constants.Tilesize)

				opts.GeoM.Translate(g.cam.X, g.cam.Y)

				screen.DrawImage(img, &opts)

				opts.GeoM.Reset()
			}
		}
	}

	playerFrame := 0
	activeAnim := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnim != nil {
		playerFrame = activeAnim.CurrentFrame()
	}

	cx := float64(g.playerSpriteSheet.Tilesize) / 2
	cy := float64(g.playerSpriteSheet.Tilesize) / 2

	opts.GeoM.Translate(-cx, -cy)

	if g.player.Flip {
		opts.GeoM.Scale(-1, 1)
		opts.GeoM.Translate(0, 0)
	}

	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.cam.X, g.cam.Y)

	// draw the player
	screen.DrawImage(
		// grab a subimage of the spritesheet
		g.player.Img.SubImage(
			g.playerSpriteSheet.Rect(playerFrame),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Tilemap (Ebiten Demo)")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
