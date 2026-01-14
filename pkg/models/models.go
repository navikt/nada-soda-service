package models

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
	ID                 string   `json:"id"`
	Table              string   `json:"table"`
	Test               string   `json:"test"`
	Outcome            string   `json:"outcome"`
	Definition         string   `json:"definition"`
	Metrics            []string `json:"metrics"`
	ResourceAttributes []string `json:"resourceAttributes"`
	Time               string   `json:"time"`
	Column             string   `json:"column"`
	Type               string   `json:"type"`
	Filter             string   `json:"filter"`
}
