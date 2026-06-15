# Go + MPI no Cluster Atlântica

Programa de teste para verificar o funcionamento de **Go com MPI** no cluster Atlântica (PUCRS).

## Pré-requisitos

O cluster já possui o ambiente necessário:

| Componente | Versão | Caminho |
|---|---|---|
| Go | 1.23.2 | `/LADAPPs/Go/go-1.23.2/bin/go` |
| OpenMPI | 4.1.1 | `/LADAPPs/OpenMPI/openmpi-4.1.1/` |

## Passo a passo

### 1. Clonar o repositório

```bash
git clone https://github.com/mvneves/lad-go-mpi.git
cd lad-go-mpi
```

### 2. Configurar o ambiente

```bash
source setup.sh
```

Este script configura o `PKG_CONFIG_PATH` e verifica se Go e OpenMPI estão acessíveis.

Se preferir configurar manualmente:

```bash
export PKG_CONFIG_PATH=/LADAPPs/OpenMPI/openmpi-4.1.1/lib/pkgconfig
```

### 3. Compilar

```bash
make build
```

Ou diretamente:

```bash
go mod tidy
go build -o teste-mpi
```

Se a compilação não produzir saída, funcionou.

### 4. Alocar recursos no SLURM

Antes de executar, é necessário reservar recursos no cluster.

**Um nó, 4 processos:**

```bash
salloc -N 1 -n 4
```

**Dois nós, 8 processos:**

```bash
salloc -N 2 -n 8
```

A saída será parecida com:

```text
salloc: Granted job allocation 32806
```

### 5. Executar

Com a alocação ativa, execute o programa com `mpirun`:

```bash
mpirun -np 4 ./teste-mpi
```

Saída esperada (1 nó, 4 processos):

```text
host=atlantica01 pid=12345 rank=0 size=4
host=atlantica01 pid=12346 rank=1 size=4
host=atlantica01 pid=12347 rank=2 size=4
host=atlantica01 pid=12348 rank=3 size=4
```

Saída esperada (2 nós, 8 processos):

```text
host=atlantica01 pid=12345 rank=0 size=8
host=atlantica01 pid=12346 rank=1 size=8
host=atlantica01 pid=12347 rank=2 size=8
host=atlantica01 pid=12348 rank=3 size=8
host=atlantica03 pid=22345 rank=4 size=8
host=atlantica03 pid=22346 rank=5 size=8
host=atlantica03 pid=22347 rank=6 size=8
host=atlantica03 pid=22348 rank=7 size=8
```

A ordem das linhas pode variar — isso é normal em programas paralelos.

### 6. Liberar a alocação

Ao terminar, saia do shell alocado:

```bash
exit
```

A mensagem esperada é:

```text
salloc: Relinquishing job allocation 32806
```

Para verificar se há alocações ativas:

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

Isso significa que o MPI tentou iniciar mais processos do que os recursos alocados. Certifique-se de que o número em `-np` corresponde ao `-n` usado no `salloc`.

### Métodos não encontrados (`undefined: mpi.Start`, etc.)

Certifique-se de usar a API correta do pacote `github.com/mnlphlp/gompi`:

```go
mpi.Init()
defer mpi.Finalize()

world := mpi.NewComm(true)
world.GetRank()
world.GetSize()
```
