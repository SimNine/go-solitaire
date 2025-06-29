package animation

import (
	"urffer.xyz/go-solitaire/src/util"
)

type Animation struct {
	StartingPos    util.Pos[float64]
	TargetPos      util.Pos[float64]
	CurrPos        func() util.Pos[float64]
	OnFinishAction func()
	Update         func()
}

func (a *Animation) UnitDelta() util.Pos[float64] {
	delta := a.TargetPos.Sub(a.CurrPos())
	if delta.X > 0 {
		delta.X = 1
	} else if delta.X < 0 {
		delta.X = -1
	}

	if delta.Y > 0 {
		delta.Y = 1
	} else if delta.Y < 0 {
		delta.Y = -1
	}

	return delta
}
