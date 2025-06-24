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

func (c *CardStack) Draw(screen *ebiten.Image) {
	// Draw a placeholder for the base of the stack
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(c.BasePos.ToFloatTuple())
	screen.DrawImage(placeholderImage, op)

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
}
