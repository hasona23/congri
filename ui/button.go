package ui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hasona23/game/utils"
)

// text orientation in button
type Orientation int

const (
	TopLeft Orientation = iota
	Middle
	BottomRight
)
const (
	DefaultBorderThickness = 1
	TextPaddingFactor      = 0.25
	TextBottomRightFactor  = 1.25
)

type Event func(b *Button)
type Button struct {
	sprite      *ebiten.Image
	Text        Label
	rect        utils.Rect
	clickEvents []Event
	hoverEvents []Event
	Style       ButtonStyle
}

func NewSpriteButton(sprite *ebiten.Image, text string, x, y float32, fontSize, scale int, fontFile []byte, textColor color.Color) *Button {

	button := &Button{
		Text:   *NewLabel(text, x, y, fontFile, fontSize, textColor),
		sprite: sprite,
		Style: ButtonStyle{
			Pos:   utils.Vec2{X: x, Y: y},
			Scale: scale,
			Color: color.Transparent,
		},
		clickEvents: []Event{},
		hoverEvents: []Event{},
	}
	button.rect = utils.NewRect(int(x), int(y), scale*sprite.Bounds().Dx(), scale*sprite.Bounds().Dy())
	button.Style.defaultTextColor = textColor
	return button
}

// normal square/rectangle button with border background color
func NewButton(txt string, x, y float32, fontSize int, scale int, fontFile []byte, textColor, backColor, bordercolor color.Color) *Button {

	button := &Button{
		Text: *NewLabel(txt, x, y, fontFile, fontSize, color.Black),
		Style: ButtonStyle{
			Pos:             utils.Vec2{X: x, Y: y},
			Scale:           scale,
			BorderThickness: DefaultBorderThickness,
			Color:           color.Transparent,
		},
		clickEvents: []Event{},
		hoverEvents: []Event{},
	}
	button.SetDefaultColors(textColor, backColor, bordercolor)
	button.DefaultColor()
	return button
}

// Button ButtonStyle is made Primarily for when there is no sprite or make effects for sprite
type ButtonStyle struct {
	Pos                utils.Vec2
	Color              color.Color
	BorderColor        color.Color
	BackColor          color.Color
	defaultBackColor   color.Color
	defaultBorderColor color.Color
	defaultTextColor   color.Color
	Scale              int
	BorderThickness    int
	TextOrientation    Orientation
}

// Draw button

func (b *Button) SetDefaultColors(textColor, backColor, borderColor color.Color) {
	b.Style.defaultTextColor = textColor
	b.Style.defaultBackColor = backColor
	b.Style.defaultBorderColor = borderColor
}

// check if button is being hovered on by the mouse cursor
func (b *Button) IsHover() bool {
	return b.rect.Contains(ebiten.CursorPosition())
}

// check if button is pressed by mouse
func (b *Button) IsPressed() bool {
	return (b.IsHover() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft))
}
func (b *Button) AddClickEvent(event func(b *Button)) {
	b.clickEvents = append(b.clickEvents, event)
}
func (b *Button) AddHoverEvent(event func(b *Button)) {
	b.hoverEvents = append(b.hoverEvents, event)
}

// calls all the click events
func (b *Button) OnClick() {
	for _, event := range b.clickEvents {
		event(b)
	}
}
func (b *Button) OnHover() {
	for _, event := range b.hoverEvents {
		event(b)
	}
}
func (b *Button) Centre() {
	b.UpdateRect()
	b.Style.Pos.X -= float32(b.rect.Width) / 2
	b.Style.Pos.Y -= float32(b.rect.Height) / 2
}

