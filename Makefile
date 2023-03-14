proto:
	protoc --proto_path=. "user.proto" --go_out="." --go_opt=paths=source_relative --go-grpc_opt=require_unimplemented_servers=false  "--go-grpc_out=." --go-grpc_opt=paths=source_relative
	protoc *.proto \
        --go-grpc_out=. \
        --go-grpc_opt=paths=source_relative \
        --proto_path=.
