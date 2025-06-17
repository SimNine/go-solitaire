package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"urffer.xyz/go-solitaire/src/game"
	"urffer.xyz/go-solitaire/src/util"
)

type Game struct {
	wordsPos    util.Pos
	xIncreasing bool
	yIncreasing bool

	windowSize       util.Dims
	windowRenderDims util.Dims

	testCard *game.Card
}

func (g *Game) Init() {
	ebiten.SetWindowTitle("Solitaire")
	ebiten.SetWindowSize(g.windowSize.X, g.windowSize.Y)
}

func (g *Game) Update() error {
	if g.xIncreasing {
		g.wordsPos.X += 1
	} else {
		g.wordsPos.X -= 1
	}
	if g.yIncreasing {
		g.wordsPos.Y += 1
	} else {
		g.wordsPos.Y -= 1
	}

	if g.wordsPos.X > g.windowRenderDims.X {
		g.xIncreasing = false
	} else if g.wordsPos.X <= 0 {
		g.xIncreasing = true
	}
	if g.wordsPos.Y > g.windowRenderDims.Y {
		g.yIncreasing = false
	} else if g.wordsPos.Y <= 0 {
		g.yIncreasing = true
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrintAt(screen, "ayoooo", g.wordsPos.X, g.wordsPos.Y)
	g.testCard.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.windowRenderDims.X, g.windowRenderDims.Y
}

func main() {

	game := &Game{
		wordsPos:         util.Pos{X: 0, Y: 0},
		xIncreasing:      true,
		yIncreasing:      true,
		windowSize:       util.Dims{X: 600, Y: 200},
		windowRenderDims: util.Dims{X: 300, Y: 100},
		testCard:         game.MakeCard(game.Five, game.Club),
	}
	game.Init()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
