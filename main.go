package main

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jakecoffman/cp"
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
	player      *entities.Player
	tilemapJSON *tilemap.TilemapJSON
	tilesets    []tileset.Tileset
	cam         *camera.Camera
	colliders   []image.Rectangle
}

func NewGame() *Game {
	space := cp.NewSpace()
	space.SetGravity(cp.Vector{X: 0, Y: 800})
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

	playerImgRun, _, err := ebitenutil.NewImageFromFile("assets/ninja/Run32x32.png")
	if err != nil {
		log.Fatal(err)
	}

	playerImgJump, _, err := ebitenutil.NewImageFromFile("assets/ninja/Jump32x32.png")
	if err != nil {
		log.Fatal(err)
	}

	animationIdle := spritesheet.NewSpriteSheet(11, 0, constants.Tilesize*2)
	animationRun := spritesheet.NewSpriteSheet(12, 0, constants.Tilesize*2)
	animationJump := spritesheet.NewSpriteSheet(1, 0, constants.Tilesize*2)

	return &Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				X: 30,
				Y: 180,
				W: 32,
				H: 32,
			},
			Health: 100,
			Animations: map[entities.PlayerState]*animations.Animation{
				entities.Idle:    animations.NewAnimation(0, 10, 1, 5.0, animationIdle, playerImgIdle),
				entities.Running: animations.NewAnimation(0, 11, 1, 5.0, animationRun, playerImgRun),
				entities.Jumping: animations.NewAnimation(0, 0, 1, 5.0, animationJump, playerImgJump),
			},
		},
		tilemapJSON: tilemapJSON,
		tilesets:    tilesets,
		cam:         camera.NewCamera(0.0, 0.0),
		colliders:   tilemapJSON.GetColliders(),
	}
}

func CheckCollisionHorizontal(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(
				int(sprite.X),
				int(sprite.Y),
				int(sprite.X)+constants.Tilesize,
				int(sprite.Y)+constants.Tilesize,
			),
		) {
			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - constants.Tilesize
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

func CheckCollisionVertical(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(
				int(sprite.X),
				int(sprite.Y),
				int(sprite.X)+constants.Tilesize,
				int(sprite.Y)+constants.Tilesize,
			),
		) {
			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - constants.Tilesize
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y)
			}
		}
	}
}

func (g *Game) Update() error {

	movingHorizontal := false

	g.player.Dx = 0.0

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.Dx = -2
		g.player.Flip = true
		movingHorizontal = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.Dx = 2
		g.player.Flip = false
		movingHorizontal = true
	}

	g.player.X += g.player.Dx

	CheckCollisionHorizontal(g.player.Sprite, g.colliders)

	activeAnim := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnim != nil {
		activeAnim.Update()
	}

	const dt = 1.0 / 60.0

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.player.Dy = jumpImpulse
		g.player.State = "jumping"
	}

	g.player.Dy += gravity * dt
	g.player.Y += g.player.Dy * dt

	if g.player.Dy > 300 {
		g.player.Dy = 300
		if movingHorizontal {
			g.player.State = "running"
		} else {
			g.player.State = "idle"
		}
	} else {
		g.player.State = "jumping"
	}

	CheckCollisionVertical(g.player.Sprite, g.colliders)

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

	cx := float64(activeAnim.Spritesheet.Tilesize) / 2
	cy := float64(activeAnim.Spritesheet.Tilesize) / 2

	opts.GeoM.Translate(-cx, -cy)

	if g.player.Flip {
		opts.GeoM.Scale(-1, 1)
		opts.GeoM.Translate(0, 0)
	}

	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.cam.X, g.cam.Y)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("DX: %f DY: %f", g.player.Dx, g.player.Dy), 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("X: %f Y: %f", g.player.X, g.player.Y), 0, 20)

	// draw the player
	screen.DrawImage(
		// grab a subimage of the spritesheet
		activeAnim.Img.SubImage(
			activeAnim.Spritesheet.Rect(playerFrame),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

	// for _, collider := range g.colliders {
	// 	vector.StrokeRect(
	// 		screen,
	// 		float32(collider.Min.X)+float32(g.cam.X),
	// 		float32(collider.Min.Y)+float32(g.cam.Y),
	// 		float32(collider.Dx()),
	// 		float32(collider.Dy()),
	// 		1.0,
	// 		color.RGBA{255, 0, 0, 255},
	// 		true,
	// 	)
	// }
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
