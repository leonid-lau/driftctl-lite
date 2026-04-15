package output

import (
	"fmt"
	"strings"
)

// FormatFlag implements the pflag.Value interface so Format can be used
// directly as a cobra/pflag flag value.
type FormatFlag struct {
	Value Format
}

// String returns the current format value as a string.
func (f *FormatFlag) String() string {
	return string(f.Value)
}

// Set validates and assigns the format from a raw string flag value.
func (f *FormatFlag) Set(s string) error {
	switch Format(strings.ToLower(s)) {
	case FormatText:
		f.Value = FormatText
	case FormatJSON:
		f.Value = FormatJSON
	default:
		return fmt.Errorf("invalid output format %q: must be one of [text, json]", s)
	}
	return nil
}

// Type returns the type name used in help text.
func (f *FormatFlag) Type() string {
	return "format"
}

// DefaultFormatFlag returns a FormatFlag pre-set to the text format.
func DefaultFormatFlag() *FormatFlag {
	return &FormatFlag{Value: FormatText}
}
