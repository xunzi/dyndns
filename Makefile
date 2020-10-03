test:
	cd client ; go test -v
build: test
	go build -o bin/dyndns client/dyndns.go
