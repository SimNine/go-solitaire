package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"urffer.xyz/go-solitaire/src/game"
	"urffer.xyz/go-solitaire/src/util"
)

type Game struct {
	windowSize       util.Dims
	windowRenderDims util.Dims

	board *game.Board
}

func (g *Game) Init() {
	// Set the window size and title
	ebiten.SetWindowTitle("Solitaire")
	ebiten.SetWindowSize(g.windowSize.X, g.windowSize.Y)
}

func (g *Game) Update() error {
	pos := util.MakePosFromTuple(ebiten.CursorPosition())
	g.board.SetCusrorPos(pos)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.board.GrabCard()
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.board.ReleaseCard()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.board.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.windowRenderDims.X, g.windowRenderDims.Y
}

func main() {
	// Initialize the game assets
	game.InitCards()
	game.InitCardStackBkg()

	// Create the game instance, init, and run it
	ebitengineGame := &Game{
		windowSize:       util.Dims{X: 1000, Y: 800},
		windowRenderDims: util.Dims{X: 1000, Y: 800},
		board:            game.NewBoard(),
	}
	ebitengineGame.Init()
	if err := ebiten.RunGame(ebitengineGame); err != nil {
		log.Fatal(err)
	}
}
