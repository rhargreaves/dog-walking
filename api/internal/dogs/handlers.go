package dogs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/rhargreaves/dog-walking/api/internal/common"
)

type DogHandler interface {
	ListDogs(c *gin.Context)
	GetDog(c *gin.Context)
	CreateDog(c *gin.Context)
	UpdateDog(c *gin.Context)
	DeleteDog(c *gin.Context)
	UploadDogPhoto(c *gin.Context)
}

type dogHandler struct {
	dogRepository DogRepository
}

func NewDogHandler(dogRepository DogRepository) DogHandler {
	return &dogHandler{dogRepository: dogRepository}
}

func (h *dogHandler) CreateDog(c *gin.Context) {
	var dog Dog
	if err := c.ShouldBindJSON(&dog); err != nil {
		handleBindError(c, err)
		return
	}

	if err := h.dogRepository.Create(&dog); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dog)
}

func (h *dogHandler) ListDogs(c *gin.Context) {
	dogs, err := h.dogRepository.List()
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dogs)
}

func (h *dogHandler) GetDog(c *gin.Context) {
	id := c.Param("id")
	dog, err := h.dogRepository.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dog)
}

func (h *dogHandler) UpdateDog(c *gin.Context) {
	id := c.Param("id")
	var dog Dog
	if err := c.ShouldBindJSON(&dog); err != nil {
		handleBindError(c, err)
		return
	}

	if err := h.dogRepository.Update(id, &dog); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dog)
}

func (h *dogHandler) DeleteDog(c *gin.Context) {
	id := c.Param("id")
	if err := h.dogRepository.Delete(id); err != nil {
		handleError(c, err)
		return
	}
}

func createS3Session() (*session.Session, error) {
	useLocalStack := os.Getenv("USE_LOCALSTACK") == "true"
	region := os.Getenv("AWS_REGION")
	if useLocalStack {
		return session.NewSession(&aws.Config{
			Region:      &region,
			Endpoint:    aws.String(os.Getenv("AWS_S3_ENDPOINT_URL")),
			Credentials: credentials.NewStaticCredentials("test", "test", ""),
		})
	}
	return session.NewSession(&aws.Config{
		Region: &region,
	})
}

func (h *dogHandler) UploadDogPhoto(c *gin.Context) {
	id := c.Param("id")
	_, err := h.dogRepository.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}

	bucketName := os.Getenv("DOG_IMAGES_BUCKET")
	if bucketName == "" {
		handleError(c, errors.New("DOG_IMAGES_BUCKET environment variable not set"))
		return
	}

	sess, err := createS3Session()
	if err != nil {
		handleError(c, err)
		return
	}

	s3Client := s3.New(sess)

	fileBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		handleError(c, err)
		return
	}

	key := fmt.Sprintf("%s", id)
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(c.GetHeader("Content-Type")),
	})

	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func handleError(c *gin.Context, err error) {
	if errors.Is(err, ErrDogNotFound) {
		c.Error(common.APIError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}
	c.Error(err)
}

func handleBindError(c *gin.Context, err error) {
	c.Error(common.APIError{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	})
}
