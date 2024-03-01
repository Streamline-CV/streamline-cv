package main

import (
	"fmt"
	"github.com/Streamline-CV/streamline-cv/cmd/reviewer/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

var inputFile, outputFile string

var rootCmd = &cobra.Command{
	Use:   "streamlinecv",
	Short: "Streamline CV",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from the root command")
	},
}

var reviewCmd = &cobra.Command{
	Use:   "reviewer",
	Short: "Streamline CV reviewer",
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		if inputFile == "" || outputFile == "" {
			log.Fatal().Msg("You must specify both a CV config file and an RDF output file.")
		}
		err := app.CreateReview(inputFile, outputFile)
		if err != nil {
			log.Fatal().Msgf("Failed doing review: %e", err)
		}
	},
}

var pdfCheckCmd = &cobra.Command{
	Use:   "pdfcheck",
	Short: "Streamline CV pdf checker",
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		if inputFile == "" || outputFile == "" {
			log.Fatal().Msg("You must specify both a CV config file and an RDF output file.")
		}
		log.Info().Msgf("Running pdf check")
		err := app.PdfCheck(inputFile, outputFile)
		if err != nil {
			log.Fatal().Msgf("Failed doing check: %e", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(reviewCmd)
	reviewCmd.Flags().StringVarP(&inputFile, "config", "c", "config.yaml", "The path to the CV config file")
	reviewCmd.Flags().StringVarP(&outputFile, "outfile", "o", "rdf.json", "The name of the RDF output file")

	rootCmd.AddCommand(pdfCheckCmd)
	pdfCheckCmd.Flags().StringVarP(&inputFile, "config", "c", "cv.pdf", "The path to the CV config file")
	pdfCheckCmd.Flags().StringVarP(&outputFile, "outfile", "o", "checks-rdf.json", "The name of the RDF output file")

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Failed to execute command")
	}
}
