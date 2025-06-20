package inner

import (
	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/config"
	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/features"
	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/primitives"
)

func Apply(minutiae *primitives.GenericList[*features.FeatureMinutia], mask *primitives.BooleanMatrix) {
	for e := minutiae.Front(); e != nil; {
		minutia := e.Value.(*features.FeatureMinutia)

		arrow := primitives.FloatAngle(minutia.Direction).ToVector().Multiply(-config.Config.MaskDisplacement).Round()

		if !mask.GetPointWithFallback(minutia.Position.Plus(arrow), false) {
			next := e.Next()
			minutiae.Remove(e)
			e = next
		} else {
			e = e.Next()
		}
	}

}
