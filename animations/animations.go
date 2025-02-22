package animations

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wenealves10/game-ebiten-engine/spritesheet"
)

type Animation struct {
	First        int
	Last         int
	Step         int
	SpeedInTps   float32
	frameCounter float32
	frame        int
	Spritesheet  *spritesheet.SpriteSheet
	Img          *ebiten.Image
}

func (a *Animation) Update() {
	a.frameCounter -= 1.0
	if a.frameCounter < 0.0 {
		a.frameCounter = a.SpeedInTps
		a.frame += a.Step
		if a.frame > a.Last {
			a.frame = a.First
		}
	}
}

func (a *Animation) CurrentFrame() int {
	return a.frame
}

func NewAnimation(first, last, step int, speed float32, spriteSheet *spritesheet.SpriteSheet, img *ebiten.Image) *Animation {
	return &Animation{
		First:        first,
		Last:         last,
		Step:         step,
		SpeedInTps:   speed,
		frameCounter: speed,
		frame:        first,
		Spritesheet:  spriteSheet,
		Img:          img,
	}
}
