package utils

import "github.com/hajimehoshi/ebiten/v2"

type SpriteEffect uint // Flip horizontal and vertical

const (
	None SpriteEffect = iota
	FlipHorizontal
	FlipVertical
)

// an interface for sprite and animated sprite
type Spriter interface {
	GetImg() *ebiten.Image
	//Just makes the flipping and rotation
	//
	//Happens before Drawing and Translation
	SetSpriteOP(op *ebiten.DrawImageOptions, rotation float32)
	GetHeight() int
	GetWidth() int
}
