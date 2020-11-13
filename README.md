# gometrics

metrics golang client library.

## metrics

\# | TYPE         | Meaning
---|--------------|--------
1  | RT           | 平均响应时间
2  | QPS          | 业务量(次数)
3  | SUCCESS_RATE | 成功率
4  | FAIL_RATE    | 失败率
5  | HIT_RATE     | 命中率
6  | CUR          | 瞬时值

## HB

心跳

\# | TYPE | Meaning
---|------|--------
1  | HB   | 一次心跳

## Client Usage

### 准备参数

1. 通过.env环境文件设置，优先级最高。在当前目录下创建.env文件，设定一些参数， eg.

    ```properties
    # 应用名称，默认使用当前pid
    APP_NAME=bingoohuangapp
    # 写入指标日志的间隔时间，默认1s
    METRICS_INTERVAL=1s
    # 写入心跳日志的间隔时间，默认20s
    HB_INTERVAL=20s
    # Metrics对象的处理容量，默认1000，来不及处理时，超额扔弃处理
    CHAN_CAP=1000
    # 日志存放的目录，默认/tmp/log/metrics
    LOG_PATH=/var/log/footstone/metrics
    # 日志文件最大保留天数
    MAX_BACKUPS=7
    ```

1. 通过命令行环境变量设置

    eg. `APP_NAME=demo demoproc`

1. 通过命令行指定环境文件名

    eg. `ENV_FILE=testdata/golden.env demoproc`

### RT 平均响应时间

```go
import (
	"github.com/bingoohuang/gometrics/metric"
)

func YourBusinessDemo1() {
    defer metric.RT("key1", "key2", "key3").Record()

    // business logic
}

func YourBusinessDemo2() {
    rt := metric.RT("key1", "key2", "key3")

    // business logic
    start := time.Now()
    // ...
    rt.RecordSince(start)
}
```

### QPS 业务量(次数)

```go
func YourBusinessDemoQPS() {
    metric.QPS("key1", "key2", "key3").Record(1 /* 业务量 */ )
}
```

or in simplified way:

```go
func YourBusinessDemoQPS() {
    metric.QPS1("key1", "key2", "key3")
}
```

### SUCCESS_RATE 成功率

```go
func YourBusinessDemoSuccessRate() {
    sr := metric.SuccessRate("key1", "key2", "key3")
    defer sr.IncrTotal()

    // business logic
    sr.IncrSuccess()
}
```

### FAIL_RATE 失败率

```go
func YourBusinessDemoFailRate() {
    fr := metric.FailRate("key1", "key2", "key3")
    defer fr.IncrTotal()

    // business logic
    fr.IncrFail()
}
```

### HIT_RATE 命中率

```go
func YourBusinessDemoHitRate() {
    fr := metric.HitRate("key1", "key2", "key3")
    defer fr.IncrTotal()

    // business logic
    fr.IncrHit()
}
```

### CUR 瞬时值

```go
func YourBusinessDemoCur() {
    // business logic
    metric.Cur("key1", "key2", "key3").Record(100)
    // business logic
}
```

### Demo

1. build `go install -ldflags="-s -w" ./...`
1. or build for linux

    - `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags="-s -w" ./...`
    - `upx ~/go/bin/linux_amd64/gometricsdemo`
    - `bssh scp -H A-gw-test2 ~/go/bin/linux_amd64/gometricsdemo r:./bingoohuang/gometrics`

1. run `ENV_FILE=testdata/golden.env gometricsdemo`

```bash
$ tail -f /tmp/metricslog/metrics-hb.bingoohuangapp.log
{"time":"20200205162411000","key":"bingoohuangapp.hb","hostname":"192.168.10.101","logtype":"HB","v1":1,"v2":0,"min":0,"max":0}
{"time":"20200205162431000","key":"bingoohuangapp.hb","hostname":"192.168.10.101","logtype":"HB","v1":1,"v2":0,"min":0,"max":0}
{"time":"20200205162451000","key":"bingoohuangapp.hb","hostname":"192.168.10.101","logtype":"HB","v1":1,"v2":0,"min":0,"max":0}
```

```bash
$ tail -f /tmp/metricslog/metrics-key.bingoohuangapp.log
{"time":"20200205162628000","key":"key1#key2#key3","hostname":"192.168.10.101","logtype":"FAIL_RATE","v1":0,"v2":2,"min":0,"max":100}
{"time":"20200205162628000","key":"key1#key2#key3","hostname":"192.168.10.101","logtype":"HIT_RATE","v1":1,"v2":2,"min":0,"max":100}
{"time":"20200205162628000","key":"key1#key2#key3","hostname":"192.168.10.101","logtype":"CUR","v1":100,"v2":0,"min":0,"max":0}
{"time":"20200205162628000","key":"key1#key2#key3","hostname":"192.168.10.101","logtype":"RT","v1":193,"v2":1,"min":0,"max":811}
```

## benchmark

```bash
$ go test -bench=.  ./...
WARN[0000] loading env file error open .env: no such file or directory
INFO[0000] log file /tmp/log/metrics/metrics-key.44739.log created
INFO[0000] log file /tmp/log/metrics/metrics-hb.44739.log created
/Users/bingoo/GitHub/gometrics/metric
goos: darwin
goarch: amd64
pkg: github.com/bingoohuang/gometrics/metric
BenchmarkRT-12                   1803442               655 ns/op
BenchmarkQPS-12                  2232487               538 ns/op
BenchmarkSuccessRate-12          2175163               552 ns/op
BenchmarkFailRate-12             2246766               516 ns/op
BenchmarkHitRate-12              2110915               597 ns/op
BenchmarkCur-12                  3023659               388 ns/op
PASS
ok      github.com/bingoohuang/gometrics/metric 11.385s
```

## cloc

```bash
# bingoo @ 192 in ~/GitHub/gometrics on git:master o [16:54:30]
$ go get -u github.com/hhatto/gocloc/cmd/gocloc

# bingoo @ 192 in ~/GitHub/gometrics on git:master x [16:54:49]
$ gocloc .
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Go                              12            230             82            788
XML                              4              0              0            289
Markdown                         3             46              0            134
-------------------------------------------------------------------------------
TOTAL                           19            276             82           1211
-------------------------------------------------------------------------------

# bingoo @ 192 in ~/GitHub/gometrics on git:master o [16:55:03]
$ date
2020年 2月 5日 星期三 16时56分33秒 CST
```
