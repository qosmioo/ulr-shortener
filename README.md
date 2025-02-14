## Технологии
- **PostgreSQL** и **Redis** в качестве хранилища ссылок
- **docker-compose** развертывание приложения
- **gRPC** и **HTTP** протоколы
- **Zap** логгирование
- **Чистая архитектура**: Delivery -> Usecase -> Storage 

## Структура проекта: 

```
Url Shortener
├── cmd
│   └── main.go
├── config
│   ├── config.go
│   └── config.yaml
├── internal
│   ├── delivery
│   │   ├── grpc.go
│   │   └── http.go
│   ├── storage
│   │   ├── postgres.go
│   │   ├── redis.go
│   │   ├── repository.go
│   │   ├── postgres_test.go      
│   │   └── redis_test.go        
│   └── usecase
│       ├── url_shortener.go
│       └── url_shortener_test.go
├── pkg
│   ├── shortener
│   │   └── shortener.go
│   └── validator
│       └── validator.go
├── proto
│   ├── shortener.proto
│   └── shortener_grpc.pb.go
├── Dockerfile
├── docker-compose.yml
├── .env
├── Makefile
├── README.md
├── go.mod
├── go.sum
```

## Запуск приложения

### Запуск 
```bash
make run
```
### Примеры использования

#### Создание короткой ссылки по http
```bash
curl -X POST -H "Content-Type: application/json" -d '{"original_url": "https://www.google.com"}' http://localhost:8080/v1/urls
```

#### Получение оригинальной ссылки по http
```bash
curl http://localhost:8080/v1/urls/1234567890
```

#### Создание короткой ссылки по grpc (с ипользованием grpcurl)
```bash
grpcurl -plaintext -d '{"longUrl": "http://example.com"}' localhost:8001 url_shortener.GRPCHandler/ShortenUrl
```

#### Получение оригинальной ссылки по grpc (с ипользованием grpcurl)
```bash
grpcurl -plaintext -d '{"shortUrl": "1234567890"}' localhost:8001 url_shortener.GRPCHandler/GetUrl
```

### Остановка
```bash
make down
```
#### Важно: по умолчанию используется Redis, для использования PostgreSQL необходимо переопределить переменную type в config.yaml на postgres

## Запуск тестов    

```bash
make test
```

## Что можно улучшить
- Добавить тесты для HTTP сервера
- Добавить тесты для gRPC сервера
- Валидировать входные данные и ициниализировать только нужное хранилище
- Добавить request id Middleware

