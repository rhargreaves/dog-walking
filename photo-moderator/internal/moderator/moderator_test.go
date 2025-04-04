package moderator

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rhargreaves/dog-walking/photo-moderator/internal/domain"
	aws_mocks "github.com/rhargreaves/dog-walking/photo-moderator/internal/moderator/aws/mocks"
	"github.com/rhargreaves/dog-walking/photo-moderator/internal/moderator/breed_detector"
	breed_detector_mocks "github.com/rhargreaves/dog-walking/photo-moderator/internal/moderator/breed_detector/mocks"
	"github.com/stretchr/testify/require"
)

func TestModeratePhoto_ApprovesLabradorDog(t *testing.T) {
	dogTableName := "dog-table"
	approvedPhotosBucket := "approved-photos-bucket"
	pendingPhotosBucket := "pending-photos-bucket"
	dogId := "1"
	hash := "123"

	breedDetector := &breed_detector_mocks.MockBreedDetector{}
	breedDetector.DetectBreedFunc = func(id string) (*domain.BreedDetectionResult, error) {
		return &domain.BreedDetectionResult{Breed: "Labrador", Confidence: 0.95}, nil
	}

	var photoStatus string
	dynamodbClient := &aws_mocks.MockDynamoDB{}
	dynamodbClient.UpdateItemFunc = func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
		photoStatus = *input.ExpressionAttributeValues[":photoStatus"].S
		return &dynamodb.UpdateItemOutput{}, nil
	}

	s3Client := &aws_mocks.MockS3{}
	s3Client.PutObjectFunc = func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
		return &s3.PutObjectOutput{}, nil
	}
	s3Client.CopyObjectFunc = func(input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error) {
		result := &s3.CopyObjectResult{ETag: aws.String(hash)}
		return &s3.CopyObjectOutput{CopyObjectResult: result}, nil
	}
	s3Client.DeleteObjectFunc = func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
		return &s3.DeleteObjectOutput{}, nil
	}

	moderator := NewModerator(dogTableName, approvedPhotosBucket, breedDetector, dynamodbClient, s3Client)
	err := moderator.ModeratePhoto(pendingPhotosBucket, dogId)
	require.NoError(t, err)
	require.Equal(t, "approved", photoStatus)
}

func TestModeratePhoto_ApprovesDogWhenBreedIsNonSpecific(t *testing.T) {
	dogTableName := "dog-table"
	approvedPhotosBucket := "approved-photos-bucket"
	pendingPhotosBucket := "pending-photos-bucket"
	dogId := "1"
	hash := "123"

	breedDetector := &breed_detector_mocks.MockBreedDetector{}
	breedDetector.DetectBreedFunc = func(id string) (*domain.BreedDetectionResult, error) {
		return nil, breed_detector.ErrNoSpecificBreedDetected
	}

	var photoStatus string
	var breed string = "existing-value"
	dynamodbClient := &aws_mocks.MockDynamoDB{}
	dynamodbClient.UpdateItemFunc = func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
		photoStatus = *input.ExpressionAttributeValues[":photoStatus"].S
		if breedAttrValue := input.ExpressionAttributeValues[":breed"]; breedAttrValue != nil {
			breed = *breedAttrValue.S
		}
		return &dynamodb.UpdateItemOutput{}, nil
	}

	s3Client := &aws_mocks.MockS3{}
	s3Client.PutObjectFunc = func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
		return &s3.PutObjectOutput{}, nil
	}
	s3Client.CopyObjectFunc = func(input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error) {
		result := &s3.CopyObjectResult{ETag: aws.String(hash)}
		return &s3.CopyObjectOutput{CopyObjectResult: result}, nil
	}
	s3Client.DeleteObjectFunc = func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
		return &s3.DeleteObjectOutput{}, nil
	}

	moderator := NewModerator(dogTableName, approvedPhotosBucket, breedDetector, dynamodbClient, s3Client)
	err := moderator.ModeratePhoto(pendingPhotosBucket, dogId)
	require.NoError(t, err)
	require.Equal(t, "approved", photoStatus)
	require.Equal(t, breed, "existing-value")
}
