package main

import (
	"encoding/json"
	"github.com/Streamline-CV/streamline-cv/pkg/assistant"
	"github.com/Streamline-CV/streamline-cv/pkg/differ"
	"github.com/Streamline-CV/streamline-cv/pkg/reporting"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
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

func init() {
	var inputFile, outputFile string

	var reviewCmd = &cobra.Command{
		Use:   "reviewer",
		Short: "Streamline CV reviewer",
		Run: func(cmd *cobra.Command, args []string) {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
			if inputFile == "" || outputFile == "" {
				log.Fatal().Msg("You must specify both a CV config file and an RDF output file.")
			}
			err := CreateReview(inputFile, outputFile)
			if err != nil {
				log.Fatal().Msgf("Failed doing review: %e", err)
			}
		},
	}
	reviewCmd.Flags().StringVarP(&inputFile, "config", "c", "config.yaml", "The path to the CV config file")
	reviewCmd.Flags().StringVarP(&outputFile, "outfile", "o", "rdf.json", "The name of the RDF output file")

	rootCmd.AddCommand(reviewCmd)
}
