package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewBoard() *Board {
	sampleCard := MakeCard(Ace, Heart)
	sampleCard.IsShown = false

	return &Board{
		suitPiles:      make(map[Suit][]Card),
		workingStacks:  [7][]Card{},
		drawPile:       []Card{},
		overturnedPile: []Card{},
		sampleCard:     sampleCard,
	}
}

type Board struct {
	suitPiles      map[Suit]([]Card)
	workingStacks  [7][]Card
	drawPile       []Card
	overturnedPile []Card

	sampleCard *Card
}

func (b *Board) Draw(screen *ebiten.Image) {
	// Fill the background with the board color
	screen.Fill(color.RGBA{
		R: 0,
		G: 75,
		B: 0,
		A: 255,
	})

	// Draw all cards on the screen
	b.sampleCard.Draw(screen)
}
