package game

import (
	"image/color"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/text/language"
	"urffer.xyz/go-solitaire/src/util"
)

const DEFAULT_CARD_WIDTH = 50
const DEFAULT_CARD_HEIGHT = 70

const DEFAULT_NUMBER_SIZE = 20.0
const DEFAULT_SUIT_SIZE = 40.0

var numberTextface *text.GoTextFace = nil
var suitTextface *text.GoTextFace = nil

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

	// Load font file
	reader, err := os.Open("unifont-16.0.04.otf")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	// Convert the font file to a GoTextFaceSource
	font, err := text.NewGoTextFaceSource(reader)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new text face if not already set
	if numberTextface == nil {
		numberTextface = &text.GoTextFace{
			Source:    font,
			Direction: text.DirectionLeftToRight,
			Size:      DEFAULT_NUMBER_SIZE,
			Language:  language.English,
		}
	}
	if suitTextface == nil {
		suitTextface = &text.GoTextFace{
			Source:    font,
			Direction: text.DirectionLeftToRight,
			Size:      DEFAULT_SUIT_SIZE,
			Language:  language.English,
		}
	}

	// Draw the number on the card in each corner
	numberSymbol := NumberSymbols[number]
	textOps := &text.DrawOptions{}
	text.Draw(
		image,
		numberSymbol,
		numberTextface,
		textOps,
	)
	textOps.GeoM.Rotate(math.Pi)
	textOps.GeoM.Translate(
		float64(DEFAULT_CARD_WIDTH),
		float64(DEFAULT_CARD_HEIGHT),
	)
	text.Draw(
		image,
		numberSymbol,
		numberTextface,
		textOps,
	)

	// Draw the suit in the center of the card
	suitSymbol := SuitSymbols[suit]
	suitOps := &text.DrawOptions{}
	suitOps.GeoM.Translate(
		float64(DEFAULT_CARD_WIDTH)/2-suitTextface.Size/4,
		float64(DEFAULT_CARD_HEIGHT)/2-suitTextface.Size/2,
	)
	text.Draw(
		image,
		suitSymbol,
		suitTextface,
		suitOps,
	)

	// Return the card with the complete image
	return &Card{
		Number: number,
		Suit:   suit,
		image:  image,
		pos:    util.Pos{X: 50, Y: 100},
	}
}

type Card struct {
	Number Number
	Suit   Suit

	image *ebiten.Image
	pos   util.Pos
}

func (c *Card) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.pos.X), float64(c.pos.Y))
	screen.DrawImage(c.image, op)
}
