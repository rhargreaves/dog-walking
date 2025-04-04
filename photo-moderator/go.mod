module github.com/rhargreaves/dog-walking/photo-moderator

go 1.23.1

require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/aws/aws-sdk-go v1.55.6
	github.com/rhargreaves/dog-walking/shared v0.0.0
	github.com/stretchr/testify v1.10.0
)

replace github.com/rhargreaves/dog-walking/shared => ../shared

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
