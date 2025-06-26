package game

import (
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/text/language"
	"urffer.xyz/go-solitaire/src/util"
)

const DEFAULT_CARD_WIDTH = 120
const DEFAULT_CARD_HEIGHT = 168

const DEFAULT_NUMBER_SIZE = 30.0
const DEFAULT_SUIT_SIZE = 60.0

var numberTextface *text.GoTextFace = nil
var suitTextface *text.GoTextFace = nil
var cardBackImage *ebiten.Image = nil
var cardBlankImage *ebiten.Image = nil

func InitCardsAssets() {
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
	cardBackImage, err = util.LoadEbitenImageFromFile("assets/card_back.png")
	if err != nil {
		log.Fatal(err)
	}

	// Scale the card back image to fit the default card dimensions
	cardBackImage = util.ScaleEbitenImage(
		cardBackImage,
		util.Dims{X: DEFAULT_CARD_WIDTH, Y: DEFAULT_CARD_HEIGHT},
	)

	// Load the blank card image
	cardBlankImage, err = util.LoadEbitenImageFromFile("assets/card_blank.png")
	if err != nil {
		log.Fatal(err)
	}

	// Scale the blank card image to fit the default card dimensions
	cardBlankImage = util.ScaleEbitenImage(
		cardBlankImage,
		util.Dims{X: DEFAULT_CARD_WIDTH, Y: DEFAULT_CARD_HEIGHT},
	)

	// Load suit images
	suitImagePaths := map[Suit]string{
		Heart:   "assets/suit_heart.png",
		Diamond: "assets/suit_diamond.png",
		Club:    "assets/suit_club.png",
		Spade:   "assets/suit_spade.png",
	}
	for suit, imagePath := range suitImagePaths {
		image, err := util.LoadEbitenImageFromFile(imagePath)
		if err != nil {
			log.Fatalf("Failed to load suit image for %s: %v", suit, err)
		}
		SuitImages[suit] = image
	}

	// Load number images
	numberImagePaths := map[Number]string{
		Ace:   "assets/num_1-ace.png",
		Two:   "assets/num_2.png",
		Three: "assets/num_3.png",
		Four:  "assets/num_4.png",
		Five:  "assets/num_5.png",
	}
	for number, imagePath := range numberImagePaths {
		image, err := util.LoadEbitenImageFromFile(imagePath)
		if err != nil {
			log.Fatalf("Failed to load number image for %s: %v", number, err)
		}
		// Scale the number image to fit the default card dimensions
		image = util.ScaleEbitenImage(
			image,
			util.Dims{X: DEFAULT_CARD_WIDTH, Y: DEFAULT_CARD_HEIGHT},
		)
		NumberImages[number] = image
	}

	// Generate remaining number images
	for _, number := range []Number{
		Six, Seven, Eight, Nine, Ten, Jack, Queen, King,
	} {
		numberSymbol := NumberSymbols[number]
		numberOps := &text.DrawOptions{}
		numberOps.GeoM.Scale(4.0, 4.0)
		newImage := ebiten.NewImage(DEFAULT_CARD_WIDTH, DEFAULT_CARD_HEIGHT)
		text.Draw(
			newImage,
			numberSymbol,
			numberTextface,
			numberOps,
		)
		NumberImages[number] = newImage
	}
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
	image := ebiten.NewImageFromImage(cardBlankImage)

	// Draw the suit in the center of the card
	suitImage := SuitImages[suit]
	suitOps := &ebiten.DrawImageOptions{}
	if suit == Heart || suit == Diamond {
		// Red suits
		suitOps.ColorScale.SetR(1.0)
		suitOps.ColorScale.SetG(0.0)
		suitOps.ColorScale.SetB(0.0)
	} else {
		// Black suits
		suitOps.ColorScale.SetR(0.0)
		suitOps.ColorScale.SetG(0.0)
		suitOps.ColorScale.SetB(0.0)
	}
	suitImage = util.ScaleEbitenImage(
		suitImage,
		util.Dims{X: suitImage.Bounds().Dx() / 4, Y: suitImage.Bounds().Dy() / 4},
	)
	suitOps.GeoM.Translate(
		float64(DEFAULT_CARD_WIDTH)/2.0-float64(suitImage.Bounds().Dx())/2.0,
		float64(DEFAULT_CARD_HEIGHT)/2.0-float64(suitTextface.Size)/2.0,
	)
	image.DrawImage(suitImage, suitOps)

	// Draw the number on the card in each corner
	if numberImage := NumberImages[number]; numberImage != nil {
		numberOps := &ebiten.DrawImageOptions{}
		if suit == Heart || suit == Diamond {
			// Red suits
			numberOps.ColorScale.SetR(1.0)
			numberOps.ColorScale.SetG(0.0)
			numberOps.ColorScale.SetB(0.0)
		} else {
			// Black suits
			numberOps.ColorScale.SetR(0.0)
			numberOps.ColorScale.SetG(0.0)
			numberOps.ColorScale.SetB(0.0)
		}
		numberImage = util.ScaleEbitenImage(
			numberImage,
			util.Dims{X: numberImage.Bounds().Dx() / 4, Y: numberImage.Bounds().Dy() / 4},
		)
		image.DrawImage(numberImage, numberOps)
		numberOps.GeoM.Rotate(math.Pi)
		numberOps.GeoM.Translate(
			float64(DEFAULT_CARD_WIDTH),
			float64(DEFAULT_CARD_HEIGHT),
		)
		image.DrawImage(numberImage, numberOps)
	} else {
		log.Fatalf("No image found for number %s", number)
	}

	// Return the card with the complete image
	return &Card{
		Number:  number,
		Suit:    suit,
		IsShown: true,
		image:   image,

		pos: util.Pos{X: 50, Y: 100},
	}
}

func (c *Card) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.pos.X), float64(c.pos.Y))
	if !c.IsShown {
		// If the card is not shown, draw the card back
		screen.DrawImage(cardBackImage, op)
	} else {
		screen.DrawImage(c.image, op)
	}
}

func (c *Card) String() string {
	return NumberSymbols[c.Number] + SuitSymbols[c.Suit]
}

func (c *Card) Contains(pos util.Pos) bool {
	return pos.X >= c.pos.X && pos.X <= c.pos.X+DEFAULT_CARD_WIDTH &&
		pos.Y >= c.pos.Y && pos.Y <= c.pos.Y+DEFAULT_CARD_HEIGHT
}
