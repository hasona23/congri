package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"
	"os"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hasona23/game/particles"
	"github.com/hasona23/game/ui"
	"github.com/hasona23/game/utils"
)

const (
	SPAWN_TIME = 3
)

type Game struct {
	cam          utils.Cam
	Tilemap      *Tilemap
	entities     map[string][]Entity
	particles    []*particles.ParticleSystem
	enemySpawner utils.Timer
	score        int
	state        State
	ui           map[State]*ui.UILayout
}

func (g *Game) Init() {
	font, err := os.ReadFile("./pixelFont.ttf")
	if err != nil {
		log.Fatal("Error opening font file: err")
	}
	g.cam = *utils.NewCamera(0, 0)
	g.Tilemap = NewTilemap()
	g.entities = make(map[string][]Entity)
	g.AddEntity(NewPlayer(5, 5))
	g.enemySpawner = utils.NewTimer(SPAWN_TIME)
	g.state = States.Menu
	g.score = 0
	//UI===============================
	//main menu
	g.ui = make(map[State]*ui.UILayout)
	menuLayout := ui.NewUILayout("menu")
	startbtn := ui.NewButton("Start", 120, 120, 16, 2, font, color.White, color.Black, color.Black)
	startbtn.AddClickEvent(func(b *ui.Button) { g.state = States.Main })
	exitbtn := ui.NewButton("Exit", 120, 160, 16, 2, font, color.White, color.Black, color.Black)
	exitbtn.AddClickEvent(func(b *ui.Button) { os.Exit(0) })
	menuLayout.AddButton("startbtn", startbtn)
	menuLayout.AddButton("exitbtn", exitbtn)
	menuLayout.ApplyHoverToAllButtons(onhover)
	g.ui[States.Menu] = menuLayout
	//game ui
	mainLayout := ui.NewUILayout("main")
	hpBar := ui.NewBar(5, 5, g.entities["player"][0].(*Player).hp, 8, utils.Point{X: 1, Y: 1}, color.RGBA{255, 0, 100, 255}, color.Gray{123})
	mainLayout.AddBar("hp", hpBar)
	score := ui.NewLabel(fmt.Sprintf("Score:%v", g.score), 5, 15, font, 16, color.RGBA{0, 100, 255, 255})
	mainLayout.AddLabel("score", score)
	g.ui[States.Main] = mainLayout

}
func onhover(b *ui.Button) {
	b.Style.BorderColor = color.White
}

func (g *Game) Update() error {
	switch g.state {
	case States.Menu:
		g.ui[States.Menu].Update()
	case States.Main:
		g.ui[States.Main].Update()
		bar, _ := g.ui[States.Main].GetBar("hp")
		bar.SetValue(g.entities["player"][0].(*Player).hp)
		player := g.entities["player"][0].(*Player)
		label, _ := g.ui[States.Main].GetLabel("score")
		label.SetText(fmt.Sprintf("score: %v", g.score))
		g.cam.FollowTarget(player.Pos.X, player.Pos.Y, 320, 240, 2)
		g.cam.Constrain(g.Tilemap.GetWidth(), g.Tilemap.GetHieght(), 320, 240)
		g.enemySpawner.UpdateTimer()
		if g.enemySpawner.Ticked() {
			x := rand.Float32() * GRID_SIZE * TILE_SIZE
			y := rand.Float32() * GRID_SIZE * TILE_SIZE
			particlesSystem := particles.NewParticleSystem(
				particles.WithArea(utils.NewRect(int(x), int(y), 32, 32)),
				particles.WithName("spawn"),
				particles.WithMotionType(particles.Circular),
				particles.WithShrinking(0.2),
				particles.WithModelParticle(*particles.NewParticle(particles.WithColor(color.RGBA{255, 0, 0, 255}), particles.WithScale(16),
					particles.WithSpeed(0.5))))
			particlesSystem.Spawn(10)

			g.particles = append(g.particles, particlesSystem)
		}
		for _, entities := range g.entities {
			for i := range entities {
				if !entities[i].IsDestroyed() {
					entities[i].Update()
				}
			}
		}
		for k := range g.entities {
			g.entities[k] = slices.DeleteFunc(g.entities[k], func(e Entity) bool { return e.IsDestroyed() })
		}

		g.particles = slices.DeleteFunc(g.particles, func(ps *particles.ParticleSystem) bool {
			if len(ps.Particles) == 0 && ps.Name == "spawn" {
				x, y := ps.Area.Centre()
				n := rand.Float32() * 100
				if n > 25 {
					NewBomber(utils.Vec2{X: float32(x), Y: float32(y)})
				} else {
					NewSniper(utils.Vec2{X: float32(x), Y: float32(y)})
				}

				return true
			}
			return false
		})
		for i := range g.particles {
			g.particles[i].Update()
		}
		if len(g.entities["player"]) == 0 {
			game.Init()
			game.state = States.Main
		}

	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch {
	case g.state == States.Menu:
		screen.Fill(color.Black)
		g.ui[States.Menu].Draw(screen)
	case g.state == States.Main || g.state == States.Pause:
		screen.Fill(color.RGBA{100, 50, 120, 255})
		g.Tilemap.Draw(screen)
		for _, ps := range g.particles {
			ps.DrawCam(screen, g.cam)
		}
		for _, entities := range g.entities {
			for i := range entities {
				entities[i].Draw(screen)
			}
		}
		g.ui[States.Main].Draw(screen)
	case g.state == States.Pause:

	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

var game Game

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Survive")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	game.Init()
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
