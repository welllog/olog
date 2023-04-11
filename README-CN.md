<p align="center">
    <br> <a href="README.md">English</a> | 中文
</p>

# olog
* olog是一个轻量级、高性能、开箱即用的日志库，完全依赖于Go标准库。
* 支持以JSON和纯文本格式输出日志。
* 支持设置上下文处理函数，以从上下文中检索字段以进行输出。
* 支持五个日志级别：DEBUG、INFO、WARN、ERROR和FATAL，对应于“debug”、“info”、“warn”、“error”和“fatal”标签。用户还可以定义自己的语义标签，如“slow”和“stat”。
* 提供输出开关控制，除了FATAL之外的所有日志级别都可以控制输出。
* 增强了对调用者打印的控制。用户可以在全局设置中禁用调用者打印以提高性能，但可以在打印某些关键日志时启用调用者打印支持。
* 为用户提供了一个Logger接口，以便轻松构建自己的日志记录器。

### 代码示例
```go
        Debug("hello world")
	Debugw("hello", Field{Key: "name", Value: "bob"})
	Info("hello world")
	Infow("hello", Field{Key: "name", Value: "linda"}, Field{Key: "age", Value: 18})
	Warnf("hello %s", "world")
	Warnw("hello", Field{Key: "order_no", Value: "AWESDDF"})
	Error("hello world")
	Errorw("hello world", Field{Key: "success", Value: true})
	Log(DEBUG, WithTag("trace"), WithPrintMsg("hello world"), WithCaller(false),
		WithFields(Field{Key: "price", Value: 32.5}))
	Fatal("fatal exit")
```
json输出如下：
```json
{"@timestamp":"2023-04-11T16:54:19+08:00","level":"debug","caller":"olog/log_test.go:501","content":"hello world"}
{"@timestamp":"2023-04-11T16:54:19+08:00","level":"debug","caller":"olog/log_test.go:502","content":"hello","name":"bob"}
{"@timestamp":"2023-04-11T16:54:19+08:00","level":"info","caller":"olog/log_test.go:503","content":"hello world"}
{"@timestamp":"2023-04-11T16:54:19+08:00","level":"info","caller":"olog/log_test.go:504","content":"hello","name":"linda","age":18}
{"@timestamp":"2023-04-11T16:54:19+08:00","level":"warn","caller":"olog/log_test.go:505","content":"hello world"}
{"@timestamp":"2023-04-11T16:54:19+08:00","level":"warn","caller":"olog/log_test.go:506","content":"hello","order_no":"AWESDDF"}
{"@timestamp":"2023-04-11T16:54:19+08:00","level":"error","caller":"olog/log_test.go:507","content":"hello world"}
{"@timestamp":"2023-04-11T16:54:19+08:00","level":"error","caller":"olog/log_test.go:508","content":"hello world","success":true}
{"@timestamp":"2023-04-11T16:54:19+08:00","level":"trace","content":"hello world","price":32.5}
{"@timestamp":"2023-04-11T16:54:19+08:00","level":"fatal","caller":"olog/log_test.go:511","content":"fatal exit"}
```
plain输出如下：
![plain](plain.webp)

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
	logger.Debug("test 2")
```

### 实现自己的logger
```
type CustomLogger struct {
	Logger
}

func (l *CustomLogger) Slow(a ...any) {
	l.Log(WARN, WithPrint(a...), WithTag("slow"), WithCallerSkipOne)
}

func (l *CustomLogger) Stat(a ...any) {
	Log(INFO, WithPrint(a...), WithTag("stat"), WithCallerSkipOne)
}

func (l *CustomLogger) Debug(a ...any) {
	Log(DEBUG, WithPrint(a...), WithCallerSkipOne, WithCaller(true))
}
```

### 日志内容输出
目前日志内容默认输出到控制台。
如果需要输出内容到文件中，需要设置日志的Writer,可以通过将文件指针传递给NewWriter函数来构造一个Writer。
如果想要实现更强大的输出控制如切割日志文件,可以使用[github.com/lestrrat-go/file-rotatelogs](https://github.com/lestrrat-go/file-rotatelogs)来构造Writer。
自主实现Write方法时需要注意参数[]byte不应该超出该方法的作用域，否则可能会导致数据并发问题并导致混乱。

### 性能
记录一条消息和3个字段(禁用caller输出,并发测试)：
![bench](bench.webp)
