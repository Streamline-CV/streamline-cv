package differ

import (
	"github.com/Streamline-CV/streamline-cv/api"
	"github.com/google/go-cmp/cmp"
	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDiffReportDiffChangesContentExtended(t *testing.T) {
	testCases := []struct {
		name            string
		sourceYAML      string
		targetYAML      string
		expectedChanges []api.Change // Assume this includes fields for path, source value, and target value
	}{
		{
			name: "Nested Addition",
			sourceYAML: `
a:
  b:
    c: 1`,
			targetYAML: `
a:
  b:
    c: 1
    d: 2`,
			expectedChanges: []api.Change{
				{
					Path:       []string{"a", "b"},
					ChangeType: api.Add,
					Source:     nil,
					Target: &api.ChangeValue{
						Name:        "d",
						Value:       "2",
						Line:        5,
						ColumnStart: 8,
						ColumnEnd:   9,
					},
				},
			},
		},
		{
			name: "List Modification",
			sourceYAML: `
a:
  - 1
  - 2`,
			targetYAML: `
a:
  - 2
  - 3`,
			expectedChanges: []api.Change{
				{
					Path:       []string{"a"},
					ChangeType: api.Remove,
					Source: &api.ChangeValue{
						Value:       "1",
						Line:        3,
						ColumnStart: 5,
						ColumnEnd:   6,
					},
					Target: nil,
				},
				{
					Path:       []string{"a"},
					ChangeType: api.Add,
					Source:     nil,
					Target: &api.ChangeValue{
						Value:       "3",
						Line:        4,
						ColumnStart: 5,
						ColumnEnd:   6,
					},
				},
			},
		},
		{
			name: "Complex Structure Change",
			sourceYAML: `
a:
  b:
    - c: 1
    - d: 2
e: 3`,
			targetYAML: `
a:
  b:
    - c: 2
f: 4`,
			expectedChanges: []api.Change{
				{
					Path: []string{},
					Source: &api.ChangeValue{
						Name:        "e",
						Value:       "3",
						Line:        6,
						ColumnStart: 4,
						ColumnEnd:   5,
					},
					Target:     nil,
					ChangeType: api.Remove,
				},
				{
					Path:   []string{},
					Source: nil,
					Target: &api.ChangeValue{
						Name:        "f",
						Value:       "4",
						Line:        5,
						ColumnStart: 4,
						ColumnEnd:   5,
					},
					ChangeType: api.Add,
				},
				{
					Path: []string{"a", "b"},
					Source: &api.ChangeValue{
						Name:        "c",
						Value:       "1",
						Line:        4,
						ColumnStart: 10,
						ColumnEnd:   11,
					},
					Target:     nil,
					ChangeType: api.Remove,
				},
				{
					Path: []string{"a", "b"},
					Source: &api.ChangeValue{
						Name:        "d",
						Value:       "2",
						Line:        5,
						ColumnStart: 10,
						ColumnEnd:   11,
					},
					Target:     nil,
					ChangeType: api.Remove,
				},
				{
					Path:   []string{"a", "b"},
					Source: nil,
					Target: &api.ChangeValue{
						Name:        "c",
						Value:       "2",
						Line:        4,
						ColumnStart: 10,
						ColumnEnd:   11,
					},
					ChangeType: api.Add,
				},
			},
		},
		{
			name: "Simple List Addition",
			sourceYAML: `
a:
  - b
  - c`,
			targetYAML: `
a:
  - b
  - c
  - d`,
			expectedChanges: []api.Change{
				{
					Path:       []string{"a"},
					ChangeType: api.Add,
					Source:     nil,
					Target: &api.ChangeValue{
						Value:       "d",
						Line:        5,
						ColumnStart: 5,
						ColumnEnd:   6,
					},
				},
			},
		},
		{
			name: "Nested Object Value Change",
			sourceYAML: `
config:
  database:
    host: db.example.com
    port: 5432
    username: user1`,
			targetYAML: `
config:
  database:
    host: db.example.com
    port: 5432
    username: user2`,
			expectedChanges: []api.Change{
				{
					Path:       []string{"config", "database", "username"},
					ChangeType: api.Modify,
					Source: &api.ChangeValue{
						Value:       "user1",
						Line:        6,
						ColumnStart: 15,
						ColumnEnd:   20,
					},
					Target: &api.ChangeValue{
						Value:       "user2",
						Line:        6,
						ColumnStart: 15,
						ColumnEnd:   20,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reporter := DiffReporter{
				Source: []byte(dedent.Dedent(tc.sourceYAML)),
				Target: []byte(dedent.Dedent(tc.targetYAML)),
			}

			report, err := reporter.diff()

			assert.NoError(t, err, "Error should not occur when diffing YAML content")

			assert.NotNil(t, report, "Report should not be nil")

			if diff := cmp.Diff(tc.expectedChanges, report.Changes); diff != "" {
				t.Errorf("DiffReporter.diff() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
