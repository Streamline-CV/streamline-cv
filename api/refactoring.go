package api

type Severity string

const (
	INFO  Severity = "Info"
	WARN  Severity = "Warning"
	ERROR Severity = "Error"
)

type SuggestionReporting struct {
	Suggestions []Suggestion
}

type Suggestion struct {
	Path        []string
	Line        int
	ColumnStart int
	ColumnEnd   int
	Value       string
	Comment     string
	Severity    Severity
}

type CheckReporting struct {
	Checks []Check
}

type Check struct {
	CheckId  string
	Message  string
	Severity Severity
}
