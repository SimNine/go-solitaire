package game

import (
	"image/color"

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
	Cards     []*Card
	RenderAll bool
	BasePos   util.Pos
}

func (c *CardStack) GetTopCard() *Card {
	if len(c.Cards) == 0 {
		return nil
	}
	return c.Cards[len(c.Cards)-1]
}

func (c *CardStack) Draw(screen *ebiten.Image) {
	if c.RenderAll {
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
		op.GeoM.Translate(c.BasePos.ToFloatTuple())
		screen.DrawImage(placeholderImage, op)
	}
}

func (c *CardStack) TranslateTo(pos util.Pos) {
	// Translate the base position of the stack to a new position
	c.BasePos = pos
	for i, card := range c.Cards {
		card.pos = pos.TranslatePos(util.Pos{X: 0, Y: i * DEFAULT_CARD_INTERPILE_SPACING})
	}
}

func (c *CardStack) splitDeckAtIndex(index int) *CardStack {
	if index < 0 || index >= len(c.Cards) {
		return nil // Invalid index
	}
	// Create a new stack with the cards from this index to the end
	newStack := &CardStack{
		Cards:     c.Cards[index:],
		RenderAll: c.RenderAll,
		BasePos:   c.BasePos.Translate(0, index*DEFAULT_CARD_INTERPILE_SPACING),
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
