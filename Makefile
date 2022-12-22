gen_proto:
	protoc --proto_path=proto proto/*.proto --go_out=user --go-grpc_out=user
	protoc --proto_path=proto proto/*.proto --go_out=gateway --go-grpc_out=gateway