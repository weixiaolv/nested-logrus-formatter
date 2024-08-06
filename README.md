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

benchmark

```sh
# new prof
make bench

# show prof
make cpuprof
make memprof
```

benchmark result

```log
enchmarkDefaultFormat-4           550015   2121 ns/op  684 B/op  14 allocs/op
BenchmarkNestedLogrusFormatter-4  685206   2140 ns/op  620 B/op  10 allocs/op
BenchmarkNoFields-4              1616091  726.0 ns/op  328 B/op   7 allocs/op
```

NestedLogrusFormatter a little faster than default format. NoFields means no interface{} to string convert, reduce allocs/op.
