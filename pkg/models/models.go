package models

import (
	"encoding/json"
	"strings"
)

type SodaReport struct {
	GCPProject              string  `json:"gcpProject"`
	Dataset                 string  `json:"dataset"`
	SlackChannel            string  `json:"slackChannel"`
	SlackNotifyOnPassedScan *string `json:"slackNotifyOnScanPassed"`
	DockerImage             string  `json:"dockerImage"`

	Results     []TestResult `json:"testResults"`
	ConfigError *string      `json:"configError"`
}

type TestResult struct {
	ID                 string         `json:"id"`
	Table              string         `json:"table"`
	Test               string         `json:"test"`
	Outcome            string         `json:"outcome"`
	Definition         string         `json:"definition"`
	Metrics            []string       `json:"metrics,omitempty"`
	MetricValues       map[string]any `json:"metricValues,omitempty"`
	ResourceAttributes []string       `json:"resourceAttributes,omitempty"`
	Time               string         `json:"time"`
	Column             string         `json:"column"`
	Type               string         `json:"type"`
	Filter             *string        `json:"filter"`
}

// rawTestResult mirrors TestResult but keeps Metrics as raw JSON so we can
// handle both v3 (array of strings) and v4 (map of metric name → value).
type rawTestResult struct {
	ID                 string          `json:"id"`
	Table              string          `json:"table"`
	Test               string          `json:"test"`
	Outcome            string          `json:"outcome"`
	Definition         string          `json:"definition"`
	Metrics            json.RawMessage `json:"metrics"`
	ResourceAttributes []string        `json:"resourceAttributes"`
	Time               string          `json:"time"`
	Column             string          `json:"column"`
	Type               string          `json:"type"`
	Filter             *string         `json:"filter"`
}

func (t *TestResult) UnmarshalJSON(data []byte) error {
	var raw rawTestResult
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	t.ID = raw.ID
	t.Table = raw.Table
	t.Test = raw.Test
	t.Outcome = normalizeOutcome(raw.Outcome)
	t.Definition = raw.Definition
	t.ResourceAttributes = raw.ResourceAttributes
	t.Time = raw.Time
	t.Column = raw.Column
	t.Type = raw.Type
	t.Filter = raw.Filter

	if len(raw.Metrics) > 0 {
		// v3 sends metrics as an array of strings: ["missing_count"]
		var metricNames []string
		if err := json.Unmarshal(raw.Metrics, &metricNames); err == nil {
			t.Metrics = metricNames
		} else {
			// v4 sends metrics as a map of name → value: {"missing_count": 42}
			var metricValues map[string]any
			if err := json.Unmarshal(raw.Metrics, &metricValues); err == nil {
				t.MetricValues = metricValues
				for k := range metricValues {
					t.Metrics = append(t.Metrics, k)
				}
			}
		}
	}

	return nil
}

// normalizeOutcome maps both v3 lowercase ("pass", "fail", "warn") and
// v4 uppercase enum names ("PASSED", "FAILED", "WARN") to canonical lowercase.
func normalizeOutcome(outcome string) string {
	switch strings.ToUpper(outcome) {
	case "PASS", "PASSED":
		return "pass"
	case "FAIL", "FAILED":
		return "fail"
	case "WARN", "WARNED":
		return "warn"
	case "ERROR":
		return "error"
	default:
		return strings.ToLower(outcome)
	}
}
