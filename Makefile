DIR=./bin
BINARY=nba-now

make: build

build: test
	GOARCH=amd64 GOOS=linux go build -o ${DIR}/${BINARY}-linux ./cmd/main.go
	GOARCH=amd64 GOOS=windows go build -o ${DIR}/${BINARY}-win ./cmd/main.go
	GOARCH=arm64 GOOS=darwin go build -o ${DIR}/${BINARY}-mac-arm ./cmd/main.go

test:
	go test ./...

clean:
	echo "Cleaning up binaries"
	rm ./bin/*


.PHONY: build test