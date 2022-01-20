package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/vikinatora/rarity-cloud-function/handlers"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	response, err := handlers.GetPolymorphs(ctx, request)
	
	return response, err
}

func main() {
	lambda.Start(handler)
}
