package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hasona23/game/utils"
)

const (
	HP = 100
)

type Entity interface {
	Update()
	Draw(screen *ebiten.Image)
	IsDestroyed() bool
	Type() string
}
type DynamicEntity struct {
	Pos       utils.Vec2
	Dir       utils.Vec2
	speed     float32
	etype     string
	color     color.Color
	Destroyed bool
}

func (g *Game) AddEntity(e Entity) {
	g.entities[e.Type()] = append(g.entities[e.Type()], e)
}
func (e DynamicEntity) GetNearTiles() map[utils.Vec2]*Tile {
	positions := []utils.Vec2{
		//same pos
		{X: e.Pos.X - TILE_SIZE, Y: e.Pos.Y},             // Left
		{X: e.Pos.X + TILE_SIZE, Y: e.Pos.Y},             // Right
		{X: e.Pos.X, Y: e.Pos.Y - TILE_SIZE},             // Up
		{X: e.Pos.X, Y: e.Pos.Y + TILE_SIZE},             // Down
		{X: e.Pos.X + TILE_SIZE, Y: e.Pos.Y - TILE_SIZE}, // Diagonal top-right
		{X: e.Pos.X - TILE_SIZE, Y: e.Pos.Y + TILE_SIZE}, // Diagonal bottom-left
		{X: e.Pos.X + TILE_SIZE, Y: e.Pos.Y + TILE_SIZE}, // Diagonal top-left
		{X: e.Pos.X - TILE_SIZE, Y: e.Pos.Y - TILE_SIZE}, // Diagonal bottom-right
	}
	tiles := make(map[utils.Vec2]*Tile)
	for _, pos := range positions {
		if tile := game.Tilemap.GetTile(pos); tile != nil {
			tiles[pos] = tile
		}
	}
	return tiles
}
