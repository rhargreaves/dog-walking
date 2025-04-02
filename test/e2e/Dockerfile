FROM golang:1.23.4-alpine
RUN apk add --no-cache docker docker-compose
WORKDIR /e2e-test
COPY . .
RUN go mod download && \
	go install gotest.tools/gotestsum@latest
CMD ["gotestsum", "--format", "testname", "./..."]
