package aws_session

import "os"

func IsLocal() bool {
	return os.Getenv("USE_LOCALSTACK") == "true"
}
