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
## Get order by uid:
```
localhost:4000/order/:uid
```
(uid in this example: order[0-100])

## Stress Tests
### WRK
First launch:
```
wrk -t4 -c4 -d30s  http://localhost:4000/order/order1

Running 30s test @ http://localhost:4000/order/order1
  4 threads and 4 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   146.77us  414.00us  11.75ms   97.26%
    Req/Sec    10.59k     0.89k   13.13k    63.21%
  1268740 requests in 30.10s, 1.38GB read
Requests/sec:  42151.00
Transfer/sec:     46.87MB
```
Last launch:
```
wrk -t4 -c8192 -d30s  http://localhost:4000/order/order1

Running 30s test @ http://localhost:4000/order/order1
  4 threads and 8192 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   162.60ms   58.32ms 521.04ms   83.33%
    Req/Sec    12.40k     2.16k   18.69k    70.32%
  1465947 requests in 30.05s, 1.59GB read
  Socket errors: connect 8, read 0, write 0, timeout 0
Requests/sec:  48777.95
Transfer/sec:     54.24MB
```
Limits on the number of file descriptors: ulimit -n 8192
![wrk](./images/wrk-4-threads.png?raw=true)
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
