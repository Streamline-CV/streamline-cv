package main

import (
	"encoding/json"
	"github.com/Streamline-CV/streamline-cv/pkg/checker"
	"github.com/Streamline-CV/streamline-cv/pkg/reporting"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"

	//import checkers to initialize them
	_ "github.com/Streamline-CV/streamline-cv/pkg/checker/pdf"
)

func Check(inputFile string, outputFile string) error {

	checks, err := checker.Registry.RunAllChecks(inputFile)
	rdf, err := reporting.ChecksToRdf(checks)
	jsonData, err := json.Marshal(rdf)
	if err != nil {
		log.Fatal().Msgf("Failed writing rdf to json: %s", err)
	}
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		log.Fatal().Msgf("Failed to write file: %e", err)
	}
	return nil
}

func init() {
	var inputFile, outputFile string

	var pdfCheckCmd = &cobra.Command{
		Use:   "pdfcheck",
		Short: "Streamline CV pdf checker",
		Run: func(cmd *cobra.Command, args []string) {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
			if inputFile == "" || outputFile == "" {
				log.Fatal().Msg("You must specify both a CV config file and an RDF output file.")
			}
			log.Info().Msgf("Running pdf check")
			err := Check(inputFile, outputFile)
			if err != nil {
				log.Fatal().Msgf("Failed doing check: %e", err)
			}
		},
	}
	pdfCheckCmd.Flags().StringVarP(&inputFile, "pdf", "c", "cv.pdf", "The path to the CV config file")
	pdfCheckCmd.Flags().StringVarP(&outputFile, "outfile", "o", "checks-rdf.json", "The name of the RDF output file")

	rootCmd.AddCommand(pdfCheckCmd)
}
