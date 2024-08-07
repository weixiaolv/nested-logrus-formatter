package tests

import (
	"bufio"
	"io"
	"os"
	"testing"

	formatter "github.com/weixiaolv/nested-logrus-formatter"
)

func BenchmarkDefaultFormat(b *testing.B) {
	log := getLogger()
	log.SetOutput(io.Discard)

	l := log.WithField("module", "mymodule").WithField("api", "/api/hello")

	for i := 0; i < b.N; i++ {
		l.Info("test")
	}
}

func BenchmarkNestedLogrusFormatter(b *testing.B) {
	log := getLogger()
	log.SetOutput(io.Discard)
	log.SetFormatter(&formatter.Formatter{
		NoColors:   true,
		ModuleName: "module",
	})

	l := log.WithField("module", "mymodule").WithField("api", "/api/hello")

	for i := 0; i < b.N; i++ {
		l.Info("test")
	}
}

func BenchmarkToFile(b *testing.B) {
	log := getLogger()
	log.SetFormatter(&formatter.Formatter{
		NoColors:   true,
		ModuleName: "module",
	})

	file, err := os.Create("log.out")
	if err != nil {
		b.Fatal(err)
	}
	defer file.Close()
	buf := bufio.NewWriter(file)
	log.SetOutput(buf)

	l := log.WithField("module", "mymodule").WithField("api", "/api/hello")

	for i := 0; i < b.N; i++ {
		l.Info("test")
	}
	buf.Flush()
}

func BenchmarkNoFields(b *testing.B) {
	log := getLogger()
	log.SetFormatter(&formatter.Formatter{
		NoColors:         true,
		ModuleName:       "module",
		NoUppercaseLevel: true,
	})
	log.SetOutput(io.Discard)

	l := log

	for i := 0; i < b.N; i++ {
		l.Info(" [mymodule] [api:/api/hello] test")
	}
}
