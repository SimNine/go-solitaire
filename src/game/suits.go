package game

import "github.com/hajimehoshi/ebiten/v2"

type Suit string

const (
	Spade   Suit = "spade"
	Diamond Suit = "diamond"
	Club    Suit = "club"
	Heart   Suit = "heart"
)

var SuitSymbols = map[Suit]string{
	Spade:   "♠",
	Diamond: "♦",
	Club:    "♣",
	Heart:   "♥",
}

var SuitImages = map[Suit]*ebiten.Image{}

func (s Suit) IsOppositeColor(other Suit) bool {
	switch s {
	case Spade, Club:
		return other == Heart || other == Diamond
	case Heart, Diamond:
		return other == Spade || other == Club
	default:
		return false
	}
}
