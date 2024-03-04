package checker

import "github.com/Streamline-CV/streamline-cv/api"

type CheckTarget string

type Checker interface {
	Applies(filePath string) bool
	RunCheck(filePath string) ([]api.Check, error)
}
