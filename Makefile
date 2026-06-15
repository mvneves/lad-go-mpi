BINARY = teste-mpi
NP     ?= 4
NODES  ?= 1

.PHONY: deps build run run-slurm clean

deps:
	go mod tidy

build: deps
	go build -o $(BINARY)

run: build
	mpirun --oversubscribe -np $(NP) ./$(BINARY)

run-slurm: build
	@echo "Solicitando alocação: $(NODES) nó(s), $(NP) tarefa(s)..."
	salloc -N $(NODES) -n $(NP) mpirun -np $(NP) ./$(BINARY)

clean:
	rm -f $(BINARY)
