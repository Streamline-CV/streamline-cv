package reporting

import (
	"github.com/Streamline-CV/streamline-cv/api"
	"github.com/reviewdog/reviewdog/proto/rdf"
)

func ToRdf(refactoring api.SuggestionReporting) (*rdf.DiagnosticResult, error) {
	result := rdf.DiagnosticResult{
		Severity: rdf.Severity_WARNING,
		Source: &rdf.Source{
			Name: "Streamline AI assistant",
		},
	}

	if len(refactoring.Suggestions) == 0 {
		return &result, nil
	}

	for _, suggestion := range refactoring.Suggestions {
		suggestionRange := &rdf.Range{
			Start: &rdf.Position{Line: int32(suggestion.Line), Column: int32(suggestion.ColumnStart)},
			End:   &rdf.Position{Line: int32(suggestion.Line), Column: int32(suggestion.ColumnEnd)},
		}
		diagnostic := rdf.Diagnostic{
			Severity: mapSeverity(suggestion.Severity),
			Location: &rdf.Location{
				Path:  "config.yaml",
				Range: suggestionRange,
			},
			Message: suggestion.Comment,
			Suggestions: []*rdf.Suggestion{
				{
					Range: suggestionRange,
					Text:  suggestion.Value,
				},
			},
		}

		result.Diagnostics = append(result.Diagnostics, &diagnostic)
	}

	return &result, nil
}

func ChecksToRdf(checkReporting api.CheckReporting) (*rdf.DiagnosticResult, error) {
	result := rdf.DiagnosticResult{
		Severity: rdf.Severity_INFO,
		Source: &rdf.Source{
			Name: "Streamline AI assistant",
		},
	}

	if len(checkReporting.Checks) == 0 {
		return &result, nil
	}
	var failed = false
	for _, check := range checkReporting.Checks {
		severity := mapSeverity(check.Severity)
		if severity == rdf.Severity_ERROR {
			failed = true
		}
		diagnostic := rdf.Diagnostic{
			Severity: severity,
			Message:  check.Message,
			Location: &rdf.Location{
				Path: "CV.pdf",
			},
		}

		result.Diagnostics = append(result.Diagnostics, &diagnostic)
	}

	if failed {
		result.Severity = rdf.Severity_ERROR
	}

	return &result, nil
}

func mapSeverity(severity string) rdf.Severity {
	switch severity {
	case "INFO":
		return rdf.Severity_INFO
	case "WARNING":
		return rdf.Severity_WARNING
	case "ERROR":
		return rdf.Severity_ERROR
	default:
		return rdf.Severity_UNKNOWN_SEVERITY
	}
}
