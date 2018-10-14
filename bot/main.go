package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	AnnoyanceTriggers = []string{"bot", "Bot", "BOT",}
	BotId             = "c1d5a8381af0c8b21fe3bcb2b6"
	// TODO Get value at runtime
	Token   = "<substitute real token at deploy time>"
	ChatURL = (&url.URL{
		Scheme:   "https",
		Host:     "api.groupme.com",
		Path:     "v3/bots/post",
		RawQuery: fmt.Sprintf("token=%s", Token),
	}).String()
	EmptyResponse = events.APIGatewayProxyResponse{}
	RandFunc      = rand.Float64
	IlyaRate      = 0.02
)

type GroupMeMessage struct {
	BotId string `json:"bot_id"`
	Text  string `json:"text"`
}

type GroupMeCallback struct {
	Name       string `json:"name"`
	SenderType string `json:"sender_type"`
	Text       string `json:"text"`
}

// TODO Reabstract and test in better isolation
func makePost(text string) (events.APIGatewayProxyResponse, error) {
	botMsg, err := json.Marshal(GroupMeMessage{
		BotId: BotId,
		Text:  text,
	})
	if err != nil {
		return EmptyResponse, err
	}

	resp, err := http.Post(ChatURL, "application/json", bytes.NewBuffer(botMsg))
	if err != nil {
		return EmptyResponse, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: resp.StatusCode,
	}, nil
}

func samHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println(request.Body)

	var inMsg *GroupMeCallback
	err := json.Unmarshal([]byte(request.Body), &inMsg)
	if err != nil {
		return EmptyResponse, err
	}

	if inMsg.SenderType != "user" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
		}, nil
	}

	if RandFunc() < IlyaRate {
		return makePost("When I talk English, I’m still computing.")
	} else {
		for _, t := range AnnoyanceTriggers {
			if strings.Contains(inMsg.Text, t) {
				return makePost("хватит надоедать мне!")
			}
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(samHandler)
}
