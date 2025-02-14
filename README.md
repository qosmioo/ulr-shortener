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