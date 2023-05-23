# Run nats-streaming + postgreSQL:
```
docker-compose up -d
```
# Run reader:
```
go run ./cmd/reader/main.go
```
# Run sender:
```
go run ./cmd/sender/main.go
```
# Get order by uid from cache:
```
localhost:4000/order/:uid