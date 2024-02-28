package api

type FileState struct {
	Source []byte
	Target []byte
}

type ChangeValue struct {
	Name        string
	Value       string
	Line        int
	ColumnStart int
	ColumnEnd   int
}

type ChangeType string

const (
	Add         ChangeType = "Add"
	Remove      ChangeType = "Remove"
	Modify      ChangeType = "Modify"
	OrderChange ChangeType = "OrderChange"
	Unknown     ChangeType = "Unknown"
)

type Change struct {
	Path       []string
	Source     *ChangeValue
	Target     *ChangeValue
	ChangeType ChangeType
}

type ChangeReport struct {
	Changes []Change
}
