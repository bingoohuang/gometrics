# gometrics

metrics golang client 

\# | TYPE  | Meaning
---|---|---
1 |RT| 平均响应时间
2| QPS| 业务量(次数)
3 | SUCCESS_RATE | 成功率
4 | FAIL_RATE|失败率
5 | HIT_RATE| 命中率
6 | CUR | 瞬时值

## Client Usage

### 准备参数

1. 通过.env环境文件设置，优先级最高。在当前目录下创建.env文件，设定一些参数， eg.
        
    ```dotenv
    # 应用名称，默认使用当前pid
    APP_NAME=bingoohuangapp
    # 写入日志的间隔时间，默认1s
    INTERVAL=1s
    # 指标日志的处理容量，默认1000
    CHAN_CAP=1000
    # 指标日志存放的目录，默认/var/log/metrics
    LOG_PATH=/tmp/metricslog
    # 指标日志文件最大保留天数
    MAX_BACKUPS=7
    ```

1. 通过命令行环境变量设置

    eg. `APP_NAME=demo demoproc`

1. 通过命令行指定环境文件名

    eg. `ENV_FILE=testdata/golden.env demoproc`

### RT 平均响应时间

```go
func YourBusinessDemo1() {
    defer gometrics.RT("key1", "key2", "key3").Record()
    
    // business logic
}

func YourBusinessDemo2() {
    rt := gometrics.RT("key1", "key2", "key3")
    
    // business logic
    start := time.Now()
    // ...
    rt.RecordSince(start)
}
```

### QPS 业务量(次数)

```go
func YourBusinessDemoQPS() {
    gometrics.QPS("key1", "key2", "key3").Record(1 /* 业务量 */ )
}

```

### SUCCESS_RATE 成功率

```go
func YourBusinessDemoSuccessRate() {
    sr := gometrics.SuccessRate("key1", "key2", "key3")
    defer sr.IncrTotal()

    // business logic
    sr.IncrSuccess()
}
```

### FAIL_RATE 失败率

```go
func YourBusinessDemoFailRate() {
    fr := gometrics.FailRate("key1", "key2", "key3")
    defer fr.IncrTotal()

    // business logic
    fr.IncrFail()
}
```

### HIT_RATE 命中率

```go
func YourBusinessDemoHitRate() {
    fr := gometrics.HitRate("key1", "key2", "key3")
    defer fr.IncrTotal()

    // business logic
    fr.IncrHit()
}
```

### CUR 瞬时值

```go
func YourBusinessDemoCur() {
    // business logic
    gometrics.Cur("key1", "key2", "key3").Record(100)
    // business logic
}
```

### Demo

1. build `go fmt ./...;goimports -w .;golangci-lint run --enable-all;golint . ;go install -ldflags="-s -w" ./...`
1. run ` ENV_FILE=testdata/golden.env gometricsdemo`
