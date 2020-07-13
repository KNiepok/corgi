package api

import "fmt"

var (
	errorIntro = block{
		BlockType: "section",
		Text: text{
			TextType: "mrkdwn",
			Text:     ":fearful: I'm not able to understand your request!",
		},
	}
	example = block{
		BlockType: "section",
		Text: text{
			TextType: "mrkdwn",
			Text:     "You should one of those formats:```daily @ HH:mm (ie. daily @ 17:10)``` ```daily @ HH (ie. daily @ 10)``` ```weekly @ DAY HH:mm (ie. weekly @ MON 12:00)``` ```weekly @ DAY HH (ie. weekly @ SUN 9)```",
		},
	}

	subscribed = block{
		BlockType: "section",
		Text: text{
			TextType: "mrkdwn",
			Text:     "All good, you've been subscribed! To change subscription just use call */subscribe* again. If you want to remove subscription, just call */unsubscribe!*",
		},
	}
	unsubscribed = block{
		BlockType: "section",
		Text: text{
			TextType: "mrkdwn",
			Text:     "Sorry to hear that! You'be been unsubscribed. Use */subscribe* to enroll again. See you around!",
		},
	}
	unexpectedError = block{
		BlockType: "section",
		Text: text{
			TextType: "mrkdwn",
			Text:     "Woops, something went wrong. Try again later or contact a maintainer!",
		},
	}
)

func subscriptionErrorExplanation(input string) block {
	msg := fmt.Sprintf("You tried to subscribe with: ```%s ```", input)
	if len(input) == 0 {
		msg = "You tried to subscribe but forgot to tell when to receive notifications."
	}
	return block{
		BlockType: "section",
		Text: text{
			TextType: "mrkdwn",
			Text:     msg,
		},
	}
}

type slackResponse struct {
	Blocks []block `json:"blocks"`
}

type block struct {
	BlockType string `json:"type"`
	Text      text   `json:"text"`
}

type text struct {
	TextType string `json:"type"`
	Text     string `json:"text"`
}
