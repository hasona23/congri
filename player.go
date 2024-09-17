package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hasona23/game/utils"
)

const (
	PLAYER_RECT_SIZE = 16
	ACCELRATION      = 0.75
	PLAYER_FIRERATE  = 0.75
)

type Player struct {
	DynamicEntity
	hp       int
	mana     int
	fireRate utils.Timer
}

var directions []utils.Vec2 = []utils.Vec2{{X: 1, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}, {X: 0, Y: -1}, {X: -1, Y: -1}, {X: 1, Y: -1}, {X: -1, Y: 1}}

// Update implements Entity.
func NewPlayer(x, y float32) *Player {
	return &Player{DynamicEntity: DynamicEntity{utils.Vec2{X: 4, Y: 4}, utils.Vec2{X: x, Y: y}, 1, "player", color.RGBA{128, 0, 129, 255}, false}, hp: HP, mana: 100, fireRate: utils.NewTimer(PLAYER_FIRERATE)}

}
func (p Player) IsDestroyed() bool {
	return p.Destroyed || p.hp <= 0
}
func (p Player) Type() string {
	return p.etype
}
func (p Player) Draw(screen *ebiten.Image) {
	//	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%v", p.hp), int(p.Pos.X)+int(game.cam.X)-2, int(p.Pos.Y)-12+int(game.cam.Y))
	vector.DrawFilledRect(screen, p.Pos.X+game.cam.X, p.Pos.Y+game.cam.Y, PLAYER_RECT_SIZE, PLAYER_RECT_SIZE, p.color, false)
	// Draw nearby tile boundaries
	/*for _, tile := range p.GetNearTiles() {
		if tile.Variant == Rigid {
			vector.DrawFilledRect(screen, tile.X+game.cam.X, tile.Y+game.cam.Y, TILE_SIZE, TILE_SIZE, color.RGBA{255, 0, 0, 100}, false)
		}
	}*/
}
func lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}
func (p *Player) Update() {
	p.fireRate.UpdateTimer()
	p.mana = int(math.Min(math.Max(0, float64(p.mana)), 100))
	//fmt.Println(p.fireRate.GetCurrentTime())
	p.Dir.X = float32(math.Round(float64(lerp(p.Dir.X, 0, ACCELRATION))))
	p.Dir.Y = float32(math.Round(float64(lerp(p.Dir.Y, 0, ACCELRATION))))
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.Dir.Y = float32(math.Round(float64(lerp(p.Dir.Y, -1, ACCELRATION))))
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.Dir.Y = float32(math.Round(float64(lerp(p.Dir.Y, 1, ACCELRATION))))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.Dir.X = float32(math.Round(float64(lerp(p.Dir.X, 1, ACCELRATION))))
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.Dir.X = float32(math.Round(float64(lerp(p.Dir.X, -1, ACCELRATION))))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) && p.mana >= 100 {

		for _, dir := range directions {
			NewBullet("player", p.Pos, dir, 2)
		}
		p.mana = 0

	}

	if (inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || ebiten.IsKeyPressed(ebiten.KeyE)) && p.fireRate.Ticked() {
		x, y := ebiten.CursorPosition()
		x -= int(game.cam.X)
		y -= int(game.cam.Y)
		NewBullet("player", p.Pos, utils.Vec2{X: float32(x) - (p.Pos.X), Y: float32(y) - (p.Pos.Y)}, 2)
	}
	p.Dir.NormalizeDir()
	dx := int(math.Round(float64(p.Dir.X * p.speed)))
	p.horizontalCollision(dx)
	dy := int(math.Round(float64(p.Dir.Y * p.speed)))
	p.verticalCollision(dy)
	p.constraintMovemnt()
	// fmt.Printf("Velocity:%2v\n", p.Dir.X*p.speed)
	for _, b := range game.entities["bullet"] {
		if b.(*Bullet).rect().Collide(p.rect()) && b.(*Bullet).Shooter == "sniper" {
			p.hp -= 30
			b.(*Bullet).Destroyed = true
			fmt.Println("I GOT SHOT")
		}
	}
}
func (p *Player) horizontalCollision(dx int) {
	collisions := map[string]bool{"right": false, "left": false}

	rect := p.rect()
	rect.X += dx
	for _, tile := range p.GetNearTiles() {
		if tile.Variant == Rigid {

			tRect := utils.NewRect(int(tile.X), int(tile.Y), TILE_SIZE, TILE_SIZE)
			if rect.Collide(tRect) {
				if dx > 0.0 {
					collisions["right"] = true
				} else if dx < 0.0 {
					collisions["left"] = true
				}
				//fmt.Println("IsColliding:", collisions)
			}
		}
	}
	if !collisions["right"] && dx > 0 || !collisions["left"] && dx < 0 {

		p.Pos.X += float32(dx)
	}

}
func (p *Player) verticalCollision(dy int) {
	collisions := map[string]bool{"up": false, "down": false}

	rect := p.rect()
	rect.Y += dy
	for _, tile := range p.GetNearTiles() {
		if tile.Variant == Rigid {

			tRect := utils.NewRect(int(tile.X), int(tile.Y), TILE_SIZE, TILE_SIZE)
			if rect.Collide(tRect) {
				if dy > 0.0 {
					collisions["down"] = true
				} else if dy < 0.0 {
					collisions["up"] = true
				}
				//fmt.Println("IsColliding:", collisions)
			}
		}
	}

	if !collisions["down"] && dy > 0 || !collisions["up"] && dy < 0 {
		p.Pos.Y += float32(dy)
	}

}

func (p *Player) constraintMovemnt() {
	p.Pos.X = float32(math.Min(float64(game.Tilemap.GetWidth()-PLAYER_RECT_SIZE), math.Max(float64(p.Pos.X), 0)))
	p.Pos.Y = float32(math.Min(float64(game.Tilemap.GetHieght()-PLAYER_RECT_SIZE), math.Max(float64(p.Pos.Y), 0)))

}
func (p Player) rect() utils.Rect {
	return utils.NewRect(int(p.Pos.X), int(p.Pos.Y), PLAYER_RECT_SIZE, PLAYER_RECT_SIZE)
}
