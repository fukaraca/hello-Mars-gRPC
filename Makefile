generate:
	protoc --proto_path=pb/ pb/hm.proto --go_opt=paths=source_relative --go_out=./pb --go-grpc_opt=paths=source_relative --go-grpc_out=./pb
run:
	go run ./server/main.go
	go run ./client/main.go