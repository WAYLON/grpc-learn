protoc -I. \
-I third_party \
--go_out=proto --go_opt=paths=source_relative \
--go-grpc_out=proto --go-grpc_opt=paths=source_relative \
--openapiv2_out=proto --openapiv2_opt logtostderr=true \
--validate_out=lang=go:proto --validate_opt=paths=source_relative \
--grpc-gateway_out=proto --grpc-gateway_opt=paths=source_relative \
pb/helloworld.proto


protoc -I. \
-I third_party \
--go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
--openapiv2_out=. --openapiv2_opt logtostderr=true \
--validate_out=lang=go:. --validate_opt=paths=source_relative \
--grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
pb/helloworld.proto
