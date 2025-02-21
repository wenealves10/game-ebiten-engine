package tileset

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/wenealves10/game-ebiten-engine/constants"
)

type Tileset interface {
	Img(id int) *ebiten.Image
}

type UniformTilesetJSON struct {
	Path string `json:"image"`
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

	if strings.Contains(path, "buildings") {
		var dynTilesetJSON DynTilesetJSON
		err = json.Unmarshal(contents, &dynTilesetJSON)
		if err != nil {
			return nil, err
		}

		dynTileset := &DynTileset{}
		dynTileset.gid = gid
		dynTileset.imgs = make([]*ebiten.Image, 0)

		for _, tileJSON := range dynTilesetJSON.Tiles {
			tileJSONPath := tileJSON.Path
			tileJSONPath = filepath.Clean(tileJSONPath)
			tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
			tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
			tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
			tileJSONPath = filepath.Join("assets/", tileJSONPath)

			img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
			if err != nil {
				return nil, err
			}

			dynTileset.imgs = append(dynTileset.imgs, img)
		}

		return dynTileset, nil
	}

	var uniformTilesetJSON UniformTilesetJSON
	err = json.Unmarshal(contents, &uniformTilesetJSON)
	if err != nil {
		return nil, err
	}

	uniformTileset := UniformTileset{}

	uniformTilesetJSONPath := uniformTilesetJSON.Path
	fmt.Println(uniformTilesetJSONPath)
	uniformTilesetJSONPath = filepath.Clean(uniformTilesetJSONPath)
	uniformTilesetJSONPath = strings.ReplaceAll(uniformTilesetJSONPath, "\\", "/")
	uniformTilesetJSONPath = strings.TrimPrefix(uniformTilesetJSONPath, "../")
	uniformTilesetJSONPath = strings.TrimPrefix(uniformTilesetJSONPath, "../")
	uniformTilesetJSONPath = filepath.Join("assets/", uniformTilesetJSONPath)

	fmt.Println(uniformTilesetJSONPath)

	img, _, err := ebitenutil.NewImageFromFile(uniformTilesetJSONPath)
	if err != nil {
		return nil, err
	}

	uniformTileset.img = img
	uniformTileset.gid = gid

	return &uniformTileset, nil
}
