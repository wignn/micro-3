gen:
	@protoc \
		--proto_path=proto "proto/order.proto" \
		--go_out=genproto --go_opt=paths=source_relative \
  		--go-grpc_out=genproto --go-grpc_opt=paths=source_relative