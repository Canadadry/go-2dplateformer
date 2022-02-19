package game

import (
	"app/pkg/atlas"
	"app/pkg/ecs"
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

type Game struct {
	world        *ecs.World
	runnerImage  *ebiten.Image
	screenWidth  int
	screenHeight int
}

func New(w, h int) (*Game, error) {
	f, err := os.Open(atlasDefinition)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	a, err := atlas.Load(f)
	if err != nil {
		return nil, err
	}
	imgFile, err := os.Open(atlasTexture)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	game := &Game{
		world:        ecs.New(),
		runnerImage:  ebiten.NewImageFromImage(img),
		screenWidth:  w,
		screenHeight: h,
	}
	game.world.AddSystem(&Painter{Game: game})
	game.world.AddEntity(map[ecs.ComponentKind]interface{}{
		ComponentKindPaintable: &Paintable{
			Frames:        []string{"alienBeige_walk1", "alienBeige_walk2"},
			FrameDuration: 5,
			X:             -30.0,
			Y:             -30.0,
			Atlas:         a,
			Texture:       game.runnerImage,
		},
	})
	game.world.AddEntity(map[ecs.ComponentKind]interface{}{
		ComponentKindPaintable: &Paintable{
			Frames:        []string{"alienBlue_swim1", "alienBlue_swim2"},
			FrameDuration: 20,
			X:             30.0,
			Y:             30.0,
			Atlas:         a,
			Texture:       game.runnerImage,
		},
	})
	return game, nil
}

func (g *Game) Update() error {
	g.world.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.world.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}
