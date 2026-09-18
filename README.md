# K8s-Autoscaling-POC

PoC de microsserviços em Go (`orders` e `kitchen`, comunicando-se via gRPC) para validar o comportamento do **Horizontal Pod Autoscaler (HPA)** do Kubernetes sob carga, com testes de carga em k6 e pipeline de CI/CD via GitHub Actions.

## Arquitetura

| Serviço  | Portas          | Papel                                                   |
|----------|-----------------|----------------------------------------------------------|
| `orders` | 8080 (gRPC), 8081 (HTTP) | Recebe e armazena pedidos (em memória)          |
| `kitchen`| 9000 (HTTP)     | Consulta/cria pedidos via cliente gRPC de `orders`       |

## Pré-requisitos

- Go 1.26+
- Docker e Docker Compose
- kubectl + um cluster local (minikube ou kind)
- [k6](https://k6.io/) (para os testes de carga)
- `protoc` + plugins Go (apenas se for regenerar os `.proto`)

## Rodando localmente (sem containers)

```bash
make run-orders    # sobe orders em :8080 (gRPC) e :8081 (HTTP)
make run-kitchen   # sobe kitchen em :9000 (HTTP)
```

## Rodando com Docker Compose

```bash
docker compose up --build
```

- `orders` fica disponível em `localhost:8080` (gRPC) e `localhost:8081` (HTTP)
- `kitchen` fica disponível em `localhost:9000` (HTTP)

## Rodando no Kubernetes

Os manifests (`Deployment`, `Service` e `HorizontalPodAutoscaler`) esperam as imagens `guilhermemarques18/orders` e `guilhermemarques18/kitchen` já disponíveis para o cluster, já que usam `imagePullPolicy: Never`.

1. Construa as imagens localmente e carregue-as no cluster:

   ```bash
   docker build -t guilhermemarques18/orders:latest -f services/orders/Dockerfile .
   docker build -t guilhermemarques18/kitchen:latest -f services/kitchen/Dockerfile .

   # kind:
   kind load docker-image guilhermemarques18/orders:latest guilhermemarques18/kitchen:latest

   # minikube:
   minikube image load guilhermemarques18/orders:latest
   minikube image load guilhermemarques18/kitchen:latest
   ```

2. Aplique os manifests:

   ```bash
   kubectl apply -f k8s/
   ```

3. Acompanhe o autoscaling:

   ```bash
   kubectl get hpa -w
   kubectl top pods
   ```

> **Nota:** o cluster precisa ter o [metrics-server](https://github.com/kubernetes-sigs/metrics-server) instalado — sem ele o HPA não consegue ler métricas de CPU/memória e as réplicas ficam travadas no mínimo.

## Testes de carga (k6)

```bash
# smoke test rápido
k6 run --env BASE_URL=http://localhost:9000 k6/smoke-test.js

# teste de carga progressiva (ramping arrival rate)
k6 run --env BASE_URL=http://kitchen:9000 k6/load-test.js
```

## CI/CD

O workflow `.github/workflows/go.yml` roda automaticamente em push/PR para `main`:

1. **build**: `go build` + `go test -race` em todos os módulos
2. **docker**: build e push das imagens de `orders` e `kitchen` para o Docker Hub (tags `latest` e `<sha>`), executado apenas em push

Requer os secrets `DOCKERHUB_USERNAME` e `DOCKERHUB_TOKEN` configurados no repositório.

## Variáveis de ambiente

| Serviço  | Variável                     | Padrão              | Descrição                                   |
|----------|------------------------------|----------------------|----------------------------------------------|
| orders   | `HTTP_ADDR`                  | `:8081`              | Porta do servidor HTTP                        |
| orders   | `GRPC_ADDR`                  | `:8080`              | Porta do servidor gRPC                        |
| orders   | `ORDER_WORK_ITERATIONS_CREATE`| `20000`             | Iterações de hash sintéticas em `CreateOrder` |
| orders   | `ORDER_WORK_ITERATIONS_GET`  | `5000`               | Iterações de hash sintéticas em `GetOrders`   |
| kitchen  | `HTTP_ADDR`                  | `:9000`              | Porta do servidor HTTP                        |
| kitchen  | `ORDERS_GRPC_ADDR`           | `localhost:8080`     | Endereço gRPC do serviço `orders`             |