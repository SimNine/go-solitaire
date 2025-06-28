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
			card.IsShown = false
			*currStack = append(*currStack, card)
		}
		workingStacks[i] = &CardStack{
			Cards:    *currStack,
			IsSpread: true,
			BasePos: POS_DRAW_PILE.Translate(
				i*(DEFAULT_CARD_SPACING+DEFAULT_CARD_WIDTH),
				DEFAULT_CARD_HEIGHT+DEFAULT_CARD_SPACING,
			),
		}
		workingStacks[i].GetTopCard().IsShown = true
		workingStacks[i].BasePos = util.Pos{
			X: DEFAULT_CARD_SPACING + i*(DEFAULT_CARD_WIDTH+DEFAULT_CARD_SPACING),
			Y: DEFAULT_CARD_SPACING + DEFAULT_CARD_HEIGHT + DEFAULT_CARD_SPACING,
		}
		workingStacks[i].repositionCards()
	}

	// Put the rest of the deck into the draw pile
	drawPileCards := deck
	for _, card := range drawPileCards {
		card.pos = POS_DRAW_PILE
		card.IsShown = false
	}
	drawPile := &CardStack{
		Cards:    drawPileCards,
		IsSpread: false,
		BasePos:  POS_DRAW_PILE,
	}

	// Create empty suit piles
	suitPiles := [4]*CardStack{}
	for i := 0; i < 4; i++ {
		suitPiles[i] = &CardStack{
			Cards:    []*Card{},
			IsSpread: false, // Suit piles only show the top card
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
			Cards:    []*Card{},
			IsSpread: false,
			BasePos:  POS_OVERTURNED_PILE,
		},
	}
}

type Board struct {
	suitPiles      [4]*CardStack
	workingStacks  [7]*CardStack
	drawPile       *CardStack
	overturnedPile *CardStack

	heldCardStack      *CardStack
	heldCardResetStack *CardStack
	heldCardOffset     util.Pos

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

	// Draw the held card stack if it exists
	if b.heldCardStack != nil {
		b.heldCardStack.Draw(screen)
	}
}

func (b *Board) SetCusrorPos(pos util.Pos) {
	b.cursorPos = pos

	if b.heldCardStack != nil {
		b.heldCardStack.TranslateTo(b.cursorPos.TranslatePos(b.heldCardOffset))
	}
}

func (b *Board) MouseDown() {
	// Try picking cards up from one of the working stacks
	for _, stack := range b.workingStacks {
		if newStack := stack.SplitDeckAtPos(b.cursorPos); newStack != nil {
			if !newStack.Cards[0].IsShown {
				log.Println("Cannot pick up a stack of cards where the bottom card is not shown.")
				stack.AppendStack(newStack)
			} else {
				log.Println("Sub-stack picked up")
				b.heldCardStack = newStack
				b.heldCardResetStack = stack
				b.heldCardOffset = b.heldCardStack.BasePos.Sub(b.cursorPos)
			}
			return
		}
	}

	// Try picking cards up from one of the suit piles
	for _, stack := range b.suitPiles {
		if newStack := stack.SplitDeckAtPos(b.cursorPos); newStack != nil {
			log.Println("Card grabbed from suit pile:", newStack)
			b.heldCardStack = newStack
			b.heldCardResetStack = stack
			b.heldCardOffset = b.heldCardStack.BasePos.Sub(b.cursorPos)
			return
		}
	}

	// Try picking a card up from the draw pile
	if b.drawPile.BaseCardContains(b.cursorPos) {
		if topCard := b.drawPile.GetTopCard(); topCard != nil {
			log.Println("Card grabbed from draw pile")
			if newStack := b.drawPile.SplitDeckAtPos(b.cursorPos); newStack != nil {
				b.heldCardStack = newStack
				b.heldCardStack.Cards[0].IsShown = true
				b.heldCardResetStack = b.overturnedPile
				b.heldCardOffset = b.heldCardStack.BasePos.Sub(b.cursorPos)
				return
			}
		} else if topCard == nil {
			if len(b.overturnedPile.Cards) > 0 {
				b.overturnedPile.Reverse()
				replenishStack := b.overturnedPile.splitDeckAtIndex(0)
				b.drawPile.AppendStack(replenishStack)
				for _, card := range b.drawPile.Cards {
					card.IsShown = false
				}
			}
		}
	}

	// Try picking a card up from the overturned pile
	if topCard := b.overturnedPile.GetTopCard(); topCard != nil && topCard.Contains(b.cursorPos) {
		log.Println("Card grabbed from overturned pile")
		if newStack := b.overturnedPile.SplitDeckAtPos(b.cursorPos); newStack != nil {
			b.heldCardStack = newStack
			b.heldCardResetStack = b.overturnedPile
			b.heldCardOffset = b.heldCardStack.BasePos.Sub(b.cursorPos)
			return
		}
	}

	log.Println("No card grabbed, cursor not over a working stack or no cards available.")
}

