package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
)

var (
	EmptyResponse = events.APIGatewayProxyResponse{}
)

type SAMHandler struct {
	gm         GroupMe
	responders []Responder
}

func NewSAMHandler(
	gm GroupMe,
	responders ...Responder) *SAMHandler {

	return &SAMHandler{
		gm:         gm,
		responders: responders,
	}
}

func (s *SAMHandler) Handle(
	request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Println(request.Body)

	inMsg, err := GroupMeParse(request.Body)
	if err != nil {
		return EmptyResponse, err
	}

	if inMsg.SenderType != "user" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
		}, nil
	}

	for _, r := range s.responders {
		if resp := r.Respond(inMsg.Text); resp != nil {
			if code, err := s.gm.Post(*resp); err != nil {
				return EmptyResponse, err
			} else {
				return events.APIGatewayProxyResponse{
					StatusCode: code,
				}, nil
			}
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
