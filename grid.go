package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hasona23/game/utils"
)

const (
	SPACING   = 1
	GRID_SIZE = 20.0 //tile
	TILE_SIZE = 32.0
)

type Variant int

const (
	Rigid Variant = iota
	Air
)

type Tile struct {
	Color color.Color
	Variant
	X, Y float32
}
type Tilemap struct {
	Tiles [GRID_SIZE * GRID_SIZE]*Tile
}

func NewTilemap() (t *Tilemap) {
	t = &Tilemap{}
	for i := range t.Tiles {
		tile := &Tile{}
		tile.X = float32(i%GRID_SIZE*TILE_SIZE) + float32(i%GRID_SIZE)*float32(SPACING)
		tile.Y = float32(i/GRID_SIZE*TILE_SIZE) + float32(i/GRID_SIZE)*float32(SPACING)

		if rand.Float32()*100 > 50 && i > 0 {
			tile.Color = color.Black
			tile.Variant = Rigid
		} else {
			tile.Color = color.White
			tile.Variant = Air
		}

		t.Tiles[i] = tile
	}
	return t
}

func (t Tilemap) Draw(screen *ebiten.Image) {
	for _, tile := range t.Tiles {
		vector.DrawFilledRect(screen, tile.X+game.cam.X, tile.Y+game.cam.Y, TILE_SIZE, TILE_SIZE, tile.Color, false)
	}
}
func (t *Tilemap) GetWidth() int {
	return GRID_SIZE*TILE_SIZE + GRID_SIZE*(SPACING)
}
func (t *Tilemap) GetHieght() int {
	return GRID_SIZE*TILE_SIZE + GRID_SIZE*(SPACING)
}
func (t Tilemap) GetTile(pos utils.Vec2) *Tile {
	// Calculate tile coordinates in world space
	tileX := int(pos.X) / (TILE_SIZE + SPACING)
	tileY := int(pos.Y) / (TILE_SIZE + SPACING)
	index := tileY*GRID_SIZE + tileX

	if index >= len(t.Tiles) || index < 0 {
		return nil
	}
	return t.Tiles[index]
}
