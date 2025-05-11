# SubPub gRPC-сервис 
Сервис реализует событийную шину (pub-sub), где можно публиковать события по ключу и подписываться на них по этому же ключу. Сервис реализован на Go с использованием gRPC.

## Как работает сервис

Сервис состоит из двух основных частей: реализации pub-sub и gRPC-сервера, который работает поверх этой шины.

### 1. Пакет `subpub`

Это реализация модели Publisher-Subscriber. Можно подписаться на события по ключу (subject) и получать сообщения, которые публикуются по этому же ключу. На один ключ можно подписывать сразу несколько подписчиков, и все они будут получать сообщения в порядке публикации (FIFO).

Для каждого подписчика создаётся своя очередь, поэтому если кто-то медленно обрабатывает сообщения, то он не блокирует остальных. Есть метод `Close(ctx)`, который завершает работу, и если переданный контекст отменён, выходим сразу, при этом работающие подписчики продолжают обрабатывать сообщения.

### 2. gRPC-сервис

Снаружи это gRPC-сервис с двумя основными методами:

* `Subscribe` — сервер стримит события, на которые подписан клиент;
* `Publish` — позволяет публиковать события по ключу.

По сути, это обёртка над subpub, и всё, что приходит в Publish, передаётся в subpub.Publish, а все подписки обрабатываются через subpub.Subscribe.


## Технологии
- Go
- gRPC
- Protocol Buffers
- Taskfile 
- Тесты на Go 



## Как склонировать и запустить

### 1. Склонировать репозиторий

```bash
git clone https://github.com/your-username/your-repo.git
cd your-repo
```

### 2. Установить зависимости

```bash
go mod tidy
```

### 3. Установить Protocol Buffers

```bash
sudo apt install -y protobuf-compiler
```

### 4. Установить Go-плагины для gRPC

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Добавить `$GOPATH/bin` в `PATH`, чтобы команды `protoc-gen-go` и `protoc-gen-go-grpc` были доступны:

```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc
```

### 5. Установить Taskfile

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

---

## Использование Taskfile

### Генерация gRPC и protobuf файлов

```bash
task gen-proto
```

Task выполнит следующую команду:

```bash
protoc proto/subpub.proto --go_out=./proto --go-grpc_out=./proto
```

---

## Запуск тестов

```bash
go test ./internal/subpub
```

---

## Запуск приложения

```bash
go run cmd/main.go
```

