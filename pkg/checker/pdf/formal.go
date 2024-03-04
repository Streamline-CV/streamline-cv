package pdf

import (
	"fmt"
	"github.com/Streamline-CV/streamline-cv/api"
	"github.com/Streamline-CV/streamline-cv/pkg/checker"
	"github.com/unidoc/unipdf/v3/model"
	"strings"
)

type PageCountChecker struct{}

func (p *PageCountChecker) Applies(filePath string) bool {
	return strings.HasSuffix(filePath, ".pdf")
}

func (p *PageCountChecker) RunCheck(filePath string) ([]api.Check, error) {
	reader, _, err := model.NewPdfReaderFromFile(filePath, nil)
	if err != nil {
		return []api.Check{}, err
	}

	numPages, err := reader.GetNumPages()
	if err != nil {
		return []api.Check{}, err
	}
	var checks []api.Check
	if numPages == 1 {
		checks = append(checks, api.Check{
			CheckId:  "PdfPageCountCheck",
			Message:  fmt.Sprintf("Pdf has exactly one page, all good!"),
			Severity: api.INFO,
		})
	} else if numPages > 1 {
		checks = append(checks, api.Check{
			CheckId:  "PdfPageCountCheck",
			Message:  fmt.Sprintf("Pdf has %d pages. It's recommended to keep it short and make it fit one page.", numPages),
			Severity: api.WARN,
		})
	} else {
		checks = append(checks, api.Check{
			CheckId:  "PdfPageCountCheck",
			Message:  fmt.Sprintf("Pdf is empty, probably something went wrong with pdf generation."),
			Severity: api.ERROR,
		})
	}
	return checks, nil
}

func newPageCountChecker() *PageCountChecker {
	return &PageCountChecker{}
}

func init() {
	checker.Registry.Register(newPageCountChecker())
}
