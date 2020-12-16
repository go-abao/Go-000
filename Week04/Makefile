build:
	go build -o bin/server cmd/server

run:
	./bin/server

rpc:
	go run test/user.go

gen:
	protoc -I. --go_out=plugins=grpc:. api/user/v1/user.proto