package game

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"urffer.xyz/go-solitaire/src/util"
)

var placeholderImage *ebiten.Image = nil

const DEFAULT_CARD_SPACING = 10

var POS_DRAW_PILE = util.Pos{
	X: DEFAULT_CARD_SPACING,
	Y: DEFAULT_CARD_SPACING,
}

// var POS_OVERTURNED_PILE = util.Pos{

// var POS_OVERTURNED_PILE = util.Pos{
// 	X: DEFAULT_CARD_SPACING + DEFAULT_CARD_WIDTH + DEFAULT_CARD_SPACING,
// 	Y: DEFAULT_CARD_SPACING,
// }

func InitBoardGfx() {
	placeholderImage = ebiten.NewImage(DEFAULT_CARD_WIDTH, DEFAULT_CARD_HEIGHT)
	placeholderImage.Fill(color.RGBA{
		R: 0,
		G: 150,
		B: 0,
		A: 255,
	})
}

func NewBoard() *Board {
	// Create a deck of cards
	deck := []*Card{}
	for _, suit := range []Suit{Heart, Diamond, Club, Spade} {
		for _, number := range []Number{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King} {
			deck = append(deck, MakeCard(number, suit))
		}
	}

	// Shuffle the deck
	for i := len(deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}

	// Distribute the cards into the working stacks
	workingStacks := [7][]*Card{}
	for i := 0; i < 7; i++ {
		currStack := &workingStacks[i]
		for j := 0; j <= i; j++ {
			var card *Card
			card, deck = deck[0], deck[1:]
			card.pos = util.Pos{
				X: DEFAULT_CARD_SPACING + i*(DEFAULT_CARD_WIDTH+DEFAULT_CARD_SPACING),
				Y: DEFAULT_CARD_SPACING + j*DEFAULT_CARD_SPACING + (DEFAULT_CARD_HEIGHT + DEFAULT_CARD_SPACING*2),
			}
			if j == i {
				card.IsShown = true // Only the last card in each stack is shown
			} else {
				card.IsShown = false // The rest are face down
			}
			*currStack = append(*currStack, card)
		}
	}

	// Put the rest of the deck into the draw pile
	drawPile := deck
	for _, card := range drawPile {
		card.pos = POS_DRAW_PILE
		card.IsShown = false
	}

	// Create the board with suit piles, working stacks, and empty draw and overturned piles
	return &Board{
		suitPiles:      make(map[Suit][]*Card),
		workingStacks:  workingStacks,
		drawPile:       drawPile,
		overturnedPile: []*Card{},
	}
}

type Board struct {
	suitPiles      map[Suit]([]*Card)
	workingStacks  [7][]*Card
	drawPile       []*Card
	overturnedPile []*Card
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

	// Render the top card of the draw pile if it exists
	if len(b.drawPile) > 0 {
		topCard := b.drawPile[0]
		topCard.Draw(screen)
	}

	// Render the overturned pile cards
	if len(b.overturnedPile) > 0 {
		// do stuff
	} else {
		// Draw a placeholder image for the overturned pile
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(10+DEFAULT_CARD_WIDTH+10, 10)
		screen.DrawImage(placeholderImage, opts)
	}
}
