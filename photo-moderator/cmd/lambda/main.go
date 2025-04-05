package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, req events.S3Event) error {
	for _, record := range req.Records {
		sourceBucket := record.S3.Bucket.Name
		sourceKey := record.S3.Object.Key
		fmt.Printf("Source bucket: %s, source key: %s\n", sourceBucket, sourceKey)

		moderator := createModerator(sourceBucket)
		photoStatus, err := moderator.ModeratePhoto(sourceBucket, sourceKey)
		if err != nil {
			return fmt.Errorf("failed to moderate photo %s: %w", sourceKey, err)
		}
		fmt.Printf("Photo status: %s\n", photoStatus)
	}
	return nil
}
