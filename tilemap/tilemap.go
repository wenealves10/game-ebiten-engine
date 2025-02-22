package tilemap

import (
	"encoding/json"
	"image"
	"os"
	"path"

	"github.com/wenealves10/game-ebiten-engine/tileset"
)

type Object struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Visible  bool    `json:"visible"`
	Rotation float64 `json:"rotation"`
	Height   float64 `json:"height"`
	Width    float64 `json:"width"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
}
type TilemapLayerJSON struct {
	Data    []int    `json:"data"`
	Width   int      `json:"width"`
	Height  int      `json:"height"`
	Name    string   `json:"name"`
	Objects []Object `json:"objects"`
}

type TilemapJSON struct {
	Layers   []*TilemapLayerJSON `json:"layers"`
	Tilesets []map[string]any    `json:"tilesets"`
}

func (t *TilemapJSON) GetColliders() []image.Rectangle {
	var colliders []image.Rectangle
	for _, layer := range t.Layers {
		for _, object := range layer.Objects {
			colliders = append(colliders, image.Rect(
				int(object.X),
				int(object.Y),
				int(object.X)+int(object.Width),
				int(object.Y)+int(object.Height),
			))
		}
	}
	return colliders
}

func (t *TilemapJSON) GetTilesetPath() ([]tileset.Tileset, error) {
	tilesets := make([]tileset.Tileset, 0)

	for _, tilesetData := range t.Tilesets {
		tilesetPath := path.Join("assets/maps", tilesetData["source"].(string))
		tileset, err := tileset.NewTileset(tilesetPath, int(tilesetData["firstgid"].(float64)))
		if err != nil {
			return nil, err
		}
		tilesets = append(tilesets, tileset)
	}

	return tilesets, nil
}

func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tilemapJSON TilemapJSON
	err = json.Unmarshal(contents, &tilemapJSON)
	if err != nil {
		return nil, err
	}

	return &tilemapJSON, nil
}
