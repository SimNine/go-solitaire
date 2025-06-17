package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"urffer.xyz/go-solitaire/src/game"
	"urffer.xyz/go-solitaire/src/util"
)

type Game struct {
	windowSize       util.Dims
	windowRenderDims util.Dims

	testCard *game.Card
}

func (g *Game) Init() {
	ebiten.SetWindowTitle("Solitaire")
	ebiten.SetWindowSize(g.windowSize.X, g.windowSize.Y)
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Fill the background with the board color
	screen.Fill(color.RGBA{
		R: 0,
		G: 75,
		B: 0,
		A: 255,
	})

	// Draw the test card
	g.testCard.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.windowRenderDims.X, g.windowRenderDims.Y
}

func main() {

	game := &Game{
		windowSize:       util.Dims{X: 800, Y: 800},
		windowRenderDims: util.Dims{X: 400, Y: 400},
		testCard:         game.MakeCard(game.Five, game.Club),
	}
	game.Init()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
