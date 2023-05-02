package models

type SodaTest struct {
	GCPProject   string `json:"gcpProject"`
	Dataset      string `json:"dataset"`
	SlackChannel string `json:"slackChannel"`

	Results []TestResult `json:"testResults"`
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
	Filter             []string `json:"filter"`
}

type BigqueryRow struct {
	ID                 string   `json:"id"`
	Project            string   `json:"project"`
	Dataset            string   `json:"dataset"`
	Table              string   `json:"table"`
	Test               string   `json:"test"`
	Outcome            string   `json:"outcome"`
	Definition         string   `json:"definition"`
	Metrics            []string `json:"metrics"`
	ResourceAttributes []string `json:"resourceAttributes"`
	Time               string   `json:"time"`
	Column             string   `json:"column"`
	Type               string   `json:"type"`
	Filter             []string `json:"filter"`
}
