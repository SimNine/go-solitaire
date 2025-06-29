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
	BaseVelocity   float64

	currVelocity float64
}

func (a *Animation) UnitDelta() util.Pos[float64] {
	fullDelta := a.TargetPos.Sub(a.StartingPos)
	currDelta := a.TargetPos.Sub(a.CurrPos())

	// Set the current velocity based on the ratio of the current delta to the full delta, scaled by the base velocity
	if fullDelta.X != 0 {
		a.currVelocity = (currDelta.X / fullDelta.X) * a.BaseVelocity
	} else if fullDelta.Y != 0 {
		a.currVelocity = (currDelta.Y / fullDelta.Y) * a.BaseVelocity
	} else {
		a.currVelocity = 0 // No movement if both deltas are zero
	}

	// Compute the step size based on the current velocity and the full delta
	return util.Pos[float64]{
		X: fullDelta.X * a.currVelocity,
		Y: fullDelta.Y * a.currVelocity,
	}
}
