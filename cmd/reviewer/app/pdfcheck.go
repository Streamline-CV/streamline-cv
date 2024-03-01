package app

import (
	"encoding/json"
	"fmt"
	"github.com/Streamline-CV/streamline-cv/api"
	"github.com/Streamline-CV/streamline-cv/pkg/reporting"
	"github.com/rs/zerolog/log"
	"github.com/unidoc/unipdf/v3/model"
	"os"
)

func PdfCheck(inputFile string, outputFile string) error {
	reader, _, err := model.NewPdfReaderFromFile(inputFile, nil)
	if err != nil {
		return err
	}

	// Get the number of pages
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
