package transparency

import "github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/features"

type ConsistentHashEntry struct {
	Key   int
	Edges []*features.IndexedEdge
}
