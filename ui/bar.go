package ui

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hasona23/game/utils"
)

// bar ui with values such as the health bar
type Bar struct {
	img   *ebiten.Image //note: try avoid having spaces in your sprite as spaces count in the width
	rect  utils.Rect
	max   int
	value int
	BarStyle
}

// controls visual aspects of bar
type BarStyle struct {
	BarColor  color.Color //color of bar itself like red in hp bar
	BackColor color.Color //color of rect behind bar like black in most bars
	Scale     utils.Point
}

// return a new BarStyle
func NewBarStyle(barColor, BackColor color.Color, scale utils.Point) BarStyle {
	return BarStyle{BarColor: barColor, BackColor: BackColor, Scale: scale}
}

// returns a new bar with an img(sprite)
func NewSpriteBar(img *ebiten.Image, x, y, value int, scale utils.Point) *Bar {
	return &Bar{img: img,
		rect:     utils.NewRect(x, y, img.Bounds().Dx(), img.Bounds().Dy()),
		value:    value,
		max:      value,
		BarStyle: NewBarStyle(color.Black, color.Black, scale),
	}
}

// returns if bar is full
func (b *Bar) IsFull() bool {
	return b.value == b.max
}

// returns if bar is empty
func (b *Bar) IsEmpty() bool {
	return b.value == 0
}

// Value is same as the width  no differnce . this returns a rectangular bar need incase your dont want sprites
// like  for testing and demo
func NewBar(x, y, value, height int, scale utils.Point, barColor, backColor color.Color) *Bar {
	return &Bar{img: nil,
		rect:     utils.NewRect(x, y, value*scale.X, height*scale.Y),
		value:    value,
		max:      value,
		BarStyle: NewBarStyle(barColor, backColor, scale),
	}
}

// this is ration between current value and max value which is used for drawing
func (b *Bar) GetRatio() float64 {
	return float64(b.value) / float64(b.max)
}

// changes position of bar
func (b *Bar) SetPos(x, y int) {
	b.rect.X = x
	b.rect.Y = y
}

// returns position of bar
func (b *Bar) GetPos() (int, int) {
	return b.rect.X, b.rect.Y
}

// returns maximum value of bar
func (b *Bar) GetMax() int {
	return b.max
}

// changes the maximum value of bar . if n <1  then max =1
func (b *Bar) SetMax(n int) {
	b.max = int(math.Max(1, float64(n)))
}

// changes current value of bar. if value is more than max it is set to max and if less than zero is set to zero
func (b *Bar) SetValue(n int) {
	b.value = int(math.Max(0, math.Min(float64(n), float64(b.max))))
}

// sets both maximum and current value both at same time
func (b *Bar) SetValueAndMax(max, value int) {
	b.SetMax(max)
	b.SetValue(value)
}

// returns the current value of bar
func (b *Bar) GetValue() int {
	return b.value
}

// changes the background color
func (b *Bar) SetBackColor(c color.Color) {
	b.BackColor = c
}

// chages the bar color
func (b *Bar) SetBarColor(c color.Color) {
	b.BarColor = c
}

// Draws the bar in a fixed position
// good for player ui or ui always present on screen
// note: img is rectangle bar if img is nil
func (b *Bar) Draw(screen *ebiten.Image) {
	if b.img == nil {
		b.drawRectBar(screen, utils.Point{X: 0, Y: 0})
	} else {
		b.drawImageBar(screen, utils.Point{X: 0, Y: 0})
	}
}

// draw bar in a dynamic pos
// suitable for moving object like enemies or want to make bar has certain pos like hp of and object like rock
// note: img is rectangle bar if img is nil
func (b *Bar) DrawCam(screen *ebiten.Image, cam utils.Cam) {
	if b.img == nil {
		b.drawRectBar(screen, utils.Point{X: int(cam.X), Y: int(cam.Y)})
	} else {
		b.drawImageBar(screen, utils.Point{X: int(cam.X), Y: int(cam.Y)})
	}
}

// draws image of bar as rect .
// used when bar img is nil
func (b *Bar) drawRectBar(screen *ebiten.Image, offset utils.Point) {
	x := float32(b.rect.X/b.Scale.X + offset.X)
	y := float32(b.rect.Y/b.Scale.Y + offset.Y)
	width := float32(b.rect.Width)
	height := float32(b.rect.Height)

	vector.DrawFilledRect(screen, x, y, width, height, b.BackColor, false)
	vector.DrawFilledRect(screen, x, y, width*float32(b.GetRatio()), height, b.BarColor, false)
}

// draws the image of bar
// used when img not nil
func (b *Bar) drawImageBar(screen *ebiten.Image, offset utils.Point) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(b.Scale.X), float64(b.Scale.Y))
	op.GeoM.Translate(float64(b.rect.X/b.Scale.X)+float64(offset.X), float64(b.rect.Y/b.Scale.Y)+float64(offset.Y))

	op.ColorScale.ScaleWithColor(b.BackColor)
	screen.DrawImage(b.img, op)

	op.ColorScale.Reset()
	subImg := b.img.SubImage(image.Rect(
		b.img.Bounds().Min.X,
		b.img.Bounds().Min.Y,
		b.img.Bounds().Min.X+int(float64(b.img.Bounds().Dx())*b.GetRatio()),
		b.img.Bounds().Max.Y,
	)).(*ebiten.Image)
	screen.DrawImage(subImg, op)
}
