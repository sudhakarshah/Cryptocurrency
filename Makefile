
all: mp2.go
	go build mp2.go

testall: parser_test.go node_test.go
	go test mp2_parser.go parser_test.go
testnode: node_test.go
	go test mp2_node.go node_test.go
