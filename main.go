package main

import (
	"app/pkg/atlas"
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

var (
	runnerImage *ebiten.Image
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("failed", err)
	}
}

func run() error {
	f, err := os.Open("ressource/spritesheet_complete.xml")
	if err != nil {
		return err
	}
	defer f.Close()
	a, err := atlas.Load(f)
	if err != nil {
		return err
	}
	imgFile, err := os.Open("ressource/spritesheet_complete.png")
	if err != nil {
		return err
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return err
	}
	runnerImage = ebiten.NewImageFromImage(img)

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebiten Demo)")
	return ebiten.RunGame(&Game{
		atlas:         a,
		frames:        []string{"alienBeige_walk1", "alienBeige_walk2"},
		frameDuration: 5,
	})
}

type Game struct {
	count         int
	frames        []string
	frameDuration int
	atlas         map[string]atlas.Frame
}

func (g *Game) Update() error {
	g.count++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	i := (g.count / g.frameDuration) % len(g.frames)
	frame := g.atlas[g.frames[i]]
	rect := image.Rect(frame.X, frame.Y, frame.X+frame.Width, frame.Y+frame.Height)
	subImg := runnerImage.SubImage(rect).(*ebiten.Image)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frame.Width)/2, -float64(frame.Height)/2)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	op.GeoM.Scale(0.25, 0.25)
	screen.DrawImage(subImg, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
