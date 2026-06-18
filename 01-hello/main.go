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

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	fmt.Printf(
		"host=%s pid=%d rank=%d size=%d\n",
		hostname,
		os.Getpid(),
		world.GetRank(),
		world.GetSize(),
	)
}
