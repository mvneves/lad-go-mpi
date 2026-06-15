package main

import (
	"fmt"
	"math/rand"

	mpi "github.com/mnlphlp/gompi"
)

func main() {
	mpi.Init()
	defer mpi.Finalize()

	world := mpi.NewComm(true)
	rank := world.GetRank()
	size := world.GetSize()

	const totalSamples = 10_000_000
	samplesPerProc := totalSamples / size

	// Cada processo usa uma seed diferente baseada no rank
	rng := rand.New(rand.NewSource(int64(rank*7919 + 42)))

	localCount := 0
	for i := 0; i < samplesPerProc; i++ {
		x := rng.Float64()
		y := rng.Float64()
		if x*x+y*y <= 1.0 {
			localCount++
		}
	}

	fmt.Printf("[rank %d] %d amostras, %d dentro do círculo\n", rank, samplesPerProc, localCount)

	// Reduce manual: todos enviam para rank 0
	if rank == 0 {
		globalCount := localCount
		for src := 1; src < size; src++ {
			buf := make([]int, 1)
			world.Recv(&buf, src, 0)
			globalCount += buf[0]
		}

		pi := 4.0 * float64(globalCount) / float64(totalSamples)
		fmt.Printf("\n=== Resultado ===\n")
		fmt.Printf("Amostras totais:  %d\n", totalSamples)
		fmt.Printf("Processos MPI:    %d\n", size)
		fmt.Printf("Pi estimado:      %.10f\n", pi)
	} else {
		world.Send([]int{localCount}, 0, 0)
	}
}
