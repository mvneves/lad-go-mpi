package main

import (
	"fmt"

	mpi "github.com/mvneves/gompi"
)

func main() {
	mpi.Init()
	defer mpi.Finalize()

	world := mpi.NewComm(true)
	rank := world.GetRank()
	size := world.GetSize()

	const chunkSize = 4
	totalSize := chunkSize * size

	// Scatter: rank 0 distribui partes iguais do array para cada processo
	local := make([]int, chunkSize)
	if rank == 0 {
		data := make([]int, totalSize)
		for i := range data {
			data[i] = i + 1
		}
		fmt.Printf("[rank 0] dados originais: %v\n", data)

		copy(local, data[:chunkSize])
		for dest := 1; dest < size; dest++ {
			start := dest * chunkSize
			world.Send(data[start:start+chunkSize], dest, 0)
		}
	} else {
		world.Recv(&local, 0, 0)
	}

	fmt.Printf("[rank %d] recebeu: %v\n", rank, local)

	// Cada processo calcula a soma local
	localSum := 0
	for _, v := range local {
		localSum += v
	}
	fmt.Printf("[rank %d] soma local: %d\n", rank, localSum)

	// Gather: rank 0 coleta os resultados parciais
	if rank == 0 {
		results := make([]int, size)
		results[0] = localSum
		for src := 1; src < size; src++ {
			buf := make([]int, 1)
			world.Recv(&buf, src, 1)
			results[src] = buf[0]
		}
		total := 0
		for _, v := range results {
			total += v
		}
		fmt.Printf("[rank 0] somas parciais: %v → total: %d\n", results, total)
	} else {
		world.Send([]int{localSum}, 0, 1)
	}
}
