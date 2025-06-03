package slack

import (
	"fmt"
	"strconv"

	"github.com/navikt/nada-soda-service/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type Client struct {
	token string
	log   *logrus.Entry
}

type testDiscrepancies struct {
	Fails    []models.TestResult
	Warnings []models.TestResult
}

func New(token string, log *logrus.Entry) *Client {
	return &Client{
		token: token,
		log:   log,
	}
}

func (s *Client) Notify(sodaTest models.SodaReport) error {
	if sodaTest.SlackChannel == "" {
		return fmt.Errorf("no Slack channel provided for dataset %v.%v", sodaTest.GCPProject, sodaTest.Dataset)
	}

	if hasDiscrepancies, discrepancies := s.findDiscrepancies(sodaTest.ConfigError, sodaTest.Results); hasDiscrepancies {
		topSection, attachments := s.createDiscrepancyMessage(sodaTest.ConfigError, discrepancies, sodaTest.GCPProject, sodaTest.Dataset)
		if err := s.postNotification(sodaTest.SlackChannel, topSection, attachments); err != nil {
			return err
		}
	} else if sodaTest.SlackNotifyOnPassedScan != nil {
		notifyPassed, err := strconv.ParseBool(*sodaTest.SlackNotifyOnPassedScan)
		if err != nil {
			s.log.Warningf("unable to parse provided boolean value '%v' for slack notification on passed soda scan: %v", sodaTest.SlackNotifyOnPassedScan, err)
			return nil
		}
		if notifyPassed {
			topSection, attachments := s.createPassedScanMessage(sodaTest.GCPProject, sodaTest.Dataset)
			if err := s.postNotification(sodaTest.SlackChannel, topSection, attachments); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Client) findDiscrepancies(configError *string, sodaResults []models.TestResult) (bool, testDiscrepancies) {
	discrepancies := testDiscrepancies{}
	for _, r := range sodaResults {
		switch r.Outcome {
		case "pass":
			continue
		case "fail", "error":
			discrepancies.Fails = append(discrepancies.Fails, r)
		default:
			discrepancies.Warnings = append(discrepancies.Warnings, r)
		}
	}

	return configError != nil || len(discrepancies.Fails) > 0 || len(discrepancies.Warnings) > 0, discrepancies
}

func (s *Client) postNotification(slackChannel string, topSection slack.Block, attachments []slack.Attachment) error {
	slackClient := slack.New(s.token)

	_, _, err := slackClient.PostMessage(slackChannel, slack.MsgOptionBlocks(topSection), slack.MsgOptionAttachments(attachments...))
	if err != nil {
		return err
	}

	return nil
}

func (s *Client) createPassedScanMessage(projectID, dataset string) (slack.Block, []slack.Attachment) {
	emoji := true
	topMessage := slack.TextBlockObject{
		Type:  "plain_text",
		Text:  "SODA scan gjennomført uten feil :checked:",
		Emoji: &emoji,
	}

	topSection := slack.NewSectionBlock(&topMessage, nil, nil)
	attachments := []slack.Attachment{
		{
			Color:      "#00ff00",
			AuthorName: fmt.Sprintf("%v.%v", projectID, dataset),
			Footer:     "SODA Bot",
		},
	}

	return topSection, attachments
}

func (s *Client) createDiscrepancyMessage(configError *string, d testDiscrepancies, projectID, dataset string) (slack.Block, []slack.Attachment) {
	emoji := true
	topMessage := slack.TextBlockObject{
		Type:  "plain_text",
		Text:  "Datakvalitetssjekk feiler :gasp:",
		Emoji: &emoji,
	}

	topSection := slack.NewSectionBlock(&topMessage, nil, nil)
	attachments := []slack.Attachment{}

	if configError != nil {
		attachments = append(attachments, slack.Attachment{
			Color:      "#ff2d00",
			AuthorName: fmt.Sprintf("%v.%v", projectID, dataset),
			Title:      "Tester har konfigurasjonsfeil",
			Text:       *configError,
			Footer:     "SODA Bot",
			Fallback:   "Tester har konfigurasjonsfeil",
		})
	}

	if len(d.Fails) > 0 {
		message := ""
		for _, e := range d.Fails {
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
		for _, w := range d.Warnings {
			line1 := ""
			if w.Column != "" {
				line1 = fmt.Sprintf("_*Tabell: %v*_ _*kolonne: %v*_\n", w.Table, w.Column)
			} else {
				line1 = fmt.Sprintf("_*Tabell: %v*_\n", w.Table)
			}
			line2 := w.Test + "\n"
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
