package utils

import (
	"math"
)

type Cam struct {
	X, Y float32
}

func NewCamera(x, y float32) *Cam {
	return &Cam{
		X: x,
		Y: y,
	}
}

// ShakeAmount is how many times you want the shaking to occur
/*func (c *Camera) Shake(strengh float32, ShakeAmount int) {
	//factors is between  -1 , 1 to maniplutate direction of x and y displacement
	factor := 0
	for range ShakeAmount {
		if rand.Int()%2 == 0 {
			factor = 1
		} else {
			factor = -1
		}
		x := rand.Float32() * strengh * float32(factor)
		y := rand.Float32() * strengh * float32(factor)
		c.X += int(x)
		c.Y += int(y)
	}
}*/

// For Smooth Movement
func (c *Cam) FollowTarget(targetX, targetY, screenWidth, screenHeight, speed float32) {
	c.X += (-targetX + screenWidth/2.0 - c.X) / speed
	c.Y += (-targetY + screenHeight/2.0 - c.Y) / speed
}

// For Sudden and fast without delay movemnt
func (c *Cam) GoTo(targetX, targetY, screenWidth, screenHeight float32) {
	c.X = (-targetX + screenWidth/2.0)
	c.Y = (-targetY + screenHeight/2.0)
}

func (c *Cam) Constrain(tilemapWidthPixels, tilemapHeightPixels, screenWidth, screenHeight int) {
	c.X = float32(math.Min(float64(c.X), 0.0))
	c.Y = float32(math.Min(float64(c.Y), 0.0))

	c.X = float32(math.Max(float64(c.X), float64(screenWidth-tilemapWidthPixels)))
	c.Y = float32(math.Max(float64(c.Y), float64(screenHeight-tilemapHeightPixels)))
}
