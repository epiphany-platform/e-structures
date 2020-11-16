build:
	go build -x -o ./output/producer github.com/mkyc/go-stucts-versioning-tests/cmd/producer

clean:
	rm -rf ./output