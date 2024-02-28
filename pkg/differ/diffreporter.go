package differ

import (
	"fmt"
	"github.com/Streamline-CV/streamline-cv/api"
	"github.com/gonvenience/ytbx"
	"github.com/homeport/dyff/pkg/dyff"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
)

type DiffReporter struct {
	Source []byte
	Target []byte
}

func (d *DiffReporter) diff() (*api.ChangeReport, error) {
	sourceFile, err := getFile(d.Source)
	if err != nil {
		return nil, err
	}
	targetFile, err := getFile(d.Target)
	if err != nil {
		return nil, err
	}
	report, err := dyff.CompareInputFiles(*sourceFile, *targetFile)
	if err != nil {
		return nil, fmt.Errorf("failed to compare files: %w", err)
	}
	var allChanges []api.Change
	for _, diff := range report.Diffs {
		var path []string
		for _, pathElement := range diff.Path.PathElements {
			path = append(path, pathElement.Name)
		}
		if path == nil {
			path = []string{}
		}
		for _, detail := range diff.Details {
			changeType := getChangeType(detail.Kind)
			var changes []api.Change
			switch changeType {
			case api.Add:
				for _, value := range getChangeValue(detail.To) {
					changes = append(changes, api.Change{
						Path:       path,
						Source:     nil,
						Target:     value,
						ChangeType: changeType,
					})
				}
			case api.Remove:
				for _, value := range getChangeValue(detail.From) {
					changes = append(changes, api.Change{
						Path:       path,
						Source:     value,
						Target:     nil,
						ChangeType: changeType,
					})
				}
			case api.Modify, api.OrderChange:
				changes = append(changes, api.Change{
					Path:       path,
					Source:     getChangeValue(detail.From)[0],
					Target:     getChangeValue(detail.To)[0],
					ChangeType: changeType,
				})
			}

			for _, change := range changes {
				allChanges = append(allChanges, change)
			}
		}
	}

	return &api.ChangeReport{
		Changes: allChanges,
	}, nil
}

func getFile(data []byte) (*ytbx.InputFile, error) {
	filePath, err := getTmpFile(data)
	if err != nil {
		return nil, err
	}
	file, err := ytbx.LoadFile(filePath)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func getTmpFile(data []byte) (string, error) {
	tempFile, err := os.CreateTemp("", "temp")
	if err != nil {
		log.Fatal().Msgf("Failed to create a temporary file %e", err)
		return "", err
	}
	_, err = tempFile.Write(data)
	if err != nil {
		log.Fatal().Msgf("Failed to write to a temporary file %e", err)
		return "", err
	}

	if err := tempFile.Close(); err != nil {
		log.Fatal().Msgf("Failed to close a temporary file %e", err)
		return "", err
	}
	return tempFile.Name(), nil
}

func getChangeValue(node *yaml.Node) []*api.ChangeValue {
	switch node.Kind {
	case yaml.SequenceNode:
		var changes []*api.ChangeValue
		for _, content := range node.Content {
			for _, subchange := range getChangeValue(content) {
				changes = append(changes, subchange)
			}
		}
		return changes
	case yaml.MappingNode:
		valNode := node.Content[1]
		return []*api.ChangeValue{{
			Name:        node.Content[0].Value,
			Value:       valNode.Value,
			Line:        valNode.Line,
			ColumnStart: valNode.Column,
			ColumnEnd:   valNode.Column + len(valNode.Value),
		},
		}
	default:
		return []*api.ChangeValue{{
			Value:       node.Value,
			Line:        node.Line,
			ColumnStart: node.Column,
			ColumnEnd:   node.Column + len(node.Value),
		},
		}
	}
}

func getChangeType(code rune) api.ChangeType {
	switch code {
	case dyff.ADDITION:
		return api.Add
	case dyff.REMOVAL:
		return api.Remove
	case dyff.MODIFICATION:
		return api.Modify
	case dyff.ORDERCHANGE:
		return api.OrderChange
	default:
		return api.Unknown
	}
}
