package main

import (
	"log"

	"mygo/internal/pkg/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.NewGame()
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("My Go - 2D Sandbox Roguelike")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
