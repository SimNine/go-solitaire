package game

type Number int

const (
	Ace   Number = 1
	Two   Number = 2
	Three Number = 3
	Four  Number = 4
	Five  Number = 5
	Six   Number = 6
	Seven Number = 7
	Eight Number = 8
	Nine  Number = 9
	Ten   Number = 10
	Jack  Number = 11
	Queen Number = 12
	King  Number = 13
)

var NumberSymbols = map[Number]string{
	Ace:   "A",
	Two:   "2",
	Three: "3",
	Four:  "4",
	Five:  "5",
	Six:   "6",
	Seven: "7",
	Eight: "8",
	Nine:  "9",
	Ten:   "10",
	Jack:  "J",
	Queen: "Q",
	King:  "K",
}
