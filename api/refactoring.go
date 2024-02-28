package api

type Refactoring struct {
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
