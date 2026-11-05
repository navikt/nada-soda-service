package slack

import (
	"testing"

	"github.com/navikt/nada-soda-service/pkg/models"
	"github.com/slack-go/slack"
)

func TestCreatePassedScanMessage(t *testing.T) {
	tests := []struct {
		name    string
		results []models.TestResult
		want    string
	}{
		{
			name:    "one table",
			results: []models.TestResult{{Table: "users"}},
			want:    "Tabell: users",
		},
		{
			name: "multiple tables",
			results: []models.TestResult{
				{Table: "users"},
				{Table: "orders"},
			},
			want: "Tabeller: users, orders",
		},
		{
			name: "duplicate tables retain first-seen order",
			results: []models.TestResult{
				{Table: "users"},
				{Table: "orders"},
				{Table: "users"},
			},
			want: "Tabeller: users, orders",
		},
		{
			name: "empty table names are skipped",
			results: []models.TestResult{
				{Table: ""},
				{Table: "users"},
				{Table: ""},
			},
			want: "Tabell: users",
		},
		{
			name:    "no table names",
			results: []models.TestResult{{Table: ""}},
			want:    "Tabell ikke oppgitt",
		},
		{
			name:    "empty report",
			results: nil,
			want:    "Tabell ikke oppgitt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			topSection, attachments := (*Client)(nil).createPassedScanMessage("project", "dataset", tt.results)

			section, ok := topSection.(*slack.SectionBlock)
			if !ok {
				t.Fatalf("top section has type %T, want *slack.SectionBlock", topSection)
			}
			if got := section.Text.Text; got != "Datakvalitetssjekker bestått" {
				t.Errorf("top text = %q, want %q", got, "Datakvalitetssjekker bestått")
			}
			if section.Text.Emoji != nil {
				t.Errorf("top text emoji = %v, want nil", *section.Text.Emoji)
			}
			if len(attachments) != 1 {
				t.Fatalf("got %d attachments, want 1", len(attachments))
			}
			attachment := attachments[0]
			if attachment.AuthorName != "project.dataset" {
				t.Errorf("author name = %q, want %q", attachment.AuthorName, "project.dataset")
			}
			if attachment.Color != "#00ff00" {
				t.Errorf("attachment color = %q, want %q", attachment.Color, "#00ff00")
			}
			if attachment.Text != tt.want {
				t.Errorf("attachment text = %q, want %q", attachment.Text, tt.want)
			}
		})
	}
}
