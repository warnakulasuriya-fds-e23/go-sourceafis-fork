package minutiae

import (
	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/features"
	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/primitives"
)

func Collect(ridges, valleys *features.Skeleton) (*primitives.GenericList[*features.FeatureMinutia], error) {
	minutiae := primitives.NewGenericList[*features.FeatureMinutia]()
	if err := collect(minutiae, ridges, features.ENDING); err != nil {
		return nil, err
	}
	if err := collect(minutiae, valleys, features.BIFURCATION); err != nil {
		return nil, err
	}

	return minutiae, nil
}

func collect(minutiae *primitives.GenericList[*features.FeatureMinutia], skeleton *features.Skeleton, t features.MinutiaType) error {
	for _, minutia := range skeleton.Minutiae {
		if len(minutia.Ridges) == 1 {
			dir, err := minutia.Ridges[0].Direction()
			if err != nil {
				return err
			}
			minutiae.PushBack(features.NewFeatureMinutia(minutia.Position, dir, t))
		}
	}

	return nil
}
