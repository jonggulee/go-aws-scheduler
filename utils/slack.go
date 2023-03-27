package utils

import (
	"fmt"

	"github.com/slack-go/slack"
)

const (
	webhookUrl   = "https://hooks.slack.com/services/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"
	successColor = "#18be52"
)

func SlackNoti(status interface{}, id, msg string) error {
	attachment := slack.Attachment{
		Fields: []slack.AttachmentField{
			{
				Title: "Instance start result",
				Value: ":white_check_mark: Success",
				Short: false,
			},
			{
				Title: "Target Instance",
				Value: "",
				Short: false,
			},
		},
	}

	// msg := slack.WebhookMessage{
	// 	Attachments: []slack.Attachment{attachment},
	// }

	err := slack.PostWebhook(webhookUrl, &msg)
	if err != nil {
		fmt.Println(err)
		return err
	}
}
