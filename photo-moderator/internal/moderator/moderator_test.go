package moderator

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	aws_mocks "github.com/rhargreaves/dog-walking/photo-moderator/internal/moderator/aws/mocks"
	"github.com/rhargreaves/dog-walking/photo-moderator/internal/moderator/breed_detector"
	breed_detector_mocks "github.com/rhargreaves/dog-walking/photo-moderator/internal/moderator/breed_detector/mocks"
	"github.com/rhargreaves/dog-walking/photo-moderator/internal/moderator/content_screener"
	content_screener_mocks "github.com/rhargreaves/dog-walking/photo-moderator/internal/moderator/content_screener/mocks"
	"github.com/stretchr/testify/require"
)

const (
	dogTableName         = "dog-table"
	approvedPhotosBucket = "approved-photos-bucket"
	pendingPhotosBucket  = "pending-photos-bucket"
	dogId                = "1"
	hash                 = "123"
)

func TestModeratePhoto_ApprovesLabradorDog(t *testing.T) {
	contentScreener := mockContentScreenerReturningSafe()
	breedDetector := mockBreedDetectorReturningLabrador()
	var dbPhotoStatus string
	dynamodbClient := mockDynamoDBClientUpdatingPhotoRecords(&dbPhotoStatus)
	s3Client := mockS3ClientReturningHash(hash)

	moderator := NewModerator(dogTableName, approvedPhotosBucket, breedDetector, dynamodbClient, s3Client, contentScreener)
	photoStatus, err := moderator.ModeratePhoto(pendingPhotosBucket, dogId)
	require.NoError(t, err)
	require.Equal(t, PhotoStatusApproved, photoStatus)
	require.Equal(t, photoStatus, dbPhotoStatus)
}

func TestModeratePhoto_ApprovesDogWhenBreedIsNonSpecific(t *testing.T) {
	contentScreener := mockContentScreenerReturningSafe()
	breedDetector := mockBreedDetectorReturningNoSpecificBreed()

	var photoStatus string
	var breed string = "existing-value"
	dynamodbClient := mockDynamoDBClientUpdatingPhotoRecordsWithBreed(&photoStatus, &breed)
	s3Client := mockS3ClientReturningHash(hash)

	moderator := NewModerator(dogTableName, approvedPhotosBucket, breedDetector, dynamodbClient, s3Client, contentScreener)
	photoStatus, err := moderator.ModeratePhoto(pendingPhotosBucket, dogId)
	require.NoError(t, err)
	require.Equal(t, PhotoStatusApproved, photoStatus)
	require.Equal(t, breed, "existing-value")
}

func TestModeratePhoto_RejectsPhotoWhenContentScreenerReturnsUnsafe(t *testing.T) {
	contentScreener := mockContentScreenerReturningUnsafeViolence()
	var dbPhotoStatus string
	dynamodbClient := mockDynamoDBClientUpdatingPhotoRecords(&dbPhotoStatus)
	s3Client := mockS3ClientReturningHash(hash)

	moderator := NewModerator(dogTableName, approvedPhotosBucket, nil, dynamodbClient, s3Client, contentScreener)
	photoStatus, err := moderator.ModeratePhoto(pendingPhotosBucket, dogId)
	require.NoError(t, err)
	require.Equal(t, PhotoStatusRejected, photoStatus)

}

func mockContentScreenerReturningUnsafeViolence() *content_screener_mocks.MockContentScreener {
	return &content_screener_mocks.MockContentScreener{
		ScreenImageFunc: func(id string) (*content_screener.ContentScreenerResult, error) {
			return &content_screener.ContentScreenerResult{IsSafe: false, Reason: "violence"}, nil
		},
	}
}

func mockContentScreenerReturningSafe() *content_screener_mocks.MockContentScreener {
	return &content_screener_mocks.MockContentScreener{
		ScreenImageFunc: func(id string) (*content_screener.ContentScreenerResult, error) {
			return &content_screener.ContentScreenerResult{IsSafe: true}, nil
		},
	}
}

func mockBreedDetectorReturningLabrador() *breed_detector_mocks.MockBreedDetector {
	return &breed_detector_mocks.MockBreedDetector{
		DetectBreedFunc: func(id string) (*breed_detector.BreedDetectionResult, error) {
			return &breed_detector.BreedDetectionResult{Breed: "Labrador", Confidence: 0.95}, nil
		},
	}
}

func mockBreedDetectorReturningNoSpecificBreed() *breed_detector_mocks.MockBreedDetector {
	return &breed_detector_mocks.MockBreedDetector{
		DetectBreedFunc: func(id string) (*breed_detector.BreedDetectionResult, error) {
			return nil, breed_detector.ErrNoSpecificBreedDetected
		},
	}
}

func mockDynamoDBClientUpdatingPhotoRecords(photoStatus *string) *aws_mocks.MockDynamoDB {
	return &aws_mocks.MockDynamoDB{
		UpdateItemFunc: func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
			*photoStatus = *input.ExpressionAttributeValues[":photoStatus"].S
			return &dynamodb.UpdateItemOutput{}, nil
		},
	}
}

func mockDynamoDBClientUpdatingPhotoRecordsWithBreed(photoStatus *string, breed *string) *aws_mocks.MockDynamoDB {
	dynamodbClient := &aws_mocks.MockDynamoDB{}
	dynamodbClient.UpdateItemFunc = func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
		*photoStatus = *input.ExpressionAttributeValues[":photoStatus"].S
		if breedAttrValue := input.ExpressionAttributeValues[":breed"]; breedAttrValue != nil {
			*breed = *breedAttrValue.S
		}
		return &dynamodb.UpdateItemOutput{}, nil
	}
	return dynamodbClient
}

func mockS3ClientReturningHash(hash string) *aws_mocks.MockS3 {
	return &aws_mocks.MockS3{
		CopyObjectFunc: func(input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error) {
			result := &s3.CopyObjectResult{ETag: aws.String(hash)}
			return &s3.CopyObjectOutput{CopyObjectResult: result}, nil
		},
		DeleteObjectFunc: func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
			return &s3.DeleteObjectOutput{}, nil
		},
		PutObjectFunc: func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
			return &s3.PutObjectOutput{}, nil
		},
	}
}
