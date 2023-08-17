gen:
	go generate ./...

fmt: gen
	go fmt ./...

test : gen
	go test ./...

cover : gen
	go test ./... -coverprofile=/tmp/coverage.out && go tool cover -html=/tmp/coverage.out

tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

build: test
