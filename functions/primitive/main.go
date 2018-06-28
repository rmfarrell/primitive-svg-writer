package main

import (
	"context"
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmfarrell/primitive-constructs/primitive"
)

var (
	lambdaPath string
)

func init() {
	lambdaPath = fmt.Sprintf("%s", os.Getenv("LAMBDA_TASK_ROOT"))
	os.Setenv("PATH", lambdaPath)
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// for key, val := range request.QueryStringParameters {

	// }
	staticImgPath := fmt.Sprintf("%s/export/honkysauce.png", lambdaPath)

	// w, h, err := getImageDimension(staticImgPath)
	// dims := fmt.Sprintf("%d:::%d", w, h)

	svg, err := primitive.WriteSvg(staticImgPath, 20)

	return events.APIGatewayProxyResponse{Body: svg, StatusCode: 200}, err
}

func main() {
	// svg, _ := primitive.WriteSvg("./honkysauce.png", 20)
	// fmt.Println(svg)
	lambda.Start(handleRequest)
}

func getImageDimension(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	defer file.Close()
	if err != nil {
		return 0, 0, fmt.Errorf("%v\n", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height, nil
}
