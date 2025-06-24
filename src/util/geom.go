package util

type Dims Pos
type Pos struct {
	X int
	Y int
}

func MakePosFromTuple(x, y int) Pos {
	return Pos{
		X: x,
		Y: y,
	}
}

func (p Pos) ToTuple() (int, int) {
	return p.X, p.Y
}

func (p Pos) ToFloatTuple() (float64, float64) {
	return float64(p.X), float64(p.Y)
}

func (p Pos) Translate(dx, dy int) Pos {
	return Pos{
		X: p.X + dx,
		Y: p.Y + dy,
	}
}

func (p Pos) TranslatePos(offset Pos) Pos {
	return p.Translate(offset.X, offset.Y)
}

func (p Pos) Sub(other Pos) Pos {
	return Pos{
		X: p.X - other.X,
		Y: p.Y - other.Y,
	}
}

func (p Pos) Eq(other Pos) bool {
	return p.X == other.X && p.Y == other.Y
}
