package differ

import (
	"github.com/rs/zerolog"
	"streamline-cv/api"
	"streamline-cv/internal/git"
)

func GetDiff(filename string, logger zerolog.Logger) (*api.ChangeReport, error) {
	logger.Info().Msgf("Getting diff for local and main")
	fileState, err := git.GetFileChange(filename)
	if err != nil {
		return nil, err
	}
	reporter := DiffReporter{
		Source: fileState.Source,
		Target: fileState.Target,
	}
	diff, err := reporter.diff()
	if err != nil {
		return nil, err
	}
	return diff, nil
}
