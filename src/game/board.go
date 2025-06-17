package game

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"urffer.xyz/go-solitaire/src/util"
)

func NewBoard() *Board {
	sampleCard := MakeCard(Ace, Heart)
	sampleCard.IsShown = false

	// Create a deck of cards
	deck := []Card{}
	for _, suit := range []Suit{Heart, Diamond, Club, Spade} {
		for _, number := range []Number{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King} {
			deck = append(deck, *MakeCard(number, suit))
		}
	}

	// Shuffle the deck
	for i := len(deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}

	// Distribute the cards into the working stacks
	workingStacks := [7][]Card{}
	for i := 0; i < 7; i++ {
		currStack := &workingStacks[i]
		for j := 0; j <= i; j++ {
			var card Card
			card, deck = deck[0], deck[1:]
			card.pos = util.Pos{
				X: 10 + i*(DEFAULT_CARD_WIDTH+10),
				Y: 10 + j*15,
			}
			if j == i {
				card.IsShown = true // Only the last card in each stack is shown
			} else {
				card.IsShown = false // The rest are face down
			}
			*currStack = append(*currStack, card)
		}
	}

	return &Board{
		suitPiles:      make(map[Suit][]Card),
		workingStacks:  workingStacks,
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
	for _, stack := range b.workingStacks {
		for _, card := range stack {
			card.Draw(screen)
		}
	}
}
