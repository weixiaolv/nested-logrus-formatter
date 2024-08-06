package tests

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	formatter "github.com/weixiaolv/nested-logrus-formatter"
)

func getLogger() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
	return log
}

func ExampleFormatter_Format_default() {
	log := getLogger()
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
	})
	log.Debug("test1")
	log.Info("test2")
	log.Warn("test3")
	log.Error("test4")

	// Output:
	// - [DEBU] test1
	// - [INFO] test2
	// - [WARN] test3
	// - [ERRO] test4
}

func TestFormatterFormatTimeStamp(t *testing.T) {
	log := getLogger()
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	output := bytes.NewBuffer([]byte{})
	log.SetOutput(output)

	log.Info("test1")
	expectedRegExp := `\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3} \[INFO\] test1`
	matchRegexp(t, expectedRegExp, output)
}

// ---- level ----

func ExampleFormatter_Format_full_level() {
	log := getLogger()
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		ShowFullLevel:   true,
	})

	log.Debug("test1")
	log.Info("test2")
	log.Warn("test3")
	log.Error("test4")

	// Output:
	// - [DEBUG] test1
	// - [INFO] test2
	// - [WARNING] test3
	// - [ERROR] test4
}

func ExampleFormatter_Format_no_uppercase_level() {
	log := getLogger()
	log.SetFormatter(&formatter.Formatter{
		NoColors:         true,
		TimestampFormat:  "-",
		NoUppercaseLevel: true,
	})

	log.Debug("test1")
	log.Info("test2")
	log.Warn("test3")
	log.Error("test4")

	// Output:
	// - [debu] test1
	// - [info] test2
	// - [warn] test3
	// - [erro] test4
}

// ---- module ----

func ExampleFormatter_Format_module() {
	log := getLogger()

	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		ModuleName:      "mod",
	})
	l := log.WithField("mod", "test_mod")

	l.Info("test1")

	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		ModuleName:      "mod",
		FieldsOrder:     []string{"test"},
	})
	l.Info("test2")

	// Output:
	// - [INFO] [test_mod] test1
	// - [INFO] [test_mod] test2

}

func ExampleFormatter_Format_module_with_fields() {
	log := getLogger()
	l := log.WithField("mod", "test_mod").WithFields(logrus.Fields{
		"c": "c2",
		"a": 1,
		"b": "b1",
	})

	// mod
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		ModuleName:      "mod",
	})
	l.Info("test1")

	// module with nospace
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		ModuleName:      "mod",
		NoFieldsSpace:   true,
	})
	l.Info("test2")

	// module with hidekeys
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		ModuleName:      "mod",
		HideKeys:        true,
	})
	l.Info("test3")

	// module with order fields
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		ModuleName:      "mod",
		FieldsOrder:     []string{"b"},
	})
	l.Info("test4")

	// module with order fields and nospace
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		ModuleName:      "mod",
		FieldsOrder:     []string{"b"},
		NoFieldsSpace:   true,
	})
	l.Info("test5")
	// Output:
	// - [INFO] [test_mod] [a:1] [b:b1] [c:c2] test1
	// - [INFO] [test_mod] [a:1][b:b1][c:c2] test2
	// - [INFO] [test_mod] [1] [b1] [c2] test3
	// - [INFO] [test_mod] [b:b1] [a:1] [c:c2] test4
	// - [INFO] [test_mod] [b:b1][a:1][c:c2] test5
}

// ---- fields ----

func ExampleFormatter_Format_fields() {
	log := getLogger()

	// no field
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		HideKeys:        false,
	})
	log.Info("test1")

	// one field
	l := log.WithField("z", 1)
	l.Info("test2")
	// two field
	ll := l.WithField("b", "b1")
	ll.Info("test3")

	// hide fields keys
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		HideKeys:        true,
	})
	ll.Info("test4")

	// Output:
	// - [INFO] test1
	// - [INFO] [z:1] test2
	// - [INFO] [b:b1] [z:1] test3
	// - [INFO] [b1] [1] test4
}

func ExampleFormatter_Format_Fields_with_sortorder() {
	log := getLogger()
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		HideKeys:        false,
	})

	// default order
	l := log.WithFields(logrus.Fields{
		"d": true,
		"b": "main",
		"c": 100,
		"a": "rest",
	})
	l.Info("test1")

	// set order
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		HideKeys:        false,
		FieldsOrder:     []string{"d", "a"},
	})
	l.Info("test2")

	// set order & hide keys
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		HideKeys:        true,
		FieldsOrder:     []string{"d", "a"},
	})
	l.Info("test3")

	// set order & nospace
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		NoFieldsSpace:   true,
		FieldsOrder:     []string{"d", "a"},
	})
	l.Info("test4")

	// Output:
	// - [INFO] [a:rest] [b:main] [c:100] [d:true] test1
	// - [INFO] [d:true] [a:rest] [b:main] [c:100] test2
	// - [INFO] [true] [rest] [main] [100] test3
	// - [INFO] [d:true][a:rest][b:main][c:100] test4
}

