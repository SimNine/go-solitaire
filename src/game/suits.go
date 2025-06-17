package game

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
