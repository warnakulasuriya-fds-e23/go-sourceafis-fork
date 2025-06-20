package matcher

import (
	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/features"
	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/templates"
)

type Probe struct {
	template *templates.SearchTemplate
	hash     map[int][]*features.IndexedEdge
}

func NewProbe(template *templates.SearchTemplate, hash map[int][]*features.IndexedEdge) *Probe {
	return &Probe{
		template: template,
		hash:     hash,
	}
}
