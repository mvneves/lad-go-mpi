# Go + MPI no Cluster Atlântica

Exemplos de programas em **Go** usando **MPI** para o cluster Atlântica (PUCRS).

## Exemplos

| # | Diretório | Conceito |
|---|-----------|----------|
| 01 | `01-hello` | Inicialização MPI, rank e size |
| 02 | `02-ping-pong` | Comunicação ponto-a-ponto (Send/Recv) |
| 03 | `03-ring` | Topologia em anel, roteamento por rank |
| 04 | `04-scatter-gather` | Distribuição e coleta de dados |
| 05 | `05-reduce-pi` | Redução (estimativa de Pi por Monte Carlo) |
| 06 | `06-workpool` | Distribuição dinâmica de tarefas |

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

Compilar todos os exemplos:

```bash
make
```

Compilar um exemplo específico:

```bash
go build -o bin/01-hello ./01-hello
```

Os binários são gerados na pasta `bin/`.

### 4. Alocar recursos no SLURM

Antes de executar, reserve recursos no cluster:

```bash
salloc -N 1 -n 4
```

Para exemplos que precisam de mais nós:

```bash
salloc -N 2 -n 8
```

### 5. Executar

Com a alocação ativa, execute o exemplo desejado:

```bash
mpirun -np 4 ./bin/01-hello
mpirun -np 2 ./bin/02-ping-pong
mpirun -np 4 ./bin/03-ring
mpirun -np 4 ./bin/04-scatter-gather
mpirun -np 4 ./bin/05-reduce-pi
mpirun -np 4 ./bin/06-workpool
```

### 6. Liberar a alocação

Ao terminar, saia do shell alocado:

```bash
exit
```

Para verificar se há alocações ativas:

```bash
squeue -u $USER
```

## Descrição dos exemplos

### 01-hello

Programa mínimo: cada processo imprime seu hostname, PID, rank e o tamanho total do comunicador. Sem comunicação entre processos.

### 02-ping-pong

Rank 0 envia uma sequência de valores para rank 1, que responde com o dobro. Demonstra **Send/Recv** bloqueante entre dois processos.

### 03-ring

Cada processo recebe um token do anterior e envia para o próximo, formando um anel. O token acumula a soma dos ranks. Demonstra **roteamento baseado em rank** e topologia lógica.

### 04-scatter-gather

Rank 0 divide um array em partes iguais e distribui para todos os processos (scatter manual). Cada um calcula a soma local e devolve o resultado para rank 0 (gather manual). Demonstra o **padrão de distribuição e coleta de dados**.

### 05-reduce-pi

Cada processo gera pontos aleatórios para estimar Pi pelo método de Monte Carlo. Os resultados parciais são enviados para rank 0, que calcula a estimativa final (reduce manual). Demonstra **paralelismo embaraçosamente paralelo** e redução.

### 06-workpool

Rank 0 atua como coordenador, distribuindo tarefas sob demanda para os demais processos. Quando um trabalhador termina, recebe a próxima tarefa disponível. Demonstra **balanceamento dinâmico de carga** (work pool pattern).

## Erros comuns

### `No package 'ompi' found`

```bash
export PKG_CONFIG_PATH=/LADAPPs/OpenMPI/openmpi-4.1.1/lib/pkgconfig
```

Ou execute `source setup.sh`.

### `not enough slots available`

Certifique-se de que o número em `-np` corresponde ao `-n` usado no `salloc`.

### Métodos não encontrados (`undefined: mpi.Start`, etc.)

Use a API correta do pacote `github.com/mnlphlp/gompi`:

```go
mpi.Init()
defer mpi.Finalize()

world := mpi.NewComm(true)
world.GetRank()
world.GetSize()
```

## Referências

- [Manual do LAD — Laboratório de Alto Desempenho (PUCRS)](https://lad-pucrs.github.io/docs/)
