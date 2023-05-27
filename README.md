# Test task L0
## Run NUTS-Streaming + PostgreSQL:
```
docker-compose up -d
```
## Run consumer:
```
go run ./cmd/consumer/main.go
```
## Run producer:
```
go run ./cmd/producer/main.go
```
## Get order by uid from cache:
```
localhost:4000/order/:uid
```
(uids in this example: order[0-100])

## Stress Tests
### WRK

```
wrk -t2 -c10 -d30s  http://localhost:4000/order/order1

Running 30s test @ http://localhost:4000/order/order1
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   364.98us  845.43us  23.81ms   94.43%
    Req/Sec    24.43k     2.31k   33.85k    70.17%
  1459650 requests in 30.03s, 1.52GB read
Requests/sec:  48607.95
Transfer/sec:     51.87MB
```

### Vegeta
```
echo "GET http://localhost:4000/order/order1" | vegeta attack -duration=30s -rate=1000 | vegeta report

Requests      [total, rate, throughput]  30000, 1000.04, 1000.03
Duration      [total, attack, wait]      29.999011235s, 29.998832912s, 178.323µs
Latencies     [mean, 50, 95, 99, max]    252.793µs, 222.682µs, 422.071µs, 625.057µs, 12.302956ms
Bytes In      [total, mean]              30000000, 1000.00
Bytes Out     [total, mean]              0, 0.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:30000  
Error Set:
```