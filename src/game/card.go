package game

import (
	"image"
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
var cardbackImage *ebiten.Image = nil

func InitCards() {
	// Load font file
	reader, err := os.Open("assets/unifont-16.0.04.otf")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	// Convert the font file to a GoTextFaceSource
	font, err := text.NewGoTextFaceSource(reader)
	if err != nil {
		log.Fatal(err)
	}

	// Create the textfaces
	numberTextface = &text.GoTextFace{
		Source:    font,
		Direction: text.DirectionLeftToRight,
		Size:      DEFAULT_NUMBER_SIZE,
		Language:  language.English,
	}
	suitTextface = &text.GoTextFace{
		Source:    font,
		Direction: text.DirectionLeftToRight,
		Size:      DEFAULT_SUIT_SIZE,
		Language:  language.English,
	}

	// Load the card back image
	reader, err = os.Open("assets/card_back.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	// Decode the image
	image, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the image to an ebiten.Image
	cardbackImage = ebiten.NewImageFromImage(image)

	// Scale the card back image to fit the default card dimensions
	ops := &ebiten.DrawImageOptions{}
	bounds := cardbackImage.Bounds().Size()
	xRatio := float64(bounds.X) / float64(DEFAULT_CARD_WIDTH)
	yRatio := float64(bounds.Y) / float64(DEFAULT_CARD_HEIGHT)
	ops.GeoM.Scale(1/xRatio, 1/yRatio)
	cardbackImageCanvas := ebiten.NewImage(DEFAULT_CARD_WIDTH, DEFAULT_CARD_HEIGHT)
	cardbackImageCanvas.DrawImage(
		cardbackImage,
		ops,
	)

	// Set the cardbackImage to the scaled image
	cardbackImage = cardbackImageCanvas
}

type Card struct {
	Number  Number
	Suit    Suit
	IsShown bool

	image *ebiten.Image
	pos   util.Pos
}

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
		Number:  number,
		Suit:    suit,
		IsShown: true,
		image:   image,
		pos:     util.Pos{X: 50, Y: 100},
	}
}

func (c *Card) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.pos.X), float64(c.pos.Y))
	if !c.IsShown {
		// If the card is not shown, draw the card back
		screen.DrawImage(cardbackImage, op)
	} else {
		screen.DrawImage(c.image, op)
	}
}
