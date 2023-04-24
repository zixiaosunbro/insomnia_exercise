
proto:
	mkdir -p proto/rule
	protoc --go_out=plugins=grpc,paths=source_relative:. proto/rule.proto
	mv proto/rule.*.go proto/rule

.PHONY: proto