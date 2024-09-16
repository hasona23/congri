package utils

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type AnimSprite struct {
	Sprite
	SpriteSheet   *ebiten.Image
	width, height int
	Animations    map[string]AnimationFrame //animations
	Current       AnimationFrame
}

// return an animasprite
func NewAnimeSprite(img *ebiten.Image, Width, Height int) *AnimSprite {
	return &AnimSprite{Sprite: *NewSprite(ebiten.NewImage(16, 16)), SpriteSheet: img, width: Width, height: Height, Animations: make(map[string]AnimationFrame)}
}
func (sprite AnimSprite) GetWidth() int {
	return sprite.Sprite.GetWidth()
}
func (sprite AnimSprite) GetHeight() int {
	return sprite.Sprite.GetHeight()
}
func (sprite AnimSprite) GetImg() *ebiten.Image {
	return sprite.Img
}
func (a AnimSprite) SetSpriteOP(op *ebiten.DrawImageOptions, rotation float32) {
	a.Sprite.SetSpriteOP(op, rotation)
}
func (a *AnimSprite) Add(animtion AnimationFrame) {
	if len(a.Animations) == 0 {
		a.Animations[animtion.Name] = animtion
		a.Current = a.Animations[animtion.Name]
	} else {
		a.Animations[animtion.Name] = animtion
	}

}

func (a *AnimSprite) playAnim() {
	if _, ok := a.Animations[a.Current.Name]; ok {
		a.Current.Update()
	} else {
		panic("AnimeSprite : no animation with name: " + a.Current.Name)
	}
	a.Sprite.Img = a.SpriteSheet.SubImage(image.Rect(a.Current.row_current*a.width,
		a.Current.col_current*a.height,
		(a.Current.row_current+1)*a.width,
		(a.Current.col_current+1)*a.height,
	)).(*ebiten.Image)
}
func (a *AnimSprite) ChangeAnim(name string) {
	if a.Current.Name != name {
		a.Current = a.Animations[name]
	}
}
func (a *AnimSprite) Animate() {
	a.playAnim()
}

// AnimeSprite Consist of many animation Frames to Player
type AnimationFrame struct {
	row_min     int
	row_max     int
	row_current int
	col_min     int
	col_max     int
	col_current int
	IsEnd       bool
	Name        string
	Timer
}

func (a AnimationFrame) IsEmpty() bool {
	return AnimationFrame{} != a
}

func NewAnimationFrame(row_min, row_max, col_min, col_max int, duration float32, name string) AnimationFrame {
	return AnimationFrame{row_min: row_min, row_max: row_max, row_current: row_min, col_min: col_min, col_max: col_max, col_current: col_min, IsEnd: false, Name: name, Timer: NewTimer(duration)}
}

// Updates Column and row number to affect the animsprite SubImage
func (a *AnimationFrame) Update() {
	a.Timer.UpdateTimer()
	//fmt.Println(a.Timer.GetCurrentTime())
	if a.Ticked() {

		a.IsEnd = false
		a.row_current++
		//fmt.Println(" Row", a.row_max, " ", a.row_current)
		if a.row_current >= a.row_max {

			a.row_current = a.row_min
			a.col_current++

		}
		//fmt.Println(" Col ", a.col_max, " ", a.col_current)
		if a.col_current >= a.col_max {

			a.col_current = a.col_min
			a.IsEnd = true

		}

	}
}
