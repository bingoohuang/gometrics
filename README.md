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

## Usage

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