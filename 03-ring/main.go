package main

import (
	"fmt"

	mpi "github.com/mnlphlp/gompi"
)

func main() {
	mpi.Init()
	defer mpi.Finalize()

	world := mpi.NewComm(true)
	rank := world.GetRank()
	size := world.GetSize()

	next := (rank + 1) % size
	prev := (rank - 1 + size) % size

	token := []int{0}

	if rank == 0 {
		token[0] = 100
		fmt.Printf("[rank %d] enviando token=%d para rank %d\n", rank, token[0], next)
		world.Send(token, next, 0)

		world.Recv(&token, prev, 0)
		fmt.Printf("[rank %d] recebeu token=%d de rank %d (após volta completa)\n", rank, token[0], prev)
	} else {
		world.Recv(&token, prev, 0)
		fmt.Printf("[rank %d] recebeu token=%d de rank %d\n", rank, token[0], prev)

		token[0] += rank
		fmt.Printf("[rank %d] enviando token=%d para rank %d\n", rank, token[0], next)
		world.Send(token, next, 0)
	}
}
