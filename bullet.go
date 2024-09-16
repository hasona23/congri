package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hasona23/game/particles"
	"github.com/hasona23/game/utils"
)

const (
	BULLET_SIZE = 8
)

type Bullet struct {
	Shooter string
	DynamicEntity
	tilesCounter int
	currentTile  *Tile
}

func NewBullet(Shooter string, pos, dir utils.Vec2, speed float32) *Bullet {
	dir.NormalizeDir()
	b := &Bullet{
		Shooter:       Shooter,
		DynamicEntity: DynamicEntity{pos, dir, speed, "bullet", color.RGBA{0, 191, 255, 255}, false},
		tilesCounter:  8,
	}
	game.AddEntity(b)
	return b
}
func (b Bullet) Type() string {
	return b.etype
}
func (b Bullet) IsDestroyed() bool {
	if b.Destroyed || b.tilesCounter <= 0 {
		particlesSystem := particles.NewParticleSystem(
			particles.WithArea(utils.NewRect(int(b.Pos.X)-8, int(b.Pos.Y)-8, 16, 16)),
			particles.WithMotionType(particles.Outward),
			particles.WithShrinking(0.075),
			particles.WithModelParticle(*particles.NewParticle(particles.WithColor(b.color), particles.WithScale(BULLET_SIZE/2),
				particles.WithSpeed(1))))
		particlesSystem.Spawn(10)
		game.particles = append(game.particles, particlesSystem)
		return true
	}
	return false
}
func (b *Bullet) Update() {
	b.Pos.X += b.Dir.X * b.speed
	b.Pos.Y += b.Dir.Y * b.speed
	b.handleCollisions()
	for _, e := range game.entities["bullet"] {
		b2 := e.(*Bullet)
		if b.rect().Collide(b2.rect()) && b2.Shooter != b.Shooter {
			b.Destroyed = true
			b2.Destroyed = true
		}
	}
}
func (b Bullet) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, b.Pos.X+game.cam.X, b.Pos.Y+game.cam.Y, BULLET_SIZE, BULLET_SIZE, b.color, false)

}

func (b Bullet) rect() utils.Rect {
	return utils.NewRect(int(b.Pos.X), int(b.Pos.Y), BULLET_SIZE, BULLET_SIZE)
}
func (b *Bullet) handleCollisions() {
	if b.Shooter == "sniper" {
		return
	}
	if tile := game.Tilemap.GetTile(b.Pos); tile != nil && tile.Variant == Rigid &&
		b.rect().Collide(utils.NewRect(int(tile.X), int(tile.Y), TILE_SIZE, TILE_SIZE)) {
		tile.Variant = Air
		tile.Color = color.White
		b.tilesCounter--
		//b.currentTile = tile
		//	fmt.Println(b.tilesCounter)
	}

}
