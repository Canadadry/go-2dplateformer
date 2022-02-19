package main

import (
	"app/pkg/atlas"
	"app/pkg/ecs"
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth     = 320
	screenHeight    = 240
	atlasDefinition = "ressource/spritesheet_complete.xml"
	atlasTexture    = "ressource/spritesheet_complete.png"
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
	f, err := os.Open(atlasDefinition)
	if err != nil {
		return err
	}
	defer f.Close()
	a, err := atlas.Load(f)
	if err != nil {
		return err
	}
	imgFile, err := os.Open(atlasTexture)
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
	game := &Game{
		world: ecs.New(),
	}
	game.world.AddSystem(&Painter{Game: game})
	game.world.AddEntity(map[ecs.ComponentKind]interface{}{
		ComponentKindPaintable: &Paintable{
			Frames:        []string{"alienBeige_walk1", "alienBeige_walk2"},
			FrameDuration: 5,
			X:             -30.0,
			Y:             -30.0,
			Atlas:         a,
		},
	})
	game.world.AddEntity(map[ecs.ComponentKind]interface{}{
		ComponentKindPaintable: &Paintable{
			Frames:        []string{"alienBlue_swim1", "alienBlue_swim2"},
			FrameDuration: 20,
			X:             30.0,
			Y:             30.0,
			Atlas:         a,
		},
	})
	return ebiten.RunGame(game)
}

const (
	ComponentKindPaintable ecs.ComponentKind = "painter"
)

type Paintable struct {
	X             float64
	Y             float64
	Count         int
	Frames        []string
	FrameDuration int
	Atlas         map[string]atlas.Frame
}

type Painter struct {
	Game *Game
}

func (p *Painter) Match(e ecs.Entity) bool {
	_, ok := e[ComponentKindPaintable]
	return ok
}

func (p *Painter) Update(e ecs.Entity) {
	cmpt := e[ComponentKindPaintable].(*Paintable)
	cmpt.Count++
}

func (p *Painter) Draw(e ecs.Entity, screen interface{}) {
	cmpt := e[ComponentKindPaintable].(*Paintable)
	i := (cmpt.Count / cmpt.FrameDuration) % len(cmpt.Frames)
	frame := cmpt.Atlas[cmpt.Frames[i]]
	rect := image.Rect(frame.X, frame.Y, frame.X+frame.Width, frame.Y+frame.Height)
	subImg := runnerImage.SubImage(rect).(*ebiten.Image)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frame.Width)/2, -float64(frame.Height)/2)
	op.GeoM.Scale(0.25, 0.25)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	op.GeoM.Translate(cmpt.X, cmpt.Y)
	screen.(*ebiten.Image).DrawImage(subImg, op)
}

type Game struct {
	world *ecs.World
}

func (g *Game) Update() error {
	g.world.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.world.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
