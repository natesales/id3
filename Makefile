all: clean build
	./id3

build:
	go build -o id3 -v

clean:
	go clean
	go fmt