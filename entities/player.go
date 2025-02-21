package entities

import "github.com/wenealves10/game-ebiten-engine/animations"

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
}

func (p *Player) ActiveAnimation(dx, dy int) *animations.Animation {
	if dx == 0 && dy == 0 {
		return p.Animations[Idle]
	}

	if dx > 0 {
		return p.Animations[Running]
	}

	if dy > 0 {
		return p.Animations[Jumping]
	}

	return nil
}
