package processor

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
	"github.com/mchmarny/vimp/internal/parser"
)

const (
	FormatUnknown   Format = iota
	FormatGrypeJSON        // grype JSON format
	FormatTrivyJSON        // trivy JSON format
	FormatSnykJSON         // snyk JSON format

	FormatUnknownName   = "unknown"
	FormatGrypeJSONName = "grype"
	FormatTrivyJSONName = "trivy"
	FormatSnykJSONName  = "snyk"
)

// Format represents the source format.
type Format int64

// String returns the string representation of the source format.
func (f Format) String() string {
	switch f {
	case FormatGrypeJSON:
		return FormatGrypeJSONName
	case FormatTrivyJSON:
		return FormatTrivyJSONName
	case FormatSnykJSON:
		return FormatSnykJSONName
	default:
		return FormatUnknownName
	}
}

// ParseFormat parses the source format.
func ParseFormat(s string) (Format, error) {
	switch s {
	case FormatGrypeJSONName:
		return FormatGrypeJSON, nil
	case FormatTrivyJSONName:
		return FormatTrivyJSON, nil
	case FormatSnykJSONName:
		return FormatSnykJSON, nil
	default:
		return FormatUnknown, fmt.Errorf("unknown format: %s", s)
	}
}

// GetFormats returns the supported source formats.
func GetFormats() []Format {
	return []Format{
		FormatGrypeJSON,
		FormatTrivyJSON,
		FormatSnykJSON,
	}
}

// GetFormatNames returns the names of the supported source formats.
func GetFormatNames() []string {
	return []string{
		FormatGrypeJSONName,
		FormatTrivyJSONName,
		FormatSnykJSONName,
	}
}

func discoverFormat(c *gabs.Container) Format {
	if c == nil {
		return FormatUnknown
	}

	// grype
	d := c.Search("descriptor", "name")
	if d.Exists() && parser.ToString(d.Data()) == "grype" {
		return FormatGrypeJSON
	}

	// trivy
	if c.ExistsP("SchemaVersion") && c.ExistsP("Results") {
		return FormatTrivyJSON
	}

	// snyk
	if c.Search("vulnerabilities").Exists() && c.Search("applications").Exists() {
		return FormatSnykJSON
	}

	return FormatUnknown
}
