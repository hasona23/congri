package ui

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hasona23/game/utils"
)

type Label struct {
	Text  string
	Pos   utils.Vec2
	Style Style
}
type Style struct {
	Font        *text.GoTextFaceSource
	Size        int
	Color       color.Color
	LineSpacing uint32
}

func NewLabel(txt string, x, y float32, fontFile []byte, size int, color color.Color) *Label {
	font, err := text.NewGoTextFaceSource(bytes.NewReader(fontFile))
	if err != nil {
		log.Fatalf("Error loading font file %v", err)
	}

	return &Label{
		Text:  txt,
		Pos:   utils.Vec2{X: x, Y: y},
		Style: NewStyle(font, size, color),
	}
}
func NewStyle(font *text.GoTextFaceSource, size int, color color.Color) Style {
	return Style{
		Font:        font,
		Size:        size,
		Color:       color,
		LineSpacing: 1,
	}
}
func (s Style) Face() *text.GoTextFace {
	return &text.GoTextFace{Source: s.Font, Size: float64(s.Size)}
}
func (l *Label) SetText(txt string) {
	l.Text = txt
}
func (l *Label) SetStyle(style Style) {
	l.Style = style
}
func (l *Label) SetPosition(x, y float32) {
	l.Pos = utils.Vec2{X: x, Y: y}
}
func (l *Label) GetText() string {
	return l.Text
}

func (l *Label) GetPosition() utils.Vec2 {
	return l.Pos
}

func (l *Label) GetStyle() Style {
	return l.Style
}
func (l *Label) Draw(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(l.Pos.X), float64(l.Pos.Y))
	op.ColorScale.ScaleWithColor(l.Style.Color)
	text.Draw(screen, l.Text, l.Style.Face(), op)
}
func (l *Label) DrawCam(screen *ebiten.Image, cam utils.Cam) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(l.Pos.X+cam.X), float64(l.Pos.Y)+float64(cam.Y))
	op.ColorScale.ScaleWithColor(l.Style.Color)
	text.Draw(screen, l.Text, l.Style.Face(), op)
}
func (l *Label) CenterText() {
	w, h := l.GetDimensions()
	l.Pos.X -= float32(w / 2)
	l.Pos.Y -= float32(h / 2)
}
func (l *Label) GetDimensions() (width, height int) {
	f := l.Style.Face()
	w, h := text.Measure(l.Text, f, float64(l.Style.LineSpacing))
	return int(w), int(h)
}
func (l *Label) SetColor(c color.Color) {
	l.Style.Color = c
}
func (l *Label) GetBounds() (x, y float32, width, height int) {
	width, height = l.GetDimensions()
	return l.Pos.X, l.Pos.Y, width, height
}
func (l *Label) SetFontSize(size int) {
	l.Style.Size = size
}
func (l *Label) Move(dx, dy float32) {
	l.Pos.X += dx
	l.Pos.Y += dy
}

func (l *Label) MoveX(dx float32) {
	l.Pos.X += dx
}

func (l *Label) MoveY(dy float32) {
	l.Pos.Y += dy
}