// ---- color ----

func TestFormatterColor(t *testing.T) {
	log := getLogger()
	output := bytes.NewBuffer([]byte{})
	log.SetOutput(output)

	l := log.WithFields(logrus.Fields{
		"mod": "test_mod",
		"b":   "world",
		"a":   "hello",
	})

	// only level color
	log.SetFormatter(&formatter.Formatter{
		TimestampFormat: "-",
		HideKeys:        false,
		ModuleName:      "mod",
		NoColors:        false,
		NoFieldsColors:  true,
	})
	l.Info("test1")
	expectedRegExp := `- \x1b\[36m\[INFO\]\x1b\[0m \[test_mod\] \[a:hello\] \[b:world\] test1`
	matchRegexp(t, expectedRegExp, output)

	// level & fields color
	log.SetFormatter(&formatter.Formatter{
		TimestampFormat: "-",
		HideKeys:        false,
		ModuleName:      "mod",
		NoColors:        false,
		NoFieldsColors:  false,
	})
	l.Info("test1")
	expectedRegExp = `- \x1b\[36m\[INFO\]\x1b\[0m \x1b\[36m\[test_mod\]\x1b\[0m \x1b\[36m\[a:hello\] \[b:world\]\x1b\[0m test1`
	matchRegexp(t, expectedRegExp, output)

	// level & fields color with order
	log.SetFormatter(&formatter.Formatter{
		TimestampFormat: "-",
		HideKeys:        false,
		ModuleName:      "mod",
		NoColors:        false,
		NoFieldsColors:  false,
		FieldsOrder:     []string{"b"},
	})
	l.Info("test1")
	expectedRegExp = `- \x1b\[36m\[INFO\]\x1b\[0m \x1b\[36m\[test_mod\]\x1b\[0m \x1b\[36m\[b:world\] \[a:hello\]\x1b\[0m test1`
	matchRegexp(t, expectedRegExp, output)
}

// ---- message ----

func ExampleFormatter_Format_message() {
	log := getLogger()
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
	})
	log.Info("test1")
	log.Info(" test2")
	log.Info("	test3")

	// trim
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		TrimMessages:    true,
	})
	log.Info(" test4 ")
	log.Info("	test5	")

	// Output:
	// - [INFO] test1
	// - [INFO]  test2
	// - [INFO] 	test3
	// - [INFO] test4
	// - [INFO] test5
}

// ---- caller ----

func TestFormatterCaller(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	log := getLogger()
	log.SetReportCaller(true)
	log.SetOutput(output)

	// caller behind the logger
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
	})
	log.Info("test1")
	expectedRegExp := `- \[INFO\] test1 \(.+\.go:[0-9]+ github.com/weixiaolv/nested-logrus-formatter.+\)\n$`
	matchRegexp(t, expectedRegExp, output)

	// caller front the logger
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		CallerFirst:     true,
	})
	log.Info("test1")
	expectedRegExp = `- \(.+\.go:[0-9]+ github.com/weixiaolv/nested-logrus-formatter.+\) \[INFO\] test1\n$`
	matchRegexp(t, expectedRegExp, output)
}

func TestFormatterCustomCaller(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	log := getLogger()
	log.SetReportCaller(true)
	log.SetOutput(output)

	// user define caller func
	log.SetFormatter(&formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		CallerFirst:     true,
		CustomCallerFormatter: func(f *runtime.Frame) string {
			fileStr := path.Base(f.File)
			arr := strings.Split(f.Function, ".")
			fnStr := arr[len(arr)-1]
			return fmt.Sprintf("[%s::%s::%d]", fileStr, fnStr, f.Line)
		},
	})
	log.Info("test1")
	expectedRegExp := `- \[formatter_test.go::TestFormatterCustomCaller::\d+\] \[INFO\] test1\n$`
	matchRegexp(t, expectedRegExp, output)
}

func matchRegexp(t *testing.T, regStr string, output *bytes.Buffer) {
	line, err := output.ReadString('\n')
	if err != nil {
		t.Errorf("Cannot read log output: %v", err)
	}
	match, err := regexp.MatchString(regStr, line)
	if err != nil {
		t.Errorf("Cannot check regexp: %v", err)
	} else if !match {
		t.Errorf("Output doesn't match\n  expected: %s\n  got: '%s'", regStr, line)
	}
}
