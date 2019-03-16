
all: mp2.go
	go build mp2.go mp2_node.go mp2_parser.go

cleanlogs:
	rm node*.log
clean:
	rm node*.log mp2
