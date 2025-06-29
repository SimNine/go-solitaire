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
	fullDelta := a.TargetPos.Sub(a.StartingPos)
	return util.Pos[float64]{
		X: fullDelta.X / 100.0,
		Y: fullDelta.Y / 100.0,
	}
}
