<p align="center">
    <br> English | <a href="README-CN.md">中文</a>
</p>

# olog
* olog is a lightweight, high-performance, out-of-the-box logging library that relies solely on the Go standard library.
* Supports outputting logs in both JSON and plain text formats.
* Supports setting context handling functions to retrieve fields from the context for output.
* Supports seven log levels: TRACE, DEBUG, INFO, NOTICE, WARN, ERROR, and FATAL, corresponding to the "trace", "debug", "info", "notice", "warn", "error", and "fatal" tags. Users can also define their own semantic tags, such as "slow" and "stat".
* Provides output switch control for all log levels except FATAL.
* Enhances control over caller printing. Users can disable caller printing in global settings to improve performance, but can enable support for caller printing when printing certain critical logs.
* Provides a Logger interface for users to easily construct their own logger.

### code example
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
The json output is as follows:
```json
{"@timestamp":"2023-04-20T16:13:00+08:00","level":"trace","caller":"olog/log_test.go:377","content":"hello world","stack":"\nolog/log_test.go:377&github.com/welllog/olog.TestPlainOutput\ntesting/testing.go:1576&testing.tRunner\nruntime/asm_amd64.s:1598&runtime.goexit"}
{"@timestamp":"2023-04-20T16:13:00+08:00","level":"trace","caller":"olog/log_test.go:378","content":"hello","name":"bob","stack":"\nolog/log_test.go:378&github.com/welllog/olog.TestPlainOutput\ntesting/testing.go:1576&testing.tRunner\nruntime/asm_amd64.s:1598&runtime.goexit"}
{"@timestamp":"2023-04-20T16:13:00+08:00","level":"debug","caller":"olog/log_test.go:379","content":"hello world"}
{"@timestamp":"2023-04-20T16:13:00+08:00","level":"info","caller":"olog/log_test.go:380","content":"hello","name":"linda","age":18}
{"@timestamp":"2023-04-20T16:13:00+08:00","level":"notice","caller":"olog/log_test.go:381","content":"hello world"}
{"@timestamp":"2023-04-20T16:13:00+08:00","level":"warn","caller":"olog/log_test.go:382","content":"hello world"}
{"@timestamp":"2023-04-20T16:13:00+08:00","level":"warn","caller":"olog/log_test.go:383","content":"hello","order_no":"AWESDDF"}
{"@timestamp":"2023-04-20T16:13:00+08:00","level":"error","caller":"olog/log_test.go:384","content":"hello world","success":true}
{"@timestamp":"2023-04-20T16:13:00+08:00","level":"print","content":"hello world","price":32.5,"stack":"\nolog/log_test.go:385&github.com/welllog/olog.TestPlainOutput"}
{"@timestamp":"2023-04-20T16:13:00+08:00","level":"fatal","caller":"olog/log_test.go:388","content":"fatal exit"}
```
The plain output is as follows:
```
2023-04-20T16:09:26+08:00	trace	olog/log_test.go:377	hello world	stack=
olog/log_test.go:377&github.com/welllog/olog.TestPlainOutput
testing/testing.go:1576&testing.tRunner
runtime/asm_amd64.s:1598&runtime.goexit
2023-04-20T16:09:26+08:00	trace	olog/log_test.go:378	hello	name=bob	stack=
olog/log_test.go:378&github.com/welllog/olog.TestPlainOutput
testing/testing.go:1576&testing.tRunner
runtime/asm_amd64.s:1598&runtime.goexit
2023-04-20T16:09:26+08:00	debug	olog/log_test.go:379	hello world
2023-04-20T16:09:26+08:00	info	olog/log_test.go:380	hello	name=linda	age=18
2023-04-20T16:09:26+08:00	notice	olog/log_test.go:381	hello world
2023-04-20T16:09:26+08:00	warn	olog/log_test.go:382	hello world
2023-04-20T16:09:26+08:00	warn	olog/log_test.go:383	hello	order_no=AWESDDF
2023-04-20T16:09:26+08:00	error	olog/log_test.go:384	hello world	success=true
2023-04-20T16:09:26+08:00	print	hello world	price=32.5	stack=
olog/log_test.go:385&github.com/welllog/olog.TestPlainOutput
2023-04-20T16:09:26+08:00	fatal	olog/log_test.go:388	fatal exit
```

### contextLogger uses
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

### Implement your own logger
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

### Log Content Output
Currently, log content is output to the console by default. 
To output content to a file, you need to set the log's Writer by constructing a Writer with the NewWriter function and passing a file pointer.
If you want to implement more powerful output control, such as log file splitting, you can use [github.com/lestrrat-go/file-rotatelogs](https://github.com/lestrrat-go/file-rotatelogs) to construct a Writer.
When implementing the Write method on your own, it is important to note that the []byte parameter should not exceed the scope of the method, otherwise data concurrency issues may occur and result in confusion.

### Performance
Log a message and 3 fields(disable caller output and RunParallel):
```
goos: darwin
goarch: arm64
pkg: github.com/welllog/olog
BenchmarkInfo
BenchmarkInfo/std.logger
BenchmarkInfo/std.logger-10         	 3740652	       284.4 ns/op	      48 B/op	       1 allocs/op
BenchmarkInfo/olog.json
BenchmarkInfo/olog.json-10          	14560014	        83.89 ns/op	      96 B/op	       1 allocs/op
BenchmarkInfo/olog.plain
BenchmarkInfo/olog.plain-10         	20607255	        58.84 ns/op	      96 B/op	       1 allocs/op
BenchmarkInfo/olog.ctx.json
BenchmarkInfo/olog.ctx.json-10      	 7297224	       150.4 ns/op	     376 B/op	       6 allocs/op
BenchmarkInfo/olog.ctx.plain
BenchmarkInfo/olog.ctx.plain-10     	 8775705	       134.7 ns/op	     376 B/op	       6 allocs/op
PASS
```