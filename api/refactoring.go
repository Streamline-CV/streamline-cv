package api

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
	Severity    string
}

type CheckReporting struct {
	Checks []Check
}

type Check struct {
	CheckId  string
	Message  string
	Severity string
}
