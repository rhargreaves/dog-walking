FROM golang:1.23.4-alpine
WORKDIR /app
COPY . .
RUN go mod download && \
	go install gotest.tools/gotestsum@latest
CMD ["gotestsum", "--format", "testname", "./..."]
