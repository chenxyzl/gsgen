#当前目录
CURRENT_DIR := $(shell pwd)

bin:
	go build -o ./bin/gsgen_tools ./gsgen_tools/main.go
model:bin
	./bin/gsgen_tools -d="$(CURRENT_DIR)/example/getter" -f=".model.go"
	./bin/gsgen_tools -d="$(CURRENT_DIR)/example/setter" -f=".model.go" -s
	./bin/gsgen_tools -d="$(CURRENT_DIR)/example/bson" -f=".model.go" -s -b
	./bin/gsgen_tools -d="$(CURRENT_DIR)/example/nest" -f=".model.go" -s -b -a="// test head annotations 1" -a="// test head annotations 2"
	./bin/gsgen_tools -d="$(CURRENT_DIR)/example/common" -f=".model.go" -b
	./bin/gsgen_tools -d="$(CURRENT_DIR)/example/with_ignore" -f=".model.go" -s -b -i="github.com/chenxyzl/gsgen/example/common.Common"


.PHONY: bin model