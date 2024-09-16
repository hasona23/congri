package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hasona23/game/particles"
	"github.com/hasona23/game/utils"
)

const (
	ENEMY_SIZE         = 16
	ENEMY_SPEED        = 1
	ENEMY_BULLET_SPEED = 1.5
)

var BomberColor = color.RGBA{255, 0, 0, 255}
var SniperColor = color.RGBA{255, 240, 0, 255}

type Bomber struct {
	DynamicEntity
}

func NewBomber(pos utils.Vec2) *Bomber {
	e := &Bomber{
		DynamicEntity: DynamicEntity{pos, utils.Vec2{X: 0, Y: 0}, ENEMY_SPEED, "enemy", BomberColor, false},
	}
	game.AddEntity(e)
	return e
}

func (e Bomber) Type() string {
	return e.etype
}
func (e Bomber) IsDestroyed() bool {
	if e.Destroyed {
		particlesSystem := particles.NewParticleSystem(
			particles.WithArea(utils.NewRect(int(e.Pos.X-8), int(e.Pos.Y-8), 16, 16)),
			particles.WithMotionType(particles.Outward),
			particles.WithShrinking(0.075),
			particles.WithModelParticle(*particles.NewParticle(particles.WithColor(e.color), particles.WithScale(ENEMY_SIZE/2),
				particles.WithSpeed(1))))
		particlesSystem.Spawn(10)
		game.particles = append(game.particles, particlesSystem)
		return true
	}
	return false
}
func (e *Bomber) Update() {
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
			e.Destroyed = true
			b.(*Bullet).Destroyed = true
			game.score++
		}
	}
	tile := game.Tilemap.GetTile(e.Pos)
	if tile.Variant == Air && e.rect().Collide(utils.NewRect(int(tile.X), int(tile.Y), TILE_SIZE, TILE_SIZE)) {
		tile.Variant = Rigid
		tile.Color = color.Black
	}

}
func (e Bomber) rect() utils.Rect {
	return utils.NewRect(int(e.Pos.X), int(e.Pos.Y), ENEMY_SIZE, ENEMY_SIZE)
}
func (e *Bomber) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, e.Pos.X+game.cam.X, e.Pos.Y+game.cam.Y, ENEMY_SIZE, ENEMY_SIZE, e.color, false)
}

type Sniper struct {
	DynamicEntity
	fireRate utils.Timer
}

func NewSniper(pos utils.Vec2) *Sniper {
	s := &Sniper{DynamicEntity: DynamicEntity{pos, utils.Vec2{X: 0, Y: 0}, 0, "enemy", SniperColor, false}, fireRate: utils.NewTimer(3)}
	game.AddEntity(s)
	return s
}
func (s Sniper) rect() utils.Rect {
	return utils.NewRect(int(s.Pos.X), int(s.Pos.Y), ENEMY_SIZE, ENEMY_SIZE)
}

func (s Sniper) Type() string {
	return s.etype
}
func (s Sniper) IsDestroyed() bool {
	if s.Destroyed {
		particlesSystem := particles.NewParticleSystem(
			particles.WithArea(utils.NewRect(int(s.Pos.X-8), int(s.Pos.Y-8), 16, 16)),
			particles.WithMotionType(particles.Outward),
			particles.WithShrinking(0.075),
			particles.WithModelParticle(*particles.NewParticle(particles.WithColor(s.color), particles.WithScale(ENEMY_SIZE/2),
				particles.WithSpeed(1))))
		particlesSystem.Spawn(10)
		game.particles = append(game.particles, particlesSystem)
		return true
	}
	return false
}
func (s *Sniper) Update() {
	player := game.entities["player"][0].(*Player)
	s.fireRate.UpdateTimer()
	if s.fireRate.Ticked() {
		b := NewBullet("sniper", s.Pos, utils.Vec2{X: float32(player.Pos.X) - (s.Pos.X), Y: float32(player.Pos.Y) - (s.Pos.Y)}, ENEMY_BULLET_SPEED)
		b.color = s.color
	}
	if s.rect().Collide(player.rect()) {
		s.Destroyed = true
		player.hp -= 20
	}
	for _, b := range game.entities["bullet"] {
		if b.(*Bullet).rect().Collide(s.rect()) && b.(*Bullet).Shooter != "sniper" {
			s.Destroyed = true
			b.(*Bullet).Destroyed = true
			game.score++
		}
	}

}

func (s Sniper) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, s.Pos.X+game.cam.X, s.Pos.Y+game.cam.Y, ENEMY_SIZE, ENEMY_SIZE, s.color, false)
}
