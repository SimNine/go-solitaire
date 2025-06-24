package game

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"urffer.xyz/go-solitaire/src/util"
)

const DEFAULT_CARD_SPACING = 10
const DEFAULT_CARD_INTERPILE_SPACING = 20

var POS_DRAW_PILE = util.Pos{
	X: DEFAULT_CARD_SPACING,
	Y: DEFAULT_CARD_SPACING,
}
var POS_OVERTURNED_PILE = POS_DRAW_PILE.Translate(
	DEFAULT_CARD_WIDTH+DEFAULT_CARD_SPACING,
	0,
)

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
	workingStacks := [7]*CardStack{}
	for i := 0; i < 7; i++ {
		currStack := &[]*Card{}
		for j := 0; j <= i; j++ {
			var card *Card
			card, deck = deck[0], deck[1:]
			card.pos = util.Pos{
				X: DEFAULT_CARD_SPACING + i*(DEFAULT_CARD_WIDTH+DEFAULT_CARD_SPACING),
				Y: DEFAULT_CARD_SPACING + DEFAULT_CARD_HEIGHT + DEFAULT_CARD_SPACING + j*DEFAULT_CARD_INTERPILE_SPACING,
			}
			if j == i {
				card.IsShown = true // Only the last card in each stack is shown
			} else {
				card.IsShown = false // The rest are face down
			}
			*currStack = append(*currStack, card)
		}
		workingStacks[i] = &CardStack{
			Cards:     *currStack,
			RenderAll: true,
			BasePos: POS_DRAW_PILE.Translate(
				i*(DEFAULT_CARD_SPACING+DEFAULT_CARD_WIDTH),
				DEFAULT_CARD_HEIGHT+DEFAULT_CARD_SPACING,
			),
		}
	}

	// Put the rest of the deck into the draw pile
	drawPileCards := deck
	for _, card := range drawPileCards {
		card.pos = POS_DRAW_PILE
		card.IsShown = false
	}
	drawPile := &CardStack{
		Cards:     drawPileCards,
		RenderAll: false,
		BasePos:   POS_DRAW_PILE,
	}

	// Create empty suit piles
	suitPiles := make(map[Suit]*CardStack)
	for i, suit := range []Suit{Heart, Diamond, Club, Spade} {
		suitPiles[suit] = &CardStack{
			Cards:     []*Card{},
			RenderAll: false, // Suit piles only show the top card
			BasePos: POS_OVERTURNED_PILE.Translate(
				(2+i)*(DEFAULT_CARD_WIDTH+DEFAULT_CARD_SPACING),
				0,
			),
		}
	}

	// Create the board with suit piles, working stacks, and empty draw and overturned piles
	return &Board{
		suitPiles:     suitPiles,
		workingStacks: workingStacks,
		drawPile:      drawPile,
		overturnedPile: &CardStack{
			Cards:     []*Card{},
			RenderAll: false,
			BasePos:   POS_OVERTURNED_PILE,
		},
	}
}

type Board struct {
	suitPiles      map[Suit](*CardStack)
	workingStacks  [7]*CardStack
	drawPile       *CardStack
	overturnedPile *CardStack

	heldCard         *Card
	heldCardResetPos util.Pos
	heldCardOffset   util.Pos

	cursorPos util.Pos
}

func (b *Board) Draw(screen *ebiten.Image) {
	// Fill the background with the board color
	screen.Fill(color.RGBA{
		R: 0,
		G: 75,
		B: 0,
		A: 255,
	})

	// Draw all working stacks
	for _, stack := range b.workingStacks {
		stack.Draw(screen)
	}

	// Draw the draw pile
	b.drawPile.Draw(screen)

	// Draw the overturned pile
	b.overturnedPile.Draw(screen)

	// Draw the suit piles
	for _, stack := range b.suitPiles {
		stack.Draw(screen)
	}
}

func (b *Board) SetCusrorPos(pos util.Pos) {
	b.cursorPos = pos

	if b.heldCard != nil {
		b.heldCard.pos = b.cursorPos.TranslatePos(b.heldCardOffset)
	}
}

func (b *Board) GrabCard() {
	for _, stack := range b.workingStacks {
		if len(stack.Cards) == 0 {
			continue // Skip empty stacks
		}
		testCard := stack.Cards[len(stack.Cards)-1]
		if testCard.Contains(b.cursorPos) {
			log.Println("Card grabbed from working stack:", testCard)
			b.heldCard = testCard
			b.heldCardResetPos = b.heldCard.pos
			b.heldCardOffset = b.heldCard.pos.Sub(b.cursorPos)
			return
		}
	}

	for _, stack := range b.suitPiles {
		if len(stack.Cards) == 0 {
			continue // Skip empty suit piles
		}
		testCard := stack.Cards[len(stack.Cards)-1]
		if testCard.Contains(b.cursorPos) {
			log.Println("Card grabbed from suit pile:", testCard)
			b.heldCard = testCard
			b.heldCardResetPos = b.heldCard.pos
			b.heldCardOffset = b.heldCard.pos.Sub(b.cursorPos)
			return
		}
	}

	if len(b.drawPile.Cards) > 0 {
		lastCard := b.drawPile.Cards[len(b.drawPile.Cards)-1]
		if lastCard.Contains(b.cursorPos) {
			log.Println("Card grabbed from draw pile:", lastCard)
			b.heldCard = lastCard
			b.heldCardResetPos = b.heldCard.pos
			b.heldCardOffset = b.heldCard.pos.Sub(b.cursorPos)
			return
		}
	}

	if len(b.overturnedPile.Cards) > 0 {
		lastCard := b.overturnedPile.Cards[len(b.overturnedPile.Cards)-1]
		if lastCard.Contains(b.cursorPos) {
			log.Println("Card grabbed from overturned pile:", lastCard)
			b.heldCard = lastCard
			b.heldCardResetPos = b.heldCard.pos
			b.heldCardOffset = b.heldCard.pos.Sub(b.cursorPos)
			return
		}
	}

	log.Println("No card grabbed, cursor not over a working stack or no cards available.")
}

func (b *Board) ReleaseCard() {
	b.heldCard = nil
}
