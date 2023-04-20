<p align="center">
    <br> <a href="README.md">English</a> | 中文
</p>

# olog
* olog是一个轻量级、高性能、开箱即用的日志库，完全依赖于Go标准库。
* 支持以JSON和纯文本格式输出日志。
* 支持设置上下文处理函数，以从上下文中检索字段以进行输出。
* 支持七个日志级别：TRACE、DEBUG、INFO、NOTICE、WARN、ERROR和FATAL，对应于“trace”、“debug”、“info”、“notice”、“warn”、“error”和“fatal”标签。用户还可以定义自己的语义标签，如“slow”和“stat”。
* 提供输出开关控制，除了FATAL之外的所有日志级别都可以控制输出。
* 增强了对调用者打印的控制。用户可以在全局设置中禁用调用者打印以提高性能，但可以在打印某些关键日志时启用调用者打印支持。
* 为用户提供了一个Logger接口，以便轻松构建自己的日志记录器。

### 代码示例
```go
    Trace("hello world")
    Tracew("hello", Field{Key: "name", Value: "bob"})
    Debug("hello world")
    Infow("hello", Field{Key: "name", Value: "linda"}, Field{Key: "age", Value: 18})
    Noticef("hello %s", "world")
    Warnf("hello %s", "world")
    Warnw("hello", Field{Key: "order_no", Value: "AWESDDF"})
    Errorw("hello world", Field{Key: "success", Value: true})
    Logf(LogOption{Level: DEBUG, LevelTag: "print", EnableCaller: EnableClose,
        Fields: []Field{{Key: "price", Value: 32.5}},
        EnableStack: EnableOpen, StackSize: 1},
        "hello world")
    Fatal("fatal exit")
```
json输出如下：
```json
{"@timestamp":"2023-04-20T18:33:42+08:00","level":"trace","caller":"olog/log_test.go:377","content":"hello world","stack":"\ngithub.com/welllog/olog.TestPlainOutput\n\t/olog/log_test.go:377\ntesting.tRunner\n\t/usr/local/opt/go/libexec/src/testing/testing.go:1576\nruntime.goexit\n\t/usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1598"}
{"@timestamp":"2023-04-20T18:33:42+08:00","level":"trace","caller":"olog/log_test.go:378","content":"hello","name":"bob","stack":"\ngithub.com/welllog/olog.TestPlainOutput\n\t/olog/log_test.go:378\ntesting.tRunner\n\t/usr/local/opt/go/libexec/src/testing/testing.go:1576\nruntime.goexit\n\t/usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1598"}
{"@timestamp":"2023-04-20T18:33:42+08:00","level":"debug","caller":"olog/log_test.go:379","content":"hello world"}
{"@timestamp":"2023-04-20T18:33:42+08:00","level":"info","caller":"olog/log_test.go:380","content":"hello","name":"linda","age":18}
{"@timestamp":"2023-04-20T18:33:42+08:00","level":"notice","caller":"olog/log_test.go:381","content":"hello world"}
{"@timestamp":"2023-04-20T18:33:42+08:00","level":"warn","caller":"olog/log_test.go:382","content":"hello world"}
{"@timestamp":"2023-04-20T18:33:42+08:00","level":"warn","caller":"olog/log_test.go:383","content":"hello","order_no":"AWESDDF"}
{"@timestamp":"2023-04-20T18:33:42+08:00","level":"error","caller":"olog/log_test.go:384","content":"hello world","success":true}
{"@timestamp":"2023-04-20T18:33:42+08:00","level":"print","content":"hello world","price":32.5,"stack":"\ngithub.com/welllog/olog.TestPlainOutput\n\t/olog/log_test.go:385"}
{"@timestamp":"2023-04-20T18:33:42+08:00","level":"fatal","caller":"olog/log_test.go:388","content":"fatal exit"}
```
plain输出如下：
```
2023-04-20T18:32:09+08:00	trace	olog/log_test.go:377	hello world	stack=
github.com/welllog/olog.TestPlainOutput
	/olog/log_test.go:377
testing.tRunner
	/usr/local/opt/go/libexec/src/testing/testing.go:1576
runtime.goexit
	/usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1598
2023-04-20T18:32:09+08:00	trace	olog/log_test.go:378	hello	name=bob	stack=
github.com/welllog/olog.TestPlainOutput
	/olog/log_test.go:378
testing.tRunner
	/usr/local/opt/go/libexec/src/testing/testing.go:1576
runtime.goexit
	/usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1598
2023-04-20T18:32:09+08:00	debug	olog/log_test.go:379	hello world
2023-04-20T18:32:09+08:00	info	olog/log_test.go:380	hello	name=linda	age=18
2023-04-20T18:32:09+08:00	notice	olog/log_test.go:381	hello world
2023-04-20T18:32:09+08:00	warn	olog/log_test.go:382	hello world
2023-04-20T18:32:09+08:00	warn	olog/log_test.go:383	hello	order_no=AWESDDF
2023-04-20T18:32:09+08:00	error	olog/log_test.go:384	hello world	success=true
2023-04-20T18:32:09+08:00	print	hello world	price=32.5	stack=
github.com/welllog/olog.TestPlainOutput
	/olog/log_test.go:385
2023-04-20T18:32:09+08:00	fatal	olog/log_test.go:388	fatal exit
```

