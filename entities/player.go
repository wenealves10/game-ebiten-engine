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
	Health     uint
	Animations map[PlayerState]*animations.Animation
}

func (p *Player) ActiveAnimation(dx, dy float64) *animations.Animation {
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
