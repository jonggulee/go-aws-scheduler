package utils

import (
	"fmt"

	"github.com/slack-go/slack"
)

const (
	WebhookUrl   = "https://hooks.slack.com/services/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"
	SuccessColor = "#18be52"
	FailedColor  = "#E96D76"
	SuccessMsg   = ":white_check_mark: Success"
	FailedMsg    = ":x: Failed"
)

var message []string
var hasError bool

func InputSlackData(msg string, isErr bool) {
	message = append(message, msg)
	hasError = isErr || hasError
}

func getSlackMessages() string {
	var messages string

	for _, mergeMsg := range message {
		messages += fmt.Sprint(mergeMsg)
	}
	return messages
}

func SendSlackMessage() {
	color := ""
	attachment := slack.Attachment{
		Fields: []slack.AttachmentField{
			{Title: "Result", Value: "", Short: false},
		},
	}

	msg := getSlackMessages()

	if hasError {
		color = FailedColor
		attachment.Fields = append(attachment.Fields, slack.AttachmentField{
			Title: "Message", Value: fmt.Sprintf("```%s```", msg), Short: false,
		})
		attachment.Fields[0].Value = FailedMsg
	}

	if !hasError {
		color = SuccessColor
		attachment.Fields = append(attachment.Fields, slack.AttachmentField{
			Title: "Message", Value: fmt.Sprintf("```%s```", msg), Short: false,
		})
		attachment.Fields[0].Value = SuccessMsg
	}

	attachment.Color = color
	message := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook(WebhookUrl, &message)
	HandleErr(err)
}
