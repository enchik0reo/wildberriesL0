# Test task L0
## Run NUTS-Streaming + PostgreSQL:
```
docker-compose up -d
```
## Run reader:
```
go run ./cmd/reader/main.go
```
## Run sender:
```
go run ./cmd/sender/main.go
```
## Get order by uid from cache:
```
localhost:4000/order/:uid
```

## Stress Tests
### WRK

```
wrk -t2 -c10 -d30s  http://localhost:4000/order/order1

Running 30s test @ http://localhost:4000/order/order1
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   381.29us    0.90ms  28.60ms   94.35%
    Req/Sec    23.86k     2.48k   31.90k    72.50%
  1425656 requests in 30.03s, 1.60GB read
Requests/sec:  47467.26
Transfer/sec:     54.59MB
```

### Vegeta
```
echo "GET http://localhost:4000/order/order1" | vegeta attack -duration=30s -rate=1000 | vegeta report

Requests      [total, rate, throughput]  30000, 1000.02, 1000.01
Duration      [total, attack, wait]      29.99967847s, 29.999327688s, 350.782µs
Latencies     [mean, 50, 95, 99, max]    297.348µs, 272.828µs, 453.958µs, 909.593µs, 15.280556ms
Bytes In      [total, mean]              32610000, 1087.00
Bytes Out     [total, mean]              0, 0.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:30000  
Error Set:
```