package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	message = flag.String("message", "", "The message to send")
	isError = flag.Bool("error", false, "Message is an error")
	isAlert = flag.Bool("alert", false, "Alert @here when message is sent")

	field  = flag.String("field", "", "Slack attachment field")
	fields = flag.String("fields", "", "Slack attachment fields")

	// optional
	webhook = flag.String("webhook", "", "Slack Webhook")
)

// Webhook is the Slack webhook to use
var (
	Webhook = os.Getenv("SLACK_WEBHOOK")
)

func main() {
	flag.Parse()
	log.SetPrefix("[Slack] ")

	if Webhook == "" {
		Webhook = *webhook
	}
	if Webhook == "" {
		log.Fatal("missing SLACK_WEBHOOK in environment, please see slack README")
	}
	if *field != "" && *fields != "" {
		log.Fatal("both 'field' and 'fields' can't be specified, use 'fields' for multiple fields")
	}
	if *message == "" {
		log.Fatal("missing --message, must be specified")
	}
	if *field != "" || *fields != "" {
		if err := loadFields(); err != nil {
			log.Fatalf("unable to load fields: %v", err)
		}
	}

	if err := slack(*message); err != nil {
		log.Fatalf("failed to send message: %v", err)
	}
}

// Field defines the structure of a Slack attachment field
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type attachment struct {
	Color  string  `json:"color"`
	Fields []Field `json:"fields"`
}

// Alert defines a slack notification/alert
type Alert struct {
	Text        string       `json:"text"`
	Err         error        `json:"-"`
	Attachments []attachment `json:"attachments"`
}

var theFields []Field

func loadFields() error {
	// TODO: implement fields parsing
	return errors.New("TODO: --field(s) support pending")
}

// trim trims a message to a specific length
func trim(message string, length int) string {
	results := strings.Split(string(message), "\n")
	if len(results) <= length {
		return message
	}
	trimmed := results[len(results)-(length-1):]
	results = append([]string{"... [Results TRIMMED for display] ..."}, trimmed...)

	return strings.Join(results, "\n")
}

func slack(message string) error {
	color := "good"
	if *isError {
		color = "danger"
		theFields = append(theFields, Field{
			Title: "Error",
			Value: "```" + trim(message, 10) + "```",
			Short: false,
		})
	}

	if *isAlert {
		message += " _cc_ <!here>"
	}

	sa := Alert{
		Text: message,
		Attachments: []attachment{
			{
				Color:  color,
				Fields: theFields,
			},
		},
	}

	jsonBody, err := json.Marshal(sa)
	if err != nil {
		return err
	}

	cl := &http.Client{Timeout: time.Second * 20}
	req, err := http.NewRequest("POST", Webhook, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("unable to use specified webhook %v", err)
	}
	resp, err := cl.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to slack: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read slack response: %v", err)
	}
	return fmt.Errorf("slack responded with invalid response [%s] %s", resp.Status, dat)
}
