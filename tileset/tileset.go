package tileset

import (
	"image"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wenealves10/game-ebiten-engine/constants"
)

type Tileset interface {
	Img(id int) *ebiten.Image
}

type UniformTilesetJSON struct {
	Path string `json:"path"`
}

type UniformTileset struct {
	img *ebiten.Image
	gid int
}

func (u *UniformTileset) Img(id int) *ebiten.Image {
	id -= u.gid

	srcX := id % 22
	srcY := id / 22

	srcX *= constants.Tilesize
	srcY *= constants.Tilesize

	return u.img.SubImage(image.Rect(srcX, srcY, srcX+constants.Tilesize, srcY+constants.Tilesize)).(*ebiten.Image)
}

type TileJSON struct {
	Id     int    `json:"id"`
	Path   string `json:"image"`
	Width  int    `json:"imagewidth"`
	Height int    `json:"imageheight"`
}

type DynTilesetJSON struct {
	Tiles []*TileJSON `json:"tiles"`
}

type DynTileset struct {
	imgs []*ebiten.Image
	gid  int
}

func (d *DynTileset) Img(id int) *ebiten.Image {
	id -= d.gid
	return d.imgs[id]
}

func NewTileset(path string, gid int) (Tileset, error) {
    contents, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    if strings.Contains(path, "")
}
