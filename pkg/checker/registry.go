package checker

import (
	"github.com/Streamline-CV/streamline-cv/api"
	"github.com/rs/zerolog/log"
	"reflect"
)

type CheckRegistry struct {
	checkers []Checker
}

var Registry = newCheckRegistry()

func newCheckRegistry() *CheckRegistry {
	return &CheckRegistry{}
}

func (r *CheckRegistry) Register(check Checker) {
	r.checkers = append(r.checkers, check)
}

func (r *CheckRegistry) RunAllChecks(filePath string) (*api.CheckReporting, error) {
	var allChecks []api.Check
	log.Info().Msgf("Checking file %s.", filePath)
	for _, checker := range r.checkers {
		if checker.Applies(filePath) {
			log.Info().Msgf("Running check %s.", reflect.TypeOf(checker))
			checks, err := checker.RunCheck(filePath)
			if err != nil {
				return nil, err
			}
			for _, check := range checks {
				log.Info().Msgf("Check result is: %s.", check)
				allChecks = append(allChecks, check)
			}
		}
	}
	return &api.CheckReporting{
		Checks: allChecks,
	}, nil
}
