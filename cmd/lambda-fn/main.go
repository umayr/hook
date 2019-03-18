package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/umayr/hook"
)

var (
	ErrUnknownAPI = fmt.Errorf("unknown api request")
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	event, exists := request.Headers["X-GitHub-Event"]
	if !exists || !hook.Events.Has(event) {
		return events.APIGatewayProxyResponse{Body: ErrUnknownAPI.Error(), StatusCode: 400}, ErrUnknownAPI
	}

	p := new(hook.Payload)
	if err := json.Unmarshal([]byte(request.Body), p); err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	h := hook.NewHook(p)
	if err := h.Perform(); err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
