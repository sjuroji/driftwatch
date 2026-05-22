// Package output provides formatters that render drift reports in
// different output formats (plain text, JSON). Use output.New to
// obtain a Formatter for a given Format constant, then call Write
// with any io.Writer and a drift.Report to produce the output.
//
// Supported formats:
//
//	FormatText  - human-readable tabular text (default CLI output)
//	FormatJSON  - machine-readable JSON suitable for CI pipelines
package output
