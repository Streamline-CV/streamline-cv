package main

import (
	"encoding/json"
	"fmt"
	"github.com/Streamline-CV/streamline-cv/api"
	"github.com/Streamline-CV/streamline-cv/pkg/reporting"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/unidoc/unipdf/v3/model"
	"os"
)

func CheckPdf(inputFile string, outputFile string) error {
	reader, _, err := model.NewPdfReaderFromFile(inputFile, nil)
	if err != nil {
		return err
	}

	numPages, err := reader.GetNumPages()
	if err != nil {
		return err
	}
	checksRefactoring := api.CheckReporting{
		Checks: []api.Check{
			{
				Message:  fmt.Sprintf("Pdf consist of %d pages", numPages),
				CheckId:  "PdfSizeCheck",
				Severity: "WARNING",
			},
		},
	}
	rdf, err := reporting.ChecksToRdf(checksRefactoring)
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
			err := CheckPdf(inputFile, outputFile)
			if err != nil {
				log.Fatal().Msgf("Failed doing check: %e", err)
			}
		},
	}
	pdfCheckCmd.Flags().StringVarP(&inputFile, "pdf", "c", "cv.pdf", "The path to the CV config file")
	pdfCheckCmd.Flags().StringVarP(&outputFile, "outfile", "o", "checks-rdf.json", "The name of the RDF output file")

	rootCmd.AddCommand(pdfCheckCmd)
}
