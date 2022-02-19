package game

import (
	"app/pkg/atlas"
	"app/pkg/ecs"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

const ComponentKindPaintable ecs.ComponentKind = "painter"

type Paintable struct {
	X             float64
	Y             float64
	Count         int
	Frames        []string
	FrameDuration int
	Atlas         map[string]atlas.Frame
	Texture       *ebiten.Image
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
	subImg := cmpt.Texture.SubImage(rect).(*ebiten.Image)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frame.Width)/2, -float64(frame.Height)/2)
	op.GeoM.Scale(0.25, 0.25)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	op.GeoM.Translate(cmpt.X, cmpt.Y)
	screen.(*ebiten.Image).DrawImage(subImg, op)
}
