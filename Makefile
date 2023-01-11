gen_user_proto:
	protoc --proto_path=proto proto/user.proto --go_out=user --go-grpc_out=user
	protoc --proto_path=proto proto/user.proto --go_out=gateway --go-grpc_out=gateway

gen_mfa_proto:
	protoc --proto_path=proto proto/mfa.proto --go_out=mfa --go-grpc_out=mfa
	protoc --proto_path=proto proto/mfa.proto --go_out=gateway --go-grpc_out=gateway