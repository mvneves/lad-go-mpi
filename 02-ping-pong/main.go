package main

import (
	"fmt"
	"os"

	mpi "github.com/mvneves/gompi"
)

func main() {
	mpi.Init()
	defer mpi.Finalize()

	world := mpi.NewComm(true)
	rank := world.GetRank()
	size := world.GetSize()

	if size < 2 {
		if rank == 0 {
			fmt.Fprintln(os.Stderr, "Este exemplo requer pelo menos 2 processos.")
		}
		return
	}

	const numPings = 10

	if rank == 0 {
		for i := 0; i < numPings; i++ {
			msg := []int{i}
			world.Send(msg, 1, 0)

			reply := make([]int, 1)
			world.Recv(&reply, 1, 1)
			fmt.Printf("[rank %d] ping %d → pong %d\n", rank, i, reply[0])
		}
	} else if rank == 1 {
		for i := 0; i < numPings; i++ {
			msg := make([]int, 1)
			world.Recv(&msg, 0, 0)

			reply := []int{msg[0] * 2}
			world.Send(reply, 0, 1)
			fmt.Printf("[rank %d] recebeu %d → respondeu %d\n", rank, msg[0], reply[0])
		}
	}
}
