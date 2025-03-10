package processor

import (
	"github.com/Jeffail/gabs/v2"
	"github.com/mchmarny/vimp/internal/converter/grype"
	"github.com/mchmarny/vimp/internal/converter/snyk"
	"github.com/mchmarny/vimp/internal/converter/trivy"
	"github.com/mchmarny/vimp/pkg/data"
	"github.com/pkg/errors"
)

// VulnerabilityMapper is a function that converts a source to a list of common vulnerability types.
type VulnerabilityMapper func(c *gabs.Container) ([]*data.Vulnerability, error)

// GetMapper returns a vulnerability converter for the given source format.
func getMapper(format Format) (VulnerabilityMapper, error) {
	switch format {
	case FormatSnykJSON:
		return snyk.Convert, nil
	case FormatTrivyJSON:
		return trivy.Convert, nil
	case FormatGrypeJSON:
		return grype.Convert, nil
	default:
		return nil, errors.Errorf("unimplemented conversion format: %s", format)
	}
}
