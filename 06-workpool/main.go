package main

import (
	"fmt"
	"time"

	mpi "github.com/mvneves/gompi"
	"github.com/mvneves/gompi/comm"
)

const (
	tagTask   = 1
	tagResult = 2

	numTasks = 20
)

func main() {
	mpi.Init()
	defer mpi.Finalize()

	world := mpi.NewComm(true)
	rank := world.GetRank()
	size := world.GetSize()

	if size < 2 {
		if rank == 0 {
			fmt.Println("Este exemplo requer pelo menos 2 processos (1 coordenador + N trabalhadores).")
		}
		return
	}

	if rank == 0 {
		coordinator(world, size)
	} else {
		worker(world, rank)
	}
}

// coordinator distribui tarefas sob demanda e coleta resultados.
// Envia taskID >= 0 para processar, taskID = -1 para encerrar.
func coordinator(world comm.Communicator, size int) {
	numWorkers := size - 1
	results := make(map[int]int)

	nextTask := 0

	// Distribui uma tarefa inicial para cada trabalhador
	for w := 1; w <= numWorkers; w++ {
		if nextTask < numTasks {
			world.Send([]int{nextTask}, w, tagTask)
			nextTask++
		} else {
			world.Send([]int{-1}, w, tagTask)
		}
	}

	// Recebe resultados e distribui tarefas restantes
	completed := 0
	for completed < numTasks {
		buf := make([]int, 3) // [workerRank, taskID, resultado]
		world.Recv(&buf, mpi.ANY_SOURCE, tagResult)
		workerRank := buf[0]
		taskID := buf[1]
		result := buf[2]

		results[taskID] = result
		completed++

		// Envia próxima tarefa ou sinal de término
		if nextTask < numTasks {
			world.Send([]int{nextTask}, workerRank, tagTask)
			nextTask++
		} else {
			world.Send([]int{-1}, workerRank, tagTask)
		}
	}

	fmt.Println("\n=== Resultados ===")
	for i := 0; i < numTasks; i++ {
		fmt.Printf("  tarefa %2d → %d² = %d\n", i, i, results[i])
	}
}

// worker recebe tarefas, processa, e envia resultados até receber taskID = -1.
func worker(world comm.Communicator, rank int) {
	for {
		task := make([]int, 1)
		world.Recv(&task, 0, tagTask)
		taskID := task[0]

		if taskID < 0 {
			fmt.Printf("[worker %d] finalizado\n", rank)
			return
		}

		// Simula processamento com custo variável
		time.Sleep(time.Duration(10+taskID*3) * time.Millisecond)
		result := taskID * taskID

		fmt.Printf("[worker %d] tarefa %d → resultado %d\n", rank, taskID, result)
		world.Send([]int{rank, taskID, result}, 0, tagResult)
	}
}
