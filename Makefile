rm-generated:
	mkdir -p proto-gen 
	rm -rf proto-gen/*

generated: rm-generated
	protoc \
		--go_out=./proto-gen \
		--go-grpc_out=./proto-gen \
		 proto/*.proto

proto: generated
	protoc-go-inject-tag -input="proto-gen/*/*.pb.go" -remove_tag_comment

.SILENT: