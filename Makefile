all:
	go build consistenthash.go
test:
	go test
bench:
	go test -bench=".*"
clean:
	rm -f consistenthash
