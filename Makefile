test:
	cd src ; go test -v
build: test
	go build -o bin/dyndns src/dyndns.go
