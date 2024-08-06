# nested-logrus-formatter

.PHONY: all
all: test demo

.PHONY: test
test:
	go test ./tests -v -count=1

cover:
	go test ./tests -coverpkg=./ -v -covermode=count -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

.PHONY: bench
bench:	
	go test ./tests -v -bench=^Benchmark -run=none -benchmem -memprofile=mem.out -cpuprofile=cpu.out

cpuprof:
	go tool pprof -http=:8080 ./cpu.out

memprof:
	go tool pprof -http=:8080 mem.out

.PHONY: demo
demo:
	go run example/main.go
