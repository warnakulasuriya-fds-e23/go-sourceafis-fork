package transparency

import (
	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/features"
	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/primitives"
)

type ConsistentSkeleton struct {
	Width, Height int
	Minutiae      []primitives.IntPoint
	Ridges        []*ConsistentSkeletonRidge
}

func NewConsistentSkeleton(skeleton *features.Skeleton) *ConsistentSkeleton {
	offsets := make(map[*features.SkeletonMinutia]int)
	var positions []primitives.IntPoint
	var ridges []*ConsistentSkeletonRidge
	var i int
	for _, minutia := range skeleton.Minutiae {
		offsets[minutia] = i
		positions = append(positions, minutia.Position)
		i++
	}

	for _, minutia := range skeleton.Minutiae {
		for _, r := range minutia.Ridges {
			if _, ok := r.Points.(*primitives.CircularList[primitives.IntPoint]); ok {
				ridge := &ConsistentSkeletonRidge{
					Start:  offsets[r.Start()],
					End:    offsets[r.End()],
					Points: r.Points,
				}
				ridges = append(ridges, ridge)
			}
		}

	}

	return &ConsistentSkeleton{
		Width:    skeleton.Size.X,
		Height:   skeleton.Size.Y,
		Minutiae: positions,
		Ridges:   ridges,
	}
}
