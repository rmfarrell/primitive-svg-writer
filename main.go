package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	// fmt.Printf("Body size = %d.\n", len(request.Body))
	// fmt.Println("Headers:")
	// for key, value := range request.Headers {
	// 	fmt.Printf("    %s: %s\n", key, value)
	// }
	settings := ""
	for key, val := range request.QueryStringParameters {

	}

	return events.APIGatewayProxyResponse{Body: "TODO: svg text", StatusCode: 200}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handleRequest)
}
