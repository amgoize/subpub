# SubPub gRPC-сервис 
Сервис реализует событийную шину (pub-sub), где можно публиковать события по ключу и подписываться на них по этому же ключу. Сервис реализован на Go с использованием gRPC.

## 🛠 Технологии
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

### 5. Установить Taskfile (один раз)

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

