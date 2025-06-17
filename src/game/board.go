package game

type Board struct {
	suitPiles      map[Suit]([]Card)
	workingStacks  [7][]Card
	drawPile       []Card
	overturnedPile []Card
}
