#当前目录
CURRENT_DIR := $(shell pwd)

bin:
	go build -o ./bin/gsgen_tools ./gsgen_tools/main.go
model:bin
	./bin/gsgen_tools -d="$(CURRENT_DIR)/example/getter" -f=".model.go"
	./bin/gsgen_tools -d="$(CURRENT_DIR)/example/setter" -f=".model.go" -s
	./bin/gsgen_tools -d="$(CURRENT_DIR)/example/bson" -f=".model.go" -s -b
	./bin/gsgen_tools -d="$(CURRENT_DIR)/example/nest" -f=".model.go" -s -b -a="// test head annotations 1" -a="// test head annotations 2"


.PHONY: bin model