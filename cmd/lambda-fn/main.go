package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/umayr/hook"
)

var (
	// APIKey for Pushover application
	APIKey string
	// UserKey for Pushover application
	UserKey string

	// ErrUnknownAPI is thrown when github invokes the lambda function for the case that's
	// not a valid event
	ErrUnknownAPI = fmt.Errorf("unknown api request")
)

func init() {
	if os.Getenv("API_KEY") == "" || os.Getenv("USER_KEY") == "" {
		panic("keys are not set for API and recipient")
	}

	APIKey, UserKey = os.Getenv("API_KEY"), os.Getenv("USER_KEY")
}

// Handler for lambda function
// It verifies the request, parses request body and invokes the hook implementation
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	event, exists := request.Headers["X-GitHub-Event"]
	if !exists || !hook.Events.Has(event) {
		return events.APIGatewayProxyResponse{Body: ErrUnknownAPI.Error(), StatusCode: 400}, ErrUnknownAPI
	}

	p := new(hook.Payload)
	if err := json.Unmarshal([]byte(request.Body), p); err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	h := hook.NewHook(p, hook.NewPushover(APIKey, UserKey))
	if err := h.Perform(); err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
