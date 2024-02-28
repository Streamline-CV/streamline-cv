package differ

import (
	"github.com/Streamline-CV/streamline-cv/api"
	"github.com/Streamline-CV/streamline-cv/internal/git"
	"github.com/rs/zerolog"
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
