package dogs

import (
	"encoding/json"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
)

func s3Event(t *testing.T, bucketName string, objectKey string) []byte {
	now := time.Now()
	event := events.S3Event{
		Records: []events.S3EventRecord{
			{
				EventVersion: "2.0",
				EventSource:  "aws:s3",
				AWSRegion:    "eu-west-1",
				EventTime:    now,
				EventName:    "ObjectCreated:Put",
				S3: events.S3Entity{
					SchemaVersion:   "1.0",
					ConfigurationID: "local-emulator",
					Bucket: events.S3Bucket{
						Name: bucketName,
						OwnerIdentity: events.S3UserIdentity{
							PrincipalID: "local-emulator",
						},
					},
					Object: events.S3Object{
						Key:       objectKey,
						Size:      0,
						ETag:      "dummy-etag",
						VersionID: "1",
						Sequencer: "dummy-sequencer",
					},
				},
			},
		},
	}
	eventJsonBytes, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Error marshaling event: %v", err)
	}

	return eventJsonBytes
}

func invokePhotoModerator(t *testing.T, dogID string) {
	event := s3Event(t, os.Getenv("PENDING_DOG_IMAGES_BUCKET"), dogID)
	dockerCompose := exec.Command("docker", "compose",
		"-f", "/proj/docker-compose.yml",
		"exec", "sam", "sam", "local", "invoke",
		"--container-host-interface", "0.0.0.0",
		"--container-host", os.Getenv("CONTAINER_HOST"),
		"--docker-volume-basedir", os.Getenv("PROJECT_ROOT"),
		"--docker-network", "dog-walking_default",
		"--skip-pull-image",
		"--event", "-",
		"PhotoModeratorFunction")
	stdin, err := dockerCompose.StdinPipe()
	require.NoError(t, err)
	go func() {
		defer stdin.Close()
		stdin.Write(event)
	}()

	output, err := dockerCompose.CombinedOutput()
	if err != nil {
		t.Logf("Photo moderator function output: %s", string(output))
	}
	require.NoError(t, err)
}
