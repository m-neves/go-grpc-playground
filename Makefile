all: gen-protoc

gen-protoc:
	protoc --go_out=. --go-grpc_out=. api/proto/*.proto

client-start:
	clear
	go run client/client.go

server-start:
	clear
	go run server/server.go
