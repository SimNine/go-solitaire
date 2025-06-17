package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"urffer.xyz/go-solitaire/src/util"
)

const DEFAULT_CARD_WIDTH = 50
const DEFAULT_CARD_HEIGHT = 70

func MakeCard(
	number Number,
	suit Suit,
) *Card {

	image := ebiten.NewImage(DEFAULT_CARD_WIDTH, DEFAULT_CARD_HEIGHT)
	image.Fill(color.RGBA{
		R: 0,
		G: 100,
		B: 0,
		A: 255,
	})

	return &Card{
		Number: number,
		Suit:   suit,
		image:  image,
		pos:    util.Pos{X: 0, Y: 0},
	}
}

type Card struct {
	Number Number
	Suit   Suit

	image *ebiten.Image
	pos   util.Pos
}

func (c *Card) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{
		R: 0,
		G: 100,
		B: 0,
		A: 255,
	})
}
