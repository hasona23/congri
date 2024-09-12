package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hasona23/game/utils"
)

const (
	ENEMY_SIZE  = 16
	ENEMY_SPEED = 1
)

type Enemy struct {
	DynamicEntity
	hp int
}

func NewEnemy(pos utils.Vec2) *Enemy {
	e := &Enemy{
		DynamicEntity: DynamicEntity{pos, utils.Vec2{X: 0, Y: 0}, ENEMY_SPEED, "enemy", color.RGBA{255, 0, 0, 255}, false},
		hp:            HP,
	}
	game.AddEntity(e)
	return e
}

func (e Enemy) Type() string {
	return e.etype
}
func (e Enemy) IsDestroyed() bool {
	return e.Destroyed || e.hp <= 0
}
func (e *Enemy) Update() {
	player := game.entities["player"][0].(*Player)

	e.Dir = utils.Vec2{X: player.Pos.X - e.Pos.X, Y: player.Pos.Y - e.Pos.Y}
	e.Dir.NormalizeDir()
	e.Pos.X += e.Dir.X * e.speed
	e.Pos.Y += e.Dir.Y * e.speed

	if e.rect().Collide(player.rect()) {
		e.Destroyed = true
		player.hp -= 20
	}
	for _, b := range game.entities["bullet"] {
		if b.(*Bullet).rect().Collide(e.rect()) {
			b.(*Bullet).Destroyed = true
			e.hp -= 50
		}
	}
	tile := game.Tilemap.GetTile(e.Pos)
	if tile.Variant == Air && e.rect().Collide(utils.NewRect(int(tile.X), int(tile.Y), TILE_SIZE, TILE_SIZE)) {
		tile.Variant = Rigid
		tile.Color = color.Black
	}
}
func (e Enemy) rect() utils.Rect {
	return utils.NewRect(int(e.Pos.X), int(e.Pos.Y), ENEMY_SIZE, ENEMY_SIZE)
}
func (e *Enemy) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, e.Pos.X+game.cam.X, e.Pos.Y+game.cam.Y, ENEMY_SIZE, ENEMY_SIZE, e.color, false)
}
