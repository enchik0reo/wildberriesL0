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
echo "GET http://localhost:4000/order/order1" | vegeta attack -duration=30s -rate=10000 | vegeta report

Requests      [total, rate, throughput]  300000, 10000.36, 10000.31
Duration      [total, attack, wait]      29.99908055s, 29.998914501s, 166.049µs
Latencies     [mean, 50, 95, 99, max]    193.004µs, 147.926µs, 438.727µs, 1.034802ms, 13.220574ms
Bytes In      [total, mean]              300000000, 1000.00
Bytes Out     [total, mean]              0, 0.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:300000  
Error Set:
```
![Plot](./images/vegeta-plot.jpg?raw=true)