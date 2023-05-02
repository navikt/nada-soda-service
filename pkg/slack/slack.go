package slack

import (
	"fmt"

	"github.com/navikt/nada-soda-service/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type SlackClient struct {
	slackToken string
	log        *logrus.Entry
}

type testDiscrepancies struct {
	Errors   []models.TestResult
	Warnings []models.TestResult
}

func New(slackToken string, log *logrus.Entry) *SlackClient {
	return &SlackClient{
		slackToken: slackToken,
		log:        log,
	}
}

func (s *SlackClient) NotifyOnTestErrors(sodaTest models.SodaTest) error {
	if sodaTest.SlackChannel == "" {
		s.log.Warningf("no slack channel provided for dataset %v.%v", sodaTest.GCPProject, sodaTest.Dataset)
		return nil
	}
	if hasDiscrepancies, discrepancies := s.findDiscrepancies(sodaTest.Results); hasDiscrepancies {
		if err := s.postNotification(discrepancies, sodaTest.GCPProject, sodaTest.Dataset, sodaTest.SlackChannel); err != nil {
			return err
		}
	}

	return nil
}

func (s *SlackClient) findDiscrepancies(sodaResults []models.TestResult) (bool, testDiscrepancies) {
	discrepancies := testDiscrepancies{}
	for _, r := range sodaResults {
		switch r.Outcome {
		case "pass":
			continue
		case "fail":
			discrepancies.Errors = append(discrepancies.Errors, r)
		default:
			discrepancies.Warnings = append(discrepancies.Warnings, r)
		}
	}

	return len(discrepancies.Errors) > 0 || len(discrepancies.Warnings) > 0, discrepancies
}

func (s *SlackClient) postNotification(d testDiscrepancies, projectID, dataset, slackChannel string) error {
	slackClient := slack.New(s.slackToken)

	topSection, attachments := s.createMessage(d, projectID, dataset)
	_, _, err := slackClient.PostMessage(slackChannel, slack.MsgOptionBlocks(topSection), slack.MsgOptionAttachments(attachments...))
	if err != nil {
		return err
	}

	return nil
}

func (s *SlackClient) createMessage(d testDiscrepancies, projectID, dataset string) (slack.Block, []slack.Attachment) {
	topMessage := slack.TextBlockObject{
		Type:  "plain_text",
		Text:  "Varsel om datakvalitetsavvik :gasp:",
		Emoji: true,
	}

	topSection := slack.NewSectionBlock(&topMessage, nil, nil)
	attachments := []slack.Attachment{}

	if len(d.Errors) > 0 {
		message := ""
		for _, e := range d.Errors {
			line1 := ""
			if e.Column != "" {
				line1 = fmt.Sprintf("_*Tabell: %v*_ _*kolonne: %v*_\n", e.Table, e.Column)
			} else {
				line1 = fmt.Sprintf("_*Tabell: %v*_\n", e.Table)
			}
			line2 := e.Test + "\n"
			message = message + line1 + line2
		}
		attachments = append(attachments, slack.Attachment{
			Color:      "#ff2d00",
			AuthorName: fmt.Sprintf("%v.%v", projectID, dataset),
			Title:      "Tester med feil",
			Text:       message,
			Footer:     "SODA Bot",
			Fallback:   "Tester med feil",
		})
	}

	if len(d.Warnings) > 0 {
		message := ""
		for _, e := range d.Errors {
			line1 := ""
			if e.Column != "" {
				line1 = fmt.Sprintf("_*Tabell: %v*_ _*kolonne: %v*_\n", e.Table, e.Column)
			} else {
				line1 = fmt.Sprintf("_*Tabell: %v*_\n", e.Table)
			}
			line2 := e.Test + "\n"
			message = message + line1 + line2
		}
		attachments = append(attachments, slack.Attachment{
			Color:      "#ffa500",
			AuthorName: fmt.Sprintf("%v.%v", projectID, dataset),
			Title:      "Tester med varslinger",
			Text:       message,
			Footer:     "SODA Bot",
			Fallback:   "Tester med varslinger",
		})
	}

	return topSection, attachments
}
