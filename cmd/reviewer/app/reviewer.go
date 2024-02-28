package app

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"os"
	"streamline-cv/pkg/assistant"
	"streamline-cv/pkg/differ"
	"streamline-cv/pkg/reporting"
)

func CreateReview(configFile string, reviewResultFile string) error {
	log.Info().Msgf("Running cli for file %s", configFile)
	changes, err := differ.GetDiff(configFile, log.Logger)
	if err != nil {
		log.Fatal().Msgf("Failed to get diff: %e", err)
	}
	aiAssistant, err := assistant.NewAiAssistant(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal().Msgf("Failed to create ai assistant: %e", err)
	}
	refactoring, err := aiAssistant.Refactor(changes)
	if err != nil {
		log.Fatal().Msgf("Failed getting refactoring: %e", err)
	}
	rdf, err := reporting.ToRdf(*refactoring)
	if err != nil {
		log.Fatal().Msgf("Failed formatting to rdf: %s", err)
	}
	jsonData, err := json.Marshal(rdf)
	if err != nil {
		log.Fatal().Msgf("Failed writing rdf to json: %s", err)
	}
	err = os.WriteFile(reviewResultFile, jsonData, 0644)
	if err != nil {
		log.Fatal().Msgf("Failed to write file: %e", err)
	}
	return nil
}