### contextLogger使用
```go
        SetDefCtxHandle(func(ctx context.Context) []Field {
		var fs []Field
		uid, ok := ctx.Value("uid").(int)
		if ok {
			fs = append(fs, Field{Key: "uid", Value: uid})
		}

		name, ok := ctx.Value("name").(string)
		if ok {
			fs = append(fs, Field{Key: "name", Value: name})
		}
		return fs
	})
        logger := WithContext(GetLogger(), context.WithValue(context.Background(), "uid", 3))
	logger.Debug("test")
		
        logger = WithContext(logger, context.WithValue(context.Background(), "name", "bob"))
        logger = WithEntries(logger, map[string]any{"requestId": "ae32fec"})
	logger.Debug("test 2")
```

### 实现自己的logger
```
type CustomLogger struct {
	Logger
}

func (l *CustomLogger) Slow(a ...any) {
    l.Log(LogOption{Level: WARN, LevelTag: "slow", CallerSkip: 1}, a...)
}

func (l *CustomLogger) Stat(a ...any) {
    l.Log(LogOption{Level: INFO, LevelTag: "stat", CallerSkip: 1}, a...)
}

func (l *CustomLogger) Debug(a ...any) {
    l.Log(LogOption{Level: DEBUG, CallerSkip: 1, EnableCaller: EnableOpen}, a...)
}
```

### 日志内容输出
目前日志内容默认输出到控制台。
如果需要输出内容到文件中，需要设置日志的Writer,可以通过将文件指针传递给NewWriter函数来构造一个Writer。
如果想要实现更强大的输出控制如切割日志文件,可以使用[github.com/lestrrat-go/file-rotatelogs](https://github.com/lestrrat-go/file-rotatelogs)来构造Writer。
自主实现Write方法时需要注意参数[]byte不应该超出该方法的作用域，否则可能会导致数据并发问题并导致混乱。

### 性能
记录一条消息和3个字段(禁用caller输出,并发测试)：
```
goos: darwin
goarch: arm64
pkg: github.com/welllog/olog
BenchmarkInfow
BenchmarkInfow/std.logger
BenchmarkInfow/std.logger-10         	 3760806	       287.4 ns/op	      48 B/op	       1 allocs/op
BenchmarkInfow/olog.json
BenchmarkInfow/olog.json-10          	14397530	        77.93 ns/op	      96 B/op	       1 allocs/op
BenchmarkInfow/olog.plain
BenchmarkInfow/olog.plain-10         	21455214	        56.19 ns/op	      96 B/op	       1 allocs/op
BenchmarkInfow/olog.ctx.json
BenchmarkInfow/olog.ctx.json-10      	14290038	        83.48 ns/op	      96 B/op	       1 allocs/op
BenchmarkInfow/olog.ctx.plain
BenchmarkInfow/olog.ctx.plain-10     	19133161	        59.64 ns/op	      96 B/op	       1 allocs/op
PASS
```