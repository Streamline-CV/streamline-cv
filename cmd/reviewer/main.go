package main

import (
	"github.com/Streamline-CV/streamline-cv/cmd/reviewer/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

var configFileName, rdfOutputFileName string

var reviewCmd = &cobra.Command{
	Use:   "reviewer",
	Short: "Streamline CV reviewer",
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		if configFileName == "" || rdfOutputFileName == "" {
			log.Fatal().Msg("You must specify both a CV config file and an RDF output file.")
		}
		err := app.CreateReview(configFileName, rdfOutputFileName)
		if err != nil {
			log.Fatal().Msgf("Failed doing review: %e", err)
		}
	},
}

func init() {
	reviewCmd.Flags().StringVarP(&configFileName, "config", "c", "config.yaml", "The path to the CV config file")
	reviewCmd.Flags().StringVarP(&rdfOutputFileName, "outfile", "o", "rdf.json", "The name of the RDF output file")
}

func main() {
	if err := reviewCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Failed to execute command")
	}
}
