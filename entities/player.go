package entities

import (
	"github.com/wenealves10/game-ebiten-engine/animations"
)

type PlayerState uint8

const (
	Idle PlayerState = iota
	Running
	Jumping
	Hitting
	DoubleJumping
	WallJumping
	Falling
)

type Player struct {
	*Sprite
	Health     uint
	Animations map[PlayerState]*animations.Animation
	Flip       bool
	State      string // "idle", "running", "jumping"
}

func (p *Player) ActiveAnimation(dx, dy int) *animations.Animation {
	if p.State == "jumping" {
		return p.Animations[Jumping]
	}

	if p.State == "running" {
		return p.Animations[Running]
	}

	return p.Animations[Idle]
}
