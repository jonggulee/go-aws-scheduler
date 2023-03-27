package utils

import (
	"fmt"

	"github.com/slack-go/slack"
)

const (
	// WebhookUrl   = "https://hooks.slack.com/services/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"

	// test
	WebhookUrl   = "https://hooks.slack.com/services/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"
	SuccessColor = "#18be52"
	FailedColor  = "#E96D76"
)

func SlackNoti(status interface{}, id, msg string, isErr bool) error {
	SuccessMsg := ":white_check_mark: Success"
	FailedMsg := ":x: Failed"
	defaultResultTitle := "Result"

	color := ""
	attachment := slack.Attachment{
		Fields: []slack.AttachmentField{
			{Title: defaultResultTitle, Value: "", Short: false},
			{Title: "Target", Value: "", Short: false},
		},
	}

	value := fmt.Sprintf("%s\n", id)

	if isErr {
		color = FailedColor
		attachment.Fields = append(attachment.Fields, slack.AttachmentField{
			Title: "Message", Value: fmt.Sprintf("```%s```", msg), Short: false,
		})
		attachment.Fields[0].Value = FailedMsg
	}

	if !isErr {
		color = SuccessColor
		attachment.Fields = append(attachment.Fields, slack.AttachmentField{
			Title: "Message", Value: fmt.Sprintf("```%s```", msg), Short: false,
		})
		attachment.Fields[0].Value = SuccessMsg
	}

	attachment.Fields[1].Value = fmt.Sprintf("```%s```", value)
	attachment.Color = color

	message := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook(WebhookUrl, &message)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
