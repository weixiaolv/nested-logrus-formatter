# nested-logrus-formatter

## Project

This formatter fork from [logrus-logstash-hook](https://github.com/bshuster-repo/logrus-logstash-hook), and add some features.

1. Add `ModuleName` Formatter define, if logrus `fields` has `ModuleName` value, logger will add `[moduleValue]` to the log after level info.
2. The module key & value will not exists fields list in log.

Example:

```log
2024-07-28T15:41:01.385631 [INFO] [db] [db:localhost:3389/customer][id:999] update customer
2024-07-28T15:41:01.445527 [INFO] [http] [uri:POST /customer][status:200][size:3879] modify customer mobile success
```

## Configuration

```go
type Formatter struct {
  // FieldsOrder - default: fields sorted alphabetically
  FieldsOrder []string

  // TimestampFormat - default: time.StampMilli = "Jan _2 15:04:05.000"
  TimestampFormat string

  // HideKeys - show [fieldValue] instead of [fieldKey:fieldValue]
  HideKeys bool

  // NoColors - disable colors
  NoColors bool

  // NoFieldsColors - apply colors only to the level, default is level + fields
  NoFieldsColors bool

  // NoFieldsSpace - no space between fields
  NoFieldsSpace bool

  // ShowFullLevel - show a full level [WARNING] instead of [WARN]
  ShowFullLevel bool

  // NoUppercaseLevel - no upper case for level value
  NoUppercaseLevel bool

  // TrimMessages - trim whitespaces on messages
  TrimMessages bool

  // CallerFirst - print caller info first
  CallerFirst bool

  // CustomCallerFormatter - set custom formatter for caller info
  CustomCallerFormatter func(*runtime.Frame) string

  // Module Name
  ModuleName string
}
```

## Usage

Get the Package:

```sh
go get github.com/weixiaolv/nested-logrus-formatter
```

Example with module:

```go
import (
  formatter "github.com/weixiaolv/nested-logrus-formatter"
  "github.com/sirupsen/logrus"
)

func main() {
  logger := logrus.New()
  logger.SetFormatter(&formatter.Formatter{
    HideKeys:    false,
    FieldsOrder: []string{"api", "status"},
    ModuleName:  "module",
  })
  log := logger.WithField("module", "main")
  
  // no module
  logger.Info("just info message")
  // Output: Jan _2 15:04:05.000 [INFO] just info message
  
  // with module
  log.WithField("api", "rest").Warn("warn message")
  // Output: Jan _2 15:04:05.000 [WARN] [main] [api:rest] warn message
}
```

We can set debug mode print caller info, product mode do not print it.

```go
var logger *logrus.Logger

func NewLogger(levelConfig string) *logrus.Logger {
  level, err := logrus.ParseLevel(levelConfig)
  if err != nil {
    fmt.Println("err=", err)
    level = logrus.InfoLevel
  }

  logger := logrus.New()
  logger.SetLevel(level)

  formatter := &formatter.Formatter{
    FieldsOrder:     []string{"api", "status"},
    ModuleName:      "mod",
    TimestampFormat: "2006-01-02 15:04:05.000",
    NoFieldsSpace:   true,
  }
  if level >= logrus.DebugLevel {
    logger.SetReportCaller(true)
    formatter.CallerFirst = true
    formatter.CustomCallerFormatter = func(f *runtime.Frame) string {
      fileStr := path.Base(f.File)
      arr := strings.Split(f.Function, ".")
      fnStr := arr[len(arr)-1]
      return fmt.Sprintf("[%s::%s()::%d]", fileStr, fnStr, f.Line)
    }
  }
  logger.SetFormatter(formatter)
  return logger
}

func main() {
  // create
  levelConfig := "debug"
  logger = NewLogger(levelConfig)

  // local module
  log := logger.WithField("mod", "test")
  
  // print log
  log.Info("log level ", levelConfig)

  // Output: while levelConfig = "debug"
  // 2024-08-07 00:05:10.378 [main.go::main()::55] [INFO] [test] log level debug

  // Output: while levelConfig = "info"
  // 2024-08-07 00:05:10.378 [INFO] [test] log level info
}
```

See more examples in the [tests](./tests/formatter_test.go) file.

## Development

run demo

```bash
# run demo:
make demo
```

auto test & coverage report

```sh
# run tests:
make test

# run coverage:
make cover
```

test case coverage 98% source code

## Benchmark

run:

```sh
make bench
```

result:

```log
BenchmarkDefaultFormat-4          507812   2731 ns/op  716 B/op  14 allocs/op
BenchmarkNestedLogrusFormatter-4  793026   1939 ns/op  644 B/op  11 allocs/op
BenchmarkToFile-4                 592366   1914 ns/op  644 B/op  11 allocs/op
BenchmarkNoFields-4              1000000   1360 ns/op  396 B/op   8 allocs/op
```

1. NestedLogrusFormatter a little faster than logrus DefaultFormat
2. Benchmark default output is io.Discard, log out to file with bufio.Writer, the performance is not bad.
3. Log out no fields means no interface{} to string convert, reduce allocs/op, 30% faster.

more info cpu profile & mem profile

```sh
make cpuprof
make memprof
```
