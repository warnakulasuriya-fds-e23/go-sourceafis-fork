package transparency

import "github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/primitives"

type ConsistentSkeletonRidge struct {
	Start, End int
	Points     primitives.List[primitives.IntPoint]
}
