package main

import (
	"image/color"
	"log"
	"math/rand/v2"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hasona23/game/particles"
	"github.com/hasona23/game/utils"
)

type Game struct {
	cam          utils.Cam
	Tilemap      *Tilemap
	entities     map[string][]Entity
	particles    []*particles.ParticleSystem
	enemySpawner utils.Timer
}

func (g *Game) Init() {
	g.cam = *utils.NewCamera(0, 0)
	g.Tilemap = NewTilemap()
	g.entities = make(map[string][]Entity)
	g.AddEntity(NewPlayer(5, 5))
	g.enemySpawner = utils.NewTimer(5)
}
func (g *Game) Update() error {

	player := g.entities["player"][0].(*Player)
	g.cam.FollowTarget(player.Pos.X, player.Pos.Y, 320, 240, 2)
	g.cam.Constrain(g.Tilemap.GetWidth(), g.Tilemap.GetHieght(), 320, 240)
	g.enemySpawner.UpdateTimerTPS()
	if g.enemySpawner.Ticked() {
		x := rand.Float32() * GRID_SIZE * TILE_SIZE
		y := rand.Float32() * GRID_SIZE * TILE_SIZE
		particlesSystem := particles.NewParticleSystem(
			particles.WithArea(utils.NewRect(int(x), int(y), 16, 16)),
			particles.WithMotionType(particles.Outward),
			particles.WithShrinking(0.075),
			particles.WithModelParticle(*particles.NewParticle(particles.WithColor(color.RGBA{255, 0, 0, 255}), particles.WithScale(8),
				particles.WithSpeed(1))))
		particlesSystem.Spawn(10)
		g.particles = append(g.particles, particlesSystem)
	}
	for _, entities := range g.entities {
		for i := range entities {
			entities[i].Update()
		}
	}
	for k := range g.entities {
		g.entities[k] = slices.DeleteFunc(g.entities[k], func(e Entity) bool { return e.IsDestroyed() })
	}

	g.particles = slices.DeleteFunc(g.particles, func(ps *particles.ParticleSystem) bool {
		if len(ps.Particles) == 0 && (ps.ModelParticle.Color == color.Color(color.RGBA{255, 0, 0, 255})) {
			NewEnemy(utils.Vec2{X: float32(ps.Area.X), Y: float32(ps.Area.Y)})
			return true
		}
		return false
	})
	for i := range g.particles {
		g.particles[i].Update()
	}
	if len(g.entities["player"]) == 0 {
		game.Init()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{100, 50, 120, 255})
	g.Tilemap.Draw(screen)
	/*for x := range int(GRID_SIZE) {
		for y := range int(GRID_SIZE) {
			tile := g.Tilemap.GetTile(utils.Vec2{X: float32(x * TILE_SIZE), Y: float32(y * TILE_SIZE)})
			if tile != nil {
				if tile.Variant == Rigid {
					vector.DrawFilledRect(screen, float32(tile.X)+g.cam.X, float32(tile.Y)+g.cam.Y, TILE_SIZE, TILE_SIZE, color.RGBA{0, 255, 0, 255}, false)
				} else {
					vector.DrawFilledRect(screen, float32(tile.X)+g.cam.X, float32(tile.Y)+g.cam.Y, TILE_SIZE, TILE_SIZE, color.RGBA{0, 0, 255, 255}, false)

				}
			}
		}
	}*/
	for _, ps := range g.particles {
		ps.DrawCam(screen, g.cam)
	}
	for _, entities := range g.entities {
		for i := range entities {
			entities[i].Draw(screen)
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

var game Game

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Survive")
	game.Init()
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
