.PHONY: build clean

build:
	GOOS=linux GOARCH=amd64 go build -o bin/linux-net-bytes-exporter main.go

clean:
	rm -rf ./bin
