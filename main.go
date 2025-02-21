package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wenealves10/game-ebiten-engine/camera"
	"github.com/wenealves10/game-ebiten-engine/constants"
	"github.com/wenealves10/game-ebiten-engine/tilemap"
	"github.com/wenealves10/game-ebiten-engine/tileset"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

type Game struct {
	tilemapJSON *tilemap.TilemapJSON
	tilesets    []tileset.Tileset
	cam         *camera.Camera
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

	return &Game{
		tilemapJSON: tilemapJSON,
		tilesets:    tilesets,
		cam:         camera.NewCamera(0.0, 0.0),
	}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// screen.Fill(color.RGBA{120, 180, 255, 255})

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
