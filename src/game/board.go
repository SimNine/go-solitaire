package game

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"urffer.xyz/go-solitaire/src/animation"
	"urffer.xyz/go-solitaire/src/util"
)

const DEFAULT_CARD_SPACING = 10
const DEFAULT_CARD_INTERPILE_SPACING = 20

var POS_DRAW_PILE = util.Pos[float64]{
	X: DEFAULT_CARD_SPACING,
	Y: DEFAULT_CARD_SPACING,
}
var POS_OVERTURNED_PILE = POS_DRAW_PILE.Translate(
	DEFAULT_CARD_WIDTH+DEFAULT_CARD_SPACING,
	0,
)

func NewBoard() *Board {
	// Create a deck of cards
	deck := &CardStack{
		Cards:    []*Card{},
		basePos:  POS_DRAW_PILE,
		isSpread: false,
	}
	for _, suit := range []Suit{Heart, Diamond, Club, Spade} {
		for _, number := range []Number{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King} {
			deck.AppendCard(MakeCard(number, suit))
		}
	}

	// Shuffle the deck
	deck.Shuffle()

	// Distribute the cards into the working stacks
	workingStacks := [7]*CardStack{}
	for i := 0; i < 7; i++ {
		workingStacks[i] = &CardStack{
			isSpread: true,
			basePos: POS_DRAW_PILE.Translate(
				float64(i*(DEFAULT_CARD_SPACING+DEFAULT_CARD_WIDTH)),
				float64(DEFAULT_CARD_HEIGHT+DEFAULT_CARD_SPACING),
			),
		}
		for j := 0; j <= i; j++ {
			cardstack := deck.splitDeckAtIndex(len(deck.Cards) - 1)
			card := cardstack.GetTopCard()
			card.IsShown = false
			workingStacks[i].AppendCard(card)
		}
		workingStacks[i].GetTopCard().IsShown = true
		workingStacks[i].TranslateTo(util.Pos[float64]{
			X: float64(DEFAULT_CARD_SPACING + i*(DEFAULT_CARD_WIDTH+DEFAULT_CARD_SPACING)),
			Y: float64(DEFAULT_CARD_SPACING + DEFAULT_CARD_HEIGHT + DEFAULT_CARD_SPACING),
		})
	}

	// Put the rest of the deck into the draw pile
	drawPile := deck
	drawPile.TranslateTo(POS_DRAW_PILE)
	drawPile.SetSpread(false)
	drawPile.SetAllShown(false)

	// Create empty suit piles
	suitPiles := [4]*CardStack{}
	for i := 0; i < 4; i++ {
		suitPiles[i] = &CardStack{
			Cards:    []*Card{},
			isSpread: false, // Suit piles only show the top card
			basePos: POS_OVERTURNED_PILE.Translate(
				float64((2+i)*(DEFAULT_CARD_WIDTH+DEFAULT_CARD_SPACING)),
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
			isSpread: false,
			basePos:  POS_OVERTURNED_PILE,
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
	heldCardOffset     util.Pos[int]

	cursorPos util.Pos[int]

	runningAnimation *animation.Animation
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

func (b *Board) Update() {
	if b.runningAnimation != nil {
		b.runningAnimation.Update()
		if b.runningAnimation.CurrPos().AlmostEq(b.runningAnimation.TargetPos, 0.01) {
			b.runningAnimation.OnFinishAction()
			b.runningAnimation = nil
		}
	} else if b.heldCardStack != nil {
		b.heldCardStack.TranslateTo(b.cursorPos.TranslatePos(b.heldCardOffset).ToFloatPos())
	}
}

func (b *Board) SetCusrorPos(pos util.Pos[int]) {
	b.cursorPos = pos
}

func (b *Board) MouseDown() {
	// If there is an ongoing animation, ignore the mouse down event
	if b.runningAnimation != nil {
		log.Println("Ignoring mouse down event, animation is running.")
		return
	}

	// Try picking cards up from one of the working stacks
	for _, stack := range b.workingStacks {
		if newStack := stack.SplitDeckAtPos(b.cursorPos.ToFloatPos()); newStack != nil {
			if !newStack.Cards[0].IsShown {
				log.Println("Cannot pick up a stack of cards where the bottom card is not shown.")
				stack.AppendStack(newStack)
			} else {
				log.Println("Sub-stack picked up")
				b.heldCardStack = newStack
				b.heldCardResetStack = stack
				b.heldCardOffset = b.heldCardStack.basePos.ToIntPos().Sub(b.cursorPos)
			}
			return
		}
	}

	// Try picking cards up from one of the suit piles
	for _, stack := range b.suitPiles {
		if newStack := stack.SplitDeckAtPos(b.cursorPos.ToFloatPos()); newStack != nil {
			log.Println("Card grabbed from suit pile:", newStack)
			b.heldCardStack = newStack
			b.heldCardResetStack = stack
			b.heldCardOffset = b.heldCardStack.basePos.ToIntPos().Sub(b.cursorPos)
			return
		}
	}

	// Try picking a card up from the draw pile
	if b.drawPile.BaseCardContains(b.cursorPos.ToFloatPos()) {
		if topCard := b.drawPile.GetTopCard(); topCard != nil {
			log.Println("Card grabbed from draw pile")
			if newStack := b.drawPile.SplitDeckAtPos(b.cursorPos.ToFloatPos()); newStack != nil {
				b.heldCardStack = newStack
				b.heldCardStack.Cards[0].IsShown = true
				b.heldCardResetStack = b.overturnedPile
				b.heldCardOffset = b.heldCardStack.basePos.ToIntPos().Sub(b.cursorPos)
				return
			}
		} else if topCard == nil {
			if len(b.overturnedPile.Cards) > 0 {
				b.overturnedPile.Reverse()
				replenishStack := b.overturnedPile.splitDeckAtIndex(0)
				b.drawPile.AppendStack(replenishStack)
				b.drawPile.SetAllShown(false)
			}
		}
	}

	// Try picking a card up from the overturned pile
	if topCard := b.overturnedPile.GetTopCard(); topCard != nil && topCard.Contains(b.cursorPos.ToFloatPos()) {
		log.Println("Card grabbed from overturned pile")
		if newStack := b.overturnedPile.SplitDeckAtPos(b.cursorPos.ToFloatPos()); newStack != nil {
			b.heldCardStack = newStack
			b.heldCardResetStack = b.overturnedPile
			b.heldCardOffset = b.heldCardStack.basePos.ToIntPos().Sub(b.cursorPos)
			return
		}
	}

	log.Println("No card grabbed, cursor not over a working stack or no cards available.")
}

func (b *Board) MouseUp() {
	// If there is an ongoing animation, ignore the mouse up event
	if b.runningAnimation != nil {
		log.Println("Ignoring mouse up event, animation is running.")
		return
	}

	// If no card is held, ignore the mouse up event
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
				if (&Card{pos: stack.basePos}).Contains(b.cursorPos.ToFloatPos()) {
					log.Println("Card dropped onto working stack:", stack)
					if newTopCard := b.heldCardResetStack.GetTopCard(); newTopCard != nil {
						newTopCard.IsShown = true
					}
					b.runningAnimation = b.heldCardStack.CreateAnimationToPos(
						stack.GetNextCardPos(),
						func() {
							stack.AppendStack(b.heldCardStack)
							b.heldCardStack = nil
							b.heldCardResetStack = nil
						},
					)
					return
				}
			}
		} else {
			if topCard.Suit.IsOppositeColor(b.heldCardStack.Cards[0].Suit) &&
				b.heldCardStack.Cards[0].Number.IsOneLessThan(topCard.Number) {
				if topCard.Contains(b.cursorPos.ToFloatPos()) {
					log.Println("Card dropped onto working stack:", stack)
					if newTopCard := b.heldCardResetStack.GetTopCard(); newTopCard != nil {
						newTopCard.IsShown = true
					}
					b.runningAnimation = b.heldCardStack.CreateAnimationToPos(
						stack.GetNextCardPos(),
						func() {
							stack.AppendStack(b.heldCardStack)
							b.heldCardStack = nil
							b.heldCardResetStack = nil
						},
					)
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
				if (&Card{pos: stack.basePos}).Contains(b.cursorPos.ToFloatPos()) {
					cardCanBePlaced = true
					log.Println("Card dropped onto suit pile:", stack)
				}
			}
		} else {
			if b.heldCardStack.Cards[0].Suit == topCard.Suit &&
				b.heldCardStack.Cards[0].Number.IsOneMoreThan(topCard.Number) {
				if topCard.Contains(b.cursorPos.ToFloatPos()) {
					cardCanBePlaced = true
				}
			}
		}

		// If the card can be placed, do stuff
		if cardCanBePlaced {
			log.Println("Card dropped onto suit pile:", stack)
			if newTopCard := b.heldCardResetStack.GetTopCard(); newTopCard != nil {
				newTopCard.IsShown = true
			}
			b.runningAnimation = b.heldCardStack.CreateAnimationToPos(
				stack.GetNextCardPos(),
				func() {
					stack.AppendStack(b.heldCardStack)
					b.heldCardStack = nil
					b.heldCardResetStack = nil
				},
			)
			return
		}
	}

	// No stack was dropped onto, so reset the held stack
	log.Println("No stack found to drop the held card onto, resetting held card stack.")
	b.runningAnimation = b.heldCardStack.CreateAnimationToPos(
		b.heldCardResetStack.GetNextCardPos(),
		func() {
			b.heldCardResetStack.AppendStack(b.heldCardStack)
			b.heldCardStack = nil
			b.heldCardResetStack = nil
		},
	)
}