func (b *Button) DefaultColor() {
	b.Style.Color = color.Transparent
	b.Style.BorderColor = b.Style.defaultBorderColor
	b.Style.BackColor = b.Style.defaultBackColor
	b.Text.Style.Color = b.Style.defaultTextColor
}
func (b *Button) UpdateRect() {

	if b.sprite != nil {
		b.rect = utils.NewRect(int(b.Style.Pos.X), int(b.Style.Pos.Y), int((b.sprite.Bounds().Dx())*b.Style.Scale),
			int((b.sprite.Bounds().Dy())*b.Style.Scale))
	} else {
		f := &text.GoTextFace{Source: b.Text.Style.Font, Size: float64(b.Text.Style.Size)}
		width, height := text.Measure(b.Text.Text, f, 1)
		b.rect = utils.NewRect(int(b.Style.Pos.X), int(b.Style.Pos.Y), int(width)*b.Style.Scale, int(height)*b.Style.Scale)
	}
}
func (b *Button) draw(screen *ebiten.Image, offset utils.Vec2) {
	b.drawButtonWithOffset(screen, offset)
	b.drawButtonTextWithOffset(screen, offset)
}

func (b *Button) Draw(screen *ebiten.Image) {
	b.draw(screen, utils.Vec2{X: 0, Y: 0})
}

func (b *Button) DrawCam(screen *ebiten.Image, cam utils.Cam) {
	b.draw(screen, utils.Vec2{X: cam.X, Y: cam.Y})
}

func (b *Button) drawButtonWithOffset(screen *ebiten.Image, offset utils.Vec2) {
	if b.sprite != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(math.Round(float64(b.rect.X/(b.Style.Scale))+float64(offset.X)), math.Round(float64(b.rect.X/int(b.Style.Scale))+float64(offset.Y)))
		op.GeoM.Scale(float64(b.Style.Scale), float64(b.Style.Scale))
		//op.ColorScale.ScaleWithColor(b.Style.Color)
		screen.DrawImage(b.sprite, op)
	} else {
		vector.DrawFilledRect(screen, float32(math.Round(float64(b.rect.X/(b.Style.Scale))+float64(offset.X))), float32(math.Round(float64(b.rect.X/(b.Style.Scale))+float64(offset.Y))), float32(b.rect.Width), float32(b.rect.Height), b.Style.BackColor, false)
		vector.StrokeRect(screen, float32(b.rect.X)+offset.X, float32(b.rect.Y)+offset.Y, float32((b.rect.Width + b.Style.BorderThickness/2)),
			float32(b.rect.Height+b.Style.BorderThickness/2),
			float32(b.Style.BorderThickness),
			b.Style.BorderColor, false)
	}
}

func (b *Button) drawButtonTextWithOffset(screen *ebiten.Image, offset utils.Vec2) {
	opText := &text.DrawOptions{}
	f := &text.GoTextFace{Source: b.Text.Style.Font, Size: float64(b.Text.Style.Size)}
	width, height := text.Measure(b.Text.Text, f, 1)

	opText.GeoM.Translate(float64(float32(b.rect.X)+offset.X), float64(float32(b.rect.Y)+offset.Y))
	if b.Style.Scale != 1 {
		b.applyTextOrientation(opText, width, height)
	}
	opText.ColorScale.ScaleWithColor(b.Text.Style.Color)
	text.Draw(screen, b.Text.Text, f, opText)
	b.UpdateRect()
}

func (b *Button) applyTextOrientation(opText *text.DrawOptions, width, height float64) {
	switch b.Style.TextOrientation {
	case Middle:
		opText.GeoM.Translate(float64(b.rect.Width)/2-width/2, float64(b.rect.Height)/2-height/2)
	case BottomRight:
		opText.GeoM.Translate(float64(b.rect.Width)-width*TextBottomRightFactor, float64(b.rect.Height)-height*TextBottomRightFactor)
	default:
		opText.GeoM.Translate(width*TextPaddingFactor, height*TextPaddingFactor)
	}
}

func (b *Button) SetText(newText string) {
	b.Text.Text = newText
	b.UpdateRect()
}
func (b *Button) SetScale(scale int) {
	b.Style.Scale = scale
	b.UpdateRect()
}
func (b *Button) RemoveClickEvent(event Event) {
	for i, e := range b.clickEvents {
		if &e == &event {
			b.clickEvents = append(b.clickEvents[:i], b.clickEvents[i+1:]...)
			break
		}
	}
}
func (b *Button) RemoveHoverEvent(event Event) {
	for i, e := range b.hoverEvents {
		if &e == &event {
			b.hoverEvents = append(b.hoverEvents[:i], b.hoverEvents[i+1:]...)
			break
		}
	}
}
