package templates

import (
	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/features"
	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/primitives"
)

type FeatureTemplate struct {
	Size     primitives.IntPoint
	Minutiae *primitives.GenericList[*features.FeatureMinutia]
}

func NewFeatureTemplate(size primitives.IntPoint, minutiae *primitives.GenericList[*features.FeatureMinutia]) *FeatureTemplate {
	return &FeatureTemplate{
		Size:     size,
		Minutiae: minutiae,
	}
}
