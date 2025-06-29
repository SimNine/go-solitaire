package game

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"urffer.xyz/go-solitaire/src/util"
)

var placeholderImage *ebiten.Image = nil

func InitCardStackBkg() {
	placeholderImage = ebiten.NewImage(DEFAULT_CARD_WIDTH, DEFAULT_CARD_HEIGHT)
	placeholderImage.Fill(color.RGBA{
		R: 0,
		G: 150,
		B: 0,
		A: 255,
	})
}

type CardStack struct {
	Cards    []*Card
	basePos  util.Pos
	isSpread bool
}

func (c *CardStack) GetTopCard() *Card {
	if len(c.Cards) == 0 {
		return nil
	}
	return c.Cards[len(c.Cards)-1]
}

func (c *CardStack) Draw(screen *ebiten.Image) {
	if c.isSpread {
		for _, card := range c.Cards {
			card.Draw(screen)
		}
	} else {
		if len(c.Cards) > 0 {
			// Draw only the top card of the stack
			topCard := c.Cards[len(c.Cards)-1]
			topCard.Draw(screen)
		}
	}

	if len(c.Cards) == 0 {
		// Draw a placeholder for the base of the stack
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(c.basePos.ToFloatTuple())
		screen.DrawImage(placeholderImage, op)
	}
}

func (c *CardStack) TranslateTo(pos util.Pos) {
	// Translate the base position of the stack to a new position
	c.basePos = pos

	// Fix the positions of all cards in the stack
	c.repositionCards()
}

func (c *CardStack) SetSpread(spread bool) {
	// Set the spread state of the stack
	c.isSpread = spread

	// Reposition cards based on the new spread state
	c.repositionCards()
}

func (c *CardStack) SetAllShown(shown bool) {
	// Set the visibility of all cards in the stack
	for _, card := range c.Cards {
		card.IsShown = shown
	}
}

func (c *CardStack) AppendStack(other *CardStack) {
	if other == nil || len(other.Cards) == 0 {
		return // Nothing to append
	}

	// Append the cards from the other stack to this stack
	c.Cards = append(c.Cards, other.Cards...)

	// Fix the positions of all cards in the stack
	c.repositionCards()
}

func (c *CardStack) AppendCard(card *Card) {
	if card == nil {
		return // Nothing to append
	}

	// Append the card to the stack
	c.Cards = append(c.Cards, card)

	// Fix the position of the new card
	c.repositionCards()
}

func (c *CardStack) Reverse() {
	// Reverse the order of cards in the stack
	for i, j := 0, len(c.Cards)-1; i < j; i, j = i+1, j-1 {
		c.Cards[i], c.Cards[j] = c.Cards[j], c.Cards[i]
	}

	// Reposition cards after reversing
	c.repositionCards()
}

func (c *CardStack) Shuffle() {
	for i := len(c.Cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		c.Cards[i], c.Cards[j] = c.Cards[j], c.Cards[i]
	}

	// Reposition cards after shuffling
	c.repositionCards()
}

func (c *CardStack) repositionCards() {
	// Reposition all cards in the stack based on the base position
	for i, card := range c.Cards {
		if c.isSpread {
			card.pos = c.basePos.Translate(0, i*DEFAULT_CARD_INTERPILE_SPACING)
		} else {
			card.pos = c.basePos
		}
	}
}

func (c *CardStack) splitDeckAtIndex(index int) *CardStack {
	if index < 0 || index >= len(c.Cards) {
		log.Println("Invalid index for splitting deck:", index)
		return nil // Invalid index
	}
	// Create a new stack with the cards from this index to the end
	newStack := &CardStack{
		Cards:    c.Cards[index:],
		isSpread: c.isSpread,
		basePos:  c.Cards[index].pos,
	}
	// Update this stack to only contain the cards before this index
	c.Cards = c.Cards[:index]
	return newStack
}

func (c *CardStack) SplitDeckAtPos(pos util.Pos) *CardStack {
	// Find the index of the card that contains the given position
	for i, card := range c.Cards {
		if card.Contains(pos) {
			if i < len(c.Cards)-1 {
				// Split deck at the current card if it's not the last one
				if !c.Cards[i+1].Contains(pos) {
					return c.splitDeckAtIndex(i)
				}
			} else {
				// Split deck on the last card
				return c.splitDeckAtIndex(i)
			}
		}
	}
	return nil // No card found at the given position
}

func (c *CardStack) BaseCardContains(pos util.Pos) bool {
	// Check if the base position of the stack contains the given position
	return pos.X >= c.basePos.X && pos.X <= c.basePos.X+DEFAULT_CARD_WIDTH &&
		pos.Y >= c.basePos.Y && pos.Y <= c.basePos.Y+DEFAULT_CARD_HEIGHT
}
