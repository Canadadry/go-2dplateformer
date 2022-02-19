package main

import (
	"app/game"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

func main() {
	g, err := game.New(screenWidth, screenHeight)
	if err != nil {
		fmt.Println("init error:", err)
	}
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebiten Demo)")
	err = ebiten.RunGame(g)
	if err != nil {
		fmt.Println("run error:", err)
	}
}
