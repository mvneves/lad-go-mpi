BINARY = teste-mpi

.PHONY: build clean

build:
	go mod tidy
	go build -o $(BINARY)

clean:
	rm -f $(BINARY)