func (b *Board) MouseUp() {
	if b.heldCardStack == nil {
		log.Println("No card held, ignoring mouse up event.")
		return
	}

	// Check if the held stack can be placed onto a working stack
	for _, stack := range b.workingStacks {
		topCard := stack.GetTopCard()

		log.Println("Checking if held stack can be placed onto working stack:", stack)

		// If the stack is empty, only a stack with a king as bottom card can be placed on it
		if topCard == nil {
			if b.heldCardStack.Cards[0].Number == King {
				// See if the stack contains the cursor position. If so, append stacks
				if (&Card{pos: stack.BasePos}).Contains(b.cursorPos) {
					log.Println("Card dropped onto working stack:", stack)
					stack.AppendStack(b.heldCardStack)
					if newTopCard := b.heldCardResetStack.GetTopCard(); newTopCard != nil {
						newTopCard.IsShown = true
					}
					b.heldCardStack = nil
					b.heldCardResetStack = nil
					return
				}
			}
		} else {
			if topCard.Suit.IsOppositeColor(b.heldCardStack.Cards[0].Suit) &&
				b.heldCardStack.Cards[0].Number.IsOneLessThan(topCard.Number) {
				if topCard.Contains(b.cursorPos) {
					log.Println("Card dropped onto working stack:", stack)
					stack.AppendStack(b.heldCardStack)
					if newTopCard := b.heldCardResetStack.GetTopCard(); newTopCard != nil {
						newTopCard.IsShown = true
					}
					b.heldCardStack = nil
					b.heldCardResetStack = nil
					return
				}
			}
		}
	}

	// Check if the held stack can be placed onto a suit pile
	for i, stack := range b.suitPiles {
		topCard := stack.GetTopCard()

		// If the held stack has more than one card, it cannot be placed onto a suit pile
		if len(b.heldCardStack.Cards) > 1 {
			log.Println("Held stack has more than one card, cannot be placed onto suit pile")
			continue
		}

		log.Println("Checking if held stack can be placed onto suit pile: ", i)

		// If the stack is empty, only an ace can be placed on it
		cardCanBePlaced := false
		if topCard == nil {
			if b.heldCardStack.Cards[0].Number == Ace {
				if (&Card{pos: stack.BasePos}).Contains(b.cursorPos) {
					cardCanBePlaced = true
					log.Println("Card dropped onto suit pile:", stack)
				}
			}
		} else {
			if b.heldCardStack.Cards[0].Suit == topCard.Suit &&
				b.heldCardStack.Cards[0].Number.IsOneMoreThan(topCard.Number) {
				if topCard.Contains(b.cursorPos) {
					cardCanBePlaced = true
				}
			}
		}

		// If the card can be placed, do stuff
		if cardCanBePlaced {
			log.Println("Card dropped onto suit pile:", stack)
			stack.AppendStack(b.heldCardStack)
			if newTopCard := b.heldCardResetStack.GetTopCard(); newTopCard != nil {
				newTopCard.IsShown = true
			}
			b.heldCardStack = nil
			b.heldCardResetStack = nil
			return
		}
	}

	// No stack was dropped onto, so reset the held stack
	log.Println("No stack found to drop the held card onto, resetting held card stack.")
	b.heldCardResetStack.AppendStack(b.heldCardStack)

	// Remove held card stack state
	b.heldCardStack = nil
	b.heldCardResetStack = nil
}
