package tests

import (
	"io"
	"testing"

	formatter "github.com/weixiaolv/nested-logrus-formatter"
)

func BenchmarkDefaultFormat(b *testing.B) {
	log := getLogger()
	log.SetOutput(io.Discard)

	l := log.WithField("api", "/api/hello")

	for i := 0; i < b.N; i++ {
		l.Info("test")
	}
}

func BenchmarkNestedLogrusFormatter(b *testing.B) {
	log := getLogger()
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		ModuleName:      "module",
	})
	log.SetOutput(io.Discard)

	l := log.WithField("module", "mymodule").WithField("api", "/api/hello")

	for i := 0; i < b.N; i++ {
		l.Info("test")
	}
}

func BenchmarkNoFields(b *testing.B) {
	log := getLogger()
	log.SetFormatter(&formatter.Formatter{
		NoColors:         true,
		TimestampFormat:  "-",
		ModuleName:       "module",
		NoUppercaseLevel: true,
	})
	log.SetOutput(io.Discard)

	l := log

	for i := 0; i < b.N; i++ {
		l.Info("test")
	}
}
