# Go + MPI no Cluster Atlântica

Programa de teste para verificar o funcionamento de **Go com MPI** no cluster Atlântica (PUCRS).

## Pré-requisitos

O cluster já possui o ambiente necessário:

| Componente | Versão | Caminho |
|---|---|---|
| Go | 1.23.2 | `/LADAPPs/Go/go-1.23.2/bin/go` |
| OpenMPI | 4.1.1 | `/LADAPPs/OpenMPI/openmpi-4.1.1/` |

## Início rápido

### 1. Clonar o repositório

```bash
git clone https://github.com/pucrs-fppd/go-mpi-cluster.git
cd go-mpi-cluster
```

### 2. Configurar o ambiente

```bash
source setup.sh
```

Este script configura o `PKG_CONFIG_PATH` e verifica se Go e OpenMPI estão acessíveis.

### 3. Compilar

```bash
make build
```

Ou diretamente:

```bash
go build -o teste-mpi
```

### 4. Executar (teste rápido)

Para um teste rápido sem alocação SLURM:

```bash
make run
```

Ou:

```bash
mpirun --oversubscribe -np 4 ./teste-mpi
```

Saída esperada:

```text
host=atlantica01 pid=12345 rank=0 size=4
host=atlantica01 pid=12346 rank=1 size=4
host=atlantica01 pid=12347 rank=2 size=4
host=atlantica01 pid=12348 rank=3 size=4
```

A ordem das linhas pode variar — isso é normal em programas paralelos.

### 5. Executar com SLURM (modo correto)

Solicite uma alocação e execute:

```bash
salloc -N 1 -n 4
mpirun -np 4 ./teste-mpi
exit
```

Para usar mais nós:

```bash
salloc -N 2 -n 8
mpirun -np 8 ./teste-mpi
exit
```

Ou via Makefile:

```bash
make run-slurm NP=8 NODES=2
```

### 6. Liberar alocação

Ao terminar, saia do shell alocado com `exit` ou `Ctrl+D`. Para verificar se há jobs ativos:

```bash
squeue -u $USER
```

## Sobre o programa

O `main.go` é um programa mínimo que:

1. Inicializa o ambiente MPI
2. Imprime o hostname, PID, rank e tamanho do comunicador
3. Finaliza o MPI

Cada execução com `mpirun -np N` inicia **N processos independentes** do mesmo binário. Cada processo recebe um identificador único chamado **rank** (de 0 a N-1).

## Erros comuns

### `No package 'ompi' found`

```bash
export PKG_CONFIG_PATH=/LADAPPs/OpenMPI/openmpi-4.1.1/lib/pkgconfig
```

Ou execute `source setup.sh`.

### `not enough slots available`

Use `--oversubscribe` para testes rápidos, ou aloque recursos com `salloc` antes de executar.

### Métodos não encontrados (`undefined: mpi.Start`, etc.)

Certifique-se de usar a API correta do pacote `github.com/mnlphlp/gompi`:

```go
mpi.Init()
defer mpi.Finalize()

world := mpi.NewComm(true)
world.GetRank()
world.GetSize()
```

## Referência rápida

```bash
source setup.sh              # configurar ambiente
make build                   # compilar
make run                     # executar local (4 processos)
make run NP=8               # executar local (8 processos)
make run-slurm NP=4 NODES=1 # executar via SLURM
make clean                   # remover binário
```
