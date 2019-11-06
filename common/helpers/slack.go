package helpers

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// PushToSlack -- sends msg to slack channel
func PushToSlack(msg string, msgType string, stack string) {

	type Payload struct {
		Text string `json:"text"`
	}

	data := Payload{
		Text: msg,
	}

	if msgType == "panic" {
		data.Text = "panic --> " + msg
	} else if msgType == "err" {
		data.Text = "err --> " + msg
	}

	data.Text += " \n\n" + stack

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://hooks.slack.com/services/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", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
}
