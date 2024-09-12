package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Img    *ebiten.Image
	Effect SpriteEffect //effect is for sprite sprite effect (flip horizontal and vertical)
}

func (sprite Sprite) GetWidth() int {
	return sprite.Img.Bounds().Dx()
}
func (sprite Sprite) GetHeight() int {
	return sprite.Img.Bounds().Dy()
}
func (sprite Sprite) GetImg() *ebiten.Image {
	return sprite.Img
}
func (sprite Sprite) FlipHorizontal(op *ebiten.DrawImageOptions) {
	s := sprite.GetWidth()
	op.GeoM.Scale(-1, 1)
	op.GeoM.Translate(float64(s), 0)
}
func (sprite Sprite) FlipVertical(op *ebiten.DrawImageOptions) {
	s := sprite.GetHeight()
	op.GeoM.Scale(1, -1)
	op.GeoM.Translate(0, float64(s))

}

// angle in radians
func (sprite Sprite) Rotate(op *ebiten.DrawImageOptions, angle float64) {
	s := sprite.Img.Bounds()
	op.GeoM.Translate(-float64(s.Dx())/2, -float64(s.Dy())/2)
	op.GeoM.Rotate(angle)
	//op.GeoM.Translate(float64(s.Dx())/2, float64(s.Dy())/2)
}
func NewSprite(img *ebiten.Image) *Sprite {
	if img == nil {
		panic("Sprite : img is nil")
	}
	return &Sprite{img, None}
}
func (s Sprite) SetSpriteOP(op *ebiten.DrawImageOptions, rotation float32) {

	switch s.Effect {
	case FlipHorizontal:
		s.FlipHorizontal(op)
	case FlipVertical:
		s.FlipVertical(op)
	}
	if rotation != 0 {
		s.Rotate(op, float64(rotation))
	}
}
