package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	formatter "github.com/weixiaolv/nested-logrus-formatter"
)

func main() {
	fmt.Print("\n--- nested-logrus-formatter ---\n\n")
	printDemo(&formatter.Formatter{
		HideKeys:        false,
		FieldsOrder:     []string{"component", "category", "req"},
		ModuleName:      "mod",
		TimestampFormat: "2006-01-02 15:04:05.000",
		NoFieldsSpace:   true,
		// NoFieldsColors:  true,
	}, "nested-logrus-formatter")

	fmt.Print("\n--- default logrus formatter ---\n\n")
	printDemo(nil, "default logrus formatter")
}

func printDemo(f logrus.Formatter, title string) {
	l := logrus.New()

	l.SetLevel(logrus.DebugLevel)

	if f != nil {
		l.SetFormatter(f)
	}

	// enable/disable file/function name
	l.SetReportCaller(false)

	l.Infof("this is %v demo", title)

	lWebServer := l.WithField("mod", "web-server")
	lWebServer.Info("starting...")
	lWebServerReq := lWebServer.WithFields(logrus.Fields{
		"req":   "GET /api/stats",
		"reqId": "#1",
	})
	lWebServerReq.Info("params: startYear=2048")
	lWebServerReq.Error("response: 400 Bad Request")

	lDbConnector := l.WithField("mod", "db-connector")
	lDbConnector.Info("connecting to db on 10.10.10.13...")
	lDbConnector.Warn("connection took 10s")

	l.Info("demo end.")
}
