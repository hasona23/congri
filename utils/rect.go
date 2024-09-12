package utils

type Rect struct {
	X, Y, Width, Height int
}

func NewRect(x, y, width, height int) Rect {
	return Rect{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

func (r Rect) Right() int {
	return r.X + r.Width
}

func (r Rect) Bottom() int {
	return r.Y + r.Height
}
func (r Rect) Left() int {
	return r.X
}
func (r Rect) Top() int {
	return r.Y
}
func (r Rect) Contains(x, y int) bool {
	return (x < r.Right() && x > r.Left() && y < r.Bottom() && y > r.Top())
}
func (r Rect) Collide(other Rect) bool {
	return r.X < other.Right() &&
		other.X < r.Right() &&
		r.Y < other.Bottom() &&
		other.Y < r.Bottom()
}
func (r Rect) Centre() (int, int) {
	return r.X + r.Width/2, r.Y + r.Height/2
}
